// ABOUTME: Tests for HTTP handlers, file serving, and live reload functionality.
// ABOUTME: Covers fileHandler, SSE connections, and client notifications.
package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileHandler_ServeHTTP(t *testing.T) {
	tempDir := createTempTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	// Create test files
	htmlContent := "<html><body><h1>Test Page</h1></body></html>"
	indexContent := "<html><body><h1>Index Page</h1></body></html>"
	notFoundContent := "<html><body><h1>404 Not Found</h1></body></html>"
	cssContent := "body { color: red; }"

	writeTestFile(t, filepath.Join(tempDir, "test.html"), htmlContent)
	writeTestFile(t, filepath.Join(tempDir, "index.html"), indexContent)
	writeTestFile(t, filepath.Join(tempDir, "404.html"), notFoundContent)
	writeTestFile(t, filepath.Join(tempDir, "style.css"), cssContent)

	// Create directory with index.html
	dirPath := filepath.Join(tempDir, "subdir")
	err := os.MkdirAll(dirPath, 0755)
	require.NoError(t, err)
	writeTestFile(t, filepath.Join(dirPath, "index.html"), indexContent)

	tests := []struct {
		name            string
		path            string
		dev             bool
		expectedStatus  int
		expectedBody    string
		expectedType    string
		checkLiveReload bool
	}{
		{
			name:            "serves HTML file in dev mode with live reload",
			path:            "/test.html",
			dev:             true,
			expectedStatus:  http.StatusOK,
			expectedBody:    htmlContent,
			expectedType:    "text/html; charset=utf-8",
			checkLiveReload: true,
		},
		{
			name:            "serves HTML file in production mode without live reload",
			path:            "/test.html",
			dev:             false,
			expectedStatus:  http.StatusOK,
			expectedBody:    htmlContent,
			expectedType:    "text/html",
			checkLiveReload: false,
		},
		{
			name:            "serves CSS file directly",
			path:            "/style.css",
			dev:             true,
			expectedStatus:  http.StatusOK,
			expectedBody:    cssContent,
			expectedType:    "text/css",
			checkLiveReload: false,
		},
		{
			name:            "serves index.html for directory request",
			path:            "/subdir/",
			dev:             true,
			expectedStatus:  http.StatusOK,
			expectedBody:    indexContent,
			expectedType:    "text/html; charset=utf-8",
			checkLiveReload: true,
		},
		{
			name:            "serves index.html for root directory",
			path:            "/",
			dev:             true,
			expectedStatus:  http.StatusOK,
			expectedBody:    indexContent,
			expectedType:    "text/html; charset=utf-8",
			checkLiveReload: true,
		},
		{
			name:            "serves 404 for non-existent file",
			path:            "/nonexistent.html",
			dev:             true,
			expectedStatus:  http.StatusOK, // Custom 404 handler returns 200
			expectedBody:    notFoundContent,
			expectedType:    "text/html; charset=utf-8",
			checkLiveReload: true,
		},
		{
			name:            "serves 404 for directory without index",
			path:            "/empty/",
			dev:             true,
			expectedStatus:  http.StatusOK, // Custom 404 handler returns 200
			expectedBody:    notFoundContent,
			expectedType:    "text/html; charset=utf-8",
			checkLiveReload: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &fileHandler{
				root:     http.Dir(tempDir),
				notFound: "404.html",
				dev:      tt.dev,
			}

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			if tt.expectedType != "" {
				assert.Contains(t, rec.Header().Get("Content-Type"), strings.Split(tt.expectedType, ";")[0])
			}

			if tt.checkLiveReload {
				assert.Contains(t, rec.Body.String(), "EventSource('/livereload')")
				assert.Contains(t, rec.Body.String(), "location.reload()")
			} else if tt.dev && strings.HasSuffix(tt.path, ".html") {
				// In dev mode, HTML should have live reload unless explicitly checked
				if !tt.checkLiveReload {
					assert.NotContains(t, rec.Body.String(), "EventSource('/livereload')")
				}
			}

			if tt.expectedBody != "" && !tt.checkLiveReload {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestFileHandler_serve404_MissingFile(t *testing.T) {
	// Test the case where the 404.html file itself is missing
	// This should trigger the fallback to http.NotFound
	tempDir := t.TempDir()

	// Create handler pointing to non-existent 404 file
	handler := &fileHandler{
		root:     http.Dir(tempDir),
		notFound: "missing-404.html", // This file doesn't exist
		dev:      true,
	}

	req := httptest.NewRequest(http.MethodGet, "/nonexistent.html", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Should fall back to standard 404 when custom 404 file is missing
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "404 page not found")
}

func TestFileHandler_serveHTML(t *testing.T) {
	tests := []struct {
		name         string
		content      string
		expectedBody string
	}{
		{
			name:         "injects live reload script",
			content:      "<html><body><h1>Test</h1></body></html>",
			expectedBody: "EventSource('/livereload')",
		},
		{
			name:         "handles content without body tag",
			content:      "<html><div>No body tag</div></html>",
			expectedBody: "<html><div>No body tag</div></html>", // Should remain unchanged
		},
		{
			name:         "handles multiple body tags",
			content:      "<html><body>First</body><body>Second</body></html>",
			expectedBody: "EventSource('/livereload')",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := &fileHandler{dev: true}

			// Create a mock file
			file := &mockFile{content: []byte(tt.content)}
			rec := httptest.NewRecorder()

			handler.serveHTML(file, rec)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "text/html; charset=utf-8", rec.Header().Get("Content-Type"))
			assert.Contains(t, rec.Body.String(), tt.expectedBody)
		})
	}
}

func TestDevServer_handleLiveReload(t *testing.T) {
	site := createTestSite(t)
	gen := createTestGenerator(t, site)
	parser := createTestParser(t, site)

	server, err := NewDevServer("", "", site, gen, parser)
	require.NoError(t, err)
	defer server.watcher.Close()

	tests := []struct {
		name        string
		timeout     time.Duration
		expectError bool
	}{
		{
			name:        "establishes SSE connection",
			timeout:     100 * time.Millisecond,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/livereload", nil)

			// Create a context with timeout
			ctx, cancel := context.WithTimeout(req.Context(), tt.timeout)
			defer cancel()
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()

			// Run handleLiveReload in a goroutine since it blocks
			done := make(chan bool)
			go func() {
				server.handleLiveReload(rec, req)
				done <- true
			}()

			// Wait for timeout or completion
			select {
			case <-done:
				// Handler completed due to context cancellation
			case <-time.After(tt.timeout + 50*time.Millisecond):
				t.Error("Handler did not respect context timeout")
			}

			// Verify SSE headers
			assert.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
			assert.Equal(t, "no-cache", rec.Header().Get("Cache-Control"))
			assert.Equal(t, "keep-alive", rec.Header().Get("Connection"))
		})
	}
}

func TestDevServer_handleLiveReload_NonFlusher(t *testing.T) {
	// Test the case where ResponseWriter doesn't implement http.Flusher
	site := createTestSite(t)
	gen := createTestGenerator(t, site)
	parser := createTestParser(t, site)

	server, err := NewDevServer("", "", site, gen, parser)
	require.NoError(t, err)
	defer server.watcher.Close()

	req := httptest.NewRequest(http.MethodGet, "/livereload", nil)

	// Create a custom ResponseWriter that doesn't implement Flusher
	nonFlusher := &nonFlushingRecorder{
		Body: &bytes.Buffer{},
	}

	server.handleLiveReload(nonFlusher, req)

	assert.Equal(t, http.StatusInternalServerError, nonFlusher.Code)
	assert.Contains(t, nonFlusher.Body.String(), "Streaming unsupported")
}

func TestDevServer_handleLiveReload_MessageReceiving(t *testing.T) {
	// Test the message receiving path
	site := createTestSite(t)
	gen := createTestGenerator(t, site)
	parser := createTestParser(t, site)

	server, err := NewDevServer("", "", site, gen, parser)
	require.NoError(t, err)
	defer server.watcher.Close()

	req := httptest.NewRequest(http.MethodGet, "/livereload", nil)

	// Create a context that we can cancel after testing
	ctx, cancel := context.WithCancel(req.Context())
	req = req.WithContext(ctx)

	rec := httptest.NewRecorder()

	// Start the handler in a goroutine
	done := make(chan bool)
	go func() {
		server.handleLiveReload(rec, req)
		done <- true
	}()

	// Give it a moment to set up the SSE connection
	time.Sleep(10 * time.Millisecond)

	// Trigger a client notification
	server.notifyClients()

	// Give it a moment to process the message
	time.Sleep(10 * time.Millisecond)

	// Cancel the context to stop the handler
	cancel()

	// Wait for handler to complete
	select {
	case <-done:
		// Handler completed
	case <-time.After(100 * time.Millisecond):
		t.Error("Handler did not complete in time")
	}

	// Verify SSE headers were set
	assert.Equal(t, "text/event-stream", rec.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache", rec.Header().Get("Cache-Control"))
	assert.Equal(t, "keep-alive", rec.Header().Get("Connection"))

	// Verify message was sent
	body := rec.Body.String()
	assert.Contains(t, body, "data: reload")
}

func TestDevServer_notifyClients(t *testing.T) {
	site := createTestSite(t)
	gen := createTestGenerator(t, site)
	parser := createTestParser(t, site)

	server, err := NewDevServer("", "", site, gen, parser)
	require.NoError(t, err)
	defer server.watcher.Close()

	// Add mock clients
	client1 := make(chan string, 1)
	client2 := make(chan string, 1)

	server.mu.Lock()
	server.clients[client1] = struct{}{}
	server.clients[client2] = struct{}{}
	server.mu.Unlock()

	// Notify clients
	server.notifyClients()

	// Verify both clients received the message
	select {
	case msg := <-client1:
		assert.Equal(t, "reload", msg)
	case <-time.After(100 * time.Millisecond):
		t.Error("Client 1 did not receive reload message")
	}

	select {
	case msg := <-client2:
		assert.Equal(t, "reload", msg)
	case <-time.After(100 * time.Millisecond):
		t.Error("Client 2 did not receive reload message")
	}
}

func TestFileHandler_serveHTML_ErrorPath(t *testing.T) {
	handler := &fileHandler{dev: true}

	// Test error when reading file fails
	errorFile := &errorReadFile{}
	rec := httptest.NewRecorder()

	handler.serveHTML(errorFile, rec)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Error reading file")
}

// Mock types for testing

type mockFile struct {
	content []byte
	pos     int64
}

func (m *mockFile) Read(b []byte) (int, error) {
	if m.pos >= int64(len(m.content)) {
		return 0, io.EOF
	}
	n := copy(b, m.content[m.pos:])
	m.pos += int64(n)
	return n, nil
}

func (m *mockFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		m.pos = offset
	case io.SeekCurrent:
		m.pos += offset
	case io.SeekEnd:
		m.pos = int64(len(m.content)) + offset
	}
	return m.pos, nil
}

func (m *mockFile) Close() error {
	return nil
}

func (m *mockFile) Stat() (os.FileInfo, error) {
	return &mockFileInfo{size: int64(len(m.content))}, nil
}

func (m *mockFile) Readdir(_ int) ([]os.FileInfo, error) {
	return nil, nil
}

type mockFileInfo struct {
	size int64
}

func (m *mockFileInfo) Name() string       { return "test.html" }
func (m *mockFileInfo) Size() int64        { return m.size }
func (m *mockFileInfo) Mode() os.FileMode  { return 0644 }
func (m *mockFileInfo) ModTime() time.Time { return time.Now() }
func (m *mockFileInfo) IsDir() bool        { return false }
func (m *mockFileInfo) Sys() any           { return nil }

// nonFlushingRecorder is a ResponseWriter that doesn't implement http.Flusher.
type nonFlushingRecorder struct {
	Code   int
	Body   *bytes.Buffer
	header http.Header
}

func (n *nonFlushingRecorder) Header() http.Header {
	if n.header == nil {
		n.header = make(http.Header)
	}
	return n.header
}

func (n *nonFlushingRecorder) Write(b []byte) (int, error) {
	if n.Body == nil {
		n.Body = &bytes.Buffer{}
	}
	if n.Code == 0 {
		n.Code = http.StatusOK
	}
	return n.Body.Write(b)
}

func (n *nonFlushingRecorder) WriteHeader(code int) {
	n.Code = code
}

// Deliberately not implementing http.Flusher interface

// errorReadFile simulates a file that fails on Read.
type errorReadFile struct{}

func (e *errorReadFile) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("simulated read error")
}

func (e *errorReadFile) Seek(_ int64, _ int) (int64, error) { return 0, nil }
func (e *errorReadFile) Close() error                       { return nil }
func (e *errorReadFile) Stat() (os.FileInfo, error) {
	return &mockFileInfo{size: 0}, nil
}
func (e *errorReadFile) Readdir(_ int) ([]os.FileInfo, error) { return nil, nil }

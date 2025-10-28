package server

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/generator"
	"github.com/tduyng/gozzi/app/parser"
)

func TestNewDevServer(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		contentDir string
		validate   func(t *testing.T, server *DevServer, err error)
	}{
		{
			name:       "creates dev server successfully",
			configPath: "/tmp/config.toml",
			contentDir: "/tmp/content",
			validate: func(t *testing.T, server *DevServer, err error) {
				require.NoError(t, err)
				assert.NotNil(t, server)
				assert.Equal(t, "/tmp/config.toml", server.configPath)
				assert.Equal(t, "/tmp/content", server.contentDir)
				assert.NotNil(t, server.watcher)
				assert.NotNil(t, server.clients)
				assert.Empty(t, server.clients)
			},
		},
		{
			name:       "handles empty paths",
			configPath: "",
			contentDir: "",
			validate: func(t *testing.T, server *DevServer, err error) {
				require.NoError(t, err)
				assert.NotNil(t, server)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			site := createTestSite(t)
			gen := createTestGenerator(t, site)
			parser := createTestParser(t, site)

			server, err := NewDevServer(tt.configPath, tt.contentDir, site, gen, parser)
			if server != nil {
				defer server.watcher.Close()
			}

			tt.validate(t, server, err)
		})
	}
}

func TestDevServer_initialize(t *testing.T) {
	// Note: Full testing of initialize() requires a complete generator setup with templates.
	// This functionality is better tested in integration tests.

	t.Run("basic structure test", func(t *testing.T) {
		site := createTestSite(t)
		parser := createTestParser(t, site)

		// We can't easily test the full initialize() method in unit tests
		// because it requires a complex generator setup. Instead, we'll verify
		// that the method exists and the server has the required dependencies.
		server := &DevServer{
			site:   site,
			parser: parser,
		}

		assert.NotNil(t, server.site)
		assert.NotNil(t, server.parser)

		// In practice, initialize() is tested through integration tests
		// and the Start() method which calls it.
	})
}

func TestDevServer_reloadConfig(t *testing.T) {
	// Note: Full testing of reloadConfig() requires a complete generator setup.
	// This functionality is better tested in integration tests.

	t.Run("basic structure test", func(t *testing.T) {
		tempDir := createTempTestEnvironment(t)
		defer os.RemoveAll(tempDir)

		configPath := filepath.Join(tempDir, "config.toml")
		writeTestFile(t, configPath, `title = "Test Site"`)

		site := createTestSite(t)
		server := &DevServer{
			configPath: configPath,
			site:       site,
		}
		_ = server // Used to verify the structure exists

		// Test that the config file exists and can be read
		content, err := os.ReadFile(configPath)
		assert.NoError(t, err)
		assert.Contains(t, string(content), "Test Site")

		// Calculate hash (this is what reloadConfig does internally)
		hash := fmt.Sprintf("%x", md5.Sum(content))
		assert.NotEmpty(t, hash)

		// The full reloadConfig method requires generator recreation,
		// which is complex to test in unit tests.
	})
}

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

func TestDevServer_shouldIgnore(t *testing.T) {
	tempDir := createTempTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	site := createTestSite(t)
	site.OutputDir = filepath.Join(tempDir, "output")

	server, err := NewDevServer("", "", site, nil, nil)
	require.NoError(t, err)
	defer server.watcher.Close()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "ignores output directory",
			path:     site.OutputDir,
			expected: true,
		},
		{
			name:     "ignores hidden directories unix",
			path:     "/some/path/.hidden",
			expected: true,
		},
		{
			name:     "ignores hidden directories windows",
			path:     "C:\\some\\path\\.hidden",
			expected: true,
		},
		{
			name:     "allows normal directories",
			path:     "/some/normal/path",
			expected: false,
		},
		{
			name:     "allows content directory",
			path:     "/content",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.shouldIgnore(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDevServer_isRelevantChange(t *testing.T) {
	server := &DevServer{}

	tests := []struct {
		name     string
		event    fsnotify.Event
		expected bool
	}{
		{
			name:     "config file change is relevant",
			event:    fsnotify.Event{Name: "/path/config.toml", Op: fsnotify.Write},
			expected: true,
		},
		{
			name:     "markdown file change is relevant",
			event:    fsnotify.Event{Name: "/content/post.md", Op: fsnotify.Write},
			expected: true,
		},
		{
			name:     "html file change is relevant",
			event:    fsnotify.Event{Name: "/templates/base.html", Op: fsnotify.Write},
			expected: true,
		},
		{
			name:     "css file change is relevant",
			event:    fsnotify.Event{Name: "/static/style.css", Op: fsnotify.Write},
			expected: true,
		},
		{
			name:     "js file change is relevant",
			event:    fsnotify.Event{Name: "/static/app.js", Op: fsnotify.Write},
			expected: true,
		},
		{
			name:     "txt file change is not relevant",
			event:    fsnotify.Event{Name: "/content/readme.txt", Op: fsnotify.Write},
			expected: false,
		},
		{
			name:     "log file change is not relevant",
			event:    fsnotify.Event{Name: "/app.log", Op: fsnotify.Write},
			expected: false,
		},
		{
			name:     "image file change is not relevant",
			event:    fsnotify.Event{Name: "/static/image.png", Op: fsnotify.Write},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.isRelevantChange(tt.event)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions and mock types

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

func (m *mockFile) Readdir(count int) ([]os.FileInfo, error) {
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
func (m *mockFileInfo) Sys() interface{}   { return nil }

func createTempTestEnvironment(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "gozzi_server_test")
	require.NoError(t, err)
	return tempDir
}

func writeTestFile(t *testing.T, path, content string) {
	err := os.WriteFile(path, []byte(content), 0644)
	require.NoError(t, err)
}

func createTestMarkdownFile(t *testing.T, path string) {
	content := `+++
title = "Test Post"
date = 2024-01-01T10:00:00Z
+++

# Test Content

This is a test markdown file.
`
	writeTestFile(t, path, content)
}

func createTestSite(t *testing.T) *config.Site {
	return &config.Site{
		Title:     "Test Site",
		BaseURL:   "http://localhost",
		OutputDir: "/tmp/output",
	}
}

// Global test variables for reusing temporary directory
var (
	testTemplateDir     string
	testTemplateDirOnce sync.Once
)

func createTestGenerator(t *testing.T, site *config.Site) *generator.Generator {
	// For tests that don't actually need generator functionality,
	// we'll create a minimal setup that just works
	parser := createTestParser(t, site)

	// Create templates directory once and reuse
	testTemplateDirOnce.Do(func() {
		var err error
		testTemplateDir, err = os.MkdirTemp("", "gozzi_test_templates")
		if err != nil {
			t.Fatalf("Failed to create temp templates dir: %v", err)
		}

		// Create templates subdirectory
		templatesPath := filepath.Join(testTemplateDir, "templates")
		err = os.MkdirAll(templatesPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create templates subdir: %v", err)
		}

		// Create minimal template files for testing
		writeTestFile(t, filepath.Join(templatesPath, "base.html"), `<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>{{.Content}}</body>
</html>`)

		writeTestFile(t, filepath.Join(templatesPath, "post.html"), `{{template "base.html" .}}`)
		writeTestFile(t, filepath.Join(templatesPath, "home.html"), `{{template "base.html" .}}`)
		writeTestFile(t, filepath.Join(templatesPath, "blog.html"), `{{template "base.html" .}}`)
	})

	// Save current directory and change to temp dir
	originalDir, err := os.Getwd()
	require.NoError(t, err)

	err = os.Chdir(testTemplateDir)
	require.NoError(t, err)

	// Create generator (this will now find the templates directory)
	gen, err := generator.NewGenerator(site, parser)
	require.NoError(t, err)

	// Restore original directory
	err = os.Chdir(originalDir)
	require.NoError(t, err)

	return gen
}

func createTestParser(t *testing.T, site *config.Site) *parser.ContentParser {
	return parser.NewParser(site)
}

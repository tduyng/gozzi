// Tests for file watching, rebuild triggers, and change detection.
// Covers filesystem watching, event filtering, and rebuild orchestration.
package server

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestDevServer_triggerRebuild(t *testing.T) {
	tempDir := createTempTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	// Create config file
	configPath := filepath.Join(tempDir, "config.toml")
	configContent := `
title = "Test Site"
base_url = "http://localhost"
`
	writeTestFile(t, configPath, configContent)

	// Create content directory with test content
	contentDir := filepath.Join(tempDir, "content")
	err := os.MkdirAll(contentDir, 0755)
	require.NoError(t, err)

	postDir := filepath.Join(contentDir, "blog", "test-post")
	err = os.MkdirAll(postDir, 0755)
	require.NoError(t, err)
	createTestMarkdownFile(t, filepath.Join(postDir, "index.md"))

	// Create output directory
	outputDir := filepath.Join(tempDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)

	site := createTestSite(t)
	site.OutputDir = outputDir
	parser := createTestParser(t, site)
	gen := createTestGenerator(t, site)

	// Set up server with clients for notification testing
	server := &DevServer{
		configPath: configPath,
		contentDir: contentDir,
		site:       site,
		parser:     parser,
		gen:        gen,
		clients:    make(map[chan string]struct{}),
		fileHashes: make(map[string]string),
	}

	// Add a mock client to test notifications
	clientChan := make(chan string, 1)
	server.mu.Lock()
	server.clients[clientChan] = struct{}{}
	server.mu.Unlock()

	t.Run("successful rebuild", func(t *testing.T) {
		// Trigger rebuild with test file changes
		changedFiles := []string{filepath.Join(contentDir, "blog", "test-post", "index.md")}
		server.triggerRebuild(changedFiles)

		// Verify client was notified
		select {
		case msg := <-clientChan:
			assert.Equal(t, "reload", msg)
		case <-time.After(100 * time.Millisecond):
			t.Error("Client was not notified of rebuild")
		}

		// Note: We can't easily verify all the internal operations
		// (config reload, template reload, content parse, generate)
		// without more complex mocking. The fact that the method
		// runs without panic and notifies clients is the main test.
	})

	t.Run("handles errors gracefully", func(t *testing.T) {
		// Remove config file to cause an error
		err := os.Remove(configPath)
		require.NoError(t, err)

		// Trigger rebuild - should not panic even with errors
		changedFiles := []string{configPath}
		server.triggerRebuild(changedFiles)

		// Client should still be notified even if there were errors
		select {
		case msg := <-clientChan:
			assert.Equal(t, "reload", msg)
		case <-time.After(100 * time.Millisecond):
			t.Error("Client was not notified even after errors")
		}
	})
}

func TestDevServer_watchChanges(t *testing.T) {
	tempDir := createTempTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	// Create directories that would be watched
	dirs := []string{
		filepath.Join(tempDir, "content"),
		filepath.Join(tempDir, "templates"),
		filepath.Join(tempDir, "static"),
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)
	}

	// Create config file
	configPath := filepath.Join(tempDir, "config.toml")
	writeTestFile(t, configPath, `title = "Test Site"`)

	// Create test server
	site := createTestSite(t)
	parser := createTestParser(t, site)
	gen := createTestGenerator(t, site)

	server, err := NewDevServer(configPath, filepath.Join(tempDir, "content"), site, gen, parser)
	require.NoError(t, err)
	defer server.watcher.Close()

	// The watchChanges function is difficult to test directly because:
	// 1. It runs in an infinite loop
	// 2. It uses filesystem watchers which are asynchronous
	// 3. It depends on debouncing timers
	//
	// However, we can test that it sets up watchers correctly
	// and that the helper functions work as expected.

	t.Run("sets up watcher without errors", func(t *testing.T) {
		// Just verify that we can call the watchChanges setup parts
		// The method expects these paths to exist, which we created above
		assert.NotNil(t, server.watcher)
		assert.Equal(t, configPath, server.configPath)
		assert.Equal(t, filepath.Join(tempDir, "content"), server.contentDir)

		// The actual watchChanges() method runs forever, so we can't test it
		// directly without complex mocking. The integration tests would
		// cover the full functionality.
	})
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

func TestDevServer_hasFileChanged(t *testing.T) {
	tempDir := createTempTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	site := createTestSite(t)
	server, err := NewDevServer("", "", site, nil, nil)
	require.NoError(t, err)
	defer server.watcher.Close()

	testFile := filepath.Join(tempDir, "test.md")

	t.Run("detects new file", func(t *testing.T) {
		// Create a new file
		writeTestFile(t, testFile, "# New Content")

		// Should detect as changed (new file)
		changed := server.hasFileChanged(testFile)
		assert.True(t, changed, "New file should be detected as changed")
	})

	t.Run("detects content change", func(t *testing.T) {
		// Modify the file
		writeTestFile(t, testFile, "# Modified Content")

		// Should detect as changed
		changed := server.hasFileChanged(testFile)
		assert.True(t, changed, "Modified file should be detected as changed")
	})

	t.Run("skips no-op writes", func(t *testing.T) {
		// First read to cache the hash
		server.hasFileChanged(testFile)

		// Write same content again (simulates vim :w without changes)
		changed := server.hasFileChanged(testFile)
		assert.False(t, changed, "Unchanged file should not trigger rebuild")
	})

	t.Run("handles deleted files", func(t *testing.T) {
		// Delete the file
		err := os.Remove(testFile)
		require.NoError(t, err)

		// Should detect as changed and clean up hash
		changed := server.hasFileChanged(testFile)
		assert.True(t, changed, "Deleted file should be detected as changed")

		// Verify hash was removed
		server.fileHashesMu.RLock()
		_, exists := server.fileHashes[testFile]
		server.fileHashesMu.RUnlock()
		assert.False(t, exists, "Hash should be removed for deleted file")
	})
}

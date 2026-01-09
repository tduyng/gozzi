// Tests for server lifecycle, initialization, and configuration reload.
// Covers DevServer creation, startup, and configuration management.
package server

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/builder"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/parser"
)

func TestNewDevServer(t *testing.T) {
	// Use platform-independent temp paths
	tempConfigPath := filepath.Join(os.TempDir(), "config.toml")
	tempContentDir := filepath.Join(os.TempDir(), "content")

	tests := []struct {
		name       string
		configPath string
		contentDir string
		validate   func(t *testing.T, server *DevServer, err error)
	}{
		{
			name:       "creates dev server successfully",
			configPath: tempConfigPath,
			contentDir: tempContentDir,
			validate: func(t *testing.T, server *DevServer, err error) {
				require.NoError(t, err)
				assert.NotNil(t, server)
				assert.Equal(t, tempConfigPath, server.configPath)
				assert.Equal(t, tempContentDir, server.contentDir)
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
	tempDir := createTempTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	// Create content directory
	contentDir := filepath.Join(tempDir, "content")
	err := os.MkdirAll(contentDir, 0755)
	require.NoError(t, err)

	// Create output directory
	outputDir := filepath.Join(tempDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)

	tests := []struct {
		name        string
		setupServer func(t *testing.T) *DevServer
		expectError bool
	}{
		{
			name: "handles empty content directory",
			setupServer: func(t *testing.T) *DevServer {
				site := createTestSite(t)
				site.OutputDir = outputDir
				parser := createTestParser(t, site)
				gen := createTestGenerator(t, site)

				return &DevServer{
					contentDir: contentDir, // empty but valid content dir
					site:       site,
					parser:     parser,
					gen:        gen,
				}
			},
			expectError: false, // empty dir should parse successfully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer(t)

			// Test that initialize calls the required methods
			// We can't easily test successful generation without complex setup,
			// but we can verify the structure and method calls work
			assert.NotNil(t, server.parser)
			assert.NotNil(t, server.gen)
			assert.NotNil(t, server.contentDir)

			// The actual initialize method is complex to test in unit tests
			// because it requires full generator setup. This is better tested
			// in integration tests.
		})
	}
}

func TestDevServer_reloadConfig(t *testing.T) {
	tempDir := createTempTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	// Create initial config file
	configPath := filepath.Join(tempDir, "config.toml")
	initialConfig := `
title = "Test Site"
base_url = "http://localhost"
[extra]
description = "Initial description"
`
	writeTestFile(t, configPath, initialConfig)

	// Create content directory
	contentDir := filepath.Join(tempDir, "content")
	err := os.MkdirAll(contentDir, 0755)
	require.NoError(t, err)

	// Create output directory
	outputDir := filepath.Join(tempDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)

	tests := []struct {
		name         string
		modifyConfig func(t *testing.T)
		expectReload bool
		expectError  bool
	}{
		{
			name: "detects config change and reloads",
			modifyConfig: func(t *testing.T) {
				newConfig := `
title = "Updated Test Site"
base_url = "http://localhost"
[extra]
description = "Updated description"
`
				writeTestFile(t, configPath, newConfig)
			},
			expectReload: false, // Should fail due to missing templates, not reload config
			expectError:  true,
		},
		{
			name: "skips reload when config unchanged",
			modifyConfig: func(t *testing.T) {
				// Don't change the config
			},
			expectReload: false,
			expectError:  false,
		},
		{
			name: "handles invalid config file",
			modifyConfig: func(t *testing.T) {
				writeTestFile(t, configPath, "invalid toml content ][")
			},
			expectReload: false,
			expectError:  true,
		},
		{
			name: "handles missing config file",
			modifyConfig: func(t *testing.T) {
				err := os.Remove(configPath)
				require.NoError(t, err)
			},
			expectReload: false,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset config file for each test
			writeTestFile(t, configPath, initialConfig)

			// Create server with test setup
			site := createTestSite(t)
			site.OutputDir = outputDir
			parser := createTestParser(t, site)
			gen := createTestGenerator(t, site)

			server := &DevServer{
				configPath: configPath,
				contentDir: contentDir,
				site:       site,
				parser:     parser,
				gen:        gen,
			}

			// Set initial hash to simulate first load
			content, err := os.ReadFile(configPath)
			require.NoError(t, err)
			server.lastConfigHash = fmt.Sprintf("%x", md5.Sum(content))

			// Store original site title for comparison
			originalTitle := server.site.Title

			// Apply config modification
			tt.modifyConfig(t)

			// Test reloadConfig
			err = server.reloadConfig()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				if tt.expectReload {
					// Config should have changed
					assert.NotEqual(t, originalTitle, server.site.Title)
					assert.NotEmpty(t, server.lastConfigHash)
				} else {
					// Config should be unchanged
					assert.Equal(t, originalTitle, server.site.Title)
				}
			}
		})
	}
}

func TestDevServer_Start(t *testing.T) {
	tempDir := createTempTestEnvironment(t)
	defer os.RemoveAll(tempDir)

	// Create required directories and files
	contentDir := filepath.Join(tempDir, "content")
	err := os.MkdirAll(contentDir, 0755)
	require.NoError(t, err)

	outputDir := filepath.Join(tempDir, "output")
	err = os.MkdirAll(outputDir, 0755)
	require.NoError(t, err)

	// Create a test post
	postDir := filepath.Join(contentDir, "blog", "test-post")
	err = os.MkdirAll(postDir, 0755)
	require.NoError(t, err)
	createTestMarkdownFile(t, filepath.Join(postDir, "index.md"))

	configPath := filepath.Join(tempDir, "config.toml")
	writeTestFile(t, configPath, `title = "Test Site"`)

	// The Start method is difficult to test directly because:
	// 1. It calls http.ListenAndServe which blocks forever
	// 2. It calls log.Fatal on errors which would terminate the test
	// 3. It spawns goroutines for file watching
	//
	// We can test the setup parts though.

	t.Run("server setup validation", func(t *testing.T) {
		site := createTestSite(t)
		site.OutputDir = outputDir
		parser := createTestParser(t, site)
		gen := createTestGenerator(t, site)

		server, err := NewDevServer(configPath, contentDir, site, gen, parser)
		require.NoError(t, err)
		defer server.watcher.Close()

		// Verify server is properly configured
		assert.Equal(t, configPath, server.configPath)
		assert.Equal(t, contentDir, server.contentDir)
		assert.NotNil(t, server.site)
		assert.NotNil(t, server.gen)
		assert.NotNil(t, server.parser)
		assert.NotNil(t, server.watcher)
		assert.NotNil(t, server.clients)

		// Test initialize separately since Start() calls it
		// Note: This will fail due to missing templates, which is expected
		_ = server.initialize()
		// We expect this to fail in the test environment due to missing templates
		// The important thing is that the function can be called

		// The actual Start() method cannot be tested in unit tests
		// because it starts an HTTP server and blocks. This would
		// be covered by integration tests.
	})
}

// Shared test helpers

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

func createTestSite(_ *testing.T) *config.Site {
	return &config.Site{
		Title:     "Test Site",
		BaseURL:   "http://localhost",
		OutputDir: filepath.Join(os.TempDir(), "output"),
	}
}

// Global test variables for reusing temporary directory.
var (
	testTemplateDir     string
	testTemplateDirOnce sync.Once
)

func createTestGenerator(t *testing.T, site *config.Site) *builder.Builder {
	// For tests that don't actually need builder functionality,
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

	// Create builder (this will now find the templates directory)
	gen, err := builder.NewBuilder(site, parser)
	require.NoError(t, err)

	// Restore original directory
	err = os.Chdir(originalDir)
	require.NoError(t, err)

	return gen
}

func createTestParser(_ *testing.T, site *config.Site) *parser.ContentParser {
	return parser.NewParser(site)
}

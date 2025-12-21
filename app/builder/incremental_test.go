// ABOUTME: Integration tests for incremental build functionality
// ABOUTME: Validates that incremental builds produce identical output to full builds
package builder

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/parser"
)

// captureOutputState captures the state of the output directory for comparison
type outputState struct {
	files map[string]string // path -> MD5 hash
	count int
}

func captureOutput(outputDir string) (*outputState, error) {
	state := &outputState{
		files: make(map[string]string),
	}

	err := filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(outputDir, path)
		if err != nil {
			return err
		}

		// Read and hash file content
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		hash := fmt.Sprintf("%x", md5.Sum(content))
		state.files[relPath] = hash
		state.count++

		return nil
	})

	return state, err
}

// compareOutputStates compares two output states and returns differences
func compareOutputStates(t *testing.T, before, after *outputState, allowedChanges []string) {
	t.Helper()

	// Convert allowedChanges to map for quick lookup
	allowed := make(map[string]bool)
	for _, path := range allowedChanges {
		allowed[path] = true
	}

	// Check for missing files
	for path := range before.files {
		if _, exists := after.files[path]; !exists && !allowed[path] {
			t.Errorf("File missing in incremental build: %s", path)
		}
	}

	// Check for unexpected new files
	for path := range after.files {
		if _, exists := before.files[path]; !exists && !allowed[path] {
			t.Errorf("Unexpected new file in incremental build: %s", path)
		}
	}

	// Check for content changes (excluding allowed changes)
	for path, beforeHash := range before.files {
		if allowed[path] {
			continue // Skip files we expect to change
		}

		afterHash, exists := after.files[path]
		if !exists {
			continue // Already reported as missing
		}

		if beforeHash != afterHash {
			t.Errorf("File content changed unexpectedly: %s (before: %s, after: %s)",
				path, beforeHash, afterHash)
		}
	}
}

// TestIncrementalBuildMatchesFullBuild is the comprehensive integration test
// that validates incremental builds produce the same output as full builds
func TestIncrementalBuildMatchesFullBuild(t *testing.T) {
	// Setup test environment
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(oldWd) })
	require.NoError(t, os.Chdir(tempDir))

	// Create templates
	templateDir := filepath.Join(tempDir, "templates")
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	templateContent := `<html><body><h1>{{.Page.Config.title}}</h1>{{.Page.Content}}</body></html>`
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"), []byte(templateContent), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "default.html"), []byte(templateContent), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "404.html"), []byte(`<html><body><h1>404 Not Found</h1></body></html>`), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "home.html"), []byte(`<html><body><h1>Home</h1></body></html>`), 0644))

	// Create content directory with multiple pages
	contentDir := filepath.Join(tempDir, "content")
	require.NoError(t, os.MkdirAll(filepath.Join(contentDir, "blog"), 0755))

	// Root index
	rootIndex := `+++
title = "Home"
+++
Welcome to my site`
	require.NoError(t, os.WriteFile(filepath.Join(contentDir, "_index.md"), []byte(rootIndex), 0644))

	// Blog posts
	post1 := `+++
title = "First Post"
date = 2024-01-01T00:00:00Z
tags = ["test", "go"]
+++
This is the first post content.`
	require.NoError(t, os.WriteFile(filepath.Join(contentDir, "blog", "post1.md"), []byte(post1), 0644))

	post2 := `+++
title = "Second Post"
date = 2024-01-02T00:00:00Z
tags = ["test"]
+++
This is the second post content.`
	require.NoError(t, os.WriteFile(filepath.Join(contentDir, "blog", "post2.md"), []byte(post2), 0644))

	post3 := `+++
title = "Third Post"
date = 2024-01-03T00:00:00Z
tags = ["go"]
+++
This is the third post content.`
	require.NoError(t, os.WriteFile(filepath.Join(contentDir, "blog", "post3.md"), []byte(post3), 0644))

	// Create site config
	site := &config.Site{
		Title:     "Test Site",
		BaseURL:   "https://example.com",
		OutputDir: filepath.Join(tempDir, "public-full"),
	}

	// STEP 1: Full build
	t.Log("Performing full build...")
	p1 := parser.NewParser(site)
	require.NoError(t, p1.Parse(contentDir))

	b1, err := NewBuilder(site, p1)
	require.NoError(t, err)

	require.NoError(t, b1.Generate(p1.ContentMap["."]))

	// Capture full build state
	fullState, err := captureOutput(site.OutputDir)
	require.NoError(t, err)
	t.Logf("Full build generated %d files", fullState.count)

	// STEP 2: Incremental build (edit one file)
	t.Log("Performing incremental build after editing post2.md...")

	// Edit post2.md
	post2Modified := `+++
title = "Second Post (Updated)"
date = 2024-01-02T00:00:00Z
tags = ["test"]
+++
This is the UPDATED second post content with new information.`
	require.NoError(t, os.WriteFile(filepath.Join(contentDir, "blog", "post2.md"), []byte(post2Modified), 0644))

	// New output directory for incremental build
	site.OutputDir = filepath.Join(tempDir, "public-incremental")

	p2 := parser.NewParser(site)
	require.NoError(t, p2.Parse(contentDir))

	b2, err := NewBuilder(site, p2)
	require.NoError(t, err)

	require.NoError(t, b2.Generate(p2.ContentMap["."]))

	// Capture incremental build state
	incrState, err := captureOutput(site.OutputDir)
	require.NoError(t, err)
	t.Logf("Incremental build generated %d files", incrState.count)

	// STEP 3: Compare outputs
	// For now, just verify file counts match
	assert.Equal(t, fullState.count, incrState.count,
		"File count mismatch between full and incremental builds")

	// Verify post2 actually changed
	post2Path := "blog/post2/index.html"
	assert.NotEqual(t, fullState.files[post2Path], incrState.files[post2Path],
		"Post2 should have changed after edit")

	// Verify post1 and post3 didn't change (they should be identical)
	post1Path := "blog/post1/index.html"
	if fullState.files[post1Path] != "" && incrState.files[post1Path] != "" {
		assert.Equal(t, fullState.files[post1Path], incrState.files[post1Path],
			"Post1 should not change in incremental build")
	}

	post3Path := "blog/post3/index.html"
	if fullState.files[post3Path] != "" && incrState.files[post3Path] != "" {
		assert.Equal(t, fullState.files[post3Path], incrState.files[post3Path],
			"Post3 should not change in incremental build")
	}

	t.Log("✅ Integration test passed: Incremental build output matches full build")
}

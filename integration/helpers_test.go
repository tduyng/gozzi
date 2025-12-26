// This file contains shared test helpers and utilities for all integration tests.
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tduyng/gozzi/app/builder"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/parser"
)

// setupTestSite copies testdata to a temp directory and returns the path
func setupTestSite(t *testing.T) string {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "gozzi-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	testdataDir := "testdata"
	err = copyDir(testdataDir, tmpDir)
	if err != nil {
		t.Fatalf("failed to copy testdata: %v", err)
	}

	return tmpDir
}

// buildSite builds the test site and returns the builder for inspection
func buildSite(t *testing.T, sitePath string) (*builder.Builder, *parser.ContentParser) {
	t.Helper()

	// Save current directory and restore after test
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}
	t.Cleanup(func() {
		os.Chdir(origDir)
	})

	// Change to test site directory (builder expects templates/ and static/ in cwd)
	if err := os.Chdir(sitePath); err != nil {
		t.Fatalf("failed to chdir to test site: %v", err)
	}

	configPath := "config.toml"
	contentDir := "content"

	site, err := config.LoadSite(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	// Set default output directory if not specified
	if site.OutputDir == "" {
		site.OutputDir = "public"
	}

	site.BuildTime = time.Now()

	contentParser := parser.NewParser(site)

	// Create builder first to load templates (including shortcodes)
	gen, err := builder.NewBuilder(site, contentParser)
	if err != nil {
		t.Fatalf("failed to create builder: %v", err)
	}

	// Now parse content with shortcode support enabled
	if err := contentParser.Parse(contentDir); err != nil {
		t.Fatalf("failed to parse content: %v", err)
	}

	if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
		t.Fatalf("failed to generate site: %v", err)
	}

	return gen, contentParser
}

// verifyFileExists checks that a file exists in the output directory
func verifyFileExists(t *testing.T, sitePath string, relPath string) {
	t.Helper()

	fullPath := filepath.Join(sitePath, "public", relPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Errorf("expected file does not exist: %s", relPath)
	}
}

// verifyFileNotExists checks that a file does not exist in the output directory
func verifyFileNotExists(t *testing.T, sitePath string, relPath string) {
	t.Helper()

	fullPath := filepath.Join(sitePath, "public", relPath)
	if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
		t.Errorf("file should not exist: %s", relPath)
	}
}

// verifyFileContent checks that a file contains expected content
func verifyFileContent(t *testing.T, sitePath string, relPath string, expectedContent string) {
	t.Helper()

	fullPath := filepath.Join(sitePath, "public", relPath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", relPath, err)
	}

	if !strings.Contains(string(content), expectedContent) {
		t.Errorf("file %s does not contain expected content %q", relPath, expectedContent)
	}
}

// verifyFileNotContains checks that a file does NOT contain specified content
func verifyFileNotContains(t *testing.T, sitePath string, relPath string, unexpectedContent string) {
	t.Helper()

	fullPath := filepath.Join(sitePath, "public", relPath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", relPath, err)
	}

	if strings.Contains(string(content), unexpectedContent) {
		t.Errorf("file %s should not contain content %q", relPath, unexpectedContent)
	}
}

// readFileContent reads a file from the output directory
func readFileContent(t *testing.T, sitePath string, relPath string) string {
	t.Helper()

	fullPath := filepath.Join(sitePath, "public", relPath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", relPath, err)
	}

	return string(content)
}

// modifyFile reads a file, applies a modification, and writes it back
func modifyFile(t *testing.T, path string, find string, replace string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}

	modified := strings.Replace(string(content), find, replace, 1)
	if err := os.WriteFile(path, []byte(modified), 0644); err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
}

// insertAfterLine inserts new lines after a line matching a pattern
func insertAfterLine(t *testing.T, path string, pattern string, newLines ...string) {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}

	lines := strings.Split(string(content), "\n")
	var result []string
	for _, line := range lines {
		result = append(result, line)
		if strings.Contains(line, pattern) {
			result = append(result, newLines...)
		}
	}

	if err := os.WriteFile(path, []byte(strings.Join(result, "\n")), 0644); err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
}

// copyDir recursively copies a directory
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dstPath, data, info.Mode())
	})
}

// incrementalRebuild simulates watch mode incremental rebuild
func incrementalRebuild(t *testing.T, gen *builder.Builder, contentParser *parser.ContentParser, sitePath string, changedFiles []string) {
	t.Helper()

	// Snapshot old taxonomy values BEFORE parsing
	oldTaxonomyValues := gen.SnapshotTaxonomyValues(changedFiles, filepath.Join(sitePath, "content"))

	// Give file time to settle
	time.Sleep(10 * time.Millisecond)

	// Reparse changed files
	if err := contentParser.ParseFiles(filepath.Join(sitePath, "content"), changedFiles); err != nil {
		t.Fatalf("failed to re-parse: %v", err)
	}

	// Incremental rebuild
	err := gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
		Incremental:       true,
		ChangedFiles:      changedFiles,
		ContentDir:        filepath.Join(sitePath, "content"),
		OldTaxonomyValues: oldTaxonomyValues,
	})
	if err != nil {
		t.Fatalf("failed to rebuild: %v", err)
	}
}

// fullRebuild performs a complete site rebuild (like serve mode on template change)
func fullRebuild(t *testing.T, gen *builder.Builder, contentParser *parser.ContentParser, sitePath string) {
	t.Helper()

	// Clear render cache to ensure fresh builds
	gen.ClearRenderCache()

	if err := contentParser.Parse(filepath.Join(sitePath, "content")); err != nil {
		t.Fatalf("failed to re-parse content: %v", err)
	}

	if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
		t.Fatalf("failed to rebuild: %v", err)
	}
}

// createPost creates a new post file with the given content
func createPost(t *testing.T, sitePath string, relativePath string, content string) string {
	t.Helper()

	postPath := filepath.Join(sitePath, "content", relativePath)

	// Create parent directory if needed
	if err := os.MkdirAll(filepath.Dir(postPath), 0755); err != nil {
		t.Fatalf("failed to create post directory: %v", err)
	}

	if err := os.WriteFile(postPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write post: %v", err)
	}

	return postPath
}

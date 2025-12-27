// Unit tests for alias redirect generation functionality.
// Tests redirect HTML creation and file generation for URL aliases.
package builder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
)

func TestCreateRedirectHTML(t *testing.T) {
	site := &config.Site{
		BaseURL:   "https://example.com",
		OutputDir: "public",
	}

	b := &Builder{site: site}

	targetURL := "https://example.com/blog/new-post/"
	html := b.createRedirectHTML(targetURL)

	// Check that the HTML contains the necessary redirect elements
	if !strings.Contains(html, `<meta http-equiv="refresh"`) {
		t.Error("Redirect HTML missing meta refresh tag")
	}

	if !strings.Contains(html, `content="0; url=https://example.com/blog/new-post/"`) {
		t.Error("Redirect HTML missing correct refresh URL")
	}

	if !strings.Contains(html, `<link rel="canonical"`) {
		t.Error("Redirect HTML missing canonical link")
	}

	if !strings.Contains(html, `href="https://example.com/blog/new-post/"`) {
		t.Error("Redirect HTML missing correct canonical href")
	}

	if !strings.Contains(html, `<a href="https://example.com/blog/new-post/">`) {
		t.Error("Redirect HTML missing fallback link")
	}
}

func TestGenerateSingleRedirect(t *testing.T) {
	tmpDir := t.TempDir()

	site := &config.Site{
		BaseURL:   "https://example.com",
		OutputDir: tmpDir,
	}

	b := &Builder{site: site}

	targetURL := "https://example.com/blog/new-post/"
	aliasPath := "/old-url"

	err := b.generateSingleRedirect(aliasPath, targetURL)
	if err != nil {
		t.Fatalf("generateSingleRedirect failed: %v", err)
	}

	// Check that the file was created
	expectedPath := filepath.Join(tmpDir, "old-url", "index.html")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Redirect file not created at %s", expectedPath)
	}

	// Read and verify the content
	content, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("Failed to read redirect file: %v", err)
	}

	if !strings.Contains(string(content), targetURL) {
		t.Error("Redirect file doesn't contain target URL")
	}
}

func TestGenerateSingleRedirect_WithoutLeadingSlash(t *testing.T) {
	tmpDir := t.TempDir()

	site := &config.Site{
		BaseURL:   "https://example.com",
		OutputDir: tmpDir,
	}

	b := &Builder{site: site}

	targetURL := "https://example.com/blog/new-post/"
	aliasPath := "old-url" // No leading slash

	err := b.generateSingleRedirect(aliasPath, targetURL)
	if err != nil {
		t.Fatalf("generateSingleRedirect failed: %v", err)
	}

	// Check that the file was created (should handle missing leading slash)
	expectedPath := filepath.Join(tmpDir, "old-url", "index.html")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Errorf("Redirect file not created at %s", expectedPath)
	}
}

func TestGenerateAliasRedirects(t *testing.T) {
	tmpDir := t.TempDir()

	site := &config.Site{
		BaseURL:   "https://example.com",
		OutputDir: tmpDir,
	}

	p := parser.NewParser(site)
	b := &Builder{
		site:   site,
		parser: p,
	}

	// Create a test node with aliases
	node := &content.Node{
		Path:      "blog/new-post.md",
		Slug:      "blog/new-post",
		Permalink: "/blog/new-post/",
		URL:       "https://example.com/blog/new-post/",
		Aliases:   []string{"/old-url", "/another-old-url", "/legacy/path"},
	}

	err := b.generateAliasRedirects(node)
	if err != nil {
		t.Fatalf("generateAliasRedirects failed: %v", err)
	}

	// Check that all alias files were created
	expectedPaths := []string{
		filepath.Join(tmpDir, "old-url", "index.html"),
		filepath.Join(tmpDir, "another-old-url", "index.html"),
		filepath.Join(tmpDir, "legacy", "path", "index.html"),
	}

	for _, path := range expectedPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Redirect file not created at %s", path)
		}

		// Verify content
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("Failed to read redirect file %s: %v", path, err)
		}

		if !strings.Contains(string(content), node.URL) {
			t.Errorf("Redirect file at %s doesn't contain target URL %s", path, node.URL)
		}
	}
}

func TestGenerateAliasRedirects_NoAliases(t *testing.T) {
	tmpDir := t.TempDir()

	site := &config.Site{
		BaseURL:   "https://example.com",
		OutputDir: tmpDir,
	}

	p := parser.NewParser(site)
	b := &Builder{
		site:   site,
		parser: p,
	}

	// Create a test node without aliases
	node := &content.Node{
		Path:      "blog/post.md",
		Slug:      "blog/post",
		Permalink: "/blog/post/",
		URL:       "https://example.com/blog/post/",
		Aliases:   []string{},
	}

	err := b.generateAliasRedirects(node)
	if err != nil {
		t.Fatalf("generateAliasRedirects failed: %v", err)
	}

	// Verify no redirect files were created
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read temp dir: %v", err)
	}

	if len(entries) > 0 {
		t.Errorf("Expected no files to be created, but found %d entries", len(entries))
	}
}

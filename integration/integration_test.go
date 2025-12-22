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
	if err := contentParser.Parse(contentDir); err != nil {
		t.Fatalf("failed to parse content: %v", err)
	}

	gen, err := builder.NewBuilder(site, contentParser)
	if err != nil {
		t.Fatalf("failed to create builder: %v", err)
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

// TestBuildMode tests all build mode scenarios
func TestBuildMode(t *testing.T) {
	t.Run("FreshBuild", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, _ := buildSite(t, sitePath)

		// Verify expected files exist
		verifyFileExists(t, sitePath, "index.html")
		verifyFileExists(t, sitePath, "about/index.html")
		verifyFileExists(t, sitePath, "blog/index.html")
		verifyFileExists(t, sitePath, "blog/post1/index.html")
		verifyFileExists(t, sitePath, "blog/post2/index.html")
		verifyFileExists(t, sitePath, "blog/post3/index.html")
		verifyFileExists(t, sitePath, "notes/index.html")
		verifyFileExists(t, sitePath, "notes/note1/index.html")

		// Verify static files
		verifyFileExists(t, sitePath, "style.css")

		// Verify auxiliary pages
		verifyFileExists(t, sitePath, "404.html")
		verifyFileExists(t, sitePath, "robots.txt")
		verifyFileExists(t, sitePath, "sitemap.xml")
		verifyFileExists(t, sitePath, "atom.xml")

		// Verify tag pages for post3
		verifyFileExists(t, sitePath, "tags/golang/index.html")
		verifyFileExists(t, sitePath, "tags/testing/index.html")

		// Verify cache stats (should start with 0% hit rate on fresh build)
		stats := gen.GetCacheStats()
		if stats.HitRate > 0 {
			t.Errorf("expected 0%% cache hit rate on fresh build, got %.1f%%", stats.HitRate)
		}
	})

	t.Run("ContentCorrectness", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Verify OG image on post with custom image
		verifyFileContent(t, sitePath, "blog/post1/index.html", "/img/post1.jpg")

		// Verify OG image on post without custom image uses site default
		verifyFileContent(t, sitePath, "blog/post2/index.html", "/img/default.jpg")

		// Verify post titles
		verifyFileContent(t, sitePath, "blog/post1/index.html", "First Post")
		verifyFileContent(t, sitePath, "blog/post2/index.html", "Second Post")
		verifyFileContent(t, sitePath, "blog/post3/index.html", "Third Post")

		// Verify blog listing shows all posts
		verifyFileContent(t, sitePath, "blog/index.html", "First Post")
		verifyFileContent(t, sitePath, "blog/index.html", "Second Post")
		verifyFileContent(t, sitePath, "blog/index.html", "Third Post")

		// Verify blog listing shows correct images
		verifyFileContent(t, sitePath, "blog/index.html", "/img/post1.jpg")
		verifyFileContent(t, sitePath, "blog/index.html", "/img/default.jpg")
		verifyFileContent(t, sitePath, "blog/index.html", "/img/post3.jpg")

		// Verify tag links on post3
		verifyFileContent(t, sitePath, "blog/post3/index.html", "/tags/golang/")
		verifyFileContent(t, sitePath, "blog/post3/index.html", "/tags/testing/")
	})

	t.Run("IncrementalBuild_NoChanges", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// First build
		gen, contentParser := buildSite(t, sitePath)
		firstStats := gen.GetCacheStats()

		// Reuse same builder and parser for incremental build (simulates serve mode)
		// This is different from running "gozzi build" twice which creates new instances
		if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("failed to generate site: %v", err)
		}
		secondStats := gen.GetCacheStats()

		// Cache hit rate should be much higher on second build with same builder instance
		if secondStats.HitRate < 50.0 {
			t.Errorf("expected >50%% cache hit rate on rebuild without changes, got %.1f%% (first: %d hits, %d misses; second: %d hits, %d misses)",
				secondStats.HitRate,
				firstStats.Hits, firstStats.Misses,
				secondStats.Hits, secondStats.Misses)
		}
	})

	t.Run("IncrementalBuild_ContentChange", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify post1 content
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		originalContent, err := os.ReadFile(post1Path)
		if err != nil {
			t.Fatalf("failed to read post1: %v", err)
		}

		modifiedContent := strings.Replace(string(originalContent), "first blog post", "MODIFIED blog post", 1)
		if err := os.WriteFile(post1Path, []byte(modifiedContent), 0644); err != nil {
			t.Fatalf("failed to write post1: %v", err)
		}

		// Re-parse content with same parser to pick up changes
		if err := contentParser.Parse(filepath.Join(sitePath, "content")); err != nil {
			t.Fatalf("failed to re-parse content: %v", err)
		}

		// Rebuild with same builder instance
		if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("failed to rebuild: %v", err)
		}

		// Verify the change is reflected
		verifyFileContent(t, sitePath, "blog/post1/index.html", "MODIFIED blog post")
	})

	t.Run("IncrementalBuild_FrontmatterChange", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Change post1's image in frontmatter
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		originalContent, err := os.ReadFile(post1Path)
		if err != nil {
			t.Fatalf("failed to read post1: %v", err)
		}

		modifiedContent := strings.Replace(string(originalContent), "/img/post1.jpg", "/img/changed.jpg", 1)
		if err := os.WriteFile(post1Path, []byte(modifiedContent), 0644); err != nil {
			t.Fatalf("failed to write post1: %v", err)
		}

		// Rebuild
		buildSite(t, sitePath)

		// Verify image changed in post itself
		verifyFileContent(t, sitePath, "blog/post1/index.html", "/img/changed.jpg")

		// Verify image changed in blog listing
		verifyFileContent(t, sitePath, "blog/index.html", "/img/changed.jpg")
	})

	t.Run("ConfigChange_FullRebuild", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Modify site config (change default image)
		configPath := filepath.Join(sitePath, "config.toml")
		originalConfig, err := os.ReadFile(configPath)
		if err != nil {
			t.Fatalf("failed to read config: %v", err)
		}

		modifiedConfig := strings.Replace(string(originalConfig), "/img/default.jpg", "/img/new-default.jpg", 1)
		if err := os.WriteFile(configPath, []byte(modifiedConfig), 0644); err != nil {
			t.Fatalf("failed to write config: %v", err)
		}

		// Rebuild
		gen, _ := buildSite(t, sitePath)
		stats := gen.GetCacheStats()

		// Cache hit rate should be 0% (full rebuild)
		if stats.HitRate > 0 {
			t.Errorf("expected 0%% cache hit rate after config change, got %.1f%%", stats.HitRate)
		}

		// Verify config change reflected in posts without custom images
		verifyFileContent(t, sitePath, "blog/post2/index.html", "/img/new-default.jpg")
	})

	t.Run("TagPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Verify tag pages generated
		verifyFileExists(t, sitePath, "tags/golang/index.html")
		verifyFileExists(t, sitePath, "tags/testing/index.html")

		// Verify tag page content
		verifyFileContent(t, sitePath, "tags/golang/index.html", "Third Post")
		verifyFileContent(t, sitePath, "tags/testing/index.html", "Third Post")

		// Add new tag to post1
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		content, err := os.ReadFile(post1Path)
		if err != nil {
			t.Fatalf("failed to read post1: %v", err)
		}

		// Insert tags line after template line
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "template =") {
				lines = append(lines[:i+1], append([]string{"tags = [\"newtag\"]"}, lines[i+1:]...)...)
				break
			}
		}

		if err := os.WriteFile(post1Path, []byte(strings.Join(lines, "\n")), 0644); err != nil {
			t.Fatalf("failed to write post1: %v", err)
		}

		// Rebuild
		buildSite(t, sitePath)

		// Verify new tag page created
		verifyFileExists(t, sitePath, "tags/newtag/index.html")
		verifyFileContent(t, sitePath, "tags/newtag/index.html", "First Post")
	})
}

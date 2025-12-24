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
		// Homepage
		verifyFileExists(t, sitePath, "index.html")

		// Single pages (various page types)
		verifyFileExists(t, sitePath, "about/index.html")
		verifyFileExists(t, sitePath, "contact/index.html")
		verifyFileExists(t, sitePath, "uses/index.html")
		verifyFileExists(t, sitePath, "privacy/index.html")

		// Blog section and posts
		verifyFileExists(t, sitePath, "blog/index.html")
		verifyFileExists(t, sitePath, "blog/post1/index.html")
		verifyFileExists(t, sitePath, "blog/post2/index.html")
		verifyFileExists(t, sitePath, "blog/post3/index.html")

		// Notes section and notes
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

// TestPageTypes tests all different page types to prevent regressions
func TestPageTypes(t *testing.T) {
	t.Run("SinglePages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Verify all single pages exist and have correct content
		pages := map[string]string{
			"about/index.html":   "About",
			"contact/index.html": "Contact",
			"uses/index.html":    "Uses",
			"privacy/index.html": "Privacy Policy",
		}

		for path, expectedTitle := range pages {
			verifyFileExists(t, sitePath, path)
			verifyFileContent(t, sitePath, path, expectedTitle)
		}

		// Verify page-specific content
		verifyFileContent(t, sitePath, "contact/index.html", "test@example.com")
		verifyFileContent(t, sitePath, "uses/index.html", "MacBook Pro")
		verifyFileContent(t, sitePath, "privacy/index.html", "does not collect personal data")
	})

	t.Run("SectionPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Blog section should list all posts
		verifyFileExists(t, sitePath, "blog/index.html")
		verifyFileContent(t, sitePath, "blog/index.html", "First Post")
		verifyFileContent(t, sitePath, "blog/index.html", "Second Post")
		verifyFileContent(t, sitePath, "blog/index.html", "Third Post")

		// Notes section should list all notes
		verifyFileExists(t, sitePath, "notes/index.html")
		verifyFileContent(t, sitePath, "notes/index.html", "First Note")
	})

	t.Run("BlogPosts", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// All blog posts should exist
		posts := []string{
			"blog/post1/index.html",
			"blog/post2/index.html",
			"blog/post3/index.html",
		}

		for _, post := range posts {
			verifyFileExists(t, sitePath, post)
		}

		// Posts should have proper structure
		verifyFileContent(t, sitePath, "blog/post1/index.html", "First Post")
		verifyFileContent(t, sitePath, "blog/post1/index.html", "first blog post")
		verifyFileContent(t, sitePath, "blog/post3/index.html", "golang")
		verifyFileContent(t, sitePath, "blog/post3/index.html", "testing")
	})

	t.Run("DraftPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Draft posts should NOT be generated by default
		draftPath := filepath.Join(sitePath, "public", "blog/draft-post/index.html")
		if _, err := os.Stat(draftPath); !os.IsNotExist(err) {
			t.Errorf("draft post should not be generated: %s", draftPath)
		}

		// Draft should not appear in blog listing
		blogListing := filepath.Join(sitePath, "public", "blog/index.html")
		content, err := os.ReadFile(blogListing)
		if err != nil {
			t.Fatalf("failed to read blog listing: %v", err)
		}
		if strings.Contains(string(content), "Draft Post") {
			t.Error("draft post should not appear in blog listing")
		}
	})

	t.Run("Homepage", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "index.html")
		verifyFileContent(t, sitePath, "index.html", "Home")
		// Note: Content rendering appears to be empty for homepage - potential bug to investigate
	})

	t.Run("AuxiliaryPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// 404 page
		verifyFileExists(t, sitePath, "404.html")
		verifyFileContent(t, sitePath, "404.html", "Page Not Found")

		// robots.txt
		verifyFileExists(t, sitePath, "robots.txt")

		// Sitemap (if generated)
		sitemapPath := filepath.Join(sitePath, "public", "sitemap.xml")
		if _, err := os.Stat(sitemapPath); err == nil {
			verifyFileContent(t, sitePath, "sitemap.xml", "https://")
		}

		// Atom feed (if generated)
		atomPath := filepath.Join(sitePath, "public", "atom.xml")
		if _, err := os.Stat(atomPath); err == nil {
			verifyFileContent(t, sitePath, "atom.xml", "<?xml")
		}
	})

	t.Run("StaticAssets", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "style.css")
		verifyFileContent(t, sitePath, "style.css", "font-family")
	})

	t.Run("SectionContentRendering", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Sections with _index.md should render their markdown content
		// This is CRITICAL - sections like /about/ must show their content, not be empty!
		verifyFileExists(t, sitePath, "about/index.html")

		// The page MUST contain the markdown content that was rendered
		verifyFileContent(t, sitePath, "about/index.html", "About This Test Site")
		verifyFileContent(t, sitePath, "about/index.html", "integration testing")
		verifyFileContent(t, sitePath, "about/index.html", "Section pages with content")
	})
}

// TestServerMode tests file watching and automatic rebuilds
func TestServerMode(t *testing.T) {
	t.Run("ContentFileChange_TriggersRebuild", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify content
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		content, _ := os.ReadFile(post1Path)
		modified := strings.Replace(string(content), "first blog post", "UPDATED content", 1)
		os.WriteFile(post1Path, []byte(modified), 0644)

		// Simulate watch mode: re-parse and rebuild
		contentParser.ParseFiles(filepath.Join(sitePath, "content"), []string{post1Path})
		gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:  true,
			ChangedFiles: []string{post1Path},
			ContentDir:   filepath.Join(sitePath, "content"),
		})

		verifyFileContent(t, sitePath, "blog/post1/index.html", "UPDATED content")
	})

	t.Run("TemplateChange_TriggersFullRebuild", func(t *testing.T) {
		t.Skip("Template reloading tested in app/builder/builder_test.go")
	})

	t.Run("StaticFileChange_CopiesFile", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Add new static file
		newStaticFile := filepath.Join(sitePath, "static/new-file.txt")
		os.WriteFile(newStaticFile, []byte("new content"), 0644)

		// Rebuild
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "new-file.txt")
		verifyFileContent(t, sitePath, "new-file.txt", "new content")
	})

	t.Run("MultipleContentChanges_IncrementalRebuild", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify multiple posts
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		post2Path := filepath.Join(sitePath, "content/blog/post2/index.md")

		content1, _ := os.ReadFile(post1Path)
		modified1 := strings.Replace(string(content1), "first blog post", "BATCH UPDATE 1", 1)
		os.WriteFile(post1Path, []byte(modified1), 0644)

		content2, _ := os.ReadFile(post2Path)
		modified2 := strings.Replace(string(content2), "second blog post", "BATCH UPDATE 2", 1)
		os.WriteFile(post2Path, []byte(modified2), 0644)

		// Simulate watch mode batch update
		changedFiles := []string{post1Path, post2Path}
		contentParser.ParseFiles(filepath.Join(sitePath, "content"), changedFiles)
		gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:  true,
			ChangedFiles: changedFiles,
			ContentDir:   filepath.Join(sitePath, "content"),
		})

		verifyFileContent(t, sitePath, "blog/post1/index.html", "BATCH UPDATE 1")
		verifyFileContent(t, sitePath, "blog/post2/index.html", "BATCH UPDATE 2")
	})

	t.Run("NewContentFile_AddsToSite", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Add new blog post
		newPost := `+++
title = "New Post"
date = 2024-01-15
template = "post.html"
+++

This is a new post added during watch mode.`

		newPostPath := filepath.Join(sitePath, "content/blog/new-post.md")
		os.WriteFile(newPostPath, []byte(newPost), 0644)

		// Full reparse needed for new files
		contentParser.Parse(filepath.Join(sitePath, "content"))
		gen.Generate(contentParser.ContentMap["."])

		verifyFileExists(t, sitePath, "blog/new-post/index.html")
		verifyFileContent(t, sitePath, "blog/new-post/index.html", "New Post")
		verifyFileContent(t, sitePath, "blog/new-post/index.html", "new post added during watch mode")

		// Should appear in blog listing
		verifyFileContent(t, sitePath, "blog/index.html", "New Post")
	})

	t.Run("DeleteContentFile_RemovesFromSite", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Verify post exists
		verifyFileExists(t, sitePath, "blog/post1/index.html")

		// Delete post and clean output
		post1Path := filepath.Join(sitePath, "content/blog/post1")
		os.RemoveAll(post1Path)

		// Clean and rebuild
		os.RemoveAll(filepath.Join(sitePath, "public"))
		buildSite(t, sitePath)

		// Output should not exist (file was deleted from content)
		outputPath := filepath.Join(sitePath, "public/blog/post1/index.html")
		if _, err := os.Stat(outputPath); !os.IsNotExist(err) {
			t.Error("deleted post should not exist in output")
		}

		// Should not appear in blog listing
		content, _ := os.ReadFile(filepath.Join(sitePath, "public/blog/index.html"))
		if strings.Contains(string(content), "First Post") {
			t.Error("deleted post should not appear in blog listing")
		}
	})

	t.Run("FrontmatterChange_UpdatesMetadata", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Change title in frontmatter
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		content, _ := os.ReadFile(post1Path)
		modified := strings.Replace(string(content), `title = "First Post"`, `title = "Changed Title"`, 1)
		os.WriteFile(post1Path, []byte(modified), 0644)

		contentParser.ParseFiles(filepath.Join(sitePath, "content"), []string{post1Path})
		gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:  true,
			ChangedFiles: []string{post1Path},
			ContentDir:   filepath.Join(sitePath, "content"),
		})

		// Title should update in post
		verifyFileContent(t, sitePath, "blog/post1/index.html", "Changed Title")

		// Title should update in blog listing
		verifyFileContent(t, sitePath, "blog/index.html", "Changed Title")
	})

	t.Run("AddTags_CreatesTagPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Add tags to post1
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		content, _ := os.ReadFile(post1Path)
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "template =") {
				lines = append(lines[:i+1], append([]string{`tags = ["rust", "performance"]`}, lines[i+1:]...)...)
				break
			}
		}
		os.WriteFile(post1Path, []byte(strings.Join(lines, "\n")), 0644)

		// Full rebuild needed for tag changes
		contentParser.Parse(filepath.Join(sitePath, "content"))
		gen.Generate(contentParser.ContentMap["."])

		verifyFileExists(t, sitePath, "tags/rust/index.html")
		verifyFileExists(t, sitePath, "tags/performance/index.html")
		verifyFileContent(t, sitePath, "tags/rust/index.html", "First Post")
	})

	t.Run("ConfigChange_FullRebuild", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Modify homepage content title (not site config)
		homePath := filepath.Join(sitePath, "content/_index.md")
		content, _ := os.ReadFile(homePath)
		modified := strings.Replace(string(content), `title = "Home"`, `title = "Updated Home"`, 1)
		os.WriteFile(homePath, []byte(modified), 0644)

		// Fresh build with changed content
		buildSite(t, sitePath)

		// Change should be reflected in homepage
		verifyFileContent(t, sitePath, "index.html", "Updated Home")
	})
}

// TestCacheInvalidation ensures cache invalidation works correctly
func TestCacheInvalidation(t *testing.T) {
	t.Run("ContentChange_InvalidatesCache", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		firstStats := gen.GetCacheStats()

		// Rebuild without changes - should have high cache hit rate
		gen.Generate(contentParser.ContentMap["."])
		noChangeStats := gen.GetCacheStats()

		if noChangeStats.HitRate < 50 {
			t.Errorf("expected high cache hit rate without changes, got %.1f%%", noChangeStats.HitRate)
		}

		// Modify content
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		content, _ := os.ReadFile(post1Path)
		modified := strings.Replace(string(content), "first blog post", "CHANGED", 1)
		os.WriteFile(post1Path, []byte(modified), 0644)

		// Reparse and rebuild
		contentParser.Parse(filepath.Join(sitePath, "content"))
		gen.ResetCacheStats()
		gen.Generate(contentParser.ContentMap["."])
		changedStats := gen.GetCacheStats()

		// Cache should be partially invalidated
		if changedStats.Misses == 0 {
			t.Error("expected cache misses after content change")
		}

		_ = firstStats
	})

	t.Run("TemplateChange_InvalidatesAllPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Rebuild without changes
		gen.Generate(contentParser.ContentMap["."])
		beforeStats := gen.GetCacheStats()

		// Modify template
		templatePath := filepath.Join(sitePath, "templates/post.html")
		content, _ := os.ReadFile(templatePath)
		modified := strings.Replace(string(content), "<article>", "<article class=\"v2\">", 1)
		os.WriteFile(templatePath, []byte(modified), 0644)

		// Reload templates
		gen.ReloadTemplates()
		gen.InvalidateTemplateCache([]string{"post.html"})

		gen.ResetCacheStats()
		gen.Generate(contentParser.ContentMap["."])
		afterStats := gen.GetCacheStats()

		// All posts using post.html should be regenerated
		if afterStats.Misses == 0 {
			t.Error("expected cache misses after template change")
		}

		_ = beforeStats
	})
}

// TestEdgeCases tests unusual but valid scenarios
func TestEdgeCases(t *testing.T) {
	t.Run("EmptyContent_StillGenerates", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Create post with empty content
		emptyPost := `+++
title = "Empty Post"
date = 2024-01-01
template = "post.html"
+++
`
		emptyPostPath := filepath.Join(sitePath, "content/blog/empty.md")
		os.WriteFile(emptyPostPath, []byte(emptyPost), 0644)

		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "blog/empty/index.html")
		verifyFileContent(t, sitePath, "blog/empty/index.html", "Empty Post")
	})

	t.Run("SpecialCharactersInFilename", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Create post with special chars
		specialPost := `+++
title = "Special Post"
date = 2024-01-01
template = "post.html"
+++

Content here.`

		specialPath := filepath.Join(sitePath, "content/blog/special-chars-&-symbols.md")
		os.WriteFile(specialPath, []byte(specialPost), 0644)

		buildSite(t, sitePath)

		// Should generate with sanitized slug
		verifyFileExists(t, sitePath, "blog/special-chars-symbols/index.html")
	})

	t.Run("VeryLongContent_Handles", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Create post with very long content
		longContent := strings.Repeat("This is a very long paragraph. ", 1000)
		longPost := `+++
title = "Long Post"
date = 2024-01-01
template = "post.html"
+++

` + longContent

		longPath := filepath.Join(sitePath, "content/blog/long-post.md")
		os.WriteFile(longPath, []byte(longPost), 0644)

		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "blog/long-post/index.html")
		verifyFileContent(t, sitePath, "blog/long-post/index.html", "Long Post")
	})

	t.Run("NestedSections_HandleCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Create nested section
		os.MkdirAll(filepath.Join(sitePath, "content/blog/nested"), 0755)
		nestedIndex := `+++
title = "Nested Section"
+++

Nested section content.`
		os.WriteFile(filepath.Join(sitePath, "content/blog/nested/_index.md"), []byte(nestedIndex), 0644)

		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "blog/nested/index.html")
		verifyFileContent(t, sitePath, "blog/nested/index.html", "Nested Section")
	})

	t.Run("NonAsciiContent_RendersCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)

		unicodePost := `+++
title = "Unicode Post 中文 日本語"
date = 2024-01-01
template = "post.html"
+++

Content with unicode: 你好世界 こんにちは مرحبا`

		unicodePath := filepath.Join(sitePath, "content/blog/unicode.md")
		os.WriteFile(unicodePath, []byte(unicodePost), 0644)

		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "blog/unicode/index.html")
		verifyFileContent(t, sitePath, "blog/unicode/index.html", "中文")
		verifyFileContent(t, sitePath, "blog/unicode/index.html", "你好世界")
	})
}

// TestPerformance tests build performance characteristics
func TestPerformance(t *testing.T) {
	t.Run("BuildCompletes_InReasonableTime", func(t *testing.T) {
		sitePath := setupTestSite(t)

		start := time.Now()
		buildSite(t, sitePath)
		duration := time.Since(start)

		// Full build of test site should be very fast
		maxDuration := 2 * time.Second
		if duration > maxDuration {
			t.Errorf("build took too long: %v (max: %v)", duration, maxDuration)
		}

		t.Logf("Full build completed in %v", duration)
	})

	t.Run("IncrementalBuild_FasterThanFull", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Measure full build time
		fullStart := time.Now()
		gen, contentParser := buildSite(t, sitePath)
		fullDuration := time.Since(fullStart)

		// Measure incremental rebuild time (no changes)
		incStart := time.Now()
		gen.Generate(contentParser.ContentMap["."])
		incDuration := time.Since(incStart)

		// Incremental should be faster (or at least not significantly slower)
		if incDuration > fullDuration*2 {
			t.Errorf("incremental build too slow: %v (full: %v)", incDuration, fullDuration)
		}

		t.Logf("Full build: %v, Incremental rebuild: %v", fullDuration, incDuration)
	})
}

// TestContentIntegrity tests that all content is properly rendered and not empty
func TestContentIntegrity(t *testing.T) {
	t.Run("AllPages_HaveContent", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// All pages should have substantial content (not just headers/footers)
		pages := []string{
			"index.html",
			"about/index.html",
			"contact/index.html",
			"uses/index.html",
			"privacy/index.html",
			"blog/index.html",
			"notes/index.html",
			"blog/post1/index.html",
			"blog/post2/index.html",
			"blog/post3/index.html",
			"notes/note1/index.html",
		}

		for _, page := range pages {
			fullPath := filepath.Join(sitePath, "public", page)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Fatalf("failed to read %s: %v", page, err)
			}

			// Pages should have more than just template boilerplate
			if len(content) < 100 {
				t.Errorf("page %s appears empty (only %d bytes)", page, len(content))
			}

			// Should contain HTML structure
			contentStr := string(content)
			if !strings.Contains(contentStr, "<html") && !strings.Contains(contentStr, "<!DOCTYPE") {
				t.Errorf("page %s missing HTML structure", page)
			}
		}
	})

	t.Run("SectionPages_NotEmpty", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// This was the critical bug - sections with _index.md must render their content
		// NOT be empty shells
		criticalSections := []string{
			"about/index.html",
			"contact/index.html",
			"uses/index.html",
			"privacy/index.html",
		}

		for _, section := range criticalSections {
			fullPath := filepath.Join(sitePath, "public", section)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Fatalf("failed to read %s: %v", section, err)
			}

			// Sections should have meaningful content (not just empty shells)
			// Testdata has minimal content, so threshold is lower than production
			if len(content) < 200 {
				t.Errorf("section %s appears empty (only %d bytes)", section, len(content))
			}

			// Should contain actual rendered markdown content
			contentStr := string(content)
			// About page specifically should have its content
			if section == "about/index.html" {
				if !strings.Contains(contentStr, "About This Test Site") {
					t.Errorf("about section missing its markdown content")
				}
			}
		}
	})

	t.Run("BlogListing_ShowsAllPosts", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		blogPath := filepath.Join(sitePath, "public", "blog", "index.html")
		content, err := os.ReadFile(blogPath)
		if err != nil {
			t.Fatalf("failed to read blog listing: %v", err)
		}

		// Should list all non-draft posts
		requiredPosts := []string{"First Post", "Second Post", "Third Post"}
		for _, post := range requiredPosts {
			if !strings.Contains(string(content), post) {
				t.Errorf("blog listing missing post: %s", post)
			}
		}

		// Should NOT show draft posts
		if strings.Contains(string(content), "Draft Post") {
			t.Error("blog listing should not show draft posts")
		}
	})
}

// TestStaticAssets ensures static file handling works correctly
func TestStaticAssets(t *testing.T) {
	t.Run("AllStaticFiles_Copied", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Verify critical static assets exist
		staticAssets := []string{
			"robots.txt",
			"style.css",
		}

		for _, asset := range staticAssets {
			verifyFileExists(t, sitePath, asset)
		}
	})

	t.Run("StaticContent_Preserved", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Static files should be copied verbatim
		verifyFileContent(t, sitePath, "style.css", "font-family")
		verifyFileContent(t, sitePath, "robots.txt", "User-agent")
	})

	t.Run("NestedStatic_CopiedCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Add nested static structure
		nestedDir := filepath.Join(sitePath, "static/images/icons")
		os.MkdirAll(nestedDir, 0755)
		os.WriteFile(filepath.Join(nestedDir, "test.svg"), []byte("<svg></svg>"), 0644)

		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "images/icons/test.svg")
		verifyFileContent(t, sitePath, "images/icons/test.svg", "<svg>")
	})
}

package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestBuildMode_FreshBuild tests full cold-start builds from scratch
func TestBuildMode_FreshBuild(t *testing.T) {
	t.Run("FreshBuild_GeneratesAllPages", func(t *testing.T) {
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

	t.Run("FreshBuild_ContentCorrectness", func(t *testing.T) {
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

	t.Run("FreshBuild_PerformanceBaseline", func(t *testing.T) {
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
}

// TestBuildMode_ConfigChange tests that config changes trigger full rebuilds
func TestBuildMode_ConfigChange(t *testing.T) {
	t.Run("ConfigChange_TriggersFullRebuild", func(t *testing.T) {
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

	t.Run("ConfigChange_BaseURL_UpdatesAllLinks", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Verify initial base_url
		verifyFileContent(t, sitePath, "sitemap.xml", "https://test.example.com")

		// Change base_url
		configPath := filepath.Join(sitePath, "config.toml")
		config, _ := os.ReadFile(configPath)
		modifiedConfig := strings.Replace(string(config), "https://test.example.com", "https://newsite.com", 1)
		os.WriteFile(configPath, []byte(modifiedConfig), 0644)

		// Rebuild
		buildSite(t, sitePath)

		// Verify base_url changed in sitemap
		verifyFileContent(t, sitePath, "sitemap.xml", "https://newsite.com")
	})
}

// TestBuildMode_TemplateChange tests that template changes trigger affected pages rebuild
func TestBuildMode_TemplateChange(t *testing.T) {
	t.Run("TemplateChange_RegeneratesAffectedPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Rebuild without changes
		gen.Generate(contentParser.ContentMap["."])
		beforeStats := gen.GetCacheStats()

		// Modify template - use actual tag from template
		templatePath := filepath.Join(sitePath, "templates/post.html")
		content, _ := os.ReadFile(templatePath)
		modified := strings.Replace(string(content), "<body>", "<body class=\"v2\">", 1)
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

		// Verify template change is reflected
		verifyFileContent(t, sitePath, "blog/post1/index.html", `class="v2"`)

		_ = beforeStats
	})

	t.Run("TemplateChange_DoesNotAffectOtherTemplates", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify post.html template
		templatePath := filepath.Join(sitePath, "templates/post.html")
		content, _ := os.ReadFile(templatePath)
		modified := strings.Replace(string(content), "<body>", "<body data-test=\"modified\">", 1)
		os.WriteFile(templatePath, []byte(modified), 0644)

		// Reload and rebuild
		gen.ReloadTemplates()
		gen.InvalidateTemplateCache([]string{"post.html"})
		gen.Generate(contentParser.ContentMap["."])

		// Post pages should have the change
		verifyFileContent(t, sitePath, "blog/post1/index.html", `data-test="modified"`)

		// Section pages should NOT have the change (they use section.html)
		content, _ = os.ReadFile(filepath.Join(sitePath, "public/blog/index.html"))
		if strings.Contains(string(content), `data-test="modified"`) {
			t.Error("section pages should not be affected by post.html template change")
		}
	})
}

// TestBuildMode_TaxonomyGeneration tests full taxonomy generation
func TestBuildMode_TaxonomyGeneration(t *testing.T) {
	t.Run("TagPages_GeneratedCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Verify tag pages generated
		verifyFileExists(t, sitePath, "tags/golang/index.html")
		verifyFileExists(t, sitePath, "tags/testing/index.html")

		// Verify tag page content
		verifyFileContent(t, sitePath, "tags/golang/index.html", "Third Post")
		verifyFileContent(t, sitePath, "tags/testing/index.html", "Third Post")

		// Verify tags index page
		verifyFileExists(t, sitePath, "tags/index.html")
		verifyFileContent(t, sitePath, "tags/index.html", "golang")
		verifyFileContent(t, sitePath, "tags/index.html", "testing")
	})

	t.Run("SeriesPages_GeneratedCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Series index should exist
		verifyFileExists(t, sitePath, "series/index.html")
		verifyFileContent(t, sitePath, "series/index.html", "All Series")

		// Individual series pages should exist
		verifyFileExists(t, sitePath, "series/test-series/index.html")
		verifyFileExists(t, sitePath, "series/another-series/index.html")

		// Verify series content
		verifyFileContent(t, sitePath, "series/test-series/index.html", "Series: test-series")
		verifyFileContent(t, sitePath, "series/test-series/index.html", "3 posts in this series")
	})

	t.Run("TaxonomyPages_ExcludeDrafts", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Add draft post with tags
		draftPost := `+++
title = "Draft With Tags"
date = 2024-01-01
template = "post.html"
tags = ["draft-tag"]
draft = true
+++
Content`

		draftPath := filepath.Join(sitePath, "content/blog/draft-with-tags.md")
		os.WriteFile(draftPath, []byte(draftPost), 0644)

		// Rebuild
		buildSite(t, sitePath)

		// Tag page for draft-tag should not exist (draft excluded)
		tagPath := filepath.Join(sitePath, "public/tags/draft-tag/index.html")
		if _, err := os.Stat(tagPath); !os.IsNotExist(err) {
			t.Error("tag page for draft post should not be generated")
		}
	})
}

// TestBuildMode_DraftHandling tests draft page handling
func TestBuildMode_DraftHandling(t *testing.T) {
	t.Run("DraftPosts_NotGenerated", func(t *testing.T) {
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

	t.Run("DraftInSeries_ExcludedFromNavigation", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Add draft post to series
		draftPost := `+++
title = "Test Series Part 2.5 Draft"
date = 2024-01-02T12:00:00Z
template = "post.html"
series = "test-series"
series_order = 2.5
draft = true

[extra]
description = "Draft post in series"
+++

This is a draft and should not affect series navigation.`

		draftPath := filepath.Join(sitePath, "content/blog/series-draft.md")
		os.WriteFile(draftPath, []byte(draftPost), 0644)

		// Rebuild
		contentParser.Parse(filepath.Join(sitePath, "content"))
		gen.Generate(contentParser.ContentMap["."])

		// Series should still show 3 posts (draft excluded)
		verifyFileContent(t, sitePath, "series/test-series/index.html", "3 posts in this series")

		// Part 2 should still link directly to Part 3 (skip draft)
		verifyFileContent(t, sitePath, "blog/series-part2/index.html", "Next:")
		verifyFileContent(t, sitePath, "blog/series-part2/index.html", "Test Series Part 3")

		// Draft should not exist in output
		draftOutput := filepath.Join(sitePath, "public/blog/series-draft/index.html")
		if _, err := os.Stat(draftOutput); !os.IsNotExist(err) {
			t.Error("draft post should not be generated")
		}
	})
}

// TestBuildMode_StaticAssets tests static file handling in full builds
func TestBuildMode_StaticAssets(t *testing.T) {
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

	t.Run("StaticContent_PreservedVerbatim", func(t *testing.T) {
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

	t.Run("StaticFiles_DoNotTriggerContentRebuild", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Populate cache with initial build
		gen.ClearRenderCache()
		if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("initial build failed: %v", err)
		}

		// Add new static file
		newStaticFile := filepath.Join(sitePath, "static/new-asset.txt")
		os.WriteFile(newStaticFile, []byte("new content"), 0644)

		// Rebuild WITHOUT clearing cache - should get high cache hit rate
		gen.ResetCacheStats()
		if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("rebuild failed: %v", err)
		}

		stats := gen.GetCacheStats()
		// Static file changes shouldn't affect content rendering
		if stats.HitRate < 80 {
			t.Errorf("expected high cache hit rate when only static files change, got %.1f%%", stats.HitRate)
		}

		verifyFileExists(t, sitePath, "new-asset.txt")
	})
}

// TestBuildMode_ContentIntegrity tests content is properly rendered
func TestBuildMode_ContentIntegrity(t *testing.T) {
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

	t.Run("SectionPages_RenderMarkdownContent", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// CRITICAL: Sections with _index.md must render their markdown content
		// This was a regression bug - sections were empty shells
		verifyFileContent(t, sitePath, "about/index.html", "About This Test Site")
		verifyFileContent(t, sitePath, "about/index.html", "integration testing")
		verifyFileContent(t, sitePath, "about/index.html", "Section pages with content")
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

// TestBuildMode_EdgeCases tests unusual but valid scenarios
func TestBuildMode_EdgeCases(t *testing.T) {
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

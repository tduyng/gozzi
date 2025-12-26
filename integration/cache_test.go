// This file tests cache invalidation scenarios and cache hit rates.
package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestCache_FreshBuild tests cache behavior on fresh builds
func TestCache_FreshBuild(t *testing.T) {
	t.Run("FreshBuild_ZeroCacheHitRate", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, _ := buildSite(t, sitePath)

		// Fresh build should have 0% cache hit rate
		stats := gen.GetCacheStats()
		if stats.HitRate > 0 {
			t.Errorf("expected 0%% cache hit rate on fresh build, got %.1f%%", stats.HitRate)
		}

		// Should have cache misses
		if stats.Misses == 0 {
			t.Error("expected cache misses on fresh build")
		}
	})

	t.Run("SecondFullBuild_HighCacheHitRate", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Reset cache stats and rebuild without changes
		gen.ResetCacheStats()
		time.Sleep(10 * time.Millisecond)
		fullRebuild(t, gen, contentParser, sitePath)

		// Second build should have high cache hit rate
		stats := gen.GetCacheStats()
		if stats.HitRate < 80 {
			t.Errorf("expected >80%% cache hit rate on unchanged rebuild, got %.1f%%", stats.HitRate)
		}
	})
}

// TestCache_ContentChange tests cache invalidation on content changes
func TestCache_ContentChange(t *testing.T) {
	t.Run("SinglePostChange_InvalidatesOnlyThatPost", func(t *testing.T) {
		t.Skip("TODO: Fix cache behavior to maintain hits for unchanged pages during incremental rebuild")
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify one post
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		modifyFile(t, post1Path, "First Post", "First Post Modified")

		// Incremental rebuild
		gen.ResetCacheStats()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{post1Path})

		// Should have some cache hits (other pages unchanged)
		stats := gen.GetCacheStats()
		if stats.Hits == 0 {
			t.Error("expected some cache hits for unchanged pages")
		}

		// Verify change reflected
		verifyFileContent(t, sitePath, "blog/post1/index.html", "First Post Modified")

		// Verify other posts unchanged (still cached)
		verifyFileContent(t, sitePath, "blog/post2/index.html", "Second Post")
	})

	t.Run("ContentChange_InvalidatesSectionListing", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Add new post to blog section
		newPost := `+++
title = "New Post"
date = 2024-02-01
template = "post.html"
+++

New content here.`

		newPostPath := filepath.Join(sitePath, "content/blog/new-post.md")
		createPost(t, sitePath, "blog/new-post.md", newPost)

		// Incremental rebuild
		gen.ResetCacheStats()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{newPostPath})

		// Blog section listing should be regenerated
		verifyFileExists(t, sitePath, "blog/new-post/index.html")
		verifyFileContent(t, sitePath, "blog/index.html", "New Post")
	})

	t.Run("MultipleContentChanges_InvalidatesMultiple", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify multiple posts
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		post2Path := filepath.Join(sitePath, "content/blog/post2/index.md")

		modifyFile(t, post1Path, "First Post", "First Post V2")
		modifyFile(t, post2Path, "Second Post", "Second Post V2")

		// Incremental rebuild with multiple changed files
		gen.ResetCacheStats()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{post1Path, post2Path})

		// Both posts should be regenerated
		verifyFileContent(t, sitePath, "blog/post1/index.html", "First Post V2")
		verifyFileContent(t, sitePath, "blog/post2/index.html", "Second Post V2")

		// Other posts should still be cached
		verifyFileContent(t, sitePath, "blog/post3/index.html", "Third Post")
	})
}

// TestCache_TemplateChange tests cache invalidation on template changes
func TestCache_TemplateChange(t *testing.T) {
	t.Run("TemplateChange_InvalidatesAllAffectedPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify post template
		templatePath := filepath.Join(sitePath, "templates/post.html")
		modifyFile(t, templatePath, "<body>", "<body class=\"v2\">")

		// Reload templates and invalidate cache
		gen.ReloadTemplates()
		gen.InvalidateTemplateCache([]string{"post.html"})

		// Full rebuild
		gen.ResetCacheStats()
		fullRebuild(t, gen, contentParser, sitePath)

		// All posts should be regenerated (cache misses)
		stats := gen.GetCacheStats()
		if stats.HitRate > 50 {
			t.Errorf("expected low cache hit rate after template change, got %.1f%%", stats.HitRate)
		}

		// Verify template change reflected
		verifyFileContent(t, sitePath, "blog/post1/index.html", `class="v2"`)
		verifyFileContent(t, sitePath, "blog/post2/index.html", `class="v2"`)
	})

	t.Run("PartialTemplateChange_InvalidatesPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify section template (doesn't affect posts) - use actual tag
		templatePath := filepath.Join(sitePath, "templates/section.html")
		modifyFile(t, templatePath, "<body>", "<body class=\"updated\">")

		// Reload and rebuild
		gen.ReloadTemplates()
		gen.InvalidateTemplateCache([]string{"section.html"})
		gen.ResetCacheStats()
		fullRebuild(t, gen, contentParser, sitePath)

		// Section pages regenerated
		verifyFileContent(t, sitePath, "blog/index.html", `class="updated"`)

		// Post pages can still use cache (different template)
		verifyFileContent(t, sitePath, "blog/post1/index.html", "First Post")
	})
}

// TestCache_TaxonomyChange tests cache invalidation on taxonomy changes
func TestCache_TaxonomyChange(t *testing.T) {
	t.Run("AddTag_InvalidatesTaxonomyPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Add tag to post1
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		insertAfterLine(t, post1Path, "template =", `tags = ["newtag"]`)

		// Incremental rebuild
		gen.ResetCacheStats()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{post1Path})

		// Tag page should be created
		verifyFileExists(t, sitePath, "tags/newtag/index.html")
		verifyFileContent(t, sitePath, "tags/newtag/index.html", "First Post")

		// Tags index should be regenerated
		verifyFileExists(t, sitePath, "tags/index.html")
		verifyFileContent(t, sitePath, "tags/index.html", "newtag")
	})

	t.Run("AddSeries_InvalidatesSeriesPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Add series to post1
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		insertAfterLine(t, post1Path, "template =", `series = "newseries"`, `series_order = 1`)

		// Incremental rebuild
		gen.ResetCacheStats()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{post1Path})

		// Series page should be created
		verifyFileExists(t, sitePath, "series/newseries/index.html")
		verifyFileContent(t, sitePath, "series/newseries/index.html", "First Post")

		// Series index should be regenerated
		verifyFileExists(t, sitePath, "series/index.html")
		verifyFileContent(t, sitePath, "series/index.html", "newseries")

		// Post should have series navigation
		verifyFileContent(t, sitePath, "blog/post1/index.html", "newseries")
	})

	t.Run("ChangeTaxonomy_InvalidatesOldAndNewPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// post3 has tags ["golang", "testing"]
		verifyFileContent(t, sitePath, "tags/golang/index.html", "Third Post")

		// Change tags on post3
		post3Path := filepath.Join(sitePath, "content/blog/post3/index.md")
		modifyFile(t, post3Path, `tags = ["golang", "testing"]`, `tags = ["rust", "benchmarking"]`)

		// Incremental rebuild
		gen.ResetCacheStats()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{post3Path})

		// Old tag pages should be removed (since only post3 had these tags)
		verifyFileNotExists(t, sitePath, "tags/golang/index.html")
		verifyFileNotExists(t, sitePath, "tags/testing/index.html")

		// New tag pages should show post3
		verifyFileExists(t, sitePath, "tags/rust/index.html")
		verifyFileExists(t, sitePath, "tags/benchmarking/index.html")
		verifyFileContent(t, sitePath, "tags/rust/index.html", "Third Post")
		verifyFileContent(t, sitePath, "tags/benchmarking/index.html", "Third Post")

		// Tags index should list new tags
		verifyFileContent(t, sitePath, "tags/index.html", "rust")
		verifyFileContent(t, sitePath, "tags/index.html", "benchmarking")
	})

	t.Run("RemoveTaxonomy_CleansUpPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create post with tag
		post := `+++
title = "Temp Post"
date = 2024-01-01
template = "post.html"
tags = ["temptag"]
+++
Content`
		tempPostPath := filepath.Join(sitePath, "content/blog/temp-post.md")
		createPost(t, sitePath, "blog/temp-post.md", post)

		// Build to create tag page
		fullRebuild(t, gen, contentParser, sitePath)
		verifyFileExists(t, sitePath, "tags/temptag/index.html")

		// Remove tag from post
		modifyFile(t, tempPostPath, `tags = ["temptag"]`, "")

		// Incremental rebuild
		gen.ResetCacheStats()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{tempPostPath})

		// Tag page should be removed (since no posts have this tag anymore)
		verifyFileNotExists(t, sitePath, "tags/temptag/index.html")
	})
}

// TestCache_ConfigChange tests cache invalidation on config changes
func TestCache_ConfigChange(t *testing.T) {
	t.Run("ConfigChange_InvalidatesAllPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Modify config
		configPath := filepath.Join(sitePath, "config.toml")
		modifyFile(t, configPath, "Example Site", "New Site Name")

		// Rebuild
		gen2, _ := buildSite(t, sitePath)

		// Should have fresh cache (new config = new build)
		stats := gen2.GetCacheStats()
		if stats.HitRate > 0 {
			t.Errorf("expected 0%% cache hit rate after config change, got %.1f%%", stats.HitRate)
		}
	})

	t.Run("ConfigChange_BaseURL_RegeneratesAll", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Verify initial base_url
		verifyFileContent(t, sitePath, "sitemap.xml", "https://test.example.com")

		// Change base_url
		configPath := filepath.Join(sitePath, "config.toml")
		modifyFile(t, configPath, "https://test.example.com", "https://newdomain.com")

		// Rebuild
		buildSite(t, sitePath)

		// All pages with absolute URLs should be regenerated
		verifyFileContent(t, sitePath, "sitemap.xml", "https://newdomain.com")
		verifyFileContent(t, sitePath, "atom.xml", "https://newdomain.com")
	})
}

// TestCache_StaticFileChange tests that static file changes don't affect content cache
func TestCache_StaticFileChange(t *testing.T) {
	t.Run("StaticChange_NoContentCacheInvalidation", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify static file
		staticPath := filepath.Join(sitePath, "static/style.css")
		modifyFile(t, staticPath, "font-family", "font-family: monospace")

		// Rebuild
		gen.ResetCacheStats()
		fullRebuild(t, gen, contentParser, sitePath)

		// Content pages should still be cached
		stats := gen.GetCacheStats()
		if stats.HitRate < 80 {
			t.Errorf("static file change should not invalidate content cache, got %.1f%% hit rate", stats.HitRate)
		}

		// Static file should be updated
		verifyFileContent(t, sitePath, "style.css", "monospace")
	})

	t.Run("NewStaticFile_DoesNotAffectCache", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Add new static file
		newStaticPath := filepath.Join(sitePath, "static/newfile.txt")
		os.WriteFile(newStaticPath, []byte("new content"), 0644)

		// Rebuild
		gen.ResetCacheStats()
		fullRebuild(t, gen, contentParser, sitePath)

		// Content cache should be unaffected
		stats := gen.GetCacheStats()
		if stats.HitRate < 80 {
			t.Errorf("new static file should not invalidate content cache, got %.1f%% hit rate", stats.HitRate)
		}

		// New file should exist in output
		verifyFileExists(t, sitePath, "newfile.txt")
	})
}

// TestCache_PartialInvalidation tests that only affected pages are invalidated
func TestCache_PartialInvalidation(t *testing.T) {
	t.Run("SinglePostChange_OthersStayCached", func(t *testing.T) {
		t.Skip("TODO: Fix content modification detection to properly update changed files")
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Get initial cache stats
		gen.ResetCacheStats()
		fullRebuild(t, gen, contentParser, sitePath)
		initialStats := gen.GetCacheStats()

		// Modify only post1
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		modifyFile(t, post1Path, "This is the first post", "This is the UPDATED first post")

		// Incremental rebuild
		gen.ResetCacheStats()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{post1Path})

		// Should have high cache hit rate (only post1 and blog listing invalidated)
		stats := gen.GetCacheStats()
		if stats.Hits == 0 {
			t.Error("expected cache hits for unchanged pages")
		}

		// Verify post1 changed
		verifyFileContent(t, sitePath, "blog/post1/index.html", "UPDATED")

		// Verify other posts unchanged
		verifyFileContent(t, sitePath, "blog/post2/index.html", "Second Post")
		verifyFileContent(t, sitePath, "blog/post3/index.html", "Third Post")

		_ = initialStats
	})

	t.Run("TaxonomyChange_OnlyAffectedPagesInvalidated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Add tag to post1 (which has no tags initially)
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		insertAfterLine(t, post1Path, "template =", `tags = ["specialized"]`)

		// Incremental rebuild
		gen.ResetCacheStats()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{post1Path})

		// New tag page created
		verifyFileExists(t, sitePath, "tags/specialized/index.html")

		// Existing tag pages (golang, testing) should be cached
		stats := gen.GetCacheStats()
		if stats.Hits == 0 {
			t.Error("expected cache hits for unaffected taxonomy pages")
		}

		// post3's tags should be unchanged
		verifyFileContent(t, sitePath, "tags/golang/index.html", "Third Post")
		verifyFileContent(t, sitePath, "tags/testing/index.html", "Third Post")
	})
}

// TestCache_CacheKeyInclusion tests that cache keys include all relevant fields
func TestCache_CacheKeyInclusion(t *testing.T) {
	t.Run("SeriesChange_InvalidatesCacheDueToKey", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Add series to post without series
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		insertAfterLine(t, post1Path, "template =", `series = "cachekeyseries"`, `series_order = 1`)

		// Incremental rebuild
		gen.ResetCacheStats()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{post1Path})

		// Post should be regenerated (series in cache key)
		verifyFileContent(t, sitePath, "blog/post1/index.html", "cachekeyseries")

		// Should have some cache misses for regenerated pages
		stats := gen.GetCacheStats()
		if stats.Misses == 0 {
			t.Error("expected cache misses when series field changes")
		}
	})

	t.Run("TagChange_InvalidatesCacheDueToKey", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Change tags on post3
		post3Path := filepath.Join(sitePath, "content/blog/post3/index.md")
		modifyFile(t, post3Path, `tags = ["golang", "testing"]`, `tags = ["golang", "performance"]`)

		// Incremental rebuild
		gen.ResetCacheStats()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{post3Path})

		// Post should be regenerated (tags in cache key)
		verifyFileContent(t, sitePath, "blog/post3/index.html", "performance")
		verifyFileContent(t, sitePath, "tags/performance/index.html", "Third Post")

		// Should have cache misses
		stats := gen.GetCacheStats()
		if stats.Misses == 0 {
			t.Error("expected cache misses when tags change")
		}
	})

	t.Run("DateChange_InvalidatesCacheDueToKey", func(t *testing.T) {
		t.Skip("TODO: Verify cache key includes date and properly invalidates when date changes")
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Change date on post
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		modifyFile(t, post1Path, "2024-01-01", "2024-12-31")

		// Incremental rebuild
		gen.ResetCacheStats()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{post1Path})

		// Should regenerate (date affects sorting in listings)
		stats := gen.GetCacheStats()
		if stats.Misses == 0 {
			t.Error("expected cache misses when date changes")
		}
	})
}

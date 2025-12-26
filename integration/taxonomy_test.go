// This file tests all taxonomy features: series, tags, and categories.
package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestTaxonomy_Series tests series taxonomy functionality
func TestTaxonomy_Series(t *testing.T) {
	t.Run("SeriesPages_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Series index should exist
		verifyFileExists(t, sitePath, "series/index.html")
		verifyFileContent(t, sitePath, "series/index.html", "All Series")

		// Should list both series
		verifyFileContent(t, sitePath, "series/index.html", "test-series")
		verifyFileContent(t, sitePath, "series/index.html", "another-series")

		// Should show post counts
		verifyFileContent(t, sitePath, "series/index.html", "3 posts")
		verifyFileContent(t, sitePath, "series/index.html", "1 posts")

		// Individual series pages should exist
		verifyFileExists(t, sitePath, "series/test-series/index.html")
		verifyFileExists(t, sitePath, "series/another-series/index.html")

		// Individual series page should show series name
		verifyFileContent(t, sitePath, "series/test-series/index.html", "Series: test-series")
		verifyFileContent(t, sitePath, "series/test-series/index.html", "3 posts in this series")
	})

	t.Run("SeriesPosts_OrderedCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Read series page content
		seriesPath := filepath.Join(sitePath, "public/series/test-series/index.html")
		content, err := os.ReadFile(seriesPath)
		if err != nil {
			t.Fatalf("failed to read series page: %v", err)
		}

		contentStr := string(content)

		// Posts should appear in order (Part 1, 2, 3)
		part1Pos := strings.Index(contentStr, "Test Series Part 1")
		part2Pos := strings.Index(contentStr, "Test Series Part 2")
		part3Pos := strings.Index(contentStr, "Test Series Part 3")

		if part1Pos == -1 || part2Pos == -1 || part3Pos == -1 {
			t.Fatal("series posts not found in series page")
		}

		if part1Pos > part2Pos || part2Pos > part3Pos {
			t.Error("series posts not in correct order (expected: Part 1 < Part 2 < Part 3)")
		}

		// Verify position numbers appear
		verifyFileContent(t, sitePath, "series/test-series/index.html", "Part 1")
		verifyFileContent(t, sitePath, "series/test-series/index.html", "Part 2")
		verifyFileContent(t, sitePath, "series/test-series/index.html", "Part 3")
	})

	t.Run("SeriesNavigation_FirstPost", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Part 1 should have series info
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", "Part 1 of 3 in")
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", "test-series")

		// Part 1 should have Next link
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", "Next:")
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", "Test Series Part 2")

		// Part 1 should NOT have Previous link
		content, err := os.ReadFile(filepath.Join(sitePath, "public/blog/series-part1/index.html"))
		if err != nil {
			t.Fatalf("failed to read part 1: %v", err)
		}
		if strings.Contains(string(content), "Previous:") {
			t.Error("Part 1 should not have Previous link")
		}
	})

	t.Run("SeriesNavigation_MiddlePost", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Part 2 should have series info
		verifyFileContent(t, sitePath, "blog/series-part2/index.html", "Part 2 of 3 in")
		verifyFileContent(t, sitePath, "blog/series-part2/index.html", "test-series")

		// Part 2 should have Previous link
		verifyFileContent(t, sitePath, "blog/series-part2/index.html", "Previous:")
		verifyFileContent(t, sitePath, "blog/series-part2/index.html", "Test Series Part 1")

		// Part 2 should have Next link
		verifyFileContent(t, sitePath, "blog/series-part2/index.html", "Next:")
		verifyFileContent(t, sitePath, "blog/series-part2/index.html", "Test Series Part 3")
	})

	t.Run("SeriesNavigation_LastPost", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Part 3 should have series info
		verifyFileContent(t, sitePath, "blog/series-part3/index.html", "Part 3 of 3 in")
		verifyFileContent(t, sitePath, "blog/series-part3/index.html", "test-series")

		// Part 3 should have Previous link
		verifyFileContent(t, sitePath, "blog/series-part3/index.html", "Previous:")
		verifyFileContent(t, sitePath, "blog/series-part3/index.html", "Test Series Part 2")

		// Part 3 should NOT have Next link
		content, err := os.ReadFile(filepath.Join(sitePath, "public/blog/series-part3/index.html"))
		if err != nil {
			t.Fatalf("failed to read part 3: %v", err)
		}
		if strings.Contains(string(content), "Next:") {
			t.Error("Part 3 should not have Next link")
		}
	})

	t.Run("MultipleSeries_IndependentNavigation", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Another series should have its own navigation
		verifyFileContent(t, sitePath, "blog/another-series-part1/index.html", "Part 1 of 1 in")
		verifyFileContent(t, sitePath, "blog/another-series-part1/index.html", "another-series")

		// Another series Part 1 should NOT have navigation to test-series
		content, err := os.ReadFile(filepath.Join(sitePath, "public/blog/another-series-part1/index.html"))
		if err != nil {
			t.Fatalf("failed to read another-series part 1: %v", err)
		}
		if strings.Contains(string(content), "test-series") {
			t.Error("another-series should not reference test-series")
		}

		// Both series should exist independently
		verifyFileExists(t, sitePath, "series/test-series/index.html")
		verifyFileExists(t, sitePath, "series/another-series/index.html")
	})

	t.Run("PostWithoutSeries_NoNavigation", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Regular posts without series should not have series navigation
		content, err := os.ReadFile(filepath.Join(sitePath, "public/blog/post1/index.html"))
		if err != nil {
			t.Fatalf("failed to read post1: %v", err)
		}

		contentStr := string(content)
		if strings.Contains(contentStr, "series-navigation") {
			t.Error("post without series should not have series navigation")
		}
		if strings.Contains(contentStr, "Part") && strings.Contains(contentStr, "of") {
			t.Error("post without series should not show part numbers")
		}
	})

	t.Run("AddPostToSeries_UpdatesPages", func(t *testing.T) {
		t.Skip("TODO: Fix series navigation to update post counts and next/prev links when series grows")
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Verify initial state
		verifyFileContent(t, sitePath, "series/test-series/index.html", "3 posts in this series")

		// Add new post to test-series
		newPost := `+++
title = "Test Series Part 4"
date = 2024-01-04
template = "post.html"
series = "test-series"
series_order = 4

[extra]
description = "Fourth part of the test series"
+++

This is part 4 added dynamically.`

		createPost(t, sitePath, "blog/series-part4.md", newPost)

		// Rebuild
		fullRebuild(t, gen, contentParser, sitePath)

		// Verify series page updated
		verifyFileContent(t, sitePath, "series/test-series/index.html", "4 posts in this series")
		verifyFileContent(t, sitePath, "series/test-series/index.html", "Test Series Part 4")

		// Verify Part 3 now has Next link
		verifyFileContent(t, sitePath, "blog/series-part3/index.html", "Part 3 of 4 in")
		verifyFileContent(t, sitePath, "blog/series-part3/index.html", "Next:")
		verifyFileContent(t, sitePath, "blog/series-part3/index.html", "Test Series Part 4")

		// Verify Part 4 has navigation
		verifyFileExists(t, sitePath, "blog/series-part4/index.html")
		verifyFileContent(t, sitePath, "blog/series-part4/index.html", "Part 4 of 4 in")
		verifyFileContent(t, sitePath, "blog/series-part4/index.html", "Previous:")
	})

	t.Run("ChangeSeriesOrder_RebuildsCorrectly", func(t *testing.T) {
		t.Skip("TODO: Fix series order changes to recalculate part numbers and navigation links")
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Change Part 2's order to 5 (make it last)
		part2Path := filepath.Join(sitePath, "content/blog/series-part2.md")
		modifyFile(t, part2Path, "series_order = 2", "series_order = 5")

		// Rebuild
		fullRebuild(t, gen, contentParser, sitePath)

		// Part 2 should now be last (Part 5)
		verifyFileContent(t, sitePath, "blog/series-part2/index.html", "Part 5 of 3 in")

		// Part 1 should now link to Part 3 as next
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", "Next:")
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", "Test Series Part 3")

		// Part 3 should have links to Part 1 and Part 2
		verifyFileContent(t, sitePath, "blog/series-part3/index.html", "Previous:")
		verifyFileContent(t, sitePath, "blog/series-part3/index.html", "Test Series Part 1")
		verifyFileContent(t, sitePath, "blog/series-part3/index.html", "Next:")
		verifyFileContent(t, sitePath, "blog/series-part3/index.html", "Test Series Part 2")
	})

	t.Run("DraftPostInSeries_ExcludedFromNavigation", func(t *testing.T) {
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

		createPost(t, sitePath, "blog/series-draft.md", draftPost)

		// Rebuild
		fullRebuild(t, gen, contentParser, sitePath)

		// Series should still show 3 posts (draft excluded)
		verifyFileContent(t, sitePath, "series/test-series/index.html", "3 posts in this series")

		// Part 2 should still link directly to Part 3 (skip draft)
		verifyFileContent(t, sitePath, "blog/series-part2/index.html", "Next:")
		verifyFileContent(t, sitePath, "blog/series-part2/index.html", "Test Series Part 3")

		// Draft should not exist in output
		verifyFileNotExists(t, sitePath, "blog/series-draft/index.html")
	})

	t.Run("SeriesPermalinks_CorrectFormat", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Series index should link to individual series correctly
		verifyFileContent(t, sitePath, "series/index.html", `href="/series/test-series/"`)
		verifyFileContent(t, sitePath, "series/index.html", `href="/series/another-series/"`)

		// Post navigation should link correctly
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", `href="/series/test-series/"`)
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", `href="/blog/series-part2/"`)
	})
}

// TestTaxonomy_Tags tests tag taxonomy functionality
func TestTaxonomy_Tags(t *testing.T) {
	t.Run("TagPages_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Tag pages should exist
		verifyFileExists(t, sitePath, "tags/golang/index.html")
		verifyFileExists(t, sitePath, "tags/testing/index.html")

		// Tag pages should list posts
		verifyFileContent(t, sitePath, "tags/golang/index.html", "Third Post")
		verifyFileContent(t, sitePath, "tags/testing/index.html", "Third Post")

		// Tags index should exist
		verifyFileExists(t, sitePath, "tags/index.html")
		verifyFileContent(t, sitePath, "tags/index.html", "golang")
		verifyFileContent(t, sitePath, "tags/index.html", "testing")
	})

	t.Run("AddTag_CreatesTagPage", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Add tag to post1
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		insertAfterLine(t, post1Path, "template =", `tags = ["database"]`)

		// Rebuild
		fullRebuild(t, gen, contentParser, sitePath)

		// Tag page should be created
		verifyFileExists(t, sitePath, "tags/database/index.html")
		verifyFileContent(t, sitePath, "tags/database/index.html", "First Post")

		// Tag should appear in tags index
		verifyFileContent(t, sitePath, "tags/index.html", "database")

		// Post should link to tag
		verifyFileContent(t, sitePath, "blog/post1/index.html", "/tags/database/")
	})

	t.Run("RemoveTag_UpdatesPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create post with tag
		post := `+++
title = "Tagged Post"
date = 2024-01-15
template = "post.html"
tags = ["removeme"]
+++
Content`
		createPost(t, sitePath, "blog/tagged.md", post)
		fullRebuild(t, gen, contentParser, sitePath)
		verifyFileExists(t, sitePath, "tags/removeme/index.html")

		// Remove tag
		modifyFile(t, filepath.Join(sitePath, "content/blog/tagged.md"), `tags = ["removeme"]`, "")

		// Clean public directory before rebuild to ensure old files don't persist
		publicDir := filepath.Join(sitePath, "public")
		os.RemoveAll(publicDir)
		os.MkdirAll(publicDir, 0755)

		// Rebuild
		fullRebuild(t, gen, contentParser, sitePath)

		// Tag page should not exist anymore since no posts have that tag
		verifyFileNotExists(t, sitePath, "tags/removeme/index.html")
	})

	t.Run("MultipleTagsOnPost_AllPagesGenerated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create post with multiple tags
		post := `+++
title = "Multi Tag Post"
date = 2024-01-20
template = "post.html"
tags = ["alpha", "beta", "gamma"]
+++
Content`
		createPost(t, sitePath, "blog/multi-tag.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		// All tag pages should exist
		verifyFileExists(t, sitePath, "tags/alpha/index.html")
		verifyFileExists(t, sitePath, "tags/beta/index.html")
		verifyFileExists(t, sitePath, "tags/gamma/index.html")

		// All should show the post
		verifyFileContent(t, sitePath, "tags/alpha/index.html", "Multi Tag Post")
		verifyFileContent(t, sitePath, "tags/beta/index.html", "Multi Tag Post")
		verifyFileContent(t, sitePath, "tags/gamma/index.html", "Multi Tag Post")

		// Post should link to all tags
		verifyFileContent(t, sitePath, "blog/multi-tag/index.html", "/tags/alpha/")
		verifyFileContent(t, sitePath, "blog/multi-tag/index.html", "/tags/beta/")
		verifyFileContent(t, sitePath, "blog/multi-tag/index.html", "/tags/gamma/")
	})

	t.Run("TagCounts_AccurateOnIndex", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create multiple posts with same tag
		dates := []string{"01", "02", "03"}
		for i := 1; i <= 3; i++ {
			post := fmt.Sprintf(`+++
title = "Popular Post %d"
date = 2024-01-%s
template = "post.html"
tags = ["popular"]
+++
Content`, i, dates[i-1])
			createPost(t, sitePath, fmt.Sprintf("blog/popular%d/index.md", i), post)
		}

		fullRebuild(t, gen, contentParser, sitePath)

		// Tag page should show all posts
		tagContent := readFileContent(t, sitePath, "tags/popular/index.html")
		for i := 1; i <= 3; i++ {
			if !strings.Contains(tagContent, fmt.Sprintf("Popular Post %d", i)) {
				t.Errorf("tag page missing Popular Post %d", i)
			}
		}
	})

	t.Run("TagsWithSpecialCharacters_Sanitized", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create post with tag containing special chars
		post := `+++
title = "Special Tag Post"
date = 2024-01-25
template = "post.html"
tags = ["C++", "Node.js"]
+++
Content`
		createPost(t, sitePath, "blog/special-tags.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		// Tags should appear in content (HTML-encoded)
		verifyFileContent(t, sitePath, "blog/special-tags/index.html", "C&#43;&#43;")
		verifyFileContent(t, sitePath, "blog/special-tags/index.html", "Node.js")
	})
}

// TestTaxonomy_CrossTaxonomy tests interactions between taxonomies
func TestTaxonomy_CrossTaxonomy(t *testing.T) {
	t.Run("PostWithSeriesAndTags_BothWork", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create post with both series and tags
		post := `+++
title = "Hybrid Post"
date = 2024-01-30
template = "post.html"
series = "hybrid-series"
series_order = 1
tags = ["hybrid", "multi"]
+++
Content`
		createPost(t, sitePath, "blog/hybrid.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		// Series page should exist
		verifyFileExists(t, sitePath, "series/hybrid-series/index.html")
		verifyFileContent(t, sitePath, "series/hybrid-series/index.html", "Hybrid Post")

		// Tag pages should exist
		verifyFileExists(t, sitePath, "tags/hybrid/index.html")
		verifyFileExists(t, sitePath, "tags/multi/index.html")
		verifyFileContent(t, sitePath, "tags/hybrid/index.html", "Hybrid Post")
		verifyFileContent(t, sitePath, "tags/multi/index.html", "Hybrid Post")

		// Post should have both series navigation and tags
		verifyFileContent(t, sitePath, "blog/hybrid/index.html", "hybrid-series")
		verifyFileContent(t, sitePath, "blog/hybrid/index.html", "/tags/hybrid/")
		verifyFileContent(t, sitePath, "blog/hybrid/index.html", "/tags/multi/")
	})

	t.Run("EmptyTaxonomy_NoPageGenerated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Only taxonomies with posts should have pages
		// Check that non-existent tags don't have pages
		nonExistentPath := filepath.Join(sitePath, "public/tags/nonexistent/index.html")
		if _, err := os.Stat(nonExistentPath); !os.IsNotExist(err) {
			t.Error("empty taxonomy should not have page")
		}
	})

	t.Run("TaxonomyOrdering_PostsByDate", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create multiple posts with same tag, different dates
		post1 := `+++
title = "Older Post"
date = 2024-01-01
template = "post.html"
tags = ["ordered"]
+++
Content`

		post2 := `+++
title = "Newer Post"
date = 2024-02-01
template = "post.html"
tags = ["ordered"]
+++
Content`

		createPost(t, sitePath, "blog/older.md", post1)
		createPost(t, sitePath, "blog/newer.md", post2)
		fullRebuild(t, gen, contentParser, sitePath)

		// Tag page should list posts in date order (newest first)
		tagContent := readFileContent(t, sitePath, "tags/ordered/index.html")
		newerPos := strings.Index(tagContent, "Newer Post")
		olderPos := strings.Index(tagContent, "Older Post")

		if newerPos == -1 || olderPos == -1 {
			t.Fatal("posts not found in tag page")
		}

		if newerPos > olderPos {
			t.Error("posts should be ordered newest first")
		}
	})
}

// TestTaxonomy_DraftExclusion tests that drafts don't appear in taxonomies
func TestTaxonomy_DraftExclusion(t *testing.T) {
	t.Run("DraftWithTags_NoTagPage", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create draft with tags
		draft := `+++
title = "Draft Post"
date = 2024-01-01
template = "post.html"
tags = ["drafttag"]
draft = true
+++
Content`
		createPost(t, sitePath, "blog/draft-tagged.md", draft)
		fullRebuild(t, gen, contentParser, sitePath)

		// Tag page should not exist or be empty
		tagPath := filepath.Join(sitePath, "public/tags/drafttag/index.html")
		if _, err := os.Stat(tagPath); !os.IsNotExist(err) {
			// If it exists, should not contain draft
			verifyFileNotContains(t, sitePath, "tags/drafttag/index.html", "Draft Post")
		}

		// Draft should not be in output
		verifyFileNotExists(t, sitePath, "blog/draft-tagged/index.html")
	})

	t.Run("DraftWithSeries_NoSeriesPage", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create draft in series
		draft := `+++
title = "Draft Series Post"
date = 2024-01-01
template = "post.html"
series = "draft-series"
series_order = 1
draft = true
+++
Content`
		createPost(t, sitePath, "blog/draft-series.md", draft)
		fullRebuild(t, gen, contentParser, sitePath)

		// Series page should not exist or be empty
		seriesPath := filepath.Join(sitePath, "public/series/draft-series/index.html")
		if _, err := os.Stat(seriesPath); !os.IsNotExist(err) {
			verifyFileNotContains(t, sitePath, "series/draft-series/index.html", "Draft Series Post")
		}
	})
}

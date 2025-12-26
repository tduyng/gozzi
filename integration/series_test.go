package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestSeriesTaxonomy tests the series feature end-to-end
func TestSeriesTaxonomy(t *testing.T) {
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

	t.Run("SeriesNavigation_OnFirstPost", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Part 1 should have series info
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", "Part 1 of 3 in")
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", "test-series")

		// Part 1 should have Next link
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", "Next:")
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", "Test Series Part 2")

		// Part 1 should NOT have Previous link (checked by reading full content)
		content, err := os.ReadFile(filepath.Join(sitePath, "public/blog/series-part1/index.html"))
		if err != nil {
			t.Fatalf("failed to read part 1: %v", err)
		}
		if strings.Contains(string(content), "Previous:") {
			t.Error("Part 1 should not have Previous link")
		}
	})

	t.Run("SeriesNavigation_OnMiddlePost", func(t *testing.T) {
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

	t.Run("SeriesNavigation_OnLastPost", func(t *testing.T) {
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

	t.Run("SeriesWithoutPosts_NotGenerated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Only series with posts should appear
		content, err := os.ReadFile(filepath.Join(sitePath, "public/series/index.html"))
		if err != nil {
			t.Fatalf("failed to read series index: %v", err)
		}

		// Should only have test-series and another-series
		contentStr := string(content)
		seriesCount := strings.Count(contentStr, "series-item")
		if seriesCount != 2 {
			t.Errorf("expected 2 series, found %d series items", seriesCount)
		}
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

		newPostPath := filepath.Join(sitePath, "content/blog/series-part4.md")
		if err := os.WriteFile(newPostPath, []byte(newPost), 0644); err != nil {
			t.Fatalf("failed to write new post: %v", err)
		}

		// Rebuild
		if err := contentParser.Parse(filepath.Join(sitePath, "content")); err != nil {
			t.Fatalf("failed to re-parse content: %v", err)
		}
		if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("failed to rebuild: %v", err)
		}

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
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Change Part 2's order to 5 (make it last)
		part2Path := filepath.Join(sitePath, "content/blog/series-part2.md")
		content, err := os.ReadFile(part2Path)
		if err != nil {
			t.Fatalf("failed to read part2: %v", err)
		}

		modifiedContent := strings.Replace(string(content), "series_order = 2", "series_order = 5", 1)
		if err := os.WriteFile(part2Path, []byte(modifiedContent), 0644); err != nil {
			t.Fatalf("failed to write part2: %v", err)
		}

		// Rebuild
		if err := contentParser.Parse(filepath.Join(sitePath, "content")); err != nil {
			t.Fatalf("failed to re-parse content: %v", err)
		}
		if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("failed to rebuild: %v", err)
		}

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

		draftPath := filepath.Join(sitePath, "content/blog/series-draft.md")
		if err := os.WriteFile(draftPath, []byte(draftPost), 0644); err != nil {
			t.Fatalf("failed to write draft: %v", err)
		}

		// Rebuild
		if err := contentParser.Parse(filepath.Join(sitePath, "content")); err != nil {
			t.Fatalf("failed to re-parse content: %v", err)
		}
		if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("failed to rebuild: %v", err)
		}

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

	t.Run("SeriesHTML_ContainsDataAttributes", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Series index should have data attributes for testing/styling
		verifyFileContent(t, sitePath, "series/index.html", `data-series="test-series"`)
		verifyFileContent(t, sitePath, "series/index.html", `data-series="another-series"`)

		// Series term page should have position attributes
		verifyFileContent(t, sitePath, "series/test-series/index.html", `data-position="1"`)
		verifyFileContent(t, sitePath, "series/test-series/index.html", `data-position="2"`)
		verifyFileContent(t, sitePath, "series/test-series/index.html", `data-position="3"`)

		// Individual posts should have series data attribute
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", `data-series="test-series"`)
	})

	t.Run("SeriesPermalinks_CorrectFormat", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Series index should link to individual series correctly
		verifyFileContent(t, sitePath, "series/index.html", `href="/series/test-series/"`)
		verifyFileContent(t, sitePath, "series/index.html", `href="/series/another-series/"`)

		// Individual series should link back to series index
		seriesContent, err := os.ReadFile(filepath.Join(sitePath, "public/series/test-series/index.html"))
		if err != nil {
			t.Fatalf("failed to read series page: %v", err)
		}

		// Posts in series should link correctly
		if !strings.Contains(string(seriesContent), "/blog/series-part1/") {
			t.Error("series page should link to post permalinks")
		}

		// Post navigation should link correctly
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", `href="/series/test-series/"`)
		verifyFileContent(t, sitePath, "blog/series-part1/index.html", `href="/blog/series-part2/"`)
	})
}

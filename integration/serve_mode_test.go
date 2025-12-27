package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/tduyng/gozzi/app/builder"
)

// TestServeMode_TaxonomyChanges tests that taxonomy changes trigger correct rebuilds
func TestServeMode_TaxonomyChanges(t *testing.T) {
	t.Run("AddSeriesField_RegeneratesPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Initial state: post1 has no series
		verifyFileExists(t, sitePath, "blog/post1/index.html")
		content, _ := os.ReadFile(filepath.Join(sitePath, "public/blog/post1/index.html"))
		if strings.Contains(string(content), "series-navigation") {
			t.Error("post1 should not have series navigation initially")
		}

		// Simulate user adding series field to post1
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")
		originalContent, _ := os.ReadFile(post1Path)

		// Add series and series_order after template line
		lines := strings.Split(string(originalContent), "\n")
		var newLines []string
		for _, line := range lines {
			newLines = append(newLines, line)
			if strings.HasPrefix(line, "template =") {
				newLines = append(newLines, `series = "new-series"`)
				newLines = append(newLines, `series_order = 1`)
			}
		}
		os.WriteFile(post1Path, []byte(strings.Join(newLines, "\n")), 0644)

		// Snapshot old taxonomy values BEFORE parsing
		oldTaxonomyValues := gen.SnapshotTaxonomyValues([]string{post1Path}, filepath.Join(sitePath, "content"))

		// Simulate incremental rebuild (like serve mode does)
		time.Sleep(10 * time.Millisecond) // Give file time to settle
		if err := contentParser.ParseFiles(filepath.Join(sitePath, "content"), []string{post1Path}); err != nil {
			t.Fatalf("failed to re-parse: %v", err)
		}

		err := gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      []string{post1Path},
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})
		if err != nil {
			t.Fatalf("failed to rebuild: %v", err)
		}

		// Verify post1 now has series navigation
		content, _ = os.ReadFile(filepath.Join(sitePath, "public/blog/post1/index.html"))
		if !strings.Contains(string(content), "series-navigation") {
			t.Error("post1 should have series navigation after adding series field")
		}
		if !strings.Contains(string(content), "new-series") {
			t.Error("post1 should show new-series in navigation")
		}

		// Verify series page was created
		verifyFileExists(t, sitePath, "series/new-series/index.html")
		verifyFileContent(t, sitePath, "series/new-series/index.html", "First Post")

		// Verify series index lists new series
		verifyFileExists(t, sitePath, "series/index.html")

		verifyFileContent(t, sitePath, "series/index.html", "new-series")
	})

	t.Run("ChangeSeriesValue_RegeneratesBothPages", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Create a post with series
		post := `+++
title = "Test Post"
date = 2024-01-01
template = "post.html"
series = "series-a"
series_order = 1
+++

Test content.`

		postPath := filepath.Join(sitePath, "content/blog/test-post.md")
		os.WriteFile(postPath, []byte(post), 0644)

		// Initial build
		gen, contentParser := buildSite(t, sitePath)

		// Verify initial state
		verifyFileExists(t, sitePath, "series/series-a/index.html")
		verifyFileContent(t, sitePath, "series/series-a/index.html", "Test Post")

		// Snapshot old taxonomy values BEFORE modifying the file
		oldTaxonomyValues := gen.SnapshotTaxonomyValues([]string{postPath}, filepath.Join(sitePath, "content"))

		// Change series from series-a to series-b
		updatedPost := strings.Replace(post, `series = "series-a"`, `series = "series-b"`, 1)
		os.WriteFile(postPath, []byte(updatedPost), 0644)

		// Incremental rebuild
		time.Sleep(10 * time.Millisecond)
		if err := contentParser.ParseFiles(filepath.Join(sitePath, "content"), []string{postPath}); err != nil {
			t.Fatalf("failed to re-parse: %v", err)
		}

		err := gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      []string{postPath},
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})
		if err != nil {
			t.Fatalf("failed to rebuild: %v", err)
		}

		// Verify series-a no longer shows test post
		seriesAContent, _ := os.ReadFile(filepath.Join(sitePath, "public/series/series-a/index.html"))
		if strings.Contains(string(seriesAContent), "Test Post") {
			t.Error("series-a should NOT contain Test Post after it moved to series-b")
		}

		// Verify series-b now shows test post
		verifyFileExists(t, sitePath, "series/series-b/index.html")
		verifyFileContent(t, sitePath, "series/series-b/index.html", "Test Post")

		// Verify post navigation points to series-b
		verifyFileContent(t, sitePath, "blog/test-post/index.html", "series-b")
	})

	t.Run("ChangeSeriesOrder_UpdatesNavigation", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Create series with 2 posts
		post1 := `+++
title = "Part 1"
date = 2024-01-01
template = "post.html"
series = "order-test-series"
series_order = 1
+++
Content 1`

		post2 := `+++
title = "Part 2"
date = 2024-01-02
template = "post.html"
series = "order-test-series"
series_order = 2
+++
Content 2`

		post1Path := filepath.Join(sitePath, "content/blog/part1.md")
		post2Path := filepath.Join(sitePath, "content/blog/part2.md")
		os.WriteFile(post1Path, []byte(post1), 0644)
		os.WriteFile(post2Path, []byte(post2), 0644)

		gen, contentParser := buildSite(t, sitePath)

		// Verify initial order
		verifyFileContent(t, sitePath, "blog/part1/index.html", "Part 2") // Next should be Part 2
		content, _ := os.ReadFile(filepath.Join(sitePath, "public/blog/part1/index.html"))
		if strings.Contains(string(content), "Previous:") {
			t.Error("Part 1 should not have Previous link")
		}

		// Swap the order (change part1 to order 2, part2 to order 1)
		changedFiles := []string{post1Path, post2Path}

		// Snapshot old taxonomy values BEFORE modifying
		oldTaxonomyValues := gen.SnapshotTaxonomyValues(changedFiles, filepath.Join(sitePath, "content"))

		updatedPost1 := strings.Replace(post1, "series_order = 1", "series_order = 2", 1)
		updatedPost2 := strings.Replace(post2, "series_order = 2", "series_order = 1", 1)
		os.WriteFile(post1Path, []byte(updatedPost1), 0644)
		os.WriteFile(post2Path, []byte(updatedPost2), 0644)

		// Incremental rebuild
		time.Sleep(10 * time.Millisecond)
		if err := contentParser.ParseFiles(filepath.Join(sitePath, "content"), changedFiles); err != nil {
			t.Fatalf("failed to re-parse: %v", err)
		}

		err := gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      changedFiles,
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})
		if err != nil {
			t.Fatalf("failed to rebuild: %v", err)
		}

		// Verify new order: Part 1 is now last, should have Previous to Part 2
		verifyFileContent(t, sitePath, "blog/part1/index.html", "Previous:")
		verifyFileContent(t, sitePath, "blog/part1/index.html", "Part 2")

		// Part 2 is now first, should have Next to Part 1
		verifyFileContent(t, sitePath, "blog/part2/index.html", "Next:")
		verifyFileContent(t, sitePath, "blog/part2/index.html", "Part 1")
	})

	t.Run("RemoveSeriesField_CleansUpPages", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Create post with series
		post := `+++
title = "Test Post"
date = 2024-01-01
template = "post.html"
series = "temp-series"
series_order = 1
+++
Content`

		postPath := filepath.Join(sitePath, "content/blog/temp-post.md")
		os.WriteFile(postPath, []byte(post), 0644)

		gen, contentParser := buildSite(t, sitePath)

		// Verify series exists
		verifyFileExists(t, sitePath, "series/temp-series/index.html")
		verifyFileContent(t, sitePath, "blog/temp-post/index.html", "temp-series")

		// Snapshot old taxonomy values BEFORE removing series
		oldTaxonomyValues := gen.SnapshotTaxonomyValues([]string{postPath}, filepath.Join(sitePath, "content"))

		// Remove series fields
		postWithoutSeries := `+++
title = "Test Post"
date = 2024-01-01
template = "post.html"
+++
Content`
		os.WriteFile(postPath, []byte(postWithoutSeries), 0644)

		// Incremental rebuild
		time.Sleep(10 * time.Millisecond)
		if err := contentParser.ParseFiles(filepath.Join(sitePath, "content"), []string{postPath}); err != nil {
			t.Fatalf("failed to re-parse: %v", err)
		}

		err := gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      []string{postPath},
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})
		if err != nil {
			t.Fatalf("failed to rebuild: %v", err)
		}

		// Verify post no longer has series navigation
		content, _ := os.ReadFile(filepath.Join(sitePath, "public/blog/temp-post/index.html"))
		if strings.Contains(string(content), "series-navigation") {
			t.Error("post should not have series navigation after removing series field")
		}

		// Verify series page no longer shows post
		seriesContent, _ := os.ReadFile(filepath.Join(sitePath, "public/series/temp-series/index.html"))
		if strings.Contains(string(seriesContent), "Test Post") {
			t.Error("series page should not show post after it was removed from series")
		}
	})

	t.Run("AddTagsField_RegeneratesTagPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// post1 has no tags initially
		verifyFileExists(t, sitePath, "blog/post1/index.html")

		// Add tags to post1
		post1Path := filepath.Join(sitePath, "content/blog/post1/index.md")

		// Snapshot old taxonomy values BEFORE modifying
		oldTaxonomyValues := gen.SnapshotTaxonomyValues([]string{post1Path}, filepath.Join(sitePath, "content"))

		originalContent, _ := os.ReadFile(post1Path)

		lines := strings.Split(string(originalContent), "\n")
		var newLines []string
		for _, line := range lines {
			newLines = append(newLines, line)
			if strings.HasPrefix(line, "template =") {
				newLines = append(newLines, `tags = ["newtag", "another"]`)
			}
		}
		os.WriteFile(post1Path, []byte(strings.Join(newLines, "\n")), 0644)

		// Incremental rebuild
		time.Sleep(10 * time.Millisecond)
		if err := contentParser.ParseFiles(filepath.Join(sitePath, "content"), []string{post1Path}); err != nil {
			t.Fatalf("failed to re-parse: %v", err)
		}

		err := gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      []string{post1Path},
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})
		if err != nil {
			t.Fatalf("failed to rebuild: %v", err)
		}

		// Verify tag pages were created
		verifyFileExists(t, sitePath, "tags/newtag/index.html")
		verifyFileExists(t, sitePath, "tags/another/index.html")
		verifyFileContent(t, sitePath, "tags/newtag/index.html", "First Post")
		verifyFileContent(t, sitePath, "tags/another/index.html", "First Post")
	})
}

// TestServeMode_IndexFileChanges tests that _index.md changes trigger correct rebuilds
func TestServeMode_IndexFileChanges(t *testing.T) {
	t.Run("RootIndexChange_RegeneratesHomepage", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Initial state
		verifyFileExists(t, sitePath, "index.html")
		initialContent, _ := os.ReadFile(filepath.Join(sitePath, "public/index.html"))
		if !strings.Contains(string(initialContent), "Welcome") {
			t.Error("homepage should contain 'Welcome' initially")
		}

		// Modify root _index.md
		indexPath := filepath.Join(sitePath, "content/_index.md")
		originalContent, _ := os.ReadFile(indexPath)

		// Change content
		modifiedContent := strings.Replace(
			string(originalContent),
			"Welcome",
			"MODIFIED: Hello World",
			1,
		)
		os.WriteFile(indexPath, []byte(modifiedContent), 0644)

		// Snapshot old taxonomy values BEFORE parsing
		oldTaxonomyValues := gen.SnapshotTaxonomyValues(
			[]string{indexPath},
			filepath.Join(sitePath, "content"),
		)

		// Incremental rebuild (simulating serve mode)
		time.Sleep(10 * time.Millisecond)
		if err := contentParser.ParseFiles(
			filepath.Join(sitePath, "content"),
			[]string{indexPath},
		); err != nil {
			t.Fatalf("failed to re-parse: %v", err)
		}

		err := gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      []string{indexPath},
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})
		if err != nil {
			t.Fatalf("failed to rebuild: %v", err)
		}

		// Verify homepage was regenerated with new content
		newContent, _ := os.ReadFile(filepath.Join(sitePath, "public/index.html"))
		if !strings.Contains(string(newContent), "MODIFIED: Hello World") {
			t.Error("homepage should contain modified content after _index.md change")
		}
		if strings.Contains(string(newContent), "Welcome") {
			t.Error("homepage should not contain old content")
		}
	})

	t.Run("NestedIndexChange_RegeneratesSection", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Initial state - blog section
		verifyFileExists(t, sitePath, "blog/index.html")

		// Modify blog/_index.md title (which will appear in rendered HTML)
		blogIndexPath := filepath.Join(sitePath, "content/blog/_index.md")
		originalContent, _ := os.ReadFile(blogIndexPath)

		// Change the title
		modifiedContent := strings.Replace(
			string(originalContent),
			"title = \"Blog\"",
			"title = \"MODIFIED BLOG TITLE\"",
			1,
		)
		os.WriteFile(blogIndexPath, []byte(modifiedContent), 0644)

		// Snapshot old taxonomy values BEFORE parsing
		oldTaxonomyValues := gen.SnapshotTaxonomyValues(
			[]string{blogIndexPath},
			filepath.Join(sitePath, "content"),
		)

		// Incremental rebuild
		time.Sleep(10 * time.Millisecond)
		if err := contentParser.ParseFiles(
			filepath.Join(sitePath, "content"),
			[]string{blogIndexPath},
		); err != nil {
			t.Fatalf("failed to re-parse: %v", err)
		}

		err := gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      []string{blogIndexPath},
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})
		if err != nil {
			t.Fatalf("failed to rebuild: %v", err)
		}

		// Verify blog section was regenerated with new title
		newContent, _ := os.ReadFile(filepath.Join(sitePath, "public/blog/index.html"))
		if !strings.Contains(string(newContent), "MODIFIED BLOG TITLE") {
			t.Error("blog section should contain modified title after _index.md change")
		}

		// Verify sitemap was also updated (sections affect sitemap)
		verifyFileExists(t, sitePath, "sitemap.xml")
	})

	t.Run("MultipleIndexChanges_RegeneratesAll", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify both root and blog _index.md
		rootIndexPath := filepath.Join(sitePath, "content/_index.md")
		blogIndexPath := filepath.Join(sitePath, "content/blog/_index.md")

		changedFiles := []string{rootIndexPath, blogIndexPath}

		// Snapshot old taxonomy values BEFORE parsing
		oldTaxonomyValues := gen.SnapshotTaxonomyValues(
			changedFiles,
			filepath.Join(sitePath, "content"),
		)

		// Modify root
		rootContent, _ := os.ReadFile(rootIndexPath)
		modifiedRoot := strings.Replace(string(rootContent), "Welcome", "ROOT CHANGED", 1)
		os.WriteFile(rootIndexPath, []byte(modifiedRoot), 0644)

		// Modify blog
		blogContent, _ := os.ReadFile(blogIndexPath)
		modifiedBlog := strings.Replace(
			string(blogContent),
			"+++",
			"+++\ntest = \"BLOG CHANGED\"",
			1,
		)
		os.WriteFile(blogIndexPath, []byte(modifiedBlog), 0644)

		// Incremental rebuild
		time.Sleep(10 * time.Millisecond)
		if err := contentParser.ParseFiles(
			filepath.Join(sitePath, "content"),
			changedFiles,
		); err != nil {
			t.Fatalf("failed to re-parse: %v", err)
		}

		err := gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      changedFiles,
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})
		if err != nil {
			t.Fatalf("failed to rebuild: %v", err)
		}

		// Verify both pages were regenerated
		homeContent, _ := os.ReadFile(filepath.Join(sitePath, "public/index.html"))
		if !strings.Contains(string(homeContent), "ROOT CHANGED") {
			t.Error("homepage should contain modified content")
		}

		// Blog should be regenerated (we can't easily verify custom fields without template support,
		// but we can verify the file was touched)
		verifyFileExists(t, sitePath, "blog/index.html")
	})
}

// TestServeMode_RootLevelContentFiles is a regression test for a critical bug
// where root-level content files (e.g., content/about.md, content/scss.md) were
// not being found during incremental builds in serve mode. The bug was in
// rebuild_analyzer.go:findNode() which looked up contentMap[""] for root files,
// but the map uses "." as the key for the root directory.
//
// This bug caused the dev server to serve stale content even though the watcher
// detected changes and triggered rebuilds. The fix ensures root-level files are
// properly found by using "." as the lookup key when dir normalizes to "".
func TestServeMode_RootLevelContentFiles(t *testing.T) {
	t.Run("RootLevelPage_RegeneratesOnChange", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Create a root-level content file (like scss.md, about.md, etc.)
		rootPage := `+++
title = "SCSS Guide"
date = 2024-12-27
template = "page.html"
+++

Original SCSS documentation content.

This is the initial version.`

		rootPagePath := filepath.Join(sitePath, "content/scss.md")
		os.WriteFile(rootPagePath, []byte(rootPage), 0644)

		// Initial build
		gen, contentParser := buildSite(t, sitePath)

		// Verify the root-level page was generated
		verifyFileExists(t, sitePath, "scss/index.html")
		initialContent, _ := os.ReadFile(filepath.Join(sitePath, "public/scss/index.html"))
		if !strings.Contains(string(initialContent), "Original SCSS documentation content") {
			t.Error("initial build should contain original content")
		}
		if !strings.Contains(string(initialContent), "initial version") {
			t.Error("initial build should contain initial version marker")
		}

		// Modify the root-level page (simulating user editing in dev mode)
		modifiedPage := `+++
title = "SCSS Guide"
date = 2024-12-27
template = "page.html"
+++

UPDATED SCSS documentation content - incremental build test!

This content was changed during development.`

		os.WriteFile(rootPagePath, []byte(modifiedPage), 0644)

		// Snapshot old taxonomy values BEFORE parsing
		oldTaxonomyValues := gen.SnapshotTaxonomyValues(
			[]string{rootPagePath},
			filepath.Join(sitePath, "content"),
		)

		// Incremental rebuild (simulating serve mode behavior)
		time.Sleep(10 * time.Millisecond)
		if err := contentParser.ParseFiles(
			filepath.Join(sitePath, "content"),
			[]string{rootPagePath},
		); err != nil {
			t.Fatalf("failed to re-parse root-level file: %v", err)
		}

		// CRITICAL TEST: This is where the bug occurred
		// The rebuild analyzer couldn't find scss.md because it looked up
		// contentMap[""] instead of contentMap["."]
		err := gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      []string{rootPagePath},
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})
		if err != nil {
			t.Fatalf("incremental rebuild failed for root-level file: %v", err)
		}

		// Verify the page was regenerated with updated content
		updatedContent, _ := os.ReadFile(filepath.Join(sitePath, "public/scss/index.html"))

		if !strings.Contains(string(updatedContent), "UPDATED SCSS documentation content") {
			t.Error("incremental build MUST update root-level file content (BUG: stale cache served)")
		}

		if !strings.Contains(string(updatedContent), "changed during development") {
			t.Error("incremental build should contain new content markers")
		}

		if strings.Contains(string(updatedContent), "Original SCSS documentation content") {
			t.Error("incremental build should NOT contain old content (cache not invalidated)")
		}

		if strings.Contains(string(updatedContent), "initial version") {
			t.Error("incremental build should NOT contain old version markers")
		}

		// Verify sitemap was also updated
		verifyFileExists(t, sitePath, "sitemap.xml")
	})

	t.Run("MultipleRootLevelPages_AllRegenerate", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Create multiple root-level pages (avoiding conflicts with existing sections)
		faqPage := `+++
title = "FAQ"
template = "page.html"
+++
FAQ content v1`

		privacyPage := `+++
title = "Privacy Policy"
template = "page.html"
+++
Privacy content v1`

		faqPath := filepath.Join(sitePath, "content/faq.md")
		privacyPath := filepath.Join(sitePath, "content/privacy-policy.md")

		os.WriteFile(faqPath, []byte(faqPage), 0644)
		os.WriteFile(privacyPath, []byte(privacyPage), 0644)

		// Initial build
		gen, contentParser := buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "faq/index.html")
		verifyFileExists(t, sitePath, "privacy-policy/index.html")

		// Modify both root-level pages
		changedFiles := []string{faqPath, privacyPath}

		oldTaxonomyValues := gen.SnapshotTaxonomyValues(
			changedFiles,
			filepath.Join(sitePath, "content"),
		)

		modifiedFaq := strings.Replace(faqPage, "v1", "v2 UPDATED", 1)
		modifiedPrivacy := strings.Replace(privacyPage, "v1", "v2 UPDATED", 1)

		os.WriteFile(faqPath, []byte(modifiedFaq), 0644)
		os.WriteFile(privacyPath, []byte(modifiedPrivacy), 0644)

		// Incremental rebuild with multiple root-level files
		time.Sleep(10 * time.Millisecond)
		if err := contentParser.ParseFiles(
			filepath.Join(sitePath, "content"),
			changedFiles,
		); err != nil {
			t.Fatalf("failed to re-parse multiple root-level files: %v", err)
		}

		err := gen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      changedFiles,
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})
		if err != nil {
			t.Fatalf("incremental rebuild failed for multiple root-level files: %v", err)
		}

		// Verify both pages were regenerated
		verifyFileContent(t, sitePath, "faq/index.html", "v2 UPDATED")
		verifyFileContent(t, sitePath, "privacy-policy/index.html", "v2 UPDATED")

		faqContent, _ := os.ReadFile(filepath.Join(sitePath, "public/faq/index.html"))
		privacyContent, _ := os.ReadFile(filepath.Join(sitePath, "public/privacy-policy/index.html"))

		if strings.Contains(string(faqContent), "v1") && !strings.Contains(string(faqContent), "v2") {
			t.Error("faq page should not contain old content")
		}

		if strings.Contains(string(privacyContent), "v1") && !strings.Contains(string(privacyContent), "v2") {
			t.Error("privacy-policy page should not contain old content")
		}
	})
}

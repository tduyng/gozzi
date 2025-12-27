package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tduyng/gozzi/app/builder"
)

// TestHomepageCache_WithConfiguredSections tests that homepage cache invalidates
// only when configured sections change.
func TestHomepageCache_WithConfiguredSections(t *testing.T) {
	t.Run("BlogPostChange_InvalidatesHomepageCache", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Add homepage_cache_sections config at top level (before [extra])
		configPath := filepath.Join(sitePath, "config.toml")
		insertAfterLine(t, configPath, "language_code", "homepage_cache_sections = [\"blog\"]")

		_, contentParser := buildSite(t, sitePath)

		// Create a blog post
		postPath := createPost(t, sitePath, "blog/test-post/index.md", `+++
title = "Test Post"
description = "Original description"
date = 2024-12-27
template = "post.html"
+++

Test content.
`)

		// Full rebuild
		if err := contentParser.Parse(filepath.Join(sitePath, "content")); err != nil {
			t.Fatalf("failed to parse: %v", err)
		}
		newGen, err := builder.NewBuilder(contentParser.Site, contentParser)
		if err != nil {
			t.Fatalf("failed to create builder: %v", err)
		}
		if err := newGen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("failed to build: %v", err)
		}

		// Get initial homepage modification time
		homePath := filepath.Join(sitePath, "public/index.html")
		initialStat, err := os.Stat(homePath)
		if err != nil {
			t.Fatalf("failed to stat homepage: %v", err)
		}
		initialModTime := initialStat.ModTime()

		// Wait to ensure different timestamps
		time.Sleep(100 * time.Millisecond)

		// Modify blog post description
		modifyFile(t, postPath, "Original description", "Updated description")

		// Incremental rebuild
		oldTaxonomyValues := newGen.SnapshotTaxonomyValues(
			[]string{postPath},
			filepath.Join(sitePath, "content"),
		)

		if err := contentParser.ParseFiles(filepath.Join(sitePath, "content"), []string{postPath}); err != nil {
			t.Fatalf("failed to re-parse: %v", err)
		}

		err = newGen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      []string{postPath},
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})
		if err != nil {
			t.Fatalf("incremental rebuild failed: %v", err)
		}

		// Verify homepage was regenerated (modification time changed)
		// This happens because blog is in homepage_cache_sections
		newStat, err := os.Stat(homePath)
		if err != nil {
			t.Fatalf("failed to stat homepage after rebuild: %v", err)
		}
		newModTime := newStat.ModTime()

		if initialModTime.Equal(newModTime) {
			t.Error("homepage should be regenerated when blog post changes (blog is in homepage_cache_sections)")
		}
	})

	t.Run("NonConfiguredSectionChange_DoesNotInvalidateCache", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Add homepage_cache_sections config - only blog, not notes
		configPath := filepath.Join(sitePath, "config.toml")
		insertAfterLine(t, configPath, "language_code", "homepage_cache_sections = [\"blog\"]")

		_, contentParser := buildSite(t, sitePath)

		// Create a notes page (not in homepage_cache_sections)
		notesDir := filepath.Join(sitePath, "content/notes")
		os.MkdirAll(notesDir, 0755)

		notePath := createPost(t, sitePath, "notes/test-note.md", `+++
title = "Test Note"
description = "Original note"
date = 2024-12-27
template = "note.html"
+++

Note content.
`)

		// Full rebuild
		if err := contentParser.Parse(filepath.Join(sitePath, "content")); err != nil {
			t.Fatalf("failed to parse: %v", err)
		}
		newGen, err := builder.NewBuilder(contentParser.Site, contentParser)
		if err != nil {
			t.Fatalf("failed to create builder: %v", err)
		}
		if err := newGen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("failed to build: %v", err)
		}

		// Get initial homepage modification time
		homePath := filepath.Join(sitePath, "public/index.html")
		initialStat, err := os.Stat(homePath)
		if err != nil {
			t.Fatalf("failed to stat homepage: %v", err)
		}
		initialModTime := initialStat.ModTime()

		// Wait to ensure different timestamps
		time.Sleep(100 * time.Millisecond)

		// Modify note description
		modifyFile(t, notePath, "Original note", "Updated note")

		// Incremental rebuild
		oldTaxonomyValues := newGen.SnapshotTaxonomyValues(
			[]string{notePath},
			filepath.Join(sitePath, "content"),
		)

		if err := contentParser.ParseFiles(filepath.Join(sitePath, "content"), []string{notePath}); err != nil {
			t.Fatalf("failed to re-parse: %v", err)
		}

		err = newGen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      []string{notePath},
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})
		if err != nil {
			t.Fatalf("incremental rebuild failed: %v", err)
		}

		// Verify homepage was NOT regenerated
		// This is because notes is NOT in homepage_cache_sections
		newStat, err := os.Stat(homePath)
		if err != nil {
			t.Fatalf("failed to stat homepage after rebuild: %v", err)
		}
		newModTime := newStat.ModTime()

		if !initialModTime.Equal(newModTime) {
			t.Error("homepage should NOT be regenerated when notes change (notes not in homepage_cache_sections)")
		}
	})
}

// TestHomepageCache_WithoutConfig tests fallback behavior when config is not set.
func TestHomepageCache_WithoutConfig(t *testing.T) {
	t.Run("AnySectionChange_InvalidatesHomepageCache", func(t *testing.T) {
		sitePath := setupTestSite(t)
		// Don't add homepage_cache_sections config - test fallback behavior

		_, contentParser := buildSite(t, sitePath)

		// Create posts in multiple sections
		blogPath := createPost(t, sitePath, "blog/test-post/index.md", `+++
title = "Test Post"
description = "Blog post"
date = 2024-12-27
template = "post.html"
+++
Content.
`)

		notesDir := filepath.Join(sitePath, "content/notes")
		os.MkdirAll(notesDir, 0755)
		notePath := createPost(t, sitePath, "notes/test-note.md", `+++
title = "Test Note"
description = "Original note"
date = 2024-12-27
template = "note.html"
+++
Note.
`)

		// Full rebuild
		if err := contentParser.Parse(filepath.Join(sitePath, "content")); err != nil {
			t.Fatalf("failed to parse: %v", err)
		}
		newGen, err := builder.NewBuilder(contentParser.Site, contentParser)
		if err != nil {
			t.Fatalf("failed to create builder: %v", err)
		}
		if err := newGen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("failed to build: %v", err)
		}

		// Test 1: Modify blog post
		homePath := filepath.Join(sitePath, "public/index.html")
		initialStat, _ := os.Stat(homePath)
		initialModTime := initialStat.ModTime()
		time.Sleep(100 * time.Millisecond)

		modifyFile(t, blogPath, "Blog post", "Updated blog post")

		oldTaxonomyValues := newGen.SnapshotTaxonomyValues([]string{blogPath}, filepath.Join(sitePath, "content"))
		contentParser.ParseFiles(filepath.Join(sitePath, "content"), []string{blogPath})
		newGen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      []string{blogPath},
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})

		stat1, _ := os.Stat(homePath)
		if initialModTime.Equal(stat1.ModTime()) {
			t.Error("homepage should be regenerated when blog post changes (no config = all sections)")
		}

		// Test 2: Modify note
		time.Sleep(100 * time.Millisecond)
		modTime2 := stat1.ModTime()

		modifyFile(t, notePath, "Original note", "Updated note")

		oldTaxonomyValues = newGen.SnapshotTaxonomyValues([]string{notePath}, filepath.Join(sitePath, "content"))
		contentParser.ParseFiles(filepath.Join(sitePath, "content"), []string{notePath})
		newGen.GenerateWithOptions(contentParser.ContentMap["."], builder.GenerateOptions{
			Incremental:       true,
			ChangedFiles:      []string{notePath},
			ContentDir:        filepath.Join(sitePath, "content"),
			OldTaxonomyValues: oldTaxonomyValues,
		})

		stat2, _ := os.Stat(homePath)
		if modTime2.Equal(stat2.ModTime()) {
			t.Error("homepage should be regenerated when note changes (no config = all sections)")
		}
	})
}

// TestHomepageCache_ConfigValidation tests edge cases in config.
func TestHomepageCache_ConfigValidation(t *testing.T) {
	t.Run("EmptyConfig_UsesAllSections", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Add empty homepage_cache_sections config
		configPath := filepath.Join(sitePath, "config.toml")
		insertAfterLine(t, configPath, "language_code", "homepage_cache_sections = []")

		_, contentParser := buildSite(t, sitePath)

		// Verify that empty config loads correctly
		if contentParser.Site.HomepageCacheSections == nil {
			t.Error("empty config should result in empty slice, not nil")
		}
		if len(contentParser.Site.HomepageCacheSections) != 0 {
			t.Errorf("empty config should have 0 sections, got %d", len(contentParser.Site.HomepageCacheSections))
		}
	})

	t.Run("NonExistentSection_GracefullyHandled", func(t *testing.T) {
		sitePath := setupTestSite(t)

		// Add config with non-existent section
		configPath := filepath.Join(sitePath, "config.toml")
		insertAfterLine(t, configPath, "language_code", "homepage_cache_sections = [\"nonexistent\", \"blog\"]")

		// Should not crash
		gen, contentParser := buildSite(t, sitePath)

		if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("build should succeed even with non-existent section: %v", err)
		}
	})
}

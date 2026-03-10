// Tests for alias/redirect functionality
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestAliases_PageRedirects tests that aliases generate proper redirect files for pages
func TestAliases_PageRedirects(t *testing.T) {
	t.Parallel()
	t.Run("SingleAlias_Created", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "New Post"
date = 2024-01-01
template = "post.html"
aliases = ["/old-post"]
+++

This post has moved.
`
		createPost(t, sitePath, "blog/new-post.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		// Check that the redirect file exists
		redirectPath := filepath.Join(sitePath, "public", "old-post", "index.html")
		if _, err := os.Stat(redirectPath); os.IsNotExist(err) {
			t.Fatalf("Redirect file not created at %s", redirectPath)
		}

		// Verify redirect content
		content, err := os.ReadFile(redirectPath)
		if err != nil {
			t.Fatalf("Failed to read redirect file: %v", err)
		}

		contentStr := string(content)
		t.Logf("Redirect content:\n%s", contentStr)

		if !strings.Contains(contentStr, `<meta http-equiv="refresh"`) {
			t.Error("Redirect missing meta refresh tag")
		}

		if !strings.Contains(contentStr, "/blog/new-post/") {
			t.Errorf("Redirect doesn't point to correct URL. Content: %s", contentStr)
		}

		if !strings.Contains(contentStr, `<link rel="canonical"`) {
			t.Error("Redirect missing canonical link")
		}
	})

	t.Run("MultipleAliases_AllCreated", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Renamed Post"
date = 2024-01-01
template = "post.html"
aliases = ["/old-url", "/legacy/post", "/archive/old-post"]
+++

This post has been renamed multiple times.
`
		createPost(t, sitePath, "blog/renamed.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		// Check all redirect files
		redirectPaths := []string{
			filepath.Join(sitePath, "public", "old-url", "index.html"),
			filepath.Join(sitePath, "public", "legacy", "post", "index.html"),
			filepath.Join(sitePath, "public", "archive", "old-post", "index.html"),
		}

		for _, path := range redirectPaths {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("Redirect file not created at %s", path)
			}

			// Verify all redirects point to the correct page
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read redirect file %s: %v", path, err)
			}

			if !strings.Contains(string(content), "/blog/renamed/") {
				t.Errorf("Redirect at %s doesn't point to /blog/renamed/", path)
			}
		}
	})

	t.Run("NoAliases_NoRedirects", func(t *testing.T) {
		t.Parallel()
		sitePath := setupReadOnlyTestSite(t)
		buildSite(t, sitePath)

		// Verify that blog posts without aliases don't create extra redirect files
		// post1.md has no aliases defined in testdata
		redirectPath := filepath.Join(sitePath, "public", "old-post1", "index.html")
		if _, err := os.Stat(redirectPath); !os.IsNotExist(err) {
			t.Error("Unexpected redirect file created for post without aliases")
		}
	})
}

// TestAliases_SectionRedirects tests that aliases work for section pages
func TestAliases_SectionRedirects(t *testing.T) {
	t.Parallel()
	t.Run("SectionAlias_Created", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create a section with aliases
		section := `+++
title = "Articles"
template = "blog.html"
aliases = ["/posts", "/writings"]
+++

This is the articles section.
`
		sectionPath := filepath.Join(sitePath, "content", "articles", "_index.md")
		os.MkdirAll(filepath.Dir(sectionPath), 0755)
		os.WriteFile(sectionPath, []byte(section), 0644)

		fullRebuild(t, gen, contentParser, sitePath)

		// Check redirect files
		redirectPaths := []string{
			filepath.Join(sitePath, "public", "posts", "index.html"),
			filepath.Join(sitePath, "public", "writings", "index.html"),
		}

		for _, path := range redirectPaths {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("Section redirect file not created at %s", path)
			}

			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("Failed to read redirect file %s: %v", path, err)
			}

			contentStr := string(content)
			t.Logf("Section redirect content at %s:\n%s", path, contentStr)

			if !strings.Contains(contentStr, "/articles/") {
				t.Errorf("Section redirect at %s doesn't point to /articles/. Content: %s", path, contentStr)
			}
		}
	})
}

// TestAliases_SnapshotRedirects snapshot test for redirect HTML format
func TestAliases_SnapshotRedirects(t *testing.T) {
	t.Parallel()
	t.Run("RedirectHTML_Format", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Snapshot Test"
date = 2024-01-01
template = "post.html"
aliases = ["/old-snapshot"]
+++

Testing redirect HTML format.
`
		createPost(t, sitePath, "blog/snapshot-test.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		SnapshotFiles(t, "Aliases_RedirectHTML", sitePath, []string{
			"old-snapshot/index.html",
		})
	})
}

// TestAliases_IncrementalBuild tests that aliases work correctly in incremental builds
func TestAliases_IncrementalBuild(t *testing.T) {
	t.Parallel()
	t.Run("AddAlias_IncrementalBuild", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Initially create post without alias
		postPath := filepath.Join(sitePath, "content", "blog", "evolving.md")
		post := `+++
title = "Evolving Post"
date = 2024-01-01
template = "post.html"
+++

Initial content.
`
		os.WriteFile(postPath, []byte(post), 0644)
		fullRebuild(t, gen, contentParser, sitePath)

		// Verify no redirect exists yet
		redirectPath := filepath.Join(sitePath, "public", "old-evolving", "index.html")
		if _, err := os.Stat(redirectPath); !os.IsNotExist(err) {
			t.Error("Redirect shouldn't exist yet")
		}

		// Update post to add alias
		postWithAlias := `+++
title = "Evolving Post"
date = 2024-01-01
template = "post.html"
aliases = ["/old-evolving"]
+++

Updated content with alias.
`
		os.WriteFile(postPath, []byte(postWithAlias), 0644)
		fullRebuild(t, gen, contentParser, sitePath)

		// Verify redirect was created
		if _, err := os.Stat(redirectPath); os.IsNotExist(err) {
			t.Error("Redirect file should have been created after adding alias")
		}

		content, err := os.ReadFile(redirectPath)
		if err != nil {
			t.Fatalf("Failed to read redirect file: %v", err)
		}

		if !strings.Contains(string(content), "/blog/evolving/") {
			t.Error("Redirect doesn't point to correct URL")
		}
	})
}

// TestAliases_SelfReferencingPrevention tests that aliases matching canonical permalinks are skipped
func TestAliases_SelfReferencingPrevention(t *testing.T) {
	t.Parallel()
	t.Run("DatePrefixFolder_SelfReferencing", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Simulate the exact bug scenario:
		// File: content/blog/2025-12-02-neovim-git-tools/index.md
		// Canonical URL: /blog/neovim-git-tools/ (date prefix removed)
		// Aliases: ["/blog/neovim-git", "/blog/neovim-git-tools"]
		// The second alias matches the canonical URL - should be skipped
		post := `+++
title = "Neovim git integration"
date = 2025-12-02
template = "post.html"
aliases = ["/blog/neovim-git", "/blog/neovim-git-tools"]
+++

Testing self-referencing alias prevention.
`
		// Create with date prefix folder (will be stripped for URL)
		createPost(t, sitePath, "blog/2025-12-02-neovim-git-tools/index.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		// The canonical page should exist at /blog/neovim-git-tools/
		canonicalPath := filepath.Join(sitePath, "public", "blog", "neovim-git-tools", "index.html")
		canonicalContent, err := os.ReadFile(canonicalPath)
		if err != nil {
			t.Fatalf("Canonical page not created at %s: %v", canonicalPath, err)
		}

		// Verify it's the actual post content, not a redirect
		if strings.Contains(string(canonicalContent), `<meta http-equiv="refresh"`) {
			t.Error("Canonical page should not be a redirect!")
		}

		if !strings.Contains(string(canonicalContent), "Neovim git integration") {
			t.Error("Canonical page should contain actual post content")
		}

		// The valid alias /blog/neovim-git should create a redirect
		validRedirectPath := filepath.Join(sitePath, "public", "blog", "neovim-git", "index.html")
		redirectContent, err := os.ReadFile(validRedirectPath)
		if err != nil {
			t.Fatalf("Valid redirect not created at %s: %v", validRedirectPath, err)
		}

		// Verify it's a redirect pointing to the canonical URL
		if !strings.Contains(string(redirectContent), `<meta http-equiv="refresh"`) {
			t.Error("Valid alias should create a redirect")
		}

		if !strings.Contains(string(redirectContent), "/blog/neovim-git-tools/") {
			t.Error("Redirect should point to canonical URL /blog/neovim-git-tools/")
		}
	})

	t.Run("ExactMatch_WithTrailingSlash", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Test when alias has trailing slash matching canonical
		post := `+++
title = "Test Post"
date = 2024-01-01
template = "post.html"
aliases = ["/blog/test-post/", "/old-test"]
+++

Testing alias with trailing slash.
`
		createPost(t, sitePath, "blog/test-post.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		// Canonical should be the real page
		canonicalPath := filepath.Join(sitePath, "public", "blog", "test-post", "index.html")
		canonicalContent, err := os.ReadFile(canonicalPath)
		if err != nil {
			t.Fatalf("Canonical page not found: %v", err)
		}

		if strings.Contains(string(canonicalContent), `<meta http-equiv="refresh"`) {
			t.Error("Canonical page was overwritten by self-referencing redirect")
		}

		// Valid alias should work
		validPath := filepath.Join(sitePath, "public", "old-test", "index.html")
		if _, err := os.Stat(validPath); os.IsNotExist(err) {
			t.Error("Valid alias redirect was not created")
		}
	})
}

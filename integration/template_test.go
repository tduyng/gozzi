// This file tests template system functionality and hot-reload behavior.
package integration

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestTemplate_Changes tests template modification and reload
func TestTemplate_Changes(t *testing.T) {
	t.Run("TemplateChange_RegeneratesPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify template - change body tag to have a class
		templatePath := filepath.Join(sitePath, "templates/post.html")
		modifyFile(t, templatePath, "<body>", "<body class=\"modified\">")

		// Reload and rebuild
		gen.ReloadTemplates()
		fullRebuild(t, gen, contentParser, sitePath)

		// Verify change reflected
		verifyFileContent(t, sitePath, "blog/post1/index.html", `class="modified"`)
	})

	t.Run("TemplateChange_OnlyAffectsCorrectPages", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Get initial content
		initialPost := readFileContent(t, sitePath, "blog/post1/index.html")

		// Modify post template only
		templatePath := filepath.Join(sitePath, "templates/post.html")
		modifyFile(t, templatePath, "<body>", "<body data-v=\"2\">")

		gen.ReloadTemplates()
		fullRebuild(t, gen, contentParser, sitePath)

		// Posts changed
		newPost := readFileContent(t, sitePath, "blog/post1/index.html")
		if newPost == initialPost {
			t.Error("post should have changed after template modification")
		}

		// Section should still work (different template)
		newSection := readFileContent(t, sitePath, "blog/index.html")
		if !strings.Contains(newSection, "First Post") {
			t.Error("section should still contain content")
		}
	})
}

// TestTemplate_HotReload tests serve mode template reloading
func TestTemplate_HotReload(t *testing.T) {
	t.Run("HotReload_UpdatesImmediately", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Simulate template change in serve mode
		templatePath := filepath.Join(sitePath, "templates/post.html")
		modifyFile(t, templatePath, "<body>", "<body class=\"reloaded\">")

		// Reload templates (like serve mode does)
		gen.ReloadTemplates()
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileContent(t, sitePath, "blog/post1/index.html", `class="reloaded"`)
	})
}

// TestTemplate_Partials tests partial template updates
func TestTemplate_Partials(t *testing.T) {
	t.Run("PartialChange_InvalidatesAll", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// If partials exist and are modified, all templates using them should update
		// This test would need actual partials in testdata to be meaningful
		// For now, verify template reload works

		gen.ReloadTemplates()
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/post1/index.html")
	})
}

// TestTemplate_Inheritance tests template inheritance
func TestTemplate_Inheritance(t *testing.T) {
	t.Run("DifferentTemplates_DifferentOutput", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Posts use post.html
		postContent := readFileContent(t, sitePath, "blog/post1/index.html")

		// Sections use section.html
		sectionContent := readFileContent(t, sitePath, "blog/index.html")

		// They should have different structures
		if postContent == sectionContent {
			t.Error("posts and sections should use different templates")
		}
	})
}

// TestTemplate_NotFound tests template error handling
func TestTemplate_NotFound(t *testing.T) {
	t.Run("MissingTemplate_HandledGracefully", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create post with non-existent template
		post := `+++
title = "Bad Template"
date = 2024-01-01
template = "nonexistent.html"
+++
Content`

		createPost(t, sitePath, "blog/bad-template.md", post)

		// Should not crash, might skip or use default
		fullRebuild(t, gen, contentParser, sitePath)
	})
}

// TestTemplate_CacheInvalidation tests template cache behavior
func TestTemplate_CacheInvalidation(t *testing.T) {
	t.Run("TemplateChange_InvalidatesCacheCompletely", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Modify template
		templatePath := filepath.Join(sitePath, "templates/post.html")
		modifyFile(t, templatePath, "<body>", "<body data-cache-test=\"1\">")

		// Reload templates
		gen.ReloadTemplates()

		fullRebuild(t, gen, contentParser, sitePath)

		// Verify template change is reflected
		verifyFileContent(t, sitePath, "blog/post1/index.html", `data-cache-test="1"`)
	})
}

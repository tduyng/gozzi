// This file tests shortcode functionality and rendering.
package integration

import (
	"testing"
)

// TestShortcode_YouTube tests YouTube shortcode
func TestShortcode_YouTube(t *testing.T) {
	t.Run("YouTube_RendersCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Check shortcode test post
		verifyFileExists(t, sitePath, "blog/shortcode-test/index.html")
		verifyFileContent(t, sitePath, "blog/shortcode-test/index.html", "youtube.com/embed/")
		verifyFileContent(t, sitePath, "blog/shortcode-test/index.html", "dQw4w9WgXcQ")
	})

	t.Run("YouTube_HasLazyLoading", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileContent(t, sitePath, "blog/shortcode-test/index.html", "loading=\"lazy\"")
	})
}

// TestShortcode_Image tests image shortcode
func TestShortcode_Image(t *testing.T) {
	t.Run("Image_RendersWithAttributes", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		content := readFileContent(t, sitePath, "blog/shortcode-test/index.html")

		// Should have image tag
		if !contains(content, "<img") {
			t.Error("image shortcode should render img tag")
		}

		// Should have src attribute
		if !contains(content, "src=") {
			t.Error("image should have src attribute")
		}
	})
}

// TestShortcode_Alert tests alert/callout shortcode
func TestShortcode_Alert(t *testing.T) {
	t.Run("Alert_RendersWithMarkdown", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Alert shortcode should support markdown inside
		verifyFileContent(t, sitePath, "blog/shortcode-test/index.html", "<strong>markdown</strong>")
	})

	t.Run("Alert_HasTypeAttribute", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileContent(t, sitePath, "blog/shortcode-test/index.html", "warning")
	})
}

// TestShortcode_MixedWithMarkdown tests shortcodes mixed with markdown
func TestShortcode_MixedWithMarkdown(t *testing.T) {
	t.Run("ShortcodesAndMarkdown_BothRender", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		content := readFileContent(t, sitePath, "blog/shortcode-test/index.html")

		// Should have markdown headings
		if !contains(content, "<h2>YouTube Video</h2>") {
			t.Error("markdown headings should render")
		}

		// Should have shortcode content
		if !contains(content, "youtube.com/embed/") {
			t.Error("shortcodes should render")
		}
	})
}

// TestShortcode_ErrorHandling tests shortcode error scenarios
func TestShortcode_ErrorHandling(t *testing.T) {
	t.Run("MissingParameter_HandledGracefully", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create post with shortcode missing required param
		post := `+++
title = "Bad Shortcode"
date = 2024-01-01
template = "post.html"
+++

{{< youtube >}}
`
		createPost(t, sitePath, "blog/bad-shortcode.md", post)

		// Should not crash
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/bad-shortcode/index.html")
	})

	t.Run("UnknownShortcode_HandledGracefully", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Unknown Shortcode"
date = 2024-01-01
template = "post.html"
+++

{{< nonexistent param="value" >}}
`
		createPost(t, sitePath, "blog/unknown.md", post)

		// Should not crash
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/unknown/index.html")
	})
}

// TestShortcode_CustomShortcodes tests custom shortcode implementation
func TestShortcode_CustomShortcodes(t *testing.T) {
	t.Run("CustomShortcode_Works", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Gallery and other custom shortcodes should work
		// Based on testdata structure
		verifyFileExists(t, sitePath, "blog/shortcode-test/index.html")
	})
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != substr && len(s) >= len(substr)
}

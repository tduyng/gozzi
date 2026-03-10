// Tests for shortcode functionality and rendering
package integration

import (
	"testing"
)

// TestShortcode_YouTube tests YouTube shortcode
func TestShortcode_YouTube(t *testing.T) {
	t.Parallel()
	t.Run("YouTube_RendersCorrectly", func(t *testing.T) {
		t.Parallel()
		sitePath := setupReadOnlyTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "Shortcode_YouTube", sitePath, []string{
			"blog/shortcode-test/index.html",
		})
	})
}

// TestShortcode_Image tests image shortcode
func TestShortcode_Image(t *testing.T) {
	t.Parallel()
	t.Run("Image_RendersWithAttributes", func(t *testing.T) {
		t.Parallel()
		sitePath := setupReadOnlyTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "Shortcode_Image", sitePath, []string{
			"blog/shortcode-test/index.html",
		})
	})
}

// TestShortcode_Alert tests alert/callout shortcode
func TestShortcode_Alert(t *testing.T) {
	t.Parallel()
	t.Run("Alert_RendersWithMarkdown", func(t *testing.T) {
		t.Parallel()
		sitePath := setupReadOnlyTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "Shortcode_Alert", sitePath, []string{
			"blog/shortcode-test/index.html",
		})
	})

	t.Run("Alert_HasTypeAttribute", func(t *testing.T) {
		t.Parallel()
		sitePath := setupReadOnlyTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "Shortcode_AlertType", sitePath, []string{
			"blog/shortcode-test/index.html",
		})
	})
}

// TestShortcode_MixedWithMarkdown tests shortcodes mixed with markdown
func TestShortcode_MixedWithMarkdown(t *testing.T) {
	t.Parallel()
	t.Run("ShortcodesAndMarkdown_BothRender", func(t *testing.T) {
		t.Parallel()
		sitePath := setupReadOnlyTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "Shortcode_Mixed", sitePath, []string{
			"blog/shortcode-test/index.html",
		})
	})
}

// TestShortcode_ErrorHandling tests shortcode error scenarios
func TestShortcode_ErrorHandling(t *testing.T) {
	t.Parallel()
	t.Run("MissingShortcode_PassesThrough", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Missing Shortcode Test"
date = 2024-01-01
template = "post.html"
+++

{{< nonexistent param="value" >}}
`
		createPost(t, sitePath, "blog/missing-shortcode.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		SnapshotFiles(t, "Shortcode_Missing", sitePath, []string{
			"blog/missing-shortcode/index.html",
		})
	})
}

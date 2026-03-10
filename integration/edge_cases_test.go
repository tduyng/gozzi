// This file tests edge cases and robustness scenarios.
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestEdgeCase_EmptyContent tests handling of empty or minimal content
func TestEdgeCase_EmptyContent(t *testing.T) {
	t.Parallel()
	t.Run("EmptyPost_StillGenerates", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Empty Post"
date = 2024-01-01
template = "post.html"
+++
`
		createPost(t, sitePath, "blog/empty.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/empty/index.html")
		verifyFileContent(t, sitePath, "blog/empty/index.html", "Empty Post")
	})

	t.Run("WhitespaceOnly_HandledCorrectly", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Whitespace Post"
date = 2024-01-01
template = "post.html"
+++

   
   
`
		createPost(t, sitePath, "blog/whitespace.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/whitespace/index.html")
	})
}

// TestEdgeCase_LongContent tests handling of very long content
func TestEdgeCase_LongContent(t *testing.T) {
	t.Parallel()
	t.Run("VeryLongContent_Handles", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		longContent := strings.Repeat("Lorem ipsum dolor sit amet. ", 10000)
		post := `+++
title = "Long Post"
date = 2024-01-01
template = "post.html"
+++

` + longContent

		createPost(t, sitePath, "blog/long.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/long/index.html")
		verifyFileContent(t, sitePath, "blog/long/index.html", "Long Post")
	})
}

// TestEdgeCase_SpecialCharacters tests filename sanitization
func TestEdgeCase_SpecialCharacters(t *testing.T) {
	t.Parallel()
	t.Run("SpecialCharsInFilename_Sanitized", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Special Post"
date = 2024-01-01
template = "post.html"
+++

Content.`

		createPost(t, sitePath, "blog/special-&-chars.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		// Should sanitize to safe slug
		verifyFileExists(t, sitePath, "blog/special-chars/index.html")
	})

	t.Run("UnicodeFilename_HandledCorrectly", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Unicode Post"
date = 2024-01-01
template = "post.html"
+++

Content.`

		createPost(t, sitePath, "blog/文章.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		// Should handle unicode filenames
		// Output path depends on implementation
		verifyFileExists(t, sitePath, "blog/index.html")
	})
}

// TestEdgeCase_UnicodeContent tests unicode in content
func TestEdgeCase_UnicodeContent(t *testing.T) {
	t.Parallel()
	t.Run("Chinese_RendersCorrectly", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "中文标题"
date = 2024-01-01
template = "post.html"
+++

这是中文内容。你好世界！`

		createPost(t, sitePath, "blog/chinese.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/chinese/index.html")
		verifyFileContent(t, sitePath, "blog/chinese/index.html", "中文")
		verifyFileContent(t, sitePath, "blog/chinese/index.html", "你好世界")
	})

	t.Run("Japanese_RendersCorrectly", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "日本語タイトル"
date = 2024-01-01
template = "post.html"
+++

これは日本語です。こんにちは！`

		createPost(t, sitePath, "blog/japanese.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/japanese/index.html")
		verifyFileContent(t, sitePath, "blog/japanese/index.html", "日本語")
	})

	t.Run("Arabic_RendersCorrectly", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "عربي"
date = 2024-01-01
template = "post.html"
+++

مرحبا بالعالم`

		createPost(t, sitePath, "blog/arabic.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/arabic/index.html")
		verifyFileContent(t, sitePath, "blog/arabic/index.html", "مرحبا")
	})

	t.Run("Emoji_RendersCorrectly", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Emoji Post 🚀"
date = 2024-01-01
template = "post.html"
+++

Hello 👋 World 🌍!`

		createPost(t, sitePath, "blog/emoji.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/emoji/index.html")
		verifyFileContent(t, sitePath, "blog/emoji/index.html", "🚀")
		verifyFileContent(t, sitePath, "blog/emoji/index.html", "👋")
	})
}

// TestEdgeCase_NestedSections tests deep nesting
func TestEdgeCase_NestedSections(t *testing.T) {
	t.Parallel()
	t.Run("DeepNesting_Works", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create deeply nested post
		os.MkdirAll(filepath.Join(sitePath, "content/blog/a/b/c/d"), 0755)
		post := `+++
title = "Deep Post"
date = 2024-01-01
template = "post.html"
+++

Deep content.`

		os.WriteFile(filepath.Join(sitePath, "content/blog/a/b/c/d/deep.md"), []byte(post), 0644)
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/a/b/c/d/deep/index.html")
	})
}

// TestEdgeCase_MissingFrontmatter tests posts with incomplete frontmatter
func TestEdgeCase_MissingFrontmatter(t *testing.T) {
	t.Parallel()
	t.Run("NoDate_UsesDefault", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "No Date Post"
template = "post.html"
+++

Content.`

		createPost(t, sitePath, "blog/no-date.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/no-date/index.html")
	})

	t.Run("InvalidDate_HandledGracefully", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Invalid Date"
date = "not-a-date"
template = "post.html"
+++

Content.`

		createPost(t, sitePath, "blog/bad-date.md", post)

		// Should not crash
		fullRebuild(t, gen, contentParser, sitePath)
	})
}

// TestEdgeCase_MalformedContent tests robustness
func TestEdgeCase_MalformedContent(t *testing.T) {
	t.Parallel()
	t.Run("MissingClosingFrontmatter_Handled", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Malformed"
date = 2024-01-01
template = "post.html"

Content without closing +++`

		createPost(t, sitePath, "blog/malformed.md", post)

		// Should handle gracefully (might skip file or use defaults)
		fullRebuild(t, gen, contentParser, sitePath)
	})
}

// TestEdgeCase_BinaryFiles tests handling of non-text files
func TestEdgeCase_BinaryFiles(t *testing.T) {
	t.Parallel()
	t.Run("BinaryInContent_Ignored", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Put binary file in content directory
		binaryPath := filepath.Join(sitePath, "content/blog/image.png")
		os.WriteFile(binaryPath, []byte{0x89, 0x50, 0x4E, 0x47}, 0644)

		// Should not crash
		fullRebuild(t, gen, contentParser, sitePath)
	})
}

// TestEdgeCase_LargeNumberOfFiles tests scalability
func TestEdgeCase_LargeNumberOfFiles(t *testing.T) {
	t.Parallel()
	t.Run("ManyPosts_HandlesCorrectly", func(t *testing.T) {
		t.Parallel()
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create 100 posts
		for i := 1; i <= 100; i++ {
			post := `+++
title = "Post ` + string(rune('0'+i%10)) + `"
date = 2024-01-01
template = "post.html"
+++

Content ` + string(rune('0'+i%10)) + `.`

			createPost(t, sitePath, "blog/auto-"+string(rune('0'+i%10))+string(rune('0'+(i/10)%10))+".md", post)
		}

		fullRebuild(t, gen, contentParser, sitePath)

		// Should complete without crashing
		verifyFileExists(t, sitePath, "blog/index.html")
	})
}

// Tests for all page types and content rendering scenarios
package integration

import (
	"path/filepath"
	"testing"
)

// TestContent_PageTypes tests different page type generation
func TestContent_PageTypes(t *testing.T) {
	t.Run("Homepage_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "PageTypes_Homepage", sitePath, []string{
			"index.html",
		})
	})

	t.Run("SinglePages_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "PageTypes_SinglePages", sitePath, []string{
			"about/index.html",
			"contact/index.html",
			"uses/index.html",
			"privacy/index.html",
		})
	})

	t.Run("BlogPosts_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "PageTypes_BlogPosts", sitePath, []string{
			"blog/post1/index.html",
			"blog/post2/index.html",
			"blog/post3/index.html",
		})
	})

	t.Run("Notes_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "PageTypes_Notes", sitePath, []string{
			"notes/note1/index.html",
		})
	})

	t.Run("SectionPages_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "PageTypes_Sections", sitePath, []string{
			"blog/index.html",
			"notes/index.html",
		})
	})
}

// TestContent_MarkdownRendering tests markdown to HTML conversion
func TestContent_MarkdownRendering(t *testing.T) {
	t.Run("Headings_Rendered", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Headings Test"
date = 2024-01-01
template = "post.html"
+++

# Heading 1
## Heading 2
### Heading 3
`
		createPost(t, sitePath, "blog/headings.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		SnapshotFiles(t, "MarkdownRendering_Headings", sitePath, []string{
			"blog/headings/index.html",
		})
	})

	t.Run("Links_Rendered", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Links Test"
date = 2024-01-01
template = "post.html"
+++

[Link Text](https://example.com)
`
		createPost(t, sitePath, "blog/links.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		SnapshotFiles(t, "MarkdownRendering_Links", sitePath, []string{
			"blog/links/index.html",
		})
	})

	t.Run("CodeBlocks_Rendered", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Code Test"
date = 2024-01-01
template = "post.html"
+++

` + "```go\nfunc main() {}\n```" + `
`
		createPost(t, sitePath, "blog/code.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		SnapshotFiles(t, "MarkdownRendering_CodeBlocks", sitePath, []string{
			"blog/code/index.html",
		})
	})

	t.Run("Lists_Rendered", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Lists Test"
date = 2024-01-01
template = "post.html"
+++

- Item 1
- Item 2
- Item 3
`
		createPost(t, sitePath, "blog/lists.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		SnapshotFiles(t, "MarkdownRendering_Lists", sitePath, []string{
			"blog/lists/index.html",
		})
	})

	t.Run("Emphasis_Rendered", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Emphasis Test"
date = 2024-01-01
template = "post.html"
+++

**Bold text** and *italic text*
`
		createPost(t, sitePath, "blog/emphasis.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		SnapshotFiles(t, "MarkdownRendering_Emphasis", sitePath, []string{
			"blog/emphasis/index.html",
		})
	})
}

// TestContent_Frontmatter tests frontmatter extraction and usage
func TestContent_Frontmatter(t *testing.T) {
	t.Run("Title_ExtractedCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "Frontmatter_Title", sitePath, []string{
			"blog/post1/index.html",
			"blog/post2/index.html",
		})
	})

	t.Run("Date_FormattedCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "Frontmatter_Date", sitePath, []string{
			"blog/post1/index.html",
		})
	})

	t.Run("ExtraFields_Available", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Extra Fields Test"
date = 2024-01-01
template = "post.html"

[extra]
custom_field = "custom value"
+++

Content
`
		createPost(t, sitePath, "blog/extra.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		SnapshotFiles(t, "Frontmatter_ExtraFields", sitePath, []string{
			"blog/extra/index.html",
		})
	})

	t.Run("MissingFields_UseDefaults", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Minimal Post"
date = 2024-01-01
template = "post.html"
+++

Minimal content
`
		createPost(t, sitePath, "blog/minimal.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		SnapshotFiles(t, "Frontmatter_Defaults", sitePath, []string{
			"blog/minimal/index.html",
		})
	})
}

// TestContent_AuxiliaryPages tests sitemap, feed, and other auxiliary content
func TestContent_AuxiliaryPages(t *testing.T) {
	t.Run("Sitemap_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "Auxiliary_Sitemap", sitePath, []string{
			"sitemap.xml",
		})
	})

	t.Run("AtomFeed_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "Auxiliary_AtomFeed", sitePath, []string{
			"atom.xml",
		})
	})

	t.Run("404Page_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "Auxiliary_404", sitePath, []string{
			"404.html",
		})
	})
}

// TestContent_Sorting tests that content is sorted correctly
func TestContent_Sorting(t *testing.T) {
	t.Run("BlogPosts_SortedByDate", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create posts with different dates
		for i, date := range []string{"2024-01-10", "2024-01-05", "2024-01-15"} {
			post := `+++
title = "Post ` + string(rune('A'+i)) + `"
date = ` + date + `
template = "post.html"
+++
Content ` + string(rune('A'+i))
			createPost(t, sitePath, filepath.Join("blog", "dated-"+string(rune('a'+i))+".md"), post)
		}

		fullRebuild(t, gen, contentParser, sitePath)

		// Blog index should show posts in reverse chronological order
		SnapshotFiles(t, "Sorting_BlogByDate", sitePath, []string{
			"blog/index.html",
		})
	})

	t.Run("SectionListings_OrderedCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		SnapshotFiles(t, "Sorting_SectionListings", sitePath, []string{
			"blog/index.html",
			"notes/index.html",
		})
	})
}

// TestContent_DraftHandling tests draft post behavior
func TestContent_DraftHandling(t *testing.T) {
	t.Run("DraftPosts_ExcludedByDefault", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Draft post should not be generated
		verifyFileNotExists(t, sitePath, "blog/draft-post/index.html")

		// Draft post should not appear in listings
		SnapshotFiles(t, "Draft_ExcludedFromListings", sitePath, []string{
			"blog/index.html",
		})
	})
}

// This file tests all page types and content rendering scenarios.
package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestContent_PageTypes tests different page type generation
func TestContent_PageTypes(t *testing.T) {
	t.Run("Homepage_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "index.html")
		verifyFileContent(t, sitePath, "index.html", "<!DOCTYPE html>")
	})

	t.Run("SinglePages_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		pages := []string{"about", "contact", "uses", "privacy"}
		for _, page := range pages {
			verifyFileExists(t, sitePath, page+"/index.html")
		}
	})

	t.Run("BlogPosts_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		posts := []string{"post1", "post2", "post3"}
		for _, post := range posts {
			verifyFileExists(t, sitePath, "blog/"+post+"/index.html")
		}
	})

	t.Run("Notes_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "notes/note1/index.html")
		verifyFileContent(t, sitePath, "notes/note1/index.html", "This is a simple note in the notes section.")
	})

	t.Run("SectionPages_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "blog/index.html")
		verifyFileExists(t, sitePath, "notes/index.html")
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

		// Snapshot testing - captures entire HTML output
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

		verifyFileContent(t, sitePath, "blog/links/index.html", `<a href="https://example.com">Link Text</a>`)
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

		// Code blocks are rendered with syntax highlighting (contains style attributes)
		verifyFileContent(t, sitePath, "blog/code/index.html", "<pre style=")
		verifyFileContent(t, sitePath, "blog/code/index.html", ">main</span>")
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

		verifyFileContent(t, sitePath, "blog/lists/index.html", "<ul>")
		verifyFileContent(t, sitePath, "blog/lists/index.html", "<li>Item 1</li>")
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

		verifyFileContent(t, sitePath, "blog/emphasis/index.html", "<strong>Bold text</strong>")
		verifyFileContent(t, sitePath, "blog/emphasis/index.html", "<em>italic text</em>")
	})
}

// TestContent_Frontmatter tests frontmatter extraction and usage
func TestContent_Frontmatter(t *testing.T) {
	t.Run("Title_ExtractedCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileContent(t, sitePath, "blog/post1/index.html", "First Post")
		verifyFileContent(t, sitePath, "blog/post2/index.html", "Second Post")
	})

	t.Run("Date_FormattedCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Dates should appear in output (exact format depends on templates)
		verifyFileContent(t, sitePath, "blog/post1/index.html", "2024")
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

		verifyFileExists(t, sitePath, "blog/extra/index.html")
	})

	t.Run("MissingFields_UseDefaults", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Minimal Post"
date = 2024-01-01
template = "post.html"
+++

Content
`
		createPost(t, sitePath, "blog/minimal.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/minimal/index.html")
		verifyFileContent(t, sitePath, "blog/minimal/index.html", "Minimal Post")
	})
}

// TestContent_SectionPages tests section page rendering (CRITICAL - was a bug)
func TestContent_SectionPages(t *testing.T) {
	t.Run("SectionWithIndex_RendersContent", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// CRITICAL: Sections with _index.md must render their markdown content
		verifyFileContent(t, sitePath, "about/index.html", "About This Test Site")
		verifyFileContent(t, sitePath, "about/index.html", "integration testing")
	})

	t.Run("SectionWithoutIndex_ListsPosts", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Blog section should list posts
		verifyFileContent(t, sitePath, "blog/index.html", "First Post")
		verifyFileContent(t, sitePath, "blog/index.html", "Second Post")
		verifyFileContent(t, sitePath, "blog/index.html", "Third Post")
	})

	t.Run("NestedSections_Work", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create nested section
		os.MkdirAll(filepath.Join(sitePath, "content/blog/category"), 0755)
		nestedIndex := `+++
title = "Category Section"
+++

Category description.`
		os.WriteFile(filepath.Join(sitePath, "content/blog/category/_index.md"), []byte(nestedIndex), 0644)

		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/category/index.html")
		verifyFileContent(t, sitePath, "blog/category/index.html", "Category Section")
	})
}

// TestContent_BlogListing tests blog section listing
func TestContent_BlogListing(t *testing.T) {
	t.Run("AllPosts_Listed", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		content := readFileContent(t, sitePath, "blog/index.html")

		requiredPosts := []string{"First Post", "Second Post", "Third Post"}
		for _, post := range requiredPosts {
			if !strings.Contains(content, post) {
				t.Errorf("blog listing missing post: %s", post)
			}
		}
	})

	t.Run("PostImages_Shown", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileContent(t, sitePath, "blog/index.html", "/img/post1.jpg")
		verifyFileContent(t, sitePath, "blog/index.html", "/img/default.jpg")
	})

	t.Run("PostMetadata_Shown", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Dates should appear
		verifyFileContent(t, sitePath, "blog/index.html", "2024")
	})

	t.Run("DraftsExcluded_FromListing", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		content := readFileContent(t, sitePath, "blog/index.html")
		if strings.Contains(content, "Draft Post") {
			t.Error("blog listing should not show draft posts")
		}
	})
}

// TestContent_DraftHandling tests draft exclusion
func TestContent_DraftHandling(t *testing.T) {
	t.Run("Drafts_NotGenerated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileNotExists(t, sitePath, "blog/draft-post/index.html")
	})

	t.Run("Drafts_NotInListings", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		draft := `+++
title = "New Draft"
date = 2024-01-01
template = "post.html"
draft = true
+++
Content`
		createPost(t, sitePath, "blog/new-draft.md", draft)
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileNotExists(t, sitePath, "blog/new-draft/index.html")
		verifyFileNotContains(t, sitePath, "blog/index.html", "New Draft")
	})
}

// TestContent_NestedSections tests nested content structure
func TestContent_NestedSections(t *testing.T) {
	t.Run("NestedPost_GeneratedCorrectly", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create nested post
		os.MkdirAll(filepath.Join(sitePath, "content/blog/nested"), 0755)
		nestedPost := `+++
title = "Nested Post"
date = 2024-01-01
template = "post.html"
+++

Nested content.`
		os.WriteFile(filepath.Join(sitePath, "content/blog/nested/post.md"), []byte(nestedPost), 0644)

		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/nested/post/index.html")
		verifyFileContent(t, sitePath, "blog/nested/post/index.html", "Nested Post")
	})

	t.Run("DeepNesting_Works", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create deeply nested structure
		os.MkdirAll(filepath.Join(sitePath, "content/blog/a/b/c"), 0755)
		deepPost := `+++
title = "Deep Post"
date = 2024-01-01
template = "post.html"
+++

Deep content.`
		os.WriteFile(filepath.Join(sitePath, "content/blog/a/b/c/deep.md"), []byte(deepPost), 0644)

		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/a/b/c/deep/index.html")
	})
}

// TestContent_Metadata tests metadata propagation
func TestContent_Metadata(t *testing.T) {
	t.Run("OGImage_CustomImage", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileContent(t, sitePath, "blog/post1/index.html", "/img/post1.jpg")
	})

	t.Run("OGImage_DefaultFallback", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileContent(t, sitePath, "blog/post2/index.html", "/img/default.jpg")
	})

	t.Run("Title_InMetaTags", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileContent(t, sitePath, "blog/post1/index.html", "First Post")
	})

	t.Run("Description_InMetaTags", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		post := `+++
title = "Meta Test"
date = 2024-01-01
template = "post.html"

[extra]
description = "Custom description here"
+++

Content`
		createPost(t, sitePath, "blog/meta.md", post)
		fullRebuild(t, gen, contentParser, sitePath)

		verifyFileExists(t, sitePath, "blog/meta/index.html")
	})
}

// TestContent_Permalinks tests URL generation
func TestContent_Permalinks(t *testing.T) {
	t.Run("BlogPost_CorrectPermalink", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Posts should be at /blog/slug/
		verifyFileExists(t, sitePath, "blog/post1/index.html")
		verifyFileExists(t, sitePath, "blog/post2/index.html")
	})

	t.Run("SinglePage_CorrectPermalink", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Single pages at /page-name/
		verifyFileExists(t, sitePath, "about/index.html")
		verifyFileExists(t, sitePath, "contact/index.html")
	})

	t.Run("Section_CorrectPermalink", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Sections at /section/
		verifyFileExists(t, sitePath, "blog/index.html")
		verifyFileExists(t, sitePath, "notes/index.html")
	})

	t.Run("InternalLinks_Correct", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		// Blog listing should link to posts correctly
		verifyFileContent(t, sitePath, "blog/index.html", "/blog/post1/")
		verifyFileContent(t, sitePath, "blog/index.html", "/blog/post2/")
	})
}

// TestContent_AuxiliaryPages tests generated auxiliary pages
func TestContent_AuxiliaryPages(t *testing.T) {
	t.Run("404Page_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "404.html")
	})

	t.Run("RobotsTxt_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "robots.txt")
		verifyFileContent(t, sitePath, "robots.txt", "User-agent")
	})

	t.Run("Sitemap_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "sitemap.xml")
		verifyFileContent(t, sitePath, "sitemap.xml", "<?xml")
		verifyFileContent(t, sitePath, "sitemap.xml", "https://test.example.com")
	})

	t.Run("AtomFeed_Generated", func(t *testing.T) {
		sitePath := setupTestSite(t)
		buildSite(t, sitePath)

		verifyFileExists(t, sitePath, "atom.xml")
		verifyFileContent(t, sitePath, "atom.xml", "<?xml")
	})
}

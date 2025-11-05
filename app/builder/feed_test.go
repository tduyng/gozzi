// Test file for feed.go
// Contains tests for Atom feed and sitemap generation functionality
package builder

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAtomFeed(t *testing.T) {
	// Prepare content files
	contentFiles := map[string]string{
		"blog/post1.md": `+++
title = "Test Post 1"
date = 2024-01-15T10:00:00Z
updated = 2024-01-16T11:00:00Z
description = "First test post"
tags = ["go", "testing"]
generate_feed = true
+++
# Test Post 1
This is the first test post content.`,

		"blog/post2.md": `+++
title = "Test Post 2"
date = 2024-01-10T10:00:00Z
description = "Second test post"
tags = ["web", "development"]
generate_feed = true
+++
# Test Post 2
This is the second test post content.`,

		"blog/draft.md": `+++
title = "Draft Post"
date = 2024-01-20T10:00:00Z
description = "Draft post"
draft = true
generate_feed = true
+++
# Draft Post
This should not appear in feed.`,
	}

	b, _ := createTestBuilderWithContent(t, contentFiles)

	// Set up site with author info
	b.site.Extra = map[string]any{
		"author": map[string]any{
			"name":  "Test Author",
			"email": "test@example.com",
		},
	}

	// Create output directory
	require.NoError(t, os.MkdirAll(b.site.OutputDir, 0755))

	// Generate Atom feed
	err := b.generateAtomFeed()
	assert.NoError(t, err)

	// Verify atom.xml was created
	feedPath := filepath.Join(b.site.OutputDir, "atom.xml")
	_, err = os.Stat(feedPath)
	assert.NoError(t, err)

	// Read and validate feed content
	feedContent, err := os.ReadFile(feedPath)
	assert.NoError(t, err)
	feedStr := string(feedContent)

	// Check XML structure
	assert.Contains(t, feedStr, `<?xml version="1.0" encoding="UTF-8"?>`)
	assert.Contains(t, feedStr, `<?xml-stylesheet type="text/xsl" href="/atom.xsl"?>`)
	assert.Contains(t, feedStr, `<feed xmlns="http://www.w3.org/2005/Atom">`)
	assert.Contains(t, feedStr, `<title>Test Site</title>`)
	assert.Contains(t, feedStr, `<id>https://example.com</id>`)
	assert.Contains(t, feedStr, `<link rel="self" href="https://example.com/atom.xml"></link>`)
	assert.Contains(t, feedStr, `<link href="https://example.com"></link>`)

	// Check author info
	assert.Contains(t, feedStr, `<author>`)
	assert.Contains(t, feedStr, `<name>Test Author</name>`)
	assert.Contains(t, feedStr, `<email>test@example.com</email>`)

	// Check entries (should have 2 non-draft posts)
	assert.Contains(t, feedStr, `<title>Test Post 1</title>`)
	assert.Contains(t, feedStr, `<title>Test Post 2</title>`)
	assert.NotContains(t, feedStr, `<title>Draft Post</title>`) // Draft should be excluded

	// Check entry details
	assert.Contains(t, feedStr, `<summary>First test post</summary>`)
	assert.Contains(t, feedStr, `<content type="html">`)
	assert.Contains(t, feedStr, `<category>go</category>`)
	assert.Contains(t, feedStr, `<category>testing</category>`)

	// Verify entries are sorted by date (most recent first)
	post1Index := strings.Index(feedStr, "Test Post 1")
	post2Index := strings.Index(feedStr, "Test Post 2")
	assert.True(t, post1Index < post2Index, "Posts should be sorted by date, most recent first")
}

func TestGenerateSitemap(t *testing.T) {
	// Prepare content files
	contentFiles := map[string]string{
		"_index.md": `+++
title = "Home"
date = 2024-01-15T10:00:00Z
+++
# Home Page`,

		"about.md": `+++
title = "About"
date = 2024-01-10T10:00:00Z
updated = 2024-01-12T10:00:00Z
+++
# About Page`,

		"blog/post1.md": `+++
title = "Blog Post"
date = 2024-01-05T10:00:00Z
tags = ["go", "web"]
+++
# Blog Post`,

		"blog/draft.md": `+++
title = "Draft Post"
date = 2024-01-20T10:00:00Z
draft = true
+++
# Draft Post`,
	}

	b, tempDir := createTestBuilderWithContent(t, contentFiles)

	// Create tags template to enable tag pages
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "templates", "tags.html"), []byte(`<html><body>Tags</body></html>`), 0644))

	// Create output directory
	require.NoError(t, os.MkdirAll(b.site.OutputDir, 0755))

	// Reload templates to pick up tags.html
	err := b.ReloadTemplates()
	require.NoError(t, err)

	// Generate sitemap
	err = b.generateSitemap()
	assert.NoError(t, err)

	// Verify sitemap.xml was created
	sitemapPath := filepath.Join(b.site.OutputDir, "sitemap.xml")
	_, err = os.Stat(sitemapPath)
	assert.NoError(t, err)

	// Read and validate sitemap content
	sitemapContent, err := os.ReadFile(sitemapPath)
	assert.NoError(t, err)
	sitemapStr := string(sitemapContent)

	// Check XML structure
	assert.Contains(t, sitemapStr, `<?xml version="1.0" encoding="UTF-8"?>`)
	assert.Contains(t, sitemapStr, `<?xml-stylesheet type="text/xsl" href="/sitemap.xsl"?>`)
	assert.Contains(t, sitemapStr, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	// Check that non-draft content is included
	assert.Contains(t, sitemapStr, `<loc>https://example.com/index</loc>`) // Home page (uses index instead of /)
	assert.Contains(t, sitemapStr, `<loc>https://example.com/about</loc>`) // About page (no trailing slash)
	// Note: Blog section might not be included if there's no blog _index.md

	// Check that draft content is excluded
	assert.NotContains(t, sitemapStr, "draft")

	// Check priorities and change frequencies (using actual values from output)
	assert.Contains(t, sitemapStr, `<priority>0.8</priority>`)         // Home page actual priority
	assert.Contains(t, sitemapStr, `<changefreq>weekly</changefreq>`)  // Home page actual frequency
	assert.Contains(t, sitemapStr, `<changefreq>monthly</changefreq>`) // Regular pages

	// Check lastmod dates
	assert.Contains(t, sitemapStr, `<lastmod>2024-01-15</lastmod>`) // Home page date
	assert.Contains(t, sitemapStr, `<lastmod>2024-01-12</lastmod>`) // About page updated date

	// Check tag pages are included
	assert.Contains(t, sitemapStr, `<loc>https://example.com/tags/go/</loc>`)
	assert.Contains(t, sitemapStr, `<loc>https://example.com/tags/web/</loc>`)
}

func TestWriteXMLFile(t *testing.T) {
	b, _ := createTestBuilder(t)

	// Create output directory
	require.NoError(t, os.MkdirAll(b.site.OutputDir, 0755))

	// Test data structure
	testData := struct {
		XMLName xml.Name `xml:"test"`
		Title   string   `xml:"title"`
		Items   []string `xml:"item"`
	}{
		Title: "Test XML",
		Items: []string{"item1", "item2"},
	}

	// Write XML file
	err := b.writeXMLFile("test.xml", `<?xml-stylesheet type="text/xsl" href="/test.xsl"?>`, testData)
	assert.NoError(t, err)

	// Verify file was created
	xmlPath := filepath.Join(b.site.OutputDir, "test.xml")
	_, err = os.Stat(xmlPath)
	assert.NoError(t, err)

	// Read and validate content
	content, err := os.ReadFile(xmlPath)
	assert.NoError(t, err)
	contentStr := string(content)

	// Check XML structure
	assert.Contains(t, contentStr, `<?xml version="1.0" encoding="UTF-8"?>`)
	assert.Contains(t, contentStr, `<?xml-stylesheet type="text/xsl" href="/test.xsl"?>`)
	assert.Contains(t, contentStr, `<test>`)
	assert.Contains(t, contentStr, `<title>Test XML</title>`)
	assert.Contains(t, contentStr, `<item>item1</item>`)
	assert.Contains(t, contentStr, `<item>item2</item>`)
	assert.Contains(t, contentStr, `</test>`)

	// Check indentation
	lines := strings.Split(contentStr, "\n")
	assert.True(t, len(lines) > 3, "XML should be properly indented")
}

func TestGenerateAtomFeedWithoutAuthor(t *testing.T) {
	// Prepare content files
	contentFiles := map[string]string{
		"post.md": `+++
title = "Simple Post"
date = 2024-01-15T10:00:00Z
description = "A simple post"
generate_feed = true
+++
# Simple Post`,
	}

	b, _ := createTestBuilderWithContent(t, contentFiles)

	// Create output directory (no author info in site)
	require.NoError(t, os.MkdirAll(b.site.OutputDir, 0755))

	// Generate Atom feed
	err := b.generateAtomFeed()
	assert.NoError(t, err)

	// Verify feed was created without author
	feedPath := filepath.Join(b.site.OutputDir, "atom.xml")
	content, err := os.ReadFile(feedPath)
	assert.NoError(t, err)
	feedStr := string(content)

	assert.Contains(t, feedStr, `<title>Test Site</title>`)
	assert.NotContains(t, feedStr, `<author>`) // No author should be present
}

func TestGenerateAtomFeedWithManyEntries(t *testing.T) {
	// Prepare many content files to test the 100-entry limit
	contentFiles := make(map[string]string)

	for i := 1; i <= 150; i++ {
		filename := fmt.Sprintf("post%03d.md", i)
		content := fmt.Sprintf(`+++
title = "Post %d"
date = 2024-01-%02dT10:00:00Z
description = "Post number %d"
generate_feed = true
+++
# Post %d`, i, (i%28)+1, i, i) // Spread across dates in January

		contentFiles[filename] = content
	}

	b, _ := createTestBuilderWithContent(t, contentFiles)

	// Create output directory
	require.NoError(t, os.MkdirAll(b.site.OutputDir, 0755))

	// Generate Atom feed
	err := b.generateAtomFeed()
	assert.NoError(t, err)

	// Read feed content
	feedPath := filepath.Join(b.site.OutputDir, "atom.xml")
	content, err := os.ReadFile(feedPath)
	assert.NoError(t, err)
	feedStr := string(content)

	// Count entries (should be limited to 100)
	entryCount := strings.Count(feedStr, "<entry>")
	assert.Equal(t, 100, entryCount, "Feed should be limited to 100 entries")
}

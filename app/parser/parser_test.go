package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
)

// Helper function for debugging.
func getKeys(m map[string]*content.Node) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestNewParser(t *testing.T) {
	tests := []struct {
		name     string
		site     *config.Site
		validate func(t *testing.T, p *ContentParser)
	}{
		{
			name: "creates parser with basic site config",
			site: &config.Site{
				Title:   "Test Site",
				BaseURL: "https://example.com",
			},
			validate: func(t *testing.T, p *ContentParser) {
				assert.NotNil(t, p.Site)
				assert.Equal(t, "Test Site", p.Site.Title)
				assert.Equal(t, "https://example.com", p.Site.BaseURL)
				assert.NotNil(t, p.ContentMap)
				assert.NotNil(t, p.Tags)
				assert.NotNil(t, p.md)
			},
		},
		{
			name: "creates parser with extra config",
			site: &config.Site{
				Title:   "Blog",
				BaseURL: "https://blog.example.com",
				Extra: map[string]any{
					"author": "John Doe",
					"img":    "/default.png",
				},
			},
			validate: func(t *testing.T, p *ContentParser) {
				assert.Equal(t, "Blog", p.Site.Title)
				assert.NotNil(t, p.Site.Extra)
				assert.Equal(t, "John Doe", p.Site.Extra["author"])
				assert.Equal(t, "/default.png", p.Site.Extra["img"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(tt.site)
			require.NotNil(t, p)
			tt.validate(t, p)
		})
	}
}

func TestGetOrCreateSection(t *testing.T) {
	site := &config.Site{
		Title:   "Test Site",
		BaseURL: "https://example.com",
	}
	p := NewParser(site)

	tests := []struct {
		name     string
		dir      string
		validate func(t *testing.T, node *content.Node)
	}{
		{
			name: "creates root section",
			dir:  ".",
			validate: func(t *testing.T, node *content.Node) {
				assert.Equal(t, ".", node.Path)
				assert.Equal(t, "", node.Slug)
				assert.Equal(t, content.NodeTypeSection, node.Type)
				assert.Nil(t, node.Parent)
			},
		},
		{
			name: "creates nested section",
			dir:  "blog",
			validate: func(t *testing.T, node *content.Node) {
				assert.Equal(t, "blog", node.Path)
				assert.Equal(t, "blog", node.Slug)
				assert.Equal(t, content.NodeTypeSection, node.Type)
				assert.NotNil(t, node.Parent)
				assert.Equal(t, ".", node.Parent.Path)
			},
		},
		{
			name: "creates deeply nested section",
			dir:  "blog/tech/2024",
			validate: func(t *testing.T, node *content.Node) {
				assert.Equal(t, "blog/tech/2024", node.Path)
				assert.Equal(t, "blog/tech/2024", node.Slug)
				assert.Equal(t, content.NodeTypeSection, node.Type)
				assert.NotNil(t, node.Parent)
				assert.Equal(t, "blog/tech", node.Parent.Path)
			},
		},
		{
			name: "returns existing section",
			dir:  "blog", // Should return the already created one
			validate: func(t *testing.T, node *content.Node) {
				assert.Equal(t, "blog", node.Path)
				assert.Equal(t, "blog", node.Slug)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := p.GetOrCreateSection(tt.dir)
			require.NotNil(t, node)
			tt.validate(t, node)

			// Verify it was stored in ContentMap
			stored, exists := p.ContentMap[tt.dir]
			assert.True(t, exists)
			assert.Equal(t, node, stored)
		})
	}
}

func TestCalculateReadStats(t *testing.T) {
	tests := []struct {
		name              string
		content           string
		expectedWordCount int
		expectedReadTime  int
	}{
		{
			name:              "empty content",
			content:           "",
			expectedWordCount: 0,
			expectedReadTime:  1, // minimum 1 minute
		},
		{
			name:              "simple text",
			content:           "Hello world this is a test",
			expectedWordCount: 6,
			expectedReadTime:  1,
		},
		{
			name:              "content with HTML tags",
			content:           "<p>Hello <strong>world</strong> this is a <em>test</em></p>",
			expectedWordCount: 6,
			expectedReadTime:  1,
		},
		{
			name:              "long content requiring multiple minutes",
			content:           strings.Repeat("word ", 440), // 440 words = 2 minutes
			expectedWordCount: 440,
			expectedReadTime:  2,
		},
		{
			name:              "content with complex HTML",
			content:           `<div class="content"><h1>Title</h1><p>This is a paragraph with <a href="#">links</a> and <code>code</code>.</p></div>`,
			expectedWordCount: 10,
			expectedReadTime:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wordCount, readTime := calculateReadStats(tt.content)
			assert.Equal(t, tt.expectedWordCount, wordCount)
			assert.Equal(t, tt.expectedReadTime, readTime)
		})
	}
}

func TestResolveImgURL(t *testing.T) {
	site := &config.Site{
		Title:   "Test Site",
		BaseURL: "https://example.com",
		Extra: map[string]any{
			"img": "/default-site.png",
		},
	}
	p := NewParser(site)

	tests := []struct {
		name        string
		frontMatter *config.FrontMatter
		slug        string
		expected    string
	}{
		{
			name: "uses frontmatter img with absolute path",
			frontMatter: &config.FrontMatter{
				Extra: map[string]any{
					"img": "/custom.png",
				},
			},
			slug:     "blog/post1",
			expected: "https://example.com/custom.png",
		},
		{
			name: "uses frontmatter img with relative path",
			frontMatter: &config.FrontMatter{
				Extra: map[string]any{
					"img": "cover.jpg",
				},
			},
			slug:     "blog/post1",
			expected: "https://example.com/blog/post1/cover.jpg",
		},
		{
			name: "falls back to site default img",
			frontMatter: &config.FrontMatter{
				Extra: map[string]any{},
			},
			slug:     "blog/post1",
			expected: "https://example.com/default-site.png",
		},
		{
			name:        "handles nil frontmatter extra",
			frontMatter: &config.FrontMatter{Extra: nil},
			slug:        "blog/post1",
			expected:    "https://example.com/default-site.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := p.resolveImgURL(tt.frontMatter, tt.slug)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildPermalink(t *testing.T) {
	tests := []struct {
		name     string
		slug     string
		expected string
	}{
		{
			name:     "empty slug returns root",
			slug:     "",
			expected: "/",
		},
		{
			name:     "index slug returns root",
			slug:     "index",
			expected: "/",
		},
		{
			name:     "simple slug",
			slug:     "about",
			expected: "/about/",
		},
		{
			name:     "nested slug",
			slug:     "blog/post1",
			expected: "/blog/post1/",
		},
		{
			name:     "slug with leading slash",
			slug:     "/blog/post1",
			expected: "/blog/post1/",
		},
		{
			name:     "slug with trailing slash",
			slug:     "blog/post1/",
			expected: "/blog/post1/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildPermalink(tt.slug)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		slug     string
		expected string
	}{
		{
			name:     "simple URL construction",
			baseURL:  "https://example.com",
			slug:     "about",
			expected: "https://example.com/about",
		},
		{
			name:     "nested slug URL",
			baseURL:  "https://blog.example.com",
			slug:     "blog/post1",
			expected: "https://blog.example.com/blog/post1",
		},
		{
			name:     "empty slug",
			baseURL:  "https://example.com",
			slug:     "",
			expected: "https://example.com/",
		},
		{
			name:     "base URL with trailing slash",
			baseURL:  "https://example.com/",
			slug:     "about",
			expected: "https://example.com//about",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildURL(tt.baseURL, tt.slug)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseTags(t *testing.T) {
	site := &config.Site{
		Title:   "Test Site",
		BaseURL: "https://example.com",
	}
	p := NewParser(site)

	// Create test page node
	pageNode := &content.Node{
		Path: "blog/test-post.md",
		Type: content.NodeTypePage,
	}

	tests := []struct {
		name         string
		frontMatter  *config.FrontMatter
		expectedTags []string
		validate     func(t *testing.T, p *ContentParser)
	}{
		{
			name: "parses single tag",
			frontMatter: &config.FrontMatter{
				Tags: []string{"go"},
			},
			expectedTags: []string{"go"},
			validate: func(t *testing.T, p *ContentParser) {
				entry, exists := p.Tags["go"]
				assert.True(t, exists)
				assert.Equal(t, 1, entry.Count)
				assert.Contains(t, entry.Pages, pageNode)
			},
		},
		{
			name: "parses multiple tags",
			frontMatter: &config.FrontMatter{
				Tags: []string{"Go", "Programming", "Tutorial"},
			},
			expectedTags: []string{"go", "programming", "tutorial"},
			validate: func(t *testing.T, p *ContentParser) {
				for _, tag := range []string{"go", "programming", "tutorial"} {
					entry, exists := p.Tags[tag]
					assert.True(t, exists, "tag %s should exist", tag)
					assert.Equal(t, 1, entry.Count)
					assert.Contains(t, entry.Pages, pageNode)
				}
			},
		},
		{
			name: "handles duplicate tags in same page",
			frontMatter: &config.FrontMatter{
				Tags: []string{"go", "Go", "GO"},
			},
			expectedTags: []string{"go"},
			validate: func(t *testing.T, p *ContentParser) {
				entry, exists := p.Tags["go"]
				assert.True(t, exists)
				assert.Equal(t, 1, entry.Count) // Should only be counted once
				assert.Contains(t, entry.Pages, pageNode)
			},
		},
		{
			name: "ignores empty tags",
			frontMatter: &config.FrontMatter{
				Tags: []string{"go", "", "  ", "programming"},
			},
			expectedTags: []string{"go", "programming"},
			validate: func(t *testing.T, p *ContentParser) {
				assert.Equal(t, 2, len(p.Tags))
				_, emptyExists := p.Tags[""]
				assert.False(t, emptyExists)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset tags for each test
			p.Tags = make(map[string]*TagEntry)

			p.parseTags(tt.frontMatter, pageNode)

			assert.Equal(t, len(tt.expectedTags), len(p.Tags))
			tt.validate(t, p)
		})
	}
}

func TestGetMarkdownProcessor(t *testing.T) {
	site := &config.Site{
		Title:   "Test Site",
		BaseURL: "https://example.com",
	}
	p := NewParser(site)

	md := p.GetMarkdownProcessor()
	assert.NotNil(t, md)

	// Test that it can process basic markdown
	var buf strings.Builder
	err := md.Convert([]byte("# Hello World"), &buf)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "<h1")
	assert.Contains(t, buf.String(), "Hello World")
}

// TestParse tests the Parse method with actual file system using TOML frontmatter.
func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, contentDir string)
		validate func(t *testing.T, p *ContentParser)
	}{
		{
			name: "parses content directory with TOML frontmatter",
			setup: func(t *testing.T, contentDir string) {
				// Create directory structure (matching gozzi's expected structure)
				require.NoError(t, os.MkdirAll(filepath.Join(contentDir, "blog", "first-post"), 0755))

				// Create _index.md for root section with TOML frontmatter
				indexContent := `+++
title = "Home"
description = "Welcome page"
+++
# Welcome to the site
This is the home page with some content.`
				require.NoError(t, os.WriteFile(filepath.Join(contentDir, "_index.md"), []byte(indexContent), 0644))

				// Create blog section index with TOML frontmatter
				blogIndexContent := `+++
title = "Blog"
description = "My blog posts"
+++
# Blog Section
All my blog posts are here.`
				require.NoError(t, os.WriteFile(filepath.Join(contentDir, "blog", "_index.md"), []byte(blogIndexContent), 0644))

				// Create a blog post with TOML frontmatter (in subdirectory)
				postContent := `+++
title = "My First Post"
date = 2024-01-15
tags = ["go", "programming", "tutorial"]
draft = false
+++
# My First Post
This is my first blog post about Go programming.

Let's learn Go together!`
				require.NoError(t, os.WriteFile(filepath.Join(contentDir, "blog", "first-post", "index.md"), []byte(postContent), 0644))
			},
			validate: func(t *testing.T, p *ContentParser) {
				// Verify sections were created
				rootSection, exists := p.ContentMap["."]
				assert.True(t, exists)
				assert.Equal(t, content.NodeTypeSection, rootSection.Type)
				assert.Contains(t, string(rootSection.Content), "Welcome to the site")
				assert.Contains(t, string(rootSection.Content), "home page")

				blogSection, exists := p.ContentMap["blog"]
				assert.True(t, exists)
				assert.Equal(t, content.NodeTypeSection, blogSection.Type)
				assert.Contains(t, string(blogSection.Content), "Blog Section")
				assert.Contains(t, string(blogSection.Content), "blog posts are here")

				// Debug logging
				t.Logf("Root section children: %d", len(rootSection.Children))
				t.Logf("Blog section children: %d", len(blogSection.Children))
				for path, node := range p.ContentMap {
					t.Logf("ContentMap[%s]: Type=%d, Children=%d", path, node.Type, len(node.Children))
				}

				// Verify blog post was parsed and added to blog section
				if len(blogSection.Children) == 0 {
					t.Fatalf("Expected 1 child in blog section, but got 0. ContentMap keys: %v", getKeys(p.ContentMap))
				}
				assert.Equal(t, 1, len(blogSection.Children))
				blogPost := blogSection.Children[0]
				assert.Equal(t, content.NodeTypePage, blogPost.Type)
				assert.Contains(t, string(blogPost.Content), "My First Post")
				assert.Contains(t, string(blogPost.Content), "Go programming")
				assert.Contains(t, string(blogPost.Content), "learn Go together")

				// Verify tags were parsed
				assert.Equal(t, 3, len(p.Tags))
				for _, tag := range []string{"go", "programming", "tutorial"} {
					tagEntry, exists := p.Tags[tag]
					assert.True(t, exists, "tag %s should exist", tag)
					assert.Equal(t, 1, tagEntry.Count)
					assert.Contains(t, tagEntry.Pages, blogPost)
				}

				// Verify pagination was built
				assert.NotNil(t, blogSection.Children)
			},
		},
		{
			name: "handles empty content directory",
			setup: func(t *testing.T, contentDir string) {
				// Create empty content directory
			},
			validate: func(t *testing.T, p *ContentParser) {
				// Should have empty ContentMap
				assert.Equal(t, 0, len(p.ContentMap))
				assert.Equal(t, 0, len(p.Tags))
			},
		},
		{
			name: "skips draft content",
			setup: func(t *testing.T, contentDir string) {
				require.NoError(t, os.MkdirAll(filepath.Join(contentDir, "blog"), 0755))

				// Create draft post
				draftContent := `+++
title = "Draft Post"
draft = true
tags = ["draft"]
+++
# Draft Post
This should not be included.`
				require.NoError(t, os.WriteFile(filepath.Join(contentDir, "blog", "draft.md"), []byte(draftContent), 0644))

				// Create published post
				publishedContent := `+++
title = "Published Post"
draft = false
tags = ["published"]
+++
# Published Post
This should be included.`
				require.NoError(t, os.WriteFile(filepath.Join(contentDir, "blog", "published.md"), []byte(publishedContent), 0644))
			},
			validate: func(t *testing.T, p *ContentParser) {
				// Should only have the root section (created for published post)
				rootSection := p.ContentMap["."]
				assert.NotNil(t, rootSection)

				// Should have only one tag (from published post)
				assert.Equal(t, 1, len(p.Tags))
				_, draftExists := p.Tags["draft"]
				assert.False(t, draftExists)
				_, publishedExists := p.Tags["published"]
				assert.True(t, publishedExists)
			},
		},
		{
			name: "handles nested directory structure",
			setup: func(t *testing.T, contentDir string) {
				// Create nested structure
				require.NoError(t, os.MkdirAll(filepath.Join(contentDir, "blog", "tech", "2024"), 0755))

				// Create nested post
				nestedContent := `+++
title = "Nested Post"
tags = ["nested", "tech"]
+++
# Nested Post
This is in a nested directory.`
				require.NoError(t, os.WriteFile(filepath.Join(contentDir, "blog", "tech", "2024", "nested-post.md"), []byte(nestedContent), 0644))
			},
			validate: func(t *testing.T, p *ContentParser) {
				// Should create all necessary sections
				rootSection := p.ContentMap["."]
				assert.NotNil(t, rootSection)

				// The nested post should be in blog/tech/2024 section (not blog/tech)
				// because the file is blog/tech/2024/nested-post.md (not index.md)
				blogTech2024Section := p.ContentMap["blog/tech/2024"]
				assert.NotNil(t, blogTech2024Section)

				// Verify the post was added to the correct parent
				assert.Equal(t, 1, len(blogTech2024Section.Children))
				nestedPost := blogTech2024Section.Children[0]
				assert.Contains(t, string(nestedPost.Content), "Nested Post")
				assert.Contains(t, string(nestedPost.Content), "nested directory")

				// Verify tags
				assert.Equal(t, 2, len(p.Tags))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory structure
			tempDir := t.TempDir()
			contentDir := filepath.Join(tempDir, "content")
			require.NoError(t, os.MkdirAll(contentDir, 0755))

			// Setup test content
			tt.setup(t, contentDir)

			// Create parser and parse
			site := &config.Site{
				Title:   "Test Site",
				BaseURL: "https://example.com",
			}
			p := NewParser(site)

			err := p.Parse(contentDir)
			assert.NoError(t, err)

			// Validate results
			tt.validate(t, p)
		})
	}
}

// TestParseSection tests the parseSection method.
func TestParseSection(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		dir      string
		validate func(t *testing.T, p *ContentParser, err error)
	}{
		{
			name: "parses section with TOML frontmatter",
			content: `+++
title = "Test Section"
description = "A test section"
[extra]
custom = "value"
+++
# Test Section
This is a test section with some content.`,
			dir: "test",
			validate: func(t *testing.T, p *ContentParser, err error) {
				assert.NoError(t, err)

				section := p.ContentMap["test"]
				assert.NotNil(t, section)
				assert.Equal(t, content.NodeTypeSection, section.Type)
				assert.Contains(t, string(section.Content), "Test Section")
				assert.Contains(t, string(section.Content), "test section")

				// Check config was merged properly
				assert.Equal(t, "Test Section", section.Config["title"])
				assert.Equal(t, "A test section", section.Config["description"])
			},
		},
		{
			name: "handles section with invalid frontmatter gracefully",
			content: `+++
invalid toml syntax [[[
+++
# Content`,
			dir: "invalid",
			validate: func(t *testing.T, p *ContentParser, err error) {
				assert.Error(t, err)
				// Should not create invalid sections
				assert.Nil(t, p.ContentMap["invalid"])
			},
		},
		{
			name: "skips draft sections",
			content: `+++
title = "Draft Section"
draft = true
+++
# Draft Section
This should be skipped.`,
			dir: "draft",
			validate: func(t *testing.T, p *ContentParser, err error) {
				assert.Error(t, err)
				// ContentMap should remain empty since no sections were added
				assert.Equal(t, 0, len(p.ContentMap))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tempDir := t.TempDir()
			sectionPath := filepath.Join(tempDir, "_index.md")
			require.NoError(t, os.WriteFile(sectionPath, []byte(tt.content), 0644))

			site := &config.Site{
				Title:   "Test Site",
				BaseURL: "https://example.com",
			}
			p := NewParser(site)

			// Test
			err := p.parseSection(sectionPath, tt.dir)

			// Validate
			tt.validate(t, p, err)
		})
	}
}

// TestParsePage tests the parsePage method.
func TestParsePage(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		dir      string
		validate func(t *testing.T, p *ContentParser, err error)
	}{
		{
			name: "parses page with TOML frontmatter and tags",
			content: `+++
title = "Test Page"
date = 2024-01-15
tags = ["go", "programming"]
[extra]
author = "Test Author"
+++
# Test Page
This is a test page with some content.

Here's a second paragraph.`,
			dir: "blog",
			validate: func(t *testing.T, p *ContentParser, err error) {
				assert.NoError(t, err)

				// Should create blog section and add page to it
				blogSection := p.ContentMap["blog"]
				assert.NotNil(t, blogSection)
				assert.Equal(t, 1, len(blogSection.Children))

				// Check the page
				page := blogSection.Children[0]
				assert.Equal(t, content.NodeTypePage, page.Type)
				assert.Contains(t, string(page.Content), "Test Page")
				assert.Contains(t, string(page.Content), "test page")
				assert.Contains(t, string(page.Content), "second paragraph")

				// Check config
				assert.Equal(t, "Test Page", page.Config["title"])

				// Check extra config is nested
				extra, exists := page.Config["extra"]
				assert.True(t, exists)
				extraMap, ok := extra.(map[string]any)
				assert.True(t, ok)
				assert.Equal(t, "Test Author", extraMap["author"])

				// Check tags were parsed
				assert.Equal(t, 2, len(p.Tags))
				for _, tag := range []string{"go", "programming"} {
					tagEntry, exists := p.Tags[tag]
					assert.True(t, exists)
					assert.Equal(t, 1, tagEntry.Count)
					assert.Contains(t, tagEntry.Pages, page)
				}

				// Check read stats
				assert.True(t, page.WordCount > 0)
				assert.True(t, page.ReadTime > 0)
			},
		},
		{
			name: "skips draft pages",
			content: `+++
title = "Draft Page"
draft = true
tags = ["draft"]
+++
# Draft Page
This should be skipped.`,
			dir: "blog",
			validate: func(t *testing.T, p *ContentParser, err error) {
				assert.Error(t, err)
				// ContentMap should remain empty since no pages or sections were added
				assert.Equal(t, 0, len(p.ContentMap))
				assert.Equal(t, 0, len(p.Tags))
			},
		},
		{
			name: "handles page without tags",
			content: `+++
title = "No Tags Page"
date = 2024-01-15
+++
# No Tags Page
This page has no tags.`,
			dir: "blog",
			validate: func(t *testing.T, p *ContentParser, err error) {
				assert.NoError(t, err)

				// Should still create the page
				parentSection := p.ContentMap["."]
				assert.NotNil(t, parentSection)
				assert.Equal(t, 1, len(parentSection.Children))

				// Should have no tags
				assert.Equal(t, 0, len(p.Tags))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tempDir := t.TempDir()
			pagePath := filepath.Join(tempDir, "test-page.md")
			require.NoError(t, os.WriteFile(pagePath, []byte(tt.content), 0644))

			site := &config.Site{
				Title:   "Test Site",
				BaseURL: "https://example.com",
			}
			p := NewParser(site)

			// Test
			err := p.parsePage(pagePath, tt.dir)

			// Validate
			tt.validate(t, p, err)
		})
	}
}

func TestTagEntry(t *testing.T) {
	// Test TagEntry functionality
	entry := &TagEntry{
		Seen: make(map[string]struct{}),
	}

	page1 := &content.Node{Path: "blog/post1.md"}
	page2 := &content.Node{Path: "blog/post2.md"}

	// Add first page
	entry.Pages = append(entry.Pages, page1)
	entry.Seen[page1.Path] = struct{}{}
	entry.Count = len(entry.Pages)

	assert.Equal(t, 1, entry.Count)
	assert.Contains(t, entry.Pages, page1)

	// Add second page
	entry.Pages = append(entry.Pages, page2)
	entry.Seen[page2.Path] = struct{}{}
	entry.Count = len(entry.Pages)

	assert.Equal(t, 2, entry.Count)
	assert.Contains(t, entry.Pages, page1)
	assert.Contains(t, entry.Pages, page2)

	// Verify seen tracking
	_, exists := entry.Seen[page1.Path]
	assert.True(t, exists)
	_, exists = entry.Seen[page2.Path]
	assert.True(t, exists)
}

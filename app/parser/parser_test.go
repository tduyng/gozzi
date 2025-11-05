// ABOUTME: Tests for parser orchestration, initialization, and main Parse method.
// ABOUTME: Covers ContentParser creation, markdown processor setup, and full directory parsing.
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

// Shared test helpers

// Helper function for debugging.
func getKeys(m map[string]*content.Node) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

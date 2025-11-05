// ABOUTME: Tests for page parsing functionality.
// ABOUTME: Covers parsePage, read stats calculation, image URL resolution, and URL building.
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

// ABOUTME: Tests for tag parsing and tracking functionality.
// ABOUTME: Covers parseTags method and TagEntry operations.
package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
)

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

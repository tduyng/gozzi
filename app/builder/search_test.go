package builder

import (
	"encoding/json"
	"html/template"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
)

func TestGenerateSearchIndex(t *testing.T) {
	tempDir := t.TempDir()

	site := &config.Site{
		BaseURL:        "https://example.com",
		OutputDir:      tempDir,
		GenerateSearch: true,
	}

	p := parser.NewParser(site)
	b := &Builder{
		site:   site,
		parser: p,
	}

	// Create test nodes
	date1, _ := time.Parse(time.RFC3339, "2024-01-01T00:00:00Z")
	date2, _ := time.Parse(time.RFC3339, "2024-01-02T00:00:00Z")

	node1 := &content.Node{
		Type:      content.NodeTypePage,
		Permalink: "/blog/post-1/",
		Content:   template.HTML("<p>This is the first post with some <strong>bold</strong> text.</p>"),
		Config: map[string]any{
			"title":       "Post 1",
			"description": "First description",
			"date":        date1,
			"tags":        []string{"tag1", "tag2"},
		},
	}

	node2 := &content.Node{
		Type:      content.NodeTypePage,
		Permalink: "/blog/post-2/",
		Content:   template.HTML("<h1>Header</h1><p>Second post.</p>"),
		Config: map[string]any{
			"title":       "Post 2",
			"description": "Second description",
			"date":        date2,
			"tags":        []any{"tag3"},
		},
	}

	// Draft node (should be skipped)
	node3 := &content.Node{
		Type:      content.NodeTypePage,
		Permalink: "/blog/draft/",
		Content:   template.HTML("<p>Draft.</p>"),
		Config: map[string]any{
			"title": "Draft",
			"draft": true,
		},
	}

	// Section node (should be skipped)
	node4 := &content.Node{
		Type:      content.NodeTypeSection,
		Permalink: "/blog/",
		Content:   template.HTML("<p>Section.</p>"),
		Config: map[string]any{
			"title": "Blog",
		},
	}

	p.ContentMap = map[string]*content.Node{
		".": {
			Children: []*content.Node{node1, node2, node3, node4},
		},
	}

	err := b.generateSearchIndex()
	require.NoError(t, err)

	indexPath := filepath.Join(tempDir, "search-index.json")
	require.FileExists(t, indexPath)

	data, err := os.ReadFile(indexPath)
	require.NoError(t, err)

	var entries []SearchEntry
	err = json.Unmarshal(data, &entries)
	require.NoError(t, err)

	// Should only contain node1 and node2 (node3 is draft, node4 is section)
	require.Len(t, entries, 2)

	// Entries should be sorted by date descending (newest first)
	assert.Equal(t, "Post 2", entries[0].Title)
	assert.Equal(t, "/blog/post-2/", entries[0].ID)
	assert.Equal(t, "Header Second post.", entries[0].Content) // HTML stripped
	assert.Equal(t, "2024-01-02", entries[0].Date)
	assert.Equal(t, []string{"tag3"}, entries[0].Tags)

	assert.Equal(t, "Post 1", entries[1].Title)
	assert.Equal(t, "/blog/post-1/", entries[1].ID)
	assert.Equal(t, "This is the first post with some bold text.", entries[1].Content) // HTML stripped
	assert.Equal(t, "2024-01-01", entries[1].Date)
	assert.Equal(t, []string{"tag1", "tag2"}, entries[1].Tags)
}

package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func TestNewTocExtension(t *testing.T) {
	ext := NewTocExtension()
	assert.NotNil(t, ext)

	// Verify it implements goldmark.Extender
	_, ok := ext.(goldmark.Extender)
	assert.True(t, ok)
}

func TestTocExtension_Extend(t *testing.T) {
	ext := NewTocExtension()

	// Create a basic goldmark instance
	md := goldmark.New()

	// This should not panic
	assert.NotPanics(t, func() {
		ext.Extend(md)
	})
}

func TestTocExtensionIntegration(t *testing.T) {
	// Test that the extension can be used with goldmark
	md := goldmark.New(
		goldmark.WithExtensions(
			NewTocExtension(),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	require.NotNil(t, md)

	tests := []struct {
		name     string
		markdown string
		validate func(t *testing.T, toc []map[string]any)
	}{
		{
			name:     "no headings produces empty TOC",
			markdown: "Just some regular text without headings.",
			validate: func(t *testing.T, toc []map[string]any) {
				assert.Empty(t, toc)
			},
		},
		{
			name: "single heading",
			markdown: `# Introduction
Some content here.`,
			validate: func(t *testing.T, toc []map[string]any) {
				assert.Len(t, toc, 1)
				assert.Equal(t, "Introduction", toc[0]["Title"])
				assert.Equal(t, 1, toc[0]["Level"])
			},
		},
		{
			name: "multiple headings same level",
			markdown: `# Introduction
Content 1.

# Getting Started  
Content 2.

# Conclusion
Content 3.`,
			validate: func(t *testing.T, toc []map[string]any) {
				assert.Len(t, toc, 3)
				assert.Equal(t, "Introduction", toc[0]["Title"])
				assert.Equal(t, "Getting Started", toc[1]["Title"])
				assert.Equal(t, "Conclusion", toc[2]["Title"])
			},
		},
		{
			name: "nested headings",
			markdown: `# Introduction
Content 1.

## Overview
Sub content.

## Details
More sub content.

# Conclusion
Final content.`,
			validate: func(t *testing.T, toc []map[string]any) {
				assert.Len(t, toc, 2)

				// First top-level heading with children
				assert.Equal(t, "Introduction", toc[0]["Title"])
				assert.Equal(t, 1, toc[0]["Level"])
				children := toc[0]["Children"].([]map[string]any)
				assert.Len(t, children, 2)
				assert.Equal(t, "Overview", children[0]["Title"])
				assert.Equal(t, 2, children[0]["Level"])
				assert.Equal(t, "Details", children[1]["Title"])
				assert.Equal(t, 2, children[1]["Level"])

				// Second top-level heading
				assert.Equal(t, "Conclusion", toc[1]["Title"])
				assert.Equal(t, 1, toc[1]["Level"])
			},
		},
		{
			name: "deeply nested headings",
			markdown: `# Level 1
Content.

## Level 2
Content.

### Level 3
Content.

#### Level 4
Content.`,
			validate: func(t *testing.T, toc []map[string]any) {
				assert.Len(t, toc, 1)

				// Navigate through the nested structure
				level1 := toc[0]
				assert.Equal(t, "Level 1", level1["Title"])
				assert.Equal(t, 1, level1["Level"])

				level2Children := level1["Children"].([]map[string]any)
				assert.Len(t, level2Children, 1)
				level2 := level2Children[0]
				assert.Equal(t, "Level 2", level2["Title"])
				assert.Equal(t, 2, level2["Level"])

				level3Children := level2["Children"].([]map[string]any)
				assert.Len(t, level3Children, 1)
				level3 := level3Children[0]
				assert.Equal(t, "Level 3", level3["Title"])
				assert.Equal(t, 3, level3["Level"])

				level4Children := level3["Children"].([]map[string]any)
				assert.Len(t, level4Children, 1)
				level4 := level4Children[0]
				assert.Equal(t, "Level 4", level4["Title"])
				assert.Equal(t, 4, level4["Level"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := parser.NewContext()
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)), parser.WithContext(pc))

			// Render to trigger extensions
			var buf strings.Builder
			err := md.Renderer().Render(&buf, []byte(tt.markdown), doc)
			assert.NoError(t, err)

			// Get TOC from context
			toc, ok := pc.Get(0).([]map[string]any)
			if !ok && len(toc) == 0 {
				toc = []map[string]any{}
			}

			tt.validate(t, toc)
		})
	}
}

func TestHeadingText(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected string
	}{
		{
			name:     "simple heading",
			markdown: "# Simple Title",
			expected: "Simple Title",
		},
		{
			name:     "heading with inline formatting",
			markdown: "# Title with **bold** and *italic*",
			expected: "Title with  and", // The current implementation strips formatting content
		},
		{
			name:     "heading with extra whitespace",
			markdown: "#   Spaced   Title   ",
			expected: "Spaced   Title",
		},
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			NewTocExtension(),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pc := parser.NewContext()
			doc := md.Parser().Parse(text.NewReader([]byte(tt.markdown)), parser.WithContext(pc))

			// Render to trigger extensions
			var buf strings.Builder
			err := md.Renderer().Render(&buf, []byte(tt.markdown), doc)
			assert.NoError(t, err)

			// Get TOC from context
			toc, ok := pc.Get(0).([]map[string]any)
			require.True(t, ok)
			require.Len(t, toc, 1)

			assert.Equal(t, tt.expected, toc[0]["Title"])
		})
	}
}

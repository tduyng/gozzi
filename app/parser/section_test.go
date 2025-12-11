// Tests for section parsing functionality.
// Covers GetOrCreateSection and parseSection methods for _index.md files.
package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
)

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
				assert.NoError(t, err) // Drafts are silently skipped, not errors
				// ContentMap should remain empty since draft section was skipped
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

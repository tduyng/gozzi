package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/parser"
)

func TestSummary_AutoGeneration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		content         string
		config          string
		expectedSummary string
		description     string
	}{
		{
			name: "Auto_Generate_TwoSentences",
			content: `+++
title = "Test Post"
date = 2024-01-15
+++

This is the first sentence of the post. This is the second sentence. This is the third sentence that should not appear.

More content here.
`,
			config: `
base_url = "https://test.example.com/"
title = "Test Site"
summary_length = 2
`,
			expectedSummary: "This is the first sentence of the post. This is the second sentence.",
			description:     "Should auto-generate summary from first 2 sentences",
		},
		{
			name: "Manual_Override_WithDescription",
			content: `+++
title = "Test Post"
date = 2024-01-15
description = "Custom manual summary"
+++

This is the first sentence. This is the second sentence. This is the third sentence.
`,
			config: `
base_url = "https://test.example.com/"
title = "Test Site"
summary_length = 2
`,
			expectedSummary: "Custom manual summary",
			description:     "Should use description as manual override",
		},
		{
			name: "Auto_Generate_OneSentence",
			content: `+++
title = "Test Post"
date = 2024-01-15
+++

This is a single sentence post. Additional content follows here.
`,
			config: `
base_url = "https://test.example.com/"
title = "Test Site"
summary_length = 1
`,
			expectedSummary: "This is a single sentence post.",
			description:     "Should extract only first sentence when configured",
		},
		{
			name: "Auto_Generate_WithFormatting",
			content: `+++
title = "Test Post"
date = 2024-01-15
+++

This is a **bold** introduction. This has *italic* text in it.

More content.
`,
			config: `
base_url = "https://test.example.com/"
title = "Test Site"
summary_length = 2
`,
			expectedSummary: "This is a bold introduction. This has italic text in it.",
			description:     "Should strip HTML formatting from summary",
		},
		{
			name: "Fallback_NoSentences",
			content: `+++
title = "Test Post"
date = 2024-01-15
+++

This is a very long paragraph without proper sentence endings that goes on and on and on and continues for quite a while without any punctuation marks to indicate sentence boundaries
`,
			config: `
base_url = "https://test.example.com/"
title = "Test Site"
summary_length = 2
`,
			expectedSummary: "", // Will be truncated to 150 chars + "..."
			description:     "Should fallback to character limit when no sentences found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary test directory
			tempDir := t.TempDir()

			// Write config file
			configPath := filepath.Join(tempDir, "config.toml")
			if err := os.WriteFile(configPath, []byte(tt.config), 0644); err != nil {
				t.Fatalf("failed to write config: %v", err)
			}

			// Load config
			site, err := config.LoadSite(configPath)
			if err != nil {
				t.Fatalf("failed to load config: %v", err)
			}

			// Create content directory
			contentDir := filepath.Join(tempDir, "content", "blog")
			if err := os.MkdirAll(contentDir, 0755); err != nil {
				t.Fatalf("failed to create content directory: %v", err)
			}

			// Write test content file
			contentPath := filepath.Join(contentDir, "test.md")
			if err := os.WriteFile(contentPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write content file: %v", err)
			}

			// Parse content
			p := parser.NewParser(site)
			if err := p.Parse(filepath.Join(tempDir, "content")); err != nil {
				t.Fatalf("failed to parse content: %v", err)
			}

			// Get the parsed page
			blogSection, exists := p.ContentMap["blog"]
			if !exists || len(blogSection.Children) == 0 {
				t.Fatalf("expected blog section with children, got none")
			}

			page := blogSection.Children[0]

			// Check summary
			summary := string(page.Summary)

			if tt.expectedSummary != "" {
				if summary != tt.expectedSummary {
					t.Errorf("Summary mismatch:\nExpected: %q\nGot:      %q", tt.expectedSummary, summary)
				}
			} else {
				// For fallback test, just check it's truncated
				if !strings.HasSuffix(summary, "...") {
					t.Errorf("Expected fallback summary to end with '...', got: %q", summary)
				}
				if len(summary) > 153 { // 150 + "..."
					t.Errorf("Expected fallback summary to be ~150 chars, got %d: %q", len(summary), summary)
				}
			}
		})
	}
}

func TestSummary_EmptyContent(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	configContent := `
base_url = "https://test.example.com/"
title = "Test Site"
`
	configPath := filepath.Join(tempDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	site, err := config.LoadSite(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	contentDir := filepath.Join(tempDir, "content", "blog")
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		t.Fatalf("failed to create content directory: %v", err)
	}

	// Content with no body
	content := `+++
title = "Empty Post"
date = 2024-01-15
+++

`
	contentPath := filepath.Join(contentDir, "empty.md")
	if err := os.WriteFile(contentPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write content file: %v", err)
	}

	p := parser.NewParser(site)
	if err := p.Parse(filepath.Join(tempDir, "content")); err != nil {
		t.Fatalf("failed to parse content: %v", err)
	}

	blogSection, exists := p.ContentMap["blog"]
	if !exists || len(blogSection.Children) == 0 {
		t.Fatalf("expected blog section with children")
	}

	page := blogSection.Children[0]
	summary := string(page.Summary)

	if summary != "" {
		t.Errorf("Expected empty summary for empty content, got: %q", summary)
	}
}

func TestSummary_DefaultConfiguration(t *testing.T) {
	t.Parallel()
	tempDir := t.TempDir()

	// Config without summary_length (should default to 2)
	configContent := `
base_url = "https://test.example.com/"
title = "Test Site"
`
	configPath := filepath.Join(tempDir, "config.toml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	site, err := config.LoadSite(configPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	contentDir := filepath.Join(tempDir, "content", "blog")
	if err := os.MkdirAll(contentDir, 0755); err != nil {
		t.Fatalf("failed to create content directory: %v", err)
	}

	content := `+++
title = "Test Post"
date = 2024-01-15
+++

First sentence here. Second sentence here. Third sentence here. Fourth sentence here.
`
	contentPath := filepath.Join(contentDir, "test.md")
	if err := os.WriteFile(contentPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write content file: %v", err)
	}

	p := parser.NewParser(site)
	if err := p.Parse(filepath.Join(tempDir, "content")); err != nil {
		t.Fatalf("failed to parse content: %v", err)
	}

	blogSection, exists := p.ContentMap["blog"]
	if !exists || len(blogSection.Children) == 0 {
		t.Fatalf("expected blog section with children")
	}

	page := blogSection.Children[0]
	summary := string(page.Summary)

	// Should default to 2 sentences
	expected := "First sentence here. Second sentence here."
	if summary != expected {
		t.Errorf("Expected default 2-sentence summary:\nExpected: %q\nGot:      %q", expected, summary)
	}
}

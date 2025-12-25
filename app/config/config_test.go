package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSite(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantSite *Site
		wantErr  bool
	}{
		{
			name: "valid_config",
			content: `base_url = "https://example.com"
feed_url = "https://example.com/feed.xml"  
title = "Test Site"
description = "A test site"
theme = "default"
img = "/images/default.jpg"
language = "en"
output_dir = "public"
generate_feed = true

[extra]
author = "Test Author"`,
			wantSite: &Site{
				BaseURL:      "https://example.com",
				FeedURL:      "https://example.com/feed.xml",
				Title:        "Test Site",
				Description:  "A test site",
				Theme:        "default",
				Img:          "/images/default.jpg",
				Lang:         "en",
				OutputDir:    "public",
				GenerateFeed: true,
				Extra: map[string]any{
					"author": "Test Author",
				},
			},
			wantErr: false,
		},
		{
			name: "minimal_config",
			content: `base_url = "https://minimal.com"
title = "Minimal Site"`,
			wantSite: &Site{
				BaseURL: "https://minimal.com",
				Title:   "Minimal Site",
				Extra:   nil, // Extra can be nil when not initialized
			},
			wantErr: false,
		},
		{
			name:     "invalid_toml",
			content:  `invalid toml content [[[`,
			wantSite: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.toml")
			require.NoError(t, writeFile(configPath, tt.content))

			site, err := LoadSite(configPath)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, site)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantSite.BaseURL, site.BaseURL)
			assert.Equal(t, tt.wantSite.Title, site.Title)
			assert.Equal(t, tt.wantSite.Description, site.Description)
			assert.Equal(t, tt.wantSite.GenerateFeed, site.GenerateFeed)

			if tt.wantSite.Extra != nil {
				assert.Equal(t, tt.wantSite.Extra, site.Extra)
			}
		})
	}
}

func TestLoadFrontMatter(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		wantFrontMatter *FrontMatter
		wantBody        string
		wantErr         bool
	}{
		{
			name: "valid_front_matter",
			content: `+++
title = "Test Post"
description = "A test post"
date = 2024-01-15T10:30:00Z
updated = 2024-01-16T12:00:00Z
tags = ["test", "golang"]
template = "post.html"
draft = false
featured = true
generate_feed = true
language = "en"
img = "/images/test.jpg"

[extra]
author = "Test Author"
+++

# Test Content

This is the body content.`,
			wantFrontMatter: &FrontMatter{
				Title:        "Test Post",
				Description:  "A test post",
				Date:         parseTime(t, "2024-01-15T10:30:00Z"),
				Updated:      parseTime(t, "2024-01-16T12:00:00Z"),
				Tags:         []string{"test", "golang"},
				Template:     "post.html",
				Draft:        false,
				Featured:     true,
				GenerateFeed: true,
				Lang:         "en",
				Img:          "/images/test.jpg",
				Extra: map[string]any{
					"author": "Test Author",
				},
			},
			wantBody: "# Test Content\n\nThis is the body content.",
			wantErr:  false,
		},
		{
			name: "no_front_matter",
			content: `# No Front Matter

Just regular markdown content.`,
			wantFrontMatter: nil,
			wantBody:        "# No Front Matter\n\nJust regular markdown content.",
			wantErr:         false,
		},
		{
			name: "empty_front_matter",
			content: `+++
+++

# Content Only

Body without front matter data.`,
			wantFrontMatter: &FrontMatter{
				Extra: map[string]any{},
			},
			wantBody: "# Content Only\n\nBody without front matter data.",
			wantErr:  false,
		},
		{
			name: "invalid_front_matter",
			content: `+++
invalid toml [[[
+++

# Content`,
			wantFrontMatter: nil,
			wantBody:        "",
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frontMatter, body, err := LoadFrontMatter([]byte(tt.content))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantBody, string(body))

			if tt.wantFrontMatter == nil {
				assert.Nil(t, frontMatter)
			} else {
				require.NotNil(t, frontMatter)
				assert.Equal(t, tt.wantFrontMatter.Title, frontMatter.Title)
				assert.Equal(t, tt.wantFrontMatter.Description, frontMatter.Description)
				assert.Equal(t, tt.wantFrontMatter.Tags, frontMatter.Tags)
				assert.Equal(t, tt.wantFrontMatter.Draft, frontMatter.Draft)
				assert.Equal(t, tt.wantFrontMatter.Featured, frontMatter.Featured)

				if !tt.wantFrontMatter.Date.IsZero() {
					assert.Equal(t, tt.wantFrontMatter.Date, frontMatter.Date)
				}
			}
		})
	}
}

func TestSiteToConfig(t *testing.T) {
	tests := []struct {
		name   string
		site   *Site
		expect map[string]any
	}{
		{
			name: "complete_site",
			site: &Site{
				BaseURL:      "https://example.com",
				FeedURL:      "https://example.com/feed.xml",
				Title:        "Test Site",
				Description:  "Test Description",
				OutputDir:    "public",
				Lang:         "en",
				GenerateFeed: true,
				Extra: map[string]any{
					"author": "Test Author",
				},
			},
			expect: map[string]any{
				"base_url":      "https://example.com",
				"feed_url":      "https://example.com/feed.xml",
				"title":         "Test Site",
				"description":   "Test Description",
				"output_dir":    "public",
				"lang":          "en",
				"generate_feed": true,
				"extra": map[string]any{
					"author": "Test Author",
				},
			},
		},
		{
			name: "defaults_applied",
			site: &Site{
				BaseURL: "https://example.com",
				Title:   "Test Site",
			},
			expect: map[string]any{
				"base_url":      "https://example.com",
				"feed_url":      "https://example.com/atom.xml",
				"title":         "Test Site",
				"description":   "",
				"output_dir":    "public",
				"lang":          "en",
				"generate_feed": false,
				"extra":         map[string]any{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.site.ToConfig()

			for key, expectedValue := range tt.expect {
				assert.Equal(t, expectedValue, config[key], "key: %s", key)
			}
		})
	}
}

func TestFrontMatterToConfig(t *testing.T) {
	tests := []struct {
		name        string
		frontMatter *FrontMatter
		expect      map[string]any
	}{
		{
			name: "complete_front_matter",
			frontMatter: &FrontMatter{
				Title:        "Test Post",
				Description:  "A test post description",
				Template:     "post.html",
				GenerateFeed: true,
				Lang:         "en",
				Featured:     true,
				Draft:        false,
				Tags:         []string{"golang", "testing"},
				Date:         parseTime(t, "2024-01-15T10:30:00Z"),
				Updated:      parseTime(t, "2024-01-16T12:00:00Z"),
				Extra: map[string]any{
					"author":   "Test Author",
					"category": "technical",
				},
			},
			expect: map[string]any{
				"title":         "Test Post",
				"description":   "A test post description",
				"template":      "post.html",
				"generate_feed": true,
				"lang":          "en",
				"featured":      true,
				"draft":         false,
				"tags":          []string{"golang", "testing"},
				"categories":    []string(nil),
				"series":        "",
				"series_order":  0,
				"date":          parseTime(t, "2024-01-15T10:30:00Z"),
				"updated":       parseTime(t, "2024-01-16T12:00:00Z"),
				"extra": map[string]any{
					"author":   "Test Author",
					"category": "technical",
				},
			},
		},
		{
			name: "minimal_front_matter",
			frontMatter: &FrontMatter{
				Title: "Simple Post",
			},
			expect: map[string]any{
				"title":         "Simple Post",
				"description":   "",
				"template":      "",
				"generate_feed": false,
				"lang":          "",
				"featured":      false,
				"draft":         false,
				"tags":          []string(nil),
				"categories":    []string(nil),
				"series":        "",
				"series_order":  0,
				"date":          time.Time{},
				"updated":       time.Time{},
				"extra":         map[string]any{},
			},
		},
		{
			name: "draft_post",
			frontMatter: &FrontMatter{
				Title: "Draft Post",
				Draft: true,
				Tags:  []string{"draft", "wip"},
			},
			expect: map[string]any{
				"title":         "Draft Post",
				"description":   "",
				"template":      "",
				"generate_feed": false,
				"lang":          "",
				"featured":      false,
				"draft":         true,
				"tags":          []string{"draft", "wip"},
				"categories":    []string(nil),
				"series":        "",
				"series_order":  0,
				"date":          time.Time{},
				"updated":       time.Time{},
				"extra":         map[string]any{},
			},
		},
		{
			name: "featured_post_with_extra",
			frontMatter: &FrontMatter{
				Title:    "Featured Article",
				Featured: true,
				Extra: map[string]any{
					"reading_time": 5,
					"series":       "Go Testing",
					"nested": map[string]any{
						"level1": "value1",
					},
				},
			},
			expect: map[string]any{
				"title":         "Featured Article",
				"description":   "",
				"template":      "",
				"generate_feed": false,
				"lang":          "",
				"featured":      true,
				"draft":         false,
				"tags":          []string(nil),
				"categories":    []string(nil),
				"series":        "",
				"series_order":  0,
				"date":          time.Time{},
				"updated":       time.Time{},
				"extra": map[string]any{
					"reading_time": 5,
					"series":       "Go Testing",
					"nested": map[string]any{
						"level1": "value1",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.frontMatter.ToConfig()

			for key, expectedValue := range tt.expect {
				assert.Equal(t, expectedValue, config[key], "key: %s", key)
			}

			// Verify all expected keys are present
			assert.Len(t, config, len(tt.expect), "config should have all expected keys")
		})
	}
}

func TestMergeConfigs(t *testing.T) {
	siteConfig := map[string]any{
		"base_url":      "https://example.com",
		"title":         "Site Title",
		"description":   "Site Description",
		"generate_feed": false,
		"extra": map[string]any{
			"author": "Site Author",
			"social": map[string]any{
				"twitter": "@site",
			},
		},
	}

	sectionConfig := map[string]any{
		"title":         "Section Title",
		"generate_feed": true,
		"extra": map[string]any{
			"author": "Section Author",
			"social": map[string]any{
				"github": "section",
			},
		},
	}

	pageConfig := map[string]any{
		"title":       "Page Title",
		"description": "Page Description",
		"extra": map[string]any{
			"custom": "page value",
		},
	}

	merged := MergeConfigs(siteConfig, sectionConfig, pageConfig)

	// Page config should override
	assert.Equal(t, "Page Title", merged["title"])
	assert.Equal(t, "Page Description", merged["description"])

	// Section config should override site
	assert.Equal(t, true, merged["generate_feed"])

	// Site config should remain for unoverridden values
	assert.Equal(t, "https://example.com", merged["base_url"])

	// Extra should be deeply merged
	extra := merged["extra"].(map[string]any)
	assert.Equal(t, "Section Author", extra["author"]) // Section overrides site
	assert.Equal(t, "page value", extra["custom"])     // Page adds new

	// Check if social exists before accessing
	if socialVal, exists := extra["social"]; exists {
		social := socialVal.(map[string]any)
		assert.Equal(t, "@site", social["twitter"])  // From site
		assert.Equal(t, "section", social["github"]) // From section
	}
}

func TestGetBool(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]any
		key      string
		expected bool
	}{
		{
			name:     "true_value",
			config:   map[string]any{"generate_feed": true},
			key:      "generate_feed",
			expected: true,
		},
		{
			name:     "false_value",
			config:   map[string]any{"draft": false},
			key:      "draft",
			expected: false,
		},
		{
			name:     "missing_key",
			config:   map[string]any{},
			key:      "missing",
			expected: false,
		},
		{
			name:     "non_bool_value",
			config:   map[string]any{"title": "Test"},
			key:      "title",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBool(tt.config, tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMergeExtra(t *testing.T) {
	tests := []struct {
		name     string
		base     map[string]any
		extra    map[string]any
		expected map[string]any
	}{
		{
			name: "simple_merge",
			base: map[string]any{
				"author": "Base Author",
				"year":   2024,
			},
			extra: map[string]any{
				"author":  "Extra Author",
				"version": "1.0",
			},
			expected: map[string]any{
				"author":  "Extra Author",
				"year":    2024,
				"version": "1.0",
			},
		},
		{
			name: "deep_map_merge",
			base: map[string]any{
				"social": map[string]any{
					"twitter": "@base",
					"email":   "base@example.com",
				},
			},
			extra: map[string]any{
				"social": map[string]any{
					"twitter": "@extra",
					"github":  "extra",
				},
			},
			expected: map[string]any{
				"social": map[string]any{
					"twitter": "@extra",
					"email":   "base@example.com",
					"github":  "extra",
				},
			},
		},
		{
			name: "array_merge",
			base: map[string]any{
				"tags": []any{"go", "web"},
			},
			extra: map[string]any{
				"tags": []any{"testing", "cli"},
			},
			expected: map[string]any{
				"tags": []any{"go", "web", "testing", "cli"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of base since MergeExtra modifies it
			baseCopy := copyMap(tt.base)
			result := MergeExtra(baseCopy, tt.extra)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to deep copy a map.
func copyMap(original map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range original {
		switch val := v.(type) {
		case map[string]any:
			result[k] = copyMap(val)
		case []any:
			copySlice := make([]any, len(val))
			copy(copySlice, val)
			result[k] = copySlice
		default:
			result[k] = v
		}
	}
	return result
}

// Helper functions.
func parseTime(t *testing.T, timeStr string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, timeStr)
	require.NoError(t, err)
	return parsed
}

func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func TestLoadSite_MinifyOptions(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantCSS  bool
		wantHTML bool
		wantErr  bool
	}{
		{
			name: "both_minify_enabled",
			content: `base_url = "https://example.com"
title = "Test"
minify_css = true
minify_html = true`,
			wantCSS:  true,
			wantHTML: true,
			wantErr:  false,
		},
		{
			name: "only_css_minify",
			content: `base_url = "https://example.com"
title = "Test"
minify_css = true`,
			wantCSS:  true,
			wantHTML: false,
			wantErr:  false,
		},
		{
			name: "only_html_minify",
			content: `base_url = "https://example.com"
title = "Test"
minify_html = true`,
			wantCSS:  false,
			wantHTML: true,
			wantErr:  false,
		},
		{
			name: "default_no_minify",
			content: `base_url = "https://example.com"
title = "Test"`,
			wantCSS:  false,
			wantHTML: false,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.toml")
			require.NoError(t, writeFile(configPath, tt.content))

			site, err := LoadSite(configPath)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantCSS, site.MinifyCSS, "MinifyCSS should match")
			assert.Equal(t, tt.wantHTML, site.MinifyHTML, "MinifyHTML should match")
		})
	}
}

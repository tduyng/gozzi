package content

import (
	"html/template"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContentNode(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		parent   *Node
		expected *Node
	}{
		{
			name:   "root_node_no_parent",
			path:   "content/index.md",
			parent: nil,
			expected: &Node{
				Path:     "content/index.md",
				Slug:     "content",
				Parent:   nil,
				Children: []*Node{},
			},
		},
		{
			name: "child_node_with_parent",
			path: "content/blog/post.md",
			parent: &Node{
				Slug: "content",
			},
			expected: &Node{
				Path: "content/blog/post.md",
				Slug: "content/post",
				Parent: &Node{
					Slug: "content",
				},
				Children: []*Node{},
			},
		},
		{
			name:   "date_prefixed_post",
			path:   "content/blog/2024-01-15-my-post.md",
			parent: nil,
			expected: &Node{
				Path:     "content/blog/2024-01-15-my-post.md",
				Slug:     "my-post",
				Parent:   nil,
				Children: []*Node{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewContentNode(tt.path, tt.parent)

			assert.Equal(t, tt.expected.Path, result.Path)
			assert.Equal(t, tt.expected.Slug, result.Slug)
			assert.NotNil(t, result.Children)
			assert.Len(t, result.Children, 0)

			if tt.expected.Parent != nil {
				require.NotNil(t, result.Parent)
				assert.Equal(t, tt.expected.Parent.Slug, result.Parent.Slug)
			} else {
				assert.Nil(t, result.Parent)
			}
		})
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		parent   *Node
		expected string
	}{
		{
			name:     "simple_file",
			path:     "hello-world.md",
			parent:   nil,
			expected: "hello-world",
		},
		{
			name:     "file_with_underscores",
			path:     "hello_world_post.md",
			parent:   nil,
			expected: "hello-world-post",
		},
		{
			name:     "file_with_spaces_and_special_chars",
			path:     "Hello World! & More.md",
			parent:   nil,
			expected: "hello-world-more",
		},
		{
			name:     "date_prefixed_file",
			path:     "2024-01-15-my-awesome-post.md",
			parent:   nil,
			expected: "my-awesome-post",
		},
		{
			name:     "date_prefixed_with_underscores",
			path:     "2024_01_15_another_post.md",
			parent:   nil,
			expected: "another-post",
		},
		{
			name:     "mixed_date_format",
			path:     "2024-1-5-short-date.md",
			parent:   nil,
			expected: "short-date",
		},
		{
			name: "child_with_parent_slug",
			path: "awesome-post.md",
			parent: &Node{
				Slug: "blog",
			},
			expected: "blog/awesome-post",
		},
		{
			name: "nested_parent_slug",
			path: "final-post.md",
			parent: &Node{
				Slug: "blog/2024",
			},
			expected: "blog/2024/final-post",
		},
		{
			name:     "index_file_uses_parent_dir",
			path:     "content/blog/index.md",
			parent:   nil,
			expected: "blog",
		},
		{
			name:     "index_file_in_root",
			path:     "index.md",
			parent:   nil,
			expected: "index",
		},
		{
			name:     "multiple_dashes_cleaned",
			path:     "hello---world--test.md",
			parent:   nil,
			expected: "hello-world-test",
		},
		{
			name:     "leading_trailing_dashes_trimmed",
			path:     "-hello-world-.md",
			parent:   nil,
			expected: "hello-world",
		},
		{
			name: "parent_with_empty_slug",
			path: "test.md",
			parent: &Node{
				Slug: "",
			},
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateSlug(tt.path, tt.parent)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestExtractBaseName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple_file",
			path:     "hello.md",
			expected: "hello",
		},
		{
			name:     "file_in_directory",
			path:     "content/hello.md",
			expected: "hello",
		},
		{
			name:     "nested_directory",
			path:     "content/blog/post.md",
			expected: "post",
		},
		{
			name:     "index_file_uses_parent_dir",
			path:     "content/blog/index.md",
			expected: "blog",
		},
		{
			name:     "index_file_in_root",
			path:     "index.md",
			expected: "index",
		},
		{
			name:     "index_file_deep_nested",
			path:     "content/blog/2024/index.md",
			expected: "2024",
		},
		{
			name:     "file_without_extension",
			path:     "content/hello",
			expected: "hello",
		},
		{
			name:     "multiple_extensions",
			path:     "content/file.tar.gz",
			expected: "file.tar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractBaseName(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNode_ToMap(t *testing.T) {
	// Create test nodes
	parent := &Node{
		Type: NodeTypeSection,
		Slug: "parent",
	}

	child := &Node{
		Type: NodeTypePage,
		Slug: "child",
	}

	node := &Node{
		Type:      NodeTypePage,
		Path:      "content/test.md",
		Slug:      "test",
		Permalink: "https://example.com/test/",
		URL:       "/test/",
		Config: map[string]any{
			"title": "Test Page",
		},
		Content:   template.HTML("<p>Test content</p>"),
		Parent:    parent,
		Children:  []*Node{child},
		Lower:     nil,
		Higher:    nil,
		WordCount: 150,
		ReadTime:  1,
		Toc: []map[string]any{
			{"title": "Heading 1", "level": 1},
		},
	}

	result := node.ToMap()

	// Verify all fields are correctly mapped
	assert.Equal(t, NodeTypePage, result["Type"])
	assert.Equal(t, "content/test.md", result["Path"])
	assert.Equal(t, "test", result["Slug"])
	assert.Equal(t, "https://example.com/test/", result["Permalink"])
	assert.Equal(t, "/test/", result["URL"])
	assert.Equal(t, map[string]any{"title": "Test Page"}, result["Config"])
	assert.Equal(t, template.HTML("<p>Test content</p>"), result["Content"])
	assert.Equal(t, parent, result["Parent"])
	assert.Equal(t, []*Node{child}, result["Children"])
	assert.Nil(t, result["Lower"])
	assert.Nil(t, result["Higher"])
	assert.Equal(t, 150, result["WordCount"])
	assert.Equal(t, 1, result["ReadTime"])
	assert.Equal(t, []map[string]any{{"title": "Heading 1", "level": 1}}, result["Toc"])
}

func TestNode_TemplateChain(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Node
		expected []string
	}{
		{
			name: "node_without_parent_or_template",
			setup: func() *Node {
				return &Node{
					Config: map[string]any{},
				}
			},
			expected: []string{"default.html"},
		},
		{
			name: "node_with_template_config",
			setup: func() *Node {
				return &Node{
					Config: map[string]any{
						"template": "custom.html",
					},
				}
			},
			expected: []string{"custom.html", "default.html"},
		},
		{
			name: "node_with_parent",
			setup: func() *Node {
				parent := &Node{
					Config: map[string]any{
						"template": "parent.html",
					},
				}
				return &Node{
					Parent: parent,
					Config: map[string]any{},
				}
			},
			expected: []string{"parent.html", "default.html", "default.html"},
		},
		{
			name: "node_with_parent_and_own_template",
			setup: func() *Node {
				parent := &Node{
					Config: map[string]any{
						"template": "parent.html",
					},
				}
				return &Node{
					Parent: parent,
					Config: map[string]any{
						"template": "child.html",
					},
				}
			},
			expected: []string{"child.html", "parent.html", "default.html", "default.html"},
		},
		{
			name: "nested_parent_hierarchy",
			setup: func() *Node {
				grandparent := &Node{
					Config: map[string]any{
						"template": "grandparent.html",
					},
				}
				parent := &Node{
					Parent: grandparent,
					Config: map[string]any{
						"template": "parent.html",
					},
				}
				return &Node{
					Parent: parent,
					Config: map[string]any{
						"template": "child.html",
					},
				}
			},
			expected: []string{"child.html", "parent.html", "grandparent.html", "default.html", "default.html", "default.html"},
		},
		{
			name: "empty_template_string_ignored",
			setup: func() *Node {
				return &Node{
					Config: map[string]any{
						"template": "",
					},
				}
			},
			expected: []string{"default.html"},
		},
		{
			name: "non_string_template_ignored",
			setup: func() *Node {
				return &Node{
					Config: map[string]any{
						"template": 123,
					},
				}
			},
			expected: []string{"default.html"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := tt.setup()
			result := node.TemplateChain()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNodeType_Constants(t *testing.T) {
	// Test that the NodeType constants have expected values
	assert.Equal(t, NodeType(0), NodeTypeSection)
	assert.Equal(t, NodeType(1), NodeTypePage)
}

func TestStripDatePrefixFromPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single_date_prefix",
			input:    "books/2026-01-01-harry-potter-7/img/cover.webp",
			expected: filepath.Join("books", "harry-potter-7", "img", "cover.webp"),
		},
		{
			name:     "multiple_date_prefixes",
			input:    "2024-05-10-blog/2026-01-01-post/img/photo.jpg",
			expected: filepath.Join("blog", "post", "img", "photo.jpg"),
		},
		{
			name:     "no_date_prefix",
			input:    "blog/post/img/cover.webp",
			expected: filepath.Join("blog", "post", "img", "cover.webp"),
		},
		{
			name:     "date_with_underscores",
			input:    "2024_05_10_blog/2026_01_01_post/file.txt",
			expected: filepath.Join("blog", "post", "file.txt"),
		},
		{
			name:     "mixed_date_formats",
			input:    "2024-5-10-blog/2026_1_1_post/data.json",
			expected: filepath.Join("blog", "post", "data.json"),
		},
		{
			name:     "deep_nesting_with_dates",
			input:    "docs/2025-12-15-tutorial/assets/images/diagram.svg",
			expected: filepath.Join("docs", "tutorial", "assets", "images", "diagram.svg"),
		},
		{
			name:     "only_date_prefix",
			input:    "2026-01-01-file.txt",
			expected: "file.txt",
		},
		{
			name:     "empty_string",
			input:    "",
			expected: "",
		},
		{
			name:     "single_component_no_date",
			input:    "file.txt",
			expected: "file.txt",
		},
		{
			name:     "single_component_with_date",
			input:    "2026-01-01-file.txt",
			expected: "file.txt",
		},
		{
			name:     "windows_style_paths",
			input:    "blog/2026-01-01-post/img/cover.webp",
			expected: filepath.Join("blog", "post", "img", "cover.webp"),
		},
		{
			name:     "cross_platform_path_normalization",
			input:    filepath.Join("blog", "2026-01-01-post", "img", "cover.webp"),
			expected: filepath.Join("blog", "post", "img", "cover.webp"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripDatePrefixFromPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRegexPatterns(t *testing.T) {
	t.Run("datePrefixRe", func(t *testing.T) {
		tests := []struct {
			input    string
			expected bool
		}{
			{"2024-01-15-post.md", true},
			{"2024_01_15_post.md", true},
			{"2024-1-5-post.md", true},
			{"2024_1_5_post.md", true},
			{"24-01-15-post.md", false},
			{"post-2024-01-15.md", false},
			{"no-date.md", false},
		}

		for _, tt := range tests {
			result := datePrefixRe.MatchString(tt.input)
			assert.Equal(t, tt.expected, result, "input: %s", tt.input)
		}
	})

	t.Run("slugCleanerRe", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"hello world", "hello-world"},
			{"hello@world", "hello-world"},
			{"hello123world", "hello123world"},
			{"HELLO-WORLD", "hello-world"},
		}

		for _, tt := range tests {
			result := slugCleanerRe.ReplaceAllString(strings.ToLower(tt.input), "-")
			assert.Equal(t, tt.expected, result, "input: %s", tt.input)
		}
	})

	t.Run("multiDashRe", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{"hello--world", "hello-world"},
			{"hello---world", "hello-world"},
			{"hello----world", "hello-world"},
			{"hello-world", "hello-world"},
		}

		for _, tt := range tests {
			result := multiDashRe.ReplaceAllString(tt.input, "-")
			assert.Equal(t, tt.expected, result, "input: %s", tt.input)
		}
	})
}

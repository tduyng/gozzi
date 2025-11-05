// ABOUTME: Tests for site-specific template functions requiring generator context.
package funcs

import (
	"html/template"
	"os"
	"path/filepath"
	"testing"

	"github.com/tduyng/gozzi/app/content"
	"github.com/yuin/goldmark"
)

func TestAssetPath(t *testing.T) {
	ctx := &SiteContext{
		BaseURL: "https://example.com",
	}
	sf := NewSiteFuncs(ctx)

	tests := []struct {
		name     string
		relPath  string
		context  any
		expected string
	}{
		{
			name:     "absolute http URL passes through",
			relPath:  "http://cdn.example.com/style.css",
			context:  nil,
			expected: "http://cdn.example.com/style.css",
		},
		{
			name:     "absolute https URL passes through",
			relPath:  "https://cdn.example.com/script.js",
			context:  nil,
			expected: "https://cdn.example.com/script.js",
		},
		{
			name:     "root-relative path",
			relPath:  "/css/main.css",
			context:  nil,
			expected: "https://example.com/css/main.css",
		},
		{
			name:     "root-relative path without leading slash",
			relPath:  "css/main.css",
			context:  nil,
			expected: "https://example.com/css/main.css",
		},
		{
			name:    "context-relative path with content node",
			relPath: "image.png",
			context: &content.Node{
				Slug: "blog/my-post",
			},
			expected: "https://example.com/blog/my-post/image.png",
		},
		{
			name:    "context-relative with trailing slash in slug",
			relPath: "thumbnail.jpg",
			context: &content.Node{
				Slug: "blog/another-post/",
			},
			expected: "https://example.com/blog/another-post/thumbnail.jpg",
		},
		{
			name:     "non-content context falls back to root-relative",
			relPath:  "assets/icon.svg",
			context:  "not a content node",
			expected: "https://example.com/assets/icon.svg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sf.AssetPath(tt.relPath, tt.context)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestAssetPathWithTrailingSlashInBaseURL(t *testing.T) {
	ctx := &SiteContext{
		BaseURL: "https://example.com/",
	}
	sf := NewSiteFuncs(ctx)

	result, err := sf.AssetPath("/css/main.css", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "https://example.com/css/main.css"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestGetSection(t *testing.T) {
	blogNode := &content.Node{Slug: "blog"}
	notesNode := &content.Node{Slug: "notes"}
	rootNode := &content.Node{Slug: "."}

	ctx := &SiteContext{
		ContentMap: map[string]*content.Node{
			"blog":  blogNode,
			"notes": notesNode,
			".":     rootNode,
		},
	}
	sf := NewSiteFuncs(ctx)

	tests := []struct {
		name        string
		path        string
		expected    *content.Node
		expectError bool
	}{
		{
			name:     "section path without _index.md",
			path:     "blog",
			expected: blogNode,
		},
		{
			name:     "section path with trailing slash",
			path:     "notes/",
			expected: notesNode,
		},
		{
			name:     "section path with _index.md",
			path:     "blog/_index.md",
			expected: blogNode,
		},
		{
			name:     "root section",
			path:     "/",
			expected: rootNode,
		},
		{
			name:        "non-existent section",
			path:        "doesnotexist",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sf.GetSection(tt.path)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestRenderMarkdown(t *testing.T) {
	md := goldmark.New()
	ctx := &SiteContext{
		Markdown: md,
	}
	sf := NewSiteFuncs(ctx)

	tests := []struct {
		name     string
		input    string
		expected template.HTML
	}{
		{
			name:     "simple paragraph",
			input:    "Hello world",
			expected: template.HTML("<p>Hello world</p>\n"),
		},
		{
			name:     "heading",
			input:    "# My Title",
			expected: template.HTML("<h1>My Title</h1>\n"),
		},
		{
			name:     "bold text",
			input:    "**bold**",
			expected: template.HTML("<p><strong>bold</strong></p>\n"),
		},
		{
			name:     "link",
			input:    "[example](https://example.com)",
			expected: template.HTML("<p><a href=\"https://example.com\">example</a></p>\n"),
		},
		{
			name:     "empty input",
			input:    "",
			expected: template.HTML(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sf.RenderMarkdown(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestRenderMarkdownWithoutProcessor(t *testing.T) {
	ctx := &SiteContext{
		Markdown: nil,
	}
	sf := NewSiteFuncs(ctx)

	_, err := sf.RenderMarkdown("# Test")
	if err == nil {
		t.Fatal("expected error when markdown processor is nil")
	}
}

func TestLoadData(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.html")
	testContent := "<div>Hello <strong>World</strong></div>"

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Run("load existing file", func(t *testing.T) {
		result, err := LoadData(testFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if string(result) != testContent {
			t.Errorf("expected %q, got %q", testContent, result)
		}

		// Verify it's template.HTML (unescaped)
		if _, ok := any(result).(template.HTML); !ok {
			t.Error("result should be template.HTML type")
		}
	})

	t.Run("load non-existent file", func(t *testing.T) {
		_, err := LoadData(filepath.Join(tmpDir, "nonexistent.html"))
		if err == nil {
			t.Fatal("expected error for non-existent file")
		}
	})
}

func TestLoadAttribute(t *testing.T) {
	// Create temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	testContent := `<script>alert("xss")</script>`

	err := os.WriteFile(testFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	t.Run("load and escape existing file", func(t *testing.T) {
		result, err := LoadAttribute(testFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Should be HTML-escaped
		expected := template.HTML(`&lt;script&gt;alert(&#34;xss&#34;)&lt;/script&gt;`)
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}

		// Verify it's template.HTML type
		if _, ok := any(result).(template.HTML); !ok {
			t.Error("result should be template.HTML type")
		}
	})

	t.Run("load non-existent file", func(t *testing.T) {
		_, err := LoadAttribute(filepath.Join(tmpDir, "nonexistent.txt"))
		if err == nil {
			t.Fatal("expected error for non-existent file")
		}
	})
}

func TestSafeHTML(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "HTML tags",
			input: "<div>Hello</div>",
		},
		{
			name:  "script tag",
			input: "<script>alert('test')</script>",
		},
		{
			name:  "plain text",
			input: "just plain text",
		},
		{
			name:  "empty string",
			input: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SafeHTML(tt.input)

			// Should return as template.HTML without escaping
			if string(result) != tt.input {
				t.Errorf("expected %q, got %q", tt.input, result)
			}

			// Verify it's template.HTML type
			if _, ok := any(result).(template.HTML); !ok {
				t.Error("result should be template.HTML type")
			}
		})
	}
}

func TestNewSiteFuncs(t *testing.T) {
	ctx := &SiteContext{
		BaseURL: "https://example.com",
	}

	sf := NewSiteFuncs(ctx)

	if sf == nil {
		t.Fatal("expected non-nil SiteFuncs")
	}

	if sf.ctx != ctx {
		t.Error("expected context to be set correctly")
	}
}

func TestMacroRenderer(t *testing.T) {
	t.Run("create new macro renderer", func(t *testing.T) {
		tmpl := template.New("test")
		mr := NewMacroRenderer(tmpl)

		if mr == nil {
			t.Fatal("expected non-nil MacroRenderer")
		}

		if mr.templates != tmpl {
			t.Error("expected templates to be set correctly")
		}
	})

	t.Run("render pagination with missing template", func(t *testing.T) {
		tmpl := template.New("test")
		mr := NewMacroRenderer(tmpl)

		renderFn := mr.RenderPagination(map[string]any{})
		_, err := renderFn(map[string]any{})

		if err == nil {
			t.Fatal("expected error when pagination template not found")
		}
	})

	t.Run("render pagination with valid template", func(t *testing.T) {
		tmpl := template.New("test")
		paginationTmpl := tmpl.New("macros/pagination.html")
		_, err := paginationTmpl.Parse(`<div>Page {{ .Page.Current }}</div>`)
		if err != nil {
			t.Fatalf("failed to parse template: %v", err)
		}

		mr := NewMacroRenderer(tmpl)
		renderFn := mr.RenderPagination(map[string]any{"title": "Test Site"})

		result, err := renderFn(map[string]any{
			"Page": map[string]any{
				"Current": 1,
			},
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := template.HTML("<div>Page 1</div>")
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})
}

package generator

import (
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

func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name           string
		site           *config.Site
		setupTemplates func(t *testing.T) string
		expectError    bool
		validate       func(t *testing.T, gen *Generator)
	}{
		{
			name: "creates generator with basic configuration",
			site: &config.Site{
				Title:     "Test Site",
				BaseURL:   "https://example.com",
				OutputDir: "output",
			},
			setupTemplates: func(t *testing.T) string {
				tempDir := t.TempDir()
				// Change to temp directory so templates path works
				oldWd, _ := os.Getwd()
				t.Cleanup(func() { os.Chdir(oldWd) })
				os.Chdir(tempDir)

				// Create templates directory with basic template
				templateDir := filepath.Join(tempDir, "templates")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				templateContent := `<html><head><title>{{.Page.Title}}</title></head><body>{{.Page.Content}}</body></html>`
				require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"), []byte(templateContent), 0644))

				return tempDir
			},
			expectError: false,
			validate: func(t *testing.T, gen *Generator) {
				assert.NotNil(t, gen.site)
				assert.NotNil(t, gen.parser)
				assert.NotNil(t, gen.templ)
				assert.Equal(t, "Test Site", gen.site.Title)
			},
		},
		{
			name: "handles missing templates directory gracefully",
			site: &config.Site{
				Title:     "Test Site",
				BaseURL:   "https://example.com",
				OutputDir: "output",
			},
			setupTemplates: func(t *testing.T) string {
				tempDir := t.TempDir()
				oldWd, _ := os.Getwd()
				t.Cleanup(func() { os.Chdir(oldWd) })
				os.Chdir(tempDir)
				// Don't create templates directory
				return tempDir
			},
			expectError: true,
			validate:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupTemplates != nil {
				tt.setupTemplates(t)
			}

			p := parser.NewParser(tt.site)
			gen, err := NewGenerator(tt.site, p)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, gen)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, gen)
				if tt.validate != nil {
					tt.validate(t, gen)
				}
			}
		})
	}
}

func TestReloadTemplates(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Create templates directory
	templateDir := filepath.Join(tempDir, "templates")
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	// Create initial template
	templateContent := `<html><body>{{.Page.Title}}</body></html>`
	templatePath := filepath.Join(templateDir, "post.html")
	require.NoError(t, os.WriteFile(templatePath, []byte(templateContent), 0644))

	site := &config.Site{Title: "Test Site", BaseURL: "https://example.com", OutputDir: "output"}
	p := parser.NewParser(site)
	gen, err := NewGenerator(site, p)
	require.NoError(t, err)

	// Verify initial template works
	tpl := gen.templ.Lookup("post.html")
	require.NotNil(t, tpl)

	// Modify template
	newContent := `<html><body><h1>{{.Page.Title}}</h1></body></html>`
	require.NoError(t, os.WriteFile(templatePath, []byte(newContent), 0644))

	// Reload templates
	err = gen.ReloadTemplates()
	assert.NoError(t, err)

	// Verify template was reloaded
	tpl = gen.templ.Lookup("post.html")
	assert.NotNil(t, tpl)
}

func TestHasTemplate(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Create templates directory
	templateDir := filepath.Join(tempDir, "templates")
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	templateContent := `<html><body>{{.Page.Title}}</body></html>`
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"), []byte(templateContent), 0644))

	site := &config.Site{Title: "Test Site", BaseURL: "https://example.com", OutputDir: "output"}
	p := parser.NewParser(site)
	gen, err := NewGenerator(site, p)
	require.NoError(t, err)

	tests := []struct {
		name     string
		template string
		expected bool
	}{
		{"existing template", "post.html", true},
		{"non-existing template", "nonexistent.html", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.hasTemplate(tt.template)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildTagPermalink(t *testing.T) {
	site := &config.Site{Title: "Test Site", BaseURL: "https://example.com"}
	p := parser.NewParser(site)
	gen := &Generator{site: site, parser: p}

	tests := []struct {
		name     string
		tag      string
		expected string
	}{
		{"simple tag", "go", "/tags/go/"},
		{"tag with spaces", "machine learning", "/tags/machine-learning/"},
		{"tag with special chars", "C++", "/tags/c/"},
		{"mixed case tag", "JavaScript", "/tags/javascript/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.buildTagPermalink(tt.tag)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuildTagURL(t *testing.T) {
	site := &config.Site{Title: "Test Site", BaseURL: "https://example.com"}
	p := parser.NewParser(site)
	gen := &Generator{site: site, parser: p}

	tests := []struct {
		name     string
		tagLink  string
		expected string
	}{
		{"simple tag link", "/tags/go/", "https://example.com/tags/go/"},
		{"nested tag link", "/tags/machine-learning/", "https://example.com/tags/machine-learning/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.buildTagURL(tt.tagLink)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWalkNodes(t *testing.T) {
	// Create a content hierarchy
	root := &content.Node{
		Path: ".",
		Type: content.NodeTypeSection,
		Children: []*content.Node{
			{
				Path: "blog",
				Type: content.NodeTypeSection,
				Children: []*content.Node{
					{Path: "blog/post1.md", Type: content.NodeTypePage},
					{Path: "blog/post2.md", Type: content.NodeTypePage},
				},
			},
			{Path: "about.md", Type: content.NodeTypePage},
		},
	}

	site := &config.Site{Title: "Test Site", BaseURL: "https://example.com"}
	p := parser.NewParser(site)
	gen := &Generator{site: site, parser: p}

	var visitedPaths []string
	gen.walkNodes(root, func(n *content.Node) {
		visitedPaths = append(visitedPaths, n.Path)
	})

	expected := []string{".", "blog", "blog/post1.md", "blog/post2.md", "about.md"}
	assert.Equal(t, expected, visitedPaths)
}

func TestCopyFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file
	srcContent := "test content"
	srcPath := filepath.Join(tempDir, "source.txt")
	require.NoError(t, os.WriteFile(srcPath, []byte(srcContent), 0644))

	// Copy to destination
	dstPath := filepath.Join(tempDir, "subdir", "dest.txt")
	err := copyFile(srcPath, dstPath)
	assert.NoError(t, err)

	// Verify destination exists and has correct content
	content, err := os.ReadFile(dstPath)
	assert.NoError(t, err)
	assert.Equal(t, srcContent, string(content))
}

func TestCopyDir(t *testing.T) {
	tempDir := t.TempDir()

	// Create source directory structure
	srcDir := filepath.Join(tempDir, "src")
	require.NoError(t, os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(srcDir, "subdir", "file2.txt"), []byte("content2"), 0644))

	// Copy to destination
	dstDir := filepath.Join(tempDir, "dst")
	err := copyDir(srcDir, dstDir)
	assert.NoError(t, err)

	// Verify files were copied
	content1, err := os.ReadFile(filepath.Join(dstDir, "file1.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "content1", string(content1))

	content2, err := os.ReadFile(filepath.Join(dstDir, "subdir", "file2.txt"))
	assert.NoError(t, err)
	assert.Equal(t, "content2", string(content2))
}

func TestGenerateRobotsTxt(t *testing.T) {
	tempDir := t.TempDir()
	site := &config.Site{
		Title:     "Test Site",
		BaseURL:   "https://example.com",
		OutputDir: tempDir,
	}
	p := parser.NewParser(site)
	gen := &Generator{site: site, parser: p}

	err := gen.generateRobotsTxt()
	assert.NoError(t, err)

	// Verify robots.txt was created
	content, err := os.ReadFile(filepath.Join(tempDir, "robots.txt"))
	assert.NoError(t, err)

	expected := `User-agent: *
Allow: /
Sitemap: https://example.com/sitemap.xml
`
	assert.Equal(t, expected, string(content))
}

func TestGenerate404Page(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Create templates directory with 404 template
	templateDir := filepath.Join(tempDir, "templates")
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	template404 := `<html><body><h1>404 - {{.Page.Title}}</h1></body></html>`
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "404.html"), []byte(template404), 0644))

	site := &config.Site{
		Title:     "Test Site",
		BaseURL:   "https://example.com",
		OutputDir: tempDir,
	}
	p := parser.NewParser(site)
	gen, err := NewGenerator(site, p)
	require.NoError(t, err)

	err = gen.generate404Page()
	assert.NoError(t, err)

	// Verify 404.html was created
	content, err := os.ReadFile(filepath.Join(tempDir, "404.html"))
	assert.NoError(t, err)
	assert.Contains(t, string(content), "404 - Page Not Found")
}

func TestProcessNode(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Create templates
	templateDir := filepath.Join(tempDir, "templates")
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	pageTemplate := `<html><body><h1>{{.Page.Title}}</h1>{{.Page.Content}}</body></html>`
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"), []byte(pageTemplate), 0644))

	site := &config.Site{
		Title:     "Test Site",
		BaseURL:   "https://example.com",
		OutputDir: tempDir,
	}
	p := parser.NewParser(site)
	gen, err := NewGenerator(site, p)
	require.NoError(t, err)

	tests := []struct {
		name string
		node *content.Node
	}{
		{
			name: "process page node",
			node: &content.Node{
				Type:    content.NodeTypePage,
				Slug:    "test-post",
				Content: template.HTML("<p>Test content</p>"),
				Config: map[string]any{
					"title":    "Test Post",
					"template": "post.html",
					"date":     time.Now(),
				},
				Parent: &content.Node{
					Type: content.NodeTypeSection,
					Slug: "blog",
				},
			},
		},
		{
			name: "process section node",
			node: &content.Node{
				Type:    content.NodeTypeSection,
				Slug:    "blog",
				Content: template.HTML("<p>Blog section</p>"),
				Config: map[string]any{
					"title":    "Blog",
					"template": "post.html",
					"date":     time.Now(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := gen.processNode(tt.node)
			assert.NoError(t, err)

			// Verify output file was created
			expectedPath := filepath.Join(tempDir, tt.node.Slug, "index.html")
			_, err = os.Stat(expectedPath)
			assert.NoError(t, err)
		})
	}
}

func TestGetLastMod(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	updateTime := time.Date(2024, 2, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		node     *content.Node
		expected string
	}{
		{
			name: "uses updated date when available",
			node: &content.Node{
				Config: map[string]any{
					"date":    testTime,
					"updated": updateTime,
				},
			},
			expected: "2024-02-15",
		},
		{
			name: "falls back to date when no updated",
			node: &content.Node{
				Config: map[string]any{
					"date": testTime,
				},
			},
			expected: "2024-01-15",
		},
		{
			name: "uses current time when no dates",
			node: &content.Node{
				Config: map[string]any{},
			},
			expected: time.Now().UTC().Format("2006-01-02"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getLastMod(tt.node)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper function to create a test generator with minimal setup
func createTestGenerator(t *testing.T) (*Generator, string) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(oldWd) })
	os.Chdir(tempDir)

	// Create templates directory with basic template
	templateDir := filepath.Join(tempDir, "templates")
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	templateContent := `<html><body>{{.Page.Title}}</body></html>`
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"), []byte(templateContent), 0644))

	site := &config.Site{
		Title:     "Test Site",
		BaseURL:   "https://example.com",
		OutputDir: filepath.Join(tempDir, "output"),
	}
	p := parser.NewParser(site)
	gen, err := NewGenerator(site, p)
	require.NoError(t, err)

	return gen, tempDir
}

func TestCopyStaticAssets(t *testing.T) {
	gen, tempDir := createTestGenerator(t)

	// Create static directory with files
	staticDir := filepath.Join(tempDir, "static")
	require.NoError(t, os.MkdirAll(filepath.Join(staticDir, "css"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(staticDir, "style.css"), []byte("body { margin: 0; }"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(staticDir, "css", "main.css"), []byte(".header { color: blue; }"), 0644))

	err := gen.copyStaticAssets()
	assert.NoError(t, err)

	// Verify files were copied to output
	content, err := os.ReadFile(filepath.Join(gen.site.OutputDir, "style.css"))
	assert.NoError(t, err)
	assert.Equal(t, "body { margin: 0; }", string(content))

	content, err = os.ReadFile(filepath.Join(gen.site.OutputDir, "css", "main.css"))
	assert.NoError(t, err)
	assert.Equal(t, ".header { color: blue; }", string(content))
}

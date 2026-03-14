package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/parser"
)

func TestNewBuilder(t *testing.T) {
	tests := []struct {
		name           string
		site           *config.Site
		setupTemplates func(t *testing.T) string
		expectError    bool
		validate       func(t *testing.T, b *Builder)
	}{
		{
			name: "creates builder with basic configuration",
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

				// Create templates directory with essential templates
				templateDir := filepath.Join(tempDir, "templates")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				templateContent := `<html><body>{{.Page.Title}}</body></html>`
				require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"), []byte(templateContent), 0644))
				require.NoError(t, os.WriteFile(filepath.Join(templateDir, "default.html"), []byte(templateContent), 0644))
				require.NoError(t, os.WriteFile(filepath.Join(templateDir, "404.html"), []byte(`<html><body><h1>404 Not Found</h1></body></html>`), 0644))
				require.NoError(t, os.WriteFile(filepath.Join(templateDir, "home.html"), []byte(templateContent), 0644))

				return tempDir
			},
			expectError: false,
			validate: func(t *testing.T, b *Builder) {
				assert.NotNil(t, b.site)
				assert.NotNil(t, b.parser)
				assert.NotNil(t, b.templ)
				assert.Equal(t, "Test Site", b.site.Title)
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
			b, err := NewBuilder(tt.site, p)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, b)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, b)
				if tt.validate != nil {
					tt.validate(t, b)
				}
			}
		})
	}
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T, b *Builder, tempDir string)
	}{
		{
			name: "successful generation with basic content",
			validate: func(t *testing.T, b *Builder, tempDir string) {
				// Create essential templates that Generate method expects
				templateDir := filepath.Join(tempDir, "templates")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				templates := map[string]string{
					"home.html":    `<html><head><title>{{.Page.Title}}</title></head><body>{{.Page.Content}}</body></html>`,
					"404.html":     `<html><head><title>404 Not Found</title></head><body><h1>Page Not Found</h1></body></html>`,
					"default.html": `<html><head><title>{{.Page.Title}}</title></head><body>{{.Page.Content}}</body></html>`,
				}

				for name, content := range templates {
					require.NoError(t, os.WriteFile(filepath.Join(templateDir, name), []byte(content), 0644))
				}

				// Reload templates after creating them
				err := b.ReloadTemplates()
				require.NoError(t, err)

				// Create minimal content
				contentFiles := map[string]string{
					"_index.md": `+++
title = "Home"
date = 2024-01-15T10:00:00Z
+++
# Home Page`,
				}

				// Create content files
				contentDir := filepath.Join(tempDir, "content")
				require.NoError(t, os.MkdirAll(contentDir, 0755))
				for filename, content := range contentFiles {
					filePath := filepath.Join(contentDir, filename)
					require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))
				}

				// Parse content
				err = b.parser.Parse(contentDir)
				require.NoError(t, err)

				// Get content root
				contentRoot, exists := b.parser.ContentMap["."]
				require.True(t, exists, "Parser should have content root")

				// Generate the site
				err = b.Generate(contentRoot)
				assert.NoError(t, err)

				// Verify basic files were generated
				outputDir := b.site.OutputDir
				assert.FileExists(t, filepath.Join(outputDir, "index.html"))
				assert.FileExists(t, filepath.Join(outputDir, "404.html"))
				assert.FileExists(t, filepath.Join(outputDir, "robots.txt"))
				assert.FileExists(t, filepath.Join(outputDir, "atom.xml"))
				assert.FileExists(t, filepath.Join(outputDir, "sitemap.xml"))

				// Verify robots.txt content
				robotsContent, err := os.ReadFile(filepath.Join(outputDir, "robots.txt"))
				assert.NoError(t, err)
				assert.Contains(t, string(robotsContent), "User-agent: *")
				assert.Contains(t, string(robotsContent), b.site.BaseURL+"/sitemap.xml")
			},
		},
		{
			name: "preserves existing output directory (incremental build)",
			validate: func(t *testing.T, b *Builder, tempDir string) {
				// Create essential templates
				templateDir := filepath.Join(tempDir, "templates")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				templates := map[string]string{
					"home.html":    `<html><head><title>{{.Page.Title}}</title></head><body>{{.Page.Content}}</body></html>`,
					"404.html":     `<html><head><title>404 Not Found</title></head><body><h1>Page Not Found</h1></body></html>`,
					"default.html": `<html><head><title>{{.Page.Title}}</title></head><body>{{.Page.Content}}</body></html>`,
				}

				for name, content := range templates {
					require.NoError(t, os.WriteFile(filepath.Join(templateDir, name), []byte(content), 0644))
				}

				// Reload templates
				err := b.ReloadTemplates()
				require.NoError(t, err)

				// Create existing output with old files
				outputDir := b.site.OutputDir
				require.NoError(t, os.MkdirAll(outputDir, 0755))
				oldFile := filepath.Join(outputDir, "old-file.html")
				require.NoError(t, os.WriteFile(oldFile, []byte("old content"), 0644))

				// Create minimal content
				contentDir := filepath.Join(tempDir, "content")
				require.NoError(t, os.MkdirAll(contentDir, 0755))
				require.NoError(t, os.WriteFile(filepath.Join(contentDir, "_index.md"), []byte(`+++
title = "Home"
date = 2024-01-15T10:00:00Z
+++
# Home`), 0644))

				// Parse and generate
				err = b.parser.Parse(contentDir)
				require.NoError(t, err)

				contentRoot, exists := b.parser.ContentMap["."]
				require.True(t, exists)

				err = b.Generate(contentRoot)
				assert.NoError(t, err)

				// Note: With incremental builds, old files are NOT automatically cleaned
				// This is intentional to preserve user-generated files
				assert.FileExists(t, oldFile, "Old files should be preserved for incremental builds")
				// Verify new file was generated
				assert.FileExists(t, filepath.Join(outputDir, "index.html"))
			},
		},
		{
			name: "handles static assets",
			validate: func(t *testing.T, b *Builder, tempDir string) {
				// Create essential templates
				templateDir := filepath.Join(tempDir, "templates")
				require.NoError(t, os.MkdirAll(templateDir, 0755))

				templates := map[string]string{
					"home.html":    `<html><head><title>{{.Page.Title}}</title></head><body>{{.Page.Content}}</body></html>`,
					"404.html":     `<html><head><title>404 Not Found</title></head><body><h1>Page Not Found</h1></body></html>`,
					"default.html": `<html><head><title>{{.Page.Title}}</title></head><body>{{.Page.Content}}</body></html>`,
				}

				for name, content := range templates {
					require.NoError(t, os.WriteFile(filepath.Join(templateDir, name), []byte(content), 0644))
				}

				// Reload templates
				err := b.ReloadTemplates()
				require.NoError(t, err)

				// Create static assets
				staticDir := filepath.Join(tempDir, "static")
				require.NoError(t, os.MkdirAll(filepath.Join(staticDir, "css"), 0755))
				require.NoError(t, os.WriteFile(filepath.Join(staticDir, "style.css"), []byte("body { margin: 0; }"), 0644))
				require.NoError(t, os.WriteFile(filepath.Join(staticDir, "css", "main.css"), []byte(".header { color: blue; }"), 0644))

				// Create content
				contentDir := filepath.Join(tempDir, "content")
				require.NoError(t, os.MkdirAll(contentDir, 0755))
				require.NoError(t, os.WriteFile(filepath.Join(contentDir, "_index.md"), []byte(`+++
title = "Home"
date = 2024-01-15T10:00:00Z
+++
# Home`), 0644))

				// Parse and generate
				err = b.parser.Parse(contentDir)
				require.NoError(t, err)

				contentRoot, exists := b.parser.ContentMap["."]
				require.True(t, exists)

				err = b.Generate(contentRoot)
				assert.NoError(t, err)

				// Verify static assets were copied
				outputDir := b.site.OutputDir
				content, err := os.ReadFile(filepath.Join(outputDir, "style.css"))
				assert.NoError(t, err)
				assert.Equal(t, "body { margin: 0; }", string(content))

				content, err = os.ReadFile(filepath.Join(outputDir, "css", "main.css"))
				assert.NoError(t, err)
				assert.Equal(t, ".header { color: blue; }", string(content))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd)
			os.Chdir(tempDir)

			// Create templates directory with basic templates
			templateDir := filepath.Join(tempDir, "templates")
			require.NoError(t, os.MkdirAll(templateDir, 0755))

			// Create all necessary templates
			templates := map[string]string{
				"404.html":  `<html><body><h1>404 - {{.Page.Title}}</h1></body></html>`,
				"post.html": `<html><body><h1>{{.Page.Title}}</h1><div>{{.Page.Content}}</div></body></html>`,
			}
			for name, content := range templates {
				require.NoError(t, os.WriteFile(filepath.Join(templateDir, name), []byte(content), 0644))
			}

			// Create site and builder
			site := &config.Site{
				Title:     "Test Site",
				BaseURL:   "https://example.com",
				OutputDir: filepath.Join(tempDir, "output"),
			}

			p := parser.NewParser(site)
			b, err := NewBuilder(site, p)
			require.NoError(t, err)

			// Run test validation
			if tt.validate != nil {
				tt.validate(t, b, tempDir)
			}
		})
	}
}

func TestGenerateWithExistingOutput(t *testing.T) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Create templates and content
	templateDir := filepath.Join(tempDir, "templates")
	require.NoError(t, os.MkdirAll(templateDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"), []byte(`<html><body>{{.Page.Title}}</body></html>`), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "default.html"), []byte(`<html><body>{{.Page.Title}}</body></html>`), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "404.html"), []byte(`<html><body>404 Not Found</body></html>`), 0644))

	contentDir := filepath.Join(tempDir, "content")
	require.NoError(t, os.MkdirAll(contentDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(contentDir, "_index.md"), []byte(`+++
title = "Home"
date = 2024-01-15T10:00:00Z
+++
# Home`), 0644))

	site := &config.Site{
		Title:     "Test Site",
		BaseURL:   "https://example.com",
		OutputDir: filepath.Join(tempDir, "output"),
	}

	// Create existing output directory with old files
	require.NoError(t, os.MkdirAll(site.OutputDir, 0755))
	oldFile := filepath.Join(site.OutputDir, "old-file.html")
	require.NoError(t, os.WriteFile(oldFile, []byte("old content"), 0644))

	p := parser.NewParser(site)
	err := p.Parse(contentDir)
	require.NoError(t, err)

	b, err := NewBuilder(site, p)
	require.NoError(t, err)

	// Get content root
	contentRoot, exists := p.ContentMap["."]
	require.True(t, exists, "Parser should have content root")

	// Generate should preserve existing output and add new files (incremental build)
	err = b.Generate(contentRoot)
	assert.NoError(t, err)

	// Verify old file was preserved (incremental build keeps user files)
	assert.FileExists(t, oldFile, "Old files should be preserved for incremental builds")

	// Verify new files were generated
	assert.FileExists(t, filepath.Join(site.OutputDir, "index.html"))
}

func TestGenerateWithConcurrentProcessing(t *testing.T) {
	// Create a larger site to test concurrent processing
	contentFiles := make(map[string]string)

	// Create many posts to trigger concurrent processing
	for i := 1; i <= 20; i++ {
		filename := fmt.Sprintf("blog/post%02d.md", i)
		content := fmt.Sprintf(`+++
title = "Post %d"
date = 2024-01-%02dT10:00:00Z
description = "Post number %d"
+++
# Post %d
Content for post %d`, i, (i%28)+1, i, i, i)
		contentFiles[filename] = content
	}

	// Add section index
	contentFiles["blog/_index.md"] = `+++
title = "Blog"
date = 2024-01-01T10:00:00Z
+++
# Blog Section`

	b, _ := createTestBuilderWithContent(t, contentFiles)

	// Get content root
	contentRoot, exists := b.parser.ContentMap["."]
	require.True(t, exists, "Parser should have content root")

	// Generate the site (should use concurrent processing)
	err := b.Generate(contentRoot)
	assert.NoError(t, err)

	// Verify all posts were generated
	for i := 1; i <= 20; i++ {
		expectedPath := filepath.Join(b.site.OutputDir, "blog", fmt.Sprintf("post%02d", i), "index.html")
		assert.FileExists(t, expectedPath)
	}

	// Verify section was generated
	assert.FileExists(t, filepath.Join(b.site.OutputDir, "blog", "index.html"))
}

func TestMinifyCSS(t *testing.T) {
	tests := []struct {
		name       string
		minifyCSS  bool
		cssContent string
		checkSize  bool
	}{
		{
			name:      "minify_enabled",
			minifyCSS: true,
			cssContent: `
body {
    margin: 0;
    padding: 0;
    font-family: sans-serif;
}

/* Comment */
.container {
    max-width: 1200px;
}`,
			checkSize: true,
		},
		{
			name:      "minify_disabled",
			minifyCSS: false,
			cssContent: `
body {
    margin: 0;
}`,
			checkSize: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputDir := filepath.Join(tmpDir, "public")
			staticDir := filepath.Join(tmpDir, "static")
			cssDir := filepath.Join(staticDir, "css")
			templateDir := filepath.Join(tmpDir, "templates")
			contentDir := filepath.Join(tmpDir, "content")

			require.NoError(t, os.MkdirAll(cssDir, 0755))
			require.NoError(t, os.MkdirAll(templateDir, 0755))
			require.NoError(t, os.MkdirAll(contentDir, 0755))

			cssPath := filepath.Join(cssDir, "style.css")
			require.NoError(t, os.WriteFile(cssPath, []byte(tt.cssContent), 0644))

			require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"),
				[]byte(`<html><body>{{.Page.Title}}</body></html>`), 0644))
			require.NoError(t, os.WriteFile(filepath.Join(templateDir, "404.html"),
				[]byte(`<html><body>404</body></html>`), 0644))

			cfg := &config.Site{
				BaseURL:   "https://example.com",
				Title:     "Test",
				OutputDir: outputDir,
				MinifyCSS: tt.minifyCSS,
			}

			require.NoError(t, os.Chdir(tmpDir))
			defer func() { _ = os.Chdir("../../..") }()

			p := parser.NewParser(cfg)
			builder, err := NewBuilder(cfg, p)
			require.NoError(t, err)

			err = builder.copyStaticAssets()
			require.NoError(t, err)

			outputCSS := filepath.Join(outputDir, "css", "style.css")
			assert.FileExists(t, outputCSS)

			originalSize := len(tt.cssContent)
			outputContent, err := os.ReadFile(outputCSS)
			require.NoError(t, err)
			outputSize := len(outputContent)

			if tt.checkSize {
				assert.Less(t, outputSize, originalSize, "Minified CSS should be smaller")
				assert.NotContains(t, string(outputContent), "/* Comment */")
			} else {
				assert.Equal(t, originalSize, outputSize, "CSS size should be unchanged when minify disabled")
			}
		})
	}
}

func TestMinifyHTML(t *testing.T) {
	tests := []struct {
		name       string
		minifyHTML bool
		checkSize  bool
	}{
		{
			name:       "minify_enabled",
			minifyHTML: true,
			checkSize:  true,
		},
		{
			name:       "minify_disabled",
			minifyHTML: false,
			checkSize:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputDir := filepath.Join(tmpDir, "public")
			templateDir := filepath.Join(tmpDir, "templates")
			contentDir := filepath.Join(tmpDir, "content")

			require.NoError(t, os.MkdirAll(templateDir, 0755))
			require.NoError(t, os.MkdirAll(contentDir, 0755))

			templateContent := `<!DOCTYPE html>
<html>
<head>
    <title>{{.Page.Title}}</title>
</head>
<body>
    <h1>{{.Page.Title}}</h1>
    <div class="content">
        {{.Page.Body}}
    </div>
</body>
</html>`
			require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"),
				[]byte(templateContent), 0644))
			require.NoError(t, os.WriteFile(filepath.Join(templateDir, "404.html"),
				[]byte(`<html><body>404</body></html>`), 0644))

			postContent := `+++
title = "Test Post"
date = 2025-01-01T00:00:00Z
template = "post.html"
+++

Test content here`
			require.NoError(t, os.WriteFile(filepath.Join(contentDir, "_index.md"),
				[]byte(postContent), 0644))

			cfg := &config.Site{
				BaseURL:    "https://example.com",
				Title:      "Test",
				OutputDir:  outputDir,
				MinifyHTML: tt.minifyHTML,
			}

			require.NoError(t, os.Chdir(tmpDir))
			defer func() { _ = os.Chdir("../../..") }()

			p := parser.NewParser(cfg)
			err := p.Parse("content")
			require.NoError(t, err)

			builder, err := NewBuilder(cfg, p)
			require.NoError(t, err)

			contentRoot, exists := p.ContentMap["."]
			require.True(t, exists)

			err = builder.Generate(contentRoot)
			require.NoError(t, err)

			outputHTML := filepath.Join(outputDir, "index.html")
			assert.FileExists(t, outputHTML)

			content, err := os.ReadFile(outputHTML)
			require.NoError(t, err)

			if tt.checkSize {
				assert.NotContains(t, string(content), "\n    <")
				assert.Contains(t, string(content), "<html><head>")
			} else {
				assert.Contains(t, string(content), "\n    <")
			}
		})
	}
}

func TestMinifyJS(t *testing.T) {
	tests := []struct {
		name      string
		minifyJS  bool
		jsContent string
		checkSize bool
	}{
		{
			name:     "minify_enabled",
			minifyJS: true,
			jsContent: `
function hello() {
    console.log('Hello World');
}

// This is a comment
const x = 1;
const arr = [1, 2, 3];

/* Multi-line
   comment */
const obj = {
    name: 'test',
    value: 42
};`,
			checkSize: true,
		},
		{
			name:     "minify_disabled",
			minifyJS: false,
			jsContent: `
const x = 1;
// Comment
console.log('test');`,
			checkSize: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputDir := filepath.Join(tmpDir, "public")
			staticDir := filepath.Join(tmpDir, "static")
			jsDir := filepath.Join(staticDir, "js")
			templateDir := filepath.Join(tmpDir, "templates")
			contentDir := filepath.Join(tmpDir, "content")

			require.NoError(t, os.MkdirAll(jsDir, 0755))
			require.NoError(t, os.MkdirAll(templateDir, 0755))
			require.NoError(t, os.MkdirAll(contentDir, 0755))

			jsPath := filepath.Join(jsDir, "script.js")
			require.NoError(t, os.WriteFile(jsPath, []byte(tt.jsContent), 0644))

			require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"),
				[]byte(`<html><body>{{.Page.Title}}</body></html>`), 0644))
			require.NoError(t, os.WriteFile(filepath.Join(templateDir, "404.html"),
				[]byte(`<html><body>404</body></html>`), 0644))

			cfg := &config.Site{
				BaseURL:   "https://example.com",
				Title:     "Test",
				OutputDir: outputDir,
				MinifyJS:  tt.minifyJS,
			}

			require.NoError(t, os.Chdir(tmpDir))
			defer func() { _ = os.Chdir("../../..") }()

			p := parser.NewParser(cfg)
			builder, err := NewBuilder(cfg, p)
			require.NoError(t, err)

			err = builder.copyStaticAssets()
			require.NoError(t, err)

			outputJS := filepath.Join(outputDir, "js", "script.js")
			assert.FileExists(t, outputJS)

			originalSize := len(tt.jsContent)
			outputContent, err := os.ReadFile(outputJS)
			require.NoError(t, err)
			outputSize := len(outputContent)

			if tt.checkSize {
				assert.Less(t, outputSize, originalSize, "Minified JS should be smaller")
				assert.NotContains(t, string(outputContent), "// This is a comment")
				assert.NotContains(t, string(outputContent), "/* Multi-line")
			} else {
				assert.Equal(t, originalSize, outputSize, "JS size should be unchanged when minify disabled")
				assert.Contains(t, string(outputContent), "// Comment")
			}
		})
	}
}

func TestMinifyJSON(t *testing.T) {
	tests := []struct {
		name        string
		minifyJSON  bool
		jsonContent string
		checkSize   bool
	}{
		{
			name:       "minify_enabled",
			minifyJSON: true,
			jsonContent: `{
    "name": "Test User",
    "email": "test@example.com",
    "active": true,
    "tags": ["dev", "golang"]
}`,
			checkSize: true,
		},
		{
			name:       "minify_disabled",
			minifyJSON: false,
			jsonContent: `{
    "key": "value"
}`,
			checkSize: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputDir := filepath.Join(tmpDir, "public")
			staticDir := filepath.Join(tmpDir, "static")
			jsonDir := filepath.Join(staticDir, "data")
			templateDir := filepath.Join(tmpDir, "templates")
			contentDir := filepath.Join(tmpDir, "content")

			require.NoError(t, os.MkdirAll(jsonDir, 0755))
			require.NoError(t, os.MkdirAll(templateDir, 0755))
			require.NoError(t, os.MkdirAll(contentDir, 0755))

			jsonPath := filepath.Join(jsonDir, "data.json")
			require.NoError(t, os.WriteFile(jsonPath, []byte(tt.jsonContent), 0644))

			require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"),
				[]byte(`<html><body>{{.Page.Title}}</body></html>`), 0644))
			require.NoError(t, os.WriteFile(filepath.Join(templateDir, "404.html"),
				[]byte(`<html><body>404</body></html>`), 0644))

			cfg := &config.Site{
				BaseURL:    "https://example.com",
				Title:      "Test",
				OutputDir:  outputDir,
				MinifyJSON: tt.minifyJSON,
			}

			require.NoError(t, os.Chdir(tmpDir))
			defer func() { _ = os.Chdir("../../..") }()

			p := parser.NewParser(cfg)
			builder, err := NewBuilder(cfg, p)
			require.NoError(t, err)

			err = builder.copyStaticAssets()
			require.NoError(t, err)

			outputJSON := filepath.Join(outputDir, "data", "data.json")
			assert.FileExists(t, outputJSON)

			originalSize := len(tt.jsonContent)
			outputContent, err := os.ReadFile(outputJSON)
			require.NoError(t, err)
			outputSize := len(outputContent)

			if tt.checkSize {
				assert.Less(t, outputSize, originalSize, "Minified JSON should be smaller")
				assert.NotContains(t, string(outputContent), "\n    ")
			} else {
				assert.Equal(t, originalSize, outputSize, "JSON size should be unchanged when minify disabled")
			}
		})
	}
}

func TestMinifySVG(t *testing.T) {
	tests := []struct {
		name       string
		minifySVG  bool
		svgContent string
		checkSize  bool
	}{
		{
			name:      "minify_enabled",
			minifySVG: true,
			svgContent: `<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100">
    <!-- Icon -->
    <circle cx="50" cy="50" r="40" fill="blue" />
    <rect x="20" y="20" width="60" height="60" fill="red" />
</svg>`,
			checkSize: true,
		},
		{
			name:      "minify_disabled",
			minifySVG: false,
			svgContent: `<svg>
    <circle cx="10" cy="10" r="5" />
</svg>`,
			checkSize: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputDir := filepath.Join(tmpDir, "public")
			staticDir := filepath.Join(tmpDir, "static")
			svgDir := filepath.Join(staticDir, "images")
			templateDir := filepath.Join(tmpDir, "templates")
			contentDir := filepath.Join(tmpDir, "content")

			require.NoError(t, os.MkdirAll(svgDir, 0755))
			require.NoError(t, os.MkdirAll(templateDir, 0755))
			require.NoError(t, os.MkdirAll(contentDir, 0755))

			svgPath := filepath.Join(svgDir, "icon.svg")
			require.NoError(t, os.WriteFile(svgPath, []byte(tt.svgContent), 0644))

			require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"),
				[]byte(`<html><body>{{.Page.Title}}</body></html>`), 0644))
			require.NoError(t, os.WriteFile(filepath.Join(templateDir, "404.html"),
				[]byte(`<html><body>404</body></html>`), 0644))

			cfg := &config.Site{
				BaseURL:   "https://example.com",
				Title:     "Test",
				OutputDir: outputDir,
				MinifySVG: tt.minifySVG,
			}

			require.NoError(t, os.Chdir(tmpDir))
			defer func() { _ = os.Chdir("../../..") }()

			p := parser.NewParser(cfg)
			builder, err := NewBuilder(cfg, p)
			require.NoError(t, err)

			err = builder.copyStaticAssets()
			require.NoError(t, err)

			outputSVG := filepath.Join(outputDir, "images", "icon.svg")
			assert.FileExists(t, outputSVG)

			originalSize := len(tt.svgContent)
			outputContent, err := os.ReadFile(outputSVG)
			require.NoError(t, err)
			outputSize := len(outputContent)

			if tt.checkSize {
				assert.LessOrEqual(t, outputSize, originalSize, "Minified SVG should be smaller or equal")
				assert.NotContains(t, string(outputContent), "<!-- Icon -->")
			} else {
				assert.Equal(t, originalSize, outputSize, "SVG size should be unchanged when minify disabled")
			}
		})
	}
}

func TestMinifyXML(t *testing.T) {
	tests := []struct {
		name       string
		minifyXML  bool
		xmlContent string
		checkSize  bool
	}{
		{
			name:      "minify_enabled",
			minifyXML: true,
			xmlContent: `<?xml version="1.0" encoding="UTF-8"?>
<config>
    <setting name="timeout" value="30" />
    <setting name="retries" value="3" />
    <database>
        <host>localhost</host>
        <port>5432</port>
    </database>
</config>`,
			checkSize: true,
		},
		{
			name:      "minify_disabled",
			minifyXML: false,
			xmlContent: `<root>
    <item>Test</item>
</root>`,
			checkSize: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputDir := filepath.Join(tmpDir, "public")
			staticDir := filepath.Join(tmpDir, "static")
			xmlDir := filepath.Join(staticDir, "data")
			templateDir := filepath.Join(tmpDir, "templates")
			contentDir := filepath.Join(tmpDir, "content")

			require.NoError(t, os.MkdirAll(xmlDir, 0755))
			require.NoError(t, os.MkdirAll(templateDir, 0755))
			require.NoError(t, os.MkdirAll(contentDir, 0755))

			xmlPath := filepath.Join(xmlDir, "config.xml")
			require.NoError(t, os.WriteFile(xmlPath, []byte(tt.xmlContent), 0644))

			require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"),
				[]byte(`<html><body>{{.Page.Title}}</body></html>`), 0644))
			require.NoError(t, os.WriteFile(filepath.Join(templateDir, "404.html"),
				[]byte(`<html><body>404</body></html>`), 0644))

			cfg := &config.Site{
				BaseURL:   "https://example.com",
				Title:     "Test",
				OutputDir: outputDir,
				MinifyXML: tt.minifyXML,
			}

			require.NoError(t, os.Chdir(tmpDir))
			defer func() { _ = os.Chdir("../../..") }()

			p := parser.NewParser(cfg)
			builder, err := NewBuilder(cfg, p)
			require.NoError(t, err)

			err = builder.copyStaticAssets()
			require.NoError(t, err)

			outputXML := filepath.Join(outputDir, "data", "config.xml")
			assert.FileExists(t, outputXML)

			originalSize := len(tt.xmlContent)
			outputContent, err := os.ReadFile(outputXML)
			require.NoError(t, err)
			outputSize := len(outputContent)

			if tt.checkSize {
				assert.LessOrEqual(t, outputSize, originalSize, "Minified XML should be smaller or equal")
				assert.NotContains(t, string(outputContent), "\n    ")
			} else {
				assert.Equal(t, originalSize, outputSize, "XML size should be unchanged when minify disabled")
			}
		})
	}
}

func TestSectionWithAssets(t *testing.T) {
	contentFiles := map[string]string{
		"docs/_index.md": `+++
title = "Documentation"
template = "default.html"
+++
# Documentation`,
	}

	b, tempDir := createTestBuilderWithContent(t, contentFiles)

	imgDir := filepath.Join(tempDir, "content", "docs", "img")
	require.NoError(t, os.MkdirAll(imgDir, 0755))
	require.NoError(t, os.WriteFile(filepath.Join(imgDir, "diagram.svg"), []byte("<svg></svg>"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(imgDir, "photo.webp"), []byte("webp data"), 0644))

	contentDir := filepath.Join(tempDir, "content")
	err := b.parser.Parse(contentDir)
	require.NoError(t, err)

	contentRoot, exists := b.parser.ContentMap["docs"]
	require.True(t, exists)

	err = b.Generate(contentRoot)
	assert.NoError(t, err)

	outputDir := b.site.OutputDir
	assert.FileExists(t, filepath.Join(outputDir, "docs", "index.html"))
	assert.FileExists(t, filepath.Join(outputDir, "docs", "img", "diagram.svg"))
	assert.FileExists(t, filepath.Join(outputDir, "docs", "img", "photo.webp"))
}

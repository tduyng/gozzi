package generator

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/config"
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

	// Reload templates to ensure all templates are loaded
	err = gen.ReloadTemplates()
	require.NoError(t, err)
}

// Helper function to create a fully initialized test generator with parsed content
func createTestGeneratorWithContent(t *testing.T, contentFiles map[string]string) (*Generator, string) {
	tempDir := t.TempDir()
	oldWd, _ := os.Getwd()
	t.Cleanup(func() { os.Chdir(oldWd) })
	os.Chdir(tempDir)

	// Create templates directory with basic template
	templateDir := filepath.Join(tempDir, "templates")
	require.NoError(t, os.MkdirAll(templateDir, 0755))

	templateContent := `<html><body>{{.Page.Title}}</body></html>`
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "post.html"), []byte(templateContent), 0644))

	// Create essential fallback templates
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "default.html"), []byte(templateContent), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "404.html"), []byte(`<html><body>404 Not Found</body></html>`), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(templateDir, "home.html"), []byte(templateContent), 0644))

	// Create content directory and files
	contentDir := filepath.Join(tempDir, "content")
	require.NoError(t, os.MkdirAll(contentDir, 0755))

	for filename, content := range contentFiles {
		filePath := filepath.Join(contentDir, filename)
		fileDir := filepath.Dir(filePath)
		require.NoError(t, os.MkdirAll(fileDir, 0755))
		require.NoError(t, os.WriteFile(filePath, []byte(content), 0644))
	}

	site := &config.Site{
		Title:     "Test Site",
		BaseURL:   "https://example.com",
		OutputDir: filepath.Join(tempDir, "output"),
	}

	// Create and initialize parser with content - this is the key!
	p := parser.NewParser(site)
	err := p.Parse(contentDir)
	require.NoError(t, err)

	gen, err := NewGenerator(site, p)
	require.NoError(t, err)

	// Reload templates to ensure all templates are loaded
	err = gen.ReloadTemplates()
	require.NoError(t, err)

	return gen, tempDir
}

// createTestGenerator creates a minimal test generator without content files
func createTestGenerator(t *testing.T) (*Generator, string) {
	return createTestGeneratorWithContent(t, map[string]string{})
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

func TestGenerateAtomFeed(t *testing.T) {
	// Prepare content files
	contentFiles := map[string]string{
		"blog/post1.md": `+++
title = "Test Post 1"
date = 2024-01-15T10:00:00Z
updated = 2024-01-16T11:00:00Z
description = "First test post"
tags = ["go", "testing"]
generate_feed = true
+++
# Test Post 1
This is the first test post content.`,

		"blog/post2.md": `+++
title = "Test Post 2"
date = 2024-01-10T10:00:00Z
description = "Second test post"
tags = ["web", "development"]
generate_feed = true
+++
# Test Post 2
This is the second test post content.`,

		"blog/draft.md": `+++
title = "Draft Post"
date = 2024-01-20T10:00:00Z
description = "Draft post"
draft = true
generate_feed = true
+++
# Draft Post
This should not appear in feed.`,
	}

	gen, _ := createTestGeneratorWithContent(t, contentFiles)

	// Set up site with author info
	gen.site.Extra = map[string]any{
		"author": map[string]any{
			"name":  "Test Author",
			"email": "test@example.com",
		},
	}

	// Create output directory
	require.NoError(t, os.MkdirAll(gen.site.OutputDir, 0755))

	// Generate Atom feed
	err := gen.generateAtomFeed()
	assert.NoError(t, err)

	// Verify atom.xml was created
	feedPath := filepath.Join(gen.site.OutputDir, "atom.xml")
	_, err = os.Stat(feedPath)
	assert.NoError(t, err)

	// Read and validate feed content
	feedContent, err := os.ReadFile(feedPath)
	assert.NoError(t, err)
	feedStr := string(feedContent)

	// Check XML structure
	assert.Contains(t, feedStr, `<?xml version="1.0" encoding="UTF-8"?>`)
	assert.Contains(t, feedStr, `<?xml-stylesheet type="text/xsl" href="/atom.xsl"?>`)
	assert.Contains(t, feedStr, `<feed xmlns="http://www.w3.org/2005/Atom">`)
	assert.Contains(t, feedStr, `<title>Test Site</title>`)
	assert.Contains(t, feedStr, `<id>https://example.com</id>`)
	assert.Contains(t, feedStr, `<link rel="self" href="https://example.com/atom.xml"></link>`)
	assert.Contains(t, feedStr, `<link href="https://example.com"></link>`)

	// Check author info
	assert.Contains(t, feedStr, `<author>`)
	assert.Contains(t, feedStr, `<name>Test Author</name>`)
	assert.Contains(t, feedStr, `<email>test@example.com</email>`)

	// Check entries (should have 2 non-draft posts)
	assert.Contains(t, feedStr, `<title>Test Post 1</title>`)
	assert.Contains(t, feedStr, `<title>Test Post 2</title>`)
	assert.NotContains(t, feedStr, `<title>Draft Post</title>`) // Draft should be excluded

	// Check entry details
	assert.Contains(t, feedStr, `<summary>First test post</summary>`)
	assert.Contains(t, feedStr, `<content type="html">`)
	assert.Contains(t, feedStr, `<category>go</category>`)
	assert.Contains(t, feedStr, `<category>testing</category>`)

	// Verify entries are sorted by date (most recent first)
	post1Index := strings.Index(feedStr, "Test Post 1")
	post2Index := strings.Index(feedStr, "Test Post 2")
	assert.True(t, post1Index < post2Index, "Posts should be sorted by date, most recent first")
}

func TestGenerateSitemap(t *testing.T) {
	// Prepare content files
	contentFiles := map[string]string{
		"_index.md": `+++
title = "Home"
date = 2024-01-15T10:00:00Z
+++
# Home Page`,

		"about.md": `+++
title = "About"
date = 2024-01-10T10:00:00Z
updated = 2024-01-12T10:00:00Z
+++
# About Page`,

		"blog/post1.md": `+++
title = "Blog Post"
date = 2024-01-05T10:00:00Z
tags = ["go", "web"]
+++
# Blog Post`,

		"blog/draft.md": `+++
title = "Draft Post"
date = 2024-01-20T10:00:00Z
draft = true
+++
# Draft Post`,
	}

	gen, tempDir := createTestGeneratorWithContent(t, contentFiles)

	// Create tags template to enable tag pages
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "templates", "tags.html"), []byte(`<html><body>Tags</body></html>`), 0644))

	// Create output directory
	require.NoError(t, os.MkdirAll(gen.site.OutputDir, 0755))

	// Reload templates to pick up tags.html
	err := gen.ReloadTemplates()
	require.NoError(t, err)

	// Generate sitemap
	err = gen.generateSitemap()
	assert.NoError(t, err)

	// Verify sitemap.xml was created
	sitemapPath := filepath.Join(gen.site.OutputDir, "sitemap.xml")
	_, err = os.Stat(sitemapPath)
	assert.NoError(t, err)

	// Read and validate sitemap content
	sitemapContent, err := os.ReadFile(sitemapPath)
	assert.NoError(t, err)
	sitemapStr := string(sitemapContent)

	// Check XML structure
	assert.Contains(t, sitemapStr, `<?xml version="1.0" encoding="UTF-8"?>`)
	assert.Contains(t, sitemapStr, `<?xml-stylesheet type="text/xsl" href="/sitemap.xsl"?>`)
	assert.Contains(t, sitemapStr, `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	// Check that non-draft content is included
	assert.Contains(t, sitemapStr, `<loc>https://example.com/index</loc>`) // Home page (uses index instead of /)
	assert.Contains(t, sitemapStr, `<loc>https://example.com/about</loc>`) // About page (no trailing slash)
	// Note: Blog section might not be included if there's no blog _index.md

	// Check that draft content is excluded
	assert.NotContains(t, sitemapStr, "draft")

	// Check priorities and change frequencies (using actual values from output)
	assert.Contains(t, sitemapStr, `<priority>0.8</priority>`)         // Home page actual priority
	assert.Contains(t, sitemapStr, `<changefreq>weekly</changefreq>`)  // Home page actual frequency
	assert.Contains(t, sitemapStr, `<changefreq>monthly</changefreq>`) // Regular pages

	// Check lastmod dates
	assert.Contains(t, sitemapStr, `<lastmod>2024-01-15</lastmod>`) // Home page date
	assert.Contains(t, sitemapStr, `<lastmod>2024-01-12</lastmod>`) // About page updated date

	// Check tag pages are included
	assert.Contains(t, sitemapStr, `<loc>https://example.com/tags/go/</loc>`)
	assert.Contains(t, sitemapStr, `<loc>https://example.com/tags/web/</loc>`)
}

func TestWriteXMLFile(t *testing.T) {
	gen, _ := createTestGenerator(t)

	// Create output directory
	require.NoError(t, os.MkdirAll(gen.site.OutputDir, 0755))

	// Test data structure
	testData := struct {
		XMLName xml.Name `xml:"test"`
		Title   string   `xml:"title"`
		Items   []string `xml:"item"`
	}{
		Title: "Test XML",
		Items: []string{"item1", "item2"},
	}

	// Write XML file
	err := gen.writeXMLFile("test.xml", `<?xml-stylesheet type="text/xsl" href="/test.xsl"?>`, testData)
	assert.NoError(t, err)

	// Verify file was created
	xmlPath := filepath.Join(gen.site.OutputDir, "test.xml")
	_, err = os.Stat(xmlPath)
	assert.NoError(t, err)

	// Read and validate content
	content, err := os.ReadFile(xmlPath)
	assert.NoError(t, err)
	contentStr := string(content)

	// Check XML structure
	assert.Contains(t, contentStr, `<?xml version="1.0" encoding="UTF-8"?>`)
	assert.Contains(t, contentStr, `<?xml-stylesheet type="text/xsl" href="/test.xsl"?>`)
	assert.Contains(t, contentStr, `<test>`)
	assert.Contains(t, contentStr, `<title>Test XML</title>`)
	assert.Contains(t, contentStr, `<item>item1</item>`)
	assert.Contains(t, contentStr, `<item>item2</item>`)
	assert.Contains(t, contentStr, `</test>`)

	// Check indentation
	lines := strings.Split(contentStr, "\n")
	assert.True(t, len(lines) > 3, "XML should be properly indented")
}

func TestGenerateAtomFeedWithoutAuthor(t *testing.T) {
	// Prepare content files
	contentFiles := map[string]string{
		"post.md": `+++
title = "Simple Post"
date = 2024-01-15T10:00:00Z
description = "A simple post"
generate_feed = true
+++
# Simple Post`,
	}

	gen, _ := createTestGeneratorWithContent(t, contentFiles)

	// Create output directory (no author info in site)
	require.NoError(t, os.MkdirAll(gen.site.OutputDir, 0755))

	// Generate Atom feed
	err := gen.generateAtomFeed()
	assert.NoError(t, err)

	// Verify feed was created without author
	feedPath := filepath.Join(gen.site.OutputDir, "atom.xml")
	content, err := os.ReadFile(feedPath)
	assert.NoError(t, err)
	feedStr := string(content)

	assert.Contains(t, feedStr, `<title>Test Site</title>`)
	assert.NotContains(t, feedStr, `<author>`) // No author should be present
}

func TestGenerateAtomFeedWithManyEntries(t *testing.T) {
	// Prepare many content files to test the 100-entry limit
	contentFiles := make(map[string]string)

	for i := 1; i <= 150; i++ {
		filename := fmt.Sprintf("post%03d.md", i)
		content := fmt.Sprintf(`+++
title = "Post %d"
date = 2024-01-%02dT10:00:00Z
description = "Post number %d"
generate_feed = true
+++
# Post %d`, i, (i%28)+1, i, i) // Spread across dates in January

		contentFiles[filename] = content
	}

	gen, _ := createTestGeneratorWithContent(t, contentFiles)

	// Create output directory
	require.NoError(t, os.MkdirAll(gen.site.OutputDir, 0755))

	// Generate Atom feed
	err := gen.generateAtomFeed()
	assert.NoError(t, err)

	// Read feed content
	feedPath := filepath.Join(gen.site.OutputDir, "atom.xml")
	content, err := os.ReadFile(feedPath)
	assert.NoError(t, err)
	feedStr := string(content)

	// Count entries (should be limited to 100)
	entryCount := strings.Count(feedStr, "<entry>")
	assert.Equal(t, 100, entryCount, "Feed should be limited to 100 entries")
}

func TestGenerate(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T, gen *Generator, tempDir string)
	}{
		{
			name: "successful generation with basic content",
			validate: func(t *testing.T, gen *Generator, tempDir string) {
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
				err := gen.ReloadTemplates()
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
				err = gen.parser.Parse(contentDir)
				require.NoError(t, err)

				// Get content root
				contentRoot, exists := gen.parser.ContentMap["."]
				require.True(t, exists, "Parser should have content root")

				// Generate the site
				err = gen.Generate(contentRoot)
				assert.NoError(t, err)

				// Verify basic files were generated
				outputDir := gen.site.OutputDir
				assert.FileExists(t, filepath.Join(outputDir, "index.html"))
				assert.FileExists(t, filepath.Join(outputDir, "404.html"))
				assert.FileExists(t, filepath.Join(outputDir, "robots.txt"))
				assert.FileExists(t, filepath.Join(outputDir, "atom.xml"))
				assert.FileExists(t, filepath.Join(outputDir, "sitemap.xml"))

				// Verify robots.txt content
				robotsContent, err := os.ReadFile(filepath.Join(outputDir, "robots.txt"))
				assert.NoError(t, err)
				assert.Contains(t, string(robotsContent), "User-agent: *")
				assert.Contains(t, string(robotsContent), gen.site.BaseURL+"/sitemap.xml")
			},
		},
		{
			name: "cleans existing output directory",
			validate: func(t *testing.T, gen *Generator, tempDir string) {
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
				err := gen.ReloadTemplates()
				require.NoError(t, err)

				// Create existing output with old files
				outputDir := gen.site.OutputDir
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
				err = gen.parser.Parse(contentDir)
				require.NoError(t, err)

				contentRoot, exists := gen.parser.ContentMap["."]
				require.True(t, exists)

				err = gen.Generate(contentRoot)
				assert.NoError(t, err)

				// Verify old file was cleaned
				assert.NoFileExists(t, oldFile)
				// Verify new file was generated
				assert.FileExists(t, filepath.Join(outputDir, "index.html"))
			},
		},
		{
			name: "handles static assets",
			validate: func(t *testing.T, gen *Generator, tempDir string) {
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
				err := gen.ReloadTemplates()
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
				err = gen.parser.Parse(contentDir)
				require.NoError(t, err)

				contentRoot, exists := gen.parser.ContentMap["."]
				require.True(t, exists)

				err = gen.Generate(contentRoot)
				assert.NoError(t, err)

				// Verify static assets were copied
				outputDir := gen.site.OutputDir
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

			// Create site and generator
			site := &config.Site{
				Title:     "Test Site",
				BaseURL:   "https://example.com",
				OutputDir: filepath.Join(tempDir, "output"),
			}

			p := parser.NewParser(site)
			gen, err := NewGenerator(site, p)
			require.NoError(t, err)

			// Run test validation
			if tt.validate != nil {
				tt.validate(t, gen, tempDir)
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

	gen, err := NewGenerator(site, p)
	require.NoError(t, err)

	// Get content root
	contentRoot, exists := p.ContentMap["."]
	require.True(t, exists, "Parser should have content root")

	// Generate should clean old output and create new
	err = gen.Generate(contentRoot)
	assert.NoError(t, err)

	// Verify old file was cleaned
	assert.NoFileExists(t, oldFile)

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

	gen, _ := createTestGeneratorWithContent(t, contentFiles)

	// Get content root
	contentRoot, exists := gen.parser.ContentMap["."]
	require.True(t, exists, "Parser should have content root")

	// Generate the site (should use concurrent processing)
	err := gen.Generate(contentRoot)
	assert.NoError(t, err)

	// Verify all posts were generated
	for i := 1; i <= 20; i++ {
		expectedPath := filepath.Join(gen.site.OutputDir, "blog", fmt.Sprintf("post%02d", i), "index.html")
		assert.FileExists(t, expectedPath)
	}

	// Verify section was generated
	assert.FileExists(t, filepath.Join(gen.site.OutputDir, "blog", "index.html"))
}

package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/tduyng/gozzi/app"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
	tplengine "github.com/tduyng/gozzi/app/template"
	"github.com/tduyng/gozzi/app/template/funcs"
)

// Generator handles site generation including templates, feeds, and static files.
type Generator struct {
	site   *config.Site
	templ  *template.Template
	parser *parser.ContentParser
	engine *tplengine.Engine
	mu     sync.Mutex
}

// NewGenerator creates a new Generator with loaded templates.
func NewGenerator(site *config.Site, parser *parser.ContentParser) (*Generator, error) {
	// Initialize template engine
	engine := tplengine.NewEngine(&tplengine.EngineConfig{
		BaseURL:         site.BaseURL,
		ContentMap:      parser.ContentMap,
		Markdown:        parser.GetMarkdownProcessor(),
		StrictTemplates: site.StrictTemplates,
	})

	gen := &Generator{
		site:   site,
		parser: parser,
		engine: engine,
	}

	tmpl, err := gen.loadTemplates()
	if err != nil {
		return nil, app.WrapWithContext(err, app.ErrTemplate, app.ErrorContext{
			Operation: "template_initialization",
			Component: "generator",
		})
	}

	gen.templ = tmpl
	return gen, nil
}

// ReloadTemplates reloads all templates from disk for development mode.
func (g *Generator) ReloadTemplates() error {
	tmpl, err := g.loadTemplates()
	if err != nil {
		return app.WrapWithContext(err, app.ErrTemplate, app.ErrorContext{
			Operation: "template_reload",
			Component: "generator",
		})
	}

	g.mu.Lock()
	g.templ = tmpl
	g.mu.Unlock()
	return nil
}

// loadTemplates loads all templates with the engine's function map.
func (g *Generator) loadTemplates() (*template.Template, error) {
	// Get base function map from engine
	funcMap := g.engine.CreateFuncMap()

	// Add a placeholder for pagination function (will be replaced after templates are loaded)
	funcMap["pagination"] = func(data map[string]any) (template.HTML, error) {
		return "", fmt.Errorf("pagination not yet initialized")
	}

	tmpl := template.New("").Funcs(funcMap)

	err := filepath.WalkDir("templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
				Operation: "template_walk",
				Component: "generator",
				Path:      path,
			})
		}

		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel("templates", path)
		if err != nil {
			return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
				Operation: "get_relative_path",
				Component: "generator",
				Path:      path,
			})
		}

		templateName := filepath.ToSlash(relPath)

		content, err := os.ReadFile(path)
		if err != nil {
			return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
				Operation: "read_template",
				Component: "generator",
				Path:      path,
			})
		}

		_, err = tmpl.New(templateName).Parse(string(content))
		if err != nil {
			return app.WrapWithContext(err, app.ErrTemplate, app.ErrorContext{
				Operation: "parse_template",
				Component: "generator",
				Path:      templateName,
			})
		}

		return nil
	})
	if err != nil {
		return nil, app.WrapWithContext(err, app.ErrTemplate, app.ErrorContext{
			Operation: "load_templates",
			Component: "generator",
		})
	}

	// Add pagination macro after templates are loaded
	macroRenderer := funcs.NewMacroRenderer(tmpl)
	tmpl = tmpl.Funcs(template.FuncMap{
		"pagination": macroRenderer.RenderPagination(g.site.ToConfig()),
	})

	return tmpl, nil
}

// Generate processes the content tree and generates the complete static site.
func (g *Generator) Generate(contentRoot *content.Node) error {
	if err := os.RemoveAll(g.site.OutputDir); err != nil {
		return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
			Operation: "clean_output_directory",
			Component: "generator",
			Path:      g.site.OutputDir,
		})
	}

	sem := make(chan struct{}, runtime.NumCPU()*2)
	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	g.walkNodes(contentRoot, func(n *content.Node) {
		wg.Add(1)
		sem <- struct{}{}

		go func(node *content.Node) {
			defer func() {
				<-sem
				wg.Done()
			}()

			if err := g.processNode(node); err != nil {
				errChan <- err
			}
		}(n)
	})

	wg.Wait()
	close(errChan)

	for err := range errChan {
		fmt.Printf("Processing error: %v\n", err)
	}

	if err := g.generate404Page(); err != nil {
		return app.WrapWithContext(err, app.ErrContent, app.ErrorContext{
			Operation: "generate_404_page",
			Component: "generator",
		})
	}

	if err := g.generateTagPages(); err != nil {
		return app.WrapWithContext(err, app.ErrContent, app.ErrorContext{
			Operation: "generate_tag_pages",
			Component: "generator",
		})
	}

	if err := g.generateRobotsTxt(); err != nil {
		return app.WrapWithContext(err, app.ErrContent, app.ErrorContext{
			Operation: "generate_robots_txt",
			Component: "generator",
		})
	}

	if err := g.generateAtomFeed(); err != nil {
		return app.WrapWithContext(err, app.ErrContent, app.ErrorContext{
			Operation: "generate_atom_feed",
			Component: "generator",
		})
	}

	if err := g.generateSitemap(); err != nil {
		return app.WrapWithContext(err, app.ErrContent, app.ErrorContext{
			Operation: "generate_sitemap",
			Component: "generator",
		})
	}

	return g.copyStaticAssets()
}

func (g *Generator) processNode(node *content.Node) error {
	switch node.Type {
	case content.NodeTypeSection:
		return g.generateSection(node)
	case content.NodeTypePage:
		return g.generatePage(node)
	}
	return nil
}

func (g *Generator) generateSection(node *content.Node) error {
	outputPath := filepath.Join(g.site.OutputDir, node.Slug, "index.html")
	nodeMap := node.ToMap()

	data := map[string]any{
		"Site": map[string]any{
			"Config": g.site.ToConfig(),
		},
		"Config":  node.Config,
		"Page":    nodeMap,
		"Section": nodeMap,
	}
	return g.renderTemplate(node, outputPath, data)
}

func (g *Generator) generatePage(node *content.Node) error {
	outputPath := filepath.Join(
		g.site.OutputDir,
		node.Slug,
		"index.html",
	)

	if node.Config["assets"] != "" {
		if err := g.copyPageAssets(node); err != nil {
			return err
		}
	}
	nodeMap := node.ToMap()
	parentMap := node.Parent.ToMap()

	data := map[string]any{
		"Site": map[string]any{
			"Config": g.site.ToConfig(),
		},
		"Config": node.Config,
		"Page":   nodeMap, "Section": parentMap,
	}

	return g.renderTemplate(node, outputPath, data)
}

func (g *Generator) generateTagPages() error {
	tagsTemplateExists := g.hasTemplate("tags.html")
	tagTemplateExists := g.hasTemplate("tag.html")

	if !tagsTemplateExists && !tagTemplateExists {
		return nil
	}

	tagsDir := filepath.Join(g.site.OutputDir, "tags")
	if err := os.MkdirAll(tagsDir, 0755); err != nil {
		return err
	}

	if tagsTemplateExists {
		if err := g.generateTagsIndex(); err != nil {
			return err
		}
	}

	if tagTemplateExists {
		for tag, pages := range g.parser.Tags {
			if err := g.generateTagPage(tag, pages); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Generator) generateTagsIndex() error {
	tags := make([]map[string]any, 0, len(g.parser.Tags))
	for tag, entry := range g.parser.Tags {
		permalink := g.buildTagPermalink(tag)
		tags = append(tags, map[string]any{
			"Name":      tag,
			"Count":     entry.Count,
			"Permalink": permalink,
			"URL":       g.buildTagURL(permalink),
		})
	}

	slices.SortFunc(tags, func(a, b map[string]any) int {
		return strings.Compare(a["Name"].(string), b["Name"].(string))
	})

	data := map[string]any{
		"Site": map[string]any{
			"Config": g.site.ToConfig(),
		},
		"Page": map[string]any{
			"Title": "Tags",
			"Tags":  tags,
			"Path":  "/tags",
		},
	}

	return g.renderTemplate(nil,
		filepath.Join(g.site.OutputDir, "tags", "index.html"),
		data, "tags.html")
}

func (g *Generator) generateTagPage(tag string, entry *parser.TagEntry) error {
	sortedPages := make([]*content.Node, len(entry.Pages))
	copy(sortedPages, entry.Pages)
	slices.SortFunc(sortedPages, func(a, b *content.Node) int {
		dateA := a.Config["date"].(time.Time)
		dateB := b.Config["date"].(time.Time)
		// Sort descending (newest first), so reverse comparison
		return dateB.Compare(dateA)
	})

	data := map[string]any{
		"Site": map[string]any{
			"Config": g.site.ToConfig(),
		},
		"Page": map[string]any{
			"Title":     fmt.Sprintf("Tag: %s", tag),
			"Tag":       tag,
			"Pages":     sortedPages,
			"Permalink": g.buildTagPermalink(tag),
			"Path":      g.buildTagPermalink(tag),
		},
	}

	outputPath := filepath.Join(
		g.site.OutputDir,
		"tags",
		funcs.Urlize(tag),
		"index.html",
	)

	return g.renderTemplate(nil, outputPath, data, "tag.html")
}

func (g *Generator) renderTemplate(node *content.Node, outputPath string, data any, templateNames ...string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	var tpl *template.Template

	if node != nil {
		for _, tplName := range node.TemplateChain() {
			tpl = g.templ.Lookup(tplName)
			if tpl != nil {
				break
			}
		}
	} else if len(templateNames) > 0 {
		for _, tplName := range templateNames {
			tpl = g.templ.Lookup(tplName)
			if tpl != nil {
				break
			}
		}
	} else {
		tpl = g.templ.Lookup("404.html")
	}

	if tpl == nil {
		return app.WrapWithContext(fmt.Errorf("no template found"), app.ErrTemplate, app.ErrorContext{
			Operation: "find_template",
			Component: "generator",
			Path:      outputPath,
		})
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
			Operation: "create_output_directory",
			Component: "generator",
			Path:      filepath.Dir(outputPath),
		})
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return app.WrapWithContext(err, app.ErrTemplate, app.ErrorContext{
			Operation: "execute_template",
			Component: "generator",
			Path:      outputPath,
		})
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
			Operation: "write_html_output",
			Component: "generator",
			Path:      outputPath,
		})
	}

	return nil
}

func (g *Generator) copyPageAssets(node *content.Node) error {
	assetsValue, exists := node.Config["assets"]
	if !exists {
		return nil // No assets to copy
	}

	assets, ok := assetsValue.(string)
	if !ok {
		return nil // Assets value is not a string
	}

	if _, err := os.Stat(assets); os.IsNotExist(err) {
		return nil // Skip missing assets
	}

	dest := filepath.Join(
		g.site.OutputDir,
		node.Slug,
		filepath.Base(assets),
	)

	return copyDir(assets, dest)
}

func (g *Generator) walkNodes(node *content.Node, fn func(*content.Node)) {
	fn(node)
	for _, child := range node.Children {
		g.walkNodes(child, fn)
	}
}

func (g *Generator) copyStaticAssets() error {
	return filepath.WalkDir("static", func(srcPath string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel("static", srcPath)
		if err != nil {
			return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
				Operation: "get_relative_path",
				Component: "generator",
				Path:      srcPath,
			})
		}

		destPath := filepath.Join(g.site.OutputDir, relPath)
		return copyFile(srcPath, destPath)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
			Operation: "create_directory",
			Component: "generator",
			Path:      filepath.Dir(dst),
		})
	}

	in, err := os.Open(src)
	if err != nil {
		return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
			Operation: "open_source_file",
			Component: "generator",
			Path:      src,
		})
	}
	defer func() {
		_ = in.Close()
	}()

	out, err := os.Create(dst)
	if err != nil {
		return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
			Operation: "create_destination_file",
			Component: "generator",
			Path:      dst,
		})
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err = io.Copy(out, in); err != nil {
		return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
			Operation: "copy_file_content",
			Component: "generator",
			Path:      fmt.Sprintf("%s -> %s", src, dst),
		})
	}
	return nil
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
				Operation: "walk_directory",
				Component: "generator",
				Path:      path,
			})
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
				Operation: "get_relative_path",
				Component: "generator",
				Path:      path,
			})
		}

		target := filepath.Join(dst, relPath)
		if d.IsDir() {
			if err := os.MkdirAll(target, 0755); err != nil {
				return app.WrapWithContext(err, app.ErrFileSystem, app.ErrorContext{
					Operation: "create_directory",
					Component: "generator",
					Path:      target,
				})
			}
			return nil
		}

		return copyFile(path, target)
	})
}

func (g *Generator) generate404Page() error {
	tpl := g.templ.Lookup("404.html")
	if tpl == nil {
		return nil
	}

	outputPath := filepath.Join(g.site.OutputDir, "404.html")
	data := map[string]any{
		"Site": map[string]any{
			"Config": g.site.ToConfig(),
		},
		"Page": map[string]any{
			"Title": "Page Not Found",
			"URL":   path.Join(g.site.BaseURL, "404.html"),
		},
	}

	return g.renderTemplate(nil, outputPath, data)
}

func (g *Generator) buildTagPermalink(tag string) string {
	return path.Join("/tags", funcs.Urlize(tag)) + "/"
}

func (g *Generator) buildTagURL(tagLink string) string {
	return g.site.BaseURL + tagLink
}

func (g *Generator) hasTemplate(name string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.templ.Lookup(name) != nil
}

func (g *Generator) generateRobotsTxt() error {
	content := fmt.Sprintf(`User-agent: *
Allow: /
Sitemap: %s/sitemap.xml
`, g.site.BaseURL)

	return os.WriteFile(
		filepath.Join(g.site.OutputDir, "robots.txt"),
		[]byte(content),
		0644,
	)
}

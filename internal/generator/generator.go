package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/content"
	"github.com/tduyng/gozzi/internal/parser"
)

type Generator struct {
	site   *config.Site
	templ  *template.Template
	parser *parser.ContentParser
	mu     sync.Mutex
}

func NewGenerator(site *config.Site, parser *parser.ContentParser) (*Generator, error) {
	gen := &Generator{
		site:   site,
		parser: parser,
	}

	tmpl, err := template.New("").Funcs(gen.CreateFuncMap()).ParseGlob("templates/*.html")
	if err != nil {
		return nil, fmt.Errorf("template parsing failed: %w", err)
	}

	gen.templ = tmpl
	return gen, nil
}

func (g *Generator) Generate(contentRoot *content.Node) error {
	if err := os.RemoveAll(g.site.OutputDir); err != nil {
		return fmt.Errorf("failed to clean output: %w", err)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	g.walkNodes(contentRoot, func(n *content.Node) {
		wg.Add(1)
		go func(node *content.Node) {
			defer wg.Done()
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
		return fmt.Errorf("failed to generate 404 page: %w", err)
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

	node.Permalink = g.buildPermalink(node)
	node.URL = g.buildURL(node)

	if node.Higher != nil {
		node.Higher.URL = g.buildURL(node.Higher)
		node.Higher.Permalink = g.buildPermalink(node.Higher)
	}
	if node.Lower != nil {
		node.Lower.URL = g.buildURL(node.Lower)
		node.Lower.Permalink = g.buildPermalink(node.Lower)
	}
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
		node.Parent.Slug,
		node.Slug,
		"index.html",
	)

	if node.Config["assets"] != "" {
		if err := g.copyPageAssets(node); err != nil {
			return err
		}
	}
	node.Permalink = g.buildPermalink(node)
	node.URL = g.buildURL(node)
	node.Parent.Permalink = g.buildPermalink(node.Parent)
	node.Parent.URL = g.buildURL(node.Parent)
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

func (g *Generator) renderTemplate(node *content.Node, outputPath string, data any) error {
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
	} else {
		tpl = g.templ.Lookup("404.html")
	}

	if tpl == nil {
		return fmt.Errorf("no template found for path: %s", outputPath)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template execution: %w", err)
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write html output: %w", err)
	}

	return nil
}

func (g *Generator) copyPageAssets(node *content.Node) error {
	assets := node.Config["assets"].(string)
	if _, err := os.Stat(assets); os.IsNotExist(err) {
		return nil // Skip missing assets
	}

	dest := filepath.Join(
		g.site.OutputDir,
		node.Parent.Slug,
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
	return filepath.Walk("static", func(srcPath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel("static", srcPath)
		if err != nil {
			return fmt.Errorf("get relative path: %w", err)
		}

		destPath := filepath.Join(g.site.OutputDir, relPath)
		return copyFile(srcPath, destPath)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source: %w", err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create dest: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("copy content: %w", err)
	}
	return nil
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return fmt.Errorf("get relative path: %w", err)
		}

		target := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		return copyFile(path, target)
	})
}

func (g *Generator) buildPermalink(node *content.Node) string {
	parts := []string{}
	for n := node; n != nil; n = n.Parent {
		if n.Slug != "/" {
			parts = append([]string{n.Slug}, parts...)
		}
	}
	return path.Join("/", path.Join(parts...)) + "/"
}

func (g *Generator) buildURL(node *content.Node) string {
	return path.Join(g.site.BaseURL, node.Permalink) + "/"
}

func (g *Generator) ReloadTemplates() error {
	tmpl, err := template.New("").Funcs(g.CreateFuncMap()).ParseGlob("templates/*.html")
	if err != nil {
		return fmt.Errorf("template reload failed: %w", err)
	}
	g.mu.Lock()
	g.templ = tmpl
	g.mu.Unlock()
	return nil
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

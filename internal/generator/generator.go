package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/content"
)

type Generator struct {
	site  *config.Site
	templ *template.Template
	mu    sync.Mutex
}

func NewGenerator(site *config.Site) (*Generator, error) {
	gen := &Generator{
		site: site,
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
	var outputPath string
	if node.Path == "." {
		outputPath = filepath.Join(g.site.OutputDir, "index.html")
	} else {
		outputPath = filepath.Join(g.site.OutputDir, node.Slug, "index.html")
	}
	tmplName := getTemplateName(node)
	return g.renderTemplate(tmplName, outputPath, map[string]any{
		"Site":    g.site,
		"Section": node,
		"Pages":   node.Children,
	})
}

func (g *Generator) generatePage(node *content.Node) error {
	basePath := g.site.OutputDir
	if node.Parent != nil {
		basePath = filepath.Join(basePath, node.Parent.Slug)
	}
	outputPath := filepath.Join(basePath, node.Slug, "index.html")
	tmplName := getTemplateName(node)
	if node.PageMeta.Assets != "" {
		if _, err := os.Stat(node.PageMeta.Assets); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("error checking assets directory: %w", err)
			}
		} else {
			assetFolderName := filepath.Base(node.PageMeta.Assets)
			var dest string
			if node.Parent != nil && node.Slug == node.Parent.Slug {
				dest = filepath.Join(g.site.OutputDir, node.Parent.Slug, assetFolderName)
			} else {
				dest = filepath.Join(g.site.OutputDir, node.Parent.Slug, node.Slug, assetFolderName)
			}
			if err := copyDir(node.PageMeta.Assets, dest); err != nil {
				return fmt.Errorf("failed to copy assets: %w", err)
			}
		}
	}

	return g.renderTemplate(tmplName, outputPath, map[string]any{
		"Site":    g.site,
		"Page":    node,
		"Section": node.Parent,
	})
}

func (g *Generator) renderTemplate(name, outputPath string, data any) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	tmpl := g.templ.Lookup(name)
	if tmpl == nil {
		return fmt.Errorf("template %q not found", name)
	}

	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	return os.WriteFile(outputPath, buf.Bytes(), 0644)
}

func (g *Generator) walkNodes(node *content.Node, fn func(*content.Node)) {
	fn(node)
	for _, child := range node.Children {
		g.walkNodes(child, fn)
	}
}

func (g *Generator) copyStaticAssets() error {
	return filepath.Walk("static", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel("static", path)
		dest := filepath.Join(g.site.OutputDir, relPath)

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}

		return copyFile(path, dest)
	})
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}

		return copyFile(path, target)
	})
}

func getTemplateName(node *content.Node) string {
	if custom, ok := node.Config["template"].(string); ok && custom != "" {
		return custom
	}
	if node.Type == content.NodeTypeSection {
		return "default.html"
	}
	return "post.html"
}

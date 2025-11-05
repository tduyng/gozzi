// Package builder handles site generation including templates, feeds, and static files.
// Main orchestrator for building static sites from content and templates.
package builder

import (
	"fmt"
	"html/template"
	"os"
	"runtime"
	"sync"

	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
	tplengine "github.com/tduyng/gozzi/app/template"
	"github.com/tduyng/gozzi/shared"
)

// Builder handles site generation including templates, feeds, and static files.
type Builder struct {
	site   *config.Site
	templ  *template.Template
	parser *parser.ContentParser
	engine *tplengine.Engine
	mu     sync.Mutex
}

// NewBuilder creates a new Builder with loaded templates.
func NewBuilder(site *config.Site, parser *parser.ContentParser) (*Builder, error) {
	// Initialize template engine
	engine := tplengine.NewEngine(&tplengine.EngineConfig{
		BaseURL:         site.BaseURL,
		ContentMap:      parser.ContentMap,
		Markdown:        parser.GetMarkdownProcessor(),
		StrictTemplates: site.StrictTemplates,
	})

	b := &Builder{
		site:   site,
		parser: parser,
		engine: engine,
	}

	tmpl, err := b.loadTemplates()
	if err != nil {
		return nil, shared.WrapWithContext(err, shared.ErrTemplate, shared.ErrorContext{
			Operation: "template_initialization",
			Component: "builder",
		})
	}

	b.templ = tmpl
	return b, nil
}

// ReloadTemplates reloads all templates from disk for development mode.
func (b *Builder) ReloadTemplates() error {
	tmpl, err := b.loadTemplates()
	if err != nil {
		return shared.WrapWithContext(err, shared.ErrTemplate, shared.ErrorContext{
			Operation: "template_reload",
			Component: "builder",
		})
	}

	b.mu.Lock()
	b.templ = tmpl
	b.mu.Unlock()
	return nil
}

// Generate processes the content tree and generates the complete static site.
func (b *Builder) Generate(contentRoot *content.Node) error {
	if err := os.RemoveAll(b.site.OutputDir); err != nil {
		return shared.WrapWithContext(err, shared.ErrFileSystem, shared.ErrorContext{
			Operation: "clean_output_directory",
			Component: "builder",
			Path:      b.site.OutputDir,
		})
	}

	sem := make(chan struct{}, runtime.NumCPU()*2)
	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	b.walkNodes(contentRoot, func(n *content.Node) {
		wg.Add(1)
		sem <- struct{}{}

		go func(node *content.Node) {
			defer func() {
				<-sem
				wg.Done()
			}()

			if err := b.processNode(node); err != nil {
				errChan <- err
			}
		}(n)
	})

	wg.Wait()
	close(errChan)

	for err := range errChan {
		fmt.Printf("Processing error: %v\n", err)
	}

	if err := b.generate404Page(); err != nil {
		return shared.WrapWithContext(err, shared.ErrContent, shared.ErrorContext{
			Operation: "generate_404_page",
			Component: "builder",
		})
	}

	if err := b.generateTagPages(); err != nil {
		return shared.WrapWithContext(err, shared.ErrContent, shared.ErrorContext{
			Operation: "generate_tag_pages",
			Component: "builder",
		})
	}

	if err := b.generateRobotsTxt(); err != nil {
		return shared.WrapWithContext(err, shared.ErrContent, shared.ErrorContext{
			Operation: "generate_robots_txt",
			Component: "builder",
		})
	}

	if err := b.generateAtomFeed(); err != nil {
		return shared.WrapWithContext(err, shared.ErrContent, shared.ErrorContext{
			Operation: "generate_atom_feed",
			Component: "builder",
		})
	}

	if err := b.generateSitemap(); err != nil {
		return shared.WrapWithContext(err, shared.ErrContent, shared.ErrorContext{
			Operation: "generate_sitemap",
			Component: "builder",
		})
	}

	return b.copyStaticAssets()
}

func (b *Builder) processNode(node *content.Node) error {
	switch node.Type {
	case content.NodeTypeSection:
		return b.generateSection(node)
	case content.NodeTypePage:
		return b.generatePage(node)
	}
	return nil
}

func (b *Builder) walkNodes(node *content.Node, fn func(*content.Node)) {
	fn(node)
	for _, child := range node.Children {
		b.walkNodes(child, fn)
	}
}

func (b *Builder) hasTemplate(name string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.templ.Lookup(name) != nil
}

func (b *Builder) generateRobotsTxt() error {
	content := fmt.Sprintf(`User-agent: *
Allow: /
Sitemap: %s/sitemap.xml
`, b.site.BaseURL)

	return os.WriteFile(
		b.site.OutputDir+"/robots.txt",
		[]byte(content),
		0644,
	)
}

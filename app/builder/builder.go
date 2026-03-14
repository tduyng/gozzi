package builder

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
	tplengine "github.com/tduyng/gozzi/app/template"
	"github.com/tduyng/gozzi/app/utils"
)

type Builder struct {
	site             *config.Site
	templ            *template.Template
	parser           *parser.ContentParser
	engine           *tplengine.Engine
	mu               sync.Mutex
	cachedTaxonomies map[string]any
	taxonomiesOnce   sync.Once
}

// GenerateOptions configures how the site generation should run.
type GenerateOptions struct {
	ContentDir        string
	OldTaxonomyValues map[string]map[string]any // Pre-snapshotted taxonomy values: relPath -> {field -> value}
}

// NewBuilder creates a new Builder with loaded templates.
func NewBuilder(site *config.Site, parser *parser.ContentParser) (*Builder, error) {
	engine := tplengine.NewEngine(&tplengine.EngineConfig{
		BaseURL:         site.BaseURL,
		ContentMap:      parser.ContentMap,
		Markdown:        parser.GetMarkdownProcessor(),
		StrictTemplates: site.StrictTemplates,
		I18n:            site.I18n,
	})

	b := &Builder{
		site:   site,
		parser: parser,
		engine: engine,
	}

	tmpl, err := b.loadTemplates()
	if err != nil {
		return nil, utils.WrapWithContext(err, utils.ErrTemplate, utils.ErrorContext{
			Operation: "template_initialization",
			Component: "builder",
		})
	}

	b.templ = tmpl

	// Initialize shortcode support in markdown processor
	parser.SetShortcodeTemplates(tmpl)

	return b, nil
}

// ReloadTemplates reloads all templates from disk for development mode.
func (b *Builder) ReloadTemplates() error {
	tmpl, err := b.loadTemplates()
	if err != nil {
		return utils.WrapWithContext(err, utils.ErrTemplate, utils.ErrorContext{
			Operation: "template_reload",
			Component: "builder",
		})
	}

	b.mu.Lock()
	b.templ = tmpl
	b.mu.Unlock()

	// Update shortcodes in markdown processor
	b.parser.SetShortcodeTemplates(tmpl)

	return nil
}

// Generate processes the content tree and generates the complete static site.
func (b *Builder) Generate(contentRoot *content.Node) error {
	return b.fullGenerate(contentRoot, false)
}

// GenerateWithOptions processes the content tree with specific generation options.
// For backwards compatibility, this always performs a full build.
func (b *Builder) GenerateWithOptions(contentRoot *content.Node, _ GenerateOptions) error {
	return b.fullGenerate(contentRoot, false)
}

// GenerateClean performs a full rebuild and cleans the output directory first.
// Use this after file deletions to remove stale content.
func (b *Builder) GenerateClean(contentRoot *content.Node) error {
	return b.fullGenerate(contentRoot, true)
}

// fullGenerate performs a complete site rebuild.
// If cleanOutput is true, it removes stale files from previous builds.
func (b *Builder) fullGenerate(contentRoot *content.Node, cleanOutput bool) error {
	// Reset taxonomies cache for each build
	b.taxonomiesOnce = sync.Once{}
	b.cachedTaxonomies = nil

	// Clean output directory only when needed (e.g., after file deletions)
	if cleanOutput {
		if err := b.cleanOutputDir(); err != nil {
			return err
		}
	}

	if contentRoot == nil {
		if err := os.MkdirAll(b.site.OutputDir, 0755); err != nil {
			return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
				Operation: "create_output_directory",
				Component: "builder",
				Path:      b.site.OutputDir,
			})
		}
		return b.copyStaticAssets()
	}

	if err := os.MkdirAll(b.site.OutputDir, 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_output_directory",
			Component: "builder",
			Path:      b.site.OutputDir,
		})
	}

	// Count nodes to properly size the error channel
	nodeCount := b.countNodes(contentRoot)

	sem := make(chan struct{}, runtime.NumCPU()*2)
	var wg sync.WaitGroup
	// Size error channel to hold all potential errors (one per node)
	// This prevents goroutine deadlock if many nodes fail processing
	errChan := make(chan error, nodeCount)

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
		return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
			Operation: "generate_404_page",
			Component: "builder",
		})
	}

	if err := b.generateTaxonomyPages(); err != nil {
		return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
			Operation: "generate_taxonomy_pages",
			Component: "builder",
		})
	}

	if err := b.generateRobotsTxt(); err != nil {
		return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
			Operation: "generate_robots_txt",
			Component: "builder",
		})
	}

	if err := b.generateAtomFeed(); err != nil {
		return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
			Operation: "generate_atom_feed",
			Component: "builder",
		})
	}

	if err := b.generateSitemap(); err != nil {
		return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
			Operation: "generate_sitemap",
			Component: "builder",
		})
	}

	if b.site.GenerateSearch {
		if err := b.generateSearchIndex(); err != nil {
			return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
				Operation: "generate_search_index",
				Component: "builder",
			})
		}
	}

	return b.copyStaticAssets()
}

// cleanOutputDir removes all files in output directory to ensure clean rebuild.
func (b *Builder) cleanOutputDir() error {
	outputDir := b.site.OutputDir

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		return nil
	}

	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "read_output_directory",
			Component: "builder",
			Path:      outputDir,
		})
	}

	for _, entry := range entries {
		entryPath := filepath.Join(outputDir, entry.Name())

		if entry.IsDir() {
			if err := os.RemoveAll(entryPath); err != nil {
				log.Printf("Warning: Could not remove directory %s: %v", entryPath, err)
			}
		} else {
			if err := os.Remove(entryPath); err != nil {
				log.Printf("Warning: Could not remove file %s: %v", entryPath, err)
			}
		}
	}

	return nil
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
	if node == nil {
		return
	}
	fn(node)
	for _, child := range node.Children {
		b.walkNodes(child, fn)
	}
}

func (b *Builder) countNodes(node *content.Node) int {
	if node == nil {
		return 0
	}
	count := 1
	for _, child := range node.Children {
		count += b.countNodes(child)
	}
	return count
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

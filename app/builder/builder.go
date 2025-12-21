package builder

import (
	"fmt"
	"html/template"
	"os"
	"runtime"
	"sync"

	"github.com/tduyng/gozzi/app/cache"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
	tplengine "github.com/tduyng/gozzi/app/template"
	"github.com/tduyng/gozzi/app/utils"
)

type Builder struct {
	site        *config.Site
	templ       *template.Template
	parser      *parser.ContentParser
	engine      *tplengine.Engine
	renderCache *cache.RenderCache // Template render cache for incremental builds
	mu          sync.Mutex
}

// GenerateOptions configures how the site generation should run
type GenerateOptions struct {
	// Incremental mode only generates changed pages instead of full site rebuild
	Incremental bool

	// ChangedFiles lists the files that changed (for incremental mode)
	ChangedFiles []string

	// ContentDir is the base content directory (needed for relative path calculation)
	ContentDir string
}

// NewBuilder creates a new Builder with loaded templates.
func NewBuilder(site *config.Site, parser *parser.ContentParser) (*Builder, error) {
	engine := tplengine.NewEngine(&tplengine.EngineConfig{
		BaseURL:         site.BaseURL,
		ContentMap:      parser.ContentMap,
		Markdown:        parser.GetMarkdownProcessor(),
		StrictTemplates: site.StrictTemplates,
	})

	b := &Builder{
		site:        site,
		parser:      parser,
		engine:      engine,
		renderCache: cache.NewRenderCache(),
	}

	tmpl, err := b.loadTemplates()
	if err != nil {
		return nil, utils.WrapWithContext(err, utils.ErrTemplate, utils.ErrorContext{
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
		return utils.WrapWithContext(err, utils.ErrTemplate, utils.ErrorContext{
			Operation: "template_reload",
			Component: "builder",
		})
	}

	b.mu.Lock()
	b.templ = tmpl
	// Note: Don't clear render cache here - let caller do selective invalidation
	// for better incremental build performance
	b.mu.Unlock()
	return nil
}

// InvalidateTemplateCache invalidates cached renders for specific templates.
// This is more efficient than clearing the entire cache when only specific templates changed.
func (b *Builder) InvalidateTemplateCache(templateNames []string) int {
	b.mu.Lock()
	defer b.mu.Unlock()

	totalInvalidated := 0
	for _, tplName := range templateNames {
		count := b.renderCache.InvalidateTemplate(tplName)
		totalInvalidated += count
	}
	return totalInvalidated
}

// Generate processes the content tree and generates the complete static site.
// This is a convenience wrapper around GenerateWithOptions with default settings.
func (b *Builder) Generate(contentRoot *content.Node) error {
	return b.GenerateWithOptions(contentRoot, GenerateOptions{
		Incremental: false, // Default to full build for backwards compatibility
	})
}

// GenerateWithOptions processes the content tree with specific generation options.
// Supports both full builds and incremental builds based on options.
func (b *Builder) GenerateWithOptions(contentRoot *content.Node, opts GenerateOptions) error {
	// Route to appropriate generation strategy
	if opts.Incremental {
		return b.incrementalGenerate(contentRoot, opts)
	}
	return b.fullGenerate(contentRoot)
}

// fullGenerate performs a complete site rebuild (original Generate logic)
func (b *Builder) fullGenerate(contentRoot *content.Node) error {
	// Silently handle nil contentRoot (empty content directory)
	if contentRoot == nil {
		// Just generate auxiliary pages (404, robots.txt, etc.)
		// Skip content-dependent pages
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
		return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
			Operation: "generate_404_page",
			Component: "builder",
		})
	}

	if err := b.generateTagPages(); err != nil {
		return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
			Operation: "generate_tag_pages",
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

	return b.copyStaticAssets()
}

// incrementalGenerate performs a selective rebuild (only changed content)
func (b *Builder) incrementalGenerate(contentRoot *content.Node, opts GenerateOptions) error {
	if contentRoot == nil {
		return b.fullGenerate(contentRoot)
	}

	// Create output directory
	if err := os.MkdirAll(b.site.OutputDir, 0755); err != nil {
		return utils.WrapWithContext(err, utils.ErrFileSystem, utils.ErrorContext{
			Operation: "create_output_directory",
			Component: "builder",
			Path:      b.site.OutputDir,
		})
	}

	// Default contentDir to "content" if not specified
	contentDir := opts.ContentDir
	if contentDir == "" {
		contentDir = "content"
	}

	// Analyze what changed to determine what needs regeneration
	tracker := NewChangeTracker(b.parser.ContentMap, b.parser)
	tracker.AnalyzeChanges(opts.ChangedFiles, contentDir)

	// CRITICAL: Get changed nodes AFTER parsing completes!
	// ParseFiles creates NEW node objects, so we need fresh pointers to updated nodes
	// Otherwise we'll be rendering old node data with old content
	changedNodes := tracker.GetChangedNodesAfterParse(b.parser.ContentMap)

	// Process changed nodes in parallel (same pattern as fullGenerate)
	sem := make(chan struct{}, runtime.NumCPU()*2)
	var wg sync.WaitGroup
	errChan := make(chan error, 100)

	for _, node := range changedNodes {
		wg.Add(1)
		sem <- struct{}{}

		go func(n *content.Node) {
			defer func() {
				<-sem
				wg.Done()
			}()

			if err := b.processNode(n); err != nil {
				errChan <- err
			}
		}(node)
	}

	wg.Wait()
	close(errChan)

	// Report any errors
	for err := range errChan {
		fmt.Printf("Processing error: %v\n", err)
	}

	// Selective auxiliary page generation
	// Only regenerate what's affected by the changes

	// 404 page rarely needs regeneration (only on template changes)
	// Skip for content-only changes
	if tracker.ShouldRegenerate404() {
		if err := b.generate404Page(); err != nil {
			return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
				Operation: "generate_404_page",
				Component: "builder",
			})
		}
	}

	// Homepage: regenerate if blog posts changed (featured posts display on homepage)
	if tracker.ShouldRegenerateHome() {
		if homeNode, exists := b.parser.ContentMap["."]; exists {
			if err := b.processNode(homeNode); err != nil {
				return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
					Operation: "regenerate_homepage",
					Component: "builder",
				})
			}
		}
	}

	// Tag pages: only regenerate affected tags
	if tracker.GetAffectedTagsCount() > 0 {
		if err := b.generateSelectiveTags(tracker.GetAffectedTags()); err != nil {
			return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
				Operation: "generate_selective_tags",
				Component: "builder",
			})
		}
	}

	// Robots.txt rarely changes
	// Skip unless config changed
	if tracker.ShouldRegenerateRobots() {
		if err := b.generateRobotsTxt(); err != nil {
			return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
				Operation: "generate_robots_txt",
				Component: "builder",
			})
		}
	}

	// Feed: only if blog posts changed
	if tracker.ShouldRegenerateFeed() {
		if err := b.generateAtomFeed(); err != nil {
			return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
				Operation: "generate_atom_feed",
				Component: "builder",
			})
		}
	}

	// Sitemap: regenerate if any content changed
	if tracker.ShouldRegenerateSitemap() {
		if err := b.generateSitemap(); err != nil {
			return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
				Operation: "generate_sitemap",
				Component: "builder",
			})
		}
	}

	// Static assets: skip unless explicitly changed
	// This is a conservative approach - could be optimized further
	return b.copyStaticAssets()
}

// GetCacheStats returns the current render cache statistics.
func (b *Builder) GetCacheStats() cache.RenderCacheStats {
	return b.renderCache.Stats()
}

// ResetCacheStats resets the hit/miss counters for per-build statistics.
func (b *Builder) ResetCacheStats() {
	b.renderCache.ResetStats()
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

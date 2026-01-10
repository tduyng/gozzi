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

// GenerateOptions configures how the site generation should run.
type GenerateOptions struct {
	Incremental       bool
	ChangedFiles      []string
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

// InvalidateTemplateCache invalidates cached renders for specific templates.
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
func (b *Builder) Generate(contentRoot *content.Node) error {
	return b.GenerateWithOptions(contentRoot, GenerateOptions{
		Incremental: false, // Default to full build for backwards compatibility
	})
}

// GenerateWithOptions processes the content tree with specific generation options.
func (b *Builder) GenerateWithOptions(contentRoot *content.Node, opts GenerateOptions) error {
	if opts.Incremental {
		return b.incrementalGenerate(contentRoot, opts)
	}
	return b.fullGenerate(contentRoot)
}

// fullGenerate performs a complete site rebuild.
func (b *Builder) fullGenerate(contentRoot *content.Node) error {
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

	return b.copyStaticAssets()
}

// incrementalGenerate performs a selective rebuild.
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

	contentDir := opts.ContentDir
	if contentDir == "" {
		contentDir = "content"
	}

	analyzer := NewRebuildAnalyzer(b.parser.ContentMap, b.parser, contentDir)

	scope := analyzer.AnalyzeChanges(opts.ChangedFiles, opts.OldTaxonomyValues)

	changedNodes := scope.GetNodes()

	sem := make(chan struct{}, runtime.NumCPU()*2)
	var wg sync.WaitGroup
	// Size error channel to match number of changed nodes
	// This prevents goroutine deadlock if many nodes fail processing
	errChan := make(chan error, len(changedNodes))

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

	for err := range errChan {
		fmt.Printf("Processing error: %v\n", err)
	}

	if scope.NeedsGlobal("404") {
		if err := b.generate404Page(); err != nil {
			return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
				Operation: "generate_404_page",
				Component: "builder",
			})
		}
	}

	if scope.NeedsGlobal("home") {
		if homeNode, exists := b.parser.ContentMap["."]; exists {
			if err := b.processNode(homeNode); err != nil {
				return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
					Operation: "regenerate_homepage",
					Component: "builder",
				})
			}
		}
	}

	if scope.NeedsGlobal("blog-listing") {
		if blogNode, exists := b.parser.ContentMap["blog"]; exists {
			if err := b.processNode(blogNode); err != nil {
				return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
					Operation: "regenerate_blog_listing",
					Component: "builder",
				})
			}
		}
	}

	// Regenerate affected taxonomy pages
	affectedTaxonomies := scope.GetAffectedTaxonomies()
	for taxonomyName, termSlugs := range affectedTaxonomies {
		if err := b.generateSelectiveTaxonomies(taxonomyName, termSlugs); err != nil {
			return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
				Operation: fmt.Sprintf("generate_selective_%s", taxonomyName),
				Component: "builder",
			})
		}
	}

	if scope.NeedsGlobal("robots") {
		if err := b.generateRobotsTxt(); err != nil {
			return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
				Operation: "generate_robots_txt",
				Component: "builder",
			})
		}
	}

	if scope.NeedsGlobal("feed") {
		if err := b.generateAtomFeed(); err != nil {
			return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
				Operation: "generate_atom_feed",
				Component: "builder",
			})
		}
	}

	if scope.NeedsGlobal("sitemap") {
		if err := b.generateSitemap(); err != nil {
			return utils.WrapWithContext(err, utils.ErrContent, utils.ErrorContext{
				Operation: "generate_sitemap",
				Component: "builder",
			})
		}
	}

	// NOTE: Static file copying is handled by the watcher in serve mode.
	// The watcher explicitly copies individual changed static files before calling
	// incrementalGenerate (see watcher.go lines 247-254).
	// Full rebuilds (via fullGenerate) still copy all static files.
	// This avoids redundant copying of potentially thousands of static files.
	return nil
}

// SnapshotTaxonomyValues captures current taxonomy values for changed files.
// This should be called BEFORE ParseFiles() to preserve old values for comparison.
func (b *Builder) SnapshotTaxonomyValues(changedFiles []string, contentDir string) map[string]map[string]any {
	analyzer := NewRebuildAnalyzer(b.parser.ContentMap, b.parser, contentDir)
	return analyzer.SnapshotTaxonomyValues(changedFiles)
}

// GetCacheStats returns the current render cache statistics.
func (b *Builder) GetCacheStats() cache.RenderCacheStats {
	return b.renderCache.Stats()
}

// ResetCacheStats resets the hit/miss counters.
func (b *Builder) ResetCacheStats() {
	b.renderCache.ResetStats()
}

// ClearRenderCache clears all cached rendered content.
func (b *Builder) ClearRenderCache() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.renderCache.Clear()
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

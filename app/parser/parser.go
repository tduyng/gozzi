// Package parser provides content parsing and markdown processing with concurrent file processing.
package parser

import (
	"context"
	"html/template"
	"io/fs"
	"math"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/tduyng/gozzi/app/cache"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/markdown"
	"github.com/tduyng/gozzi/app/paginate"
	"github.com/tduyng/gozzi/app/utils"
	"github.com/yuin/goldmark"
	highlight "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// ParseStats tracks incremental parsing statistics
type ParseStats struct {
	FilesSkipped atomic.Uint64 // Files skipped due to unchanged content
	FilesParsed  atomic.Uint64 // Files actually parsed
	TotalFiles   atomic.Uint64 // Total files encountered
}

// ContentParser parses markdown content files and builds the content tree.
type ContentParser struct {
	Site               *config.Site
	ContentMap         map[string]*content.Node
	Tags               map[string]*TagEntry // Deprecated: Use Taxonomies["tags"] instead
	Taxonomies         map[string]*Taxonomy // All taxonomies (tags, categories, series, custom)
	mu                 sync.Mutex
	md                 goldmark.Markdown
	shortcodeProcessor *markdown.ShortcodeProcessor
	hashCache          *cache.HashCache // Content hash cache for incremental parsing
	stats              *ParseStats      // Statistics for monitoring
}

// NewParser creates a new ContentParser with the given site configuration.
func NewParser(cfg *config.Site) *ContentParser {
	syntaxTheme := cfg.SyntaxTheme
	if syntaxTheme == "" {
		syntaxTheme = "dracula"
	}

	return &ContentParser{
		Site:       cfg,
		ContentMap: make(map[string]*content.Node),
		Tags:       make(map[string]*TagEntry),
		Taxonomies: make(map[string]*Taxonomy),
		hashCache:  cache.NewHashCache(),
		stats:      &ParseStats{},
		md: goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Footnote,
				highlight.NewHighlighting(
					highlight.WithGuessLanguage(true),
					highlight.WithStyle(syntaxTheme),
				),
				markdown.NewMathExtension(),
				markdown.NewMermaidExtension(),
				markdown.NewTocExtension(),
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		),
	}
}

// Parse walks the content directory and parses all markdown files.
func (p *ContentParser) Parse(rootDir string) error {
	p.mu.Lock()
	for k := range p.ContentMap {
		delete(p.ContentMap, k)
	}
	for k := range p.Taxonomies {
		delete(p.Taxonomies, k)
	}
	for k := range p.Tags {
		delete(p.Tags, k)
	}
	p.mu.Unlock()

	var files []string
	_ = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if filepath.Base(path) == "_index.md" || filepath.Ext(path) == ".md" {
			files = append(files, path)
		}
		return nil
	})

	// Sort files for deterministic parsing order
	slices.Sort(files)

	return p.parseFiles(rootDir, files)
}

// ParseFiles parses only the specified files for incremental rebuilds.
func (p *ContentParser) ParseFiles(rootDir string, files []string) error {
	var mdFiles []string
	for _, f := range files {
		if filepath.Ext(f) == ".md" || filepath.Base(f) == "_index.md" {
			mdFiles = append(mdFiles, f)
		}
	}

	if len(mdFiles) == 0 {
		return nil
	}

	// Sort files for deterministic parsing order
	slices.Sort(mdFiles)

	return p.parseFiles(rootDir, mdFiles)
}

// parseFiles is the internal implementation that parses a list of files.
func (p *ContentParser) parseFiles(rootDir string, files []string) error {
	ctx := context.Background()
	pool := utils.NewWorkerPool(ctx)

	_ = pool.ProcessFiles(files, func(ctx context.Context, filePath string) error {
		relPath, _ := filepath.Rel(rootDir, filePath)
		dir := filepath.Dir(relPath)

		switch {
		case filepath.Base(filePath) == "_index.md":
			_ = p.parseSection(filePath, dir)
		case filepath.Ext(filePath) == ".md":
			_ = p.parsePage(filePath, dir)
		}
		return nil
	})

	paginator := paginate.New(p.ContentMap)
	paginator.BuildLinks()

	// Sort children for deterministic output
	p.sortChildren()

	return nil
}

func calculateReadStats(content string) (int, int) {
	re := regexp.MustCompile(`<[^>]*>`)
	plainText := re.ReplaceAllString(content, " ")

	words := strings.Fields(plainText)
	wordCount := len(words)

	readTime := max(int(math.Ceil(float64(wordCount)/220)), 1)

	return wordCount, readTime
}

// GetMarkdownProcessor returns the configured goldmark markdown processor.
func (p *ContentParser) GetMarkdownProcessor() goldmark.Markdown {
	return p.md
}

// SetShortcodeTemplates updates the markdown processor with shortcode support.
func (p *ContentParser) SetShortcodeTemplates(templates *template.Template) {
	syntaxTheme := p.Site.SyntaxTheme
	if syntaxTheme == "" {
		syntaxTheme = "dracula"
	}

	// Initialize shortcode processor
	p.shortcodeProcessor = markdown.NewShortcodeProcessor(templates)

	p.md = goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			highlight.NewHighlighting(
				highlight.WithGuessLanguage(true),
				highlight.WithStyle(syntaxTheme),
			),
			markdown.NewMathExtension(),
			markdown.NewMermaidExtension(),
			markdown.NewTocExtension(),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)
}

func buildPermalink(slug string) string {
	permalink := strings.Trim(slug, "/")
	if permalink == "" || permalink == "index" {
		return "/"
	}
	return "/" + permalink + "/"
}

func buildURL(baseURL, slug string) string {
	return baseURL + "/" + slug
}

// ResetStats resets parsing statistics
func (p *ContentParser) ResetStats() {
	p.stats.FilesSkipped.Store(0)
	p.stats.FilesParsed.Store(0)
	p.stats.TotalFiles.Store(0)
}

// GetHashCache returns the hash cache for external access
func (p *ContentParser) GetHashCache() *cache.HashCache {
	return p.hashCache
}

// sortChildren recursively sorts all Children slices for deterministic output
func (p *ContentParser) sortChildren() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, node := range p.ContentMap {
		p.sortNodeChildren(node)
	}
}

// sortNodeChildren recursively sorts a node's children by path
func (p *ContentParser) sortNodeChildren(node *content.Node) {
	if node == nil || len(node.Children) == 0 {
		return
	}

	// Sort children by path for deterministic order
	slices.SortFunc(node.Children, func(a, b *content.Node) int {
		if a.Path < b.Path {
			return -1
		}
		if a.Path > b.Path {
			return 1
		}
		return 0
	})

	// Recursively sort grandchildren
	for _, child := range node.Children {
		p.sortNodeChildren(child)
	}
}

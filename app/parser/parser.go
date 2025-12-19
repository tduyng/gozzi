// Package parser provides content parsing and markdown processing with concurrent file processing.
package parser

import (
	"context"
	"io/fs"
	"math"
	"path/filepath"
	"regexp"
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
	Site       *config.Site
	ContentMap map[string]*content.Node
	Tags       map[string]*TagEntry
	mu         sync.Mutex
	md         goldmark.Markdown
	hashCache  *cache.HashCache // Content hash cache for incremental parsing
	stats      *ParseStats      // Statistics for monitoring
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
	// Clear the existing map instead of creating a new one
	// This preserves the reference held by the template engine
	p.mu.Lock()
	for k := range p.ContentMap {
		delete(p.ContentMap, k)
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

	return p.parseFiles(rootDir, files)
}

// ParseFiles parses only the specified files for incremental rebuilds.
// This is used by the dev server to avoid re-parsing the entire content directory.
func (p *ContentParser) ParseFiles(rootDir string, files []string) error {
	// Filter to only markdown files
	var mdFiles []string
	for _, f := range files {
		if filepath.Ext(f) == ".md" || filepath.Base(f) == "_index.md" {
			mdFiles = append(mdFiles, f)
		}
	}

	if len(mdFiles) == 0 {
		return nil
	}

	return p.parseFiles(rootDir, mdFiles)
}

// parseFiles is the internal implementation that parses a list of files
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

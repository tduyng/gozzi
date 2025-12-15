// Package parser provides content parser orchestration and markdown processing.
// Main entry point for parsing content directory with concurrent file processing.
package parser

import (
	"context"
	"io/fs"
	"math"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

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

// ContentParser parses markdown content files and builds the content tree.
type ContentParser struct {
	Site       *config.Site
	ContentMap map[string]*content.Node
	Tags       map[string]*TagEntry
	mu         sync.Mutex
	md         goldmark.Markdown
}

// NewParser creates a new ContentParser with the given site configuration.
func NewParser(cfg *config.Site) *ContentParser {
	// Default to "dracula" theme if not specified
	syntaxTheme := cfg.SyntaxTheme
	if syntaxTheme == "" {
		syntaxTheme = "dracula"
	}

	return &ContentParser{
		Site:       cfg,
		ContentMap: make(map[string]*content.Node),
		Tags:       make(map[string]*TagEntry),
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
	p.ContentMap = make(map[string]*content.Node) // Reset ContentMap
	p.mu.Unlock()

	// Collect all files to parse
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

	// Create worker pool and process files
	ctx := context.Background()
	pool := utils.NewWorkerPool(ctx)

	// Process all markdown files concurrently
	// Note: Errors are silently ignored to preserve old behavior
	_ = pool.ProcessFiles(files, func(ctx context.Context, filePath string) error {
		relPath, _ := filepath.Rel(rootDir, filePath)
		dir := filepath.Dir(relPath)

		switch {
		case filepath.Base(filePath) == "_index.md":
			_ = p.parseSection(filePath, dir)
		case filepath.Ext(filePath) == ".md":
			_ = p.parsePage(filePath, dir)
		}
		// Preserve old behavior: silently ignore errors
		return nil
	})

	paginator := paginate.New(p.ContentMap)
	paginator.BuildLinks()
	return nil
}

func calculateReadStats(content string) (int, int) {
	// Strip HTML tags
	re := regexp.MustCompile(`<[^>]*>`)
	plainText := re.ReplaceAllString(content, " ")

	// Count words
	words := strings.Fields(plainText)
	wordCount := len(words)

	// Calculate read time (220 words/minute)
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

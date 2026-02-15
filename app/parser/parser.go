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

type ParseStats struct {
	FilesSkipped atomic.Uint64
	FilesParsed  atomic.Uint64
	TotalFiles   atomic.Uint64
}

type ContentParser struct {
	Site               *config.Site
	ContentMap         map[string]*content.Node
	Tags               map[string]*TagEntry
	Taxonomies         map[string]*Taxonomy
	mu                 sync.Mutex
	md                 goldmark.Markdown
	shortcodeProcessor *markdown.ShortcodeProcessor
	stats              *ParseStats
}

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

	slices.Sort(files)

	return p.parseFiles(rootDir, files)
}

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

	slices.Sort(mdFiles)

	return p.parseFiles(rootDir, mdFiles)
}

func (p *ContentParser) parseFiles(rootDir string, files []string) error {
	ctx := context.Background()
	pool := utils.NewWorkerPool(ctx)

	_ = pool.ProcessFiles(files, func(ctx context.Context, filePath string) error {
		relPath, _ := filepath.Rel(rootDir, filePath)
		dir := filepath.Dir(relPath)
		dir = filepath.ToSlash(dir)

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

func (p *ContentParser) GetMarkdownProcessor() goldmark.Markdown {
	return p.md
}

func (p *ContentParser) SetShortcodeTemplates(templates *template.Template) {
	syntaxTheme := p.Site.SyntaxTheme
	if syntaxTheme == "" {
		syntaxTheme = "dracula"
	}

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

func (p *ContentParser) detectLanguage(path, dir string, fm *config.FrontMatter) string {
	if fm != nil && fm.Lang != "" {
		return fm.Lang
	}

	basename := filepath.Base(path)
	if ext := filepath.Ext(basename); ext == ".md" {
		nameWithoutExt := strings.TrimSuffix(basename, ext)
		parts := strings.Split(nameWithoutExt, ".")
		if len(parts) >= 2 {
			potentialLang := parts[len(parts)-1]
			if p.Site.I18n != nil && p.Site.I18n.GetLanguage(potentialLang) != nil {
				return potentialLang
			}
		}
	}

	if dir != "." && dir != "" {
		parts := strings.Split(filepath.ToSlash(dir), "/")
		if len(parts) > 0 {
			firstDir := parts[0]
			if p.Site.I18n != nil && p.Site.I18n.GetLanguage(firstDir) != nil {
				return firstDir
			}
		}
	}

	if p.Site.I18n != nil {
		return p.Site.I18n.GetDefaultLanguage().Code
	}

	if p.Site.Lang != "" {
		return p.Site.Lang
	}

	return "en"
}

func (p *ContentParser) ResetStats() {
	p.stats.FilesSkipped.Store(0)
	p.stats.FilesParsed.Store(0)
	p.stats.TotalFiles.Store(0)
}

func (p *ContentParser) sortChildren() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, node := range p.ContentMap {
		p.sortNodeChildren(node)
	}
}

func (p *ContentParser) sortNodeChildren(node *content.Node) {
	if node == nil || len(node.Children) == 0 {
		return
	}

	slices.SortFunc(node.Children, func(a, b *content.Node) int {
		if a.Path < b.Path {
			return -1
		}
		if a.Path > b.Path {
			return 1
		}
		return 0
	})

	for _, child := range node.Children {
		p.sortNodeChildren(child)
	}
}

package parser

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/tduyng/gozzi/app"
	"github.com/tduyng/gozzi/app/concurrent"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/paginate"
	"github.com/yuin/goldmark"
	highlight "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/mermaid"
)

// ContentParser parses markdown content files and builds the content tree.
type ContentParser struct {
	Site       *config.Site
	ContentMap map[string]*content.Node
	Tags       map[string]*TagEntry
	mu         sync.Mutex
	md         goldmark.Markdown
}

// TagEntry tracks pages associated with a specific tag.
type TagEntry struct {
	Pages []*content.Node
	Count int
	Seen  map[string]struct{} // Track page paths
}

// NewParser creates a new ContentParser with the given site configuration.
func NewParser(cfg *config.Site) *ContentParser {
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
					highlight.WithStyle("dracula"),
				),
				&mermaid.Extender{},
				NewMathExtension(),
				NewTocExtension(),
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
	pool := concurrent.NewWorkerPool(ctx)

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

func (p *ContentParser) parseSection(path, dir string) error {
	mdContent, err := os.ReadFile(path)
	if err != nil {
		return app.WrapWithContext(app.ErrFileSystem, err, app.ErrorContext{
			Operation: "read_section_file",
			Component: "content_parser",
			Path:      path,
		})
	}

	frontMatter, contentPart, err := config.LoadFrontMatter(mdContent)
	if err != nil {
		return app.WrapWithContext(app.ErrContent, err, app.ErrorContext{
			Operation: "parse_section_frontmatter",
			Component: "content_parser",
			Path:      path,
		})
	}
	if frontMatter.Draft {
		return app.WrapWithContext(app.ErrContent, fmt.Errorf("section is marked as draft"), app.ErrorContext{
			Operation: "validate_section_draft_status",
			Component: "content_parser",
			Path:      path,
		})
	}

	pc := parser.NewContext()
	doc := p.md.Parser().Parse(text.NewReader(contentPart), parser.WithContext(pc))
	toc, _ := pc.Get(0).([]map[string]any)

	var htmlBuf bytes.Buffer
	if err := p.md.Renderer().Render(&htmlBuf, contentPart, doc); err != nil {
		return app.WrapWithContext(app.ErrContent, err, app.ErrorContext{
			Operation: "render_section_markdown",
			Component: "content_parser",
			Path:      path,
		})
	}
	sectionConfig := frontMatter.ToConfig()
	mergedConfig := config.MergeConfigs(p.Site.ToConfig(), sectionConfig, nil)
	slug := content.GenerateSlug(path, nil)
	mergedConfig["assets"] = filepath.Join(filepath.Dir(path), "img")
	mergedConfig["img_url"] = p.resolveImgURL(frontMatter, slug)

	p.mu.Lock()
	defer p.mu.Unlock()

	node := p.GetOrCreateSection(dir)
	node.Type = content.NodeTypeSection
	node.Config = mergedConfig
	node.Content = template.HTML(htmlBuf.String())
	node.Permalink = buildPermalink(slug)
	node.URL = buildURL(p.Site.BaseURL, slug)

	wordCount, readTime := calculateReadStats(string(contentPart))
	node.WordCount = wordCount
	node.ReadTime = readTime
	node.Path = strings.TrimPrefix(path, "content/")
	node.Toc = toc

	return nil
}

func (p *ContentParser) parsePage(path, dir string) error {
	mdContent, err := os.ReadFile(path)
	if err != nil {
		return app.WrapWithContext(app.ErrFileSystem, err, app.ErrorContext{
			Operation: "read_page_file",
			Component: "content_parser",
			Path:      path,
		})
	}

	pageConfig, contentPart, err := config.LoadFrontMatter(mdContent)
	if err != nil {
		return app.WrapWithContext(app.ErrContent, err, app.ErrorContext{
			Operation: "parse_page_frontmatter",
			Component: "content_parser",
			Path:      path,
		})
	}
	if pageConfig.Draft {
		return app.WrapWithContext(app.ErrContent, fmt.Errorf("page is marked as draft"), app.ErrorContext{
			Operation: "validate_page_draft_status",
			Component: "content_parser",
			Path:      path,
		})
	}

	pc := parser.NewContext()
	doc := p.md.Parser().Parse(text.NewReader(contentPart), parser.WithContext(pc))
	toc, _ := pc.Get(0).([]map[string]any)

	var htmlBuf bytes.Buffer
	if err := p.md.Renderer().Render(&htmlBuf, contentPart, doc); err != nil {
		return app.WrapWithContext(app.ErrContent, err, app.ErrorContext{
			Operation: "render_page_markdown",
			Component: "content_parser",
			Path:      path,
		})
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	var sectionConfig map[string]any
	if secNode, exists := p.ContentMap[dir]; exists {
		sectionConfig = secNode.Config
	}

	mergedConfig := config.MergeConfigs(p.Site.ToConfig(), sectionConfig, pageConfig.ToConfig())

	// For index.md files, the page should be created at the directory level,
	// with parent being the directory's parent, not the directory itself
	var parent *content.Node
	var pagePath string
	if filepath.Base(path) == "index.md" {
		// For blog/first-post/index.md: parent = blog, pagePath = blog/first-post
		parent = p.GetOrCreateSection(filepath.Dir(dir))
		pagePath = dir
	} else {
		// For blog/post.md: parent = blog, pagePath = blog/post
		parent = p.GetOrCreateSection(dir)
		pagePath = strings.TrimSuffix(path, "content/")
	}

	slug := content.GenerateSlug(pagePath, parent)
	permalink := buildPermalink(slug)

	wordCount, readTime := calculateReadStats(string(contentPart))
	mergedConfig["assets"] = filepath.Join(filepath.Dir(path), "img")
	mergedConfig["img_url"] = p.resolveImgURL(pageConfig, slug)

	pageNode := &content.Node{
		Path:      pagePath,
		Slug:      slug,
		Permalink: permalink,
		URL:       buildURL(p.Site.BaseURL, slug),
		Type:      content.NodeTypePage,
		Parent:    parent,
		Config:    mergedConfig,
		Content:   template.HTML(htmlBuf.String()),
		WordCount: wordCount,
		ReadTime:  readTime,
		Toc:       toc,
	}

	parent.Children = append(parent.Children, pageNode)

	// handle tags
	if len(pageConfig.Tags) > 0 {
		p.parseTags(pageConfig, pageNode)
	}

	return nil
}

func (p *ContentParser) resolveImgURL(fm *config.FrontMatter, slug string) string {
	img := ""
	if fm.Extra != nil {
		if val, ok := fm.Extra["img"]; ok {
			img = fmt.Sprintf("%v", val)
		}
	}

	if img == "" && p.Site.Extra != nil {
		if val, ok := p.Site.Extra["img"]; ok {
			img = fmt.Sprintf("%v", val)
		}
	}

	baseURL := strings.TrimSuffix(p.Site.BaseURL, "/")
	if strings.HasPrefix(img, "/") {
		return baseURL + img
	}

	return baseURL + filepath.Join("/", slug, img)
}

// GetOrCreateSection retrieves or creates a section node for the given directory.
func (p *ContentParser) GetOrCreateSection(dir string) *content.Node {
	if node, exists := p.ContentMap[dir]; exists {
		return node
	}

	var parent *content.Node
	var sectionSlug string

	if dir == "." { // Handle root section
		sectionSlug = ""
	} else {
		parentDir := filepath.Dir(dir)
		parent = p.GetOrCreateSection(parentDir)
		baseName := filepath.Base(dir)
		sectionSlug = content.GenerateSlug(baseName, nil)

		// Combine with parent slug
		if parent.Slug != "" {
			sectionSlug = parent.Slug + "/" + sectionSlug
		}
	}

	node := content.NewContentNode(dir, parent)
	node.Type = content.NodeTypeSection
	node.Slug = sectionSlug // Override generated slug

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	p.ContentMap[dir] = node
	return node
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

func (p *ContentParser) parseTags(pageConfig *config.FrontMatter, pageNode *content.Node) {
	uniqueTags := make(map[string]bool)
	for _, rawTag := range pageConfig.Tags {
		tag := strings.ToLower(strings.TrimSpace(rawTag))
		if tag == "" {
			continue
		}

		if _, exists := uniqueTags[tag]; exists {
			continue
		}
		uniqueTags[tag] = true

		// Get or create tag entry
		entry, exists := p.Tags[tag]
		if !exists {
			entry = &TagEntry{
				Seen: make(map[string]struct{}),
			}
			p.Tags[tag] = entry
		}

		// Add page if not already present
		if _, exists := entry.Seen[pageNode.Path]; !exists {
			entry.Pages = append(entry.Pages, pageNode)
			entry.Seen[pageNode.Path] = struct{}{}
			entry.Count = len(entry.Pages)
		}
	}
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

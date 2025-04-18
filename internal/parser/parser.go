package parser

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"

	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/content"
	"github.com/tduyng/gozzi/internal/paginate"
	"github.com/yuin/goldmark"
	highlight "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/mermaid"
)

type ContentParser struct {
	Site       *config.Site
	ContentMap map[string]*content.Node
	Tags       map[string]*TagEntry
	mu         sync.Mutex
	md         goldmark.Markdown
	workerPool chan struct{}
}

type TagEntry struct {
	Pages []*content.Node
	Count int
	Seen  map[string]struct{} // Track page paths
}

func NewParser(cfg *config.Site) *ContentParser {
	return &ContentParser{
		Site:       cfg,
		ContentMap: make(map[string]*content.Node),
		Tags:       make(map[string]*TagEntry),
		workerPool: make(chan struct{}, runtime.NumCPU()*2),
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

func (p *ContentParser) Parse(rootDir string) error {
	p.mu.Lock()
	p.ContentMap = make(map[string]*content.Node) // Reset ContentMap
	p.mu.Unlock()

	var wg sync.WaitGroup
	fileChan := make(chan string, 100)

	go func() {
		filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if filepath.Base(path) == "_index.md" || filepath.Ext(path) == ".md" {
				fileChan <- path
			}
			return nil
		})
		close(fileChan)
	}()

	for path := range fileChan {
		wg.Add(1)
		p.workerPool <- struct{}{}

		go func(path string) {
			defer func() {
				<-p.workerPool
				wg.Done()
			}()

			relPath, _ := filepath.Rel(rootDir, path)
			dir := filepath.Dir(relPath)

			switch {
			case filepath.Base(path) == "_index.md":
				p.parseSection(path, dir)
			case filepath.Ext(path) == ".md":
				p.parsePage(path, dir)
			}
		}(path)
	}

	wg.Wait()

	paginator := paginate.New(p.ContentMap)
	paginator.BuildLinks()
	return nil
}

func (p *ContentParser) parseSection(path, dir string) error {
	mdContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	frontMatter, contentPart, err := config.LoadFrontMatter(mdContent)
	if err != nil || frontMatter.Draft {
		return err
	}

	pc := parser.NewContext()
	doc := p.md.Parser().Parse(text.NewReader(contentPart), parser.WithContext(pc))
	toc, _ := pc.Get(0).([]map[string]any)

	var htmlBuf bytes.Buffer
	if err := p.md.Renderer().Render(&htmlBuf, contentPart, doc); err != nil {
		return fmt.Errorf("markdown rendering failed: %w", err)
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
	node.URL = buildURL(p.Site.BaseURL, node.Permalink)

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
		return err
	}

	pageConfig, contentPart, err := config.LoadFrontMatter(mdContent)
	if err != nil || pageConfig.Draft {
		return err
	}

	pc := parser.NewContext()
	doc := p.md.Parser().Parse(text.NewReader(contentPart), parser.WithContext(pc))
	toc, _ := pc.Get(0).([]map[string]any)

	var htmlBuf bytes.Buffer
	if err := p.md.Renderer().Render(&htmlBuf, contentPart, doc); err != nil {
		return fmt.Errorf("markdown rendering failed: %w", err)
	}

	var sectionConfig map[string]any
	if secNode, exists := p.ContentMap[dir]; exists {
		sectionConfig = secNode.Config
	}

	mergedConfig := config.MergeConfigs(p.Site.ToConfig(), sectionConfig, pageConfig.ToConfig())

	p.mu.Lock()
	defer p.mu.Unlock()

	parent := p.GetOrCreateSection(filepath.Dir(dir))
	slug := content.GenerateSlug(path, parent)
	permalink := buildPermalink(slug)

	wordCount, readTime := calculateReadStats(string(contentPart))
	mergedConfig["assets"] = filepath.Join(filepath.Dir(path), "img")
	mergedConfig["img_url"] = p.resolveImgURL(pageConfig, slug)

	pageNode := &content.Node{
		Path:      strings.TrimSuffix(path, "content/"),
		Slug:      slug,
		Permalink: permalink,
		URL:       buildURL(p.Site.BaseURL, permalink),
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
	if permalink == "" {
		return "/"
	}
	return "/" + permalink + "/"
}

func buildURL(baseURL, permalink string) string {
	return baseURL + permalink
}

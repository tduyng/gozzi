package parser

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/tduyng/gozzi/internal/config"
	"github.com/tduyng/gozzi/internal/content"
	"github.com/yuin/goldmark"
	highlight "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

type ContentParser struct {
	Site       *config.Site
	ContentMap map[string]*content.Node
	mu         sync.Mutex
	md         goldmark.Markdown
}

func NewParser(cfg *config.Site) *ContentParser {
	return &ContentParser{
		Site:       cfg,
		ContentMap: make(map[string]*content.Node),
		md: goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Footnote,
				highlight.NewHighlighting(
					highlight.WithGuessLanguage(true),
					highlight.WithStyle("dracula"),
				),
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

	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(rootDir, path)
		dir := filepath.Dir(relPath)

		switch {
		case filepath.Base(path) == "_index.md":
			return p.parseSection(path, dir)
		case filepath.Ext(path) == ".md":
			return p.parsePage(path, dir)
		}

		return nil
	})
}

func (p *ContentParser) parseSection(path, dir string) error {
	mkdown, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	frontMatter, contentPart, err := config.LoadFrontMatter(mkdown)
	if err != nil || frontMatter.Draft {
		return err
	}

	var buf bytes.Buffer
	if err := p.md.Convert(contentPart, &buf); err != nil {
		return fmt.Errorf("markdown conversion failed: %w", err)
	}
	sectionConfig := frontMatter.ToConfig()
	mergedConfig := config.MergeConfigs(p.Site.ToConfig(), sectionConfig, nil)
	mergedConfig["assets"] = filepath.Join(filepath.Dir(path), "img")
	mergedConfig["img"] = p.resolveImgURL(frontMatter, path)

	p.mu.Lock()
	defer p.mu.Unlock()

	node := p.GetOrCreateSection(dir)
	node.Type = content.NodeTypeSection
	node.Config = sectionConfig
	node.Content = template.HTML(buf.String())

	var pages []*content.Node
	for _, child := range node.Children {
		if child.Type == content.NodeTypePage {
			pages = append(pages, child)
		}
	}

	// Sort pages by date (newest first)
	sort.SliceStable(pages, func(i, j int) bool {
		dateI := pages[i].Config["date"].(time.Time)
		dateJ := pages[j].Config["date"].(time.Time)
		return dateI.After(dateJ)
	})

	// Set pagination links
	for i := range pages {
		if i > 0 {
			pages[i].Higher = pages[i-1] // Older post
		}
		if i < len(pages)-1 {
			pages[i].Lower = pages[i+1] // Newer post
		}
		pages[i].Parent = node
	}

	return nil
}

func (p *ContentParser) parsePage(path, dir string) error {
	mkdownContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	pageConfig, contentPart, err := config.LoadFrontMatter(mkdownContent)
	if err != nil || pageConfig.Draft {
		return err
	}

	var buf bytes.Buffer
	if err := p.md.Convert(contentPart, &buf); err != nil {
		return fmt.Errorf("markdown conversion failed: %w", err)
	}

	var sectionConfig map[string]any
	if secNode, exists := p.ContentMap[dir]; exists {
		sectionConfig = secNode.Config
	}

	mergedConfig := config.MergeConfigs(p.Site.ToConfig(), sectionConfig, pageConfig.ToConfig())
	mergedConfig["assets"] = filepath.Join(filepath.Dir(path), "img")
	mergedConfig["img"] = p.resolveImgURL(pageConfig, path)

	p.mu.Lock()
	defer p.mu.Unlock()

	parent := p.GetOrCreateSection(filepath.Dir(dir))
	slug := content.GenerateSlug(path, parent)

	pageNode := &content.Node{
		Path:    path,
		Slug:    slug,
		Type:    content.NodeTypePage,
		Parent:  parent,
		Config:  mergedConfig,
		Content: template.HTML(buf.String()),
	}

	parent.Children = append(parent.Children, pageNode)

	return nil
}

func (p *ContentParser) resolveImgURL(fm *config.FrontMatter, path string) string {
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

	mdDir := filepath.Dir(path)
	return baseURL + filepath.Join("/", mdDir, img)
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

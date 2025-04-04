package parser

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"

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
					highlight.WithStyle("github"),
				),
			),
			goldmark.WithRendererOptions(
				html.WithUnsafe(),
			),
		),
	}
}

func (p *ContentParser) Parse(rootDir string) error {
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

	sectionConfig, err := config.LoadFrontMatter(mkdown)
	if err != nil {
		return err
	}
	if sectionConfig.Draft {
		return nil
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	node := p.GetOrCreateSection(dir)
	node.Type = content.NodeTypeSection
	node.Config = config.MergeConfigs(p.Site.ToConfig(), sectionConfig.ToConfig(), nil)

	return nil
}

func (p *ContentParser) parsePage(path, dir string) error {
	mkdownContent, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	parts := bytes.SplitN(mkdownContent, []byte("+++"), 3)
	var pageConfig *config.FrontMatter
	var contentPart []byte
	if len(parts) < 3 {
		pageConfig = &config.FrontMatter{}
		contentPart = mkdownContent
	} else {
		pageConfig, err = config.LoadFrontMatter(mkdownContent)
		if err != nil {
			return err
		}
		contentPart = parts[2]
	}
	if pageConfig.Draft {
		return nil
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

	pageNode := content.NewContentNode(path, nil)
	pageNode.Type = content.NodeTypePage
	pageNode.Config = mergedConfig
	pageNode.Content = template.HTML(buf.String())

	sectionDir := filepath.Dir(dir)
	parent := p.GetOrCreateSection(sectionDir)
	pageNode.Parent = parent
	parent.Children = append(parent.Children, pageNode)
	parent.Section.Children = append(parent.Section.Children, pageNode)

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
		if node.Section == nil {
			node.Section = &content.Section{}
		}
		return node
	}

	var parent *content.Node
	if dir != "." {
		parentDir := filepath.Dir(dir)
		parent = p.GetOrCreateSection(parentDir)
	}

	node := content.NewContentNode(dir, parent)
	node.Type = content.NodeTypeSection
	node.Section = &content.Section{}

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	p.ContentMap[dir] = node
	return node
}

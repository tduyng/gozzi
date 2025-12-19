// Package parser provides page parsing for individual markdown files.
package parser

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/utils"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func (p *ContentParser) parsePage(path, dir string) error {
	mdContent, err := os.ReadFile(path)
	if err != nil {
		return utils.WrapWithContext(utils.ErrFileSystem, err, utils.ErrorContext{
			Operation: "read_page_file",
			Component: "content_parser",
			Path:      path,
		})
	}

	// Incremental parsing: Check if content has actually changed
	p.stats.TotalFiles.Add(1)
	if !p.hashCache.HasChanged(path, mdContent) {
		p.stats.FilesSkipped.Add(1)
		return nil // Content unchanged, skip parsing
	}
	p.stats.FilesParsed.Add(1)

	pageConfig, contentPart, err := config.LoadFrontMatter(mdContent)
	if err != nil {
		return utils.WrapWithContext(utils.ErrContent, err, utils.ErrorContext{
			Operation: "parse_page_frontmatter",
			Component: "content_parser",
			Path:      path,
		})
	}

	if pageConfig.Draft && !p.Site.BuildDrafts {
		return nil
	}

	pc := parser.NewContext()
	doc := p.md.Parser().Parse(text.NewReader(contentPart), parser.WithContext(pc))
	toc, _ := pc.Get(0).([]map[string]any)

	var htmlBuf bytes.Buffer
	if err := p.md.Renderer().Render(&htmlBuf, contentPart, doc); err != nil {
		return utils.WrapWithContext(utils.ErrContent, err, utils.ErrorContext{
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

	var parent *content.Node
	var pagePath string
	if filepath.Base(path) == "index.md" {
		parent = p.GetOrCreateSection(filepath.Dir(dir))
		pagePath = dir
	} else {
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

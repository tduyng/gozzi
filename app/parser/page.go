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
	"github.com/tduyng/gozzi/app/summary"
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

	p.stats.TotalFiles.Add(1)
	if !p.hashCache.HasChanged(path, mdContent) {
		p.stats.FilesSkipped.Add(1)
		return nil
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

	// Process shortcodes if processor is available
	if p.shortcodeProcessor != nil {
		contentPart, err = p.shortcodeProcessor.Process(contentPart)
		if err != nil {
			return utils.WrapWithContext(utils.ErrContent, err, utils.ErrorContext{
				Operation: "process_shortcodes",
				Component: "content_parser",
				Path:      path,
			})
		}
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

	// Detect and set language for this page
	lang := p.detectLanguage(path, dir, pageConfig)
	mergedConfig["lang"] = lang

	var parent *content.Node
	var pagePath string
	if filepath.Base(path) == "index.md" {
		parent = p.GetOrCreateSection(filepath.Dir(dir))
		pagePath = dir
	} else {
		parent = p.GetOrCreateSection(dir)
		// Handle both absolute and relative paths
		if filepath.IsAbs(path) {
			// For absolute paths, construct the path from content/ + dir + filename
			// to match the format used during initial build
			pagePath = filepath.Join("content", dir, filepath.Base(path))
		} else {
			pagePath = strings.TrimSuffix(path, "content/")
		}
	}

	slug := content.GenerateSlug(pagePath, parent)
	permalink := buildPermalink(slug)

	wordCount, readTime := calculateReadStats(string(contentPart))
	mergedConfig["assets"] = filepath.Join(filepath.Dir(path), "img")
	mergedConfig["img_url"] = p.resolveImgURL(pageConfig, slug)

	// Generate summary
	summaryGen := summary.New()
	if p.Site.SummaryLength > 0 {
		summaryGen.SentenceCount = p.Site.SummaryLength
	}
	summaryText := summaryGen.Generate(pageConfig.Description, template.HTML(htmlBuf.String()))

	// Extract aliases from frontmatter
	aliases := pageConfig.Aliases
	if aliases == nil {
		aliases = []string{}
	}

	pageNode := &content.Node{
		Path:      pagePath,
		Slug:      slug,
		Permalink: permalink,
		URL:       buildURL(p.Site.BaseURL, slug),
		Type:      content.NodeTypePage,
		Parent:    parent,
		Config:    mergedConfig,
		Content:   template.HTML(htmlBuf.String()),
		Summary:   template.HTML(summaryText),
		WordCount: wordCount,
		ReadTime:  readTime,
		Toc:       toc,
		Aliases:   aliases,
	}

	found := false
	var targetNode *content.Node
	for _, child := range parent.Children {
		if child.Path == pageNode.Path {
			// This ensures any pointers to this node (from RebuildAnalyzer, Builder, etc.)
			// see the updated data
			p.RemovePageFromAllTaxonomies(child)

			// Update all fields of the existing node
			child.Slug = pageNode.Slug
			child.Permalink = pageNode.Permalink
			child.URL = pageNode.URL
			child.Config = pageNode.Config
			child.Content = pageNode.Content
			child.Summary = pageNode.Summary
			child.WordCount = pageNode.WordCount
			child.ReadTime = pageNode.ReadTime
			child.Toc = pageNode.Toc
			child.Aliases = pageNode.Aliases
			// Note: Path, Type, and Parent should not change

			targetNode = child
			found = true
			break
		}
	}
	if !found {
		parent.Children = append(parent.Children, pageNode)
		targetNode = pageNode
	}

	// Parse all taxonomies (tags, categories, series, custom)
	p.ParseTaxonomies(pageConfig, targetNode)

	// Maintain backwards compatibility with Tags field
	if len(pageConfig.Tags) > 0 {
		p.parseTags(pageConfig, targetNode)
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

	return baseURL + filepath.ToSlash(filepath.Join("/", slug, img))
}

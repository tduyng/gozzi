// Package parser provides section parsing for _index.md files.
package parser

import (
	"bytes"
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

func (p *ContentParser) parseSection(path, dir string) error {
	mdContent, err := os.ReadFile(path)
	if err != nil {
		return utils.WrapWithContext(utils.ErrFileSystem, err, utils.ErrorContext{
			Operation: "read_section_file",
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

	frontMatter, contentPart, err := config.LoadFrontMatter(mdContent)
	if err != nil {
		return utils.WrapWithContext(utils.ErrContent, err, utils.ErrorContext{
			Operation: "parse_section_frontmatter",
			Component: "content_parser",
			Path:      path,
		})
	}

	if frontMatter.Draft && !p.Site.BuildDrafts {
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
			Operation: "render_section_markdown",
			Component: "content_parser",
			Path:      path,
		})
	}
	sectionConfig := frontMatter.ToConfig()
	mergedConfig := config.MergeConfigs(p.Site.ToConfig(), sectionConfig, nil)

	// Detect and set language for this section
	lang := p.detectLanguage(path, dir, frontMatter)
	mergedConfig["lang"] = lang

	p.mu.Lock()
	defer p.mu.Unlock()

	existingNode := p.GetOrCreateSection(dir)

	// Use the existing node's slug (which is correctly calculated from the directory)
	slug := existingNode.Slug
	mergedConfig["assets"] = filepath.Join(filepath.Dir(path), "img")
	mergedConfig["img_url"] = p.resolveImgURL(frontMatter, slug)

	wordCount, readTime := calculateReadStats(string(contentPart))

	// Extract aliases from frontmatter
	aliases := frontMatter.Aliases
	if aliases == nil {
		aliases = []string{}
	}

	newNode := &content.Node{
		Type:      content.NodeTypeSection,
		Config:    mergedConfig,
		Content:   template.HTML(htmlBuf.String()),
		Permalink: buildPermalink(slug),
		URL:       buildURL(p.Site.BaseURL, slug),
		WordCount: wordCount,
		ReadTime:  readTime,
		Path:      strings.TrimPrefix(path, "content/"),
		Toc:       toc,
		Slug:      existingNode.Slug,
		Parent:    existingNode.Parent,
		Children:  existingNode.Children,
		Higher:    existingNode.Higher,
		Lower:     existingNode.Lower,
		Aliases:   aliases,
	}

	p.ContentMap[dir] = newNode

	if newNode.Parent != nil {
		for i, child := range newNode.Parent.Children {
			if child == existingNode {
				newNode.Parent.Children[i] = newNode
				break
			}
		}
	}

	for _, child := range newNode.Children {
		child.Parent = newNode
	}

	return nil
}

// GetOrCreateSection retrieves or creates a section node.
func (p *ContentParser) GetOrCreateSection(dir string) *content.Node {
	if node, exists := p.ContentMap[dir]; exists {
		return node
	}

	var parent *content.Node
	var sectionSlug string

	if dir == "." {
		sectionSlug = ""
	} else {
		parentDir := filepath.Dir(dir)
		parent = p.GetOrCreateSection(parentDir)
		baseName := filepath.Base(dir)
		sectionSlug = content.GenerateSlug(baseName, nil)

		if parent.Slug != "" {
			sectionSlug = parent.Slug + "/" + sectionSlug
		}
	}

	node := content.NewContentNode(dir, parent)
	node.Type = content.NodeTypeSection
	node.Slug = sectionSlug
	node.Permalink = buildPermalink(sectionSlug)
	node.URL = buildURL(p.Site.BaseURL, sectionSlug)

	if parent != nil {
		parent.Children = append(parent.Children, node)
	}

	p.ContentMap[dir] = node
	return node
}

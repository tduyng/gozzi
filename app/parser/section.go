// Package parser provides section parsing logic for _index.md files.
// Handles section frontmatter, content rendering, and section node creation.
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

	frontMatter, contentPart, err := config.LoadFrontMatter(mdContent)
	if err != nil {
		return utils.WrapWithContext(utils.ErrContent, err, utils.ErrorContext{
			Operation: "parse_section_frontmatter",
			Component: "content_parser",
			Path:      path,
		})
	}
	// Skip draft sections unless BuildDrafts is enabled
	if frontMatter.Draft && !p.Site.BuildDrafts {
		return nil // Silently skip draft sections
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

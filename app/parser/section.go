package parser

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/utils"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

func (p *ContentParser) parseSection(path, dir string) error {
	p.stats.TotalFiles.Add(1)
	p.stats.FilesParsed.Add(1)

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

	if frontMatter == nil {
		frontMatter = &config.FrontMatter{}
	}

	if frontMatter.Draft && !p.Site.BuildDrafts {
		return nil
	}

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

	p.mu.Lock()
	defer p.mu.Unlock()

	existingNode := p.getOrCreateSection(dir)
	slug := existingNode.Slug
	permalink := buildPermalink(slug)

	htmlContent := htmlBuf.String()
	htmlContent = rewriteRelativePaths(htmlContent, slug)

	sectionConfig := frontMatter.ToConfig()
	mergedConfig := config.MergeConfigs(p.Site.ToConfig(), sectionConfig, nil)

	lang := p.detectLanguage(path, dir, frontMatter)
	mergedConfig["lang"] = lang
	mergedConfig["assets"] = filepath.Join(filepath.Dir(path), "img")
	mergedConfig["img_url"] = p.resolveImgURL(frontMatter, slug)

	wordCount, readTime := calculateReadStats(string(contentPart))

	aliases := frontMatter.Aliases

	newNode := &content.Node{
		Type:      content.NodeTypeSection,
		Config:    mergedConfig,
		Content:   template.HTML(htmlContent),
		Permalink: permalink,
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

func (p *ContentParser) GetOrCreateSection(dir string) *content.Node {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.getOrCreateSection(dir)
}

func (p *ContentParser) getOrCreateSection(dir string) *content.Node {
	if node, exists := p.ContentMap[dir]; exists {
		return node
	}

	var parent *content.Node
	var sectionSlug string

	if dir == "." {
		sectionSlug = ""
	} else {
		parentDir := path.Dir(dir)
		parent = p.getOrCreateSection(parentDir)
		baseName := path.Base(dir)
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

var (
	srcRelRegex   = regexp.MustCompile(`(src|href)=["']([^/][^"':]*)["']`)
	protoRelRegex = regexp.MustCompile(`^(http|https|mailto|tel):`)
)

func rewriteRelativePaths(html string, sectionDir string) string {
	sectionDir = strings.TrimPrefix(sectionDir, "./")
	if sectionDir == "." || sectionDir == "" {
		return html
	}

	// Ensure leading slash and trailing slash for base
	basePath := "/" + strings.Trim(sectionDir, "/") + "/"

	return srcRelRegex.ReplaceAllStringFunc(html, func(match string) string {
		submatch := srcRelRegex.FindStringSubmatch(match)
		if len(submatch) < 3 {
			return match
		}

		attr := submatch[1]
		val := submatch[2]

		// Skip if it has a protocol or is a fragment/absolute path
		if protoRelRegex.MatchString(val) || strings.HasPrefix(val, "/") || strings.HasPrefix(val, "#") {
			return match
		}

		return fmt.Sprintf(`%s="%s"`, attr, path.Join(basePath, val))
	})
}

// Site-specific template functions that require generator context for asset paths and content access.
package funcs

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/tduyng/gozzi/app/content"
	"github.com/yuin/goldmark"
)

// SiteContext provides access to site-wide data needed by template functions.
type SiteContext struct {
	BaseURL    string
	ContentMap map[string]*content.Node
	Markdown   goldmark.Markdown
}

// SiteFuncs holds template functions that require site context.
type SiteFuncs struct {
	ctx *SiteContext
}

// NewSiteFuncs creates a new SiteFuncs with the given context.
func NewSiteFuncs(ctx *SiteContext) *SiteFuncs {
	return &SiteFuncs{ctx: ctx}
}

// AssetPath generates the full URL for an asset path.
func (sf *SiteFuncs) AssetPath(relPath string, context any) (string, error) {
	baseURL := strings.TrimSuffix(sf.ctx.BaseURL, "/")

	// Absolute URLs pass through
	if strings.HasPrefix(relPath, "http://") || strings.HasPrefix(relPath, "https://") {
		return relPath, nil
	}

	// Root-relative paths
	if strings.HasPrefix(relPath, "/") {
		return baseURL + relPath, nil
	}

	// Context-relative paths
	if ctx, ok := context.(*content.Node); ok {
		pagePath := ctx.Slug
		if !strings.HasSuffix(pagePath, "/") {
			pagePath += "/"
		}
		return baseURL + "/" + pagePath + relPath, nil
	}

	// Default to root-relative
	return baseURL + "/" + relPath, nil
}

// GetSection retrieves a section node by path.
func (sf *SiteFuncs) GetSection(path string) (*content.Node, error) {
	// Normalize path to directory
	if !strings.Contains(path, "_index.md") {
		if strings.HasSuffix(path, "/") {
			path = path + "_index.md"
		} else {
			path = path + "/_index.md"
		}
	}

	sectionDir := strings.TrimSuffix(path, "/_index.md")
	if sectionDir == "" {
		sectionDir = "."
	}

	if node, ok := sf.ctx.ContentMap[sectionDir]; ok {
		return node, nil
	}

	return nil, fmt.Errorf("section not found: %s", path)
}

// RenderMarkdown renders markdown content to HTML.
func (sf *SiteFuncs) RenderMarkdown(input string) (template.HTML, error) {
	if sf.ctx.Markdown == nil {
		return "", fmt.Errorf("markdown processor not initialized")
	}

	var buf bytes.Buffer
	if err := sf.ctx.Markdown.Convert([]byte(input), &buf); err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}
	return template.HTML(buf.String()), nil
}

// LoadData loads file content as HTML (unescaped).
func LoadData(path string) (template.HTML, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load file %s: %w", path, err)
	}
	return template.HTML(string(content)), nil
}

// LoadAttribute loads file content as HTML-escaped attribute value.
func LoadAttribute(path string) (template.HTML, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load file %s: %w", path, err)
	}
	escaped := template.HTMLEscapeString(string(content))
	return template.HTML(escaped), nil
}

// SafeHTML marks a string as safe HTML (bypasses auto-escaping).
func SafeHTML(s string) template.HTML {
	return template.HTML(s)
}

// MacroRenderer holds methods for rendering template macros.
type MacroRenderer struct {
	templates *template.Template
}

// NewMacroRenderer creates a new macro renderer.
func NewMacroRenderer(templates *template.Template) *MacroRenderer {
	return &MacroRenderer{templates: templates}
}

// RenderPagination renders the pagination macro.
func (mr *MacroRenderer) RenderPagination(siteConfig map[string]any) func(data map[string]any) (template.HTML, error) {
	return func(data map[string]any) (template.HTML, error) {
		var buf bytes.Buffer
		tpl := mr.templates.Lookup("macros/pagination.html")
		if tpl == nil {
			return "", fmt.Errorf("pagination template not found")
		}

		err := tpl.Execute(&buf, map[string]any{
			"Page": data["Page"],
			"Site": map[string]any{
				"Config": siteConfig,
			},
		})
		if err != nil {
			return "", fmt.Errorf("failed to render pagination: %w", err)
		}
		return template.HTML(buf.String()), nil
	}
}

// Package funcs provides site-specific template functions requiring generator context.
package funcs

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/i18n"
	"github.com/yuin/goldmark"
)

// SiteContext provides access to site-wide data needed by template functions.
type SiteContext struct {
	BaseURL    string
	ContentMap map[string]*content.Node
	Markdown   goldmark.Markdown
	I18n       *i18n.I18n
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

// Translate returns a translation for the given key in the specified language.
// Usage in templates:
//
//	{{ i18n "section.contact" }}  - uses current page language from context
//	{{ i18n "hello" "fr" }}       - explicit language
func (sf *SiteFuncs) Translate(key string, langOrContext ...any) (string, error) {
	if sf.ctx.I18n == nil {
		return key, nil // If i18n not configured, return the key itself
	}

	var langCode string

	// Try to extract language from arguments
	if len(langOrContext) > 0 {
		// Check if first argument is a string (explicit language code)
		if lang, ok := langOrContext[0].(string); ok {
			langCode = lang
		} else if node, ok := langOrContext[0].(*content.Node); ok {
			// Try to get language from node config
			if lang, exists := node.Config["lang"]; exists {
				if langStr, ok := lang.(string); ok {
					langCode = langStr
				}
			}
		} else if configMap, ok := langOrContext[0].(map[string]any); ok {
			// Try to get language from config map
			if lang, exists := configMap["lang"]; exists {
				if langStr, ok := lang.(string); ok {
					langCode = langStr
				}
			}
		}
	}

	// If no language specified, use default
	if langCode == "" {
		langCode = sf.ctx.I18n.GetDefaultLanguage().Code
	}

	result, err := sf.ctx.I18n.Translate(langCode, key)
	if err != nil {
		// If translation fails, return the key itself (graceful fallback)
		return key, nil
	}
	return result, nil
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

// LangURL generates a URL for a specific language version of a page.
// Usage in templates:
//
//	{{ langURL "fr" .Page }}           - Get current page in French
//	{{ langURL "en" .Page.Permalink }} - Get specific path in English
func (sf *SiteFuncs) LangURL(langCode string, pageOrPath any) string {
	if sf.ctx.I18n == nil || !sf.ctx.I18n.IsEnabled() {
		// If i18n not enabled, return the original URL
		if node, ok := pageOrPath.(*content.Node); ok {
			return node.Permalink
		}
		if path, ok := pageOrPath.(string); ok {
			return path
		}
		return "/"
	}

	// Validate language code
	if sf.ctx.I18n.GetLanguage(langCode) == nil {
		// Invalid language code, return original
		if node, ok := pageOrPath.(*content.Node); ok {
			return node.Permalink
		}
		if path, ok := pageOrPath.(string); ok {
			return path
		}
		return "/"
	}

	var currentLang string
	var pathWithoutLang string

	// Extract current language and path
	if node, ok := pageOrPath.(*content.Node); ok {
		if lang, exists := node.Config["lang"]; exists {
			if langStr, ok := lang.(string); ok {
				currentLang = langStr
			}
		}
		pathWithoutLang = node.Permalink
	} else if pageMap, ok := pageOrPath.(map[string]any); ok {
		// Handle Page as a map (template data format)
		if config, exists := pageMap["Config"]; exists {
			if configMap, ok := config.(map[string]any); ok {
				if lang, exists := configMap["lang"]; exists {
					if langStr, ok := lang.(string); ok {
						currentLang = langStr
					}
				}
			}
		}
		if permalink, exists := pageMap["Permalink"]; exists {
			if permalinkStr, ok := permalink.(string); ok {
				pathWithoutLang = permalinkStr
			}
		}
	} else if path, ok := pageOrPath.(string); ok {
		pathWithoutLang = path
		// Try to detect language from path
		if strings.HasPrefix(path, "/") {
			path = strings.TrimPrefix(path, "/")
		}
		parts := strings.Split(path, "/")
		if len(parts) > 0 && sf.ctx.I18n.GetLanguage(parts[0]) != nil {
			currentLang = parts[0]
		}
	} else {
		return "/"
	}

	// Remove current language prefix from path
	if currentLang != "" && strings.HasPrefix(pathWithoutLang, "/"+currentLang+"/") {
		pathWithoutLang = strings.TrimPrefix(pathWithoutLang, "/"+currentLang)
	} else if currentLang != "" && pathWithoutLang == "/"+currentLang {
		pathWithoutLang = "/"
	}

	// Add new language prefix
	if pathWithoutLang == "/" || pathWithoutLang == "" {
		return "/" + langCode + "/"
	}

	// Ensure path starts with /
	if !strings.HasPrefix(pathWithoutLang, "/") {
		pathWithoutLang = "/" + pathWithoutLang
	}

	return "/" + langCode + pathWithoutLang
}

// CurrentLang returns the current language code from context.
// Usage in templates: {{ currentLang .Page }}
func (sf *SiteFuncs) CurrentLang(context any) string {
	if node, ok := context.(*content.Node); ok {
		if lang, exists := node.Config["lang"]; exists {
			if langStr, ok := lang.(string); ok {
				return langStr
			}
		}
	} else if configMap, ok := context.(map[string]any); ok {
		if lang, exists := configMap["lang"]; exists {
			if langStr, ok := lang.(string); ok {
				return langStr
			}
		}
	}

	// Fallback to default language
	if sf.ctx.I18n != nil {
		return sf.ctx.I18n.GetDefaultLanguage().Code
	}

	return "en"
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

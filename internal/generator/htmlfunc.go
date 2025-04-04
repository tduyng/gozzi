package generator

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/tduyng/gozzi/internal/content"
)

func (g *Generator) CreateFuncMap() template.FuncMap {
	return template.FuncMap{
		"asset":       g.assetPath,
		"default":     defaultValue,
		"format_date": formatDate,
		"get_section": g.getSection,
		"limit":       limit,
		"load":        loadDataToHTML,
		"priority":    priority,
		"reverse":     reverse,
		"safe":        safeHTML,
		"urlize":      urlize,
	}
}

func urlize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = regexp.MustCompile(`[^a-z0-9-]`).ReplaceAllString(s, "")
	s = regexp.MustCompile(`-+`).ReplaceAllString(s, "-")

	return strings.Trim(s, "-")
}

func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

func loadDataToHTML(path string) template.HTML {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file %s: %v", path, err)
		return ""
	}
	return template.HTML(string(content))
}

func (g *Generator) assetPath(relPath string, context any) string {
	baseURL := strings.TrimSuffix(g.site.BaseURL, "/")

	if strings.HasPrefix(relPath, "http://") || strings.HasPrefix(relPath, "https://") {
		return relPath
	}

	if strings.HasPrefix(relPath, "/") {
		return baseURL + relPath
	}
	if ctx, ok := context.(*content.Node); ok {
		pagePath := ctx.Slug
		if !strings.HasSuffix(pagePath, "/") {
			pagePath += "/"
		}
		return baseURL + pagePath + relPath
	}
	return baseURL + "/" + relPath
}

func priority(vals ...any) string {
	for _, v := range vals {
		if v == nil {
			continue
		}
		if s, ok := v.(string); ok && s != "" {
			return s
		}
		s := fmt.Sprint(v)
		if s != "" && s != "<nil>" {
			return s
		}
	}
	return ""
}

func (g *Generator) getSection(path string) *content.Node {
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
	if node, ok := g.parser.ContentMap[sectionDir]; ok {
		return node
	}
	return g.parser.GetOrCreateSection(sectionDir)
}

func formatDate(t time.Time, format ...string) string {
	if len(format) > 0 && format[0] != "" {
		return t.Format(format[0])
	}
	return t.Format("2006-01-02")
}

func limit(max any, items []*content.Node) []*content.Node {
	var m int
	switch v := max.(type) {
	case int:
		m = v
	case int64:
		m = int(v)
	default:
		m = 0
	}
	if m > len(items) {
		m = len(items)
	}
	return items[:m]
}

func reverse(items []*content.Node) []*content.Node {
	reversed := make([]*content.Node, len(items))
	for i := range items {
		reversed[len(items)-1-i] = items[i]
	}
	return reversed
}

func defaultValue(val any, def string) string {
	if val == nil {
		return def
	}
	if s, ok := val.(string); ok {
		if s == "" {
			return def
		}
		return s
	}
	s := fmt.Sprint(val)
	if s == "" || s == "<nil>" {
		return def
	}
	return s
}

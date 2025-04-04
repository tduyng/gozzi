package generator

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/tduyng/gozzi/internal/content"
)

func (g *Generator) CreateFuncMap() template.FuncMap {
	return template.FuncMap{
		"urlize":   urlize,
		"safe":     safeHTML,
		"load":     loadDataToHTML,
		"asset":    g.assetPath,
		"priority": priority,
		"section":  g.getSection,
		"get":      getConfig,
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

func getConfig(ctx map[string]any, key string) any {
	if val, ok := ctx[key]; ok {
		return val
	}
	return nil
}

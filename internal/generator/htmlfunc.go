package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tduyng/gozzi/internal/content"
)

func (g *Generator) CreateFuncMap() template.FuncMap {
	return template.FuncMap{
		"asset":       g.assetPath,
		"default":     defaultValue,
		"date":        formatDate,
		"get_section": g.getSection,
		"limit":       limit,
		"load":        loadDataToHTML,
		"priority":    priority,
		"reverse":     reverse,
		"safe":        safeHTML,
		"urlize":      urlize,
		"group_by":    g.groupBy,
		"pagination":  g.renderPagination,
		"start_with":  startWith,
	}
}

type Group struct {
	Key   string
	Items []*content.Node
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

func formatDate(t time.Time, layout ...string) string {
	if t.IsZero() {
		return ""
	}
	if len(layout) > 0 {
		return t.Format(layout[0])
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

func (g *Generator) groupBy(key string, nodes []*content.Node) []Group {
	groups := make(map[string][]*content.Node)

	for _, node := range nodes {
		// Extract date from node config
		dateVal, ok := node.Config["date"]
		if !ok {
			continue
		}

		var t time.Time
		switch v := dateVal.(type) {
		case time.Time:
			t = v
		case string:
			parsed, err := time.Parse(time.RFC3339, v)
			if err != nil {
				continue
			}
			t = parsed
		}

		// Get grouping key
		var groupKey string
		switch key {
		case "year":
			groupKey = strconv.Itoa(t.Year())
		case "month":
			groupKey = t.Format("2006-01")
		case "day":
			groupKey = t.Format("2006-01-02")
		default:
			continue
		}

		groups[groupKey] = append(groups[groupKey], node)
	}

	// Convert to sorted slice
	var result []Group
	for k, items := range groups {
		result = append(result, Group{Key: k, Items: items})
	}

	// Sort descending chronological order
	sort.Slice(result, func(i, j int) bool {
		return result[i].Key > result[j].Key
	})

	return result
}

func (g *Generator) renderPagination(data map[string]any) template.HTML {
	var buf bytes.Buffer
	tpl := g.templ.Lookup("macros/pagination.html")
	if tpl == nil {
		return ""
	}

	err := tpl.Execute(&buf, map[string]any{
		"Page": data["Page"],
		"Site": map[string]any{
			"Config": g.site.ToConfig(),
		},
	})
	if err != nil {
		return template.HTML(fmt.Sprintf("<!-- Pagination error: %v -->", err))
	}
	return template.HTML(buf.String())
}

func startWith(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

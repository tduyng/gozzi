package generator

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"reflect"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/tduyng/gozzi/app/content"
)

func (g *Generator) CreateFuncMap() template.FuncMap {
	return template.FuncMap{
		"add":            addNumbers,
		"and":            andLogic,
		"asset":          g.assetPath,
		"contains":       contains,
		"date":           formatDate,
		"default":        defaultValue,
		"dict":           createDictionary,
		"eq":             eq,
		"first":          firstElement,
		"get_section":    g.getSection,
		"group_by":       g.groupBy,
		"has_prefix":     strings.HasPrefix,
		"has_suffix":     strings.HasSuffix,
		"join":           strings.Join,
		"last":           lastElement,
		"limit":          limit,
		"load":           loadDataToHTML,
		"load_attribute": loadAttributeToHTML,
		"lower":          strings.ToLower,
		"markdown":       g.renderMarkdown,
		"ne":             ne,
		"now":            time.Now,
		"or":             orLogic,
		"pluralize":      pluralize,
		"priority":       priority,
		"replace":        strings.ReplaceAll,
		"reverse":        reverse,
		"safe":           safeHTML,
		"split":          strings.Split,
		"sub":            func(a, b int) int { return a - b },
		"to_date":        parseDate,
		"trim":           strings.TrimSpace,
		"upper":          strings.ToUpper,
		"urlize":         urlize,
		"where":          where,
		// macros
		"pagination": g.renderPagination,
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

func loadAttributeToHTML(path string) template.HTML {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Error reading file %s: %v", path, err)
		return ""
	}
	escaped := template.HTMLEscapeString(string(content))
	return template.HTML(escaped)
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
		return baseURL + "/" + pagePath + relPath
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

func eq(a, b any) bool {
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func ne(a, b any) bool {
	return !eq(a, b)
}

func firstElement(items any) any {
	switch reflect.TypeOf(items).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(items)
		if s.Len() == 0 {
			return nil
		}
		return s.Index(0).Interface()
	}
	return nil
}

func lastElement(items any) any {
	switch reflect.TypeOf(items).Kind() {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(items)
		if s.Len() == 0 {
			return nil
		}
		return s.Index(s.Len() - 1).Interface()
	}
	return nil
}

func andLogic(values ...any) bool {
	for _, v := range values {
		if !isTruthy(v) {
			return false
		}
	}
	return true
}

func orLogic(values ...any) bool {
	return slices.ContainsFunc(values, isTruthy)
}

func isTruthy(v any) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case int:
		return val != 0
	case string:
		return val != ""
	case []any:
		return len(val) > 0
	default:
		return true
	}
}

func contains(haystack, needle any) bool {
	switch h := haystack.(type) {
	case string:
		n := fmt.Sprintf("%v", needle)
		return strings.Contains(h, n)
	case []any:
		for _, item := range h {
			if eq(item, needle) {
				return true
			}
		}
	}
	return false
}

func addNumbers(a, b any) any {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			return a + b
		case float64:
			return float64(a) + b
		}
	case float64:
		switch b := b.(type) {
		case int:
			return a + float64(b)
		case float64:
			return a + b
		}
	}
	return 0
}

func parseDate(layout, value string) time.Time {
	t, _ := time.Parse(layout, value)
	return t
}

func pluralize(singular string, count any) string {
	var c int
	switch v := count.(type) {
	case int:
		c = v
	case int64:
		c = int(v)
	default:
		return singular
	}

	if c == 1 {
		return singular
	}
	return singular + "s" // Simple pluralization, extend for irregular forms
}

func createDictionary(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("invalid number of arguments for dict")
	}

	dict := make(map[string]any)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func (g *Generator) renderMarkdown(input string) template.HTML {
	var buf bytes.Buffer

	md := g.parser.GetMarkdownProcessor()

	if err := md.Convert([]byte(input), &buf); err != nil {
		return template.HTML("")
	}
	return template.HTML(buf.String())
}

func where(sections []any, field string, value any) []any {
	var result []any
	for _, s := range sections {
		sectionMap, ok := s.(map[string]any)
		if !ok {
			continue
		}

		fieldValue, exists := sectionMap[field]
		if exists && fmt.Sprintf("%v", fieldValue) == fmt.Sprintf("%v", value) {
			result = append(result, s)
		}
	}
	return result
}

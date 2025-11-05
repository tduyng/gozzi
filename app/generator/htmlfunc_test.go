package generator

import (
	"html/template"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/parser"
	tplengine "github.com/tduyng/gozzi/app/template"
)

func TestUrlize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple string", "hello world", "hello-world"},
		{"mixed case", "Hello World", "hello-world"},
		{"special characters", "Hello, World!", "hello-world"},
		{"multiple spaces", "hello    world", "hello-world"},
		{"multiple dashes", "hello---world", "hello-world"},
		{"leading/trailing spaces", "  hello world  ", "hello-world"},
		{"numbers", "test 123 post", "test-123-post"},
		{"underscores", "hello_world_test", "helloworldtest"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := urlize(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSafeHTML(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected template.HTML
	}{
		{"simple text", "hello", template.HTML("hello")},
		{"html tags", "<p>hello</p>", template.HTML("<p>hello</p>")},
		{"complex html", "<div class=\"test\">content</div>", template.HTML("<div class=\"test\">content</div>")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := safeHTML(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadDataToHTML(t *testing.T) {
	tempDir := t.TempDir()

	// Create test file
	testContent := "<p>Test content</p>"
	testFile := filepath.Join(tempDir, "test.html")
	require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))

	tests := []struct {
		name     string
		path     string
		expected template.HTML
	}{
		{"existing file", testFile, template.HTML(testContent)},
		{"non-existing file", filepath.Join(tempDir, "nonexistent.html"), template.HTML("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := loadDataToHTML(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadAttributeToHTML(t *testing.T) {
	tempDir := t.TempDir()

	// Create test file with content that needs escaping
	testContent := "<script>alert('test')</script>"
	testFile := filepath.Join(tempDir, "test.js")
	require.NoError(t, os.WriteFile(testFile, []byte(testContent), 0644))

	tests := []struct {
		name     string
		path     string
		expected template.HTML
	}{
		{"existing file with html", testFile, template.HTML("&lt;script&gt;alert(&#39;test&#39;)&lt;/script&gt;")},
		{"non-existing file", filepath.Join(tempDir, "nonexistent.js"), template.HTML("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := loadAttributeToHTML(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAssetPath(t *testing.T) {
	site := &config.Site{
		BaseURL: "https://example.com",
	}
	p := parser.NewParser(site)
	gen := &Generator{site: site, parser: p}

	tests := []struct {
		name     string
		relPath  string
		context  any
		expected string
	}{
		{"absolute http url", "http://other.com/asset.css", nil, "http://other.com/asset.css"},
		{"absolute https url", "https://other.com/asset.css", nil, "https://other.com/asset.css"},
		{"root relative path", "/css/style.css", nil, "https://example.com/css/style.css"},
		{"relative path with node context", "style.css", &content.Node{Slug: "blog"}, "https://example.com/blog/style.css"},
		{"relative path with node context and slash", "style.css", &content.Node{Slug: "blog/"}, "https://example.com/blog/style.css"},
		{"relative path without context", "style.css", nil, "https://example.com/style.css"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.assetPath(tt.relPath, tt.context)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPriority(t *testing.T) {
	tests := []struct {
		name     string
		values   []any
		expected string
	}{
		{"first non-empty string", []any{"", "second", "third"}, "second"},
		{"first non-nil value", []any{nil, "value", "other"}, "value"},
		{"all nil", []any{nil, nil, nil}, ""},
		{"mixed types", []any{nil, 0, "value"}, "0"},
		{"empty strings", []any{"", "", ""}, ""},
		{"non-string values", []any{123, 456}, "123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := priority(tt.values...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatDate(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 12, 30, 45, 0, time.UTC)

	tests := []struct {
		name     string
		time     time.Time
		layout   []string
		expected string
	}{
		{"default format", testTime, nil, "2024-01-15"},
		{"custom format", testTime, []string{"2006-01-02 15:04:05"}, "2024-01-15 12:30:45"},
		{"month year format", testTime, []string{"January 2006"}, "January 2024"},
		{"zero time", time.Time{}, nil, ""},
		{"zero time with layout", time.Time{}, []string{"2006-01-02"}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDate(tt.time, tt.layout...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLimit(t *testing.T) {
	nodes := []*content.Node{
		{Slug: "first"},
		{Slug: "second"},
		{Slug: "third"},
		{Slug: "fourth"},
	}

	tests := []struct {
		name     string
		max      any
		items    []*content.Node
		expected int
	}{
		{"limit to 2", 2, nodes, 2},
		{"limit to 3", 3, nodes, 3},
		{"limit greater than length", 10, nodes, 4},
		{"limit to 0", 0, nodes, 0},
		{"int64 limit", int64(2), nodes, 2},
		{"invalid type", "invalid", nodes, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := limit(tt.max, tt.items)
			assert.Len(t, result, tt.expected)
			if tt.expected > 0 {
				assert.Equal(t, tt.items[:tt.expected], result)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	nodes := []*content.Node{
		{Slug: "first"},
		{Slug: "second"},
		{Slug: "third"},
	}

	result := reverse(nodes)

	assert.Len(t, result, 3)
	assert.Equal(t, "third", result[0].Slug)
	assert.Equal(t, "second", result[1].Slug)
	assert.Equal(t, "first", result[2].Slug)

	// Verify original slice is unchanged
	assert.Equal(t, "first", nodes[0].Slug)
}

func TestDefaultValue(t *testing.T) {
	tests := []struct {
		name     string
		val      any
		def      string
		expected string
	}{
		{"nil value", nil, "default", "default"},
		{"empty string", "", "default", "default"},
		{"non-empty string", "value", "default", "value"},
		{"integer", 123, "default", "123"},
		{"boolean true", true, "default", "true"},
		{"boolean false", false, "default", "false"},
		{"zero value", 0, "default", "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defaultValue(tt.val, tt.def)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGroupBy(t *testing.T) {
	site := &config.Site{BaseURL: "https://example.com"}
	p := parser.NewParser(site)
	gen := &Generator{site: site, parser: p}

	date1 := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)
	date3 := time.Date(2023, 12, 10, 0, 0, 0, 0, time.UTC)

	nodes := []*content.Node{
		{Slug: "post1", Config: map[string]any{"date": date1}},
		{Slug: "post2", Config: map[string]any{"date": date2}},
		{Slug: "post3", Config: map[string]any{"date": date3}},
		{Slug: "post4", Config: map[string]any{"date": "2024-01-16T00:00:00Z"}},
		{Slug: "post5", Config: map[string]any{}}, // No date
	}

	tests := []struct {
		name     string
		key      string
		expected []string
	}{
		{"group by year", "year", []string{"2024", "2023"}},
		{"group by month", "month", []string{"2024-02", "2024-01", "2023-12"}},
		{"group by day", "day", []string{"2024-02-20", "2024-01-16", "2024-01-15", "2023-12-10"}},
		{"invalid key", "invalid", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.groupBy(tt.key, nodes)

			var keys []string
			for _, group := range result {
				keys = append(keys, group.Key)
			}

			if tt.expected == nil {
				assert.Nil(t, keys)
			} else {
				assert.Equal(t, tt.expected, keys)
			}
		})
	}
}

func TestEq(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"equal strings", "hello", "hello", true},
		{"different strings", "hello", "world", false},
		{"equal numbers", 123, 123, true},
		{"different numbers", 123, 456, false},
		{"mixed types equal", 123, "123", true},
		{"mixed types different", 123, "456", false},
		{"nil values", nil, nil, true},
		{"nil and string", nil, "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := eq(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNe(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected bool
	}{
		{"equal strings", "hello", "hello", false},
		{"different strings", "hello", "world", true},
		{"equal numbers", 123, 123, false},
		{"different numbers", 123, 456, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ne(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFirstElement(t *testing.T) {
	tests := []struct {
		name     string
		items    any
		expected any
	}{
		{"string slice", []string{"first", "second", "third"}, "first"},
		{"int slice", []int{1, 2, 3}, 1},
		{"empty slice", []string{}, nil},
		{"array", [3]string{"first", "second", "third"}, "first"},
		{"non-slice", "not a slice", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := firstElement(tt.items)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLastElement(t *testing.T) {
	tests := []struct {
		name     string
		items    any
		expected any
	}{
		{"string slice", []string{"first", "second", "third"}, "third"},
		{"int slice", []int{1, 2, 3}, 3},
		{"empty slice", []string{}, nil},
		{"array", [3]string{"first", "second", "third"}, "third"},
		{"non-slice", "not a slice", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := lastElement(tt.items)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAndLogic(t *testing.T) {
	tests := []struct {
		name     string
		values   []any
		expected bool
	}{
		{"all true", []any{true, "non-empty", 1}, true},
		{"contains false", []any{true, false, 1}, false},
		{"contains empty string", []any{true, "", 1}, false},
		{"contains zero", []any{true, "test", 0}, false},
		{"contains nil", []any{true, "test", nil}, false},
		{"empty values", []any{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := andLogic(tt.values...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestOrLogic(t *testing.T) {
	tests := []struct {
		name     string
		values   []any
		expected bool
	}{
		{"contains true", []any{false, true, 0}, true},
		{"all false", []any{false, "", 0, nil}, false},
		{"contains non-empty string", []any{false, "test", 0}, true},
		{"contains non-zero number", []any{false, "", 1}, true},
		{"empty values", []any{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := orLogic(tt.values...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		{"nil", nil, false},
		{"true", true, true},
		{"false", false, false},
		{"non-zero int", 1, true},
		{"zero int", 0, false},
		{"non-empty string", "test", true},
		{"empty string", "", false},
		{"non-empty slice", []any{1, 2}, true},
		{"empty slice", []any{}, false},
		{"other type", struct{}{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTruthy(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		haystack any
		needle   any
		expected bool
	}{
		{"string contains substring", "hello world", "world", true},
		{"string does not contain", "hello world", "test", false},
		{"slice contains item", []any{"hello", "world", "test"}, "world", true},
		{"slice does not contain", []any{"hello", "world"}, "test", false},
		{"slice contains number", []any{1, 2, 3}, 2, true},
		{"other type", 123, "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := contains(tt.haystack, tt.needle)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAddNumbers(t *testing.T) {
	tests := []struct {
		name     string
		a, b     any
		expected any
	}{
		{"int + int", 5, 3, 8},
		{"int + float64", 5, 3.5, 8.5},
		{"float64 + int", 5.5, 3, 8.5},
		{"float64 + float64", 5.5, 3.5, 9.0},
		{"invalid types", "string", 3, 0},
		{"int + string", 5, "string", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addNumbers(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseDate(t *testing.T) {
	tests := []struct {
		name     string
		layout   string
		value    string
		expected string
	}{
		{"valid date", "2006-01-02", "2024-01-15", "2024-01-15 00:00:00 +0000 UTC"},
		{"invalid date", "2006-01-02", "invalid", "0001-01-01 00:00:00 +0000 UTC"},
		{"datetime layout", "2006-01-02 15:04:05", "2024-01-15 12:30:45", "2024-01-15 12:30:45 +0000 UTC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDate(tt.layout, tt.value)
			assert.Equal(t, tt.expected, result.String())
		})
	}
}

func TestPluralize(t *testing.T) {
	tests := []struct {
		name     string
		singular string
		count    any
		expected string
	}{
		{"singular count", "item", 1, "item"},
		{"plural count", "item", 2, "items"},
		{"zero count", "item", 0, "items"},
		{"int64 count", "item", int64(5), "items"},
		{"invalid count type", "item", "invalid", "item"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pluralize(tt.singular, tt.count)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateDictionary(t *testing.T) {
	tests := []struct {
		name        string
		values      []any
		expected    map[string]any
		expectError bool
	}{
		{
			name:     "valid dictionary",
			values:   []any{"key1", "value1", "key2", "value2"},
			expected: map[string]any{"key1": "value1", "key2": "value2"},
		},
		{
			name:        "odd number of arguments",
			values:      []any{"key1", "value1", "key2"},
			expectError: true,
		},
		{
			name:        "non-string key",
			values:      []any{123, "value1", "key2", "value2"},
			expectError: true,
		},
		{
			name:     "empty dictionary",
			values:   []any{},
			expected: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := createDictionary(tt.values...)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestWhere(t *testing.T) {
	sections := []any{
		map[string]any{"name": "blog", "type": "section"},
		map[string]any{"name": "about", "type": "page"},
		map[string]any{"name": "contact", "type": "page"},
		map[string]any{"name": "projects", "type": "section"},
		"invalid", // Not a map
	}

	tests := []struct {
		name     string
		field    string
		value    any
		expected int
	}{
		{"filter by type section", "type", "section", 2},
		{"filter by type page", "type", "page", 2},
		{"filter by name", "name", "blog", 1},
		{"no matches", "type", "nonexistent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := where(sections, tt.field, tt.value)
			assert.Len(t, result, tt.expected)
		})
	}
}

func TestCreateFuncMap(t *testing.T) {
	site := &config.Site{BaseURL: "https://example.com"}
	p := parser.NewParser(site)

	// Initialize just the engine for testing CreateFuncMap
	// without needing full template loading
	engine := tplengine.NewEngine(&tplengine.EngineConfig{
		BaseURL:         site.BaseURL,
		ContentMap:      p.ContentMap,
		Markdown:        p.GetMarkdownProcessor(),
		StrictTemplates: site.StrictTemplates,
	})

	gen := &Generator{
		site:   site,
		parser: p,
		engine: engine,
	}

	funcMap := gen.CreateFuncMap()

	// Test that all expected functions are present
	expectedFunctions := []string{
		"add", "and", "asset", "contains", "date", "default", "dict",
		"eq", "first", "get_section", "group_by", "has_prefix", "has_suffix",
		"join", "last", "limit", "load", "load_attribute", "lower", "markdown",
		"ne", "now", "or", "pluralize", "priority", "replace", "reverse",
		"safe", "split", "sub", "to_date", "trim", "upper", "urlize",
		"where", "pagination",
	}

	for _, funcName := range expectedFunctions {
		t.Run("has function "+funcName, func(t *testing.T) {
			_, exists := funcMap[funcName]
			assert.True(t, exists, "Function %s should exist in funcMap", funcName)
		})
	}

	// Test a few functions work correctly
	t.Run("test add function", func(t *testing.T) {
		addFunc := funcMap["add"].(func(any, any) (any, error))
		result, err := addFunc(5, 3)
		require.NoError(t, err)
		assert.Equal(t, 8, result)
	})

	t.Run("test eq function", func(t *testing.T) {
		eqFunc := funcMap["eq"].(func(any, any) bool)
		result := eqFunc("test", "test")
		assert.True(t, result)
	})
}

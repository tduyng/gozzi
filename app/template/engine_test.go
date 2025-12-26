// Tests for template engine function registry and engine initialization.
// Verifies function registration, FuncMap building, and custom function handling.
package template

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/template/funcs"
	"github.com/yuin/goldmark"
)

func TestNewFuncRegistry(t *testing.T) {
	r := NewFuncRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
	if r.funcs == nil {
		t.Fatal("expected initialized funcs map")
	}
}

func TestFuncRegistry_Register(t *testing.T) {
	r := NewFuncRegistry()

	// Register a new function
	err := r.Register("test", FuncDef{
		Fn:          func() string { return "test" },
		Description: "Test function",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Attempt to register duplicate
	err = r.Register("test", FuncDef{
		Fn:          func() string { return "test2" },
		Description: "Duplicate",
	})
	if err == nil {
		t.Fatal("expected error for duplicate registration")
	}
	if err.Error() != `function "test" already registered` {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestFuncRegistry_MustRegister(t *testing.T) {
	r := NewFuncRegistry()

	// Should not panic for new function
	r.MustRegister("test", FuncDef{
		Fn:          func() string { return "test" },
		Description: "Test function",
	})

	// Should panic for duplicate
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for duplicate registration")
		}
	}()
	r.MustRegister("test", FuncDef{
		Fn:          func() string { return "test2" },
		Description: "Duplicate",
	})
}

func TestFuncRegistry_Get(t *testing.T) {
	r := NewFuncRegistry()
	testFn := func() string { return "test" }

	r.MustRegister("test", FuncDef{
		Fn:          testFn,
		Description: "Test function",
	})

	// Get existing function
	def, ok := r.Get("test")
	if !ok {
		t.Fatal("expected function to be found")
	}
	if def.Description != "Test function" {
		t.Errorf("expected description 'Test function', got %q", def.Description)
	}

	// Get non-existent function
	_, ok = r.Get("nonexistent")
	if ok {
		t.Fatal("expected function not to be found")
	}
}

func TestFuncRegistry_List(t *testing.T) {
	r := NewFuncRegistry()

	// Empty registry
	if len(r.List()) != 0 {
		t.Errorf("expected empty list, got %d items", len(r.List()))
	}

	// Add functions
	r.MustRegister("func1", FuncDef{Fn: func() {}, Description: "First"})
	r.MustRegister("func2", FuncDef{Fn: func() {}, Description: "Second"})
	r.MustRegister("func3", FuncDef{Fn: func() {}, Description: "Third"})

	names := r.List()
	if len(names) != 3 {
		t.Fatalf("expected 3 functions, got %d", len(names))
	}

	// Check all names are present (order doesn't matter)
	expected := map[string]bool{"func1": true, "func2": true, "func3": true}
	for _, name := range names {
		if !expected[name] {
			t.Errorf("unexpected function name: %q", name)
		}
		delete(expected, name)
	}
	if len(expected) > 0 {
		t.Errorf("missing function names: %v", expected)
	}
}

func TestFuncRegistry_BuildFuncMap_WithoutContext(t *testing.T) {
	r := NewFuncRegistry()

	testFn := func() string { return "test" }
	r.MustRegister("test", FuncDef{
		Fn:          testFn,
		RequiresCtx: false,
		Description: "Test function",
	})

	// Build without context
	fm := r.BuildFuncMap(nil)
	if len(fm) != 1 {
		t.Fatalf("expected 1 function, got %d", len(fm))
	}

	if _, ok := fm["test"]; !ok {
		t.Fatal("expected 'test' function in map")
	}
}

func TestFuncRegistry_BuildFuncMap_RequiresContext(t *testing.T) {
	r := NewFuncRegistry()

	// Register context-requiring function
	r.MustRegister("asset", FuncDef{
		Fn:          nil,
		RequiresCtx: true,
		Description: "Asset function",
	})

	// Build without context - should not include function
	fm := r.BuildFuncMap(nil)
	if _, ok := fm["asset"]; ok {
		t.Fatal("expected 'asset' function not in map without context")
	}

	// Build with context - should include function
	ctx := &funcs.SiteContext{
		BaseURL:    "https://example.com",
		ContentMap: make(map[string]*content.Node),
		Markdown:   goldmark.New(),
	}
	fm = r.BuildFuncMap(ctx)
	if _, ok := fm["asset"]; !ok {
		t.Fatal("expected 'asset' function in map with context")
	}
}

func TestCreateDefaultRegistry(t *testing.T) {
	r := CreateDefaultRegistry()

	// Check that all expected functions are registered
	expectedFuncs := []string{
		// Core
		"add", "sub", "eq", "ne", "and", "or",
		// Collections
		"first", "last", "contains", "reverse", "concat", "sort_by", "limit", "where", "group_by", "related_posts",
		// Strings
		"lower", "upper", "trim", "replace", "split", "join",
		"has_prefix", "has_suffix", "starts_with", "ends_with", "urlize", "default", "priority", "pluralize",
		// Date
		"date", "to_date", "now",
		// Data
		"dict", "safe", "load", "load_attribute",
		// Site-specific
		"asset", "get_section", "markdown", "i18n",
		// I18n URL generation
		"langURL", "currentLang",
	}

	registeredFuncs := r.List()
	if len(registeredFuncs) != len(expectedFuncs) {
		t.Errorf("expected %d functions, got %d", len(expectedFuncs), len(registeredFuncs))
	}

	for _, name := range expectedFuncs {
		if _, ok := r.Get(name); !ok {
			t.Errorf("expected function %q to be registered", name)
		}
	}
}

func TestNewEngine(t *testing.T) {
	config := &EngineConfig{
		BaseURL:         "https://example.com",
		ContentMap:      make(map[string]*content.Node),
		Markdown:        goldmark.New(),
		StrictTemplates: true,
	}

	engine := NewEngine(config)
	if engine == nil {
		t.Fatal("expected non-nil engine")
	}
	if engine.registry == nil {
		t.Fatal("expected non-nil registry")
	}
	if engine.config != config {
		t.Fatal("expected config to be set")
	}
}

func TestEngine_CreateFuncMap(t *testing.T) {
	config := &EngineConfig{
		BaseURL:    "https://example.com",
		ContentMap: make(map[string]*content.Node),
		Markdown:   goldmark.New(),
	}

	engine := NewEngine(config)
	fm := engine.CreateFuncMap()

	// Should be a valid template.FuncMap
	if fm == nil {
		t.Fatal("expected non-nil FuncMap")
	}

	// Check some expected functions exist
	expectedFuncs := []string{"add", "lower", "date", "dict", "asset"}
	for _, name := range expectedFuncs {
		if _, ok := fm[name]; !ok {
			t.Errorf("expected function %q in FuncMap", name)
		}
	}
}

func TestEngine_RegisterCustomFunc(t *testing.T) {
	config := &EngineConfig{
		BaseURL:    "https://example.com",
		ContentMap: make(map[string]*content.Node),
		Markdown:   goldmark.New(),
	}

	engine := NewEngine(config)

	// Register custom function
	customFn := func(s string) string { return "custom: " + s }
	err := engine.RegisterCustomFunc("custom", customFn, "Custom test function")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it's in the registry
	def, ok := engine.registry.Get("custom")
	if !ok {
		t.Fatal("expected custom function to be registered")
	}
	if def.Description != "Custom test function" {
		t.Errorf("unexpected description: %q", def.Description)
	}

	// Verify it's in the FuncMap
	fm := engine.CreateFuncMap()
	if _, ok := fm["custom"]; !ok {
		t.Fatal("expected custom function in FuncMap")
	}

	// Try to register duplicate
	err = engine.RegisterCustomFunc("custom", customFn, "Duplicate")
	if err == nil {
		t.Fatal("expected error for duplicate custom function")
	}
}

func TestEngine_ListFunctions(t *testing.T) {
	config := &EngineConfig{
		BaseURL:    "https://example.com",
		ContentMap: make(map[string]*content.Node),
		Markdown:   goldmark.New(),
	}

	engine := NewEngine(config)
	funcs := engine.ListFunctions()

	// Should have default functions
	if len(funcs) == 0 {
		t.Fatal("expected default functions in list")
	}

	// Add custom function
	engine.RegisterCustomFunc("custom", func() {}, "Custom")
	funcsAfter := engine.ListFunctions()

	if len(funcsAfter) != len(funcs)+1 {
		t.Errorf("expected %d functions after adding custom, got %d", len(funcs)+1, len(funcsAfter))
	}
}

func TestEngine_Integration(t *testing.T) {
	// Create a complete engine and test template execution
	config := &EngineConfig{
		BaseURL: "https://example.com",
		ContentMap: map[string]*content.Node{
			"blog": {
				Path:   "content/blog/_index.md",
				Slug:   "blog",
				Config: map[string]any{"title": "Blog"},
			},
		},
		Markdown: goldmark.New(),
	}

	engine := NewEngine(config)
	fm := engine.CreateFuncMap()

	// Create a simple template
	tmpl, err := template.New("test").Funcs(fm).Parse(`
		{{- $sum := add 2 3 -}}
		Sum: {{ $sum }}
		Upper: {{ upper "hello" }}
		Date: {{ date now "2006-01-02" }}
	`)
	if err != nil {
		t.Fatalf("failed to parse template: %v", err)
	}

	// Execute template
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		t.Fatalf("failed to execute template: %v", err)
	}
}

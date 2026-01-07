// Package template provides template function registry and management.
package template

import (
	"fmt"
	"html/template"

	"github.com/tduyng/gozzi/app/content"
	"github.com/tduyng/gozzi/app/i18n"
	"github.com/tduyng/gozzi/app/template/funcs"
	"github.com/yuin/goldmark"
)

// FuncRegistry manages template functions with metadata.
type FuncRegistry struct {
	funcs map[string]FuncDef
}

// FuncDef defines a template function with metadata.
type FuncDef struct {
	Fn          any
	RequiresCtx bool // Requires site context
	Description string
}

// NewFuncRegistry creates a new function registry.
func NewFuncRegistry() *FuncRegistry {
	return &FuncRegistry{
		funcs: make(map[string]FuncDef),
	}
}

// Register adds a function to the registry.
func (r *FuncRegistry) Register(name string, def FuncDef) error {
	if _, exists := r.funcs[name]; exists {
		return fmt.Errorf("function %q already registered", name)
	}
	r.funcs[name] = def
	return nil
}

// MustRegister registers a function or panics if it already exists.
func (r *FuncRegistry) MustRegister(name string, def FuncDef) {
	if err := r.Register(name, def); err != nil {
		panic(err)
	}
}

// Get retrieves a function definition.
func (r *FuncRegistry) Get(name string) (FuncDef, bool) {
	def, ok := r.funcs[name]
	return def, ok
}

// List returns all registered function names.
func (r *FuncRegistry) List() []string {
	names := make([]string, 0, len(r.funcs))
	for name := range r.funcs {
		names = append(names, name)
	}
	return names
}

// BuildFuncMap creates a template.FuncMap with optional site context.
func (r *FuncRegistry) BuildFuncMap(ctx *funcs.SiteContext) template.FuncMap {
	funcMap := make(template.FuncMap)

	var siteFuncs *funcs.SiteFuncs
	if ctx != nil {
		siteFuncs = funcs.NewSiteFuncs(ctx)
	}

	for name, def := range r.funcs {
		if def.RequiresCtx && siteFuncs != nil {
			// Bind site-specific functions
			switch name {
			case "asset":
				funcMap[name] = siteFuncs.AssetPath
			case "get_section":
				funcMap[name] = siteFuncs.GetSection
			case "markdown":
				funcMap[name] = siteFuncs.RenderMarkdown
			case "i18n":
				funcMap[name] = siteFuncs.Translate
			case "langURL":
				funcMap[name] = siteFuncs.LangURL
			case "currentLang":
				funcMap[name] = siteFuncs.CurrentLang
			}
		} else if !def.RequiresCtx {
			funcMap[name] = def.Fn
		}
	}

	return funcMap
}

// CreateDefaultRegistry creates a registry with all built-in functions.
func CreateDefaultRegistry() *FuncRegistry {
	r := NewFuncRegistry()

	r.MustRegister("add", FuncDef{Fn: funcs.Add, Description: "Add two numbers"})
	r.MustRegister("sub", FuncDef{Fn: funcs.Sub, Description: "Subtract b from a"})
	r.MustRegister("eq", FuncDef{Fn: funcs.Eq, Description: "Test equality"})
	r.MustRegister("ne", FuncDef{Fn: funcs.Ne, Description: "Test inequality"})
	r.MustRegister("gt", FuncDef{Fn: funcs.Gt, Description: "Test greater than (a > b)"})
	r.MustRegister("ge", FuncDef{Fn: funcs.Ge, Description: "Test greater or equal (a >= b)"})
	r.MustRegister("lt", FuncDef{Fn: funcs.Lt, Description: "Test less than (a < b)"})
	r.MustRegister("le", FuncDef{Fn: funcs.Le, Description: "Test less or equal (a <= b)"})
	r.MustRegister("and", FuncDef{Fn: funcs.And, Description: "Logical AND"})
	r.MustRegister("or", FuncDef{Fn: funcs.Or, Description: "Logical OR"})
	r.MustRegister("cond", FuncDef{Fn: funcs.Cond, Description: "Ternary conditional (if-then-else)"})

	r.MustRegister("first", FuncDef{Fn: funcs.First, Description: "Get first element"})
	r.MustRegister("last", FuncDef{Fn: funcs.Last, Description: "Get last element"})
	r.MustRegister("len", FuncDef{Fn: funcs.Len, Description: "Get length of slice, array, map, or string"})
	r.MustRegister("slice", FuncDef{Fn: funcs.Slice, Description: "Create a slice from arguments"})
	r.MustRegister("contains", FuncDef{Fn: funcs.Contains, Description: "Check if collection contains value"})
	r.MustRegister("reverse", FuncDef{Fn: funcs.Reverse, Description: "Reverse a slice of nodes"})
	r.MustRegister("concat", FuncDef{Fn: funcs.Concat, Description: "Merge multiple slices into one"})
	r.MustRegister("sort_by", FuncDef{Fn: funcs.SortBy, Description: "Sort nodes by field"})
	r.MustRegister("limit", FuncDef{Fn: funcs.Limit, Description: "Limit number of items"})
	r.MustRegister("where", FuncDef{Fn: funcs.Where, Description: "Filter items by field value"})
	r.MustRegister("group_by", FuncDef{Fn: funcs.GroupBy, Description: "Group nodes by date field"})
	r.MustRegister("related_posts", FuncDef{
		Fn:          funcs.RelatedPosts,
		Description: "Find related posts using tag-based scoring",
	})

	r.MustRegister("lower", FuncDef{Fn: funcs.Lower, Description: "Convert to lowercase"})
	r.MustRegister("upper", FuncDef{Fn: funcs.Upper, Description: "Convert to uppercase"})
	r.MustRegister("trim", FuncDef{Fn: funcs.Trim, Description: "Trim whitespace"})
	r.MustRegister("replace", FuncDef{Fn: funcs.Replace, Description: "Replace substring"})
	r.MustRegister("split", FuncDef{Fn: funcs.Split, Description: "Split string"})
	r.MustRegister("join", FuncDef{Fn: funcs.Join, Description: "Join strings"})
	r.MustRegister("has_prefix", FuncDef{Fn: funcs.HasPrefix, Description: "Check string prefix"})
	r.MustRegister("has_suffix", FuncDef{Fn: funcs.HasSuffix, Description: "Check string suffix"})
	r.MustRegister("starts_with", FuncDef{Fn: funcs.StartsWith, Description: "Check string prefix (alias for has_prefix)"})
	r.MustRegister("ends_with", FuncDef{Fn: funcs.EndsWith, Description: "Check string suffix (alias for has_suffix)"})
	r.MustRegister("urlize", FuncDef{Fn: funcs.Urlize, Description: "Convert to URL slug"})
	r.MustRegister("default", FuncDef{Fn: funcs.Default, Description: "Default value if empty"})
	r.MustRegister("priority", FuncDef{Fn: funcs.Priority, Description: "First non-empty value"})
	r.MustRegister("pluralize", FuncDef{Fn: funcs.Pluralize, Description: "Pluralize word based on count"})

	r.MustRegister("date", FuncDef{Fn: funcs.FormatDate, Description: "Format date"})
	r.MustRegister("to_date", FuncDef{Fn: funcs.ParseDate, Description: "Parse date string"})
	r.MustRegister("now", FuncDef{Fn: funcs.Now, Description: "Current time"})

	r.MustRegister("dict", FuncDef{Fn: funcs.Dict, Description: "Create map from key-value pairs"})
	r.MustRegister("safe", FuncDef{Fn: funcs.SafeHTML, Description: "Mark HTML as safe"})
	r.MustRegister("load", FuncDef{Fn: funcs.LoadData, Description: "Load file as HTML"})
	r.MustRegister("load_attribute", FuncDef{Fn: funcs.LoadAttribute, Description: "Load file as escaped attribute"})

	r.MustRegister("asset", FuncDef{Fn: nil, RequiresCtx: true, Description: "Generate asset URL"})
	r.MustRegister("get_section", FuncDef{Fn: nil, RequiresCtx: true, Description: "Get section by path"})
	r.MustRegister("markdown", FuncDef{Fn: nil, RequiresCtx: true, Description: "Render markdown to HTML"})
	r.MustRegister("i18n", FuncDef{Fn: nil, RequiresCtx: true, Description: "Translate key to current language"})
	r.MustRegister("langURL", FuncDef{Fn: nil, RequiresCtx: true, Description: "Generate language-specific URL"})
	r.MustRegister("currentLang", FuncDef{Fn: nil, RequiresCtx: true, Description: "Get current language code"})

	return r
}

// EngineConfig holds configuration for the template engine.
type EngineConfig struct {
	BaseURL         string
	ContentMap      map[string]*content.Node
	Markdown        goldmark.Markdown
	StrictTemplates bool // Error on undefined variables
	I18n            *i18n.I18n
}

// Engine manages template loading and execution.
type Engine struct {
	registry *FuncRegistry
	config   *EngineConfig
}

// NewEngine creates a new template engine.
func NewEngine(config *EngineConfig) *Engine {
	return &Engine{
		registry: CreateDefaultRegistry(),
		config:   config,
	}
}

// CreateFuncMap creates a template.FuncMap with all registered functions.
func (e *Engine) CreateFuncMap() template.FuncMap {
	ctx := &funcs.SiteContext{
		BaseURL:    e.config.BaseURL,
		ContentMap: e.config.ContentMap,
		Markdown:   e.config.Markdown,
		I18n:       e.config.I18n,
	}
	return e.registry.BuildFuncMap(ctx)
}

// RegisterCustomFunc adds a custom function to the registry.
func (e *Engine) RegisterCustomFunc(name string, fn any, description string) error {
	return e.registry.Register(name, FuncDef{
		Fn:          fn,
		RequiresCtx: false,
		Description: description,
	})
}

// ListFunctions returns all registered function names.
func (e *Engine) ListFunctions() []string {
	return e.registry.List()
}

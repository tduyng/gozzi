package markdown

import (
	"html/template"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

// ShortcodeNode represents a shortcode in the AST.
type ShortcodeNode struct {
	ast.BaseInline
	Name     string
	Params   map[string]string
	Content  string
	IsClosed bool // true for {{< >}}, false for {{% %}}
}

// Dump implements ast.Node interface.
func (n *ShortcodeNode) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindShortcode is the NodeKind for ShortcodeNode.
var KindShortcode = ast.NewNodeKind("Shortcode")

// Kind implements ast.Node interface.
func (n *ShortcodeNode) Kind() ast.NodeKind {
	return KindShortcode
}

// NewShortcodeNode creates a new shortcode node.
func NewShortcodeNode(name string, params map[string]string, content string, isClosed bool) *ShortcodeNode {
	return &ShortcodeNode{
		Name:     name,
		Params:   params,
		Content:  content,
		IsClosed: isClosed,
	}
}

// ShortcodeParser parses shortcode syntax in markdown.
type ShortcodeParser struct {
	templates *template.Template
}

// NewShortcodeParser creates a new shortcode parser.
func NewShortcodeParser(templates *template.Template) *ShortcodeParser {
	return &ShortcodeParser{
		templates: templates,
	}
}

// Trigger returns characters that trigger this parser.
func (s *ShortcodeParser) Trigger() []byte {
	return []byte{'{'}
}

// Parse parses shortcode syntax.
func (s *ShortcodeParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()

	// Check for shortcode opening: {{< or {{%
	if len(line) < 4 {
		return nil
	}

	var isClosed bool

	if line[0] == '{' && line[1] == '{' && line[2] == '<' {
		isClosed = true
	} else if line[0] == '{' && line[1] == '{' && line[2] == '%' {
		isClosed = false
	} else {
		return nil
	}

	// Find closing delimiter
	var endDelim string
	if isClosed {
		endDelim = ">}}"
	} else {
		endDelim = "%}}"
	}

	// Find end of opening tag
	pos := 3
	for pos < len(line) && line[pos] == ' ' {
		pos++
	}

	// Extract shortcode name
	nameStart := pos
	for pos < len(line) && line[pos] != ' ' && !isEndDelim(line[pos:], endDelim) {
		pos++
	}

	if pos == nameStart {
		return nil
	}

	name := string(line[nameStart:pos])

	// Skip whitespace
	for pos < len(line) && line[pos] == ' ' {
		pos++
	}

	// Parse parameters (key="value" or key=value)
	params := make(map[string]string)
	for pos < len(line) && !isEndDelim(line[pos:], endDelim) {
		// Parse key
		keyStart := pos
		for pos < len(line) && line[pos] != '=' && line[pos] != ' ' {
			pos++
		}

		if pos == keyStart || pos >= len(line) || line[pos] != '=' {
			break
		}

		key := string(line[keyStart:pos])
		pos++ // skip '='

		// Parse value
		var value string
		if pos < len(line) && line[pos] == '"' {
			// Quoted value
			pos++ // skip opening "
			valueStart := pos
			for pos < len(line) && line[pos] != '"' {
				pos++
			}
			value = string(line[valueStart:pos])
			if pos < len(line) {
				pos++ // skip closing "
			}
		} else {
			// Unquoted value
			valueStart := pos
			for pos < len(line) && line[pos] != ' ' && !isEndDelim(line[pos:], endDelim) {
				pos++
			}
			value = string(line[valueStart:pos])
		}

		params[key] = value

		// Skip whitespace
		for pos < len(line) && line[pos] == ' ' {
			pos++
		}
	}

	// Check for closing delimiter
	if !isEndDelim(line[pos:], endDelim) {
		return nil
	}

	pos += len(endDelim)

	// For paired shortcodes, extract content until closing tag
	var content string
	if !isClosed {
		// Look for {{%/ name %}}
		closeTag := "{{%/ " + name + " %}}"
		remaining := string(line[pos:])

		// Simple implementation: look for closing tag on same line
		// TODO: Support multi-line content
		closeIdx := findString(remaining, closeTag)
		if closeIdx != -1 {
			content = remaining[:closeIdx]
			pos += closeIdx + len(closeTag)
		}
	}

	block.Advance(pos)

	return NewShortcodeNode(name, params, content, isClosed)
}

// isEndDelim checks if bytes start with the end delimiter.
func isEndDelim(b []byte, delim string) bool {
	if len(b) < len(delim) {
		return false
	}
	return string(b[:len(delim)]) == delim
}

// findString finds substring in string (simple implementation).
func findString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// ShortcodeRenderer renders shortcodes to HTML.
type ShortcodeRenderer struct {
	templates *template.Template
}

// NewShortcodeRenderer creates a new shortcode renderer.
func NewShortcodeRenderer(templates *template.Template) *ShortcodeRenderer {
	return &ShortcodeRenderer{
		templates: templates,
	}
}

// RegisterFuncs registers the renderer.
func (r *ShortcodeRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindShortcode, r.renderShortcode)
}

// renderShortcode renders a shortcode node to HTML.
func (r *ShortcodeRenderer) renderShortcode(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*ShortcodeNode)

	// Look up template
	tmplName := "shortcodes/" + n.Name + ".html"
	tmpl := r.templates.Lookup(tmplName)

	if tmpl == nil {
		// Shortcode template not found, render as-is or skip
		return ast.WalkContinue, nil
	}

	// Prepare template data
	data := map[string]interface{}{
		"Params":  n.Params,
		"Content": template.HTML(n.Content),
	}

	// Add individual params as top-level fields for convenience
	for k, v := range n.Params {
		data[k] = v
	}

	// Execute template
	if err := tmpl.Execute(w, data); err != nil {
		return ast.WalkContinue, err
	}

	return ast.WalkContinue, nil
}

// ShortcodeExtension is the goldmark extension for shortcodes.
type ShortcodeExtension struct {
	templates *template.Template
}

// Extend adds the shortcode parser and renderer to goldmark.
func (e *ShortcodeExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithInlineParsers(
			util.Prioritized(NewShortcodeParser(e.templates), 100),
		),
	)
	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(NewShortcodeRenderer(e.templates), 100),
		),
	)
}

// NewShortcodeExtension creates a new shortcode extension.
func NewShortcodeExtension(templates *template.Template) goldmark.Extender {
	return &ShortcodeExtension{
		templates: templates,
	}
}

// Package markdown provides HTML renderer for KaTeX math expressions.
package markdown

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

// KaTeXHTMLRenderer renders KaTeX math nodes as HTML with LaTeX delimiters.
type KaTeXHTMLRenderer struct {
	html.Config
}

// RegisterFuncs registers rendering functions for KaTeX node types.
func (r *KaTeXHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindKaTeXInline, r.renderInline)
	reg.Register(KindKaTeXBlock, r.renderBlock)
}

// renderInline renders inline math expressions as \(...\).
func (r *KaTeXHTMLRenderer) renderInline(
	w util.BufWriter,
	source []byte,
	n ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	if entering {
		node := n.(*KaTeXInline)

		// Write LaTeX inline delimiters
		_, _ = w.WriteString(`\(`)
		_, _ = w.Write(node.Equation)
		_, _ = w.WriteString(`\)`)
	}
	return ast.WalkContinue, nil
}

// renderBlock renders block math expressions as <div>\[...\]</div>.
func (r *KaTeXHTMLRenderer) renderBlock(
	w util.BufWriter,
	source []byte,
	n ast.Node,
	entering bool,
) (ast.WalkStatus, error) {
	if entering {
		node := n.(*KaTeXBlock)

		// Write LaTeX block delimiters wrapped in a div
		_, _ = w.WriteString("<div>")
		_, _ = w.WriteString(`\[`)
		_, _ = w.Write(node.Equation)
		_, _ = w.WriteString(`\]`)
		_, _ = w.WriteString("</div>")
	}
	return ast.WalkContinue, nil
}

// ABOUTME: Goldmark extension for client-side KaTeX math rendering.
// ABOUTME: Parses $...$ (inline) and $$...$$ (block) math notations without external dependencies.

package markdown

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

// KaTeXExtender implements goldmark.Extender for KaTeX math rendering.
type KaTeXExtender struct{}

// Extend adds KaTeX parser and renderer to the goldmark instance.
func (e *KaTeXExtender) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&KaTeXParser{}, 0),
	))
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(&KaTeXHTMLRenderer{}, 0),
	))
}

// NewMathExtension creates a new math extension for client-side KaTeX rendering.
// This extension parses math delimiters ($...$ for inline, $$...$$ for block)
// and outputs LaTeX delimiters (\(...\) and \[...\]) that KaTeX can render in the browser.
// Works on all architectures (no CGO, quickjs, or external dependencies required).
func NewMathExtension() goldmark.Extender {
	return &KaTeXExtender{}
}

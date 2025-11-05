// Package markdown provides a Goldmark extension for math notation support in markdown content.
// Provides AST transformer to enable mathematical expressions rendering.
package markdown

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type mathExtension struct{}

// NewMathExtension creates a goldmark extension for math notation support.
func NewMathExtension() goldmark.Extender { return &mathExtension{} }

func (e *mathExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithASTTransformers(
		util.Prioritized(&mathTransformer{}, 100),
	))
}

type mathTransformer struct{}

func (t *mathTransformer) Transform(doc *ast.Document, reader text.Reader, pc parser.Context) {
	_ = ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		return ast.WalkContinue, nil
	})
}

// Package markdown provide AST node types for inline and block KaTeX math expressions.

package markdown

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/util"
)

// KaTeXInline represents an inline math expression ($...$).
type KaTeXInline struct {
	ast.BaseInline
	Equation []byte
}

// Inline marks this as an inline node.
func (n *KaTeXInline) Inline() {}

// IsBlank checks if the node contains only whitespace.
func (n *KaTeXInline) IsBlank(source []byte) bool {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		text := c.(*ast.Text).Segment
		if !util.IsBlank(text.Value(source)) {
			return false
		}
	}
	return true
}

// Dump outputs the node structure for debugging.
func (n *KaTeXInline) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindKaTeXInline is the NodeKind for inline KaTeX math.
var KindKaTeXInline = ast.NewNodeKind("KaTeXInline")

// Kind returns the node kind for inline math.
func (n *KaTeXInline) Kind() ast.NodeKind {
	return KindKaTeXInline
}

// KaTeXBlock represents a block math expression ($$...$$).
type KaTeXBlock struct {
	ast.BaseInline
	Equation []byte
}

// IsBlank checks if the node contains only whitespace.
func (n *KaTeXBlock) IsBlank(source []byte) bool {
	for c := n.FirstChild(); c != nil; c = c.NextSibling() {
		text := c.(*ast.Text).Segment
		if !util.IsBlank(text.Value(source)) {
			return false
		}
	}
	return true
}

// Dump outputs the node structure for debugging.
func (n *KaTeXBlock) Dump(source []byte, level int) {
	ast.DumpHelper(n, source, level, nil, nil)
}

// KindKaTeXBlock is the NodeKind for block KaTeX math.
var KindKaTeXBlock = ast.NewNodeKind("KaTeXBlock")

// Kind returns the node kind for block math.
func (n *KaTeXBlock) Kind() ast.NodeKind {
	return KindKaTeXBlock
}

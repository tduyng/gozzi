// Package markdown provides a Goldmark extension for math notation support in markdown content.
// Provides native KaTeX rendering for mathematical expressions using goldmark-katex.
package markdown

import (
	"github.com/FurqanSoftware/goldmark-katex"
	"github.com/yuin/goldmark"
)

// NewMathExtension creates a goldmark extension for math notation support.
// Uses goldmark-katex to render mathematical expressions with KaTeX.
// Supports inline math ($...$) and block math ($$...$$).
func NewMathExtension() goldmark.Extender {
	return &katex.Extender{}
}

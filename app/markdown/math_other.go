//go:build !amd64

// Package markdown provides a Goldmark extension for math notation support in markdown content.
// This version is for non-amd64 architectures (like ARM) where quickjs is not supported.
// Math rendering is disabled on these platforms.
package markdown

import (
	"github.com/yuin/goldmark"
)

// noopExtender is a no-op extension for architectures that don't support quickjs.
type noopExtender struct{}

func (e *noopExtender) Extend(m goldmark.Markdown) {
	// No-op: math rendering not available on this architecture
}

// NewMathExtension creates a no-op math extension for non-amd64 architectures.
// Math notation will not be rendered on ARM and other non-amd64 platforms.
func NewMathExtension() goldmark.Extender {
	return &noopExtender{}
}

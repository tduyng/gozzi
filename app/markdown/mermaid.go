// Package markdown provide mermaid diagram rendering extension for goldmark markdown processor.
// Provides client-side rendering for mermaid code blocks using MermaidJS.
package markdown

import (
	"github.com/yuin/goldmark"
	"go.abhg.dev/goldmark/mermaid"
)

// NewMermaidExtension creates a new mermaid diagram rendering extension.
func NewMermaidExtension() goldmark.Extender {
	return &mermaid.Extender{
		// Use client-side rendering (no CLI dependency required)
		RenderMode: mermaid.RenderModeClient,

		// Container tag for diagrams
		ContainerTag: "pre",

		// Automatically include MermaidJS script from CDN
		// Set to true if users want to manage the script themselves
		NoScript: false,

		// Use latest version from CDN (can be overridden)
		// Leave empty to use default (latest from cdn.jsdelivr.net)
		MermaidURL: "",
	}
}

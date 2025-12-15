// Package markdown provides mermaid diagram rendering extension for goldmark.
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

		// Don't auto-inject script - users manage MermaidJS in their templates
		// This gives users full control over version, CDN choice, and configuration
		NoScript: true,

		// MermaidURL not used when NoScript is true
		MermaidURL: "",
	}
}

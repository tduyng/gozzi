//go:build amd64

package markdown

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func TestMermaidExtension_WithMath(t *testing.T) {
	// Test that mermaid works alongside KaTeX (amd64 only)
	md := goldmark.New(
		goldmark.WithExtensions(
			NewMathExtension(),
			NewMermaidExtension(),
		),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	input := `# Mixed Content

Math: $E = mc^2$

` + "```mermaid" + `
graph TD
    A-->B
` + "```" + `

More text.`

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	output := buf.String()

	// Should contain KaTeX math
	if !strings.Contains(output, "class=\"katex\"") {
		t.Error("Expected KaTeX math in output")
	}

	// Should contain Mermaid diagram
	if !strings.Contains(output, "class=\"mermaid\"") {
		t.Error("Expected Mermaid diagram in output")
	}

	// Both should coexist
	if !strings.Contains(output, "--&gt;") { // HTML entities
		t.Error("Expected mermaid content in output")
	}
}

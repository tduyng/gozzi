package markdown

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func TestNewMermaidExtension(t *testing.T) {
	ext := NewMermaidExtension()
	if ext == nil {
		t.Fatal("NewMermaidExtension() returned nil")
	}
}

func TestMermaidExtension_Extend(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(NewMermaidExtension()),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	if md == nil {
		t.Fatal("Failed to create goldmark instance with mermaid extension")
	}
}

func TestMermaidExtensionIntegration(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(NewMermaidExtension()),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	tests := []struct {
		name            string
		input           string
		expectedContain []string
		notContain      []string
	}{
		{
			name:  "simple flowchart",
			input: "```mermaid\ngraph TD\n    A-->B\n```",
			expectedContain: []string{
				"<pre",
				"class=\"mermaid\"",
				"graph TD",
				"--&gt;",  // HTML entities are escaped
				"<script", // MermaidJS auto-injected
				"mermaid.initialize",
			},
			notContain: []string{},
		},
		{
			name:  "sequence diagram",
			input: "```mermaid\nsequenceDiagram\n    Alice->>Bob: Hello\n```",
			expectedContain: []string{
				"<pre",
				"class=\"mermaid\"",
				"sequenceDiagram",
				"-&gt;&gt;", // HTML entities are escaped
			},
			notContain: []string{},
		},
		{
			name:  "gantt chart",
			input: "```mermaid\ngantt\n    title Project\n    section Tasks\n    Task 1: 2024-01-01, 30d\n```",
			expectedContain: []string{
				"<pre",
				"class=\"mermaid\"",
				"gantt",
				"title Project",
			},
			notContain: []string{},
		},
		{
			name:  "no mermaid in regular code block",
			input: "```js\nconsole.log('hello')\n```",
			expectedContain: []string{
				"<pre>",
				"<code",
				"console.log",
			},
			notContain: []string{
				"class=\"mermaid\"",
			},
		},
		{
			name:  "multiple mermaid blocks",
			input: "```mermaid\ngraph LR\n    A-->B\n```\n\nSome text\n\n```mermaid\ngraph TD\n    C-->D\n```",
			expectedContain: []string{
				"class=\"mermaid\"",
				"--&gt;", // HTML entities are escaped
				"Some text",
			},
			notContain: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := md.Convert([]byte(tt.input), &buf); err != nil {
				t.Fatalf("Conversion failed: %v", err)
			}

			output := buf.String()

			// Check expected content
			for _, expected := range tt.expectedContain {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, got:\n%s", expected, output)
				}
			}

			// Check for content that should not be present
			for _, notExpected := range tt.notContain {
				if strings.Contains(output, notExpected) {
					t.Errorf("Expected output NOT to contain %q, but it did:\n%s", notExpected, output)
				}
			}
		})
	}
}

func TestMermaidExtension_MixedContent(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(NewMermaidExtension()),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	input := `# Diagram Example

Here's a flowchart:

` + "```mermaid" + `
graph LR
    Start --> End
` + "```" + `

And regular text after.

` + "```js" + `
console.log('not mermaid');
` + "```"

	var buf bytes.Buffer
	if err := md.Convert([]byte(input), &buf); err != nil {
		t.Fatalf("Conversion failed: %v", err)
	}

	output := buf.String()

	// Should contain heading
	if !strings.Contains(output, "<h1") {
		t.Error("Expected H1 heading in output")
	}

	// Should contain mermaid diagram
	if !strings.Contains(output, "class=\"mermaid\"") {
		t.Error("Expected mermaid diagram in output")
	}

	if !strings.Contains(output, "--&gt;") { // HTML entities
		t.Error("Expected mermaid content in output")
	}

	// Should contain regular code block
	if !strings.Contains(output, "console.log") {
		t.Error("Expected regular code block in output")
	}

	// Should contain regular text
	if !strings.Contains(output, "regular text") {
		t.Error("Expected regular text in output")
	}

	// Should include MermaidJS script (auto-injected)
	if !strings.Contains(output, "<script") {
		t.Error("Expected MermaidJS <script> tag in output")
	}

	if !strings.Contains(output, "mermaid.initialize") {
		t.Error("Expected mermaid.initialize() call in output")
	}
}

func TestMermaidExtension_WithMath(t *testing.T) {
	// Test that mermaid works alongside KaTeX
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

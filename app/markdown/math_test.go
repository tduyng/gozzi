// ABOUTME: This file contains tests for the client-side KaTeX math extension.
// ABOUTME: Tests verify that math delimiters are preserved for browser-side rendering.

package markdown

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

func TestNewMathExtension(t *testing.T) {
	ext := NewMathExtension()
	assert.NotNil(t, ext)

	// Verify it implements goldmark.Extender
	assert.Implements(t, (*goldmark.Extender)(nil), ext)
}

func TestMathExtension_Extend(t *testing.T) {
	ext := NewMathExtension()

	// Create a basic goldmark instance
	md := goldmark.New()

	// This should not panic
	assert.NotPanics(t, func() {
		ext.Extend(md)
	})
}

func TestMathExtensionIntegration(t *testing.T) {
	// Test that the extension can be used with goldmark
	md := goldmark.New(
		goldmark.WithExtensions(
			NewMathExtension(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	require.NotNil(t, md)

	// Basic markdown rendering should still work
	var buf strings.Builder
	err := md.Convert([]byte("# Hello World"), &buf)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "<h1")
	assert.Contains(t, buf.String(), "Hello World")
}

func TestMathExtension_InlineMath(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(
			NewMathExtension(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name:  "simple inline math",
			input: "The equation $E = mc^2$ is famous.",
			contains: []string{
				`\(E = mc^2\)`,
			},
		},
		{
			name:  "multiple inline math",
			input: "We have $a = b$ and $c = d$ here.",
			contains: []string{
				`\(a = b\)`,
				`\(c = d\)`,
			},
		},
		{
			name:  "inline math with Greek letters",
			input: "The value of $\\pi$ is approximately 3.14.",
			contains: []string{
				`\(\pi\)`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			err := md.Convert([]byte(tt.input), &buf)
			require.NoError(t, err)

			output := buf.String()
			for _, expected := range tt.contains {
				assert.Contains(t, output, expected, "Output should contain: %s", expected)
			}
		})
	}
}

func TestMathExtension_BlockMath(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(
			NewMathExtension(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	tests := []struct {
		name     string
		input    string
		contains []string
	}{
		{
			name: "simple block math",
			input: `Block equation:

$$
A = \pi r^2
$$`,
			contains: []string{
				`<div>\[`,
				`A = \pi r^2`,
				`\]</div>`,
			},
		},
		{
			name: "block math with fractions",
			input: `The integral:

$$
\int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}
$$`,
			contains: []string{
				`<div>\[`,
				`\int_{-\infty}^{\infty}`,
				`\]</div>`,
			},
		},
		{
			name: "block math with matrix",
			input: `Matrix form:

$$
\begin{bmatrix}
a & b \\
c & d
\end{bmatrix}
$$`,
			contains: []string{
				`<div>\[`,
				`\begin{bmatrix}`,
				`\]</div>`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf strings.Builder
			err := md.Convert([]byte(tt.input), &buf)
			require.NoError(t, err)

			output := buf.String()
			for _, expected := range tt.contains {
				assert.Contains(t, output, expected, "Output should contain: %s", expected)
			}
		})
	}
}

func TestMathExtension_NoMath(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(
			NewMathExtension(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	input := `# Regular Document

This is just regular markdown with no math.

- Item 1
- Item 2

**Bold text** and *italic text*.`

	var buf strings.Builder
	err := md.Convert([]byte(input), &buf)
	require.NoError(t, err)

	output := buf.String()

	// Should not contain math classes when there's no math
	assert.NotContains(t, output, "math-inline")
	assert.NotContains(t, output, "math-block")
	assert.Contains(t, output, "<h1")
	assert.Contains(t, output, "<li")
	assert.Contains(t, output, "<strong")
	assert.Contains(t, output, "<em")
}

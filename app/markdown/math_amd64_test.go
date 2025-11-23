//go:build amd64

package markdown

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

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
				"katex",
				"E = mc^2",
			},
		},
		{
			name:  "multiple inline math",
			input: "We have $a = b$ and $c = d$ here.",
			contains: []string{
				"katex",
				"a = b",
				"c = d",
			},
		},
		{
			name:  "inline math with Greek letters",
			input: "The value of $\\pi$ is approximately 3.14.",
			contains: []string{
				"katex",
				"π", // KaTeX renders \pi as π
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
				"katex",
				"display",
				"A = \\pi r^2",
			},
		},
		{
			name: "block math with fractions",
			input: `The integral:

$$
\int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}
$$`,
			contains: []string{
				"katex",
				"display",
				"\\int",
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
				"katex",
				"display",
				"bmatrix",
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

func TestMathExtension_MixedContent(t *testing.T) {
	md := goldmark.New(
		goldmark.WithExtensions(
			NewMathExtension(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	input := `# Mathematical Document

This document contains both inline math like $E = mc^2$ and block math:

$$
F = ma
$$

Regular text continues here with more inline math $v = \frac{d}{t}$.

## Another Section

More block equations:

$$
\sum_{i=1}^{n} i = \frac{n(n+1)}{2}
$$

The end.`

	var buf strings.Builder
	err := md.Convert([]byte(input), &buf)
	require.NoError(t, err)

	output := buf.String()

	// Check for multiple instances of KaTeX
	assert.Contains(t, output, "katex")
	assert.Contains(t, output, "<h1")
	assert.Contains(t, output, "<h2")
	assert.Contains(t, output, "E = mc^2")
	assert.Contains(t, output, "F = ma")
	assert.Contains(t, output, "display")
}

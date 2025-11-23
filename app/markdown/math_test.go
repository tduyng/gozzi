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

	// Should not contain katex classes when there's no math
	assert.NotContains(t, output, "katex")
	assert.Contains(t, output, "<h1")
	assert.Contains(t, output, "<li")
	assert.Contains(t, output, "<strong")
	assert.Contains(t, output, "<em")
}

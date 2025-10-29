package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
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
	)

	require.NotNil(t, md)

	// Basic markdown rendering should still work
	var buf strings.Builder
	err := md.Convert([]byte("# Hello World"), &buf)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "<h1")
	assert.Contains(t, buf.String(), "Hello World")
}

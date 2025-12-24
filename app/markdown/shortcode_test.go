package markdown

import (
	"bytes"
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func TestShortcodeNode(t *testing.T) {
	node := NewShortcodeNode("test", map[string]string{"key": "value"}, "content", true)

	assert.Equal(t, "test", node.Name)
	assert.Equal(t, "value", node.Params["key"])
	assert.Equal(t, "content", node.Content)
	assert.True(t, node.IsClosed)
	assert.Equal(t, KindShortcode, node.Kind())
}

func TestShortcodeExtension_SelfClosing(t *testing.T) {
	// Create shortcode templates
	tmpl := template.Must(template.New("").Parse(`
{{- define "shortcodes/youtube.html" -}}
<div class="video"><iframe src="https://youtube.com/embed/{{ .id }}"></iframe></div>
{{- end -}}
`))

	md := goldmark.New(
		goldmark.WithExtensions(NewShortcodeExtension(tmpl)),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	source := `# Test

This is a video: {{< youtube id="abc123" >}}

More content.`

	var buf bytes.Buffer
	err := md.Convert([]byte(source), &buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, `<div class="video">`)
	assert.Contains(t, output, `https://youtube.com/embed/abc123`)
}

func TestShortcodeExtension_Paired(t *testing.T) {
	// Create shortcode templates
	tmpl := template.Must(template.New("").Parse(`
{{- define "shortcodes/alert.html" -}}
<div class="alert alert-{{ .type }}">{{ .Content }}</div>
{{- end -}}
`))

	md := goldmark.New(
		goldmark.WithExtensions(NewShortcodeExtension(tmpl)),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	source := `# Test

{{{% alert type="warning" %}}This is a warning{{%/ alert %}}

More content.`

	var buf bytes.Buffer
	err := md.Convert([]byte(source), &buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, `<div class="alert alert-warning">`)
	assert.Contains(t, output, `This is a warning`)
}

func TestShortcodeExtension_MultipleParams(t *testing.T) {
	// Create shortcode templates
	tmpl := template.Must(template.New("").Parse(`
{{- define "shortcodes/image.html" -}}
<img src="{{ .src }}" alt="{{ .alt }}" loading="lazy">
{{- end -}}
`))

	md := goldmark.New(
		goldmark.WithExtensions(NewShortcodeExtension(tmpl)),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	source := `{{< image src="photo.jpg" alt="My photo" >}}`

	var buf bytes.Buffer
	err := md.Convert([]byte(source), &buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, `src="photo.jpg"`)
	assert.Contains(t, output, `alt="My photo"`)
	assert.Contains(t, output, `loading="lazy"`)
}

func TestShortcodeExtension_MissingTemplate(t *testing.T) {
	tmpl := template.New("")

	md := goldmark.New(
		goldmark.WithExtensions(NewShortcodeExtension(tmpl)),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	source := `{{< nonexistent id="test" >}}`

	var buf bytes.Buffer
	err := md.Convert([]byte(source), &buf)
	require.NoError(t, err)

	// Should not crash, just skip rendering
	output := buf.String()
	assert.NotContains(t, output, "nonexistent")
}

func TestShortcodeExtension_MixedWithMarkdown(t *testing.T) {
	tmpl := template.Must(template.New("").Parse(`
{{- define "shortcodes/note.html" -}}
<aside class="note">{{ .Content }}</aside>
{{- end -}}
`))

	md := goldmark.New(
		goldmark.WithExtensions(NewShortcodeExtension(tmpl)),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	source := `# Heading

Regular **markdown** content.

{{{% note %}}Important note here{{%/ note %}}

More *markdown*.`

	var buf bytes.Buffer
	err := md.Convert([]byte(source), &buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "<h1")
	assert.Contains(t, output, "<strong>markdown</strong>")
	assert.Contains(t, output, `<aside class="note">`)
	assert.Contains(t, output, "<em>markdown</em>")
}

func TestShortcodeParser_QuotedParams(t *testing.T) {
	tmpl := template.Must(template.New("").Parse(`
{{- define "shortcodes/link.html" -}}
<a href="{{ .url }}" title="{{ .title }}">{{ .text }}</a>
{{- end -}}
`))

	md := goldmark.New(
		goldmark.WithExtensions(NewShortcodeExtension(tmpl)),
		goldmark.WithParserOptions(parser.WithAutoHeadingID()),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	source := `{{< link url="https://example.com" title="Example Site" text="Click here" >}}`

	var buf bytes.Buffer
	err := md.Convert([]byte(source), &buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, `href="https://example.com"`)
	assert.Contains(t, output, `title="Example Site"`)
	assert.Contains(t, output, `Click here`)
}

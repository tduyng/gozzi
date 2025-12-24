package markdown

import (
	"html/template"
	"strings"
	"testing"
)

func TestShortcodeProcessor_SelfClosing(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse(`<div class="youtube" data-id="{{ .id }}"></div>`))
	tmpl = template.Must(tmpl.New("shortcodes/youtube.html").Parse(`<div class="youtube" data-id="{{ .id }}"></div>`))

	processor := NewShortcodeProcessor(tmpl)

	markdown := []byte(`# Test

{{< youtube id="abc123" >}}

Some text.`)

	result, err := processor.Process(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(result)
	if !strings.Contains(output, `<div class="youtube" data-id="abc123"></div>`) {
		t.Errorf("expected YouTube shortcode to be rendered, got: %s", output)
	}
	if strings.Contains(output, "{{<") {
		t.Errorf("shortcode syntax should be replaced, got: %s", output)
	}
}

func TestShortcodeProcessor_Paired(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse(`<div class="alert">{{ .Content }}</div>`))
	tmpl = template.Must(tmpl.New("shortcodes/alert.html").Parse(`<div class="alert">{{ .Content }}</div>`))

	processor := NewShortcodeProcessor(tmpl)

	markdown := []byte(`# Test

{{% alert %}}
This is content.
{%/ alert %}}

Some text.`)

	result, err := processor.Process(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(result)
	if !strings.Contains(output, `<div class="alert">`) {
		t.Errorf("expected alert shortcode to be rendered, got: %s", output)
	}
	if !strings.Contains(output, "This is content.") {
		t.Errorf("expected shortcode content to be preserved, got: %s", output)
	}
	if strings.Contains(output, "{{% ") {
		t.Errorf("shortcode syntax should be replaced, got: %s", output)
	}
}

func TestShortcodeProcessor_WithParameters(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse(`<img src="{{ .src }}" alt="{{ .alt }}">`))
	tmpl = template.Must(tmpl.New("shortcodes/image.html").Parse(`<img src="{{ .src }}" alt="{{ .alt }}">`))

	processor := NewShortcodeProcessor(tmpl)

	markdown := []byte(`{{< image src="/test.jpg" alt="Test Image" >}}`)

	result, err := processor.Process(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(result)
	if !strings.Contains(output, `src="/test.jpg"`) {
		t.Errorf("expected src parameter to be rendered, got: %s", output)
	}
	if !strings.Contains(output, `alt="Test Image"`) {
		t.Errorf("expected alt parameter to be rendered, got: %s", output)
	}
}

func TestShortcodeProcessor_MissingTemplate(t *testing.T) {
	tmpl := template.New("test")

	processor := NewShortcodeProcessor(tmpl)

	markdown := []byte(`{{< nonexistent >}}`)

	result, err := processor.Process(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// When template is missing, original shortcode should remain unchanged
	output := string(result)
	if !strings.Contains(output, "{{< nonexistent >}}") {
		t.Errorf("expected original shortcode to remain when template is missing, got: %s", output)
	}
}

func TestShortcodeProcessor_MultilineContent(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse(`<blockquote>{{ .Content }}</blockquote>`))
	tmpl = template.Must(tmpl.New("shortcodes/quote.html").Parse(`<blockquote>{{ .Content }}</blockquote>`))

	processor := NewShortcodeProcessor(tmpl)

	markdown := []byte(`{{% quote %}}
Line 1
Line 2
Line 3
{%/ quote %}}`)

	result, err := processor.Process(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(result)
	if !strings.Contains(output, "Line 1") || !strings.Contains(output, "Line 2") || !strings.Contains(output, "Line 3") {
		t.Errorf("expected all lines to be preserved, got: %s", output)
	}
}

func TestShortcodeProcessor_MismatchedTags(t *testing.T) {
	tmpl := template.Must(template.New("test").Parse(`<div>{{ .Content }}</div>`))
	tmpl = template.Must(tmpl.New("shortcodes/alert.html").Parse(`<div>{{ .Content }}</div>`))
	tmpl = template.Must(tmpl.New("shortcodes/quote.html").Parse(`<blockquote>{{ .Content }}</blockquote>`))

	processor := NewShortcodeProcessor(tmpl)

	// Mismatched opening and closing tags should not be processed
	markdown := []byte(`{{% alert %}}
Content
{%/ quote %}}`)

	result, err := processor.Process(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := string(result)
	// Should remain unchanged when tags don't match
	if !strings.Contains(output, "{{% alert %}}") {
		t.Errorf("expected original shortcode to remain when tags don't match, got: %s", output)
	}
}

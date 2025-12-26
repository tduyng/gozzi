package markdown

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// ShortcodeProcessor processes shortcodes in markdown before goldmark parsing.
type ShortcodeProcessor struct {
	templates *template.Template
	md        goldmark.Markdown
}

// NewShortcodeProcessor creates a new shortcode processor.
func NewShortcodeProcessor(templates *template.Template) *ShortcodeProcessor {
	// Create a basic markdown processor for shortcode content
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Typographer,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(), // Allow raw HTML
		),
	)

	return &ShortcodeProcessor{
		templates: templates,
		md:        md,
	}
}

// Process replaces shortcodes in markdown with rendered HTML.
func (sp *ShortcodeProcessor) Process(markdown []byte) ([]byte, error) {
	content := string(markdown)

	// Process self-closing shortcodes: {{< name key="value" >}}
	selfClosingPattern := regexp.MustCompile(`\{\{<\s*(\w+)((?:\s+\w+="[^"]*")*)\s*>\}\}`)
	content = selfClosingPattern.ReplaceAllStringFunc(content, func(match string) string {
		return sp.renderShortcode(match, true)
	})

	// Process paired shortcodes: {{% name %}}content{%/ name %}}
	// Opening: {{% ... %}}  Closing: {%/ ... %}}  (asymmetric - this is Hugo's syntax)
	// (?s) enables DOTALL mode so . matches newlines
	pairedPattern := regexp.MustCompile(`(?s)\{\{%\s*(\w+)((?:\s+\w+="[^"]*")*)\s*%\}\}(.*?)\{%/\s*(\w+)\s*%\}\}`)
	content = pairedPattern.ReplaceAllStringFunc(content, func(match string) string {
		// Verify opening and closing tags match
		matches := pairedPattern.FindStringSubmatch(match)
		if len(matches) >= 5 && matches[1] == matches[4] {
			return sp.renderShortcode(match, false)
		}
		// Tags don't match, return original
		return match
	})

	return []byte(content), nil
}

func (sp *ShortcodeProcessor) renderShortcode(match string, isSelfClosing bool) string {
	var name string
	var paramsStr string
	var content string

	if isSelfClosing {
		// Parse: {{< name key="value" >}}
		pattern := regexp.MustCompile(`\{\{<\s*(\w+)((?:\s+\w+="[^"]*")*)\s*>\}\}`)
		matches := pattern.FindStringSubmatch(match)
		if len(matches) < 3 {
			return match
		}
		name = matches[1]
		paramsStr = matches[2]
	} else {
		// Parse: {{% name key="value" %}}content{%/ name %}}
		// (?s) enables DOTALL mode so . matches newlines
		pattern := regexp.MustCompile(`(?s)\{\{%\s*(\w+)((?:\s+\w+="[^"]*")*)\s*%\}\}(.*?)\{%/\s*(\w+)\s*%\}\}`)
		matches := pattern.FindStringSubmatch(match)
		if len(matches) < 5 {
			return match
		}
		// Verify opening and closing tags match
		if matches[1] != matches[4] {
			return match
		}
		name = matches[1]
		paramsStr = matches[2]
		content = matches[3]

		// For {{% %}} shortcodes, process content as markdown
		// Trim whitespace but preserve intentional formatting
		trimmedContent := content
		// Only trim if there's significant whitespace
		if len(content) > 0 && (content[0] == '\n' || content[0] == '\r') {
			trimmedContent = content[1:]
		}
		if len(trimmedContent) > 0 && (trimmedContent[len(trimmedContent)-1] == '\n' || trimmedContent[len(trimmedContent)-1] == '\r') {
			trimmedContent = trimmedContent[:len(trimmedContent)-1]
		}

		// Convert markdown to HTML
		var htmlBuf bytes.Buffer
		if err := sp.md.Convert([]byte(trimmedContent), &htmlBuf); err == nil {
			content = htmlBuf.String()
		}
	}

	// Parse parameters
	params := make(map[string]string)
	paramPattern := regexp.MustCompile(`(\w+)="([^"]*)"`)
	paramMatches := paramPattern.FindAllStringSubmatch(paramsStr, -1)
	for _, pm := range paramMatches {
		if len(pm) >= 3 {
			params[pm[1]] = pm[2]
		}
	}

	// Look up template
	tmplName := "shortcodes/" + name + ".html"
	tmpl := sp.templates.Lookup(tmplName)
	if tmpl == nil {
		// Template not found, return original
		return match
	}

	// Prepare template data
	data := map[string]any{
		"Params":  params,
		"Content": template.HTML(content),
	}

	// Add individual params as top-level fields
	for k, v := range params {
		data[k] = v
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Sprintf("<!-- shortcode error: %v -->", err)
	}

	return buf.String()
}

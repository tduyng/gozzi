package markdown

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
)

// ShortcodeProcessor processes shortcodes in markdown before goldmark parsing.
type ShortcodeProcessor struct {
	templates *template.Template
}

// NewShortcodeProcessor creates a new shortcode processor.
func NewShortcodeProcessor(templates *template.Template) *ShortcodeProcessor {
	return &ShortcodeProcessor{
		templates: templates,
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

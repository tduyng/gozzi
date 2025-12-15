// Package template provides validation with detailed error reporting.
package template

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
)

// Error represents a template error with source context.
type Error struct {
	Template string // Template file name
	Line     int    // Line number (1-indexed)
	Column   int    // Column number (1-indexed)
	Snippet  string // Source code snippet
	Message  string // Error message
	Hint     string // Optional hint for fixing
}

// Error implements the error interface.
func (e *Error) Error() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Template error in %s", e.Template))
	if e.Line > 0 {
		sb.WriteString(fmt.Sprintf(" at line %d", e.Line))
		if e.Column > 0 {
			sb.WriteString(fmt.Sprintf(", column %d", e.Column))
		}
	}
	sb.WriteString(fmt.Sprintf(": %s\n", e.Message))

	if e.Snippet != "" {
		sb.WriteString(fmt.Sprintf("\n%s\n", e.Snippet))
	}

	if e.Hint != "" {
		sb.WriteString(fmt.Sprintf("\nHint: %s\n", e.Hint))
	}

	return sb.String()
}

// Validator validates templates before generation.
type Validator struct {
	templatesDir string
	funcMap      template.FuncMap
}

// NewValidator creates a new template validator.
func NewValidator(templatesDir string, funcMap template.FuncMap) *Validator {
	return &Validator{
		templatesDir: templatesDir,
		funcMap:      funcMap,
	}
}

// Validate checks all templates for syntax errors.
func (v *Validator) Validate() []Error {
	var errors []Error

	err := filepath.WalkDir(v.templatesDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		relPath, _ := filepath.Rel(v.templatesDir, path)
		templateErrors := v.validateTemplate(relPath, path)
		errors = append(errors, templateErrors...)

		return nil
	})

	if err != nil {
		errors = append(errors, Error{
			Template: v.templatesDir,
			Message:  fmt.Sprintf("Failed to walk templates directory: %v", err),
		})
	}

	return errors
}

func (v *Validator) validateTemplate(name, path string) []Error {
	content, err := os.ReadFile(path)
	if err != nil {
		return []Error{{
			Template: name,
			Message:  fmt.Sprintf("Failed to read template: %v", err),
		}}
	}

	tmpl := template.New(name).Funcs(v.funcMap)
	_, err = tmpl.Parse(string(content))
	if err != nil {
		return v.parseError(name, string(content), err)
	}

	return nil
}

func (v *Validator) parseError(name, content string, err error) []Error {
	errMsg := err.Error()
	lines := strings.Split(content, "\n")

	// Try to extract line number from error message
	// Go template errors often include line:column info
	var lineNum int
	var snippet string

	parts := strings.Split(errMsg, ":")
	if len(parts) >= 3 {
		_, _ = fmt.Sscanf(parts[1], "%d", &lineNum)
	}

	if lineNum > 0 && lineNum <= len(lines) {
		snippet = v.buildSnippet(lines, lineNum)
	}

	return []Error{{
		Template: name,
		Line:     lineNum,
		Message:  extractMessage(errMsg),
		Snippet:  snippet,
		Hint:     suggestFix(errMsg),
	}}
}

func (v *Validator) buildSnippet(lines []string, errorLine int) string {
	var sb strings.Builder

	start := max(0, errorLine-3)
	end := min(len(lines), errorLine+2)

	for i := start; i < end; i++ {
		lineNum := i + 1
		prefix := "  "
		if lineNum == errorLine {
			prefix = "→ "
		}
		sb.WriteString(fmt.Sprintf("%s%4d | %s\n", prefix, lineNum, lines[i]))
	}

	return sb.String()
}

func extractMessage(errMsg string) string {
	parts := strings.SplitN(errMsg, ":", 4)
	if len(parts) >= 4 {
		return strings.TrimSpace(parts[3])
	}
	return errMsg
}

func suggestFix(errMsg string) string {
	errLower := strings.ToLower(errMsg)

	switch {
	case strings.Contains(errLower, "function") && strings.Contains(errLower, "not defined"):
		return "Check that the function name is correct and registered in the template engine"
	case strings.Contains(errLower, "unexpected"):
		return "Check for missing closing tags or mismatched braces"
	case strings.Contains(errLower, "unterminated"):
		return "Check for unclosed template actions {{ }} or comments {{/* */}}"
	case strings.Contains(errLower, "nil pointer"):
		return "Ensure the data passed to the template contains all required fields"
	default:
		return ""
	}
}

type ValidationResult struct {
	Errors   []Error
	Warnings []Error
}

func (vr *ValidationResult) IsValid() bool {
	return len(vr.Errors) == 0
}

func (vr *ValidationResult) Report() string {
	if vr.IsValid() && len(vr.Warnings) == 0 {
		return "✓ All templates are valid\n"
	}

	var sb strings.Builder

	if len(vr.Errors) > 0 {
		sb.WriteString(fmt.Sprintf("✗ Found %d template error(s):\n\n", len(vr.Errors)))
		for i, err := range vr.Errors {
			sb.WriteString(fmt.Sprintf("Error %d:\n%s\n", i+1, err.Error()))
		}
	}

	if len(vr.Warnings) > 0 {
		sb.WriteString(fmt.Sprintf("⚠ Found %d template warning(s):\n\n", len(vr.Warnings)))
		for i, warn := range vr.Warnings {
			sb.WriteString(fmt.Sprintf("Warning %d:\n%s\n", i+1, warn.Error()))
		}
	}

	return sb.String()
}

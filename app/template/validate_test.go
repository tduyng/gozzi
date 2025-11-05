// Tests for template validation with source-level error reporting.
package template

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		contains []string
	}{
		{
			name: "full error with all fields",
			err: &Error{
				Template: "home.html",
				Line:     10,
				Column:   5,
				Snippet:  "  10 | {{ .Title }}",
				Message:  "undefined variable",
				Hint:     "Check variable spelling",
			},
			contains: []string{
				"home.html",
				"line 10",
				"column 5",
				"undefined variable",
				"{{ .Title }}",
				"Hint: Check variable spelling",
			},
		},
		{
			name: "error without column",
			err: &Error{
				Template: "post.html",
				Line:     5,
				Message:  "syntax error",
			},
			contains: []string{
				"post.html",
				"line 5",
				"syntax error",
			},
		},
		{
			name: "error without line number",
			err: &Error{
				Template: "base.html",
				Message:  "parse error",
			},
			contains: []string{
				"base.html",
				"parse error",
			},
		},
		{
			name: "minimal error",
			err: &Error{
				Template: "error.html",
				Message:  "something went wrong",
			},
			contains: []string{
				"error.html",
				"something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("expected error to contain %q, got:\n%s", expected, result)
				}
			}
		})
	}
}

func TestNewValidator(t *testing.T) {
	funcMap := template.FuncMap{
		"test": func() string { return "test" },
	}

	validator := NewValidator("/path/to/templates", funcMap)

	if validator == nil {
		t.Fatal("expected non-nil validator")
	}

	if validator.templatesDir != "/path/to/templates" {
		t.Errorf("expected templatesDir %q, got %q", "/path/to/templates", validator.templatesDir)
	}

	if validator.funcMap == nil {
		t.Error("expected non-nil funcMap")
	}
}

func TestValidator_ValidateTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	funcMap := template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}

	validator := NewValidator(tmpDir, funcMap)

	tests := []struct {
		name          string
		templateName  string
		content       string
		expectError   bool
		errorContains string
	}{
		{
			name:         "valid template",
			templateName: "valid.html",
			content:      "<h1>{{ .Title }}</h1>",
			expectError:  false,
		},
		{
			name:         "valid template with function",
			templateName: "with-func.html",
			content:      "<h1>{{ upper .Title }}</h1>",
			expectError:  false,
		},
		{
			name:          "invalid syntax - unclosed action",
			templateName:  "unclosed.html",
			content:       "<h1>{{ .Title </h1>",
			expectError:   true,
			errorContains: "unexpected",
		},
		{
			name:          "undefined function",
			templateName:  "undefined-func.html",
			content:       "{{ undefined_func .Title }}",
			expectError:   true,
			errorContains: "function \"undefined_func\" not defined",
		},
		{
			name:          "malformed template",
			templateName:  "malformed.html",
			content:       "{{ end }}",
			expectError:   true,
			errorContains: "end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create template file
			path := filepath.Join(tmpDir, tt.templateName)
			if err := os.WriteFile(path, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to create template file: %v", err)
			}

			errors := validator.validateTemplate(tt.templateName, path)

			if tt.expectError {
				if len(errors) == 0 {
					t.Fatal("expected error but got none")
				}

				errorMsg := errors[0].Error()
				if !strings.Contains(strings.ToLower(errorMsg), strings.ToLower(tt.errorContains)) {
					t.Errorf("expected error containing %q, got:\n%s", tt.errorContains, errorMsg)
				}
			} else {
				if len(errors) > 0 {
					t.Errorf("expected no errors, got: %v", errors)
				}
			}
		})
	}
}

func TestValidator_ValidateTemplateFileReadError(t *testing.T) {
	tmpDir := t.TempDir()
	validator := NewValidator(tmpDir, template.FuncMap{})

	// Try to validate non-existent file
	errors := validator.validateTemplate("nonexistent.html", filepath.Join(tmpDir, "nonexistent.html"))

	if len(errors) == 0 {
		t.Fatal("expected error for non-existent file")
	}

	if !strings.Contains(errors[0].Message, "Failed to read template") {
		t.Errorf("expected 'Failed to read template' error, got: %s", errors[0].Message)
	}
}

func TestValidator_Validate(t *testing.T) {
	tmpDir := t.TempDir()

	// Create valid templates
	validTemplates := map[string]string{
		"home.html":          "<h1>{{ .Title }}</h1>",
		"post.html":          "<div>{{ .Content }}</div>",
		"partials/head.html": "<head><title>{{ .Title }}</title></head>",
	}

	// Create invalid template
	invalidTemplate := "<h1>{{ .Title </h1>" // unclosed action

	// Write valid templates
	for name, content := range validTemplates {
		path := filepath.Join(tmpDir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("failed to create directory: %v", err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create template %s: %v", name, err)
		}
	}

	// Write invalid template
	invalidPath := filepath.Join(tmpDir, "invalid.html")
	if err := os.WriteFile(invalidPath, []byte(invalidTemplate), 0644); err != nil {
		t.Fatalf("failed to create invalid template: %v", err)
	}

	// Create non-HTML file (should be skipped)
	txtPath := filepath.Join(tmpDir, "readme.txt")
	if err := os.WriteFile(txtPath, []byte("not a template"), 0644); err != nil {
		t.Fatalf("failed to create txt file: %v", err)
	}

	validator := NewValidator(tmpDir, template.FuncMap{})
	errors := validator.Validate()

	// Should have exactly 1 error (from invalid.html)
	if len(errors) != 1 {
		t.Fatalf("expected 1 error, got %d: %v", len(errors), errors)
	}

	if !strings.Contains(errors[0].Template, "invalid.html") {
		t.Errorf("expected error from invalid.html, got: %s", errors[0].Template)
	}
}

func TestValidator_ValidateWithWalkError(t *testing.T) {
	// Use a non-existent directory
	validator := NewValidator("/nonexistent/directory", template.FuncMap{})
	errors := validator.Validate()

	if len(errors) == 0 {
		t.Fatal("expected error for non-existent directory")
	}

	if !strings.Contains(errors[0].Message, "Failed to walk templates directory") {
		t.Errorf("expected walk error, got: %s", errors[0].Message)
	}
}

func TestExtractMessage(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected string
	}{
		{
			name:     "template error with line and column info",
			errMsg:   "template: home.html:10:5: function \"undefined\" not defined",
			expected: "5: function \"undefined\" not defined",
		},
		{
			name:     "simple error",
			errMsg:   "parse error",
			expected: "parse error",
		},
		{
			name:     "error with line but no column",
			errMsg:   "template: post.html:5: unexpected end",
			expected: "unexpected end",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractMessage(tt.errMsg)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestSuggestFix(t *testing.T) {
	tests := []struct {
		name     string
		errMsg   string
		expected string
	}{
		{
			name:     "undefined function",
			errMsg:   "function \"test\" not defined",
			expected: "Check that the function name is correct and registered in the template engine",
		},
		{
			name:     "unexpected token",
			errMsg:   "unexpected \"}\" in command",
			expected: "Check for missing closing tags or mismatched braces",
		},
		{
			name:     "unterminated action",
			errMsg:   "unterminated quoted string",
			expected: "Check for unclosed template actions {{ }} or comments {{/* */}}",
		},
		{
			name:     "nil pointer",
			errMsg:   "nil pointer evaluating",
			expected: "Ensure the data passed to the template contains all required fields",
		},
		{
			name:     "unknown error",
			errMsg:   "some other error",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := suggestFix(tt.errMsg)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestBuildSnippet(t *testing.T) {
	validator := &Validator{}

	lines := []string{
		"line 1",
		"line 2",
		"line 3",
		"line 4 - error here",
		"line 5",
		"line 6",
		"line 7",
	}

	tests := []struct {
		name      string
		errorLine int
		contains  []string
	}{
		{
			name:      "error in middle",
			errorLine: 4,
			contains:  []string{"line 2", "line 3", "line 4 - error here", "line 5", "→"},
		},
		{
			name:      "error at start",
			errorLine: 1,
			contains:  []string{"line 1", "line 2", "→"},
		},
		{
			name:      "error at end",
			errorLine: 7,
			contains:  []string{"line 5", "line 6", "line 7", "→"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validator.buildSnippet(lines, tt.errorLine)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("expected snippet to contain %q, got:\n%s", expected, result)
				}
			}
		})
	}
}

func TestParseError(t *testing.T) {
	validator := &Validator{}

	content := `<html>
<head>
  <title>{{ .Title }}</title>
</head>
<body>
  {{ undefined_func .Content }}
</body>
</html>`

	// Create a real template error by trying to parse with undefined function
	_, err := template.New("test").Parse(content)
	if err == nil {
		// If no error, create a synthetic one for testing
		err = fmt.Errorf("template: test.html:6:5: function \"undefined_func\" not defined")
	}

	errors := validator.parseError("test.html", content, err)

	if len(errors) == 0 {
		t.Fatal("expected at least one error")
	}

	if errors[0].Template != "test.html" {
		t.Errorf("expected template name 'test.html', got %q", errors[0].Template)
	}

	// Verify the error has a message
	if errors[0].Message == "" {
		t.Error("expected non-empty error message")
	}
}

func TestValidationResult_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		result   *ValidationResult
		expected bool
	}{
		{
			name:     "no errors",
			result:   &ValidationResult{Errors: []Error{}},
			expected: true,
		},
		{
			name: "with errors",
			result: &ValidationResult{
				Errors: []Error{
					{Template: "test.html", Message: "error"},
				},
			},
			expected: false,
		},
		{
			name: "with warnings but no errors",
			result: &ValidationResult{
				Errors:   []Error{},
				Warnings: []Error{{Template: "test.html", Message: "warning"}},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.result.IsValid()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestValidationResult_Report(t *testing.T) {
	tests := []struct {
		name     string
		result   *ValidationResult
		contains []string
	}{
		{
			name:     "all valid",
			result:   &ValidationResult{},
			contains: []string{"✓ All templates are valid"},
		},
		{
			name: "with errors",
			result: &ValidationResult{
				Errors: []Error{
					{Template: "home.html", Message: "syntax error"},
					{Template: "post.html", Message: "parse error"},
				},
			},
			contains: []string{"✗ Found 2 template error(s)", "Error 1:", "Error 2:", "home.html", "post.html"},
		},
		{
			name: "with warnings",
			result: &ValidationResult{
				Warnings: []Error{
					{Template: "old.html", Message: "deprecated function"},
				},
			},
			contains: []string{"⚠ Found 1 template warning(s)", "Warning 1:", "old.html"},
		},
		{
			name: "with both errors and warnings",
			result: &ValidationResult{
				Errors: []Error{
					{Template: "error.html", Message: "error"},
				},
				Warnings: []Error{
					{Template: "warn.html", Message: "warning"},
				},
			},
			contains: []string{"✗ Found 1 template error(s)", "⚠ Found 1 template warning(s)", "error.html", "warn.html"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			report := tt.result.Report()

			for _, expected := range tt.contains {
				if !strings.Contains(report, expected) {
					t.Errorf("expected report to contain %q, got:\n%s", expected, report)
				}
			}
		})
	}
}

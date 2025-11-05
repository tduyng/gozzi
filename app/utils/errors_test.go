package utils

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorCategories(t *testing.T) {
	tests := []struct {
		name     string
		category error
		want     string
	}{
		{"config error", ErrConfig, "configuration error"},
		{"template error", ErrTemplate, "template error"},
		{"content error", ErrContent, "content error"},
		{"filesystem error", ErrFileSystem, "filesystem error"},
		{"server error", ErrServer, "server error"},
		{"parser error", ErrParser, "parser error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.category.Error(); got != tt.want {
				t.Errorf("category.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapWithContext(t *testing.T) {
	originalErr := errors.New("original error")
	ctx := ErrorContext{
		Operation: "test_operation",
		Component: "test_component",
		Path:      "/test/path",
	}

	wrappedErr := WrapWithContext(originalErr, ErrConfig, ctx)

	// Check that error can be unwrapped
	if !errors.Is(wrappedErr, ErrConfig) {
		t.Error("wrapped error should be identifiable as config error")
	}

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("wrapped error should be identifiable as original error")
	}

	// Check error message contains context
	errMsg := wrappedErr.Error()
	if errMsg == "" {
		t.Error("error message should not be empty")
	}

	// Check context can be extracted
	if extractedCtx, ok := GetErrorContext(wrappedErr); ok {
		if extractedCtx.Operation != ctx.Operation {
			t.Errorf("extracted operation = %v, want %v", extractedCtx.Operation, ctx.Operation)
		}
		if extractedCtx.Component != ctx.Component {
			t.Errorf("extracted component = %v, want %v", extractedCtx.Component, ctx.Component)
		}
		if extractedCtx.Path != ctx.Path {
			t.Errorf("extracted path = %v, want %v", extractedCtx.Path, ctx.Path)
		}
	} else {
		t.Error("should be able to extract context from contextual error")
	}
}

func TestConfigError(t *testing.T) {
	originalErr := fmt.Errorf("config file not found")
	configErr := ConfigError("loading", "/path/to/config.toml", originalErr)

	if !IsCategory(configErr, ErrConfig) {
		t.Error("config error should be identifiable as config category")
	}

	if !errors.Is(configErr, originalErr) {
		t.Error("config error should wrap original error")
	}

	ctx, ok := GetErrorContext(configErr)
	if !ok {
		t.Error("should be able to extract context")
	}

	if ctx.Operation != "loading" {
		t.Errorf("operation = %v, want %v", ctx.Operation, "loading")
	}

	if ctx.Component != "config" {
		t.Errorf("component = %v, want %v", ctx.Component, "config")
	}

	if ctx.Path != "/path/to/config.toml" {
		t.Errorf("path = %v, want %v", ctx.Path, "/path/to/config.toml")
	}
}

func TestTemplateError(t *testing.T) {
	originalErr := fmt.Errorf("template syntax error")
	templateErr := TemplateError("parsing", "home.html", originalErr)

	if !IsCategory(templateErr, ErrTemplate) {
		t.Error("template error should be identifiable as template category")
	}

	ctx, ok := GetErrorContext(templateErr)
	if !ok {
		t.Error("should be able to extract context")
	}

	if ctx.Operation != "parsing" {
		t.Errorf("operation = %v, want %v", ctx.Operation, "parsing")
	}

	if ctx.Component != "template" {
		t.Errorf("component = %v, want %v", ctx.Component, "template")
	}
}

func TestContentError(t *testing.T) {
	originalErr := fmt.Errorf("markdown parsing failed")
	contentErr := ContentError("processing", "/content/post.md", originalErr)

	if !IsCategory(contentErr, ErrContent) {
		t.Error("content error should be identifiable as content category")
	}

	ctx, ok := GetErrorContext(contentErr)
	if !ok {
		t.Error("should be able to extract context")
	}

	if ctx.Component != "content" {
		t.Errorf("component = %v, want %v", ctx.Component, "content")
	}
}

func TestFileSystemError(t *testing.T) {
	originalErr := fmt.Errorf("permission denied")
	fsErr := FileSystemError("creating", "/output/index.html", originalErr)

	if !IsCategory(fsErr, ErrFileSystem) {
		t.Error("filesystem error should be identifiable as filesystem category")
	}

	ctx, ok := GetErrorContext(fsErr)
	if !ok {
		t.Error("should be able to extract context")
	}

	if ctx.Component != "filesystem" {
		t.Errorf("component = %v, want %v", ctx.Component, "filesystem")
	}
}

func TestServerError(t *testing.T) {
	originalErr := fmt.Errorf("port already in use")
	serverErr := ServerError("starting", originalErr)

	if !IsCategory(serverErr, ErrServer) {
		t.Error("server error should be identifiable as server category")
	}

	ctx, ok := GetErrorContext(serverErr)
	if !ok {
		t.Error("should be able to extract context")
	}

	if ctx.Component != "server" {
		t.Errorf("component = %v, want %v", ctx.Component, "server")
	}
}

func TestParserError(t *testing.T) {
	originalErr := fmt.Errorf("invalid frontmatter")
	parserErr := ParserError("parsing", "/content/page.md", originalErr)

	if !IsCategory(parserErr, ErrParser) {
		t.Error("parser error should be identifiable as parser category")
	}

	ctx, ok := GetErrorContext(parserErr)
	if !ok {
		t.Error("should be able to extract context")
	}

	if ctx.Component != "parser" {
		t.Errorf("component = %v, want %v", ctx.Component, "parser")
	}
}

func TestIsCategory(t *testing.T) {
	originalErr := errors.New("test error")
	configErr := ConfigError("test", "config.toml", originalErr)

	// Should identify correct category
	if !IsCategory(configErr, ErrConfig) {
		t.Error("should identify config error correctly")
	}

	// Should not identify wrong category
	if IsCategory(configErr, ErrTemplate) {
		t.Error("should not identify config error as template error")
	}
}

func TestGetErrorContext_NonContextualError(t *testing.T) {
	plainErr := errors.New("plain error")

	ctx, ok := GetErrorContext(plainErr)
	if ok {
		t.Error("should not extract context from non-contextual error")
	}

	if ctx.Operation != "" || ctx.Component != "" || ctx.Path != "" {
		t.Error("context should be empty for non-contextual error")
	}
}

func TestWrapWithContext_NilError(t *testing.T) {
	ctx := ErrorContext{Operation: "test", Component: "test"}

	result := WrapWithContext(nil, ErrConfig, ctx)
	if result != nil {
		t.Error("wrapping nil error should return nil")
	}
}

// Package shared provides error handling types and utilities for the gozzi application.
// This package contains common error categories and contextual error wrapping functionality.
package shared

import (
	"errors"
	"fmt"
)

// Error categories for better error handling.
var (
	// ErrConfig indicates configuration-related errors.
	ErrConfig = errors.New("configuration error")

	// ErrTemplate indicates template-related errors.
	ErrTemplate = errors.New("template error")

	// ErrContent indicates content processing errors.
	ErrContent = errors.New("content error")

	// ErrFileSystem indicates file system operation errors.
	ErrFileSystem = errors.New("filesystem error")

	// ErrServer indicates server-related errors.
	ErrServer = errors.New("server error")

	// ErrParser indicates content parsing errors.
	ErrParser = errors.New("parser error")
)

// ErrorContext provides additional context for errors.
type ErrorContext struct {
	Operation string
	Component string
	Path      string
	Details   map[string]any
}

// ContextualError wraps an error with additional context.
type ContextualError struct {
	Err      error
	Category error
	Context  ErrorContext
}

func (e *ContextualError) Error() string {
	if e.Context.Path != "" {
		return fmt.Sprintf("%s in %s [%s]: %v", e.Context.Operation, e.Context.Component, e.Context.Path, e.Err)
	}
	return fmt.Sprintf("%s in %s: %v", e.Context.Operation, e.Context.Component, e.Err)
}

func (e *ContextualError) Unwrap() error {
	return e.Err
}

// Is implements error matching for both category and underlying error.
func (e *ContextualError) Is(target error) bool {
	return errors.Is(e.Category, target) || errors.Is(e.Err, target)
}

// WrapWithContext wraps an error with enhanced context.
func WrapWithContext(err error, category error, ctx ErrorContext) error {
	if err == nil {
		return nil
	}

	return &ContextualError{
		Err:      err,
		Category: category,
		Context:  ctx,
	}
}

// ConfigError creates a configuration error with context.
func ConfigError(operation, path string, err error) error {
	return WrapWithContext(err, ErrConfig, ErrorContext{
		Operation: operation,
		Component: "config",
		Path:      path,
	})
}

// TemplateError creates a template error with context.
func TemplateError(operation, templateName string, err error) error {
	return WrapWithContext(err, ErrTemplate, ErrorContext{
		Operation: operation,
		Component: "template",
		Path:      templateName,
	})
}

// ContentError creates a content error with context.
func ContentError(operation, filePath string, err error) error {
	return WrapWithContext(err, ErrContent, ErrorContext{
		Operation: operation,
		Component: "content",
		Path:      filePath,
	})
}

// FileSystemError creates a filesystem error with context.
func FileSystemError(operation, path string, err error) error {
	return WrapWithContext(err, ErrFileSystem, ErrorContext{
		Operation: operation,
		Component: "filesystem",
		Path:      path,
	})
}

// ServerError creates a server error with context.
func ServerError(operation string, err error) error {
	return WrapWithContext(err, ErrServer, ErrorContext{
		Operation: operation,
		Component: "server",
	})
}

// ParserError creates a parser error with context.
func ParserError(operation, contentPath string, err error) error {
	return WrapWithContext(err, ErrParser, ErrorContext{
		Operation: operation,
		Component: "parser",
		Path:      contentPath,
	})
}

// IsCategory checks if an error belongs to a specific category.
func IsCategory(err error, category error) bool {
	return errors.Is(err, category)
}

// GetErrorContext extracts the ErrorContext from a ContextualError if available.
func GetErrorContext(err error) (ErrorContext, bool) {
	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		return contextualErr.Context, true
	}
	return ErrorContext{}, false
}

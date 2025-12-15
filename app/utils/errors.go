// Package utils provides error handling types and contextual error wrapping.
package utils

import (
	"errors"
	"fmt"
)

var (
	ErrConfig     = errors.New("configuration error")
	ErrTemplate   = errors.New("template error")
	ErrContent    = errors.New("content error")
	ErrFileSystem = errors.New("filesystem error")
	ErrServer     = errors.New("server error")
	ErrParser     = errors.New("parser error")
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

func ConfigError(operation, path string, err error) error {
	return WrapWithContext(err, ErrConfig, ErrorContext{
		Operation: operation,
		Component: "config",
		Path:      path,
	})
}

func TemplateError(operation, templateName string, err error) error {
	return WrapWithContext(err, ErrTemplate, ErrorContext{
		Operation: operation,
		Component: "template",
		Path:      templateName,
	})
}

func ContentError(operation, filePath string, err error) error {
	return WrapWithContext(err, ErrContent, ErrorContext{
		Operation: operation,
		Component: "content",
		Path:      filePath,
	})
}

func FileSystemError(operation, path string, err error) error {
	return WrapWithContext(err, ErrFileSystem, ErrorContext{
		Operation: operation,
		Component: "filesystem",
		Path:      path,
	})
}

func ServerError(operation string, err error) error {
	return WrapWithContext(err, ErrServer, ErrorContext{
		Operation: operation,
		Component: "server",
	})
}

func ParserError(operation, contentPath string, err error) error {
	return WrapWithContext(err, ErrParser, ErrorContext{
		Operation: operation,
		Component: "parser",
		Path:      contentPath,
	})
}

func IsCategory(err error, category error) bool {
	return errors.Is(err, category)
}

func GetErrorContext(err error) (ErrorContext, bool) {
	var contextualErr *ContextualError
	if errors.As(err, &contextualErr) {
		return contextualErr.Context, true
	}
	return ErrorContext{}, false
}

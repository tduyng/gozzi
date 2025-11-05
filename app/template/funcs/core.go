// Core template functions for basic operations like arithmetic, logic, and comparisons.
package funcs

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

// Core arithmetic operations

// Add adds two numbers (int or float64).
func Add(a, b any) (any, error) {
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			return a + b, nil
		case float64:
			return float64(a) + b, nil
		default:
			return nil, fmt.Errorf("add: second argument must be int or float64, got %T", b)
		}
	case float64:
		switch b := b.(type) {
		case int:
			return a + float64(b), nil
		case float64:
			return a + b, nil
		default:
			return nil, fmt.Errorf("add: second argument must be int or float64, got %T", b)
		}
	default:
		return nil, fmt.Errorf("add: first argument must be int or float64, got %T", a)
	}
}

// Sub subtracts b from a.
func Sub(a, b int) int {
	return a - b
}

// Core comparison operations

// Eq checks if two values are equal.
func Eq(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

// Ne checks if two values are not equal.
func Ne(a, b any) bool {
	return !Eq(a, b)
}

// Core logic operations

// And returns true if all values are truthy.
func And(values ...any) bool {
	for _, v := range values {
		if !isTruthy(v) {
			return false
		}
	}
	return true
}

// Or returns true if any value is truthy.
func Or(values ...any) bool {
	return slices.ContainsFunc(values, isTruthy)
}

// isTruthy determines if a value is truthy.
func isTruthy(v any) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case int:
		return val != 0
	case string:
		return val != ""
	case []any:
		return len(val) > 0
	default:
		// Use reflection for slices and maps
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			return rv.Len() > 0
		default:
			return true
		}
	}
}

// Collection operations

// First returns the first element of a collection.
func First(items any) (any, error) {
	if items == nil {
		return nil, fmt.Errorf("first: cannot get first element of nil")
	}

	rv := reflect.ValueOf(items)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		if rv.Len() == 0 {
			return nil, fmt.Errorf("first: collection is empty")
		}
		return rv.Index(0).Interface(), nil
	default:
		return nil, fmt.Errorf("first: argument must be a slice or array, got %T", items)
	}
}

// Last returns the last element of a collection.
func Last(items any) (any, error) {
	if items == nil {
		return nil, fmt.Errorf("last: cannot get last element of nil")
	}

	rv := reflect.ValueOf(items)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		if rv.Len() == 0 {
			return nil, fmt.Errorf("last: collection is empty")
		}
		return rv.Index(rv.Len() - 1).Interface(), nil
	default:
		return nil, fmt.Errorf("last: argument must be a slice or array, got %T", items)
	}
}

// Contains checks if a collection contains a value.
func Contains(haystack, needle any) (bool, error) {
	if haystack == nil {
		return false, fmt.Errorf("contains: haystack cannot be nil")
	}

	switch h := haystack.(type) {
	case string:
		n := fmt.Sprintf("%v", needle)
		return strings.Contains(h, n), nil
	case []any:
		for _, item := range h {
			if Eq(item, needle) {
				return true, nil
			}
		}
		return false, nil
	default:
		// Use reflection for other slice types
		rv := reflect.ValueOf(haystack)
		if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
			return false, fmt.Errorf("contains: haystack must be string, slice, or array, got %T", haystack)
		}
		for i := 0; i < rv.Len(); i++ {
			if Eq(rv.Index(i).Interface(), needle) {
				return true, nil
			}
		}
		return false, nil
	}
}

// String operations

// Default returns the value if not empty, otherwise returns the default.
func Default(val any, def string) string {
	if val == nil {
		return def
	}
	if s, ok := val.(string); ok {
		if s == "" {
			return def
		}
		return s
	}
	s := fmt.Sprint(val)
	if s == "" || s == "<nil>" {
		return def
	}
	return s
}

// Priority returns the first non-empty value from the arguments.
func Priority(vals ...any) string {
	for _, v := range vals {
		if v == nil {
			continue
		}
		if s, ok := v.(string); ok && s != "" {
			return s
		}
		s := fmt.Sprint(v)
		if s != "" && s != "<nil>" {
			return s
		}
	}
	return ""
}

// Pluralize adds 's' to singular if count is not 1.
func Pluralize(singular string, count any) (string, error) {
	var c int
	switch v := count.(type) {
	case int:
		c = v
	case int64:
		c = int(v)
	default:
		return "", fmt.Errorf("pluralize: count must be int or int64, got %T", count)
	}

	if c == 1 {
		return singular, nil
	}
	return singular + "s", nil
}

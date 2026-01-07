// Package funcs provides core template functions for arithmetic, logic, and comparisons.
package funcs

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
)

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

func Sub(a, b int) int {
	return a - b
}

func Eq(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func Ne(a, b any) bool {
	return !Eq(a, b)
}

// Gt returns true if a > b
func Gt(a, b any) (bool, error) {
	return compare(a, b, func(cmp int) bool { return cmp > 0 })
}

// Ge returns true if a >= b
func Ge(a, b any) (bool, error) {
	return compare(a, b, func(cmp int) bool { return cmp >= 0 })
}

// Lt returns true if a < b
func Lt(a, b any) (bool, error) {
	return compare(a, b, func(cmp int) bool { return cmp < 0 })
}

// Le returns true if a <= b
func Le(a, b any) (bool, error) {
	return compare(a, b, func(cmp int) bool { return cmp <= 0 })
}

func compare(a, b any, op func(int) bool) (bool, error) {
	aVal, aOk := toNumber(a)
	bVal, bOk := toNumber(b)

	if !aOk {
		return false, fmt.Errorf("cannot compare: first argument is not a number (%T)", a)
	}
	if !bOk {
		return false, fmt.Errorf("cannot compare: second argument is not a number (%T)", b)
	}

	if aVal < bVal {
		return op(-1), nil
	} else if aVal > bVal {
		return op(1), nil
	}
	return op(0), nil
}

func toNumber(v any) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case int32:
		return float64(val), true
	case float64:
		return val, true
	case float32:
		return float64(val), true
	default:
		return 0, false
	}
}

func And(values ...any) bool {
	for _, v := range values {
		if !isTruthy(v) {
			return false
		}
	}
	return true
}

func Or(values ...any) bool {
	return slices.ContainsFunc(values, isTruthy)
}

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
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			return rv.Len() > 0
		default:
			return true
		}
	}
}

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

// In is an alias for Contains with reversed argument order for more readable syntax.
// Usage: {{ if in "golang" .Tags }}...{{ end }}
// This is equivalent to: {{ if contains .Tags "golang" }}...{{ end }}
func In(needle, haystack any) (bool, error) {
	return Contains(haystack, needle)
}

// Default returns the value if it's non-empty, otherwise returns the default.
// Now supports any type for both val and def (improved from string-only).
func Default(def, val any) any {
	if val == nil {
		return def
	}

	// Handle different types
	switch v := val.(type) {
	case string:
		if v == "" {
			return def
		}
		return v
	case int, int8, int16, int32, int64:
		if fmt.Sprint(v) == "0" {
			return def
		}
		return v
	case float32, float64:
		if fmt.Sprint(v) == "0" {
			return def
		}
		return v
	case bool:
		if !v {
			return def
		}
		return v
	default:
		// For slices, arrays, maps - check if empty using reflection
		rv := reflect.ValueOf(val)
		switch rv.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			if rv.Len() == 0 {
				return def
			}
			return v
		default:
			// For other types, return the value as-is
			return v
		}
	}
}

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

// Len returns the length of a slice, array, map, or string
func Len(v any) (int, error) {
	if v == nil {
		return 0, fmt.Errorf("len: cannot get length of nil")
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map, reflect.String:
		return rv.Len(), nil
	default:
		return 0, fmt.Errorf("len: argument must be slice, array, map, or string, got %T", v)
	}
}

// Slice creates a slice from the given arguments
func Slice(items ...any) []any {
	return items
}

// Cond is a ternary conditional operator: returns trueVal if condition is true, otherwise falseVal
func Cond(condition bool, trueVal, falseVal any) any {
	if condition {
		return trueVal
	}
	return falseVal
}

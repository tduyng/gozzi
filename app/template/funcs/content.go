// Package funcs provides content-specific template functions for filtering and grouping.
package funcs

import (
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/tduyng/gozzi/app/content"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Limit returns at most n items from the collection.
func Limit(max any, items any) (any, error) {
	if items == nil {
		return nil, nil
	}

	var n int
	switch v := max.(type) {
	case int:
		n = v
	case int64:
		n = int(v)
	case float64:
		n = int(v)
	default:
		return nil, fmt.Errorf("limit: first argument must be an integer, got %T", max)
	}

	if n < 0 {
		return nil, fmt.Errorf("limit: count cannot be negative")
	}

	rv := reflect.ValueOf(items)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("limit: second argument must be a slice or array, got %T", items)
	}

	l := rv.Len()
	if n > l {
		n = l
	}

	return rv.Slice(0, n).Interface(), nil
}

// Reverse returns a copy of the slice with elements in reverse order.
// Supports both []*content.Node and generic []any.
func Reverse(items any) any {
	if items == nil {
		return nil
	}

	rv := reflect.ValueOf(items)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return items
	}

	l := rv.Len()
	result := reflect.MakeSlice(rv.Type(), l, l)
	for i := 0; i < l; i++ {
		result.Index(i).Set(rv.Index(l - 1 - i))
	}

	return result.Interface()
}

func extractTime(v any) time.Time {
	if v == nil {
		return time.Time{}
	}
	switch t := v.(type) {
	case time.Time:
		return t
	case string:
		// Try parsing common formats
		if parsed, err := time.Parse(time.RFC3339, t); err == nil {
			return parsed
		}
		if parsed, err := time.Parse("2006-01-02", t); err == nil {
			return parsed
		}
	}
	return time.Time{} // Zero time
}

// Where filters a slice of items by a field value.
// Supports both map[string]any and *content.Node.
func Where(items any, field string, value any) (any, error) {
	if field == "" {
		return nil, fmt.Errorf("where: field name cannot be empty")
	}

	rv := reflect.ValueOf(items)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("where: first argument must be a slice or array, got %T", items)
	}

	result := reflect.MakeSlice(rv.Type(), 0, rv.Len())

	for i := 0; i < rv.Len(); i++ {
		item := rv.Index(i).Interface()
		var fieldValue any
		var found bool

		// Check if it's a pointer to content.Node
		if node, ok := item.(*content.Node); ok {
			// Check Config map first (front matter)
			if val, exists := node.Config[field]; exists {
				fieldValue = val
				found = true
			} else {
				// Fallback to Node struct fields using reflection
				nodeRv := reflect.ValueOf(node).Elem()
				// Try case-insensitive field name match for struct fields
				fieldName := cases.Title(language.English).String(field)
				fieldVal := nodeRv.FieldByName(fieldName)
				if fieldVal.IsValid() {
					fieldValue = fieldVal.Interface()
					found = true
				}
			}
		} else {
			// Try as map[string]any
			itemRv := reflect.ValueOf(item)
			if itemRv.Kind() == reflect.Map {
				mapVal := itemRv.MapIndex(reflect.ValueOf(field))
				if mapVal.IsValid() {
					fieldValue = mapVal.Interface()
					found = true
				}
			} else {
				// Try reflection for other struct types
				if itemRv.Kind() == reflect.Ptr {
					itemRv = itemRv.Elem()
				}
				if itemRv.Kind() == reflect.Struct {
					fieldVal := itemRv.FieldByName(cases.Title(language.English).String(field))
					if fieldVal.IsValid() {
						fieldValue = fieldVal.Interface()
						found = true
					}
				}
			}
		}

		if found && Eq(fieldValue, value) {
			result = reflect.Append(result, rv.Index(i))
		}
	}

	return result.Interface(), nil
}

// Group represents a grouped collection of content nodes.
type Group struct {
	Key   string
	Items []*content.Node
}

// GroupBy groups content nodes by a given criteria (year, month, day).
func GroupBy(criteria string, nodes []*content.Node) ([]Group, error) {
	groupsMap := make(map[string][]*content.Node)
	var keys []string

	for _, node := range nodes {
		dateVal, exists := node.Config["date"]
		if !exists {
			continue
		}

		t := extractTime(dateVal)
		if t.IsZero() {
			continue
		}

		var key string
		switch criteria {
		case "year":
			key = t.Format("2006")
		case "month":
			key = t.Format("2006-01")
		case "day":
			key = t.Format("2006-01-02")
		default:
			return nil, fmt.Errorf("invalid grouping criteria: %s", criteria)
		}

		if _, exists := groupsMap[key]; !exists {
			keys = append(keys, key)
		}
		groupsMap[key] = append(groupsMap[key], node)
	}

	// Sort keys descending
	slices.Sort(keys)
	slices.Reverse(keys)

	var groups []Group
	for _, key := range keys {
		groups = append(groups, Group{
			Key:   key,
			Items: groupsMap[key],
		})
	}

	return groups, nil
}

// Concat merges multiple slices of nodes into one.
func Concat(slices ...[]*content.Node) []*content.Node {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}

	result := make([]*content.Node, 0, totalLen)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// SortBy sorts a slice of nodes by a given field (only "date" currently supported).
func SortBy(field string, nodes []*content.Node) ([]*content.Node, error) {
	if field != "date" {
		return nil, fmt.Errorf("sort_by: only 'date' field is supported currently")
	}

	// Create a copy to avoid modifying original
	sorted := make([]*content.Node, len(nodes))
	copy(sorted, nodes)

	slices.SortFunc(sorted, func(a, b *content.Node) int {
		timeA := extractTime(a.Config["date"])
		timeB := extractTime(b.Config["date"])

		if timeA.IsZero() && timeB.IsZero() {
			return 0
		}
		if timeA.IsZero() {
			return 1 // a is less (zero), so it comes after b
		}
		if timeB.IsZero() {
			return -1 // b is less (zero), so it comes after a
		}

		// Descending order (newest first)
		if timeA.After(timeB) {
			return -1
		}
		if timeA.Before(timeB) {
			return 1
		}
		return 0
	})

	return sorted, nil
}

// RelatedPosts finds posts related to the given page based on tag overlap.
func RelatedPosts(pageData any, allPosts []*content.Node) []*content.Node {
	// Convert template map to Node if necessary
	var page *content.Node

	switch v := pageData.(type) {
	case *content.Node:
		page = v
	case map[string]any:
		// Try to find existing node first
		if perm, ok := v["Permalink"].(string); ok && perm != "" {
			for _, n := range allPosts {
				if n.Permalink == perm {
					page = n
					break
				}
			}
		}

		if page == nil {
			// Extract Config FIRST, before creating the node
			// This is critical because template context has structure:
			// { "Config": {...}, "Permalink": "...", "Series": ... }
			var cfg map[string]any
			var permalink string

			if c, ok := v["Config"].(map[string]any); ok {
				// Nested Config field exists (template context)
				cfg = c
			} else {
				// Fallback: treat entire map as config (legacy/direct usage)
				cfg = v
			}

			if pl, ok := v["Permalink"].(string); ok {
				permalink = pl
			}

			page = &content.Node{
				Config:    cfg,
				Permalink: permalink,
			}
		}
	default:
		return []*content.Node{} // Unsupported type
	}

	config := content.DefaultRelatedConfig()
	// Return 6 candidates for client randomization (will show 3)
	config.ResultLimit = 6
	config.MaxCandidates = 10
	// Disable recency penalty for tag-based related posts
	// Content relevance (tags) is more important than publication date
	config.RecencyDecayDays = 0

	finder := content.NewRelatedPostsFinder(allPosts, config)
	return finder.FindRelated(page)
}

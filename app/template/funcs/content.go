// Content-specific template functions for filtering, grouping, and manipulating content nodes.
package funcs

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/tduyng/gozzi/app/content"
)

// Limit limits the number of items returned from a slice.
func Limit(maxItems any, items []*content.Node) ([]*content.Node, error) {
	var m int
	switch v := maxItems.(type) {
	case int:
		m = v
	case int64:
		m = int(v)
	default:
		return nil, fmt.Errorf("limit: maxItems must be int or int64, got %T", maxItems)
	}

	if m < 0 {
		return nil, fmt.Errorf("limit: maxItems must be non-negative, got %d", m)
	}

	if m > len(items) {
		m = len(items)
	}
	return items[:m], nil
}

// Reverse reverses the order of content nodes.
func Reverse(items []*content.Node) []*content.Node {
	reversed := make([]*content.Node, len(items))
	for i := range items {
		reversed[len(items)-1-i] = items[i]
	}
	return reversed
}

// Where filters a slice of maps by a field value.
func Where(sections []any, field string, value any) ([]any, error) {
	if field == "" {
		return nil, fmt.Errorf("where: field name cannot be empty")
	}

	var result []any
	for _, s := range sections {
		sectionMap, ok := s.(map[string]any)
		if !ok {
			continue
		}

		fieldValue, exists := sectionMap[field]
		if exists && Eq(fieldValue, value) {
			result = append(result, s)
		}
	}
	return result, nil
}

// Group represents a grouped collection of content nodes.
type Group struct {
	Key   string
	Items []*content.Node
}

// GroupBy groups content nodes by a date field (year, month, or day).
func GroupBy(key string, nodes []*content.Node) ([]Group, error) {
	if key != "year" && key != "month" && key != "day" {
		return nil, fmt.Errorf("group_by: key must be 'year', 'month', or 'day', got %q", key)
	}

	groups := make(map[string][]*content.Node)

	for _, node := range nodes {
		// Extract date from node config
		dateVal, ok := node.Config["date"]
		if !ok {
			continue // Skip nodes without date
		}

		var t time.Time
		switch v := dateVal.(type) {
		case time.Time:
			t = v
		case string:
			parsed, err := time.Parse(time.RFC3339, v)
			if err != nil {
				continue // Skip invalid dates
			}
			t = parsed
		default:
			continue // Skip unsupported date types
		}

		// Get grouping key
		var groupKey string
		switch key {
		case "year":
			groupKey = strconv.Itoa(t.Year())
		case "month":
			groupKey = t.Format("2006-01")
		case "day":
			groupKey = t.Format("2006-01-02")
		}

		groups[groupKey] = append(groups[groupKey], node)
	}

	// Convert to sorted slice
	var result []Group
	for k, items := range groups {
		result = append(result, Group{Key: k, Items: items})
	}

	// Sort descending chronological order
	slices.SortFunc(result, func(a, b Group) int {
		return strings.Compare(b.Key, a.Key) // Reverse for descending
	})

	return result, nil
}

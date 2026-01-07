// Package funcs provides content-specific template functions for filtering and grouping.
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

// Concat merges multiple slices of content nodes into a single slice.
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

// SortBy sorts content nodes by a field in descending order (newest first).
// Currently only supports sorting by "date" field.
func SortBy(field string, nodes []*content.Node) ([]*content.Node, error) {
	if field != "date" {
		return nil, fmt.Errorf("sort_by: currently only 'date' field is supported, got %q", field)
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]*content.Node, len(nodes))
	copy(sorted, nodes)

	// Sort by date in descending order (newest first)
	slices.SortFunc(sorted, func(a, b *content.Node) int {
		aDate := extractTime(a.Config["date"])
		bDate := extractTime(b.Config["date"])

		// Handle nil dates (push to end)
		if aDate.IsZero() && bDate.IsZero() {
			return 0
		}
		if aDate.IsZero() {
			return 1 // a goes after b
		}
		if bDate.IsZero() {
			return -1 // b goes after a
		}

		// Descending order: newer dates first
		return bDate.Compare(aDate)
	})

	return sorted, nil
}

// extractTime extracts a time.Time from various date formats.
func extractTime(dateVal any) time.Time {
	switch v := dateVal.(type) {
	case time.Time:
		return v
	case string:
		// Try RFC3339 format
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t
		}
		// Try date-only format
		if t, err := time.Parse("2006-01-02", v); err == nil {
			return t
		}
	}
	return time.Time{} // Zero time
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

// RelatedPosts finds related posts using intelligent tag-based scoring with randomization.
// Returns up to 6 candidates from top 10 matches for client-side random selection.
func RelatedPosts(pageData any, allPosts []*content.Node) []*content.Node {
	// Convert template map to Node if necessary
	var page *content.Node

	switch v := pageData.(type) {
	case *content.Node:
		page = v
	case map[string]any:
		// Extract Config FIRST, before creating the node
		// This is critical because template context has structure:
		// { "Config": {...}, "Permalink": "...", "Series": ... }
		var config map[string]any
		var permalink string

		if cfg, ok := v["Config"].(map[string]any); ok {
			// Nested Config field exists (template context)
			config = cfg
		} else {
			// Fallback: treat entire map as config (legacy/direct usage)
			config = v
		}

		if pl, ok := v["Permalink"].(string); ok {
			permalink = pl
		}

		page = &content.Node{
			Config:    config,
			Permalink: permalink,
		}
	default:
		return []*content.Node{} // Unsupported type
	}

	config := content.DefaultRelatedConfig()
	// Return 6 candidates for client randomization (will show 3)
	config.ResultLimit = 6
	config.MaxCandidates = 10

	finder := content.NewRelatedPostsFinder(allPosts, config)
	return finder.FindRelated(page)
}

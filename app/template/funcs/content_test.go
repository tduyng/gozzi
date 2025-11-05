// ABOUTME: Tests for content-specific template functions for filtering, grouping, and manipulating content nodes.
// ABOUTME: Validates Limit, Reverse, Where, and GroupBy functions.
package funcs

import (
	"reflect"
	"testing"
	"time"

	"github.com/tduyng/gozzi/app/content"
)

func TestLimit(t *testing.T) {
	nodes := []*content.Node{
		{Slug: "post1"},
		{Slug: "post2"},
		{Slug: "post3"},
		{Slug: "post4"},
		{Slug: "post5"},
	}

	tests := []struct {
		name     string
		maxItems any
		items    []*content.Node
		want     int
		wantErr  bool
	}{
		{"limit 3", 3, nodes, 3, false},
		{"limit all", 5, nodes, 5, false},
		{"limit more than available", 10, nodes, 5, false},
		{"limit 0", 0, nodes, 0, false},
		{"int64", int64(2), nodes, 2, false},
		{"negative limit", -1, nodes, 0, true},
		{"invalid type", "three", nodes, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Limit(tt.maxItems, tt.items)
			if (err != nil) != tt.wantErr {
				t.Errorf("Limit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.want {
				t.Errorf("Limit() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}

func TestReverse(t *testing.T) {
	nodes := []*content.Node{
		{Slug: "post1"},
		{Slug: "post2"},
		{Slug: "post3"},
	}

	reversed := Reverse(nodes)

	// Check length
	if len(reversed) != len(nodes) {
		t.Errorf("Reverse() length = %d, want %d", len(reversed), len(nodes))
	}

	// Check order is reversed
	for i := range nodes {
		if reversed[i].Slug != nodes[len(nodes)-1-i].Slug {
			t.Errorf("Reverse()[%d].Slug = %s, want %s",
				i, reversed[i].Slug, nodes[len(nodes)-1-i].Slug)
		}
	}

	// Ensure original is not modified
	if nodes[0].Slug != "post1" {
		t.Error("Reverse() modified original slice")
	}
}

func TestReverseEmpty(t *testing.T) {
	nodes := []*content.Node{}
	reversed := Reverse(nodes)

	if len(reversed) != 0 {
		t.Errorf("Reverse() of empty slice = %d items, want 0", len(reversed))
	}
}

func TestWhere(t *testing.T) {
	sections := []any{
		map[string]any{"name": "blog", "type": "post"},
		map[string]any{"name": "about", "type": "page"},
		map[string]any{"name": "notes", "type": "post"},
		map[string]any{"name": "contact", "type": "page"},
	}

	tests := []struct {
		name    string
		field   string
		value   any
		want    int
		wantErr bool
	}{
		{"filter by type post", "type", "post", 2, false},
		{"filter by type page", "type", "page", 2, false},
		{"filter by name", "name", "blog", 1, false},
		{"no matches", "type", "article", 0, false},
		{"empty field", "", "post", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Where(sections, tt.field, tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("Where() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != tt.want {
				t.Errorf("Where() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}

func TestWhereInvalidItems(t *testing.T) {
	sections := []any{
		"not a map",
		123,
		map[string]any{"name": "valid"},
	}

	got, err := Where(sections, "name", "valid")
	if err != nil {
		t.Errorf("Where() unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("Where() = %d items, want 1 (should skip invalid items)", len(got))
	}
}

func TestGroupBy(t *testing.T) {
	nodes := []*content.Node{
		{
			Slug:   "post1",
			Config: map[string]any{"date": time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		},
		{
			Slug:   "post2",
			Config: map[string]any{"date": time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)},
		},
		{
			Slug:   "post3",
			Config: map[string]any{"date": time.Date(2023, 12, 5, 0, 0, 0, 0, time.UTC)},
		},
		{
			Slug:   "post4",
			Config: map[string]any{"date": time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC)},
		},
	}

	t.Run("group by year", func(t *testing.T) {
		groups, err := GroupBy("year", nodes)
		if err != nil {
			t.Fatalf("GroupBy() error = %v", err)
		}

		if len(groups) != 2 {
			t.Fatalf("GroupBy() returned %d groups, want 2", len(groups))
		}

		// Should be sorted descending (2024, 2023)
		if groups[0].Key != "2024" {
			t.Errorf("First group key = %s, want 2024", groups[0].Key)
		}
		if groups[1].Key != "2023" {
			t.Errorf("Second group key = %s, want 2023", groups[1].Key)
		}

		// 2024 should have 3 posts
		if len(groups[0].Items) != 3 {
			t.Errorf("2024 group has %d items, want 3", len(groups[0].Items))
		}
		// 2023 should have 1 post
		if len(groups[1].Items) != 1 {
			t.Errorf("2023 group has %d items, want 1", len(groups[1].Items))
		}
	})

	t.Run("group by month", func(t *testing.T) {
		groups, err := GroupBy("month", nodes)
		if err != nil {
			t.Fatalf("GroupBy() error = %v", err)
		}

		if len(groups) != 3 {
			t.Fatalf("GroupBy() returned %d groups, want 3", len(groups))
		}

		// Should be sorted descending
		expectedKeys := []string{"2024-02", "2024-01", "2023-12"}
		for i, expected := range expectedKeys {
			if groups[i].Key != expected {
				t.Errorf("Group %d key = %s, want %s", i, groups[i].Key, expected)
			}
		}
	})

	t.Run("group by day", func(t *testing.T) {
		groups, err := GroupBy("day", nodes)
		if err != nil {
			t.Fatalf("GroupBy() error = %v", err)
		}

		if len(groups) != 4 {
			t.Fatalf("GroupBy() returned %d groups, want 4", len(groups))
		}

		// Each group should have 1 item
		for i, g := range groups {
			if len(g.Items) != 1 {
				t.Errorf("Group %d has %d items, want 1", i, len(g.Items))
			}
		}
	})

	t.Run("invalid key", func(t *testing.T) {
		_, err := GroupBy("hour", nodes)
		if err == nil {
			t.Error("GroupBy() expected error for invalid key")
		}
	})
}

func TestGroupByStringDates(t *testing.T) {
	nodes := []*content.Node{
		{
			Slug:   "post1",
			Config: map[string]any{"date": "2024-01-15T10:00:00Z"},
		},
		{
			Slug:   "post2",
			Config: map[string]any{"date": "2024-01-20T10:00:00Z"},
		},
	}

	groups, err := GroupBy("year", nodes)
	if err != nil {
		t.Fatalf("GroupBy() error = %v", err)
	}

	if len(groups) != 1 {
		t.Fatalf("GroupBy() returned %d groups, want 1", len(groups))
	}

	if groups[0].Key != "2024" {
		t.Errorf("Group key = %s, want 2024", groups[0].Key)
	}

	if len(groups[0].Items) != 2 {
		t.Errorf("Group has %d items, want 2", len(groups[0].Items))
	}
}

func TestGroupByNoDate(t *testing.T) {
	nodes := []*content.Node{
		{
			Slug:   "post1",
			Config: map[string]any{"title": "Post 1"}, // No date field
		},
		{
			Slug:   "post2",
			Config: map[string]any{"date": time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		},
	}

	groups, err := GroupBy("year", nodes)
	if err != nil {
		t.Fatalf("GroupBy() error = %v", err)
	}

	// Should only have one group (for the node with date)
	if len(groups) != 1 {
		t.Fatalf("GroupBy() returned %d groups, want 1", len(groups))
	}

	if len(groups[0].Items) != 1 {
		t.Errorf("Group has %d items, want 1", len(groups[0].Items))
	}
}

func TestGroupType(t *testing.T) {
	// Test the Group struct itself
	group := Group{
		Key: "2024",
		Items: []*content.Node{
			{Slug: "post1"},
			{Slug: "post2"},
		},
	}

	if group.Key != "2024" {
		t.Errorf("Group.Key = %s, want 2024", group.Key)
	}

	if len(group.Items) != 2 {
		t.Errorf("Group.Items length = %d, want 2", len(group.Items))
	}
}

func TestEqHelper(t *testing.T) {
	// This test verifies Eq is used correctly by Where
	// We're testing that Where uses Eq for comparisons
	sections := []any{
		map[string]any{"count": 5},
		map[string]any{"count": "5"}, // String "5"
	}

	// Eq treats 5 and "5" as equal (via string representation)
	got, err := Where(sections, "count", 5)
	if err != nil {
		t.Fatalf("Where() error = %v", err)
	}

	// Should match both because Eq compares string representations
	if len(got) != 2 {
		t.Errorf("Where() with Eq comparison returned %d items, want 2", len(got))
	}
}

func TestReverseNilSafety(t *testing.T) {
	// Ensure Reverse handles edge cases gracefully
	var nilSlice []*content.Node
	reversed := Reverse(nilSlice)

	if reversed == nil {
		t.Error("Reverse() returned nil, expected empty slice")
	}

	if !reflect.DeepEqual(reversed, []*content.Node{}) {
		t.Errorf("Reverse() of nil = %v, want empty slice", reversed)
	}
}

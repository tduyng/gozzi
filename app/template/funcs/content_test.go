// Tests for content-specific template functions for filtering, grouping, and manipulating content nodes.
// Validates Limit, Reverse, Where, and GroupBy functions.
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

func TestConcat(t *testing.T) {
	slice1 := []*content.Node{
		{Slug: "post1"},
		{Slug: "post2"},
	}
	slice2 := []*content.Node{
		{Slug: "note1"},
		{Slug: "note2"},
		{Slug: "note3"},
	}
	slice3 := []*content.Node{
		{Slug: "page1"},
	}

	tests := []struct {
		name   string
		slices [][]*content.Node
		want   int
	}{
		{"concat two slices", [][]*content.Node{slice1, slice2}, 5},
		{"concat three slices", [][]*content.Node{slice1, slice2, slice3}, 6},
		{"concat one slice", [][]*content.Node{slice1}, 2},
		{"concat empty slices", [][]*content.Node{{}, {}}, 0},
		{"concat no slices", [][]*content.Node{}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Concat(tt.slices...)
			if len(got) != tt.want {
				t.Errorf("Concat() returned %d items, want %d", len(got), tt.want)
			}
		})
	}
}

func TestConcatOrder(t *testing.T) {
	slice1 := []*content.Node{{Slug: "a"}, {Slug: "b"}}
	slice2 := []*content.Node{{Slug: "c"}, {Slug: "d"}}

	result := Concat(slice1, slice2)

	expected := []string{"a", "b", "c", "d"}
	for i, slug := range expected {
		if result[i].Slug != slug {
			t.Errorf("Concat()[%d].Slug = %s, want %s", i, result[i].Slug, slug)
		}
	}
}

func TestConcatDoesNotModifyOriginals(t *testing.T) {
	slice1 := []*content.Node{{Slug: "post1"}}
	slice2 := []*content.Node{{Slug: "note1"}}
	originalLen1 := len(slice1)
	originalLen2 := len(slice2)

	result := Concat(slice1, slice2)

	// Modify result slice structure (not the nodes themselves)
	_ = append(result, &content.Node{Slug: "new"})

	// Check original slice structures are unmodified
	if len(slice1) != originalLen1 {
		t.Errorf("Concat() modified original slice1 length: got %d, want %d", len(slice1), originalLen1)
	}
	if len(slice2) != originalLen2 {
		t.Errorf("Concat() modified original slice2 length: got %d, want %d", len(slice2), originalLen2)
	}

	// Note: The nodes themselves are pointers, so they're shared.
	// This is expected behavior - we're only ensuring the slice structure isn't modified.
}

func TestSortBy(t *testing.T) {
	nodes := []*content.Node{
		{
			Slug:   "post1",
			Config: map[string]any{"date": time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		},
		{
			Slug:   "post2",
			Config: map[string]any{"date": time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC)},
		},
		{
			Slug:   "post3",
			Config: map[string]any{"date": time.Date(2023, 12, 5, 0, 0, 0, 0, time.UTC)},
		},
		{
			Slug:   "post4",
			Config: map[string]any{"date": time.Date(2024, 8, 10, 0, 0, 0, 0, time.UTC)},
		},
	}

	t.Run("sort by date descending", func(t *testing.T) {
		sorted, err := SortBy("date", nodes)
		if err != nil {
			t.Fatalf("SortBy() error = %v", err)
		}

		if len(sorted) != 4 {
			t.Fatalf("SortBy() returned %d items, want 4", len(sorted))
		}

		// Should be sorted descending (newest first)
		expectedOrder := []string{"post2", "post4", "post1", "post3"}
		for i, expectedSlug := range expectedOrder {
			if sorted[i].Slug != expectedSlug {
				t.Errorf("SortBy()[%d].Slug = %s, want %s", i, sorted[i].Slug, expectedSlug)
			}
		}
	})

	t.Run("invalid field", func(t *testing.T) {
		_, err := SortBy("title", nodes)
		if err == nil {
			t.Error("SortBy() expected error for invalid field")
		}
	})

	t.Run("does not modify original", func(t *testing.T) {
		original := []*content.Node{
			{Slug: "a", Config: map[string]any{"date": time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)}},
			{Slug: "b", Config: map[string]any{"date": time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)}},
		}
		originalFirst := original[0].Slug

		_, err := SortBy("date", original)
		if err != nil {
			t.Fatalf("SortBy() error = %v", err)
		}

		if original[0].Slug != originalFirst {
			t.Error("SortBy() modified original slice")
		}
	})
}

func TestSortByWithStringDates(t *testing.T) {
	nodes := []*content.Node{
		{
			Slug:   "post1",
			Config: map[string]any{"date": "2024-01-15"},
		},
		{
			Slug:   "post2",
			Config: map[string]any{"date": "2025-06-20T10:00:00Z"},
		},
		{
			Slug:   "post3",
			Config: map[string]any{"date": "2023-12-05"},
		},
	}

	sorted, err := SortBy("date", nodes)
	if err != nil {
		t.Fatalf("SortBy() error = %v", err)
	}

	expectedOrder := []string{"post2", "post1", "post3"}
	for i, expectedSlug := range expectedOrder {
		if sorted[i].Slug != expectedSlug {
			t.Errorf("SortBy()[%d].Slug = %s, want %s", i, sorted[i].Slug, expectedSlug)
		}
	}
}

func TestSortByWithMissingDates(t *testing.T) {
	nodes := []*content.Node{
		{
			Slug:   "post1",
			Config: map[string]any{"date": time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		},
		{
			Slug:   "no-date",
			Config: map[string]any{"title": "No Date Post"},
		},
		{
			Slug:   "post2",
			Config: map[string]any{"date": time.Date(2025, 6, 20, 0, 0, 0, 0, time.UTC)},
		},
	}

	sorted, err := SortBy("date", nodes)
	if err != nil {
		t.Fatalf("SortBy() error = %v", err)
	}

	// Items with valid dates should come first
	if sorted[0].Slug != "post2" {
		t.Errorf("First item = %s, want post2 (newest date)", sorted[0].Slug)
	}
	if sorted[1].Slug != "post1" {
		t.Errorf("Second item = %s, want post1", sorted[1].Slug)
	}
	// Item without date should be pushed to end
	if sorted[2].Slug != "no-date" {
		t.Errorf("Last item = %s, want no-date", sorted[2].Slug)
	}
}

func TestExtractTime(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  bool // true if valid time, false if zero time
	}{
		{"time.Time", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), true},
		{"RFC3339 string", "2024-01-15T10:00:00Z", true},
		{"date-only string", "2024-01-15", true},
		{"invalid string", "not a date", false},
		{"invalid type", 12345, false},
		{"nil", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTime(tt.input)
			isValid := !result.IsZero()
			if isValid != tt.want {
				t.Errorf("extractTime(%v) valid = %v, want %v", tt.input, isValid, tt.want)
			}
		})
	}
}

func TestRelatedPosts(t *testing.T) {
	now := time.Now()

	allPosts := []*content.Node{
		{
			Slug:      "post1",
			Permalink: "/post1/",
			Config: map[string]any{
				"tags": []string{"go", "testing"},
				"date": now,
			},
		},
		{
			Slug:      "post2",
			Permalink: "/post2/",
			Config: map[string]any{
				"tags": []string{"go", "performance"},
				"date": now.AddDate(0, 0, -10),
			},
		},
		{
			Slug:      "post3",
			Permalink: "/post3/",
			Config: map[string]any{
				"tags": []string{"go", "testing", "ci"},
				"date": now.AddDate(0, 0, -5),
			},
		},
		{
			Slug:      "post4",
			Permalink: "/post4/",
			Config: map[string]any{
				"tags": []string{"python"},
				"date": now.AddDate(0, 0, -1),
			},
		},
	}

	currentPage := allPosts[0] // post1 with tags: go, testing

	related := RelatedPosts(currentPage, allPosts)

	// Should return related posts (not empty if matches exist)
	if len(related) == 0 {
		t.Fatal("Expected related posts but got none")
	}

	// Should not include current post itself
	for _, r := range related {
		if r.Permalink == currentPage.Permalink {
			t.Error("RelatedPosts should not include current post itself")
		}
	}

	// Should not include posts with no tag overlap
	for _, r := range related {
		if r.Permalink == "/post4/" {
			t.Error("Should not include post4 (no tag overlap)")
		}
	}

	// Should return at most 6 posts (configured for client randomization)
	if len(related) > 6 {
		t.Errorf("Expected at most 6 related posts, got %d", len(related))
	}
}

func TestRelatedPosts_NoTags(t *testing.T) {
	now := time.Now()

	allPosts := []*content.Node{
		{
			Slug:      "post1",
			Permalink: "/post1/",
			Config: map[string]any{
				"date": now,
			},
		},
		{
			Slug:      "post2",
			Permalink: "/post2/",
			Config: map[string]any{
				"tags": []string{"go"},
				"date": now,
			},
		},
	}

	currentPage := allPosts[0] // post1 with no tags

	related := RelatedPosts(currentPage, allPosts)

	// Should return no related posts when current page has no tags
	if len(related) != 0 {
		t.Errorf("Expected no related posts for page with no tags, got %d", len(related))
	}
}

func TestRelatedPosts_EmptyInput(t *testing.T) {
	currentPage := &content.Node{
		Slug:      "post1",
		Permalink: "/post1/",
		Config: map[string]any{
			"tags": []string{"go"},
		},
	}

	// Empty posts list
	related := RelatedPosts(currentPage, []*content.Node{})

	if len(related) != 0 {
		t.Errorf("Expected no related posts with empty posts list, got %d", len(related))
	}
}

// TestRelatedPosts_WithTemplateContext tests the critical bug fix where
// template context (map[string]any with nested Config field) wasn't
// being parsed correctly, causing tag extraction to fail.
func TestRelatedPosts_WithTemplateContext(t *testing.T) {
	now := time.Now()

	// Simulate template context structure as passed from templates
	// Structure: { "Config": {...}, "Permalink": "...", "Series": ... }
	templateContext := map[string]any{
		"Config": map[string]any{
			"title": "Test Post",
			"tags":  []any{"blog", "go"}, // tags as []any (TOML parsing result)
			"date":  now,
		},
		"Permalink": "/blog/test-post/",
	}

	allPosts := []*content.Node{
		{
			Permalink: "/blog/post-1/",
			Config: map[string]any{
				"title": "Post 1",
				"tags":  []string{"blog", "tutorial"},
				"date":  now.AddDate(0, 0, -5),
			},
		},
		{
			Permalink: "/blog/post-2/",
			Config: map[string]any{
				"title": "Post 2",
				"tags":  []string{"go", "programming"},
				"date":  now.AddDate(0, 0, -3),
			},
		},
		{
			Permalink: "/blog/post-3/",
			Config: map[string]any{
				"title": "Post 3",
				"tags":  []string{"python"},
				"date":  now.AddDate(0, 0, -1),
			},
		},
	}

	// This should extract tags correctly from nested Config
	related := RelatedPosts(templateContext, allPosts)

	// Should find related posts (post-1 matches "blog", post-2 matches "go")
	if len(related) == 0 {
		t.Fatal("Expected to find related posts from template context, got none")
	}

	// Verify we got posts with matching tags
	foundBlogMatch := false
	foundGoMatch := false
	for _, r := range related {
		if r.Permalink == "/blog/post-1/" {
			foundBlogMatch = true
		}
		if r.Permalink == "/blog/post-2/" {
			foundGoMatch = true
		}
		// Should not include post-3 (no tag overlap)
		if r.Permalink == "/blog/post-3/" {
			t.Error("Should not include post-3 (no tag overlap with 'blog' or 'go')")
		}
	}

	if !foundBlogMatch {
		t.Error("Expected to find post-1 (matches 'blog' tag)")
	}
	if !foundGoMatch {
		t.Error("Expected to find post-2 (matches 'go' tag)")
	}
}

// TestRelatedPosts_WithLegacyDirectConfig tests backwards compatibility
// when Config is passed directly (not nested in template context)
func TestRelatedPosts_WithLegacyDirectConfig(t *testing.T) {
	now := time.Now()

	// Legacy usage: passing config directly (no "Config" wrapper)
	directConfig := map[string]any{
		"title":     "Test Post",
		"tags":      []string{"go", "testing"},
		"date":      now,
		"Permalink": "/blog/legacy-post/",
	}

	allPosts := []*content.Node{
		{
			Permalink: "/blog/related/",
			Config: map[string]any{
				"tags": []string{"go"},
				"date": now.AddDate(0, 0, -1),
			},
		},
	}

	// Should still work with legacy direct config
	related := RelatedPosts(directConfig, allPosts)

	if len(related) == 0 {
		t.Error("Expected to find related posts with legacy direct config")
	}
}

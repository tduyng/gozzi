package content

import (
	"testing"
	"time"
)

func TestRelatedPostsFinder_Basic(t *testing.T) {
	now := time.Now()

	posts := []*Node{
		createTestPost("post1", []string{"go", "testing"}, now),
		createTestPost("post2", []string{"go", "performance"}, now.AddDate(0, 0, -10)),
		createTestPost("post3", []string{"go", "testing", "ci"}, now.AddDate(0, 0, -20)),
		createTestPost("post4", []string{"python", "testing"}, now.AddDate(0, 0, -5)),
		createTestPost("post5", []string{"javascript"}, now.AddDate(0, 0, -1)),
	}

	config := DefaultRelatedConfig()
	config.ResultLimit = 3
	config.RandomBonusMax = 0 // Disable randomness for predictable tests

	finder := NewRelatedPostsFinder(posts, config)

	// Test finding related posts for post1 (tags: go, testing)
	related := finder.FindRelated(posts[0])

	if len(related) == 0 {
		t.Fatal("Expected related posts but got none")
	}

	// Should prioritize posts with more matching tags
	// post3 has 2 matching tags (go, testing)
	// post2 and post4 each have 1 matching tag
	foundPost3 := false
	for _, r := range related {
		if r.Permalink == "/post3/" {
			foundPost3 = true
		}
		// Should NOT include post5 (no tag overlap)
		if r.Permalink == "/post5/" {
			t.Error("post5 should not be related (no tag overlap)")
		}
		// Should NOT include self
		if r.Permalink == "/post1/" {
			t.Error("Should not include current post itself")
		}
	}

	if !foundPost3 {
		t.Error("Expected post3 (2 tag matches) to be in related posts")
	}
}

func TestRelatedPostsFinder_NoTags(t *testing.T) {
	now := time.Now()

	posts := []*Node{
		createTestPost("post1", []string{}, now),
		createTestPost("post2", []string{"go"}, now),
	}

	config := DefaultRelatedConfig()
	finder := NewRelatedPostsFinder(posts, config)

	// Post with no tags should have no related posts
	related := finder.FindRelated(posts[0])
	if len(related) != 0 {
		t.Errorf("Expected no related posts for post with no tags, got %d", len(related))
	}
}

func TestRelatedPostsFinder_NoMatches(t *testing.T) {
	now := time.Now()

	posts := []*Node{
		createTestPost("post1", []string{"go", "testing"}, now),
		createTestPost("post2", []string{"python", "django"}, now),
		createTestPost("post3", []string{"javascript", "react"}, now),
	}

	config := DefaultRelatedConfig()
	finder := NewRelatedPostsFinder(posts, config)

	// Post with completely different tags should have no related posts
	related := finder.FindRelated(posts[0])
	if len(related) != 0 {
		t.Errorf("Expected no related posts when no tag overlap, got %d", len(related))
	}
}

func TestRelatedPostsFinder_ResultLimit(t *testing.T) {
	now := time.Now()

	// Create 10 posts all with "go" tag
	posts := make([]*Node, 11)
	posts[0] = createTestPost("current", []string{"go"}, now)
	for i := 1; i <= 10; i++ {
		posts[i] = createTestPost("post"+string(rune('0'+i)), []string{"go"}, now.AddDate(0, 0, -i))
	}

	config := DefaultRelatedConfig()
	config.ResultLimit = 3

	finder := NewRelatedPostsFinder(posts, config)
	related := finder.FindRelated(posts[0])

	if len(related) != 3 {
		t.Errorf("Expected exactly 3 related posts, got %d", len(related))
	}
}

func TestRelatedPostsFinder_MinTagMatches(t *testing.T) {
	now := time.Now()

	posts := []*Node{
		createTestPost("post1", []string{"go", "testing", "ci"}, now),
		createTestPost("post2", []string{"go"}, now),                  // 1 match
		createTestPost("post3", []string{"go", "testing"}, now),       // 2 matches
		createTestPost("post4", []string{"go", "testing", "ci"}, now), // 3 matches
		createTestPost("post5", []string{"python"}, now),              // 0 matches
	}

	config := DefaultRelatedConfig()
	config.MinTagMatches = 2 // Require at least 2 matching tags
	config.RandomBonusMax = 0

	finder := NewRelatedPostsFinder(posts, config)
	related := finder.FindRelated(posts[0])

	// Should only include post3 and post4 (2+ tag matches)
	for _, r := range related {
		permalink := r.Permalink
		if permalink == "/post2/" || permalink == "/post5/" {
			t.Errorf("Post %s should not be included (less than 2 tag matches)", permalink)
		}
	}

	if len(related) != 2 {
		t.Errorf("Expected 2 related posts with MinTagMatches=2, got %d", len(related))
	}
}

func TestRelatedPostsFinder_RecencyBonus(t *testing.T) {
	now := time.Now()

	posts := []*Node{
		createTestPost("current", []string{"go"}, now),
		createTestPost("recent", []string{"go"}, now.AddDate(0, 0, -5)), // 5 days old
		createTestPost("old", []string{"go"}, now.AddDate(0, 0, -365)),  // 1 year old
	}

	config := DefaultRelatedConfig()
	config.ResultLimit = 2
	config.RecencyDecayDays = 30
	config.RandomBonusMax = 0 // Disable randomness

	finder := NewRelatedPostsFinder(posts, config)
	related := finder.FindRelated(posts[0])

	if len(related) == 0 {
		t.Fatal("Expected related posts")
	}

	// Recent post should rank higher than old post
	if related[0].Permalink != "/recent/" {
		t.Errorf("Expected recent post first, got %s", related[0].Permalink)
	}
}

func TestRelatedPostsFinder_TagIndexEfficiency(t *testing.T) {
	now := time.Now()

	// Create 100 posts with various tags
	posts := make([]*Node, 100)
	for i := range 100 {
		var tags []string
		if i%10 == 0 {
			tags = []string{"go", "testing"}
		} else if i%5 == 0 {
			tags = []string{"go"}
		} else {
			tags = []string{"python", "javascript", "rust"}
		}
		posts[i] = createTestPost("post"+string(rune(i)), tags, now.AddDate(0, 0, -i))
	}

	config := DefaultRelatedConfig()
	finder := NewRelatedPostsFinder(posts, config)

	// Tag index should have been built
	if len(finder.tagIndex) == 0 {
		t.Error("Tag index should not be empty")
	}

	// Finding related should be fast (O(k), not O(n²))
	startTime := time.Now()
	related := finder.FindRelated(posts[0])
	duration := time.Since(startTime)

	if duration > 10*time.Millisecond {
		t.Errorf("Finding related posts took too long: %v (expected <10ms)", duration)
	}

	if len(related) == 0 {
		t.Error("Expected some related posts")
	}
}

func TestRelatedPostsFinder_Randomization(t *testing.T) {
	now := time.Now()

	// Create posts with identical tag matches
	posts := []*Node{
		createTestPost("current", []string{"go"}, now),
		createTestPost("match1", []string{"go"}, now.AddDate(0, 0, -1)),
		createTestPost("match2", []string{"go"}, now.AddDate(0, 0, -2)),
		createTestPost("match3", []string{"go"}, now.AddDate(0, 0, -3)),
		createTestPost("match4", []string{"go"}, now.AddDate(0, 0, -4)),
	}

	config := DefaultRelatedConfig()
	config.ResultLimit = 2
	config.RandomBonusMax = 5 // High randomness

	finder := NewRelatedPostsFinder(posts, config)

	// Run multiple times and check we get different results
	results := make(map[string]int)
	runs := 50

	for range runs {
		related := finder.FindRelated(posts[0])
		for _, r := range related {
			results[r.Permalink]++
		}
	}

	// With randomization, we should see multiple different posts selected
	if len(results) < 3 {
		t.Errorf("Expected randomization to select at least 3 different posts across %d runs, got %d", runs, len(results))
	}
}

func TestExtractTags_VariousFormats(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]any
		expected []string
	}{
		{
			name:     "string slice",
			config:   map[string]any{"tags": []string{"go", "testing"}},
			expected: []string{"go", "testing"},
		},
		{
			name:     "any slice",
			config:   map[string]any{"tags": []any{"go", "testing"}},
			expected: []string{"go", "testing"},
		},
		{
			name:     "no tags",
			config:   map[string]any{},
			expected: []string{},
		},
		{
			name:     "nil config",
			config:   nil,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{Config: tt.config}
			finder := &RelatedPostsFinder{}
			tags := finder.extractTags(node)

			if len(tags) != len(tt.expected) {
				t.Errorf("Expected %d tags, got %d", len(tt.expected), len(tags))
			}
		})
	}
}

// Helper to create test posts.
func createTestPost(slug string, tags []string, date time.Time) *Node {
	return &Node{
		Slug:      slug,
		Permalink: "/" + slug + "/",
		Config: map[string]any{
			"tags": tags,
			"date": date,
		},
	}
}

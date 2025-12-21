package builder

import (
	"testing"
	"time"

	"github.com/tduyng/gozzi/app/cache"
	"github.com/tduyng/gozzi/app/config"
)

func TestCreateStableCacheKey_TagPage(t *testing.T) {
	builder := &Builder{
		site: &config.Site{
			BaseURL: "https://example.com",
			Title:   "Test Site",
		},
	}

	// Create test data for a tag page
	data1 := map[string]any{
		"Site": map[string]any{
			"Config": builder.site.ToConfig(),
		},
		"Page": map[string]any{
			"Tag": "golang",
			"Pages": []map[string]any{
				{
					"Permalink": "/blog/post1/",
					"Config": map[string]any{
						"title": "Post 1",
						"date":  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						"extra": map[string]any{
							"featured": true,
						},
					},
				},
				{
					"Permalink": "/blog/post2/",
					"Config": map[string]any{
						"title": "Post 2",
						"date":  time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						"extra": map[string]any{
							"featured": false,
						},
					},
				},
			},
		},
	}

	// Create identical data with different BuildTime
	builder.site.BuildTime = time.Now()
	key1 := builder.createStableCacheKey("tag.html", data1)
	hash1, err := cache.ComputeDataHash(key1)
	if err != nil {
		t.Fatalf("Failed to compute hash1: %v", err)
	}

	// Change BuildTime (simulates second build)
	builder.site.BuildTime = time.Now().Add(5 * time.Second)

	// Create data2 with updated Site config but same tag/pages
	data2 := map[string]any{
		"Site": map[string]any{
			"Config": builder.site.ToConfig(), // Different BuildTime!
		},
		"Page": map[string]any{
			"Tag": "golang",
			"Pages": []map[string]any{
				{
					"Permalink": "/blog/post1/",
					"Config": map[string]any{
						"title": "Post 1",
						"date":  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						"extra": map[string]any{
							"featured": true,
						},
					},
				},
				{
					"Permalink": "/blog/post2/",
					"Config": map[string]any{
						"title": "Post 2",
						"date":  time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
						"extra": map[string]any{
							"featured": false,
						},
					},
				},
			},
		},
	}

	key2 := builder.createStableCacheKey("tag.html", data2)
	hash2, err := cache.ComputeDataHash(key2)
	if err != nil {
		t.Fatalf("Failed to compute hash2: %v", err)
	}

	// Hashes should be IDENTICAL even though BuildTime changed
	if hash1 != hash2 {
		t.Errorf("Expected identical hashes for same tag data despite different BuildTime.\nHash1: %s\nHash2: %s\nKey1: %+v\nKey2: %+v",
			hash1, hash2, key1, key2)
	}
}

func TestCreateStableCacheKey_DifferentTagContent(t *testing.T) {
	builder := &Builder{
		site: &config.Site{
			BaseURL: "https://example.com",
			Title:   "Test Site",
		},
	}

	builder.site.BuildTime = time.Now()

	// Tag page with post1
	data1 := map[string]any{
		"Site": map[string]any{
			"Config": builder.site.ToConfig(),
		},
		"Page": map[string]any{
			"Tag": "golang",
			"Pages": []map[string]any{
				{
					"Permalink": "/blog/post1/",
					"Config": map[string]any{
						"title": "Post 1",
						"date":  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	// Tag page with post2 (different content)
	data2 := map[string]any{
		"Site": map[string]any{
			"Config": builder.site.ToConfig(),
		},
		"Page": map[string]any{
			"Tag": "golang",
			"Pages": []map[string]any{
				{
					"Permalink": "/blog/post2/",
					"Config": map[string]any{
						"title": "Post 2",                                    // Different title
						"date":  time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), // Different date
					},
				},
			},
		},
	}

	key1 := builder.createStableCacheKey("tag.html", data1)
	hash1, err := cache.ComputeDataHash(key1)
	if err != nil {
		t.Fatalf("Failed to compute hash1: %v", err)
	}

	key2 := builder.createStableCacheKey("tag.html", data2)
	hash2, err := cache.ComputeDataHash(key2)
	if err != nil {
		t.Fatalf("Failed to compute hash2: %v", err)
	}

	// Hashes should be DIFFERENT (different posts)
	if hash1 == hash2 {
		t.Errorf("Expected different hashes for different tag content.\nHash: %s\nKey1: %+v\nKey2: %+v",
			hash1, key1, key2)
	}
}

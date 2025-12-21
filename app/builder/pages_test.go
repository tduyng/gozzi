package builder

import (
	"fmt"
	"testing"
	"time"

	"github.com/tduyng/gozzi/app/cache"
	"github.com/tduyng/gozzi/app/config"
	"github.com/tduyng/gozzi/app/content"
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

func TestPageCacheKey_ExtraConfigChanges(t *testing.T) {
	// Test that changes to extra config (comment, toc, reaction, etc.) invalidate cache
	node1 := &content.Node{
		Path:      "blog/post1",
		Content:   "Same content",
		WordCount: 100,
		ReadTime:  1,
		Config: map[string]any{
			"title":    "Test Post",
			"date":     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			"template": "post.html",
			"extra": map[string]any{
				"comment":  true,
				"reaction": true,
				"toc":      true,
			},
		},
	}

	node2 := &content.Node{
		Path:      "blog/post1",
		Content:   "Same content",
		WordCount: 100,
		ReadTime:  1,
		Config: map[string]any{
			"title":    "Test Post",
			"date":     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			"template": "post.html",
			"extra": map[string]any{
				"comment":  false, // Changed!
				"reaction": true,
				"toc":      false, // Changed!
			},
		},
	}

	// Create cache keys
	key1 := map[string]any{
		"Path":      node1.Path,
		"Content":   string(node1.Content),
		"WordCount": node1.WordCount,
		"ReadTime":  node1.ReadTime,
		"Title":     node1.Config["title"],
		"Date":      node1.Config["date"].(time.Time).Format("2006-01-02"),
		"Template":  node1.Config["template"],
		"Extra":     fmt.Sprintf("%+v", node1.Config["extra"]),
	}

	key2 := map[string]any{
		"Path":      node2.Path,
		"Content":   string(node2.Content),
		"WordCount": node2.WordCount,
		"ReadTime":  node2.ReadTime,
		"Title":     node2.Config["title"],
		"Date":      node2.Config["date"].(time.Time).Format("2006-01-02"),
		"Template":  node2.Config["template"],
		"Extra":     fmt.Sprintf("%+v", node2.Config["extra"]),
	}

	hash1, err := cache.ComputeDataHash(key1)
	if err != nil {
		t.Fatalf("Failed to compute hash1: %v", err)
	}

	hash2, err := cache.ComputeDataHash(key2)
	if err != nil {
		t.Fatalf("Failed to compute hash2: %v", err)
	}

	// Hashes should be DIFFERENT because extra config changed
	if hash1 == hash2 {
		t.Errorf("Expected different hashes when extra config changes.\nHash: %s\nNode1 extra: %+v\nNode2 extra: %+v",
			hash1, node1.Config["extra"], node2.Config["extra"])
	}
}

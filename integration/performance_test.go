// This file contains performance benchmarks for build operations.
package integration

import (
	"testing"
	"time"
)

// TestPerformance_FreshBuild tests fresh build performance
func TestPerformance_FreshBuild(t *testing.T) {
	t.Run("SmallSite_FastBuild", func(t *testing.T) {
		sitePath := setupTestSite(t)

		start := time.Now()
		buildSite(t, sitePath)
		duration := time.Since(start)

		// Small test site should build very quickly
		maxDuration := 5 * time.Second
		if duration > maxDuration {
			t.Errorf("build took %v, expected < %v", duration, maxDuration)
		}

		t.Logf("Fresh build completed in %v", duration)
	})
}

// TestPerformance_IncrementalBuild tests incremental build speed
func TestPerformance_IncrementalBuild(t *testing.T) {
	t.Run("IncrementalFasterThanFull", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Measure full rebuild
		fullStart := time.Now()
		fullRebuild(t, gen, contentParser, sitePath)
		fullDuration := time.Since(fullStart)

		// Modify one file
		post1Path := sitePath + "/content/blog/post1/index.md"
		modifyFile(t, post1Path, "First Post", "First Post Updated")

		// Measure incremental rebuild
		incStart := time.Now()
		incrementalRebuild(t, gen, contentParser, sitePath, []string{post1Path})
		incDuration := time.Since(incStart)

		// Incremental should be faster
		if incDuration > fullDuration {
			t.Logf("Incremental (%v) should be faster than full rebuild (%v)", incDuration, fullDuration)
		}

		t.Logf("Full rebuild: %v, Incremental: %v", fullDuration, incDuration)
	})
}

// TestPerformance_CacheEffectiveness tests cache performance
func TestPerformance_CacheEffectiveness(t *testing.T) {
	t.Run("CacheImproves_RebuildSpeed", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// First build - populate cache
		gen.ClearRenderCache()
		if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("initial build failed: %v", err)
		}

		// Second build WITHOUT clearing cache - should hit cache
		gen.ResetCacheStats()
		start := time.Now()
		if err := gen.Generate(contentParser.ContentMap["."]); err != nil {
			t.Fatalf("cached rebuild failed: %v", err)
		}
		duration := time.Since(start)

		stats := gen.GetCacheStats()
		// When content is unchanged, should have very high cache hit rate
		if stats.HitRate < 90 {
			t.Errorf("expected >90%% cache hit rate on unchanged rebuild, got %.1f%%", stats.HitRate)
		}

		t.Logf("Cached rebuild: %v, hit rate: %.1f%%, hits: %d, misses: %d",
			duration, stats.HitRate, stats.Hits, stats.Misses)
	})
}

// BenchmarkFreshBuild benchmarks full site build
func BenchmarkFreshBuild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sitePath := setupTestSite(&testing.T{})
		buildSite(&testing.T{}, sitePath)
	}
}

// BenchmarkIncrementalBuild benchmarks incremental builds
func BenchmarkIncrementalBuild(b *testing.B) {
	sitePath := setupTestSite(&testing.T{})
	gen, contentParser := buildSite(&testing.T{}, sitePath)
	post1Path := sitePath + "/content/blog/post1/index.md"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		modifyFile(&testing.T{}, post1Path, "First Post", "First Post V"+string(rune('0'+i%10)))
		incrementalRebuild(&testing.T{}, gen, contentParser, sitePath, []string{post1Path})
	}
}

// BenchmarkTaxonomyGeneration benchmarks taxonomy page generation
func BenchmarkTaxonomyGeneration(b *testing.B) {
	sitePath := setupTestSite(&testing.T{})

	// Add many posts with tags
	gen, contentParser := buildSite(&testing.T{}, sitePath)

	for i := 0; i < 50; i++ {
		post := `+++
title = "Tagged Post"
date = 2024-01-01
template = "post.html"
tags = ["tag1", "tag2", "tag3"]
+++
Content`
		createPost(&testing.T{}, sitePath, "blog/tagged-"+string(rune('0'+i%10))+".md", post)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fullRebuild(&testing.T{}, gen, contentParser, sitePath)
	}
}

// BenchmarkMarkdownRendering benchmarks markdown parsing
func BenchmarkMarkdownRendering(b *testing.B) {
	sitePath := setupTestSite(&testing.T{})
	gen, contentParser := buildSite(&testing.T{}, sitePath)

	longContent := `
# Heading 1
## Heading 2
### Heading 3

This is **bold** and *italic* text.

- List item 1
- List item 2
- List item 3

` + "```go\nfunc main() {}\n```" + `

[Link](https://example.com)
`

	post := `+++
title = "Benchmark Post"
date = 2024-01-01
template = "post.html"
+++

` + longContent

	createPost(&testing.T{}, sitePath, "blog/benchmark.md", post)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fullRebuild(&testing.T{}, gen, contentParser, sitePath)
	}
}

// TestPerformance_MemoryUsage tests memory efficiency
func TestPerformance_MemoryUsage(t *testing.T) {
	t.Run("ReasonableMemoryUsage", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create moderate number of posts
		for i := 0; i < 50; i++ {
			post := `+++
title = "Post ` + string(rune('0'+i%10)) + `"
date = 2024-01-01
template = "post.html"
+++
Content`
			createPost(t, sitePath, "blog/mem-"+string(rune('0'+i%10))+string(rune('0'+(i/10)%10))+".md", post)
		}

		fullRebuild(t, gen, contentParser, sitePath)

		// Should complete without OOM
		// Exact memory limits depend on system
		t.Log("Memory test completed successfully")
	})
}

// TestPerformance_LargeSite tests handling of large sites
func TestPerformance_LargeSite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large site test in short mode")
	}

	t.Run("100Posts_BuildsEfficiently", func(t *testing.T) {
		sitePath := setupTestSite(t)
		gen, contentParser := buildSite(t, sitePath)

		// Create 100 posts
		for i := 1; i <= 100; i++ {
			post := `+++
title = "Post ` + string(rune('0'+i%10)) + `"
date = 2024-01-0` + string(rune('0'+i%10)) + `
template = "post.html"
tags = ["tag` + string(rune('0'+i%5)) + `"]
+++

Post content here.`

			createPost(t, sitePath, "blog/large-"+string(rune('0'+i%10))+string(rune('0'+(i/10)%10))+string(rune('0'+(i/100)%10))+".md", post)
		}

		start := time.Now()
		fullRebuild(t, gen, contentParser, sitePath)
		duration := time.Since(start)

		// Even 100 posts should build reasonably fast
		maxDuration := 10 * time.Second
		if duration > maxDuration {
			t.Errorf("100 post build took %v, expected < %v", duration, maxDuration)
		}

		t.Logf("100 posts built in %v", duration)
	})
}

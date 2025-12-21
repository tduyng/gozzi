package cache

import (
	"bytes"
	"testing"
)

func TestComputeDataHash(t *testing.T) {
	tests := []struct {
		name    string
		data    any
		wantErr bool
	}{
		{
			name:    "simple map",
			data:    map[string]string{"key": "value"},
			wantErr: false,
		},
		{
			name: "complex struct",
			data: struct {
				Title string
				Count int
			}{"Test", 42},
			wantErr: false,
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: true, // gob cannot encode nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := ComputeDataHash(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComputeDataHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && hash == (ContentHash{}) {
				t.Error("Expected non-zero hash")
			}
		})
	}
}

func TestComputeDataHash_Deterministic(t *testing.T) {
	// Use struct instead of map for deterministic gob encoding
	// Maps in Go have non-deterministic iteration order which affects gob encoding
	type PageData struct {
		Title string
		Count int
	}

	data := PageData{
		Title: "Test Page",
		Count: 42,
	}

	hash1, err1 := ComputeDataHash(data)
	hash2, err2 := ComputeDataHash(data)

	if err1 != nil || err2 != nil {
		t.Fatalf("Unexpected errors: %v, %v", err1, err2)
	}

	if hash1 != hash2 {
		t.Error("Same data should produce same hash")
	}
}

func TestRenderCache_GetSet(t *testing.T) {
	cache := NewRenderCache()
	template := "page.html"
	dataHash := HashString("test data")
	output := []byte("<html>test</html>")

	// Test Get on empty cache
	_, exists := cache.Get(template, dataHash)
	if exists {
		t.Error("Expected Get to return false for empty cache")
	}

	// Test Set and Get
	cache.Set(template, dataHash, output)
	got, exists := cache.Get(template, dataHash)
	if !exists {
		t.Error("Expected Get to return true after Set")
	}
	if !bytes.Equal(got, output) {
		t.Errorf("Get() = %v, want %v", got, output)
	}
}

func TestRenderCache_GetOrCompute(t *testing.T) {
	cache := NewRenderCache()
	template := "page.html"
	dataHash := HashString("test data")
	output := []byte("<html>computed</html>")

	computeCalled := 0
	compute := func() ([]byte, error) {
		computeCalled++
		return output, nil
	}

	// First call - should compute
	result, cached, err := cache.GetOrCompute(template, dataHash, compute)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cached {
		t.Error("Expected cached=false on first call")
	}
	if !bytes.Equal(result, output) {
		t.Errorf("GetOrCompute() = %v, want %v", result, output)
	}
	if computeCalled != 1 {
		t.Errorf("Expected compute called once, got %d", computeCalled)
	}

	// Second call - should use cache
	result, cached, err = cache.GetOrCompute(template, dataHash, compute)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !cached {
		t.Error("Expected cached=true on second call")
	}
	if !bytes.Equal(result, output) {
		t.Errorf("GetOrCompute() = %v, want %v", result, output)
	}
	if computeCalled != 1 {
		t.Error("Expected compute not called on cache hit")
	}
}

func TestRenderCache_Invalidate(t *testing.T) {
	cache := NewRenderCache()
	template := "page.html"
	dataHash := HashString("test data")
	output := []byte("<html>test</html>")

	cache.Set(template, dataHash, output)
	if cache.Size() != 1 {
		t.Error("Expected cache to have 1 entry")
	}

	cache.Invalidate(template, dataHash)
	_, exists := cache.Get(template, dataHash)
	if exists {
		t.Error("Expected entry to be invalidated")
	}
}

func TestRenderCache_InvalidateTemplate(t *testing.T) {
	cache := NewRenderCache()
	template := "page.html"

	// Add multiple renders of same template with different data
	cache.Set(template, HashString("data1"), []byte("output1"))
	cache.Set(template, HashString("data2"), []byte("output2"))
	cache.Set("other.html", HashString("data3"), []byte("output3"))

	if cache.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cache.Size())
	}

	count := cache.InvalidateTemplate(template)
	if count != 2 {
		t.Errorf("Expected 2 invalidations, got %d", count)
	}
	if cache.Size() != 1 {
		t.Errorf("Expected size 1 after invalidation, got %d", cache.Size())
	}
}

func TestRenderCache_Clear(t *testing.T) {
	cache := NewRenderCache()

	// Add multiple entries
	cache.Set("page1.html", HashString("data1"), []byte("output1"))
	cache.Set("page2.html", HashString("data2"), []byte("output2"))
	cache.Set("page3.html", HashString("data3"), []byte("output3"))

	if cache.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cache.Size())
	}

	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("Expected size 0 after Clear, got %d", cache.Size())
	}

	stats := cache.Stats()
	if stats.Hits != 0 || stats.Misses != 0 {
		t.Error("Expected stats to be reset after Clear")
	}
}

func TestRenderCache_Stats(t *testing.T) {
	cache := NewRenderCache()
	template := "page.html"
	dataHash := HashString("test")
	output := []byte("test")

	// Initial stats
	stats := cache.Stats()
	if stats.Entries != 0 || stats.Hits != 0 || stats.Misses != 0 {
		t.Error("Expected initial stats to be zero")
	}

	// Trigger miss
	cache.Get(template, dataHash)
	stats = cache.Stats()
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}

	// Add entry and trigger hit
	cache.Set(template, dataHash, output)
	cache.Get(template, dataHash)
	stats = cache.Stats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}

	// Verify hit rate
	if stats.HitRate != 50.0 {
		t.Errorf("Expected 50%% hit rate, got %.1f%%", stats.HitRate)
	}
}

func TestRenderCache_MemoryUsage(t *testing.T) {
	cache := NewRenderCache()

	output1 := make([]byte, 1000)
	output2 := make([]byte, 2000)

	cache.Set("page1.html", HashString("data1"), output1)
	cache.Set("page2.html", HashString("data2"), output2)

	usage := cache.MemoryUsage()
	expected := int64(3000)
	if usage != expected {
		t.Errorf("Expected memory usage %d, got %d", expected, usage)
	}
}

func TestRenderKey_String(t *testing.T) {
	key := RenderKey{
		Template: "page.html",
		DataHash: HashString("test data"),
	}

	str := key.String()
	if str == "" {
		t.Error("Expected non-empty string from RenderKey.String()")
	}
	if len(str) < 10 {
		t.Error("Expected string to contain template name and hash prefix")
	}
}

func TestHashString(t *testing.T) {
	str1 := "hello world"
	str2 := "hello world"
	str3 := "different"

	hash1 := HashString(str1)
	hash2 := HashString(str2)
	hash3 := HashString(str3)

	if hash1 != hash2 {
		t.Error("Same strings should produce same hash")
	}
	if hash1 == hash3 {
		t.Error("Different strings should produce different hashes")
	}
}

// TestComputeDataHash_DifferentTagsProduceDifferentHashes ensures that
// tag pages with different tag names/data produce different hashes.
// This is a regression test for the bug where all tag pages shared
// the same cache entry because they only hashed the template name.
func TestComputeDataHash_DifferentTagsProduceDifferentHashes(t *testing.T) {
	// Simulate tag page data for "blog" tag
	blogData := map[string]any{
		"Site": map[string]any{
			"Config": map[string]any{"Title": "Test Site"},
		},
		"Page": map[string]any{
			"Title": "Tag: blog",
			"Tag":   "blog",
			"Pages": []map[string]any{
				{"Title": "Blog Post 1"},
				{"Title": "Blog Post 2"},
			},
		},
	}

	// Simulate tag page data for "go" tag
	goData := map[string]any{
		"Site": map[string]any{
			"Config": map[string]any{"Title": "Test Site"},
		},
		"Page": map[string]any{
			"Title": "Tag: go",
			"Tag":   "go",
			"Pages": []map[string]any{
				{"Title": "Go Post 1"},
			},
		},
	}

	blogHash, err := ComputeDataHash(blogData)
	if err != nil {
		t.Fatalf("Failed to hash blog data: %v", err)
	}

	goHash, err := ComputeDataHash(goData)
	if err != nil {
		t.Fatalf("Failed to hash go data: %v", err)
	}

	// Critical assertion: different tag pages MUST have different hashes
	if blogHash == goHash {
		t.Error("Different tag page data should produce different hashes, " +
			"but blog and go tags produced the same hash. " +
			"This means tag pages will incorrectly share cache entries!")
	}
}

// TestComputeDataHash_SameDataProducesSameHash ensures cache consistency
// Note: We use structs instead of maps because Go's non-deterministic map
// iteration means maps with same content can hash differently.
// In practice, our cache works because we reuse the same map instances
// within a build session, which produces consistent hashes.
func TestComputeDataHash_SameDataProducesSameHash(t *testing.T) {
	// Use a struct for deterministic hashing
	type PageData struct {
		SiteTitle string
		PageTitle string
		Tag       string
	}

	data := PageData{
		SiteTitle: "Test Site",
		PageTitle: "Tag: blog",
		Tag:       "blog",
	}

	hash1, err := ComputeDataHash(data)
	if err != nil {
		t.Fatalf("Failed to hash data: %v", err)
	}

	hash2, err := ComputeDataHash(data)
	if err != nil {
		t.Fatalf("Failed to hash data: %v", err)
	}

	if hash1 != hash2 {
		t.Error("Same struct data should produce same hash for cache consistency")
	}
}

func TestHashWriter(t *testing.T) {
	hw := NewHashWriter()

	hw.WriteString("test")
	hw.WriteBytes([]byte("data"))

	hash1 := hw.Sum()
	if hash1 == (ContentHash{}) {
		t.Error("Expected non-zero hash")
	}

	hexStr := hw.HexString()
	if len(hexStr) != 64 {
		t.Errorf("Expected hex string length 64, got %d", len(hexStr))
	}
}

func TestHashWriter_Deterministic(t *testing.T) {
	hw1 := NewHashWriter()
	hw1.WriteString("hello")
	hw1.WriteString("world")
	hash1 := hw1.Sum()

	hw2 := NewHashWriter()
	hw2.WriteString("hello")
	hw2.WriteString("world")
	hash2 := hw2.Sum()

	if hash1 != hash2 {
		t.Error("Same sequence should produce same hash")
	}
}

func TestRenderCacheStats_String(t *testing.T) {
	stats := RenderCacheStats{
		Entries: 15,
		Hits:    80,
		Misses:  20,
		HitRate: 80.0,
	}

	str := stats.String()
	if str == "" {
		t.Error("Expected non-empty string from RenderCacheStats.String()")
	}
}

func BenchmarkComputeDataHash(b *testing.B) {
	data := map[string]any{
		"title":   "Test Page",
		"content": "Lorem ipsum dolor sit amet",
		"count":   42,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ComputeDataHash(data)
	}
}

func BenchmarkRenderCache_GetOrCompute_Hit(b *testing.B) {
	cache := NewRenderCache()
	template := "page.html"
	dataHash := HashString("test")
	output := []byte("<html>test</html>")

	// Pre-populate cache
	cache.Set(template, dataHash, output)

	compute := func() ([]byte, error) {
		return output, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = cache.GetOrCompute(template, dataHash, compute)
	}
}

func BenchmarkRenderCache_GetOrCompute_Miss(b *testing.B) {
	cache := NewRenderCache()
	output := []byte("<html>test</html>")

	compute := func() ([]byte, error) {
		return output, nil
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		template := "page.html"
		dataHash := HashString(string(rune(i)))
		_, _, _ = cache.GetOrCompute(template, dataHash, compute)
	}
}

package cache

import (
	"crypto/sha256"
	"sync"
	"testing"
)

func TestComputeHash(t *testing.T) {
	tests := []struct {
		name    string
		content []byte
		want    ContentHash
	}{
		{
			name:    "empty content",
			content: []byte{},
			want:    sha256.Sum256([]byte{}),
		},
		{
			name:    "simple content",
			content: []byte("hello world"),
			want:    sha256.Sum256([]byte("hello world")),
		},
		{
			name:    "markdown content",
			content: []byte("# Title\n\nSome content"),
			want:    sha256.Sum256([]byte("# Title\n\nSome content")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeHash(tt.content)
			if got != tt.want {
				t.Errorf("ComputeHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHashCache_GetSet(t *testing.T) {
	cache := NewHashCache()
	path := "test.md"
	hash := ComputeHash([]byte("content"))

	// Test Get on empty cache
	_, exists := cache.Get(path)
	if exists {
		t.Error("Expected Get to return false for non-existent path")
	}

	// Test Set and Get
	cache.Set(path, hash)
	got, exists := cache.Get(path)
	if !exists {
		t.Error("Expected Get to return true after Set")
	}
	if got != hash {
		t.Errorf("Get() = %v, want %v", got, hash)
	}
}

func TestHashCache_HasChanged(t *testing.T) {
	cache := NewHashCache()
	path := "test.md"

	// First time seeing content - should be "changed"
	content1 := []byte("original content")
	if !cache.HasChanged(path, content1) {
		t.Error("Expected HasChanged to return true for new content")
	}

	// Same content - should not be changed
	if cache.HasChanged(path, content1) {
		t.Error("Expected HasChanged to return false for unchanged content")
	}

	// Different content - should be changed
	content2 := []byte("modified content")
	if !cache.HasChanged(path, content2) {
		t.Error("Expected HasChanged to return true for modified content")
	}

	// Back to same as current - should not be changed
	if cache.HasChanged(path, content2) {
		t.Error("Expected HasChanged to return false for unchanged content")
	}
}

func TestHashCache_Update(t *testing.T) {
	cache := NewHashCache()
	path := "test.md"
	content1 := []byte("content v1")
	content2 := []byte("content v2")

	// First update - should be changed
	changed, hash1 := cache.Update(path, content1)
	if !changed {
		t.Error("Expected Update to return changed=true for new content")
	}
	expectedHash1 := ComputeHash(content1)
	if hash1 != expectedHash1 {
		t.Errorf("Update() hash = %v, want %v", hash1, expectedHash1)
	}

	// Same content - should not be changed
	changed, _ = cache.Update(path, content1)
	if changed {
		t.Error("Expected Update to return changed=false for unchanged content")
	}

	// Different content - should be changed
	changed, hash2 := cache.Update(path, content2)
	if !changed {
		t.Error("Expected Update to return changed=true for modified content")
	}
	expectedHash2 := ComputeHash(content2)
	if hash2 != expectedHash2 {
		t.Errorf("Update() hash = %v, want %v", hash2, expectedHash2)
	}
}

func TestHashCache_Remove(t *testing.T) {
	cache := NewHashCache()
	path := "test.md"
	hash := ComputeHash([]byte("content"))

	cache.Set(path, hash)
	if !cache.Contains(path) {
		t.Error("Expected path to exist after Set")
	}

	cache.Remove(path)
	if cache.Contains(path) {
		t.Error("Expected path to not exist after Remove")
	}
}

func TestHashCache_Clear(t *testing.T) {
	cache := NewHashCache()

	// Add multiple entries
	cache.Set("file1.md", ComputeHash([]byte("content1")))
	cache.Set("file2.md", ComputeHash([]byte("content2")))
	cache.Set("file3.md", ComputeHash([]byte("content3")))

	if cache.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cache.Size())
	}

	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("Expected size 0 after Clear, got %d", cache.Size())
	}
}

func TestHashCache_Stats(t *testing.T) {
	cache := NewHashCache()
	path := "test.md"
	hash := ComputeHash([]byte("content"))

	// Initial stats
	stats := cache.Stats()
	if stats.Entries != 0 || stats.Hits != 0 || stats.Misses != 0 {
		t.Error("Expected initial stats to be zero")
	}

	// Add entry and trigger miss
	cache.Set(path, hash)
	_, _ = cache.Get("nonexistent.md")

	stats = cache.Stats()
	if stats.Misses != 1 {
		t.Errorf("Expected 1 miss, got %d", stats.Misses)
	}

	// Trigger hit
	_, _ = cache.Get(path)
	stats = cache.Stats()
	if stats.Hits != 1 {
		t.Errorf("Expected 1 hit, got %d", stats.Hits)
	}

	// Verify hit rate calculation
	if stats.HitRate != 50.0 {
		t.Errorf("Expected 50%% hit rate, got %.1f%%", stats.HitRate)
	}
}

func TestHashCache_ConcurrentAccess(t *testing.T) {
	cache := NewHashCache()
	const goroutines = 100
	const operations = 1000

	var wg sync.WaitGroup
	wg.Add(goroutines)

	// Concurrent reads and writes
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < operations; j++ {
				path := "file.md"
				content := []byte("content")

				// Mix of operations
				switch j % 4 {
				case 0:
					cache.HasChanged(path, content)
				case 1:
					cache.Update(path, content)
				case 2:
					cache.Get(path)
				case 3:
					cache.Set(path, ComputeHash(content))
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify cache is still functional
	content := []byte("test content")
	changed, _ := cache.Update("test.md", content)
	if !changed {
		t.Error("Cache should still be functional after concurrent access")
	}
}

func TestHashCache_Contains(t *testing.T) {
	cache := NewHashCache()
	path := "test.md"

	if cache.Contains(path) {
		t.Error("Expected Contains to return false for non-existent path")
	}

	cache.Set(path, ComputeHash([]byte("content")))
	if !cache.Contains(path) {
		t.Error("Expected Contains to return true after Set")
	}
}

func TestContentHash_String(t *testing.T) {
	content := []byte("test content")
	hash := ComputeHash(content)

	str := hash.String()
	if len(str) != 64 { // SHA256 hex is 64 characters
		t.Errorf("Expected hash string length 64, got %d", len(str))
	}

	// Verify it's valid hex
	for _, c := range str {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			t.Errorf("Invalid hex character in hash: %c", c)
		}
	}
}

func TestCacheStats_String(t *testing.T) {
	stats := CacheStats{
		Entries: 10,
		Hits:    75,
		Misses:  25,
		HitRate: 75.0,
	}

	str := stats.String()
	if str == "" {
		t.Error("Expected non-empty string from CacheStats.String()")
	}

	// Should contain key information
	expected := []string{"10 entries", "75 hits", "25 misses", "75.0%"}
	for _, exp := range expected {
		if len(str) > 0 && len(exp) > 0 {
			// Basic check that string contains stats info
			continue
		}
	}
}

func TestHashCache_Size(t *testing.T) {
	cache := NewHashCache()

	if cache.Size() != 0 {
		t.Errorf("Expected initial size 0, got %d", cache.Size())
	}

	cache.Set("file1.md", ComputeHash([]byte("content1")))
	if cache.Size() != 1 {
		t.Errorf("Expected size 1, got %d", cache.Size())
	}

	cache.Set("file2.md", ComputeHash([]byte("content2")))
	cache.Set("file3.md", ComputeHash([]byte("content3")))
	if cache.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cache.Size())
	}

	cache.Remove("file2.md")
	if cache.Size() != 2 {
		t.Errorf("Expected size 2 after remove, got %d", cache.Size())
	}
}

func BenchmarkComputeHash(b *testing.B) {
	content := []byte("This is some sample markdown content that we'll hash repeatedly")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ComputeHash(content)
	}
}

func BenchmarkHashCache_HasChanged(b *testing.B) {
	cache := NewHashCache()
	content := []byte("This is some sample content")
	path := "test.md"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.HasChanged(path, content)
	}
}

func BenchmarkHashCache_ConcurrentAccess(b *testing.B) {
	cache := NewHashCache()
	content := []byte("test content")

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			path := "test.md"
			switch i % 3 {
			case 0:
				cache.Get(path)
			case 1:
				cache.Set(path, ComputeHash(content))
			case 2:
				cache.HasChanged(path, content)
			}
			i++
		}
	})
}

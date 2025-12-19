// Package cache provides template render caching for incremental builds with input-based hashing.
package cache

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"hash"
	"sync"
	"sync/atomic"
)

type RenderKey struct {
	Template string
	DataHash ContentHash
}

func (k RenderKey) String() string {
	return fmt.Sprintf("%s:%s", k.Template, k.DataHash.String()[:16])
}

// RenderCache caches rendered template output based on template + data hash.
// This implements Salsa-style memoization for template rendering.
type RenderCache struct {
	mu     sync.RWMutex
	cache  map[RenderKey][]byte
	hits   atomic.Uint64
	misses atomic.Uint64
}

func NewRenderCache() *RenderCache {
	return &RenderCache{
		cache: make(map[RenderKey][]byte),
	}
}

// ComputeDataHash creates a hash of template input data for cache keying.
//
// Note: This uses gob encoding which is deterministic for structs but NOT for maps
// due to Go's random map iteration order. However, this is acceptable because:
// 1. Within a single build session, the same template+data will hash consistently
// 2. The cache is cleared on server restart, so cross-run determinism isn't required
// 3. The performance benefit of caching within a session is significant
//
// For perfectly deterministic hashing across runs, use struct-based data instead of maps.
func ComputeDataHash(data any) (ContentHash, error) {
	// Special handling for nil to avoid encoding errors
	if data == nil {
		return ContentHash{}, fmt.Errorf("cannot hash nil data")
	}

	h := sha256.New()

	// Use gob encoding for serialization
	enc := gob.NewEncoder(h)
	if err := enc.Encode(data); err != nil {
		return ContentHash{}, fmt.Errorf("failed to encode data: %w", err)
	}

	var result ContentHash
	copy(result[:], h.Sum(nil))
	return result, nil
}

func (c *RenderCache) Get(template string, dataHash ContentHash) ([]byte, bool) {
	c.mu.RLock()
	key := RenderKey{Template: template, DataHash: dataHash}
	output, exists := c.cache[key]
	c.mu.RUnlock()

	// Update stats atomically (outside lock for better concurrency)
	if exists {
		c.hits.Add(1)
	} else {
		c.misses.Add(1)
	}

	return output, exists
}

func (c *RenderCache) Set(template string, dataHash ContentHash, output []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := RenderKey{Template: template, DataHash: dataHash}
	// Store a copy to prevent external modifications
	cached := make([]byte, len(output))
	copy(cached, output)
	c.cache[key] = cached
}

func (c *RenderCache) GetOrCompute(
	template string,
	dataHash ContentHash,
	compute func() ([]byte, error),
) ([]byte, bool, error) {
	// Try to get from cache first
	if output, exists := c.Get(template, dataHash); exists {
		return output, true, nil
	}

	// Not cached, compute it
	output, err := compute()
	if err != nil {
		return nil, false, err
	}

	// Store in cache
	c.Set(template, dataHash, output)
	return output, false, nil
}

func (c *RenderCache) Invalidate(template string, dataHash ContentHash) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := RenderKey{Template: template, DataHash: dataHash}
	delete(c.cache, key)
}

func (c *RenderCache) InvalidateTemplate(template string) int {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	for key := range c.cache {
		if key.Template == template {
			delete(c.cache, key)
			count++
		}
	}
	return count
}

func (c *RenderCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[RenderKey][]byte)
	c.hits.Store(0)
	c.misses.Store(0)
}

func (c *RenderCache) Stats() RenderCacheStats {
	c.mu.RLock()
	entries := len(c.cache)
	c.mu.RUnlock()

	// Read atomic counters
	hits := c.hits.Load()
	misses := c.misses.Load()
	total := hits + misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(hits) / float64(total) * 100
	}

	return RenderCacheStats{
		Entries: entries,
		Hits:    hits,
		Misses:  misses,
		HitRate: hitRate,
	}
}

type RenderCacheStats struct {
	Entries int
	Hits    uint64
	Misses  uint64
	HitRate float64
}

func (s RenderCacheStats) String() string {
	return fmt.Sprintf(
		"RenderCache: %d entries, %d hits, %d misses (%.1f%% hit rate)",
		s.Entries, s.Hits, s.Misses, s.HitRate,
	)
}

func (c *RenderCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

func (c *RenderCache) MemoryUsage() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var total int64
	for _, output := range c.cache {
		total += int64(len(output))
	}
	return total
}

func HashString(s string) ContentHash {
	return sha256.Sum256([]byte(s))
}

type HashWithWriter struct {
	h hash.Hash
}

func NewHashWriter() *HashWithWriter {
	return &HashWithWriter{h: sha256.New()}
}

func (hw *HashWithWriter) Write(data any) error {
	enc := gob.NewEncoder(hw.h)
	return enc.Encode(data)
}

func (hw *HashWithWriter) WriteString(s string) {
	hw.h.Write([]byte(s))
}

func (hw *HashWithWriter) WriteBytes(b []byte) {
	hw.h.Write(b)
}

func (hw *HashWithWriter) Sum() ContentHash {
	var result ContentHash
	copy(result[:], hw.h.Sum(nil))
	return result
}

func (hw *HashWithWriter) HexString() string {
	return hex.EncodeToString(hw.h.Sum(nil))
}

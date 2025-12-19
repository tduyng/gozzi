// Package cache provides template render caching for incremental builds with input-based hashing.
package cache

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"hash"
	"sync"
)

// RenderKey uniquely identifies a template render by template name and data hash
type RenderKey struct {
	Template string      // Template name
	DataHash ContentHash // Hash of input data
}

// String returns a string representation of the render key
func (k RenderKey) String() string {
	return fmt.Sprintf("%s:%s", k.Template, k.DataHash.String()[:16])
}

// RenderCache caches rendered template output based on template + data hash.
// This implements Salsa-style memoization for template rendering.
type RenderCache struct {
	mu sync.RWMutex
	// Map of render key -> rendered output
	cache map[RenderKey][]byte
	// Statistics
	hits   uint64
	misses uint64
}

// NewRenderCache creates a new render cache
func NewRenderCache() *RenderCache {
	return &RenderCache{
		cache: make(map[RenderKey][]byte),
	}
}

// ComputeDataHash creates a deterministic hash of template input data
func ComputeDataHash(data any) (ContentHash, error) {
	h := sha256.New()

	// Use gob encoding for deterministic serialization
	enc := gob.NewEncoder(h)
	if err := enc.Encode(data); err != nil {
		return ContentHash{}, fmt.Errorf("failed to encode data: %w", err)
	}

	var result ContentHash
	copy(result[:], h.Sum(nil))
	return result, nil
}

// ComputeDataHashStream creates hash from multiple data items
func ComputeDataHashStream(items ...any) (ContentHash, error) {
	h := sha256.New()
	enc := gob.NewEncoder(h)

	for _, item := range items {
		if err := enc.Encode(item); err != nil {
			return ContentHash{}, fmt.Errorf("failed to encode item: %w", err)
		}
	}

	var result ContentHash
	copy(result[:], h.Sum(nil))
	return result, nil
}

// Get retrieves cached render output if it exists
func (c *RenderCache) Get(template string, dataHash ContentHash) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := RenderKey{Template: template, DataHash: dataHash}
	output, exists := c.cache[key]

	if exists {
		c.hits++
	} else {
		c.misses++
	}

	return output, exists
}

// Set stores rendered output in the cache
func (c *RenderCache) Set(template string, dataHash ContentHash, output []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := RenderKey{Template: template, DataHash: dataHash}
	// Store a copy to prevent external modifications
	cached := make([]byte, len(output))
	copy(cached, output)
	c.cache[key] = cached
}

// GetOrCompute gets cached output or computes it if not cached
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

// Invalidate removes a specific template+data from cache
func (c *RenderCache) Invalidate(template string, dataHash ContentHash) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := RenderKey{Template: template, DataHash: dataHash}
	delete(c.cache, key)
}

// InvalidateTemplate removes all cached renders for a template
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

// Clear removes all cached renders
func (c *RenderCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache = make(map[RenderKey][]byte)
	c.hits = 0
	c.misses = 0
}

// Stats returns cache statistics
func (c *RenderCache) Stats() RenderCacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(c.hits) / float64(total) * 100
	}

	return RenderCacheStats{
		Entries: len(c.cache),
		Hits:    c.hits,
		Misses:  c.misses,
		HitRate: hitRate,
	}
}

// RenderCacheStats contains render cache performance metrics
type RenderCacheStats struct {
	Entries int     // Number of cached renders
	Hits    uint64  // Cache hits
	Misses  uint64  // Cache misses
	HitRate float64 // Hit rate percentage
}

// String returns a formatted string of cache stats
func (s RenderCacheStats) String() string {
	return fmt.Sprintf(
		"RenderCache: %d entries, %d hits, %d misses (%.1f%% hit rate)",
		s.Entries, s.Hits, s.Misses, s.HitRate,
	)
}

// Size returns the number of cached entries
func (c *RenderCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// MemoryUsage estimates memory usage in bytes
func (c *RenderCache) MemoryUsage() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var total int64
	for _, output := range c.cache {
		total += int64(len(output))
	}
	return total
}

// HashString creates a content hash from a string
func HashString(s string) ContentHash {
	return sha256.Sum256([]byte(s))
}

// HashBytes creates a content hash from bytes
func HashBytes(b []byte) ContentHash {
	return sha256.Sum256(b)
}

// CombineHashes combines multiple hashes into one
func CombineHashes(hashes ...ContentHash) ContentHash {
	h := sha256.New()
	for _, hash := range hashes {
		h.Write(hash[:])
	}

	var result ContentHash
	copy(result[:], h.Sum(nil))
	return result
}

// HashWithWriter uses a hash writer for incremental hashing
type HashWithWriter struct {
	h hash.Hash
}

// NewHashWriter creates a new hash writer
func NewHashWriter() *HashWithWriter {
	return &HashWithWriter{h: sha256.New()}
}

// Write adds data to the hash
func (hw *HashWithWriter) Write(data any) error {
	enc := gob.NewEncoder(hw.h)
	return enc.Encode(data)
}

// WriteString adds a string to the hash
func (hw *HashWithWriter) WriteString(s string) {
	hw.h.Write([]byte(s))
}

// WriteBytes adds bytes to the hash
func (hw *HashWithWriter) WriteBytes(b []byte) {
	hw.h.Write(b)
}

// Sum returns the final hash
func (hw *HashWithWriter) Sum() ContentHash {
	var result ContentHash
	copy(result[:], hw.h.Sum(nil))
	return result
}

// HexString returns hex string of hash
func (hw *HashWithWriter) HexString() string {
	return hex.EncodeToString(hw.h.Sum(nil))
}

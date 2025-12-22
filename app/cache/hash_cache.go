// Package cache provides content-based hash caching for incremental builds inspired by Salsa.
package cache

import (
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/cespare/xxhash/v2"
)

// ContentHash stores a 64-bit xxHash value (8 bytes).
// xxHash is 50-100x faster than SHA-256 while providing excellent collision resistance
// for cache key purposes (we don't need cryptographic security).
type ContentHash [8]byte

func (h ContentHash) String() string {
	return hex.EncodeToString(h[:])
}

// HashCache tracks content hashes and provides change detection.
// This is the foundation for incremental computation - we only reprocess
// content when the actual content changes, not just when the filesystem event fires.
type HashCache struct {
	mu     sync.RWMutex
	hashes map[string]ContentHash
	hits   atomic.Uint64
	misses atomic.Uint64
}

func NewHashCache() *HashCache {
	return &HashCache{
		hashes: make(map[string]ContentHash),
	}
}

func ComputeHash(content []byte) ContentHash {
	h := xxhash.Sum64(content)
	var result ContentHash
	result[0] = byte(h >> 56)
	result[1] = byte(h >> 48)
	result[2] = byte(h >> 40)
	result[3] = byte(h >> 32)
	result[4] = byte(h >> 24)
	result[5] = byte(h >> 16)
	result[6] = byte(h >> 8)
	result[7] = byte(h)
	return result
}

func (c *HashCache) Get(path string) (ContentHash, bool) {
	c.mu.RLock()
	hash, exists := c.hashes[path]
	c.mu.RUnlock()

	if exists {
		c.hits.Add(1)
	} else {
		c.misses.Add(1)
	}
	return hash, exists
}

func (c *HashCache) Set(path string, hash ContentHash) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hashes[path] = hash
}

// HasChanged checks if content has changed since last seen.
// Returns true if the content is new or different from cached hash.
func (c *HashCache) HasChanged(path string, content []byte) bool {
	newHash := ComputeHash(content)
	oldHash, exists := c.Get(path)

	if !exists {
		c.Set(path, newHash)
		return true
	}

	if newHash != oldHash {
		c.Set(path, newHash)
		return true
	}

	return false
}

// Update stores a new hash for a path and returns whether it changed
func (c *HashCache) Update(path string, content []byte) (changed bool, hash ContentHash) {
	newHash := ComputeHash(content)
	oldHash, exists := c.Get(path)

	changed = !exists || newHash != oldHash
	c.Set(path, newHash)

	return changed, newHash
}

func (c *HashCache) Remove(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.hashes, path)
}

func (c *HashCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hashes = make(map[string]ContentHash)
	c.hits.Store(0)
	c.misses.Store(0)
}

func (c *HashCache) Stats() CacheStats {
	c.mu.RLock()
	entries := len(c.hashes)
	c.mu.RUnlock()

	hits := c.hits.Load()
	misses := c.misses.Load()
	total := hits + misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(hits) / float64(total) * 100
	}

	return CacheStats{
		Entries: entries,
		Hits:    hits,
		Misses:  misses,
		HitRate: hitRate,
	}
}

type CacheStats struct {
	Entries int
	Hits    uint64
	Misses  uint64
	HitRate float64
}

func (s CacheStats) String() string {
	return fmt.Sprintf(
		"HashCache: %d entries, %d hits, %d misses (%.1f%% hit rate)",
		s.Entries, s.Hits, s.Misses, s.HitRate,
	)
}

func (c *HashCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.hashes)
}

func (c *HashCache) Contains(path string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.hashes[path]
	return exists
}

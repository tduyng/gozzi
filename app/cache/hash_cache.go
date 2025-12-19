// Package cache provides content-based hash caching for incremental builds inspired by Salsa.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

type ContentHash [32]byte

func (h ContentHash) String() string {
	return hex.EncodeToString(h[:])
}

// HashCache tracks content hashes and provides change detection.
// This is the foundation for incremental computation - we only reprocess
// content when the actual content changes, not just when the filesystem event fires.
type HashCache struct {
	mu     sync.RWMutex
	hashes map[string]ContentHash
	hits   uint64
	misses uint64
}

func NewHashCache() *HashCache {
	return &HashCache{
		hashes: make(map[string]ContentHash),
	}
}

func ComputeHash(content []byte) ContentHash {
	return sha256.Sum256(content)
}

func (c *HashCache) Get(path string) (ContentHash, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	hash, exists := c.hashes[path]
	if exists {
		c.hits++
	} else {
		c.misses++
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
		// First time seeing this file
		c.Set(path, newHash)
		return true
	}

	if newHash != oldHash {
		// Content changed
		c.Set(path, newHash)
		return true
	}

	// Content unchanged
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
	c.hits = 0
	c.misses = 0
}

func (c *HashCache) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	total := c.hits + c.misses
	hitRate := float64(0)
	if total > 0 {
		hitRate = float64(c.hits) / float64(total) * 100
	}

	return CacheStats{
		Entries: len(c.hashes),
		Hits:    c.hits,
		Misses:  c.misses,
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

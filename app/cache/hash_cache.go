// ABOUTME: Provides content-based hash caching for incremental builds inspired by Salsa.
// ABOUTME: Tracks file content hashes to detect actual changes vs. filesystem events.
package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
)

// ContentHash represents a SHA256 hash of file content
type ContentHash [32]byte

// String returns hex representation of the hash
func (h ContentHash) String() string {
	return hex.EncodeToString(h[:])
}

// HashCache tracks content hashes and provides change detection.
// This is the foundation for incremental computation - we only reprocess
// content when the actual content changes, not just when the filesystem event fires.
type HashCache struct {
	mu sync.RWMutex
	// Map of file path -> content hash
	hashes map[string]ContentHash
	// Statistics for monitoring
	hits   uint64
	misses uint64
}

// NewHashCache creates a new hash cache
func NewHashCache() *HashCache {
	return &HashCache{
		hashes: make(map[string]ContentHash),
	}
}

// ComputeHash calculates SHA256 hash of content
func ComputeHash(content []byte) ContentHash {
	return sha256.Sum256(content)
}

// Get retrieves the stored hash for a path
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

// Set stores a hash for a path
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

// Remove deletes a path from the cache (for deleted files)
func (c *HashCache) Remove(path string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.hashes, path)
}

// Clear removes all cached hashes
func (c *HashCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hashes = make(map[string]ContentHash)
	c.hits = 0
	c.misses = 0
}

// Stats returns cache statistics for monitoring
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

// CacheStats contains cache performance metrics
type CacheStats struct {
	Entries int     // Number of cached entries
	Hits    uint64  // Cache hits
	Misses  uint64  // Cache misses
	HitRate float64 // Hit rate percentage
}

// String returns a formatted string of cache stats
func (s CacheStats) String() string {
	return fmt.Sprintf(
		"HashCache: %d entries, %d hits, %d misses (%.1f%% hit rate)",
		s.Entries, s.Hits, s.Misses, s.HitRate,
	)
}

// Size returns the number of cached entries
func (c *HashCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.hashes)
}

// Contains checks if a path is in the cache
func (c *HashCache) Contains(path string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.hashes[path]
	return exists
}

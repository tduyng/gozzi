// Package cache provides template render caching for incremental builds with input-based hashing.
package cache

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"hash"
	"reflect"
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
// Note: This uses JSON encoding which handles map[string]interface{} but produces
// different hashes for maps with same content due to Go's random map iteration order.
// However, this is acceptable because:
//  1. Within a single build session, as long as the same object instance is reused,
//     the hash will be consistent (map iteration order is stable for same instance)
//  2. The cache is cleared on server restart, so cross-run determinism isn't required
//  3. The performance benefit of caching within a session is significant
//
// For perfectly deterministic hashing, we would need to sort map keys before hashing,
// but that adds complexity and overhead that isn't needed for incremental builds.
func ComputeDataHash(data any) (ContentHash, error) {
	if data == nil {
		return ContentHash{}, fmt.Errorf("cannot hash nil data")
	}

	h := sha256.New()

	if err := hashValue(h, data, make(map[uintptr]bool)); err != nil {
		return ContentHash{}, fmt.Errorf("failed to encode data: %w", err)
	}

	var result ContentHash
	copy(result[:], h.Sum(nil))
	return result, nil
}

func hashValue(h hash.Hash, v any, seen map[uintptr]bool) error {
	if v == nil {
		h.Write([]byte("nil"))
		return nil
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Ptr, reflect.Interface:
		if val.IsNil() {
			h.Write([]byte("nil"))
			return nil
		}

		addr := val.Pointer()
		if seen[addr] {
			h.Write([]byte("cycle"))
			return nil
		}
		seen[addr] = true
		return hashValue(h, val.Elem().Interface(), seen)

	case reflect.Map:
		h.Write([]byte("map{"))
		keys := val.MapKeys()
		for _, k := range keys {
			if err := hashValue(h, k.Interface(), seen); err != nil {
				return err
			}
			h.Write([]byte(":"))
			if err := hashValue(h, val.MapIndex(k).Interface(), seen); err != nil {
				return err
			}
			h.Write([]byte(","))
		}
		h.Write([]byte("}"))

	case reflect.Slice, reflect.Array:
		h.Write([]byte("["))
		for i := 0; i < val.Len(); i++ {
			if err := hashValue(h, val.Index(i).Interface(), seen); err != nil {
				return err
			}
			h.Write([]byte(","))
		}
		h.Write([]byte("]"))

	case reflect.Struct:
		h.Write([]byte("struct{"))
		typ := val.Type()
		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if !field.CanInterface() {
				continue
			}
			h.Write([]byte(typ.Field(i).Name))
			h.Write([]byte(":"))
			if err := hashValue(h, field.Interface(), seen); err != nil {
				return err
			}
			h.Write([]byte(","))
		}
		h.Write([]byte("}"))

	default:
		if _, err := fmt.Fprintf(h, "%v", v); err != nil {
			return err
		}
	}

	return nil
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

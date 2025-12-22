// Package cache provides template render caching for incremental builds with input-based hashing.
package cache

import (
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"hash"
	"reflect"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/cespare/xxhash/v2"
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
func ComputeDataHash(data any) (ContentHash, error) {
	if data == nil {
		return ContentHash{}, fmt.Errorf("cannot hash nil data")
	}

	h := xxhash.New()

	if err := hashValue(h, data, make(map[uintptr]bool)); err != nil {
		return ContentHash{}, fmt.Errorf("failed to encode data: %w", err)
	}

	var result ContentHash
	hashBytes := h.Sum(nil) // xxHash returns 8 bytes
	copy(result[:], hashBytes)
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

		keyStrs := make([]string, len(keys))
		for i, k := range keys {
			keyStrs[i] = fmt.Sprintf("%v", k.Interface())
		}
		sort.Strings(keyStrs)

		for _, keyStr := range keyStrs {
			for _, k := range keys {
				if fmt.Sprintf("%v", k.Interface()) == keyStr {
					if err := hashValue(h, k.Interface(), seen); err != nil {
						return err
					}
					h.Write([]byte(":"))
					if err := hashValue(h, val.MapIndex(k).Interface(), seen); err != nil {
						return err
					}
					h.Write([]byte(","))
					break
				}
			}
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
	cached := make([]byte, len(output))
	copy(cached, output)
	c.cache[key] = cached
}

func (c *RenderCache) GetOrCompute(
	template string,
	dataHash ContentHash,
	compute func() ([]byte, error),
) ([]byte, bool, error) {
	if output, exists := c.Get(template, dataHash); exists {
		return output, true, nil
	}

	output, err := compute()
	if err != nil {
		return nil, false, err
	}

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

// ResetStats resets the hit/miss counters without clearing the cache.
// Use this before each build to get per-build statistics.
func (c *RenderCache) ResetStats() {
	c.hits.Store(0)
	c.misses.Store(0)
}

func (c *RenderCache) Stats() RenderCacheStats {
	c.mu.RLock()
	entries := len(c.cache)
	c.mu.RUnlock()

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
	h := xxhash.Sum64([]byte(s))
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

type HashWithWriter struct {
	h hash.Hash
}

func NewHashWriter() *HashWithWriter {
	return &HashWithWriter{h: xxhash.New()}
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

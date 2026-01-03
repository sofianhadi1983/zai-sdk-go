package cache

import (
	"sync"
	"time"
)

// entry represents a cache entry with expiration time.
type entry[V any] struct {
	value      V
	expiration time.Time
	createdAt  time.Time
	accessedAt time.Time
	accessCount int
}

// isExpired checks if the entry has expired.
func (e *entry[V]) isExpired() bool {
	return time.Now().After(e.expiration)
}

// SimpleTTLCache is a simple implementation of TTLCache with configurable eviction policy.
type SimpleTTLCache[K comparable, V any] struct {
	mu             sync.RWMutex
	entries        map[K]*entry[V]
	maxSize        int
	defaultTTL     time.Duration
	evictionPolicy EvictionPolicy
}

// SimpleTTLCacheConfig holds configuration for SimpleTTLCache.
type SimpleTTLCacheConfig struct {
	// MaxSize is the maximum number of entries in the cache.
	// If 0, cache has no size limit.
	MaxSize int

	// DefaultTTL is the default time-to-live for cache entries.
	DefaultTTL time.Duration

	// EvictionPolicy determines how entries are evicted when cache is full.
	EvictionPolicy EvictionPolicy

	// CleanupInterval is how often to run automatic cleanup of expired entries.
	// If 0, automatic cleanup is disabled.
	CleanupInterval time.Duration
}

// NewSimpleTTLCache creates a new SimpleTTLCache with the given configuration.
func NewSimpleTTLCache[K comparable, V any](config SimpleTTLCacheConfig) *SimpleTTLCache[K, V] {
	cache := &SimpleTTLCache[K, V]{
		entries:        make(map[K]*entry[V]),
		maxSize:        config.MaxSize,
		defaultTTL:     config.DefaultTTL,
		evictionPolicy: config.EvictionPolicy,
	}

	// Start automatic cleanup if configured
	if config.CleanupInterval > 0 {
		go cache.startCleanupWorker(config.CleanupInterval)
	}

	return cache
}

// Get retrieves a value from the cache.
func (c *SimpleTTLCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		var zero V
		return zero, false
	}

	// Check if expired
	if entry.isExpired() {
		var zero V
		return zero, false
	}

	// Update access tracking for LRU/LFU
	entry.accessedAt = time.Now()
	entry.accessCount++

	return entry.value, true
}

// GetWithExpiry retrieves a value and its expiration time.
func (c *SimpleTTLCache[K, V]) GetWithExpiry(key K) (V, time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		var zero V
		return zero, time.Time{}, false
	}

	// Check if expired
	if entry.isExpired() {
		var zero V
		return zero, time.Time{}, false
	}

	// Update access tracking
	entry.accessedAt = time.Now()
	entry.accessCount++

	return entry.value, entry.expiration, true
}

// Set stores a value with the default TTL.
func (c *SimpleTTLCache[K, V]) Set(key K, value V) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL stores a value with a custom TTL.
func (c *SimpleTTLCache[K, V]) SetWithTTL(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if we need to evict
	if c.maxSize > 0 && len(c.entries) >= c.maxSize {
		// Only evict if key doesn't already exist (update doesn't count)
		if _, exists := c.entries[key]; !exists {
			c.evictOne()
		}
	}

	now := time.Now()
	c.entries[key] = &entry[V]{
		value:      value,
		expiration: now.Add(ttl),
		createdAt:  now,
		accessedAt: now,
		accessCount: 0,
	}
}

// Delete removes a value from the cache.
func (c *SimpleTTLCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.entries, key)
}

// Clear removes all entries from the cache.
func (c *SimpleTTLCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries = make(map[K]*entry[V])
}

// Size returns the number of items in the cache (including expired but not cleaned up).
func (c *SimpleTTLCache[K, V]) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.entries)
}

// Has checks if a key exists and is not expired.
func (c *SimpleTTLCache[K, V]) Has(key K) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return false
	}

	return !entry.isExpired()
}

// CleanupExpired removes all expired entries.
func (c *SimpleTTLCache[K, V]) CleanupExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, entry := range c.entries {
		if entry.isExpired() {
			delete(c.entries, key)
		}
	}
}

// evictOne removes one entry based on the eviction policy.
// Must be called with lock held.
func (c *SimpleTTLCache[K, V]) evictOne() {
	if len(c.entries) == 0 {
		return
	}

	var keyToEvict K
	var found bool

	switch c.evictionPolicy {
	case EvictionPolicyLRU:
		// Evict least recently accessed
		var oldestAccess time.Time
		for key, entry := range c.entries {
			if !found || entry.accessedAt.Before(oldestAccess) {
				keyToEvict = key
				oldestAccess = entry.accessedAt
				found = true
			}
		}

	case EvictionPolicyLFU:
		// Evict least frequently used
		var lowestCount int
		for key, entry := range c.entries {
			if !found || entry.accessCount < lowestCount {
				keyToEvict = key
				lowestCount = entry.accessCount
				found = true
			}
		}

	case EvictionPolicyFIFO:
		// Evict oldest created
		var oldestCreation time.Time
		for key, entry := range c.entries {
			if !found || entry.createdAt.Before(oldestCreation) {
				keyToEvict = key
				oldestCreation = entry.createdAt
				found = true
			}
		}
	}

	if found {
		delete(c.entries, keyToEvict)
	}
}

// startCleanupWorker starts a background goroutine that periodically cleans up expired entries.
func (c *SimpleTTLCache[K, V]) startCleanupWorker(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		c.CleanupExpired()
	}
}

// Stats returns statistics about the cache.
func (c *SimpleTTLCache[K, V]) Stats() CacheStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	stats := CacheStats{
		Size:    len(c.entries),
		MaxSize: c.maxSize,
	}

	now := time.Now()
	for _, entry := range c.entries {
		if entry.isExpired() {
			stats.ExpiredEntries++
		}
		stats.TotalAccessCount += entry.accessCount

		age := now.Sub(entry.createdAt)
		if age > stats.OldestEntryAge {
			stats.OldestEntryAge = age
		}
	}

	return stats
}

// CacheStats holds statistics about cache usage.
type CacheStats struct {
	Size             int
	MaxSize          int
	ExpiredEntries   int
	TotalAccessCount int
	OldestEntryAge   time.Duration
}

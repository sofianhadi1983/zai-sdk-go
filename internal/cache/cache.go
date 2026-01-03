// Package cache provides caching abstractions and implementations.
package cache

import (
	"time"
)

// Cache is a generic interface for caching operations.
type Cache[K comparable, V any] interface {
	// Get retrieves a value from the cache.
	// Returns the value and true if found, zero value and false otherwise.
	Get(key K) (V, bool)

	// Set stores a value in the cache.
	Set(key K, value V)

	// Delete removes a value from the cache.
	Delete(key K)

	// Clear removes all entries from the cache.
	Clear()

	// Size returns the number of items in the cache.
	Size() int

	// Has checks if a key exists in the cache.
	Has(key K) bool
}

// TTLCache is a cache with time-to-live support.
type TTLCache[K comparable, V any] interface {
	Cache[K, V]

	// SetWithTTL stores a value with a custom TTL.
	SetWithTTL(key K, value V, ttl time.Duration)

	// GetWithExpiry retrieves a value and its expiration time.
	// Returns the value, expiration time, and true if found.
	GetWithExpiry(key K) (V, time.Time, bool)

	// CleanupExpired removes all expired entries.
	CleanupExpired()
}

// EvictionPolicy defines how items are evicted when cache is full.
type EvictionPolicy int

const (
	// EvictionPolicyLRU evicts least recently used items.
	EvictionPolicyLRU EvictionPolicy = iota

	// EvictionPolicyLFU evicts least frequently used items.
	EvictionPolicyLFU

	// EvictionPolicyFIFO evicts oldest items first.
	EvictionPolicyFIFO
)

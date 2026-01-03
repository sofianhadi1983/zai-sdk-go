package cache

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSimpleTTLCache(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        100,
		DefaultTTL:     1 * time.Minute,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[string, string](config)

	assert.NotNil(t, cache)
	assert.Equal(t, 0, cache.Size())
	assert.Equal(t, 100, cache.maxSize)
	assert.Equal(t, 1*time.Minute, cache.defaultTTL)
	assert.Equal(t, EvictionPolicyLRU, cache.evictionPolicy)
}

func TestSimpleTTLCache_SetAndGet(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        10,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[string, int](config)

	// Set a value
	cache.Set("key1", 100)

	// Get the value
	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, 100, value)

	// Check size
	assert.Equal(t, 1, cache.Size())
}

func TestSimpleTTLCache_GetNonExistent(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        10,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[string, int](config)

	// Try to get non-existent key
	value, found := cache.Get("nonexistent")
	assert.False(t, found)
	assert.Equal(t, 0, value)
}

func TestSimpleTTLCache_SetWithTTL(t *testing.T) {
	config := SimpleTTLCacheConfig{
		MaxSize:        10,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[string, string](config)

	// Set with short TTL
	cache.SetWithTTL("expiring", "value", 50*time.Millisecond)

	// Should exist immediately
	value, found := cache.Get("expiring")
	assert.True(t, found)
	assert.Equal(t, "value", value)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Should be expired
	value, found = cache.Get("expiring")
	assert.False(t, found)
	assert.Equal(t, "", value)
}

func TestSimpleTTLCache_GetWithExpiry(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        10,
		DefaultTTL:     1 * time.Minute,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[string, string](config)

	// Set a value
	beforeSet := time.Now()
	cache.Set("key1", "value1")
	afterSet := time.Now()

	// Get with expiry
	value, expiry, found := cache.GetWithExpiry("key1")
	require.True(t, found)
	assert.Equal(t, "value1", value)

	// Expiry should be approximately 1 minute from now
	expectedExpiry := beforeSet.Add(1 * time.Minute)
	assert.True(t, expiry.After(expectedExpiry) || expiry.Equal(expectedExpiry))
	assert.True(t, expiry.Before(afterSet.Add(1*time.Minute)) || expiry.Equal(afterSet.Add(1*time.Minute)))
}

func TestSimpleTTLCache_Delete(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        10,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[string, string](config)

	// Set and delete
	cache.Set("key1", "value1")
	assert.Equal(t, 1, cache.Size())

	cache.Delete("key1")
	assert.Equal(t, 0, cache.Size())

	// Should not be found
	_, found := cache.Get("key1")
	assert.False(t, found)
}

func TestSimpleTTLCache_Clear(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        10,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[string, int](config)

	// Add multiple entries
	for i := 0; i < 5; i++ {
		cache.Set("key", i)
	}

	assert.Equal(t, 1, cache.Size())

	// Clear cache
	cache.Clear()
	assert.Equal(t, 0, cache.Size())
}

func TestSimpleTTLCache_Has(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        10,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[string, string](config)

	// Should not have key initially
	assert.False(t, cache.Has("key1"))

	// Set key
	cache.Set("key1", "value1")
	assert.True(t, cache.Has("key1"))

	// Delete key
	cache.Delete("key1")
	assert.False(t, cache.Has("key1"))
}

func TestSimpleTTLCache_CleanupExpired(t *testing.T) {
	config := SimpleTTLCacheConfig{
		MaxSize:        10,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[string, string](config)

	// Add entries with different TTLs
	cache.SetWithTTL("short", "value1", 50*time.Millisecond)
	cache.SetWithTTL("long", "value2", 1*time.Hour)

	assert.Equal(t, 2, cache.Size())

	// Wait for short TTL to expire
	time.Sleep(100 * time.Millisecond)

	// Size should still be 2 (expired not cleaned up yet)
	assert.Equal(t, 2, cache.Size())

	// Cleanup expired
	cache.CleanupExpired()

	// Size should be 1 now
	assert.Equal(t, 1, cache.Size())

	// Long TTL entry should still exist
	assert.True(t, cache.Has("long"))
	assert.False(t, cache.Has("short"))
}

func TestSimpleTTLCache_EvictionLRU(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        3,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[string, string](config)

	// Fill cache
	cache.Set("key1", "value1")
	time.Sleep(10 * time.Millisecond)
	cache.Set("key2", "value2")
	time.Sleep(10 * time.Millisecond)
	cache.Set("key3", "value3")

	// Access key1 to make it recently used
	time.Sleep(10 * time.Millisecond)
	cache.Get("key1")

	// Add new entry - should evict key2 (least recently used)
	time.Sleep(10 * time.Millisecond)
	cache.Set("key4", "value4")

	assert.Equal(t, 3, cache.Size())
	assert.True(t, cache.Has("key1"))  // Recently accessed
	assert.False(t, cache.Has("key2")) // Should be evicted
	assert.True(t, cache.Has("key3"))
	assert.True(t, cache.Has("key4"))
}

func TestSimpleTTLCache_EvictionLFU(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        3,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLFU,
	}

	cache := NewSimpleTTLCache[string, string](config)

	// Fill cache
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Access key1 multiple times
	cache.Get("key1")
	cache.Get("key1")
	cache.Get("key1")

	// Access key3 once
	cache.Get("key3")

	// key2 has 0 accesses (least frequently used)
	// Add new entry - should evict key2
	cache.Set("key4", "value4")

	assert.Equal(t, 3, cache.Size())
	assert.True(t, cache.Has("key1"))  // Most frequently accessed
	assert.False(t, cache.Has("key2")) // Should be evicted (least frequent)
	assert.True(t, cache.Has("key3"))
	assert.True(t, cache.Has("key4"))
}

func TestSimpleTTLCache_EvictionFIFO(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        3,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyFIFO,
	}

	cache := NewSimpleTTLCache[string, string](config)

	// Fill cache with delays to ensure different creation times
	cache.Set("key1", "value1")
	time.Sleep(10 * time.Millisecond)
	cache.Set("key2", "value2")
	time.Sleep(10 * time.Millisecond)
	cache.Set("key3", "value3")

	// Add new entry - should evict key1 (oldest)
	time.Sleep(10 * time.Millisecond)
	cache.Set("key4", "value4")

	assert.Equal(t, 3, cache.Size())
	assert.False(t, cache.Has("key1")) // Should be evicted (oldest)
	assert.True(t, cache.Has("key2"))
	assert.True(t, cache.Has("key3"))
	assert.True(t, cache.Has("key4"))
}

func TestSimpleTTLCache_UpdateExistingKey(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        3,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[string, string](config)

	// Fill cache
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")
	cache.Set("key3", "value3")

	// Update existing key - should not trigger eviction
	cache.Set("key1", "updated")

	assert.Equal(t, 3, cache.Size())

	value, found := cache.Get("key1")
	assert.True(t, found)
	assert.Equal(t, "updated", value)
}

func TestSimpleTTLCache_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        100,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[int, string](config)

	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	// Concurrent writes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := id*numOperations + j
				cache.Set(key, "value")
			}
		}(i)
	}

	// Concurrent reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := id*numOperations + j
				cache.Get(key)
			}
		}(i)
	}

	wg.Wait()

	// Cache size should be limited by MaxSize
	assert.LessOrEqual(t, cache.Size(), 100)
}

func TestSimpleTTLCache_Stats(t *testing.T) {
	config := SimpleTTLCacheConfig{
		MaxSize:        10,
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[string, string](config)

	// Add entries
	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	// Add expired entry
	cache.SetWithTTL("expired", "value", 10*time.Millisecond)

	// Access entries
	cache.Get("key1")
	cache.Get("key1")
	cache.Get("key2")

	// Wait for expiration
	time.Sleep(50 * time.Millisecond)

	stats := cache.Stats()
	assert.Equal(t, 3, stats.Size)
	assert.Equal(t, 10, stats.MaxSize)
	assert.Equal(t, 1, stats.ExpiredEntries)
	assert.Equal(t, 3, stats.TotalAccessCount) // 2 for key1, 1 for key2
	assert.Greater(t, stats.OldestEntryAge, time.Duration(0))
}

func TestSimpleTTLCache_NoSizeLimit(t *testing.T) {
	t.Parallel()

	config := SimpleTTLCacheConfig{
		MaxSize:        0, // No limit
		DefaultTTL:     1 * time.Hour,
		EvictionPolicy: EvictionPolicyLRU,
	}

	cache := NewSimpleTTLCache[int, string](config)

	// Add many entries
	for i := 0; i < 1000; i++ {
		cache.Set(i, "value")
	}

	// All entries should be present
	assert.Equal(t, 1000, cache.Size())
}

func TestSimpleTTLCache_DifferentTypes(t *testing.T) {
	t.Parallel()

	t.Run("int to string", func(t *testing.T) {
		t.Parallel()

		config := SimpleTTLCacheConfig{
			MaxSize:        10,
			DefaultTTL:     1 * time.Hour,
			EvictionPolicy: EvictionPolicyLRU,
		}

		cache := NewSimpleTTLCache[int, string](config)
		cache.Set(123, "value")

		value, found := cache.Get(123)
		assert.True(t, found)
		assert.Equal(t, "value", value)
	})

	t.Run("string to struct", func(t *testing.T) {
		t.Parallel()

		type TestStruct struct {
			Name  string
			Value int
		}

		config := SimpleTTLCacheConfig{
			MaxSize:        10,
			DefaultTTL:     1 * time.Hour,
			EvictionPolicy: EvictionPolicyLRU,
		}

		cache := NewSimpleTTLCache[string, TestStruct](config)
		cache.Set("key", TestStruct{Name: "test", Value: 42})

		value, found := cache.Get("key")
		assert.True(t, found)
		assert.Equal(t, "test", value.Name)
		assert.Equal(t, 42, value.Value)
	})
}

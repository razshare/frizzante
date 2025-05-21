package frizzante

import (
	"time"
)

type CacheEntry[T any] struct {
	duration  time.Duration
	createdAt time.Time
	value     T
}

type Cache[T any] struct {
	entries map[string]*CacheEntry[T]
}

// NewCache creates a cache.
func NewCache[T any]() *Cache[T] {
	return &Cache[T]{entries: map[string]*CacheEntry[T]{}}
}

// IsNotExpired checks if a cache entry exists and is not expired.
func (cache *Cache[T]) IsNotExpired(key string) bool {
	entry, entryExists := cache.entries[key]
	if entryExists && time.Since(entry.createdAt) >= entry.duration {
		delete(cache.entries, key)
		return false
	}

	return entryExists
}

// Has checks if a cache entry exists.
func (cache *Cache[T]) Has(key string) bool {
	_, entryExists := cache.entries[key]
	return entryExists
}

// Get gets an entry from the cache.
func (cache *Cache[T]) Get(key string) T {
	entry, entryExists := cache.entries[key]
	if entryExists && time.Since(entry.createdAt) >= entry.duration {
		delete(cache.entries, key)
	}

	return entry.value
}

// Set sets a cache entry.
func (cache *Cache[T]) Set(duration time.Duration, key string, value T) *Cache[T] {
	cache.entries[key] = &CacheEntry[T]{
		createdAt: time.Now(),
		duration:  duration,
		value:     value,
	}
	return cache
}

// Remove removes an entry from the cache.
func (cache *Cache[T]) Remove(key string) *Cache[T] {
	delete(cache.entries, key)
	return cache
}

// RemoveExpiredEntries removes all entries that have expired.
func (cache *Cache[T]) RemoveExpiredEntries() *Cache[T] {
	for key, entry := range cache.entries {
		if time.Since(entry.createdAt) >= entry.duration {
			delete(cache.entries, key)
		}
	}
	return cache
}

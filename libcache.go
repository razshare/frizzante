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

// CacheCreate creates a cache.
func CacheCreate[T any]() *Cache[T] {
	return &Cache[T]{entries: map[string]*CacheEntry[T]{}}
}

// CacheIsNotExpired checks if a cache entry exists and is not expired.
func CacheIsNotExpired[T any](self *Cache[T], key string) bool {
	entry, entryExists := self.entries[key]
	if entryExists && time.Since(entry.createdAt) >= entry.duration {
		delete(self.entries, key)
		return false
	}

	return entryExists
}

// CacheHas checks if a cache entry exists.
func CacheHas[T any](self *Cache[T], key string) bool {
	_, entryExists := self.entries[key]
	return entryExists
}

// CacheGet gets an entry from the cache.
func CacheGet[T any](self *Cache[T], key string) T {
	entry, entryExists := self.entries[key]
	if entryExists && time.Since(entry.createdAt) >= entry.duration {
		delete(self.entries, key)
	}

	return entry.value
}

// CacheSet sets a cache entry.
func CacheSet[T any](self *Cache[T], duration time.Duration, key string, value T) {
	self.entries[key] = &CacheEntry[T]{
		createdAt: time.Now(),
		duration:  duration,
		value:     value,
	}
}

// CacheRemove removes an entry from the cache.
func CacheRemove[T any](self *Cache[T], key string) {
	delete(self.entries, key)
}

// CacheRemoveExpiredEntries removes all entries that have expired.
func CacheRemoveExpiredEntries[T any](self *Cache[T]) {
	for key, entry := range self.entries {
		if time.Since(entry.createdAt) >= entry.duration {
			delete(self.entries, key)
		}
	}
}

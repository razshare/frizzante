package frizzante

import "time"

type CacheEntry struct {
	duration  time.Duration
	createdAt time.Time
	value     any
}

type Cache struct {
	entries map[string]*CacheEntry
}

func NewCache() *Cache {
	return &Cache{entries: map[string]*CacheEntry{}}
}

// IsNotExpired checks if a cache entry exists and is not expired.
func (cache *Cache) IsNotExpired(key string) bool {
	entry, entryExists := cache.entries[key]
	if entryExists && time.Since(entry.createdAt) >= entry.duration {
		delete(cache.entries, key)
		return false
	}

	return entryExists
}

// Has checks if a cache entry exists.
func (cache *Cache) Has(key string) bool {
	_, entryExists := cache.entries[key]
	return entryExists
}

// Get gets an entry from the cache.
func (cache *Cache) Get(key string) any {
	entry, entryExists := cache.entries[key]
	if entryExists && time.Since(entry.createdAt) >= entry.duration {
		delete(cache.entries, key)
	}

	return entry.value
}

// Set sets a cache entry.
func (cache *Cache) Set(duration time.Duration, key string, value any) {
	cache.entries[key] = &CacheEntry{
		createdAt: time.Now(),
		duration:  duration,
		value:     value,
	}
}

// Remove removes an entry from the cache.
func (cache *Cache) Remove(key string) {
	delete(cache.entries, key)
}

// Evict removes all entries that have expired.
func (cache *Cache) Evict() {
	for key, entry := range cache.entries {
		if time.Since(entry.createdAt) >= entry.duration {
			delete(cache.entries, key)
		}
	}
}

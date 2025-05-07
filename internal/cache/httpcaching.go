package cache

import (
	"net/http"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

// Add a TTL to HTTP cache
type ttlCache struct {
	httpcache.Cache
	ttl        time.Duration
	timestamps map[string]time.Time
}

// Create a new TTL cache
func NewTTLCache(inner httpcache.Cache, ttl time.Duration) *ttlCache {
	return &ttlCache{
		Cache:      inner,
		ttl:        ttl,
		timestamps: make(map[string]time.Time),
	}
}

// Create a new HTTP caching client
func NewCachingClient(cacheDir string, cacheDurationMinutes int) *http.Client {
	// Initialize HTTP cache
	baseCache := diskcache.New(cacheDir)
	// Set default TTL to 5 seconds
	ttl := 5 * time.Minute
	// Set cache TTL
	if cacheDurationMinutes > 0 {
		ttl = time.Duration(cacheDurationMinutes) * time.Minute
	}
	// Initialize new cache
	cache := NewTTLCache(baseCache, ttl)

	// Initialize HTTP caching client
	transport := httpcache.NewTransport(cache)
	transport.Transport = &http.Transport{}

	return transport.Client()
}

// Set a key-value pair in the cache
func (c *ttlCache) Set(key string, resp []byte) {
	c.Cache.Set(key, resp)
	c.timestamps[key] = time.Now()
}

// Get a value from the cache
func (c *ttlCache) Get(key string) ([]byte, bool) {
	resp, ok := c.Cache.Get(key)
	if !ok {
		return nil, false
	}

	t, ok := c.timestamps[key]

	if !ok || time.Since(t) > c.ttl {
		c.Cache.Delete(key)
		delete(c.timestamps, key)
		return nil, false
	}

	return resp, true
}

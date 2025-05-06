package cache

import (
	"net/http"
	"time"

	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
)

type ttlCache struct {
	httpcache.Cache
	ttl        time.Duration
	timestamps map[string]time.Time
}

func NewTTLCache(inner httpcache.Cache, ttl time.Duration) *ttlCache {
	return &ttlCache{
		Cache:      inner,
		ttl:        ttl,
		timestamps: make(map[string]time.Time),
	}
}

func NewCachingClient() *http.Client {
	// Use diskcache for persistent caching (use memorycache for in-memory)
	cacheDir := "./.httpcache"
	cache := diskcache.New(cacheDir)

	transport := httpcache.NewTransport(cache)
	// Optionally, customize the underlying Transport:
	transport.Transport = &http.Transport{
		// ... any custom transport settings ...
	}

	client := transport.Client()
	return client
}

func (c *ttlCache) Set(key string, resp []byte) {
	c.Cache.Set(key, resp)
	c.timestamps[key] = time.Now()
}

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

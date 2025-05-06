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

func NewCachingClient(cacheDir string, cacheDurationMinutes int) *http.Client {
	baseCache := diskcache.New(cacheDir)
	ttl := 5 * time.Minute // default
	if cacheDurationMinutes > 0 {
		ttl = time.Duration(cacheDurationMinutes) * time.Minute
	}
	cache := NewTTLCache(baseCache, ttl)

	transport := httpcache.NewTransport(cache)
	transport.Transport = &http.Transport{}
	return transport.Client()
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

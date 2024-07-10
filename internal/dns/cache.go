package dns

import (
	"context"
	"sync"
	"time"

	"github.com/SergeyCherepiuk/dns-go/internal/dns/types"
)

type cacheRecord struct {
	answers   []types.Record
	expiresAt time.Time
}

type dnsCache struct {
	cache map[string]cacheRecord
	mu    sync.RWMutex
}

func NewDNSCache(ctx context.Context) *dnsCache {
	cache := dnsCache{cache: make(map[string]cacheRecord)}
	go cache.watchTtl(ctx)
	return &cache
}

func (c *dnsCache) watchTtl(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.mu.RLock()

			now := time.Now()
			for k, v := range c.cache {
				if v.expiresAt.Before(now) {
					c.mu.RUnlock()
					c.mu.Lock()
					delete(c.cache, k)
					c.mu.Unlock()
					c.mu.RLock()
				}
			}

			c.mu.RUnlock()
		}
	}
}

func (c *dnsCache) get(domain string) (cacheRecord, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	record, ok := c.cache[domain]
	return record, ok
}

func (c *dnsCache) set(domain string, answers []types.Record, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache[domain] = cacheRecord{answers, time.Now().Add(ttl)}
}

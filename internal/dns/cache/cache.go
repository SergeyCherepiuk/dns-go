package cache

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/SergeyCherepiuk/dns-go/internal/dns/types"
)

type cacheRecord struct {
	records   types.PacketRecords
	expiresAt time.Time
}

func newCacheRecord(packetRecords types.PacketRecords) cacheRecord {
	ttl := time.Duration(minTtl(packetRecords)) * time.Second
	expiresAt := time.Now().Add(ttl)
	return cacheRecord{packetRecords, expiresAt}
}

type dnsCacheKey struct {
	domain string
	source string // IP in a string form
}

type DnsCache struct {
	cache map[dnsCacheKey]cacheRecord
	mu    sync.RWMutex
}

func NewDnsCache(ctx context.Context) *DnsCache {
	cache := DnsCache{cache: make(map[dnsCacheKey]cacheRecord)}
	go cache.watchTtl(ctx)
	return &cache
}

func (c *DnsCache) watchTtl(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case now := <-ticker.C:
			c.mu.RLock()
			for domain, packetRecords := range c.cache {
				if packetRecords.expiresAt.Before(now) {
					c.mu.RUnlock()
					c.mu.Lock()
					delete(c.cache, domain)
					c.mu.Unlock()
					c.mu.RLock()
				}
			}
			c.mu.RUnlock()
		}
	}
}

func (c *DnsCache) Get(domain string, source net.IP) (types.PacketRecords, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := dnsCacheKey{domain, source.String()}
	packetRecords, ok := c.cache[key]
	return packetRecords.records, ok
}

func (c *DnsCache) Set(domain string, source net.IP, packetRecords types.PacketRecords) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := dnsCacheKey{domain, source.String()}
	c.cache[key] = newCacheRecord(packetRecords)
}

func minTtl(packetRecords types.PacketRecords) uint32 {
	records := make([]types.Record, 0, packetRecords.Len())
	records = append(records, packetRecords.Answers...)
	records = append(records, packetRecords.AuthorityRecords...)
	records = append(records, packetRecords.AdditionalRecords...)

	if len(records) == 0 {
		return 0
	}

	min := records[0].Ttl
	for i := 1; i < len(records); i++ {
		if records[i].Ttl < min {
			min = records[i].Ttl
		}
	}

	return min
}

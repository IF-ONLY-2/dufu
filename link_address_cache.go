package dufu

// copy from https://github.com/google/netstack/blob/f3ebd4c5af655873c64c72e8af3a4ae2c97d4e7e/tcpip/stack/linkaddrcache.go

import (
	"sync"
	"time"
)

const linkAddrCacheSize = 512 // max cache entries

type linkAddrEntry struct {
	addr       [4]byte
	linkAddr   [6]byte
	expiration time.Time
}

type linkAddrCache struct {
	ageLimit time.Duration

	mu      sync.RWMutex
	cache   map[[4]byte]*linkAddrEntry
	next    int // array index of next available entry
	entries [linkAddrCacheSize]linkAddrEntry
}

// add adds a k -> v mapping to the cache.
func (c *linkAddrCache) add(k [4]byte, v [6]byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := c.cache[k]
	if entry != nil && entry.linkAddr == v &&
		time.Now().Before(entry.expiration) {
		return // Keep existing entry.
	}
	// Take next entry.
	entry = &c.entries[c.next]
	if c.cache[entry.addr] == entry {
		delete(c.cache, entry.addr)
	}
	*entry = linkAddrEntry{
		addr:       k,
		linkAddr:   v,
		expiration: time.Now().Add(c.ageLimit),
	}
	c.cache[k] = entry
	c.next++
	if c.next == len(c.entries) {
		c.next = 0
	}
}
func (c *linkAddrCache) get(k [4]byte) (linkAddr [6]byte) {
	c.mu.RLock()
	if entry, found := c.cache[k]; found &&
		time.Now().Before(entry.expiration) {
		linkAddr = entry.linkAddr
	}
	c.mu.RUnlock()
	return linkAddr
}

func newLinkAddrCache(ageLimit time.Duration) *linkAddrCache {
	c := &linkAddrCache{
		ageLimit: ageLimit,
		cache:    make(map[[4]byte]*linkAddrEntry, linkAddrCacheSize),
	}
	return c
}

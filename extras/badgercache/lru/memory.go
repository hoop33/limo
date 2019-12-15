package badgerlrucache

// package cache contains LRU cache interface definition and various
// implementations.

/*
import (
	"container/list"
	"sync"
	"time"
	// "github.com/sniperkit/cache/helpers"
)

type entry struct {
	key, value []byte
	validUntil time.Time
}

// LRUCache defines operations on a standard LRU cache.
// All operations are meant to be thread-safe.
type LRUCache interface {
	// Add caches value for key. If cache is full (not part of the
	// interface), oldest records are dropped to make room.
	Add(key, value []byte)
	// Get retrieves value corresponding to key from cache if present.
	// Returns false in second argument if no value found or any other
	// error.
	Get(key []byte) ([]byte, bool)
	// Len returns number of values currently cached.
	Len() int
}

// memCache implements LRUCache using simple in-memory map.
type memCache struct {
	maxEntries int
	ttl        time.Duration

	mu    sync.Mutex
	cache map[string]*list.Element
	ll    *list.List
}

// NewSimpleLRUCache returns new in-memory cache of size maxEntries and with
// global ttl (records that are older than ttl are assumed to be missing
// and discarded).
func NewSimpleLRUCache(maxEntries int, ttl time.Duration) LRUCache {
	return &memCache{
		maxEntries: maxEntries,
		ttl:        ttl,
		cache:      make(map[string]*list.Element),
		ll:         list.New(),
	}
}

// Add implements LRUCache.Add.
func (c *memCache) Add(key, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If already present, refresh record in the queue.
	e, ok := c.cache[string(key)]
	if ok && e.Value.(*entry).validUntil.After(time.Now()) {
		c.ll.MoveToFront(e)
		return
	}

	c.cache[string(key)] = c.ll.PushFront(&entry{key, val, time.Now().Add(c.ttl)})

	if c.ll.Len() > c.maxEntries {
		c.removeOldest()
	}
}

// Get implements LRUCache.Get.
func (c *memCache) Get(key []byte) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.cache[string(key)]
	if ok && e.Value.(*entry).validUntil.After(time.Now()) {
		c.ll.MoveToFront(e)
		return e.Value.(*entry).value, true
	}
	return nil, false
}

// Must hold c.mu.
func (c *memCache) removeOldest() {
	e := c.ll.Back()
	if e == nil {
		return
	}

	c.ll.Remove(e)
	delete(c.cache, string(e.Value.(*entry).key))
}

// Len implements LRUCache.Len.
func (c *memCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len()
}
*/

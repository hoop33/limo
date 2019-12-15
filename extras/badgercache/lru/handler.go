package badgerlrucache

/*
import (
	"bytes"
	"container/list"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"hash/fnv"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	log "github.com/sirupsen/logrus"
)
*/

/*
	Refs:
	- https://github.com/dilyevsky/httplru/blob/master/cache/badger_lru.go
	- https://github.com/dilyevsky/httplru/blob/master/main.go
*/
/*
// badgerCache implements LRU cache using BadgerDB as high-performant,
// embedded k/v store.
type Cache struct {
	maxEntries int
	ttl        time.Duration

	db *badger.DB
	mu sync.Mutex
	// TODO(dilyevsky): Perhaps LSM or B+ tree would be of use, however
	// as it is, it can already support millions of cached entries pretty
	// efficiently.
	// Anything larger would probably need to be distributed anyway.
	ll    *list.List
	cache map[uint64]*list.Element
}

// NewLRUCache returns new cache store.
func NewLRUCache(db *badger.DB, maxEntries int, ttl time.Duration) LRUCache {
	return &Cache{
		maxEntries: maxEntries,
		ttl:        ttl,
		db:         db,
		ll:         list.New(),
		cache:      make(map[uint64]*list.Element),
	}
}

// dbEntry contains value entries for BadgerDB (original key and its
// corresponding value) to be serialized into byte array.
type dbEntry struct {
	Key, Value []byte
}

// kbuf returns byte-array representation of key.
func kbuf(key uint64) []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, key)
	return buf
}

// fingerprint returns fnv hash of b.
func fingerprint(b []byte) uint64 {
	h := fnv.New64a()
	h.Write([]byte(b))
	return h.Sum64()
}

// writeToDB writes gob-encoded val for key to BadgerDB (sets record ttl).
func (c *Cache) writeToDB(key uint64, val *dbEntry, ttl time.Duration) error {
	// TODO(dilyevsky): Protobuf should be way faster.
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(val)
	if err != nil {
		return err
	}

	return c.db.Update(func(txn *badger.Txn) error {
		return txn.SetWithTTL(kbuf(key), buf.Bytes(), ttl)
	})
}

// readFromDB reads *dbEntry corresponding to key from BadgerDB.
func (c *Cache) readFromDB(key uint64) (*dbEntry, error) {
	e := new(dbEntry)
	err := c.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(kbuf(key))
		if err != nil {
			return err
		}
		val, err := item.Value()
		if err != nil {
			return err
		}

		buf := bytes.NewBuffer(val)
		decoder := gob.NewDecoder(buf)
		err = decoder.Decode(e)
		if err != nil {
			return err
		}

		return nil
	})
	return e, err
}

// removeFromDB removes record from BadgerDB based on key.
func (c *Cache) removeFromDB(key uint64) error {
	return c.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(kbuf(key))
	})
}

// Add implements LRUCache.Add.
func (c *Cache) Add(key, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	h := fingerprint(key)
	e, ok := c.cache[h]
	if ok {
		c.ll.MoveToFront(e)
		return
	}

	// If new size exceeds cache max size, attempt to remove oldest entry
	// first (can't add new values if that fails).
	if c.ll.Len() == c.maxEntries {
		if err := c.removeOldest(); err != nil {
			return
		}
	}
	if err := c.writeToDB(h, &dbEntry{key, val}, c.ttl); err != nil {
		log.Errorf("failed write to DB: %v", err)
		return
	}
	c.cache[h] = c.ll.PushFront(h)
}

// Get implements LRUCache.Get.
func (c *Cache) Get(key []byte) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	h := fingerprint(key)
	e, ok := c.cache[h]
	if !ok {
		return nil, false
	}
	dbEntry, err := c.readFromDB(h)
	if err != nil {
		return nil, false
	}
	if string(dbEntry.Key) != string(key) { // Hash collision (rare).
		return nil, false
	}
	c.ll.MoveToFront(e)
	return dbEntry.Value, true
}

// Must hold c.mu.
func (c *Cache) removeOldest() error {
	e := c.ll.Back()
	if e == nil {
		return errors.New("no entries")
	}

	h := e.Value.(uint64)
	if err := c.removeFromDB(h); err != nil {
		return err
	}
	c.ll.Remove(e)
	delete(c.cache, h)
	return nil
}

// Len implements LRUCache.Len.
func (c *Cache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len()
}
*/

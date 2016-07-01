package forestdb

import (
	"fmt"
	"sync"
)

// KVPool is a structure representing a pool of KVStores
// inside a file.  Each has been opened with it's own
// File handle, so they can be used concurrently safely.
type KVPool struct {
	closedMutex sync.RWMutex
	closed      bool
	stores      chan *KVStore
}

var PoolClosed = fmt.Errorf("pool already closed")

func NewKVPool(filename string, config *Config, kvstore string, kvconfig *KVStoreConfig, size int) (*KVPool, error) {
	rv := KVPool{}
	rv.stores = make(chan *KVStore, size)
	for i := 0; i < size; i++ {
		db, err := Open(filename, config)
		if err != nil {
			// close everything else we've already opened
			rv.Close() // ignore errors closing? and return open error?
			return nil, err
		}
		kvs, err := db.OpenKVStore(kvstore, kvconfig)
		if err != nil {
			// close the db file we just opened
			db.Close()
			// close everything else we've already opened
			rv.Close() // ignore errors closing? and return open error?
			return nil, err
		}
		rv.stores <- kvs
	}
	return &rv, nil
}

func (p *KVPool) Get() (*KVStore, error) {
	rv, ok := <-p.stores
	if !ok {
		return nil, PoolClosed
	}
	return rv, nil
}

func (p *KVPool) Return(kvs *KVStore) error {
	p.closedMutex.RLock()
	defer p.closedMutex.RUnlock()
	if !p.closed {
		p.stores <- kvs
		return nil
	}
	return PoolClosed
}

func (p *KVPool) Close() (rverr error) {
	p.closedMutex.Lock()
	if !p.closed {
		close(p.stores)
	}
	p.closed = true
	p.closedMutex.Unlock()

	for kvs := range p.stores {
		err := kvs.Close()
		if err != nil {
			if rverr == nil {
				rverr = err
			}
			// keep going try to close file
		}
		db := kvs.File()
		err = db.Close()
		if err != nil {
			if rverr == nil {
				rverr = err
			}
		}
	}
	return
}

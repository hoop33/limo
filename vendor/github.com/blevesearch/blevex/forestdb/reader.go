//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package forestdb

import (
	"sync"

	"github.com/blevesearch/bleve/index/store"
	"github.com/couchbase/goforestdb"
)

type Reader struct {
	store    *Store
	kvstore  *forestdb.KVStore
	snapshot *forestdb.KVStore
	parent   *Reader
	m        sync.Mutex // Protects the fields that follow.
	refs     int
}

func (r *Reader) addRef() *Reader {
	r.m.Lock()
	r.refs++
	r.m.Unlock()
	return r
}

func (r *Reader) decRef() (rverr error) {
	r.m.Lock()
	r.refs--
	refs := r.refs
	r.m.Unlock()

	if refs > 0 {
		return nil
	}

	rverr = r.snapshot.Close()

	if r.parent != nil {
		err := r.parent.decRef()
		if rverr != nil {
			return rverr // return first error
		}
		return err
	}

	// only the "root", nil-parent Reader will return the kvstore to the kvpool.
	// return to pool even error closing snapshot?
	if r.kvstore != nil {
		err := r.store.returnKVStore(r.kvstore)
		if rverr != nil {
			return rverr // return first error
		}
		return err
	}

	return rverr
}

func (r *Reader) Get(key []byte) ([]byte, error) {
	res, err := r.snapshot.GetKV(key)
	if err != nil && err != forestdb.RESULT_KEY_NOT_FOUND {
		return nil, err
	}
	return res, nil
}

func (r *Reader) MultiGet(keys [][]byte) ([][]byte, error) {
	return store.MultiGet(r, keys)
}

func (r *Reader) PrefixIterator(prefix []byte) store.KVIterator {
	// compute range end
	var end []byte
	for i := len(prefix) - 1; i >= 0; i-- {
		c := prefix[i]
		if c < 0xff {
			end = make([]byte, i+1)
			copy(end, prefix)
			end[i] = c + 1
			break
		}
	}
	itr, err := r.snapshot.IteratorInit(prefix, end, forestdb.ITR_NO_DELETES|forestdb.FDB_ITR_SKIP_MAX_KEY)
	rv := Iterator{
		store:    r.store,
		iterator: itr,
		valid:    err == nil,
		start:    prefix,
		parent:   r.addRef(),
	}
	rv.Seek(prefix)
	return &rv
}

func (r *Reader) RangeIterator(start, end []byte) store.KVIterator {
	opts := forestdb.ITR_NO_DELETES
	if end != nil {
		opts = opts | forestdb.FDB_ITR_SKIP_MAX_KEY
	}
	itr, err := r.snapshot.IteratorInit(start, end, opts)
	rv := Iterator{
		store:    r.store,
		iterator: itr,
		valid:    err == nil,
		start:    start,
		parent:   r.addRef(),
	}
	rv.Seek(start)
	return &rv
}

func (r *Reader) Close() error {
	return r.decRef()
}

// Reader method allows cloning of a snapshot for multi-threaded use
func (r *Reader) Reader() (store.KVReader, error) {
	snapshot, err := r.snapshot.SnapshotOpen(forestdb.SnapshotInmem)
	if err != nil {
		return nil, err
	}
	return &Reader{
		store:    r.store,
		kvstore:  nil, // dont return to pool, since this is a clone
		snapshot: snapshot,
		parent:   r.addRef(),
		refs:     1,
	}, nil
}

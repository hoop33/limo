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
	"encoding/json"
	"fmt"
	"sync"

	"github.com/blevesearch/bleve/index/store"
	"github.com/blevesearch/bleve/registry"
	"github.com/couchbase/goforestdb"
)

const Name = "forestdb"

// DefaultConcurrent is the default kvpool size.  When 0, don't use
// kvpool (0 is the default).
var DefaultConcurrent = 0

type Store struct {
	m      sync.RWMutex
	path   string
	kvpool *forestdb.KVPool // May be nil if kvpool disabled.
	mo     store.MergeOperator

	fdbConfig *forestdb.Config
	kvsConfig *forestdb.KVStoreConfig

	statsMutex  sync.Mutex
	statsHandle *forestdb.KVStore
	stats       *kvStat

	skipBatch bool
}

func New(mo store.MergeOperator, config map[string]interface{}) (store.KVStore, error) {

	path, ok := config["path"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify path")
	}

	fdbDefaultConfig := forestdb.DefaultConfig()
	fdbDefaultConfig.SetCompactionMode(forestdb.COMPACT_AUTO)
	fdbDefaultConfig.SetMultiKVInstances(false)
	fdbConfig, err := applyConfig(fdbDefaultConfig, config)
	if err != nil {
		return nil, err
	}

	kvsConfig := forestdb.DefaultKVStoreConfig()
	if cim, ok := config["create_if_missing"].(bool); ok && cim {
		kvsConfig.SetCreateIfMissing(true)
	}

	var kvpool *forestdb.KVPool
	var statsHandle *forestdb.KVStore

	numConcurrent := DefaultConcurrent
	if nc, ok := config["num_concurrent"].(float64); ok {
		numConcurrent = int(nc)
	}
	if numConcurrent > 0 {
		// request 1 extra connection in pool to be reserved for issuing
		// stats calls
		kvpool, err := forestdb.NewKVPool(path, fdbConfig, "default", kvsConfig,
			numConcurrent+1)
		if err != nil {
			return nil, err
		}

		statsHandle, err = kvpool.Get()
		if err != nil {
			kvpool.Close()
			return nil, err
		}
	} else {
		statsHandle, err = forestdb.OpenFileKVStore(path, fdbConfig, "default", kvsConfig)
		if err != nil {
			return nil, err
		}
	}

	var skipBatch bool
	if skipBatchV, ok := config["skip_batch"].(bool); ok {
		skipBatch = skipBatchV
	}

	rv := Store{
		path:        path,
		kvpool:      kvpool,
		mo:          mo,
		fdbConfig:   fdbConfig,
		kvsConfig:   kvsConfig,
		statsHandle: statsHandle,
		skipBatch:   skipBatch,
	}

	rv.stats = &kvStat{s: &rv}

	return &rv, nil
}

func (s *Store) Close() error {
	if s.statsHandle != nil {
		s.returnKVStore(s.statsHandle)
		s.statsHandle = nil
	}
	if s.kvpool != nil {
		return s.kvpool.Close()
	}
	return nil
}

func (s *Store) Reader() (store.KVReader, error) {
	kvstore, err := s.acquireKVStore()
	if err != nil {
		return nil, err
	}
	snapshot, err := kvstore.SnapshotOpen(forestdb.SnapshotInmem)
	if err != nil {
		return nil, err
	}
	return &Reader{
		store:    s,
		kvstore:  kvstore,
		snapshot: snapshot,
		parent:   nil,
		refs:     1,
	}, nil
}

func (s *Store) Stats() json.Marshaler {
	return s.stats
}

func (s *Store) StatsMap() map[string]interface{} {
	return s.stats.statsMap()
}

func (s *Store) Writer() (store.KVWriter, error) {
	kvstore, err := s.acquireKVStore()
	if err != nil {
		return nil, err
	}
	return &Writer{
		store:   s,
		kvstore: kvstore,
	}, nil
}

func (s *Store) acquireKVStore() (*forestdb.KVStore, error) {
	if s.kvpool != nil {
		return s.kvpool.Get()
	}
	return forestdb.OpenFileKVStore(s.path, s.fdbConfig, "default", s.kvsConfig)
}

func (s *Store) returnKVStore(kvs *forestdb.KVStore) error {
	if s.kvpool != nil {
		return s.kvpool.Return(kvs)
	}
	return forestdb.CloseFileKVStore(kvs)
}

func init() {
	registry.RegisterKVStore(Name, New)
}

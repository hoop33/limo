//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package leveldb

import (
	"fmt"
	"sync"

	"github.com/blevesearch/bleve/index/store"
	"github.com/blevesearch/bleve/registry"
	"github.com/jmhodges/levigo"
)

const Name = "leveldb"

type Store struct {
	path string
	opts *levigo.Options
	db   *levigo.DB
	mo   store.MergeOperator

	mergeMutex sync.Mutex
}

func New(mo store.MergeOperator, config map[string]interface{}) (store.KVStore, error) {
	path, ok := config["path"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify path")
	}

	rv := Store{
		path: path,
		opts: levigo.NewOptions(),
		mo:   mo,
	}

	_, err := applyConfig(rv.opts, config)
	if err != nil {
		return nil, err
	}

	rv.db, err = levigo.Open(rv.path, rv.opts)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

func (s *Store) Close() error {
	s.db.Close()
	s.opts.Close()
	return nil
}

func (s *Store) Reader() (store.KVReader, error) {
	snapshot := s.db.NewSnapshot()
	options := defaultReadOptions()
	options.SetSnapshot(snapshot)
	return &Reader{
		store:    s,
		snapshot: s.db.NewSnapshot(),
		options:  options,
	}, nil
}

func (s *Store) Writer() (store.KVWriter, error) {
	return &Writer{
		store:   s,
		options: defaultWriteOptions(),
	}, nil
}

func init() {
	registry.RegisterKVStore(Name, New)
}

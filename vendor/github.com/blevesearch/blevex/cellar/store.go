//  Copyright (c) 2016 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package cellar

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/blevesearch/bleve/index/store"
	"github.com/blevesearch/bleve/registry"
	"github.com/couchbaselabs/cellar"
)

const Name = "cellar"

type Store struct {
	path string
	db   *cellar.Cellar
	mo   store.MergeOperator
}

func New(mo store.MergeOperator, config map[string]interface{}) (store.KVStore, error) {
	path, ok := config["path"].(string)
	if !ok {
		return nil, fmt.Errorf("must specify path")
	}

	cellarOpts := &cellar.Options{}
	if autoMerge, ok := config["automerge"].(bool); ok {
		log.Printf("setting automatic merge to: %t", autoMerge)
		cellarOpts.AutomaticMerge = autoMerge
	}

	db, err := cellar.Open(path, cellarOpts)
	if err != nil {
		return nil, err
	}

	rv := Store{
		path: path,
		db:   db,
		mo:   mo,
	}
	return &rv, nil
}

func (bs *Store) Close() error {
	return bs.db.Close()
}

func (bs *Store) Reader() (store.KVReader, error) {
	tx, err := bs.db.Begin(false)
	if err != nil {
		return nil, err
	}
	return &Reader{
		store: bs,
		tx:    tx,
	}, nil
}

func (bs *Store) Writer() (store.KVWriter, error) {
	return &Writer{
		store: bs,
	}, nil
}

func (bs *Store) Stats() json.Marshaler {
	return &stats{
		s: bs,
	}
}

func init() {
	registry.RegisterKVStore(Name, New)
}

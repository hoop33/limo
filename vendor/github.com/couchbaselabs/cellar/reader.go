//  Copyright (c) 2016 Couchbase, Inc.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the
//  License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on an "AS
//  IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
//  express or implied. See the License for the specific language
//  governing permissions and limitations under the License.

package cellar

import (
	"fmt"

	"github.com/boltdb/bolt"
)

type reader struct {
	root      segmentList
	txs       []*bolt.Tx
	mutations []*bolt.Bucket
	deletions []*bolt.Bucket
}

func newReader(root segmentList) (*reader, error) {
	rv := &reader{
		root:      root,
		txs:       make([]*bolt.Tx, 0, len(root)),
		mutations: make([]*bolt.Bucket, 0, len(root)),
		deletions: make([]*bolt.Bucket, 0, len(root)),
	}

	for _, segment := range root {
		tx, err := segment.DB.Begin(false)
		if err != nil {
			return nil, fmt.Errorf("newReader begin '%d': %v", segment.seq, err)
		}
		rv.txs = append(rv.txs, tx)
		mutationsBucket := tx.Bucket(mutationsBucketName)
		rv.mutations = append(rv.mutations, mutationsBucket)
		deletionsBucket := tx.Bucket(deletionsBucketName)
		rv.deletions = append(rv.deletions, deletionsBucket)
	}

	return rv, nil
}

func (r *reader) Get(key []byte) []byte {
	var rv []byte
	for j, mutationsBucket := range r.mutations {
		deletionsBucket := r.deletions[j]
		v := deletionsBucket.Get(key)
		if v != nil {
			// key deleted don't look any further
			return nil
		}
		v = mutationsBucket.Get(key)
		if v != nil {
			rv = make([]byte, len(v))
			copy(rv, v)
			break
		}
	}
	return rv
}

func (r *reader) Close() error {
	var err error
	for _, tx := range r.txs {
		cerr := tx.Rollback()
		// remember first err seen, but keep closing txs
		if cerr != nil && err == nil {
			err = cerr
		}
	}
	return err
}

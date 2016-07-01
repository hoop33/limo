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
	"fmt"

	"github.com/blevesearch/bleve/index/store"
	"github.com/couchbase/goforestdb"
)

type Writer struct {
	store   *Store
	kvstore *forestdb.KVStore
}

func (w *Writer) NewBatch() store.KVBatch {
	rv := Batch{
		w:     w,
		merge: store.NewEmulatedMerge(w.store.mo),
		batch: forestdb.NewKVBatch(),
	}
	return &rv
}

func (w *Writer) NewBatchEx(options store.KVBatchOptions) ([]byte, store.KVBatch, error) {
	// NOTE: We've reverted to old, emulated batch due to MB-17558.
	//
	// rv := newBatchEx(w, options)
	// return rv.buf, rv, nil

	return make([]byte, options.TotalBytes), w.NewBatch(), nil
}

func (w *Writer) ExecuteBatch(b store.KVBatch) error {
	w.store.m.Lock()
	defer w.store.m.Unlock()

	batchex, ok := b.(*BatchEx)
	if ok {
		err := batchex.apply()
		if err != nil {
			return err
		}

		return w.kvstore.File().Commit(forestdb.COMMIT_NORMAL)
	}

	batch, ok := b.(*Batch)
	if !ok {
		return fmt.Errorf("wrong type of batch")
	}

	for key, mergeOps := range batch.merge.Merges {
		k := []byte(key)
		ob, err := w.kvstore.GetKV(k)
		if err != nil && err != forestdb.RESULT_KEY_NOT_FOUND {
			return err
		}
		mergedVal, fullMergeOk := w.store.mo.FullMerge(k, ob, mergeOps)
		if !fullMergeOk {
			return fmt.Errorf("merge operator returned failure")
		}
		batch.Set(k, mergedVal)
	}

	if w.store.skipBatch {
		if batch.skipBatchErr != nil {
			return batch.skipBatchErr
		}
		return w.kvstore.File().Commit(forestdb.COMMIT_NORMAL)
	}

	return w.kvstore.ExecuteBatch(batch.batch, forestdb.COMMIT_NORMAL)
}

func (w *Writer) Close() error {
	return w.store.returnKVStore(w.kvstore)
}

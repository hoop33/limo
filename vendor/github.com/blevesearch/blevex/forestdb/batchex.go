//  Copyright (c) 2016 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package forestdb

/*
#include <stdio.h>
#include <stdlib.h>
#include <libforestdb/forestdb.h>

fdb_status blevex_forestdb_execute_direct_batch(
    fdb_kvs_handle* handle,
    int num_sets,
    char** set_keys,
    size_t* set_keys_sizes,
    char** set_vals,
    size_t* set_vals_sizes,
    int num_deletes,
    char** delete_keys,
    size_t* delete_keys_sizes) {
    fdb_status rv;

    fdb_doc *doc;
    fdb_doc_create(&doc, NULL, 0, NULL, 0, NULL, 0);

    int i = 0;
    for (i = 0; i < num_sets; i++ ) {
        doc->key = set_keys[i];
        doc->keylen = set_keys_sizes[i];

        doc->body = set_vals[i];
        doc->bodylen = set_vals_sizes[i];

        rv = fdb_set(handle, doc);
        if (rv != FDB_RESULT_SUCCESS) {
            doc->key = NULL;
            doc->keylen = 0;

            doc->body = NULL;
            doc->bodylen = 0;

            fdb_doc_free(doc);

            return rv;
        }
    }

    doc->body = NULL;
    doc->bodylen = 0;

    for (i = 0; i < num_deletes; i++ ) {
        doc->key = delete_keys[i];
        doc->keylen = delete_keys_sizes[i];

        rv = fdb_del(handle, doc);
        if (rv != FDB_RESULT_SUCCESS) {
            fdb_doc_free(doc);
            return rv;
        }
    }

    doc->key = NULL;
    doc->keylen = 0;

    fdb_doc_free(doc);

    return FDB_RESULT_SUCCESS;
}

void blevex_forestdb_alloc_direct_batch(size_t totalBytes, size_t n, void **out) {
    out[0] = malloc(totalBytes);
    out[1] = malloc(n * sizeof(char *));
    out[2] = malloc(n * sizeof(size_t));
}

void blevex_forestdb_free_direct_batch(void **bufs) {
    free(bufs[0]);
    free(bufs[1]);
    free(bufs[2]);
}
*/
//#include <libforestdb/forestdb.h>
//#cgo LDFLAGS: -lforestdb
import "C"

import (
	"fmt"
	"reflect"
	"unsafe"

	goforestdb "github.com/couchbase/goforestdb"

	"github.com/blevesearch/bleve/index/store"
)

type BatchEx struct {
	w     *Writer
	merge *store.EmulatedMerge

	cbufs []unsafe.Pointer
	buf   []byte

	num_sets       int
	set_keys       []*C.char
	set_keys_sizes []C.size_t
	set_vals       []*C.char
	set_vals_sizes []C.size_t

	num_deletes       int
	delete_keys       []*C.char
	delete_keys_sizes []C.size_t
}

func newBatchEx(w *Writer, o store.KVBatchOptions) *BatchEx {
	s := o.NumSets + o.NumMerges
	ss := s + o.NumSets + o.NumMerges
	ssd := ss + o.NumDeletes

	cbufs := make([]unsafe.Pointer, 3)

	C.blevex_forestdb_alloc_direct_batch(C.size_t(o.TotalBytes),
		C.size_t(ssd), (*unsafe.Pointer)(&cbufs[0]))

	buf := unsafeToByteSlice(cbufs[0], o.TotalBytes)
	arr_ptr_char := unsafeToCPtrCharSlice(cbufs[1], ssd)
	arr_size_t := unsafeToCSizeTSlice(cbufs[2], ssd)

	return &BatchEx{
		w:     w,
		merge: store.NewEmulatedMerge(w.store.mo),

		cbufs:             cbufs,
		buf:               buf,
		set_keys:          arr_ptr_char[0:s],
		set_keys_sizes:    arr_size_t[0:s],
		set_vals:          arr_ptr_char[s:ss],
		set_vals_sizes:    arr_size_t[s:ss],
		delete_keys:       arr_ptr_char[ss:ssd],
		delete_keys_sizes: arr_size_t[ss:ssd],
	}
}

func (b *BatchEx) Set(key, val []byte) {
	b.set_keys[b.num_sets] = (*C.char)(unsafe.Pointer(&key[0]))
	b.set_keys_sizes[b.num_sets] = (C.size_t)(len(key))
	b.set_vals[b.num_sets] = (*C.char)(unsafe.Pointer(&val[0]))
	b.set_vals_sizes[b.num_sets] = (C.size_t)(len(val))
	b.num_sets += 1
}

func (b *BatchEx) Delete(key []byte) {
	b.delete_keys[b.num_deletes] = (*C.char)(unsafe.Pointer(&key[0]))
	b.delete_keys_sizes[b.num_deletes] = (C.size_t)(len(key))
	b.num_deletes += 1
}

func (b *BatchEx) Merge(key, val []byte) {
	b.merge.Merge(key, val)
}

func (b *BatchEx) Reset() {
	b.num_sets = 0
	b.num_deletes = 0

	b.merge = store.NewEmulatedMerge(b.w.store.mo)
}

func (b *BatchEx) Close() error {
	b.w = nil
	b.merge = nil

	C.blevex_forestdb_free_direct_batch((*unsafe.Pointer)(&b.cbufs[0]))

	b.cbufs = nil
	b.buf = nil

	b.num_sets = 0
	b.set_keys = nil
	b.set_keys_sizes = nil
	b.set_vals = nil
	b.set_vals_sizes = nil

	b.num_deletes = 0
	b.delete_keys = nil
	b.delete_keys_sizes = nil

	return nil
}

func (b *BatchEx) apply() error {
	// Hold onto final merge key/val bytes so GC doesn't collect them
	// until we're done.
	var mergeBytes [][]byte

	if len(b.merge.Merges) > 0 {
		mergeBytes = make([][]byte, 0, len(b.merge.Merges)*2)

		mergeKeys := make([][]byte, len(b.merge.Merges))
		mergeOps := make([][][]byte, len(b.merge.Merges))

		i := 0
		for key, ops := range b.merge.Merges {
			mergeKeys[i] = []byte(key)
			mergeOps[i] = ops
			i += 1
		}

		currVals, releaseVals, err := directMultiGet(b.w.kvstore, mergeKeys)
		if err != nil {
			return err
		}

		defer releaseVals()

		for i, key := range mergeKeys {
			mergedVal, ok := b.w.store.mo.FullMerge(key, currVals[i], mergeOps[i])
			if !ok {
				return fmt.Errorf("forestdb BatchEx merge operator failure,"+
					" key: %s, currVal: %q, mergeOps: %#v, i: %d",
					key, currVals[i], mergeOps[i], i)
			}

			mergeBytes = append(mergeBytes, key, mergedVal)

			b.Set(key, mergedVal)
		}
	}

	var num_sets C.int
	var set_keys **C.char
	var set_keys_sizes *C.size_t
	var set_vals **C.char
	var set_vals_sizes *C.size_t

	var num_deletes C.int
	var delete_keys **C.char
	var delete_keys_sizes *C.size_t

	if b.num_sets > 0 {
		num_sets = (C.int)(b.num_sets)
		set_keys = (**C.char)(unsafe.Pointer(&b.set_keys[0]))
		set_keys_sizes = (*C.size_t)(unsafe.Pointer(&b.set_keys_sizes[0]))
		set_vals = (**C.char)(unsafe.Pointer(&b.set_vals[0]))
		set_vals_sizes = (*C.size_t)(unsafe.Pointer(&b.set_vals_sizes[0]))
	}

	if b.num_deletes > 0 {
		num_deletes = (C.int)(b.num_deletes)
		delete_keys = (**C.char)(unsafe.Pointer(&b.delete_keys[0]))
		delete_keys_sizes = (*C.size_t)(unsafe.Pointer(&b.delete_keys_sizes[0]))
	}

	errNo := C.blevex_forestdb_execute_direct_batch(
		(*C.fdb_kvs_handle)(b.w.kvstore.Handle()),
		num_sets,
		set_keys,
		set_keys_sizes,
		set_vals,
		set_vals_sizes,
		num_deletes,
		delete_keys,
		delete_keys_sizes)
	if int(errNo) != 0 {
		return goforestdb.Error(errNo)
	}

	if mergeBytes != nil { // Ok to let GC have mergeBytes now.
		mergeBytes = nil
	}

	return nil
}

// Originally from github.com/tecbot/gorocksdb/util.go.
func unsafeToByteSlice(data unsafe.Pointer, len int) []byte {
	var value []byte

	sH := (*reflect.SliceHeader)(unsafe.Pointer(&value))
	sH.Cap, sH.Len, sH.Data = len, len, uintptr(data)

	return value
}

func unsafeToCPtrCharSlice(data unsafe.Pointer, len int) []*C.char {
	var value []*C.char

	sH := (*reflect.SliceHeader)(unsafe.Pointer(&value))
	sH.Cap, sH.Len, sH.Data = len, len, uintptr(data)

	return value
}

func unsafeToCSizeTSlice(data unsafe.Pointer, len int) []C.size_t {
	var value []C.size_t

	sH := (*reflect.SliceHeader)(unsafe.Pointer(&value))
	sH.Cap, sH.Len, sH.Data = len, len, uintptr(data)

	return value
}

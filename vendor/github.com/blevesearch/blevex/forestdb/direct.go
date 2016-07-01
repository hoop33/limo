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

fdb_status blevex_forestdb_multi_get(
    fdb_kvs_handle* handle,
    size_t num_keys,
    const char* const* keys,
    const size_t* keys_sizes,
    char** vals,
    size_t* vals_sizes) {
    fdb_status rv;

    fdb_doc *doc;
    fdb_doc_create(&doc, NULL, 0, NULL, 0, NULL, 0);

    int i;
    for (i = 0; i < num_keys; i++) {
        doc->key = (void *) keys[i];
        doc->keylen = keys_sizes[i];

        doc->body = NULL;
        doc->bodylen = 0;

        rv = fdb_get(handle, doc);
        if (rv != FDB_RESULT_SUCCESS && rv != FDB_RESULT_KEY_NOT_FOUND) {
            doc->key = NULL;
            doc->keylen = 0;

            doc->body = NULL;
            doc->bodylen = 0;

            fdb_doc_free(doc);

            return rv;
        }

        vals[i] = doc->body;
        vals_sizes[i] = doc->bodylen;
    }

    doc->key = NULL;
    doc->keylen = 0;

    doc->body = NULL;
    doc->bodylen = 0;

    fdb_doc_free(doc);

    return FDB_RESULT_SUCCESS;
}

void blevex_forestdb_free_bufs(size_t num_bufs, char** bufs) {
    int i;
    for (i = 0; i < num_bufs; i++) {
        if (bufs[i] != NULL) {
            free(bufs[i]);
            bufs[i] = NULL;
        }
    }
}
*/
import "C"

import (
	"reflect"
	"unsafe"

	goforestdb "github.com/couchbase/goforestdb"
)

func directMultiGet(kvstore *goforestdb.KVStore, keysIn [][]byte) (
	valsOut [][]byte,
	valsRelease func(),
	err error) {
	keys := make([]*C.char, len(keysIn))
	key_lens := make([]C.size_t, len(keysIn))
	vals := make([]*C.char, len(keysIn))
	val_lens := make([]C.size_t, len(keysIn))

	for i, key := range keysIn {
		keys[i] = (*C.char)(unsafe.Pointer(&key[0]))
		key_lens[i] = (C.size_t)(len(key))
	}

	valsRelease = func() {
		C.blevex_forestdb_free_bufs(
			(C.size_t)(len(vals)), (**C.char)(unsafe.Pointer(&vals[0])))

		for i := range vals {
			vals[i] = nil
		}
	}

	errNo := C.blevex_forestdb_multi_get(
		(*C.fdb_kvs_handle)(kvstore.Handle()),
		(C.size_t)(len(keys)),
		(**C.char)(unsafe.Pointer(&keys[0])),
		(*C.size_t)(unsafe.Pointer(&key_lens[0])),
		(**C.char)(unsafe.Pointer(&vals[0])),
		(*C.size_t)(unsafe.Pointer(&val_lens[0])))
	if int(errNo) != 0 {
		valsRelease()

		return nil, nil, goforestdb.Error(errNo)
	}

	valsOut = make([][]byte, len(vals))
	for i, val := range vals {
		if val != nil {
			valsOut[i] = charToByte(val, val_lens[i])
		}
	}

	return valsOut, valsRelease, nil
}

// From github.com/tecbot/gorocksdb.
func charToByte(data *C.char, len C.size_t) []byte {
	var value []byte

	sH := (*reflect.SliceHeader)(unsafe.Pointer(&value))
	sH.Cap, sH.Len, sH.Data = int(len), int(len), uintptr(unsafe.Pointer(data))

	return value
}

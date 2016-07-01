package forestdb

//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

//#include <stdlib.h>
//#include <libforestdb/forestdb.h>
import "C"

import (
	"unsafe"
)

// GetKV simplified API for key/value access to Get()
func (k *KVStore) GetKV(key []byte) ([]byte, error) {

	var kk unsafe.Pointer
	if len(key) != 0 {
		kk = unsafe.Pointer(&key[0])
	}
	lenk := len(key)

	var bodyLen C.size_t
	var bodyPointer unsafe.Pointer

	Log.Tracef("fdb_get_kv call k:%p db:%p kk:%v", k, k.db, kk)
	errNo := C.fdb_get_kv(k.db, kk, C.size_t(lenk), &bodyPointer, &bodyLen)
	Log.Tracef("fdb_get_kv retn k:%p errNo:%v body:%p len:%v", k, errNo, bodyPointer, bodyLen)
	if errNo != RESULT_SUCCESS {
		return nil, Error(errNo)
	}

	body := C.GoBytes(bodyPointer, C.int(bodyLen))
	C.fdb_free_block(bodyPointer)
	return body, nil
}

// SetKV simplified API for key/value access to Set()
func (k *KVStore) SetKV(key, value []byte) error {

	var kk, v unsafe.Pointer

	if len(key) != 0 {
		kk = unsafe.Pointer(&key[0])
	}

	if len(value) != 0 {
		v = unsafe.Pointer(&value[0])
	}

	lenk := len(key)
	lenv := len(value)

	Log.Tracef("fdb_set_kv call k:%p db:%p kk:%v v:%v", k, k.db, kk, v)
	errNo := C.fdb_set_kv(k.db, kk, C.size_t(lenk), v, C.size_t(lenv))
	Log.Tracef("fdb_set_kv retn k:%p errNo:%v", k, errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// DeleteKV simplified API for key/value access to Delete()
func (k *KVStore) DeleteKV(key []byte) error {

	var kk unsafe.Pointer
	if len(key) != 0 {
		kk = unsafe.Pointer(&key[0])
	}

	lenk := len(key)

	Log.Tracef("fdb_del_kv call k:%p db:%p kk:%v", k, k.db, kk)
	errNo := C.fdb_del_kv(k.db, kk, C.size_t(lenk))
	Log.Tracef("fdb_del_kv retn k:%p errNo:%v", k, errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

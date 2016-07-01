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
//#include "log.h"
//extern void LogCallbackInternal(int, char*, char*);
//void log_callback(int errcode, char *msg, void *ctx) {
//    LogCallbackInternal(errcode, msg, ctx);
//}
//extern void FatalErrorCallbackInternal();
//void gofatal_error_callback() {
//    FatalErrorCallbackInternal();
//}
import "C"

import "unsafe"

// KVStore handle
type KVStore struct {
	f    *File
	db   *C.fdb_kvs_handle
	name string
}

// File returns the File containing this KVStore
func (k *KVStore) File() *File {
	return k.f
}

// Handle returns the underlying fdb_kvs_handle for advanced uses.
func (k *KVStore) Handle() *C.fdb_kvs_handle {
	return k.db
}

// Close the KVStore and release related resources.
func (k *KVStore) Close() error {
	Log.Tracef("fdb_kvs_close call k:%p db:%p", k, k.db)
	errNo := C.fdb_kvs_close(k.db)
	Log.Tracef("fdb_kvs_close retn k:%p errNo:%v", k, errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// Info returns the information about a given kvstore
func (k *KVStore) Info() (*KVStoreInfo, error) {
	rv := KVStoreInfo{}
	Log.Tracef("fdb_get_kvs_info call k:%p db:%p", k, k.db)
	errNo := C.fdb_get_kvs_info(k.db, &rv.info)
	Log.Tracef("fdb_kvs_close retn k:%p errNo:%v info:%v", k, errNo, rv.info)
	if errNo != RESULT_SUCCESS {
		return nil, Error(errNo)
	}
	return &rv, nil
}

// OpsInfo returns the information about the ops on given kvstore
func (k *KVStore) OpsInfo() (*KVSOpsInfo, error) {
	rv := KVSOpsInfo{}
	Log.Tracef("fdb_get_kvs_ops_info call k:%p db:%p", k, k.db)
	errNo := C.fdb_get_kvs_ops_info(k.db, &rv.info)
	Log.Tracef("fdb_get_kvs_ops_info k:%p errNo:%v info:%v", k, errNo, rv.info)
	if errNo != RESULT_SUCCESS {
		return nil, Error(errNo)
	}
	return &rv, nil
}

// Get retrieves the metadata and doc body for a given key
func (k *KVStore) Get(doc *Doc) error {
	Log.Tracef("fdb_get call k:%p db:%p doc:%v", k, k.db, doc.doc)
	errNo := C.fdb_get(k.db, doc.doc)
	Log.Tracef("fdb_get retn k:%p errNo:%v doc:%v", k, errNo, doc.doc)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// GetMetaOnly retrieves the metadata for a given key
func (k *KVStore) GetMetaOnly(doc *Doc) error {
	Log.Tracef("fdb_get_metaonly call k:%p db:%p doc:%v", k, k.db, doc.doc)
	errNo := C.fdb_get_metaonly(k.db, doc.doc)
	Log.Tracef("fdb_get_metaonly retn k:%p errNo:%v doc:%v", k, errNo, doc.doc)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// GetBySeq retrieves the metadata and doc body for a given sequence number
func (k *KVStore) GetBySeq(doc *Doc) error {
	Log.Tracef("fdb_get_byseq call k:%p db:%p doc:%v", k, k.db, doc.doc)
	errNo := C.fdb_get_byseq(k.db, doc.doc)
	Log.Tracef("fdb_get_byseq retn k:%p errNo:%v doc:%v", k, errNo, doc.doc)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// GetMetaOnlyBySeq retrieves the metadata for a given sequence number
func (k *KVStore) GetMetaOnlyBySeq(doc *Doc) error {
	Log.Tracef("fdb_get_metaonly_byseq call k:%p db:%p doc:%v", k, k.db, doc.doc)
	errNo := C.fdb_get_metaonly_byseq(k.db, doc.doc)
	Log.Tracef("fdb_get_metaonly_byseq retn k:%p errNo:%v doc:%v", k, errNo, doc.doc)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// GetByOffset retrieves a doc's metadata and body with a given doc offset in the database file
func (k *KVStore) GetByOffset(doc *Doc) error {
	Log.Tracef("fdb_get_byoffset call k:%p db:%p doc:%v", k, k.db, doc.doc)
	errNo := C.fdb_get_byoffset(k.db, doc.doc)
	Log.Tracef("fdb_get_byoffset retn k:%p errNo:%v doc:%v", k, errNo, doc.doc)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// Set update the metadata and doc body for a given key
func (k *KVStore) Set(doc *Doc) error {
	Log.Tracef("fdb_set call k:%p db:%p doc:%v", k, k.db, doc.doc)
	errNo := C.fdb_set(k.db, doc.doc)
	Log.Tracef("fdb_set retn k:%p errNo:%v doc:%v", k, errNo, doc.doc)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// Delete deletes a key, its metadata and value
func (k *KVStore) Delete(doc *Doc) error {
	Log.Tracef("fdb_del call k:%p db:%p doc:%v", k, k.db, doc.doc)
	errNo := C.fdb_del(k.db, doc.doc)
	Log.Tracef("fdb_set retn k:%p errNo:%v doc:%v", k, errNo, doc.doc)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// Shutdown destroys all the resources (e.g., buffer cache, in-memory WAL indexes, daemon compaction thread, etc.) and then shutdown the ForestDB engine
func Shutdown() error {
	Log.Tracef("fdb_shutdown call")
	errNo := C.fdb_shutdown()
	Log.Tracef("fdb_shutdown retn errNo:%v", errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

type logContext struct {
	callback *LogCallback
	name     string
	userCtx  interface{}
}

// Hold references to log callbacks and contexts.
var logCallbacks []LogCallback
var logContexts []interface{}

func registerLogCallback(cb LogCallback, ctx interface{}) int {
	logCallbacks = append(logCallbacks, cb)
	logContexts = append(logContexts, ctx)
	return len(logCallbacks) - 1
}

func (k *KVStore) SetLogCallback(l LogCallback, userCtx interface{}) {
	var ctx C.log_context
	ctx.offset = C.int(registerLogCallback(l, userCtx))
	ctx.name = C.CString(k.name)
	C.fdb_set_log_callback(k.db, C.fdb_log_callback(C.log_callback), unsafe.Pointer(&ctx))
}

func SetFatalErrorCallback(callback FatalErrorCallback) {
	fatalErrorCallback = callback
	C.fdb_set_fatal_error_callback(C.fdb_fatal_error_callback(C.gofatal_error_callback))
}

type FatalErrorCallback func()

var fatalErrorCallback FatalErrorCallback

type LogCallback func(name string, errCode int, msg string, ctx interface{})

func LoggingLogCallback(name string, errCode int, msg string, ctx interface{}) {
	Log.Errorf("ForestDB (%s) Error Code: %d Message: %s", name, errCode, msg)
}

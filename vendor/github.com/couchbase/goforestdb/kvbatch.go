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
	"reflect"
	"unsafe"
)

type batchOp struct {
	del  bool
	k    unsafe.Pointer
	klen C.size_t
	v    unsafe.Pointer
	vlen C.size_t
}

type KVBatch struct {
	ops []*batchOp
}

func NewKVBatch() *KVBatch {
	return &KVBatch{
		ops: make([]*batchOp, 0, 100),
	}
}

func copySliceToC(in []byte) (C.size_t, unsafe.Pointer) {
	outl := C.size_t(len(in))
	outv := C.malloc(outl)
	var tmpSlice []byte
	shdr := (*reflect.SliceHeader)(unsafe.Pointer(&tmpSlice))
	shdr.Data = uintptr(outv)
	shdr.Len = int(outl)
	shdr.Cap = shdr.Len
	copy(tmpSlice, in)
	return outl, outv
}

func (b *KVBatch) Set(k, v []byte) {
	bo := batchOp{del: false}
	bo.klen, bo.k = copySliceToC(k)
	if len(v) > 0 {
		bo.vlen, bo.v = copySliceToC(v)
	}
	b.ops = append(b.ops, &bo)
}

func (b *KVBatch) Delete(k []byte) {
	bo := batchOp{del: true}
	bo.klen, bo.k = copySliceToC(k)
	b.ops = append(b.ops, &bo)
}

func (b *KVBatch) Reset() {
	for _, op := range b.ops {
		if op.klen > 0 {
			C.free(op.k)
		}
		if op.vlen > 0 {
			C.free(op.v)
		}
	}
	b.ops = b.ops[:0]
}

func (k *KVStore) ExecuteBatch(b *KVBatch, opt CommitOpt) (err error) {

	err = k.File().BeginTransaction(ISOLATION_READ_COMMITTED)
	if err != nil {
		return
	}
	// defer function to ensure that once started,
	// we either commit transaction or abort it
	defer func() {
		// if nothing went wrong, commit
		if err == nil {
			// careful to catch error here too
			err = k.File().EndTransaction(opt)
		} else {
			// caller should see error that caused abort,
			// not success or failure of abort itself
			_ = k.File().AbortTransaction()
		}
	}()

	for _, op := range b.ops {
		if op.del {
			Log.Tracef("fdb_del_kv call k:%p db:%p kk:%v", k, k.db, op.k)
			errNo := C.fdb_del_kv(k.db, op.k, op.klen)
			Log.Tracef("fdb_del_kv retn k:%p errNo:%v", k, errNo)
			if errNo != RESULT_SUCCESS {
				err = Error(errNo)
				return
			}
		} else {
			Log.Tracef("fdb_set_kv call k:%p db:%p kk:%v v:%v", k, k.db, op.k, op.v)
			errNo := C.fdb_set_kv(k.db, op.k, op.klen, op.v, op.vlen)
			Log.Tracef("fdb_set_kv retn k:%p errNo:%v", k, errNo)
			if errNo != RESULT_SUCCESS {
				err = Error(errNo)
				return
			}
		}
	}
	return
}

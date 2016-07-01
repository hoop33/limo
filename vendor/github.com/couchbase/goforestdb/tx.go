//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package forestdb

//#include <libforestdb/forestdb.h>
import "C"

// IsolationLevel is the Transaction Isolation Level
type IsolationLevel uint8

const (
	//ISOLATION_SERIALIZABLE IsolationLevel = 0x00
	//ISOLATION_REPEATABLE_READ IsolationLevel = 0x01
	ISOLATION_READ_COMMITTED   IsolationLevel = 0x02
	ISOLATION_READ_UNCOMMITTED IsolationLevel = 0x03
)

func (f *File) BeginTransaction(level IsolationLevel) error {
	Log.Tracef("fdb_begin_transaction call f:%p dbfile:%p level:%v", f, f.dbfile, level)
	errNo := C.fdb_begin_transaction(f.dbfile, C.fdb_isolation_level_t(level))
	Log.Tracef("fdb_begin_transaction retn f:%p errNo:%v", f, errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

func (f *File) EndTransaction(opt CommitOpt) error {
	Log.Tracef("fdb_end_transaction call f:%p dbfile:%p opt:%v", f, f.dbfile, opt)
	errNo := C.fdb_end_transaction(f.dbfile, C.fdb_commit_opt_t(opt))
	Log.Tracef("fdb_end_transaction retn f:%p errNo:%v", f, errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

func (f *File) AbortTransaction() error {
	Log.Tracef("fdb_abort_transaction call f:%p dbfile:%p", f, f.dbfile)
	errNo := C.fdb_abort_transaction(f.dbfile)
	Log.Tracef("fdb_abort_transaction retn f:%p errNo:%v", f, errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

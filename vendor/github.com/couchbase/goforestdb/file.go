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

// Database handle
type File struct {
	dbfile *C.fdb_file_handle
}

// Init initializes forestdb library
func Init(config *Config) error {
	if config == nil {
		config = DefaultConfig()
	}

	errNo := C.fdb_init(config.config)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// Open opens the database with a given file name
func Open(filename string, config *Config) (*File, error) {

	if config == nil {
		config = DefaultConfig()
	}

	dbname := C.CString(filename)
	defer C.free(unsafe.Pointer(dbname))

	rv := File{}
	Log.Tracef("fdb_open call rv:%p dbname:%v conf:%v", &rv, dbname, config.config)
	errNo := C.fdb_open(&rv.dbfile, dbname, config.config)
	Log.Tracef("fdb_open ret rv:%p errNo:%v rv:%v", &rv, errNo, rv)
	if errNo != RESULT_SUCCESS {
		return nil, Error(errNo)
	}
	return &rv, nil
}

// Options to be passed to Commit()
type CommitOpt uint8

const (
	// Perform commit without any options.
	COMMIT_NORMAL CommitOpt = 0x00
	// Manually flush WAL entries even though it doesn't reach the configured threshol
	COMMIT_MANUAL_WAL_FLUSH CommitOpt = 0x01
)

// Commit all pending changes into disk.
func (f *File) Commit(opt CommitOpt) error {
	Log.Tracef("fdb_commit call f:%p dbfile:%p opt:%v", f, f.dbfile, opt)
	errNo := C.fdb_commit(f.dbfile, C.fdb_commit_opt_t(opt))
	Log.Tracef("fdb_commit retn f:%p errNo:%v", f, errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// Compact the current database file and create a new compacted file
func (f *File) Compact(newfilename string) error {

	fn := C.CString(newfilename)
	defer C.free(unsafe.Pointer(fn))

	Log.Tracef("fdb_compact call f:%p dbfile:%p fn:%v", f, f.dbfile, fn)
	errNo := C.fdb_compact(f.dbfile, fn)
	Log.Tracef("fdb_compact retn f:%p errNo:%v", f, errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// CompactUpto compacts the current database file upto given snapshot marker
//and creates a new compacted file
func (f *File) CompactUpto(newfilename string, sm *SnapMarker) error {

	fn := C.CString(newfilename)
	defer C.free(unsafe.Pointer(fn))

	Log.Tracef("fdb_compact_upto call f:%p dbfile:%p fn:%v marker:%v", f, f.dbfile, fn, sm.marker)
	errNo := C.fdb_compact_upto(f.dbfile, fn, sm.marker)
	Log.Tracef("fdb_compact_upto retn f:%p errNo:%v", f, errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// EstimateSpaceUsed returns the overall disk space actively used by the current database file
func (f *File) EstimateSpaceUsed() int {
	Log.Tracef("fdb_estimate_space_used call f:%p dbfile:%p", f, f.dbfile)
	rv := int(C.fdb_estimate_space_used(f.dbfile))
	Log.Tracef("fdb_estimate_space_used retn f:%p rv:%v", f, rv)
	return rv
}

// DbInfo returns the information about a given database handle
func (f *File) Info() (*FileInfo, error) {
	rv := FileInfo{}
	Log.Tracef("fdb_get_file_info call f:%p dbfile:%p", f, f.dbfile)
	errNo := C.fdb_get_file_info(f.dbfile, &rv.info)
	Log.Tracef("fdb_get_file_info retn f:%p errNo:%v, info:%v", f, errNo, rv.info)
	if errNo != RESULT_SUCCESS {
		return nil, Error(errNo)
	}
	return &rv, nil
}

// SwitchCompactionMode changes the compaction mode of a ForestDB file
// if the mode is changed to COMPACT_AUTO the compaction
// threshold is set to the threshold passed to this API
func (f *File) SwitchCompactionMode(mode CompactOpt, threshold int) error {
	errNo := C.fdb_switch_compaction_mode(f.dbfile, C.fdb_compaction_mode_t(mode), C.size_t(threshold))
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// Close the database file
func (f *File) Close() error {
	Log.Tracef("fdb_close call f:%p dbfile:%p", f, f.dbfile)
	errNo := C.fdb_close(f.dbfile)
	Log.Tracef("fdb_close retn f:%p errNo:%v", f, errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// OpenKVStore opens the named KVStore within the File
// using the provided KVStoreConfig.  If config is
// nil the DefaultKVStoreConfig() will be used.
func (f *File) OpenKVStore(name string, config *KVStoreConfig) (*KVStore, error) {
	if config == nil {
		config = DefaultKVStoreConfig()
	}

	rv := KVStore{
		f:    f,
		name: name,
	}
	kvsname := C.CString(name)
	defer C.free(unsafe.Pointer(kvsname))
	Log.Tracef("fdb_kvs_open call f:%p dbfile:%p kvsname:%v config:%v", f, f.dbfile, kvsname, config.config)
	errNo := C.fdb_kvs_open(f.dbfile, &rv.db, kvsname, config.config)
	Log.Tracef("fdb_kvs_open retn f:%p errNo:%v db:%p rv:%p", f, errNo, rv.db, &rv)
	if errNo != RESULT_SUCCESS {
		return nil, Error(errNo)
	}
	return &rv, nil
}

// OpenKVStore opens the default KVStore within the File
// using the provided KVStoreConfig.  If config is
// nil the DefaultKVStoreConfig() will be used.
func (f *File) OpenKVStoreDefault(config *KVStoreConfig) (*KVStore, error) {
	return f.OpenKVStore("default", config)
}

func (f *File) GetKVStoreNames() ([]string, error) {
	var ninfo C.fdb_kvs_name_list
	errNo := C.fdb_get_kvs_name_list(f.dbfile, &ninfo)
	if errNo != RESULT_SUCCESS {
		return nil, Error(errNo)
	}

	size := int(ninfo.num_kvs_names)
	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(ninfo.kvs_names)),
		Len:  size,
		Cap:  size,
	}

	all := *(*[]*C.char)(unsafe.Pointer(&hdr))
	rv := make([]string, size)
	for i := 0; i < size; i++ {
		rv[i] = C.GoString(all[i])
	}

	C.fdb_free_kvs_name_list(&ninfo)
	if errNo != RESULT_SUCCESS {
		return nil, Error(errNo)
	}

	return rv, nil
}

// Destroy destroys all resources associated with a ForestDB file permanently
func Destroy(filename string, config *Config) error {

	if config == nil {
		config = DefaultConfig()
	}

	dbname := C.CString(filename)
	defer C.free(unsafe.Pointer(dbname))

	Log.Tracef("fdb_destroy call dbname:%v config:%v", dbname, config.config)
	errNo := C.fdb_destroy(dbname, config.config)
	Log.Tracef("fdb_destroy retn dbname:%v errNo:%v", dbname, errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// CloseFileKVStore closes the KVStore and its associated forestdb file db.
func CloseFileKVStore(kvs *KVStore) error {
	f := kvs.File()
	err0 := kvs.Close()
	err1 := f.Close()
	if err0 != nil {
		return err0
	}
	return err1
}

// OpenFileKVStore opens up a forestdb file db and a single KVStore
// inside that file/db.
func OpenFileKVStore(fileName string, config *Config,
	kvstoreName string, kvstoreConfig *KVStoreConfig) (
	*KVStore, error) {
	db, err := Open(fileName, config)
	if err != nil {
		return nil, err
	}
	kvs, err := db.OpenKVStore(kvstoreName, kvstoreConfig)
	if err != nil {
		// close the db file we just opened
		db.Close()
		return nil, err
	}
	return kvs, nil
}

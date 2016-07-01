package forestdb

//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

//#include <libforestdb/forestdb.h>
import "C"

import (
	"fmt"
)

// FileInfo stores information about a given file
type FileInfo struct {
	info C.fdb_file_info
}

func (i *FileInfo) Filename() string {
	return C.GoString(i.info.filename)
}

func (i *FileInfo) NewFilename() string {
	return C.GoString(i.info.new_filename)
}

func (i *FileInfo) DocCount() uint64 {
	return uint64(i.info.doc_count)
}

func (i *FileInfo) SpaceUsed() uint64 {
	return uint64(i.info.space_used)
}

func (i *FileInfo) FileSize() uint64 {
	return uint64(i.info.file_size)
}

func (i *FileInfo) String() string {
	return fmt.Sprintf("filename: %s new_filename: %s doc_count: %d space_used: %d file_size: %d", i.Filename(), i.NewFilename(), i.DocCount(), i.SpaceUsed(), i.FileSize())
}

// KVStoreInfo stores information about a given kvstore
type KVStoreInfo struct {
	info C.fdb_kvs_info
}

func (i *KVStoreInfo) Name() string {
	return C.GoString(i.info.name)
}

func (i *KVStoreInfo) LastSeqNum() SeqNum {
	return SeqNum(i.info.last_seqnum)
}

func (i *KVStoreInfo) DocCount() uint64 {
	return uint64(i.info.doc_count)
}

func (i *KVStoreInfo) String() string {
	return fmt.Sprintf("name: %s last_seqnum: %d", i.Name(), i.LastSeqNum())
}

type KVSOpsInfo struct {
	info C.fdb_kvs_ops_info
}

func (i *KVSOpsInfo) NumSets() uint64 {
	return uint64(i.info.num_sets)
}

func (i *KVSOpsInfo) NumDels() uint64 {
	return uint64(i.info.num_dels)
}

func (i *KVSOpsInfo) NumCommits() uint64 {
	return uint64(i.info.num_commits)
}

func (i *KVSOpsInfo) NumCompacts() uint64 {
	return uint64(i.info.num_compacts)
}

func (i *KVSOpsInfo) NumGets() uint64 {
	return uint64(i.info.num_gets)
}

func (i *KVSOpsInfo) NumIteratorGets() uint64 {
	return uint64(i.info.num_iterator_gets)
}

func (i *KVSOpsInfo) NumIteratorMoves() uint64 {
	return uint64(i.info.num_iterator_moves)
}

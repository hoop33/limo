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
	"unsafe"
)

type SeqNum uint64

// ForestDB doc structure definition
type Doc struct {
	doc *C.fdb_doc
}

// NewDoc creates a new FDB_DOC instance on heap with a given key, its metadata, and its doc body
func NewDoc(key, meta, body []byte) (*Doc, error) {
	var k, m, b unsafe.Pointer

	if len(key) != 0 {
		k = unsafe.Pointer(&key[0])
	}

	if len(meta) != 0 {
		m = unsafe.Pointer(&meta[0])
	}

	if len(body) != 0 {
		b = unsafe.Pointer(&body[0])
	}

	lenk := len(key)
	lenm := len(meta)
	lenb := len(body)

	rv := Doc{}

	Log.Tracef("fdb_doc_create call k:%p doc:%v m:%v b:%v", k, rv.doc, m, b)
	errNo := C.fdb_doc_create(&rv.doc,
		k, C.size_t(lenk), m, C.size_t(lenm), b, C.size_t(lenb))
	Log.Tracef("fdb_doc_create ret k:%p errNo:%v doc:%v", k, errNo, rv.doc)
	if errNo != RESULT_SUCCESS {
		return nil, Error(errNo)
	}
	return &rv, nil
}

// Update a FDB_DOC instance with a given metadata and body
// NOTE: does not update the database, just the in memory structure
func (d *Doc) Update(meta, body []byte) error {
	var m, b unsafe.Pointer

	if len(meta) != 0 {
		m = unsafe.Pointer(&meta[0])
	}

	if len(body) != 0 {
		b = unsafe.Pointer(&body[0])
	}

	lenm := len(meta)
	lenb := len(body)
	Log.Tracef("fdb_doc_update call d:%p doc:%v m:%v b:%v", d, d.doc, m, b)
	errNo := C.fdb_doc_update(&d.doc, m, C.size_t(lenm), b, C.size_t(lenb))
	Log.Tracef("fdb_doc_update retn d:%p errNo:%v doc:%v m:%v b:%v", d, errNo, d.doc, m, b)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

// Key returns the document key
func (d *Doc) Key() []byte {
	return C.GoBytes(d.doc.key, C.int(d.doc.keylen))
}

// Meta returns the document metadata
func (d *Doc) Meta() []byte {
	return C.GoBytes(d.doc.meta, C.int(d.doc.metalen))
}

// Body returns the document body
func (d *Doc) Body() []byte {
	return C.GoBytes(d.doc.body, C.int(d.doc.bodylen))
}

// SeqNum returns the document sequence number
func (d *Doc) SeqNum() SeqNum {
	return SeqNum(d.doc.seqnum)
}

// SetSeqNum sets the document sequence number
// NOTE: only to be used when initiating a sequence number lookup
func (d *Doc) SetSeqNum(sn SeqNum) {
	C.fdb_doc_set_seqnum(d.doc, C.fdb_seqnum_t(sn))
}

// Offset returns the offset position on disk
func (d *Doc) Offset() uint64 {
	return uint64(d.doc.offset)
}

// Deleted returns whether or not this document has been deleted
func (d *Doc) Deleted() bool {
	return bool(d.doc.deleted)
}

// Close releases resources allocated to this document
func (d *Doc) Close() error {
	Log.Tracef("fdb_doc_free call d:%p doc:%v", d, d.doc)
	errNo := C.fdb_doc_free(d.doc)
	Log.Tracef("fdb_doc_free retn d:%p errNo:%v", d, errNo)
	if errNo != RESULT_SUCCESS {
		return Error(errNo)
	}
	return nil
}

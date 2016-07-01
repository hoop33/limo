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
	"bytes"

	"github.com/couchbase/goforestdb"
)

type Iterator struct {
	store    *Store
	snapshot *forestdb.KVStore
	iterator *forestdb.Iterator
	curr     *forestdb.Doc
	valid    bool
	start    []byte
	parent   *Reader
}

func (i *Iterator) Seek(key []byte) {
	if key == nil {
		key = []byte{0}
	}
	if i.start != nil && bytes.Compare(key, i.start) < 0 {
		key = i.start
	}
	err := i.iterator.Seek(key, forestdb.FDB_ITR_SEEK_HIGHER)
	if err != nil {
		i.valid = false
		return
	} else {
		i.valid = true // seek can make invalid iterator valid again
	}
	if i.curr != nil {
		i.curr.Close()
	}
	i.curr, err = i.iterator.Get()
	if err != nil {
		i.valid = false
		return
	}
}

func (i *Iterator) Next() {
	err := i.iterator.Next()
	if err != nil {
		i.valid = false
		return
	}
	if i.curr != nil {
		i.curr.Close()
	}
	i.curr, err = i.iterator.Get()
	if err != nil {
		i.valid = false
	}
}

func (i *Iterator) Current() ([]byte, []byte, bool) {
	if i.Valid() {
		return i.Key(), i.Value(), true
	}
	return nil, nil, false
}

func (i *Iterator) Key() []byte {
	return i.curr.Key()
}

func (i *Iterator) Value() []byte {
	return i.curr.Body()
}

func (i *Iterator) Valid() bool {
	return i.valid
}

func (i *Iterator) Close() error {
	i.valid = false
	if i.curr != nil {
		i.curr.Close()
	}
	if i.parent != nil {
		i.parent.decRef()
	}
	return i.iterator.Close()
}

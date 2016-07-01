//  Copyright (c) 2016 Couchbase, Inc.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the
//  License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing,
//  software distributed under the License is distributed on an "AS
//  IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
//  express or implied. See the License for the specific language
//  governing permissions and limitations under the License.

package cellar

import (
	"bytes"

	"github.com/boltdb/bolt"
)

type mergeCursor struct {
	reader           *reader
	mutationsCursors []*bolt.Cursor
	deletionsCursors []*bolt.Cursor

	key       [][]byte
	val       [][]byte
	curr      int
	dkey      [][]byte
	currIsDel bool
}

func newMergeCursor(reader *reader) *mergeCursor {
	rv := &mergeCursor{
		reader:           reader,
		mutationsCursors: make([]*bolt.Cursor, 0, len(reader.root)),
		deletionsCursors: make([]*bolt.Cursor, 0, len(reader.root)),
		key:              make([][]byte, len(reader.root)),
		val:              make([][]byte, len(reader.root)),
		dkey:             make([][]byte, len(reader.root)),
	}

	for _, mutationsBucket := range reader.mutations {
		mutationsCursor := mutationsBucket.Cursor()
		rv.mutationsCursors = append(rv.mutationsCursors, mutationsCursor)
	}

	for _, deletionsBucket := range reader.deletions {
		deletionsCursor := deletionsBucket.Cursor()
		rv.deletionsCursors = append(rv.deletionsCursors, deletionsCursor)
	}

	return rv
}

func (c *mergeCursor) Seek(seek []byte) (key []byte, value []byte, deleted bool) {
	for i, cursor := range c.mutationsCursors {
		c.key[i], c.val[i] = cursor.Seek(seek)
	}
	for i, cursor := range c.deletionsCursors {
		c.dkey[i], _ = cursor.Seek(seek)
	}
	c.updateCurr()
	if c.currIsDel {
		return c.dkey[c.curr], nil, true
	}
	return c.key[c.curr], c.val[c.curr], false
}

func (c *mergeCursor) next() {
	currKey := c.key[c.curr]
	if c.currIsDel {
		currKey = c.dkey[c.curr]
	}
	// increment any cursor pointing at the
	// current key (could be more than just 1)
	for i, cursor := range c.mutationsCursors {
		if bytes.Compare(currKey, c.key[i]) == 0 {
			c.key[i], c.val[i] = cursor.Next()
		}
	}
	for i, cursor := range c.deletionsCursors {
		if bytes.Compare(currKey, c.dkey[i]) == 0 {
			c.dkey[i], _ = cursor.Next()
		}
	}
}

func (c *mergeCursor) Next() (key []byte, value []byte, deleted bool) {
	c.next()
	c.updateCurr()
	if c.currIsDel {
		return c.dkey[c.curr], nil, true
	}
	return c.key[c.curr], c.val[c.curr], false
}

func (c *mergeCursor) updateCurr() {
	// find curr (iterator index with lowest key, always prefering first seen)
	var currKey []byte
	c.curr = 0
	for i, k := range c.key {
		dk := c.dkey[i]
		if dk != nil && (currKey == nil || bytes.Compare(dk, currKey) < 0) {
			currKey = dk
			c.curr = i
			c.currIsDel = true
		}
		if k != nil && (currKey == nil || bytes.Compare(k, currKey) < 0) {
			currKey = k
			c.curr = i
			c.currIsDel = false
		}
	}
}

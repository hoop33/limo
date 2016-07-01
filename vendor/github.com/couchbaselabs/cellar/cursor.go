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

// Cursor is a tool for iterating through k/v pairs in the cellar
type Cursor struct {
	reader           *reader
	mutationsCursors []*bolt.Cursor
	deletionsCursors []*bolt.Cursor

	key  [][]byte
	val  [][]byte
	curr int
}

func newCursor(reader *reader) *Cursor {
	rv := &Cursor{
		reader:           reader,
		mutationsCursors: make([]*bolt.Cursor, 0, len(reader.root)),
		deletionsCursors: make([]*bolt.Cursor, 0, len(reader.root)),
		key:              make([][]byte, len(reader.root)),
		val:              make([][]byte, len(reader.root)),
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

// Seek moves the cursor to the specified key
func (c *Cursor) Seek(seek []byte) (key []byte, value []byte) {
	for i, cursor := range c.mutationsCursors {
		c.key[i], c.val[i] = cursor.Seek(seek)
	}
	c.updateCurr()
	for c.checkCurrDeleted() {
		c.next()
		c.updateCurr()
	}
	return c.key[c.curr], c.val[c.curr]
}

func (c *Cursor) next() {
	currKey := c.key[c.curr]
	// increment any cursor pointing at the
	// current key (could be more than just 1)
	for i, cursor := range c.mutationsCursors {
		if bytes.Compare(currKey, c.key[i]) == 0 {
			c.key[i], c.val[i] = cursor.Next()
		}
	}
}

// Next moves the cursor to the next key
func (c *Cursor) Next() (key []byte, value []byte) {
	c.next()
	c.updateCurr()
	for c.checkCurrDeleted() {
		c.next()
		c.updateCurr()
	}
	// FIXME end case where nothing else?
	return c.key[c.curr], c.val[c.curr]
}

func (c *Cursor) updateCurr() {
	// find curr (iterator index with lowest key, always prefering first seen)
	var currKey []byte
	c.curr = 0
	for i, k := range c.key {
		if k != nil && (currKey == nil || bytes.Compare(k, currKey) < 0) {
			currKey = k
			c.curr = i
		}
	}
}

func (c *Cursor) checkCurrDeleted() bool {
	currKey := c.key[c.curr]
	// seek all the deletion cusors on previous segments to the current key
	for _, deletionCursor := range c.deletionsCursors[:c.curr] {
		k, _ := deletionCursor.Seek(currKey)
		if bytes.Compare(k, currKey) == 0 {
			// one of the segments with higher priority has this key deleted
			return true
		}
	}
	return false
}

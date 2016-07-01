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

// Tx represents a cellar transaction
type Tx struct {
	cellar         *Cellar
	writable       bool
	managed        bool
	segmentBuilder *segmentBuilder
	root           segmentList
	reader         *reader
}

// Rollback will abort this transaction, none of the operations performed
// up to this point will be reflected in the cellar
func (tx *Tx) Rollback() error {
	if tx.managed {
		return ErrTxIsManaged
	}
	if tx.cellar == nil {
		return ErrTxClosed
	}
	return tx.rollback()
}

func (tx *Tx) rollback() error {
	if tx.cellar == nil {
		return nil
	}
	if tx.writable {
		err := tx.segmentBuilder.Abort()
		if err != nil {
			return err
		}
	}
	return tx.close()
}

func (tx *Tx) close() error {
	if tx.cellar == nil {
		return nil
	}
	if tx.writable {
		tx.cellar.rwlock.Unlock()
	}
	if tx.root != nil {
		for _, segment := range tx.root {
			segment.decrRef("tx closing")
		}
	}
	if tx.reader != nil {
		err := tx.reader.Close()
		if err != nil {
			return err
		}
		tx.reader = nil
	}

	// Clear all references.
	tx.cellar = nil
	return nil
}

// Commit will atomically make all the operations in the transactions a
// part of this cellar
func (tx *Tx) Commit() error {
	if tx.managed {
		return ErrTxIsManaged
	}
	if tx.cellar == nil {
		return ErrTxClosed
	} else if !tx.writable {
		return ErrTxNotWritable
	}
	// build the new segment
	newSegmentPath := tx.segmentBuilder.db.Path()
	err := tx.segmentBuilder.Build()
	if err != nil {
		return err
	}
	newsegment, err := openSegmentPath(newSegmentPath)
	if err != nil {
		return err
	}
	// make this segment live
	err = tx.cellar.pushRoot(newsegment)
	if err != nil {
		return err
	}
	return tx.close()
}

// Get will look up the specified key
// if theere is no value, nil is returned
// NOTE: an empty byte slice is a valid value, and not the same as nil
func (tx *Tx) Get(key []byte) []byte {
	return tx.reader.Get(key)
}

// Cursor returns an object which can be used to iterate k/v paris in the cellar
func (tx *Tx) Cursor() *Cursor {
	return newCursor(tx.reader)
}

// Delete will remove the key from the cellar
func (tx *Tx) Delete(key []byte) error {
	if tx.cellar == nil {
		return ErrTxClosed
	} else if !tx.writable {
		return ErrTxNotWritable
	}
	return tx.segmentBuilder.Delete(key)
}

// Put will update the value for the specified key
func (tx *Tx) Put(key []byte, value []byte) error {
	if tx.cellar == nil {
		return ErrTxClosed
	} else if !tx.writable {
		return ErrTxNotWritable
	}
	return tx.segmentBuilder.Put(key, value)
}

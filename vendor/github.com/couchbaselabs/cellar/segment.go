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
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/boltdb/bolt"
)

var segmentOpts = bolt.Options{
	ReadOnly: true,
}

const segmentPrefix = "cellar-"

var metaBucketName = []byte("x")
var mutationsBucketName = []byte("m")
var deletionsBucketName = []byte("d")
var seqKeyName = []byte("seq")

// Segment is a read-only bolt.DB, plus some extra bookkeeping
type segment struct {
	*bolt.DB
	seq uint64

	mergeInProgress uint64

	refsCond *sync.Cond
	refsLock sync.Mutex
	refs     uint64
}

func (s *segment) incrRef(reason ...string) {
	s.refsLock.Lock()
	defer s.refsLock.Unlock()
	s.refs++
	Logger.Printf("incr count for %d to %d - for %v", s.seq, s.refs, reason)
	s.refsCond.Broadcast()
}

func (s *segment) decrRef(reason ...string) {
	s.refsLock.Lock()
	defer s.refsLock.Unlock()
	s.refs--
	Logger.Printf("decr count for %d to %d - for %v", s.seq, s.refs, reason)
	s.refsCond.Broadcast()
}

func (s *segment) String() string {
	return fmt.Sprintf("{seq: %d}", s.seq)
}

func segmentFilename(seq uint64) string {
	return fmt.Sprintf("%s%016x", segmentPrefix, seq)
}

func openSegment(cellarPath string, seq uint64) (*segment, error) {
	path := fmt.Sprintf("%s%s%s", cellarPath, string(os.PathSeparator), segmentFilename(seq))
	return openSegmentPath(path)
}

func openSegmentPath(path string) (*segment, error) {
	db, err := bolt.Open(path, 0600, &segmentOpts)
	if err != nil {
		return nil, err
	}
	rv := &segment{
		DB: db,
	}
	rv.refsCond = sync.NewCond(&rv.refsLock)
	// read the sequence number out of the metadata
	err = db.View(func(tx *bolt.Tx) error {
		meta := tx.Bucket(metaBucketName)
		seqBytes := meta.Get(seqKeyName)
		segmentSeq, err := strconv.ParseUint(string(seqBytes), 16, 64)
		if err != nil {
			return err
		}
		rv.seq = segmentSeq
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rv, nil
}

func (s *segment) Seq() uint64 {
	return s.seq
}

func (s *segment) Close() error {
	s.refsLock.Lock()
	for s.refs > 0 {
		Logger.Printf("close %d - waiting for refs, currently: %d", s.seq, s.refs)
		s.refsCond.Wait()
	}
	Logger.Printf("closed %d - waiting for refs, currently: %d", s.seq, s.refs)
	s.refsLock.Unlock()

	return s.DB.Close()
}

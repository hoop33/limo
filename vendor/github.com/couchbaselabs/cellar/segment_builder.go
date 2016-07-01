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

	"github.com/boltdb/bolt"
)

type segmentBuilder struct {
	db        *bolt.DB
	tx        *bolt.Tx
	mutations *bolt.Bucket
	deletions *bolt.Bucket
	metadata  *bolt.Bucket
}

func newSegmentBuilder(cellarPath string, seq uint64) (*segmentBuilder, error) {
	path := fmt.Sprintf("%s/%s%016x", cellarPath, segmentPrefix, seq)
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, fmt.Errorf("newSegmentBuilder Open: %v", err)
	}
	// set no-sync since we're bulk loading it
	db.NoSync = true
	// try to start a write transaction
	tx, err := db.Begin(true)
	if err != nil {
		return nil, fmt.Errorf("newSegmentBuilder Begin: %v", err)
	}
	mutations, err := tx.CreateBucketIfNotExists(mutationsBucketName)
	if err != nil {
		return nil, fmt.Errorf("newSegmentBuilder CreateBucketIfNotExists '%s': %v", mutationsBucketName, err)
	}
	mutations.FillPercent = 1.0
	deletions, err := tx.CreateBucketIfNotExists(deletionsBucketName)
	if err != nil {
		return nil, fmt.Errorf("newSegmentBuilder CreateBucketIfNotExists '%s': %v", deletionsBucketName, err)
	}
	deletions.FillPercent = 1.0
	metadata, err := tx.CreateBucketIfNotExists(metaBucketName)
	if err != nil {
		return nil, fmt.Errorf("newSegmentBuilder CreateBucketIfNotExists '%s': %v", metaBucketName, err)
	}
	err = metadata.Put(seqKeyName, []byte(fmt.Sprintf("%016x", seq)))
	if err != nil {
		return nil, fmt.Errorf("newSegmentBuilder Put '%s': %v", seqKeyName, err)
	}
	return &segmentBuilder{
		db:        db,
		tx:        tx,
		mutations: mutations,
		deletions: deletions,
		metadata:  metadata,
	}, nil
}

func (s *segmentBuilder) Put(key, val []byte) error {
	err := s.mutations.Put(key, val)
	if err != nil {
		return fmt.Errorf("segmentBuilder Put: %v", err)
	}
	return nil
}

func (s *segmentBuilder) PutMetadata(key, val []byte) error {
	err := s.metadata.Put(key, val)
	if err != nil {
		return fmt.Errorf("segmentBuilder PutMetadata: %v", err)
	}
	return nil
}

func (s *segmentBuilder) Delete(key []byte) error {
	err := s.deletions.Put(key, []byte{})
	if err != nil {
		return fmt.Errorf("segmentBuilder Delete: %v", err)
	}
	return nil
}

func (s *segmentBuilder) Build() error {
	err := s.tx.Commit()
	if err != nil {
		return fmt.Errorf("segmentBuilder Build Commit: %v", err)
	}
	err = s.db.Close()
	if err != nil {
		return fmt.Errorf("segmentBuilder Build Close: %v", err)
	}
	return nil
}

func (s *segmentBuilder) Abort() error {
	cleanupPath := s.db.Path()
	err := s.tx.Rollback()
	if err != nil {
		return fmt.Errorf("segmentBuilder Abort Rollback: %v", err)
	}
	err = s.db.Close()
	if err != nil {
		return fmt.Errorf("segmentBuilder Abort Close: %v", err)
	}
	err = os.Remove(cleanupPath)
	if err != nil {
		return fmt.Errorf("segmentBuilder Abort Remove: %v", err)
	}
	return nil
}

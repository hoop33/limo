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

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"

	"github.com/boltdb/bolt"
)

// NOTE: we intentionally copy/pasted code here to avoid exposing internals
// of cellar to the rest of the world

var readOnly = bolt.Options{
	ReadOnly: true,
}

func main() {
	flag.Parse()

	cellarPath := flag.Arg(0)

	db, err := bolt.Open(fmt.Sprintf("%s/%s", cellarPath, "master.db"), 0600, &readOnly)
	if err != nil {
		log.Fatalf("error opening cellar master db: %v", err)
	}

	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("m"))
		if bucket == nil {
			log.Fatal("cellar master db contains does not contain master bucket 'm'")
		}
		root := bucket.Get([]byte("root"))
		if root == nil {
			log.Fatal("celler master bucker does not contain root key 'root'")
		}
		rootSeqs, err := parseRoot(root)
		if err != nil {
			log.Fatalf("error parsing cellar root sequences: %v", err)
		}
		fmt.Printf("Cellar root sequences: %v\n", rootSeqs)

		return nil
	})
}

func parseRoot(val []byte) ([]uint64, error) {
	var rv []uint64
	buf := bytes.NewBuffer(val)
	var next uint64
	err := binary.Read(buf, binary.BigEndian, &next)
	for err == nil {
		rv = append(rv, next)
		err = binary.Read(buf, binary.BigEndian, &next)
	}
	if err != nil && err != io.EOF {
		return nil, err
	}
	return rv, nil
}

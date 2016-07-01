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
	"encoding/binary"
	"io"
)

type segmentList []*segment

func (s segmentList) MarshalBinary() (data []byte, err error) {
	var buf bytes.Buffer
	for _, segment := range s {
		err := binary.Write(&buf, binary.BigEndian, segment.Seq())
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
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

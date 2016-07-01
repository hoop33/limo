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

import "os"

// Merge represents an ordered set of adjacent segments to be merged
// dropDeletes specifies whether or not the deletes should be dropped
// deletes can only be dropped if the result of the merge is the final segment
type Merge struct {
	cellar        *Cellar
	newSegmentSeq uint64
	sources       segmentList
	dropDeletes   bool
}

func doMerge(m *Merge) error {

	segmentBuilder, err := newSegmentBuilder(m.cellar.path, m.newSegmentSeq)
	if err != nil {
		return err
	}

	r, err := newReader(m.sources)
	if err != nil {
		return err
	}

	c := newMergeCursor(r)
	k, v, deleted := c.Seek([]byte{})
	for k != nil {
		var err error
		if deleted && !m.dropDeletes {
			err = segmentBuilder.Delete(k)
		} else if !deleted {
			err = segmentBuilder.Put(k, v)
		}
		if err != nil {
			return err
		}
		k, v, deleted = c.Next()
	}
	err = r.Close()
	if err != nil {
		return err
	}
	// decr refs on the sources		// release refs
	for _, segment := range m.sources {
		segment.decrRef("merge done")
	}

	newSegmentPath := segmentBuilder.db.Path()
	err = segmentBuilder.Build()
	if err != nil {
		return err
	}
	newsegment, err := openSegmentPath(newSegmentPath)
	if err != nil {
		return err
	}
	// make this segment live
	err = m.cellar.replaceSegments(m.sources, newsegment)
	if err != nil {
		segmentPath := newsegment.Path()
		// this segment was not made live, close it, delete it
		_ = newsegment.Close()
		// FIXME currently the merge segment will never be used
		// so we delete it, however in the future, this could change
		_ = os.RemoveAll(segmentPath)
		return err
	}

	return nil
}

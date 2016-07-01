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

// MergePolicy is anything which can prescribe a set of Merges to be done
type MergePolicy interface {
	Merges(*Cellar, segmentList) []*Merge
}

// SimpleMergePolicy has no brain at all, it simply always chooses to merge
// two consecutive segments that are not already being merged
type SimpleMergePolicy struct{}

// Merges returns the set of prescribed merge operations for this set of segments
func (s *SimpleMergePolicy) Merges(cellar *Cellar, segments segmentList) []*Merge {
	//rv := make([]*Merge, 0)
	var rv []*Merge
	consecutive := make(segmentList, 0)
	for i := len(segments) - 1; i >= 0; i-- {
		segment := segments[i]
		if segment.mergeInProgress == 0 {
			// insert, not append to keep the order the same (we're iterating reverse)
			consecutive = append(consecutive, nil)
			copy(consecutive[1:], consecutive[:])
			consecutive[0] = segment
			if len(consecutive) == 2 {
				rv = append(rv, &Merge{
					cellar:      cellar,
					sources:     consecutive,
					dropDeletes: (i == len(segments)-2), // if merging last 2 segments, we can drop deletes
				})
				consecutive = make(segmentList, 0)
			}
		} else {
			// found a segment with merge in progress, so start over
			consecutive = make(segmentList, 0)
		}
	}

	return rv
}

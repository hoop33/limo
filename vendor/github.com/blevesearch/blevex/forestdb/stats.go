//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package forestdb

import "encoding/json"

type kvStat struct {
	s *Store
}

func (k *kvStat) statsMap() map[string]interface{} {
	k.s.statsMutex.Lock()
	defer k.s.statsMutex.Unlock()

	if k.s.statsHandle == nil {
		return map[string]interface{}{}
	}

	f := k.s.statsHandle.File()
	finfo, err := f.Info()
	if err != nil {
		return map[string]interface{}{}
	}

	opsInfo, err := k.s.statsHandle.OpsInfo()
	if err != nil {
		return map[string]interface{}{}
	}

	m := map[string]interface{}{}
	m["sets"] = opsInfo.NumSets()
	m["dels"] = opsInfo.NumDels()
	m["commits"] = opsInfo.NumCommits()
	m["compacts"] = opsInfo.NumCompacts()
	m["gets"] = opsInfo.NumGets()
	m["iterator_gets"] = opsInfo.NumIteratorGets()
	m["iterator_moves"] = opsInfo.NumIteratorMoves()
	m["space_used"] = finfo.SpaceUsed()
	m["file_size"] = finfo.FileSize()

	return m
}

func (k *kvStat) MarshalJSON() ([]byte, error) {
	m := k.statsMap()
	return json.Marshal(m)
}

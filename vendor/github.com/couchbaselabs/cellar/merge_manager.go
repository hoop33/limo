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
	"sync"
	"sync/atomic"
)

type mergeManager struct {
	cellar        *Cellar
	changes       chan segmentList
	closeChan     chan struct{}
	mutex         sync.Mutex
	running       bool
	wg            sync.WaitGroup
	policy        MergePolicy
	auto          bool
	mergeWork     chan *Merge
	maxConcurrent int
}

func newMergeManager(cellar *Cellar, policy MergePolicy, auto bool, maxConcurrent int) *mergeManager {
	rv := &mergeManager{
		changes:       make(chan segmentList),
		closeChan:     make(chan struct{}),
		running:       false,
		policy:        policy,
		auto:          auto,
		mergeWork:     make(chan *Merge, 1024),
		cellar:        cellar,
		maxConcurrent: maxConcurrent,
	}

	return rv
}

func (m *mergeManager) Start() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.running {
		return fmt.Errorf("already running")
	}
	m.running = true
	// start workers
	for i := 0; i < m.maxConcurrent; i++ {
		m.wg.Add(1)
		go mergeWorker(&m.wg, m.mergeWork, m.closeChan)
	}
	m.wg.Add(1)
	go m.Run()
	return nil
}

func (m *mergeManager) Run() {
OUTER:
	for {
		select {
		case <-m.closeChan:
			break OUTER
		case newRoot, ok := <-m.changes:
			if !ok {
				break OUTER
			} else {
				// roots changed, see if any merges should be done
				merges := m.policy.Merges(m.cellar, newRoot)
				for _, merge := range merges {
					// assign this merge a new segment seq
					merge.newSegmentSeq = atomic.AddUint64(&m.cellar.seq, 1)
					// set mergeInProgress so we don't keep merging the same segments
					for _, s := range merge.sources {
						s.mergeInProgress = merge.newSegmentSeq
						// also incr ref count for each source
						s.incrRef("merge work")
					}
					m.mergeWork <- merge
				}
				// release refs
				for _, segment := range newRoot {
					segment.decrRef("merge notification")
				}
			}
		}
	}
	m.wg.Done()
}

func (m *mergeManager) Stop() error {
	Logger.Printf("staring to close merge manager")
	close(m.closeChan)
	m.wg.Wait()
	Logger.Printf("done closing merge manager")
	// drain anything left on merge work and decr refs
	close(m.mergeWork)
OUTER:
	for {
		select {
		case work, ok := <-m.mergeWork:
			if !ok {
				break OUTER
			}
			for _, s := range work.sources {
				s.decrRef("draining merge work")
			}
		}
	}
	return nil
}

func (m *mergeManager) RootChange(root segmentList) {
	if m.auto {
		m.changes <- root
	} else {
		// release refs
		for _, segment := range root {
			segment.decrRef("no automatic mergin")
		}
	}
}

func (m *mergeManager) ForceMerge(root segmentList) {
	m.changes <- root
}

func mergeWorker(wg *sync.WaitGroup, workChan chan *Merge, closeChan chan struct{}) {
OUTER:
	for {
		select {
		case <-closeChan:
			break OUTER
		case work, ok := <-workChan:
			if !ok {
				break OUTER
			}
			err := doMerge(work)
			if err != nil {
				Logger.Printf("MERGE ERROR: %v", err)
			}
		}
	}
	wg.Done()
	Logger.Printf("merge worker done")
}

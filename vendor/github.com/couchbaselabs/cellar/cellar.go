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
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/boltdb/bolt"
)

// Logger is a configurable logger used by this package
// by default output is discarded
var Logger = log.New(ioutil.Discard, "cellar ", log.LstdFlags)

const masterDbName = "master.db"
const crawlSpaceName = ".crawlspace"

var masterBucketName = []byte("m")
var rootKeyName = []byte("root")

// Options let you change configurable behavior within the cellar
type Options struct {
	AutomaticMerge bool
}

// DefaultOptions give the standard cellar behavior
var DefaultOptions = &Options{
	AutomaticMerge: true,
}

// Cellar is a merged-multi-segment(bolt) k/v store
type Cellar struct {
	path   string
	rwlock sync.Mutex // ONLY used for enforcing single writer

	seq    uint64
	master *bolt.DB

	// don't access these directly prefer getRoot/pushRoot/replaceSegments
	root     atomic.Value
	rootLock sync.RWMutex

	mergeManager *mergeManager

	stats Stats
}

// Open is used to create/open a cellar
// path should be a directory to hold the cellar contents
// if options is nil, DefaultOptions will be used
func Open(path string, options *Options) (*Cellar, error) {
	if options == nil {
		options = DefaultOptions
	}

	// make preceeding path elements if necessary
	err := os.MkdirAll(path, 0700)
	if err != nil {
		return nil, err
	}
	// make crawlspace
	err = os.MkdirAll(fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), crawlSpaceName), 0700)
	if err != nil {
		return nil, err
	}
	// map all the segments in this path
	abandonedSegmentFiles := make(map[string]uint64)
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading cellar path entries: %v", err)
	}
	for _, fileInfo := range fileInfos {
		filename := fileInfo.Name()
		// segments all have the pattern cellar-0000000000000000
		if strings.HasPrefix(filename, segmentPrefix) && len(filename) == len(segmentPrefix)+16 {
			fileseq, err := strconv.ParseUint(filename[len(segmentPrefix):len(segmentPrefix)+16], 16, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing segment filename seq: %v", err)
			}
			abandonedSegmentFiles[fileInfo.Name()] = fileseq
		}
	}

	// open the master db, creating it if necessary
	db, err := bolt.Open(fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), masterDbName), 0600, nil)
	if err != nil {
		return nil, err
	}

	rv := &Cellar{
		path:   path,
		master: db,
	}

	// read
	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("m"))
		if err != nil {
			return err
		}
		root := bucket.Get([]byte("root"))
		if root == nil {
			initialRoot := make(segmentList, 0)
			initialRootBytes, err := initialRoot.MarshalBinary()
			if err != nil {
				return err
			}
			err = bucket.Put(rootKeyName, initialRootBytes)
			if err != nil {
				return err
			}
			// initialize atomic value to empty root with correct type
			rv.root.Store(initialRoot)
		} else {
			rootSeqs, err := parseRoot(root)
			if err != nil {
				return fmt.Errorf("error parsing cellar root sequences: %v", err)
			}
			root := make(segmentList, 0)
			for _, seq := range rootSeqs {
				segment, err := openSegment(path, seq)
				if err != nil {
					return err
				}
				// remove this file name from abandonedSegmentFiles
				delete(abandonedSegmentFiles, segmentFilename(seq))
				// bump the ref count for the segment
				// this ensures a segment on the root, always has at least 1 ref
				segment.incrRef("on root")
				root = append(root, segment)
				if seq > rv.seq {
					rv.seq = seq
				}
			}
			rv.root.Store(root)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// anything left in abandonedSegmentFiles will be moved to the crawlspace
	for abandonedSegment, abandonedSegmentSeq := range abandonedSegmentFiles {
		// ensure we don't reuse seqs that were abandoned
		if abandonedSegmentSeq > rv.seq {
			rv.seq = abandonedSegmentSeq
		}
		// move it, again to prevent accidental reuse
		segmentPath := fmt.Sprintf("%s%s%s", path, string(os.PathSeparator), abandonedSegment)
		segmentCrawlSpacePath := fmt.Sprintf("%s%s%s%s%s", path, string(os.PathSeparator), crawlSpaceName, string(os.PathSeparator), abandonedSegment)
		err := os.Rename(segmentPath, segmentCrawlSpacePath)
		if err != nil {
			Logger.Printf("error moving segment %s to crawlspace: %v", abandonedSegment, err)
		}
	}

	rv.mergeManager = newMergeManager(rv, &SimpleMergePolicy{}, options.AutomaticMerge, 2)
	err = rv.mergeManager.Start()
	if err != nil {
		return nil, err
	}

	return rv, nil
}

// Begin starts a new transaction
// writable controls whether or not this transaction supports Put/Delete
func (c *Cellar) Begin(writable bool) (*Tx, error) {

	var segmentBuilder *segmentBuilder
	if writable {
		c.rwlock.Lock()

		nextSeq := atomic.AddUint64(&c.seq, 1)
		var err error
		segmentBuilder, err = newSegmentBuilder(c.path, nextSeq)
		if err != nil {
			return nil, fmt.Errorf("cellar Begin newSegmentBuilder: %v", err)
		}
	}

	root := c.getRoot("cellar begin")
	reader, err := newReader(root)
	if err != nil {
		return nil, fmt.Errorf("cellar Begin newReader: %v", err)
	}

	return &Tx{
		cellar:         c,
		writable:       writable,
		managed:        false,
		segmentBuilder: segmentBuilder,
		root:           root,
		reader:         reader,
	}, nil
}

// Close will release all resources associated with the cellar
// close may block while waiting for readers/mergers to complete
// or reach a resumable point
func (c *Cellar) Close() error {
	Logger.Printf("cellar closing")

	// set master to nil, this signals to stop accepting mutations to root
	c.rootLock.Lock()
	master := c.master
	c.master = nil
	c.rootLock.Unlock()

	var err error
	// stop the merger
	Logger.Printf("telling merge manager to stop")
	err = c.mergeManager.Stop()

	// at this point no one else is racing to change the root, get the final
	// we don't use getRoot() because that furhter incrs the refs
	// ref should be 1 for anything on root already
	c.rootLock.Lock()
	croot := c.root.Load().(segmentList)
	// now replace the root with an empty one
	nroot := make(segmentList, 0)
	c.root.Store(nroot)
	c.rootLock.Unlock()

	// now that these segments are off the root, decr refs
	for _, segment := range croot {
		segment.decrRef("cellar closing")
	}

	for _, segment := range croot {
		Logger.Printf("closing %d", segment.seq)
		cerr := segment.Close()
		if cerr != nil && err == nil {
			err = cerr
		}
	}

	// finally close master
	merr := master.Close()
	if merr != nil && err == nil {
		err = merr
	}

	return merr
}

// GoString returns the Go string representation of the cellar.
func (c *Cellar) GoString() string {
	return fmt.Sprintf("cellar.Cellar{path:%q}", c.path)
}

// Path returns the path to root of the cellar.
func (c *Cellar) Path() string {
	return c.path
}

// String returns the string representation of the cellar.
func (c *Cellar) String() string {
	return fmt.Sprintf("Cellar<%q>", c.path)
}

// Update starts a managed transaction for wriring
// the provided function is executed within the transaction
// to abort the transaction the function should return an error
// to commit the transaction the function should return nil
func (c *Cellar) Update(fn func(*Tx) error) error {
	t, err := c.Begin(true)
	if err != nil {
		return fmt.Errorf("cellar Update Begin: %v", err)
	}

	// Make sure the transaction rolls back in the event of a panic.
	defer func() {
		if t.cellar != nil {
			_ = t.rollback()
		}
	}()

	// Mark as a managed tx so that the inner function cannot manually rollback.
	t.managed = true

	// If an error is returned from the function then rollback and return error.
	err = fn(t)
	t.managed = false
	if err != nil {
		_ = t.Rollback()
		return fmt.Errorf("cellar Update Rollback: %v", err)
	}

	err = t.Commit()
	if err != nil {
		return fmt.Errorf("cellar Update Commit: %v", err)
	}
	return nil
}

// View starts a managed transaction for reading
// the provided function is executed within the transaction
func (c *Cellar) View(fn func(*Tx) error) error {
	t, err := c.Begin(false)
	if err != nil {
		return err
	}

	// Make sure the transaction rolls back in the event of a panic.
	defer func() {
		if t.cellar != nil {
			_ = t.rollback()
		}
	}()

	// Mark as a managed tx so that the inner function cannot manually rollback.
	t.managed = true

	// If an error is returned from the function then pass it through.
	err = fn(t)
	t.managed = false
	if err != nil {
		_ = t.Rollback()
		return err
	}

	if err := t.Rollback(); err != nil {
		return err
	}

	return nil
}

// getRoot does atomic load of the current root
func (c *Cellar) getRoot(reason ...string) segmentList {
	c.rootLock.RLock()
	defer c.rootLock.RUnlock()
	return c.getRootLocked(reason...)
}

// getRoot does atomic load of the current root
func (c *Cellar) getRootLocked(reason ...string) segmentList {
	rv := c.root.Load().(segmentList)
	for _, segment := range rv {
		segment.incrRef(reason...)
	}
	return rv
}

func (c *Cellar) pushRoot(seg *segment) error {
	// we need to hold the lock the entire time
	// because we don't want to ensure our update
	// reflects the current root
	c.rootLock.Lock()
	defer c.rootLock.Unlock()

	// check to see if cellar is closed
	if c.master == nil {
		return ErrTxClosed
	}

	// get current root
	croot := c.getRootLocked("cellar pushRoot")

	// make new root
	nroot := make(segmentList, len(croot)+1)
	// place new segment in front
	nroot[0] = seg
	// copy remaining semgents
	copy(nroot[1:], croot)

	nrootbytes, err := nroot.MarshalBinary()
	if err != nil {
		return err
	}

	// persist the new root to our master database
	err = c.master.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(masterBucketName)
		if err != nil {
			return err
		}
		err = bucket.Put(rootKeyName, nrootbytes)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// bump the ref count for the new segment
	// this ensures a segment on the root, always has at least 1 ref
	seg.incrRef("on root")

	// now update the live root
	c.root.Store(nroot)

	// decr the refs
	for _, segment := range croot {
		segment.decrRef("cellar pushRoot")
	}

	// increment segment count
	atomic.AddUint64(&c.stats.segments, 1)

	// notify the merge manager of the new root
	// - we use get root to ensure we incr the refs
	root := c.getRootLocked("cellar push root notify merge")
	c.mergeManager.RootChange(root)

	return nil
}

func (c *Cellar) replaceSegments(replace segmentList, newseg *segment) error {
	// we need to hold the lock the entire time
	// because we don't want to ensure our update
	// reflects the current root
	c.rootLock.Lock()
	defer c.rootLock.Unlock()

	// check to see if cellar is closed
	if c.master == nil {
		return ErrTxClosed
	}

	// get current root
	croot := c.getRootLocked("cellar replaceSegments")

	// make new root
	nroot := make(segmentList, 0)
	// iterate through current segments
	replacei := 0
	for _, segment := range croot {
		if replacei < len(replace) && segment.seq == replace[replacei].seq {
			// this is a segment being replaced
			if replacei == 0 {
				// this is the first segment being replaced
				nroot = append(nroot, newseg)
			}
			replacei++
		} else {
			// this is not a segment being replaced, copy it over
			nroot = append(nroot, segment)
		}
	}

	nrootbytes, err := nroot.MarshalBinary()
	if err != nil {
		return err
	}

	// persist the new root to our master database
	err = c.master.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(masterBucketName)
		if err != nil {
			return err
		}
		err = bucket.Put(rootKeyName, nrootbytes)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	// bump the ref count for the new segment
	// this ensures a segment on the root, always has at least 1 ref
	newseg.incrRef("on root")

	// now update the live root
	c.root.Store(nroot)

	// decr the refs
	for _, segment := range croot {
		segment.decrRef("cellar replaceSegments")
	}

	newSegmentCount := uint64(len(croot) - len(replace) + 1)
	// increment segment count
	atomic.StoreUint64(&c.stats.segments, newSegmentCount)
	atomic.AddUint64(&c.stats.mergesCompleted, 1)

	// notify the merge manager of the new root
	// - we use get root to ensure we incr the refs
	root := c.getRootLocked("cellar replace segments notify merge")
	c.mergeManager.RootChange(root)

	for _, s := range replace {
		// decr the count now that it is off the root (balances incr from pushRoot)
		s.decrRef("off root")

		Logger.Printf("about to close seq: %d, refcount: %d", s.seq, atomic.LoadUint64(&s.refs))
		// asynchronously cleanup
		go func(seg *segment) {
			segmentPath := seg.DB.Path()
			Logger.Printf("closing: %v, roots: %v", segmentPath, nroot)
			err := seg.Close()
			if err != nil {
				// FIXME anything better to do?
				Logger.Printf("err closing segment %d: %v", seg.seq, err)
			}
			err = os.RemoveAll(segmentPath)
			if err != nil {
				// FIXME anything better to do?
				Logger.Printf("err removing segment %d: %v", seg.seq, err)
			}
		}(s)
	}

	return nil
}

// ForceMerge will force the cellar to perform a merge operations
// this function does not wait for the merge to finish
func (c *Cellar) ForceMerge() {
	root := c.getRoot("cellar forceMerge")
	c.mergeManager.ForceMerge(root)
}

// Stats returns a structure containing interesting metrics about the cellar
func (c *Cellar) Stats() *Stats {
	rv := &Stats{}
	rv.mergesCompleted = atomic.LoadUint64(&c.stats.mergesCompleted)
	rv.segments = atomic.LoadUint64(&c.stats.segments)
	return rv
}

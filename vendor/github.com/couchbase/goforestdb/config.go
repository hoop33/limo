package forestdb

//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

//#include <libforestdb/forestdb.h>
import "C"

import (
	"unsafe"
)

type OpenFlags uint32

const (
	OPEN_FLAG_CREATE OpenFlags = 1
	OPEN_FLAG_RDONLY OpenFlags = 2
)

type SeqTreeOpt uint8

const (
	SEQTREE_NOT_USE SeqTreeOpt = 0
	SEQTREE_USE     SeqTreeOpt = 1
)

type DurabilityOpt uint8

const (
	DRB_NONE          DurabilityOpt = 0
	DRB_ODIRECT       DurabilityOpt = 0x1
	DRB_ASYNC         DurabilityOpt = 0x2
	DRB_ODIRECT_ASYNC DurabilityOpt = 0x3
)

type CompactOpt uint8

const (
	COMPACT_MANUAL CompactOpt = 0
	COMPACT_AUTO   CompactOpt = 1
)

// ForestDB config options
type Config struct {
	config *C.fdb_config
}

func (c *Config) ChunkSize() uint16 {
	return uint16(c.config.chunksize)
}

func (c *Config) SetChunkSize(s uint16) {
	c.config.chunksize = C.uint16_t(s)
}

func (c *Config) BlockSize() uint32 {
	return uint32(c.config.blocksize)
}

func (c *Config) SetBlockSize(s uint32) {
	c.config.blocksize = C.uint32_t(s)
}

func (c *Config) BufferCacheSize() uint64 {
	return uint64(c.config.buffercache_size)
}

func (c *Config) SetBufferCacheSize(s uint64) {
	c.config.buffercache_size = C.uint64_t(s)
}

func (c *Config) WalThreshold() uint64 {
	return uint64(c.config.wal_threshold)
}

func (c *Config) SetWalThreshold(s uint64) {
	c.config.wal_threshold = C.uint64_t(s)
}

func (c *Config) WalFlushBeforeCommit() bool {
	return bool(c.config.wal_flush_before_commit)
}

func (c *Config) SetWalFlushBeforeCommit(b bool) {
	c.config.wal_flush_before_commit = C.bool(b)
}

func (c *Config) AutoCommit() bool {
	return bool(c.config.auto_commit)
}

func (c *Config) SetAutoCommit(b bool) {
	c.config.auto_commit = C.bool(b)
}

func (c *Config) PurgingInterval() uint32 {
	return uint32(c.config.purging_interval)
}

func (c *Config) SetPurgingInterval(s uint32) {
	c.config.purging_interval = C.uint32_t(s)
}

func (c *Config) SeqTreeOpt() SeqTreeOpt {
	return SeqTreeOpt(c.config.seqtree_opt)
}

func (c *Config) SetSeqTreeOpt(o SeqTreeOpt) {
	c.config.seqtree_opt = C.fdb_seqtree_opt_t(o)
}

func (c *Config) DurabilityOpt() DurabilityOpt {
	return DurabilityOpt(c.config.durability_opt)
}

func (c *Config) SetDurabilityOpt(o DurabilityOpt) {
	c.config.durability_opt = C.fdb_durability_opt_t(o)
}

func (c *Config) OpenFlags() OpenFlags {
	return OpenFlags(c.config.flags)
}

func (c *Config) SetOpenFlags(o OpenFlags) {
	c.config.flags = C.fdb_open_flags(o)
}

func (c *Config) CompactionBufferSizeMax() uint32 {
	return uint32(c.config.compaction_buf_maxsize)
}

func (c *Config) SetCompactionBufferSizeMax(s uint32) {
	c.config.compaction_buf_maxsize = C.uint32_t(s)
}

func (c *Config) CleanupCacheOnClose() bool {
	return bool(c.config.cleanup_cache_onclose)
}

func (c *Config) SetCleanupCacheOnClose(b bool) {
	c.config.cleanup_cache_onclose = C.bool(b)
}

func (c *Config) CompressDocumentBody() bool {
	return bool(c.config.compress_document_body)
}

func (c *Config) SetCompressDocumentBody(b bool) {
	c.config.compress_document_body = C.bool(b)
}

func (c *Config) CompactionMode() CompactOpt {
	return CompactOpt(c.config.compaction_mode)
}

func (c *Config) SetCompactionMode(o CompactOpt) {
	c.config.compaction_mode = C.fdb_compaction_mode_t(o)
}

func (c *Config) CompactionThreshold() uint8 {
	return uint8(c.config.compaction_threshold)
}

func (c *Config) SetCompactionThreshold(s uint8) {
	c.config.compaction_threshold = C.uint8_t(s)
}

func (c *Config) CompactionMinimumFilesize() uint64 {
	return uint64(c.config.compaction_minimum_filesize)
}

func (c *Config) SetCompactionMinimumFilesize(s uint64) {
	c.config.compaction_minimum_filesize = C.uint64_t(s)
}

func (c *Config) CompactorSleepDuration() uint64 {
	return uint64(c.config.compactor_sleep_duration)
}

func (c *Config) SetCompactorSleepDuration(s uint64) {
	c.config.compactor_sleep_duration = C.uint64_t(s)
}

func (c *Config) MultiKVInstances() bool {
	return bool(c.config.multi_kv_instances)
}

func (c *Config) SetMultiKVInstances(multi bool) {
	c.config.multi_kv_instances = C.bool(multi)
}

func (c *Config) PrefetchDuration() uint64 {
	return uint64(c.config.prefetch_duration)
}

func (c *Config) SetPrefetchDuration(s uint64) {
	c.config.prefetch_duration = C.uint64_t(s)
}

func (c *Config) NumWalPartitions() uint16 {
	return uint16(c.config.num_wal_partitions)
}

func (c *Config) SetNumWalPartitions(s uint16) {
	c.config.num_wal_partitions = C.uint16_t(s)
}

func (c *Config) NumBcachePartitions() uint16 {
	return uint16(c.config.num_bcache_partitions)
}

func (c *Config) SetNumBcachePartitions(s uint16) {
	c.config.num_bcache_partitions = C.uint16_t(s)
}

// TODO: compaction_cb.
// TODO: compaction_cb_mask.
// TODO: compaction_cb_ctx.

func (c *Config) MaxWriterLockProb() uint8 {
	return uint8(c.config.max_writer_lock_prob)
}

func (c *Config) SetMaxWriterLockProb(s uint8) {
	c.config.max_writer_lock_prob = C.size_t(s)
}

func (c *Config) NumCompactorThreads() int {
	return int(c.config.num_compactor_threads)
}

func (c *Config) SetNumCompactorThreads(s int) {
	c.config.num_compactor_threads = C.size_t(s)
}

func (c *Config) NumBgflusherThreads() int {
	return int(c.config.num_bgflusher_threads)
}

func (c *Config) SetNumBgflusherThreads(s int) {
	c.config.num_bgflusher_threads = C.size_t(s)
}

// TODO: encryption_key

func (c *Config) NumBlockReusingThreshold() int {
	return int(c.config.block_reusing_threshold)
}

func (c *Config) SetNumBlockReusingThreshold(s int) {
	c.config.block_reusing_threshold = C.size_t(s)
}

func (c *Config) NumKeepingHeaders() int {
	return int(c.config.num_keeping_headers)
}

func (c *Config) SetNumKeepingHeaders(s int) {
	c.config.num_keeping_headers = C.size_t(s)
}

// DefaultConfig gets the default ForestDB config
func DefaultConfig() *Config {
	Log.Debugf("fdb_get_default_config call")
	config := C.fdb_get_default_config()
	Log.Debugf("fdb_get_default_config ret config:%v", config)
	return &Config{
		config: &config,
	}
}

// ForestDB KVStore config options
type KVStoreConfig struct {
	config *C.fdb_kvs_config
}

func (c *KVStoreConfig) CreateIfMissing() bool {
	return bool(c.config.create_if_missing)
}

func (c *KVStoreConfig) SetCreateIfMissing(b bool) {
	c.config.create_if_missing = C.bool(b)
}

func (c *KVStoreConfig) SetCustomCompare(comparator unsafe.Pointer) {
	c.config.custom_cmp = C.fdb_custom_cmp_variable(comparator)
}

// DefaultConfig gets the default ForestDB config
func DefaultKVStoreConfig() *KVStoreConfig {
	Log.Debugf("fdb_get_default_kvs_config call")
	config := C.fdb_get_default_kvs_config()
	Log.Debugf("fdb_get_default_kvs_config ret config:%v", config)
	return &KVStoreConfig{
		config: &config,
	}
}

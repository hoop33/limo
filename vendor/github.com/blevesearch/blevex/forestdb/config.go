package forestdb

import (
	"fmt"

	"github.com/couchbase/goforestdb"
)

func applyConfig(c *forestdb.Config, config map[string]interface{}) (
	*forestdb.Config, error) {

	if v, exists := config["autoCommit"].(bool); exists {
		c.SetAutoCommit(v)
	}
	if v, exists := config["blockSize"].(float64); exists {
		c.SetBlockSize(uint32(v))
	}
	if v, exists := config["bufferCacheSize"].(float64); exists {
		c.SetBufferCacheSize(uint64(v))
	}
	if v, exists := config["chunkSize"].(float64); exists {
		c.SetChunkSize(uint16(v))
	}
	if v, exists := config["cleanupCacheOnClose"].(bool); exists {
		c.SetCleanupCacheOnClose(v)
	}
	if v, exists := config["compactionBufferSizeMax"].(float64); exists {
		c.SetCompactionBufferSizeMax(uint32(v))
	}
	if v, exists := config["compactionMinimumFilesize"].(float64); exists {
		c.SetCompactionMinimumFilesize(uint64(v))
	}
	if v, exists := config["compactionMode"].(string); exists {
		switch v {
		case "manual":
			c.SetCompactionMode(forestdb.COMPACT_MANUAL)
		case "auto":
			c.SetCompactionMode(forestdb.COMPACT_AUTO)
		default:
			return nil, fmt.Errorf("Unknown compaction mode: %s", v)
		}

	}
	if v, exists := config["compactionThreshold"].(float64); exists {
		c.SetCompactionThreshold(uint8(v))
	}
	if v, exists := config["compactorSleepDuration"].(float64); exists {
		c.SetCompactorSleepDuration(uint64(v))
	}
	if v, exists := config["compressDocumentBody"].(bool); exists {
		c.SetCompressDocumentBody(v)
	}
	if v, exists := config["multiKVInstances"].(bool); exists {
		c.SetMultiKVInstances(v)
	}
	if v, exists := config["prefetchDuration"].(float64); exists {
		c.SetPrefetchDuration(uint64(v))
	}
	if v, exists := config["numWalPartitions"].(float64); exists {
		c.SetNumWalPartitions(uint16(v))
	}
	if v, exists := config["numBcachePartitions"].(float64); exists {
		c.SetNumBcachePartitions(uint16(v))
	}
	if v, exists := config["durabilityOpt"].(string); exists {
		switch v {
		case "none":
			c.SetDurabilityOpt(forestdb.DRB_NONE)
		case "odirect":
			c.SetDurabilityOpt(forestdb.DRB_ODIRECT)
		case "async":
			c.SetDurabilityOpt(forestdb.DRB_ASYNC)
		case "async_odirect":
			c.SetDurabilityOpt(forestdb.DRB_ODIRECT_ASYNC)
		default:
			return nil, fmt.Errorf("Unknown durability option: %s", v)
		}

	}
	if v, exists := config["openFlags"].(string); exists {
		switch v {
		case "create":
			c.SetOpenFlags(forestdb.OPEN_FLAG_CREATE)
		case "readonly":
			c.SetOpenFlags(forestdb.OPEN_FLAG_RDONLY)
		default:
			return nil, fmt.Errorf("Unknown open flag: %s", v)
		}
	}
	if v, exists := config["purgingInterval"].(float64); exists {
		c.SetPurgingInterval(uint32(v))
	}
	if v, exists := config["seqTreeOpt"].(bool); exists {
		if !v {
			c.SetSeqTreeOpt(forestdb.SEQTREE_NOT_USE)
		}
	}
	if v, exists := config["walFlushBeforeCommit"].(bool); exists {
		c.SetWalFlushBeforeCommit(v)
	}
	if v, exists := config["walThreshold"].(float64); exists {
		c.SetWalThreshold(uint64(v))
	}
	if v, exists := config["maxWriterLockProb"].(float64); exists {
		c.SetMaxWriterLockProb(uint8(v))
	}
	if v, exists := config["numCompactorThreads"].(float64); exists {
		c.SetNumCompactorThreads(int(v))
	}
	if v, exists := config["numBgflusherThreads"].(float64); exists {
		c.SetNumBgflusherThreads(int(v))
	}
	if v, exists := config["numBlockReusingThreshold"].(float64); exists {
		c.SetNumBlockReusingThreshold(int(v))
	}
	if v, exists := config["numKeepingHeaders"].(float64); exists {
		c.SetNumKeepingHeaders(int(v))
	}
	return c, nil
}

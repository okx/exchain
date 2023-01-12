//go:build rocksdb
// +build rocksdb

package db

import (
	"bufio"
	"strconv"
	"strings"

	"github.com/cosmos/gorocksdb"
	"github.com/prometheus/client_golang/prometheus"
)

type RocksDBMetrics struct {
	opts *gorocksdb.Options

	blockCacheMiss                        *prometheus.GaugeVec
	blockCacheHit                         *prometheus.GaugeVec
	blockCacheAdd                         *prometheus.GaugeVec
	blockCacheAddFailures                 *prometheus.GaugeVec
	blockCacheIndexMiss                   *prometheus.GaugeVec
	blockCacheIndexHit                    *prometheus.GaugeVec
	blockCacheIndexAdd                    *prometheus.GaugeVec
	blockCacheIndexBytesInsert            *prometheus.GaugeVec
	blockCacheIndexBytesEvict             *prometheus.GaugeVec
	blockCacheFilterMiss                  *prometheus.GaugeVec
	blockCacheFilterHit                   *prometheus.GaugeVec
	blockCacheFilterAdd                   *prometheus.GaugeVec
	blockCacheFilterBytesInsert           *prometheus.GaugeVec
	blockCacheFilterBytesEvict            *prometheus.GaugeVec
	blockCacheDataMiss                    *prometheus.GaugeVec
	blockCacheDataHit                     *prometheus.GaugeVec
	blockCacheDataAdd                     *prometheus.GaugeVec
	blockCacheDataBytesInsert             *prometheus.GaugeVec
	blockCacheBytesRead                   *prometheus.GaugeVec
	blockCacheBytesWrite                  *prometheus.GaugeVec
	bloomFilterUseful                     *prometheus.GaugeVec
	bloomFilterFullPositive               *prometheus.GaugeVec
	bloomFilterFullTruePositive           *prometheus.GaugeVec
	bloomFilterMicros                     *prometheus.GaugeVec
	persistentCacheHit                    *prometheus.GaugeVec
	persistentCacheMiss                   *prometheus.GaugeVec
	simBlockCacheHit                      *prometheus.GaugeVec
	simBlockCacheMiss                     *prometheus.GaugeVec
	memtableHit                           *prometheus.GaugeVec
	memtableMiss                          *prometheus.GaugeVec
	l0Hit                                 *prometheus.GaugeVec
	l1Hit                                 *prometheus.GaugeVec
	l2andupHit                            *prometheus.GaugeVec
	compactionKeyDropNew                  *prometheus.GaugeVec
	compactionKeyDropObsolete             *prometheus.GaugeVec
	compactionKeyDropRangeDel             *prometheus.GaugeVec
	compactionKeyDropUser                 *prometheus.GaugeVec
	compactionRangeDelDropObsolete        *prometheus.GaugeVec
	compactionOptimizedDelDropObsolete    *prometheus.GaugeVec
	compactionCancelled                   *prometheus.GaugeVec
	numberKeysWritten                     *prometheus.GaugeVec
	numberKeysRead                        *prometheus.GaugeVec
	numberKeysUpdated                     *prometheus.GaugeVec
	bytesWritten                          *prometheus.GaugeVec
	bytesRead                             *prometheus.GaugeVec
	numberDbSeek                          *prometheus.GaugeVec
	numberDbNext                          *prometheus.GaugeVec
	numberDbPrev                          *prometheus.GaugeVec
	numberDbSeekFound                     *prometheus.GaugeVec
	numberDbNextFound                     *prometheus.GaugeVec
	numberDbPrevFound                     *prometheus.GaugeVec
	dbIterBytesRead                       *prometheus.GaugeVec
	noFileCloses                          *prometheus.GaugeVec
	noFileOpens                           *prometheus.GaugeVec
	noFileErrors                          *prometheus.GaugeVec
	l0SlowdownMicros                      *prometheus.GaugeVec
	memtableCompactionMicros              *prometheus.GaugeVec
	l0NumFilesStallMicros                 *prometheus.GaugeVec
	stallMicros                           *prometheus.GaugeVec
	dbMutexWaitMicros                     *prometheus.GaugeVec
	rateLimitDelayMillis                  *prometheus.GaugeVec
	numIterators                          *prometheus.GaugeVec
	numberMultigetGet                     *prometheus.GaugeVec
	numberMultigetKeysRead                *prometheus.GaugeVec
	numberMultigetBytesRead               *prometheus.GaugeVec
	numberDeletesFiltered                 *prometheus.GaugeVec
	numberMergeFailures                   *prometheus.GaugeVec
	bloomFilterPrefixChecked              *prometheus.GaugeVec
	bloomFilterPrefixUseful               *prometheus.GaugeVec
	numberReseeksIteration                *prometheus.GaugeVec
	getupdatessinceCalls                  *prometheus.GaugeVec
	blockCachecompressedMiss              *prometheus.GaugeVec
	blockCachecompressedHit               *prometheus.GaugeVec
	blockCachecompressedAdd               *prometheus.GaugeVec
	blockCachecompressedAddFailures       *prometheus.GaugeVec
	walSynced                             *prometheus.GaugeVec
	walBytes                              *prometheus.GaugeVec
	writeSelf                             *prometheus.GaugeVec
	writeOther                            *prometheus.GaugeVec
	writeTimeout                          *prometheus.GaugeVec
	writeWal                              *prometheus.GaugeVec
	compactReadBytes                      *prometheus.GaugeVec
	compactWriteBytes                     *prometheus.GaugeVec
	flushWriteBytes                       *prometheus.GaugeVec
	compactReadMarkedBytes                *prometheus.GaugeVec
	compactReadPeriodicBytes              *prometheus.GaugeVec
	compactReadTtlBytes                   *prometheus.GaugeVec
	compactWriteMarkedBytes               *prometheus.GaugeVec
	compactWritePeriodicBytes             *prometheus.GaugeVec
	compactWriteTtlBytes                  *prometheus.GaugeVec
	numberDirectLoadTableProperties       *prometheus.GaugeVec
	numberSuperversionAcquires            *prometheus.GaugeVec
	numberSuperversionReleases            *prometheus.GaugeVec
	numberSuperversionCleanups            *prometheus.GaugeVec
	numberBlockCompressed                 *prometheus.GaugeVec
	numberBlockDecompressed               *prometheus.GaugeVec
	numberBlockNotCompressed              *prometheus.GaugeVec
	mergeOperationTimeNanos               *prometheus.GaugeVec
	filterOperationTimeNanos              *prometheus.GaugeVec
	rowCacheHit                           *prometheus.GaugeVec
	rowCacheMiss                          *prometheus.GaugeVec
	readAmpEstimateUsefulBytes            *prometheus.GaugeVec
	readAmpTotalReadBytes                 *prometheus.GaugeVec
	numberRateLimiterDrains               *prometheus.GaugeVec
	numberIterSkip                        *prometheus.GaugeVec
	blobdbNumPut                          *prometheus.GaugeVec
	blobdbNumWrite                        *prometheus.GaugeVec
	blobdbNumGet                          *prometheus.GaugeVec
	blobdbNumMultiget                     *prometheus.GaugeVec
	blobdbNumSeek                         *prometheus.GaugeVec
	blobdbNumNext                         *prometheus.GaugeVec
	blobdbNumPrev                         *prometheus.GaugeVec
	blobdbNumKeysWritten                  *prometheus.GaugeVec
	blobdbNumKeysRead                     *prometheus.GaugeVec
	blobdbBytesWritten                    *prometheus.GaugeVec
	blobdbBytesRead                       *prometheus.GaugeVec
	blobdbWriteInlined                    *prometheus.GaugeVec
	blobdbWriteInlinedTtl                 *prometheus.GaugeVec
	blobdbWriteBlob                       *prometheus.GaugeVec
	blobdbWriteBlobTtl                    *prometheus.GaugeVec
	blobdbBlobFileBytesWritten            *prometheus.GaugeVec
	blobdbBlobFileBytesRead               *prometheus.GaugeVec
	blobdbBlobFileSynced                  *prometheus.GaugeVec
	blobdbBlobIndexExpiredCount           *prometheus.GaugeVec
	blobdbBlobIndexExpiredSize            *prometheus.GaugeVec
	blobdbBlobIndexEvictedCount           *prometheus.GaugeVec
	blobdbBlobIndexEvictedSize            *prometheus.GaugeVec
	blobdbGcNumFiles                      *prometheus.GaugeVec
	blobdbGcNumNewFiles                   *prometheus.GaugeVec
	blobdbGcFailures                      *prometheus.GaugeVec
	blobdbGcNumKeysOverwritten            *prometheus.GaugeVec
	blobdbGcNumKeysExpired                *prometheus.GaugeVec
	blobdbGcNumKeysRelocated              *prometheus.GaugeVec
	blobdbGcBytesOverwritten              *prometheus.GaugeVec
	blobdbGcBytesExpired                  *prometheus.GaugeVec
	blobdbGcBytesRelocated                *prometheus.GaugeVec
	blobdbFifoNumFilesEvicted             *prometheus.GaugeVec
	blobdbFifoNumKeysEvicted              *prometheus.GaugeVec
	blobdbFifoBytesEvicted                *prometheus.GaugeVec
	txnOverheadMutexPrepare               *prometheus.GaugeVec
	txnOverheadMutexOldCommitMap          *prometheus.GaugeVec
	txnOverheadDuplicateKey               *prometheus.GaugeVec
	txnOverheadMutexSnapshot              *prometheus.GaugeVec
	txnGetTryagain                        *prometheus.GaugeVec
	numberMultigetKeysFound               *prometheus.GaugeVec
	numIteratorCreated                    *prometheus.GaugeVec
	numIteratorDeleted                    *prometheus.GaugeVec
	blockCacheCompressionDictMiss         *prometheus.GaugeVec
	blockCacheCompressionDictHit          *prometheus.GaugeVec
	blockCacheCompressionDictAdd          *prometheus.GaugeVec
	blockCacheCompressionDictBytesInsert  *prometheus.GaugeVec
	blockCacheCompressionDictBytesEvict   *prometheus.GaugeVec
	blockCacheAddRedundant                *prometheus.GaugeVec
	blockCacheIndexAddRedundant           *prometheus.GaugeVec
	blockCacheFilterAddRedundant          *prometheus.GaugeVec
	blockCacheDataAddRedundant            *prometheus.GaugeVec
	blockCacheCompressionDictAddRedundant *prometheus.GaugeVec
	filesMarkedTrash                      *prometheus.GaugeVec
	filesDeletedImmediately               *prometheus.GaugeVec
	errorHandlerBgErrroCount              *prometheus.GaugeVec
	errorHandlerBgIoErrroCount            *prometheus.GaugeVec
	errorHandlerBgRetryableIoErrroCount   *prometheus.GaugeVec
	errorHandlerAutoresumeCount           *prometheus.GaugeVec
	errorHandlerAutoresumeRetryTotalCount *prometheus.GaugeVec
	errorHandlerAutoresumeSuccessCount    *prometheus.GaugeVec
	memtablePayloadBytesAtFlush           *prometheus.GaugeVec
	memtableGarbageBytesAtFlush           *prometheus.GaugeVec
	secondaryCacheHits                    *prometheus.GaugeVec
	verifyChecksumReadBytes               *prometheus.GaugeVec
	backupReadBytes                       *prometheus.GaugeVec
	backupWriteBytes                      *prometheus.GaugeVec
	remoteCompactReadBytes                *prometheus.GaugeVec
	remoteCompactWriteBytes               *prometheus.GaugeVec
	hotFileReadBytes                      *prometheus.GaugeVec
	warmFileReadBytes                     *prometheus.GaugeVec
	coldFileReadBytes                     *prometheus.GaugeVec
	hotFileReadCount                      *prometheus.GaugeVec
	warmFileReadCount                     *prometheus.GaugeVec
	coldFileReadCount                     *prometheus.GaugeVec

	dbGetMicros                         *prometheus.Desc
	dbWriteMicros                       *prometheus.Desc
	compactionTimesMicros               *prometheus.Desc
	compactionTimesCpuMicros            *prometheus.Desc
	subcompactionSetupTimesMicros       *prometheus.Desc
	tableSyncMicros                     *prometheus.Desc
	compactionOutfileSyncMicros         *prometheus.Desc
	walFileSyncMicros                   *prometheus.Desc
	manifestFileSyncMicros              *prometheus.Desc
	tableOpenIoMicros                   *prometheus.Desc
	dbMultigetMicros                    *prometheus.Desc
	readBlockCompactionMicros           *prometheus.Desc
	readBlockGetMicros                  *prometheus.Desc
	writeRawBlockMicros                 *prometheus.Desc
	l0SlowdownCount                     *prometheus.Desc
	memtableCompactionCount             *prometheus.Desc
	numFilesStallCount                  *prometheus.Desc
	hardRateLimitDelayCount             *prometheus.Desc
	softRateLimitDelayCount             *prometheus.Desc
	numfilesInSinglecompaction          *prometheus.Desc
	dbSeekMicros                        *prometheus.Desc
	dbWriteStall                        *prometheus.Desc
	sstReadMicros                       *prometheus.Desc
	numSubcompactionsScheduled          *prometheus.Desc
	bytesPerRead                        *prometheus.Desc
	bytesPerWrite                       *prometheus.Desc
	bytesPerMultiget                    *prometheus.Desc
	bytesCompressed                     *prometheus.Desc
	bytesDecompressed                   *prometheus.Desc
	compressionTimesNanos               *prometheus.Desc
	decompressionTimesNanos             *prometheus.Desc
	readNumMergeOperands                *prometheus.Desc
	blobdbKeySize                       *prometheus.Desc
	blobdbValueSize                     *prometheus.Desc
	blobdbWriteMicros                   *prometheus.Desc
	blobdbGetMicros                     *prometheus.Desc
	blobdbMultigetMicros                *prometheus.Desc
	blobdbSeekMicros                    *prometheus.Desc
	blobdbNextMicros                    *prometheus.Desc
	blobdbPrevMicros                    *prometheus.Desc
	blobdbBlobFileWriteMicros           *prometheus.Desc
	blobdbBlobFileReadMicros            *prometheus.Desc
	blobdbBlobFileSyncMicros            *prometheus.Desc
	blobdbGcMicros                      *prometheus.Desc
	blobdbCompressionMicros             *prometheus.Desc
	blobdbDecompressionMicros           *prometheus.Desc
	dbFlushMicros                       *prometheus.Desc
	sstBatchSize                        *prometheus.Desc
	numIndexAndFilterBlocksReadPerLevel *prometheus.Desc
	numDataBlocksReadPerLevel           *prometheus.Desc
	numSstReadPerLevel                  *prometheus.Desc
	errorHandlerAutoresumeRetryCount    *prometheus.Desc
}

func NewRocksDBMetrics(opts *gorocksdb.Options) *RocksDBMetrics {
	return &RocksDBMetrics{
		opts:                                  opts,
		blockCacheMiss:                        prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_miss"}, nil),
		blockCacheHit:                         prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_hit"}, nil),
		blockCacheAdd:                         prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_add"}, nil),
		blockCacheAddFailures:                 prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_add_failures"}, nil),
		blockCacheIndexMiss:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_index_miss"}, nil),
		blockCacheIndexHit:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_index_hit"}, nil),
		blockCacheIndexAdd:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_index_add"}, nil),
		blockCacheIndexBytesInsert:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_index_bytes_insert"}, nil),
		blockCacheIndexBytesEvict:             prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_index_bytes_evict"}, nil),
		blockCacheFilterMiss:                  prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_filter_miss"}, nil),
		blockCacheFilterHit:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_filter_hit"}, nil),
		blockCacheFilterAdd:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_filter_add"}, nil),
		blockCacheFilterBytesInsert:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_filter_bytes_insert"}, nil),
		blockCacheFilterBytesEvict:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_filter_bytes_evict"}, nil),
		blockCacheDataMiss:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_data_miss"}, nil),
		blockCacheDataHit:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_data_hit"}, nil),
		blockCacheDataAdd:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_data_add"}, nil),
		blockCacheDataBytesInsert:             prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_data_bytes_insert"}, nil),
		blockCacheBytesRead:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_bytes_read"}, nil),
		blockCacheBytesWrite:                  prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_bytes_write"}, nil),
		bloomFilterUseful:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "bloom_filter_useful"}, nil),
		bloomFilterFullPositive:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "bloom_filter_full_positive"}, nil),
		bloomFilterFullTruePositive:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "bloom_filter_full_true_positive"}, nil),
		bloomFilterMicros:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "bloom_filter_micros"}, nil),
		persistentCacheHit:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "persistent_cache_hit"}, nil),
		persistentCacheMiss:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "persistent_cache_miss"}, nil),
		simBlockCacheHit:                      prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "sim_block_cache_hit"}, nil),
		simBlockCacheMiss:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "sim_block_cache_miss"}, nil),
		memtableHit:                           prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "memtable_hit"}, nil),
		memtableMiss:                          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "memtable_miss"}, nil),
		l0Hit:                                 prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "l0_hit"}, nil),
		l1Hit:                                 prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "l1_hit"}, nil),
		l2andupHit:                            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "l2andup_hit"}, nil),
		compactionKeyDropNew:                  prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compaction_key_drop_new"}, nil),
		compactionKeyDropObsolete:             prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compaction_key_drop_obsolete"}, nil),
		compactionKeyDropRangeDel:             prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compaction_key_drop_range_del"}, nil),
		compactionKeyDropUser:                 prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compaction_key_drop_user"}, nil),
		compactionRangeDelDropObsolete:        prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compaction_range_del_drop_obsolete"}, nil),
		compactionOptimizedDelDropObsolete:    prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compaction_optimized_del_drop_obsolete"}, nil),
		compactionCancelled:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compaction_cancelled"}, nil),
		numberKeysWritten:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_keys_written"}, nil),
		numberKeysRead:                        prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_keys_read"}, nil),
		numberKeysUpdated:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_keys_updated"}, nil),
		bytesWritten:                          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "bytes_written"}, nil),
		bytesRead:                             prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "bytes_read"}, nil),
		numberDbSeek:                          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_db_seek"}, nil),
		numberDbNext:                          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_db_next"}, nil),
		numberDbPrev:                          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_db_prev"}, nil),
		numberDbSeekFound:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_db_seek_found"}, nil),
		numberDbNextFound:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_db_next_found"}, nil),
		numberDbPrevFound:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_db_prev_found"}, nil),
		dbIterBytesRead:                       prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "db_iter_bytes_read"}, nil),
		noFileCloses:                          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "no_file_closes"}, nil),
		noFileOpens:                           prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "no_file_opens"}, nil),
		noFileErrors:                          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "no_file_errors"}, nil),
		l0SlowdownMicros:                      prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "l0_slowdown_micros"}, nil),
		memtableCompactionMicros:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "memtable_compaction_micros"}, nil),
		l0NumFilesStallMicros:                 prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "l0_num_files_stall_micros"}, nil),
		stallMicros:                           prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "stall_micros"}, nil),
		dbMutexWaitMicros:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "db_mutex_wait_micros"}, nil),
		rateLimitDelayMillis:                  prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "rate_limit_delay_millis"}, nil),
		numIterators:                          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "num_iterators"}, nil),
		numberMultigetGet:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_multiget_get"}, nil),
		numberMultigetKeysRead:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_multiget_keys_read"}, nil),
		numberMultigetBytesRead:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_multiget_bytes_read"}, nil),
		numberDeletesFiltered:                 prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_deletes_filtered"}, nil),
		numberMergeFailures:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_merge_failures"}, nil),
		bloomFilterPrefixChecked:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "bloom_filter_prefix_checked"}, nil),
		bloomFilterPrefixUseful:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "bloom_filter_prefix_useful"}, nil),
		numberReseeksIteration:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_reseeks_iteration"}, nil),
		getupdatessinceCalls:                  prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "getupdatessince_calls"}, nil),
		blockCachecompressedMiss:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cachecompressed_miss"}, nil),
		blockCachecompressedHit:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cachecompressed_hit"}, nil),
		blockCachecompressedAdd:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cachecompressed_add"}, nil),
		blockCachecompressedAddFailures:       prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cachecompressed_add_failures"}, nil),
		walSynced:                             prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "wal_synced"}, nil),
		walBytes:                              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "wal_bytes"}, nil),
		writeSelf:                             prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "write_self"}, nil),
		writeOther:                            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "write_other"}, nil),
		writeTimeout:                          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "write_timeout"}, nil),
		writeWal:                              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "write_wal"}, nil),
		compactReadBytes:                      prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compact_read_bytes"}, nil),
		compactWriteBytes:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compact_write_bytes"}, nil),
		flushWriteBytes:                       prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "flush_write_bytes"}, nil),
		compactReadMarkedBytes:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compact_read_marked_bytes"}, nil),
		compactReadPeriodicBytes:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compact_read_periodic_bytes"}, nil),
		compactReadTtlBytes:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compact_read_ttl_bytes"}, nil),
		compactWriteMarkedBytes:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compact_write_marked_bytes"}, nil),
		compactWritePeriodicBytes:             prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compact_write_periodic_bytes"}, nil),
		compactWriteTtlBytes:                  prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "compact_write_ttl_bytes"}, nil),
		numberDirectLoadTableProperties:       prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_direct_load_table_properties"}, nil),
		numberSuperversionAcquires:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_superversion_acquires"}, nil),
		numberSuperversionReleases:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_superversion_releases"}, nil),
		numberSuperversionCleanups:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_superversion_cleanups"}, nil),
		numberBlockCompressed:                 prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_block_compressed"}, nil),
		numberBlockDecompressed:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_block_decompressed"}, nil),
		numberBlockNotCompressed:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_block_not_compressed"}, nil),
		mergeOperationTimeNanos:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "merge_operation_time_nanos"}, nil),
		filterOperationTimeNanos:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "filter_operation_time_nanos"}, nil),
		rowCacheHit:                           prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "row_cache_hit"}, nil),
		rowCacheMiss:                          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "row_cache_miss"}, nil),
		readAmpEstimateUsefulBytes:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "read_amp_estimate_useful_bytes"}, nil),
		readAmpTotalReadBytes:                 prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "read_amp_total_read_bytes"}, nil),
		numberRateLimiterDrains:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_rate_limiter_drains"}, nil),
		numberIterSkip:                        prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_iter_skip"}, nil),
		blobdbNumPut:                          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_num_put"}, nil),
		blobdbNumWrite:                        prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_num_write"}, nil),
		blobdbNumGet:                          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_num_get"}, nil),
		blobdbNumMultiget:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_num_multiget"}, nil),
		blobdbNumSeek:                         prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_num_seek"}, nil),
		blobdbNumNext:                         prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_num_next"}, nil),
		blobdbNumPrev:                         prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_num_prev"}, nil),
		blobdbNumKeysWritten:                  prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_num_keys_written"}, nil),
		blobdbNumKeysRead:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_num_keys_read"}, nil),
		blobdbBytesWritten:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_bytes_written"}, nil),
		blobdbBytesRead:                       prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_bytes_read"}, nil),
		blobdbWriteInlined:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_write_inlined"}, nil),
		blobdbWriteInlinedTtl:                 prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_write_inlined_ttl"}, nil),
		blobdbWriteBlob:                       prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_write_blob"}, nil),
		blobdbWriteBlobTtl:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_write_blob_ttl"}, nil),
		blobdbBlobFileBytesWritten:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_blob_file_bytes_written"}, nil),
		blobdbBlobFileBytesRead:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_blob_file_bytes_read"}, nil),
		blobdbBlobFileSynced:                  prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_blob_file_synced"}, nil),
		blobdbBlobIndexExpiredCount:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_blob_index_expired_count"}, nil),
		blobdbBlobIndexExpiredSize:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_blob_index_expired_size"}, nil),
		blobdbBlobIndexEvictedCount:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_blob_index_evicted_count"}, nil),
		blobdbBlobIndexEvictedSize:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_blob_index_evicted_size"}, nil),
		blobdbGcNumFiles:                      prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_gc_num_files"}, nil),
		blobdbGcNumNewFiles:                   prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_gc_num_new_files"}, nil),
		blobdbGcFailures:                      prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_gc_failures"}, nil),
		blobdbGcNumKeysOverwritten:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_gc_num_keys_overwritten"}, nil),
		blobdbGcNumKeysExpired:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_gc_num_keys_expired"}, nil),
		blobdbGcNumKeysRelocated:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_gc_num_keys_relocated"}, nil),
		blobdbGcBytesOverwritten:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_gc_bytes_overwritten"}, nil),
		blobdbGcBytesExpired:                  prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_gc_bytes_expired"}, nil),
		blobdbGcBytesRelocated:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_gc_bytes_relocated"}, nil),
		blobdbFifoNumFilesEvicted:             prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_fifo_num_files_evicted"}, nil),
		blobdbFifoNumKeysEvicted:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_fifo_num_keys_evicted"}, nil),
		blobdbFifoBytesEvicted:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "blobdb_fifo_bytes_evicted"}, nil),
		txnOverheadMutexPrepare:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "txn_overhead_mutex_prepare"}, nil),
		txnOverheadMutexOldCommitMap:          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "txn_overhead_mutex_old_commit_map"}, nil),
		txnOverheadDuplicateKey:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "txn_overhead_duplicate_key"}, nil),
		txnOverheadMutexSnapshot:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "txn_overhead_mutex_snapshot"}, nil),
		txnGetTryagain:                        prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "txn_get_tryagain"}, nil),
		numberMultigetKeysFound:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "number_multiget_keys_found"}, nil),
		numIteratorCreated:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "num_iterator_created"}, nil),
		numIteratorDeleted:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "num_iterator_deleted"}, nil),
		blockCacheCompressionDictMiss:         prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_compression_dict_miss"}, nil),
		blockCacheCompressionDictHit:          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_compression_dict_hit"}, nil),
		blockCacheCompressionDictAdd:          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_compression_dict_add"}, nil),
		blockCacheCompressionDictBytesInsert:  prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_compression_dict_bytes_insert"}, nil),
		blockCacheCompressionDictBytesEvict:   prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_compression_dict_bytes_evict"}, nil),
		blockCacheAddRedundant:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_add_redundant"}, nil),
		blockCacheIndexAddRedundant:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_index_add_redundant"}, nil),
		blockCacheFilterAddRedundant:          prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_filter_add_redundant"}, nil),
		blockCacheDataAddRedundant:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_data_add_redundant"}, nil),
		blockCacheCompressionDictAddRedundant: prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "block_cache_compression_dict_add_redundant"}, nil),
		filesMarkedTrash:                      prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "files_marked_trash"}, nil),
		filesDeletedImmediately:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "files_deleted_immediately"}, nil),
		errorHandlerBgErrroCount:              prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "error_handler_bg_errro_count"}, nil),
		errorHandlerBgIoErrroCount:            prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "error_handler_bg_io_errro_count"}, nil),
		errorHandlerBgRetryableIoErrroCount:   prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "error_handler_bg_retryable_io_errro_count"}, nil),
		errorHandlerAutoresumeCount:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "error_handler_autoresume_count"}, nil),
		errorHandlerAutoresumeRetryTotalCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "error_handler_autoresume_retry_total_count"}, nil),
		errorHandlerAutoresumeSuccessCount:    prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "error_handler_autoresume_success_count"}, nil),
		memtablePayloadBytesAtFlush:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "memtable_payload_bytes_at_flush"}, nil),
		memtableGarbageBytesAtFlush:           prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "memtable_garbage_bytes_at_flush"}, nil),
		secondaryCacheHits:                    prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "secondary_cache_hits"}, nil),
		verifyChecksumReadBytes:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "verify_checksum_read_bytes"}, nil),
		backupReadBytes:                       prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "backup_read_bytes"}, nil),
		backupWriteBytes:                      prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "backup_write_bytes"}, nil),
		remoteCompactReadBytes:                prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "remote_compact_read_bytes"}, nil),
		remoteCompactWriteBytes:               prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "remote_compact_write_bytes"}, nil),
		hotFileReadBytes:                      prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "hot_file_read_bytes"}, nil),
		warmFileReadBytes:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "warm_file_read_bytes"}, nil),
		coldFileReadBytes:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "cold_file_read_bytes"}, nil),
		hotFileReadCount:                      prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "hot_file_read_count"}, nil),
		warmFileReadCount:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "warm_file_read_count"}, nil),
		coldFileReadCount:                     prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: "rocksdb", Name: "cold_file_read_count"}, nil),

		dbGetMicros:                         prometheus.NewDesc("rocksdb_db_get_micros", "", nil, nil),
		dbWriteMicros:                       prometheus.NewDesc("rocksdb_db_write_micros", "", nil, nil),
		compactionTimesMicros:               prometheus.NewDesc("rocksdb_compaction_times_micros", "", nil, nil),
		compactionTimesCpuMicros:            prometheus.NewDesc("rocksdb_compaction_times_cpu_micros", "", nil, nil),
		subcompactionSetupTimesMicros:       prometheus.NewDesc("rocksdb_subcompaction_setup_times_micros", "", nil, nil),
		tableSyncMicros:                     prometheus.NewDesc("rocksdb_table_sync_micros", "", nil, nil),
		compactionOutfileSyncMicros:         prometheus.NewDesc("rocksdb_compaction_outfile_sync_micros", "", nil, nil),
		walFileSyncMicros:                   prometheus.NewDesc("rocksdb_wal_file_sync_micros", "", nil, nil),
		manifestFileSyncMicros:              prometheus.NewDesc("rocksdb_manifest_file_sync_micros", "", nil, nil),
		tableOpenIoMicros:                   prometheus.NewDesc("rocksdb_table_open_io_micros", "", nil, nil),
		dbMultigetMicros:                    prometheus.NewDesc("rocksdb_db_multiget_micros", "", nil, nil),
		readBlockCompactionMicros:           prometheus.NewDesc("rocksdb_read_block_compaction_micros", "", nil, nil),
		readBlockGetMicros:                  prometheus.NewDesc("rocksdb_read_block_get_micros", "", nil, nil),
		writeRawBlockMicros:                 prometheus.NewDesc("rocksdb_write_raw_block_micros", "", nil, nil),
		l0SlowdownCount:                     prometheus.NewDesc("rocksdb_l0_slowdown_count", "", nil, nil),
		memtableCompactionCount:             prometheus.NewDesc("rocksdb_memtable_compaction_count", "", nil, nil),
		numFilesStallCount:                  prometheus.NewDesc("rocksdb_num_files_stall_count", "", nil, nil),
		hardRateLimitDelayCount:             prometheus.NewDesc("rocksdb_hard_rate_limit_delay_count", "", nil, nil),
		softRateLimitDelayCount:             prometheus.NewDesc("rocksdb_soft_rate_limit_delay_count", "", nil, nil),
		numfilesInSinglecompaction:          prometheus.NewDesc("rocksdb_numfiles_in_singlecompaction", "", nil, nil),
		dbSeekMicros:                        prometheus.NewDesc("rocksdb_db_seek_micros", "", nil, nil),
		dbWriteStall:                        prometheus.NewDesc("rocksdb_db_write_stall", "", nil, nil),
		sstReadMicros:                       prometheus.NewDesc("rocksdb_sst_read_micros", "", nil, nil),
		numSubcompactionsScheduled:          prometheus.NewDesc("rocksdb_num_subcompactions_scheduled", "", nil, nil),
		bytesPerRead:                        prometheus.NewDesc("rocksdb_bytes_per_read", "", nil, nil),
		bytesPerWrite:                       prometheus.NewDesc("rocksdb_bytes_per_write", "", nil, nil),
		bytesPerMultiget:                    prometheus.NewDesc("rocksdb_bytes_per_multiget", "", nil, nil),
		bytesCompressed:                     prometheus.NewDesc("rocksdb_bytes_compressed", "", nil, nil),
		bytesDecompressed:                   prometheus.NewDesc("rocksdb_bytes_decompressed", "", nil, nil),
		compressionTimesNanos:               prometheus.NewDesc("rocksdb_compression_times_nanos", "", nil, nil),
		decompressionTimesNanos:             prometheus.NewDesc("rocksdb_decompression_times_nanos", "", nil, nil),
		readNumMergeOperands:                prometheus.NewDesc("rocksdb_read_num_merge_operands", "", nil, nil),
		blobdbKeySize:                       prometheus.NewDesc("rocksdb_blobdb_key_size", "", nil, nil),
		blobdbValueSize:                     prometheus.NewDesc("rocksdb_blobdb_value_size", "", nil, nil),
		blobdbWriteMicros:                   prometheus.NewDesc("rocksdb_blobdb_write_micros", "", nil, nil),
		blobdbGetMicros:                     prometheus.NewDesc("rocksdb_blobdb_get_micros", "", nil, nil),
		blobdbMultigetMicros:                prometheus.NewDesc("rocksdb_blobdb_multiget_micros", "", nil, nil),
		blobdbSeekMicros:                    prometheus.NewDesc("rocksdb_blobdb_seek_micros", "", nil, nil),
		blobdbNextMicros:                    prometheus.NewDesc("rocksdb_blobdb_next_micros", "", nil, nil),
		blobdbPrevMicros:                    prometheus.NewDesc("rocksdb_blobdb_prev_micros", "", nil, nil),
		blobdbBlobFileWriteMicros:           prometheus.NewDesc("rocksdb_blobdb_blob_file_write_micros", "", nil, nil),
		blobdbBlobFileReadMicros:            prometheus.NewDesc("rocksdb_blobdb_blob_file_read_micros", "", nil, nil),
		blobdbBlobFileSyncMicros:            prometheus.NewDesc("rocksdb_blobdb_blob_file_sync_micros", "", nil, nil),
		blobdbGcMicros:                      prometheus.NewDesc("rocksdb_blobdb_gc_micros", "", nil, nil),
		blobdbCompressionMicros:             prometheus.NewDesc("rocksdb_blobdb_compression_micros", "", nil, nil),
		blobdbDecompressionMicros:           prometheus.NewDesc("rocksdb_blobdb_decompression_micros", "", nil, nil),
		dbFlushMicros:                       prometheus.NewDesc("rocksdb_db_flush_micros", "", nil, nil),
		sstBatchSize:                        prometheus.NewDesc("rocksdb_sst_batch_size", "", nil, nil),
		numIndexAndFilterBlocksReadPerLevel: prometheus.NewDesc("rocksdb_num_index_and_filter_blocks_read_per_level", "", nil, nil),
		numDataBlocksReadPerLevel:           prometheus.NewDesc("rocksdb_num_data_blocks_read_per_level", "", nil, nil),
		numSstReadPerLevel:                  prometheus.NewDesc("rocksdb_num_sst_read_per_level", "", nil, nil),
		errorHandlerAutoresumeRetryCount:    prometheus.NewDesc("rocksdb_error_handler_autoresume_retry_count", "", nil, nil),
	}
}

func (dbm *RocksDBMetrics) Describe(dc chan<- *prometheus.Desc) {
	dbm.blockCacheMiss.Describe(dc)
	dbm.blockCacheHit.Describe(dc)
	dbm.blockCacheAdd.Describe(dc)
	dbm.blockCacheAddFailures.Describe(dc)
	dbm.blockCacheIndexMiss.Describe(dc)
	dbm.blockCacheIndexHit.Describe(dc)
	dbm.blockCacheIndexAdd.Describe(dc)
	dbm.blockCacheIndexBytesInsert.Describe(dc)
	dbm.blockCacheIndexBytesEvict.Describe(dc)
	dbm.blockCacheFilterMiss.Describe(dc)
	dbm.blockCacheFilterHit.Describe(dc)
	dbm.blockCacheFilterAdd.Describe(dc)
	dbm.blockCacheFilterBytesInsert.Describe(dc)
	dbm.blockCacheFilterBytesEvict.Describe(dc)
	dbm.blockCacheDataMiss.Describe(dc)
	dbm.blockCacheDataHit.Describe(dc)
	dbm.blockCacheDataAdd.Describe(dc)
	dbm.blockCacheDataBytesInsert.Describe(dc)
	dbm.blockCacheBytesRead.Describe(dc)
	dbm.blockCacheBytesWrite.Describe(dc)
	dbm.bloomFilterUseful.Describe(dc)
	dbm.bloomFilterFullPositive.Describe(dc)
	dbm.bloomFilterFullTruePositive.Describe(dc)
	dbm.bloomFilterMicros.Describe(dc)
	dbm.persistentCacheHit.Describe(dc)
	dbm.persistentCacheMiss.Describe(dc)
	dbm.simBlockCacheHit.Describe(dc)
	dbm.simBlockCacheMiss.Describe(dc)
	dbm.memtableHit.Describe(dc)
	dbm.memtableMiss.Describe(dc)
	dbm.l0Hit.Describe(dc)
	dbm.l1Hit.Describe(dc)
	dbm.l2andupHit.Describe(dc)
	dbm.compactionKeyDropNew.Describe(dc)
	dbm.compactionKeyDropObsolete.Describe(dc)
	dbm.compactionKeyDropRangeDel.Describe(dc)
	dbm.compactionKeyDropUser.Describe(dc)
	dbm.compactionRangeDelDropObsolete.Describe(dc)
	dbm.compactionOptimizedDelDropObsolete.Describe(dc)
	dbm.compactionCancelled.Describe(dc)
	dbm.numberKeysWritten.Describe(dc)
	dbm.numberKeysRead.Describe(dc)
	dbm.numberKeysUpdated.Describe(dc)
	dbm.bytesWritten.Describe(dc)
	dbm.bytesRead.Describe(dc)
	dbm.numberDbSeek.Describe(dc)
	dbm.numberDbNext.Describe(dc)
	dbm.numberDbPrev.Describe(dc)
	dbm.numberDbSeekFound.Describe(dc)
	dbm.numberDbNextFound.Describe(dc)
	dbm.numberDbPrevFound.Describe(dc)
	dbm.dbIterBytesRead.Describe(dc)
	dbm.noFileCloses.Describe(dc)
	dbm.noFileOpens.Describe(dc)
	dbm.noFileErrors.Describe(dc)
	dbm.l0SlowdownMicros.Describe(dc)
	dbm.memtableCompactionMicros.Describe(dc)
	dbm.l0NumFilesStallMicros.Describe(dc)
	dbm.stallMicros.Describe(dc)
	dbm.dbMutexWaitMicros.Describe(dc)
	dbm.rateLimitDelayMillis.Describe(dc)
	dbm.numIterators.Describe(dc)
	dbm.numberMultigetGet.Describe(dc)
	dbm.numberMultigetKeysRead.Describe(dc)
	dbm.numberMultigetBytesRead.Describe(dc)
	dbm.numberDeletesFiltered.Describe(dc)
	dbm.numberMergeFailures.Describe(dc)
	dbm.bloomFilterPrefixChecked.Describe(dc)
	dbm.bloomFilterPrefixUseful.Describe(dc)
	dbm.numberReseeksIteration.Describe(dc)
	dbm.getupdatessinceCalls.Describe(dc)
	dbm.blockCachecompressedMiss.Describe(dc)
	dbm.blockCachecompressedHit.Describe(dc)
	dbm.blockCachecompressedAdd.Describe(dc)
	dbm.blockCachecompressedAddFailures.Describe(dc)
	dbm.walSynced.Describe(dc)
	dbm.walBytes.Describe(dc)
	dbm.writeSelf.Describe(dc)
	dbm.writeOther.Describe(dc)
	dbm.writeTimeout.Describe(dc)
	dbm.writeWal.Describe(dc)
	dbm.compactReadBytes.Describe(dc)
	dbm.compactWriteBytes.Describe(dc)
	dbm.flushWriteBytes.Describe(dc)
	dbm.compactReadMarkedBytes.Describe(dc)
	dbm.compactReadPeriodicBytes.Describe(dc)
	dbm.compactReadTtlBytes.Describe(dc)
	dbm.compactWriteMarkedBytes.Describe(dc)
	dbm.compactWritePeriodicBytes.Describe(dc)
	dbm.compactWriteTtlBytes.Describe(dc)
	dbm.numberDirectLoadTableProperties.Describe(dc)
	dbm.numberSuperversionAcquires.Describe(dc)
	dbm.numberSuperversionReleases.Describe(dc)
	dbm.numberSuperversionCleanups.Describe(dc)
	dbm.numberBlockCompressed.Describe(dc)
	dbm.numberBlockDecompressed.Describe(dc)
	dbm.numberBlockNotCompressed.Describe(dc)
	dbm.mergeOperationTimeNanos.Describe(dc)
	dbm.filterOperationTimeNanos.Describe(dc)
	dbm.rowCacheHit.Describe(dc)
	dbm.rowCacheMiss.Describe(dc)
	dbm.readAmpEstimateUsefulBytes.Describe(dc)
	dbm.readAmpTotalReadBytes.Describe(dc)
	dbm.numberRateLimiterDrains.Describe(dc)
	dbm.numberIterSkip.Describe(dc)
	dbm.blobdbNumPut.Describe(dc)
	dbm.blobdbNumWrite.Describe(dc)
	dbm.blobdbNumGet.Describe(dc)
	dbm.blobdbNumMultiget.Describe(dc)
	dbm.blobdbNumSeek.Describe(dc)
	dbm.blobdbNumNext.Describe(dc)
	dbm.blobdbNumPrev.Describe(dc)
	dbm.blobdbNumKeysWritten.Describe(dc)
	dbm.blobdbNumKeysRead.Describe(dc)
	dbm.blobdbBytesWritten.Describe(dc)
	dbm.blobdbBytesRead.Describe(dc)
	dbm.blobdbWriteInlined.Describe(dc)
	dbm.blobdbWriteInlinedTtl.Describe(dc)
	dbm.blobdbWriteBlob.Describe(dc)
	dbm.blobdbWriteBlobTtl.Describe(dc)
	dbm.blobdbBlobFileBytesWritten.Describe(dc)
	dbm.blobdbBlobFileBytesRead.Describe(dc)
	dbm.blobdbBlobFileSynced.Describe(dc)
	dbm.blobdbBlobIndexExpiredCount.Describe(dc)
	dbm.blobdbBlobIndexExpiredSize.Describe(dc)
	dbm.blobdbBlobIndexEvictedCount.Describe(dc)
	dbm.blobdbBlobIndexEvictedSize.Describe(dc)
	dbm.blobdbGcNumFiles.Describe(dc)
	dbm.blobdbGcNumNewFiles.Describe(dc)
	dbm.blobdbGcFailures.Describe(dc)
	dbm.blobdbGcNumKeysOverwritten.Describe(dc)
	dbm.blobdbGcNumKeysExpired.Describe(dc)
	dbm.blobdbGcNumKeysRelocated.Describe(dc)
	dbm.blobdbGcBytesOverwritten.Describe(dc)
	dbm.blobdbGcBytesExpired.Describe(dc)
	dbm.blobdbGcBytesRelocated.Describe(dc)
	dbm.blobdbFifoNumFilesEvicted.Describe(dc)
	dbm.blobdbFifoNumKeysEvicted.Describe(dc)
	dbm.blobdbFifoBytesEvicted.Describe(dc)
	dbm.txnOverheadMutexPrepare.Describe(dc)
	dbm.txnOverheadMutexOldCommitMap.Describe(dc)
	dbm.txnOverheadDuplicateKey.Describe(dc)
	dbm.txnOverheadMutexSnapshot.Describe(dc)
	dbm.txnGetTryagain.Describe(dc)
	dbm.numberMultigetKeysFound.Describe(dc)
	dbm.numIteratorCreated.Describe(dc)
	dbm.numIteratorDeleted.Describe(dc)
	dbm.blockCacheCompressionDictMiss.Describe(dc)
	dbm.blockCacheCompressionDictHit.Describe(dc)
	dbm.blockCacheCompressionDictAdd.Describe(dc)
	dbm.blockCacheCompressionDictBytesInsert.Describe(dc)
	dbm.blockCacheCompressionDictBytesEvict.Describe(dc)
	dbm.blockCacheAddRedundant.Describe(dc)
	dbm.blockCacheIndexAddRedundant.Describe(dc)
	dbm.blockCacheFilterAddRedundant.Describe(dc)
	dbm.blockCacheDataAddRedundant.Describe(dc)
	dbm.blockCacheCompressionDictAddRedundant.Describe(dc)
	dbm.filesMarkedTrash.Describe(dc)
	dbm.filesDeletedImmediately.Describe(dc)
	dbm.errorHandlerBgErrroCount.Describe(dc)
	dbm.errorHandlerBgIoErrroCount.Describe(dc)
	dbm.errorHandlerBgRetryableIoErrroCount.Describe(dc)
	dbm.errorHandlerAutoresumeCount.Describe(dc)
	dbm.errorHandlerAutoresumeRetryTotalCount.Describe(dc)
	dbm.errorHandlerAutoresumeSuccessCount.Describe(dc)
	dbm.memtablePayloadBytesAtFlush.Describe(dc)
	dbm.memtableGarbageBytesAtFlush.Describe(dc)
	dbm.secondaryCacheHits.Describe(dc)
	dbm.verifyChecksumReadBytes.Describe(dc)
	dbm.backupReadBytes.Describe(dc)
	dbm.backupWriteBytes.Describe(dc)
	dbm.remoteCompactReadBytes.Describe(dc)
	dbm.remoteCompactWriteBytes.Describe(dc)
	dbm.hotFileReadBytes.Describe(dc)
	dbm.warmFileReadBytes.Describe(dc)
	dbm.coldFileReadBytes.Describe(dc)
	dbm.hotFileReadCount.Describe(dc)
	dbm.warmFileReadCount.Describe(dc)
	dbm.coldFileReadCount.Describe(dc)

	dc <- dbm.dbGetMicros
	dc <- dbm.dbWriteMicros
	dc <- dbm.compactionTimesMicros
	dc <- dbm.compactionTimesCpuMicros
	dc <- dbm.subcompactionSetupTimesMicros
	dc <- dbm.tableSyncMicros
	dc <- dbm.compactionOutfileSyncMicros
	dc <- dbm.walFileSyncMicros
	dc <- dbm.manifestFileSyncMicros
	dc <- dbm.tableOpenIoMicros
	dc <- dbm.dbMultigetMicros
	dc <- dbm.readBlockCompactionMicros
	dc <- dbm.readBlockGetMicros
	dc <- dbm.writeRawBlockMicros
	dc <- dbm.l0SlowdownCount
	dc <- dbm.memtableCompactionCount
	dc <- dbm.numFilesStallCount
	dc <- dbm.hardRateLimitDelayCount
	dc <- dbm.softRateLimitDelayCount
	dc <- dbm.numfilesInSinglecompaction
	dc <- dbm.dbSeekMicros
	dc <- dbm.dbWriteStall
	dc <- dbm.sstReadMicros
	dc <- dbm.numSubcompactionsScheduled
	dc <- dbm.bytesPerRead
	dc <- dbm.bytesPerWrite
	dc <- dbm.bytesPerMultiget
	dc <- dbm.bytesCompressed
	dc <- dbm.bytesDecompressed
	dc <- dbm.compressionTimesNanos
	dc <- dbm.decompressionTimesNanos
	dc <- dbm.readNumMergeOperands
	dc <- dbm.blobdbKeySize
	dc <- dbm.blobdbValueSize
	dc <- dbm.blobdbWriteMicros
	dc <- dbm.blobdbGetMicros
	dc <- dbm.blobdbMultigetMicros
	dc <- dbm.blobdbSeekMicros
	dc <- dbm.blobdbNextMicros
	dc <- dbm.blobdbPrevMicros
	dc <- dbm.blobdbBlobFileWriteMicros
	dc <- dbm.blobdbBlobFileReadMicros
	dc <- dbm.blobdbBlobFileSyncMicros
	dc <- dbm.blobdbGcMicros
	dc <- dbm.blobdbCompressionMicros
	dc <- dbm.blobdbDecompressionMicros
	dc <- dbm.dbFlushMicros
	dc <- dbm.sstBatchSize
	dc <- dbm.numIndexAndFilterBlocksReadPerLevel
	dc <- dbm.numDataBlocksReadPerLevel
	dc <- dbm.numSstReadPerLevel
	dc <- dbm.errorHandlerAutoresumeRetryCount
}

func (dbm *RocksDBMetrics) Collect(mc chan<- prometheus.Metric) {
	stats := dbm.opts.GetStatisticsString()

	scanner := bufio.NewScanner(strings.NewReader(stats))
	for scanner.Scan() {
		fs := strings.Fields(scanner.Text())
		switch fs[0] {
		case "rocksdb.block.cache.miss":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheMiss.WithLabelValues().Set(v)
			dbm.blockCacheMiss.Collect(mc)
		case "rocksdb.block.cache.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheHit.WithLabelValues().Set(v)
			dbm.blockCacheHit.Collect(mc)
		case "rocksdb.block.cache.add":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheAdd.WithLabelValues().Set(v)
			dbm.blockCacheAdd.Collect(mc)
		case "rocksdb.block.cache.add.failures":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheAddFailures.WithLabelValues().Set(v)
			dbm.blockCacheAddFailures.Collect(mc)
		case "rocksdb.block.cache.index.miss":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheIndexMiss.WithLabelValues().Set(v)
			dbm.blockCacheIndexMiss.Collect(mc)
		case "rocksdb.block.cache.index.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheIndexHit.WithLabelValues().Set(v)
			dbm.blockCacheIndexHit.Collect(mc)
		case "rocksdb.block.cache.index.add":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheIndexAdd.WithLabelValues().Set(v)
			dbm.blockCacheIndexAdd.Collect(mc)
		case "rocksdb.block.cache.index.bytes.insert":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheIndexBytesInsert.WithLabelValues().Set(v)
			dbm.blockCacheIndexBytesInsert.Collect(mc)
		case "rocksdb.block.cache.index.bytes.evict":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheIndexBytesEvict.WithLabelValues().Set(v)
			dbm.blockCacheIndexBytesEvict.Collect(mc)
		case "rocksdb.block.cache.filter.miss":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheFilterMiss.WithLabelValues().Set(v)
			dbm.blockCacheFilterMiss.Collect(mc)
		case "rocksdb.block.cache.filter.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheFilterHit.WithLabelValues().Set(v)
			dbm.blockCacheFilterHit.Collect(mc)
		case "rocksdb.block.cache.filter.add":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheFilterAdd.WithLabelValues().Set(v)
			dbm.blockCacheFilterAdd.Collect(mc)
		case "rocksdb.block.cache.filter.bytes.insert":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheFilterBytesInsert.WithLabelValues().Set(v)
			dbm.blockCacheFilterBytesInsert.Collect(mc)
		case "rocksdb.block.cache.filter.bytes.evict":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheFilterBytesEvict.WithLabelValues().Set(v)
			dbm.blockCacheFilterBytesEvict.Collect(mc)
		case "rocksdb.block.cache.data.miss":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheDataMiss.WithLabelValues().Set(v)
			dbm.blockCacheDataMiss.Collect(mc)
		case "rocksdb.block.cache.data.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheDataHit.WithLabelValues().Set(v)
			dbm.blockCacheDataHit.Collect(mc)
		case "rocksdb.block.cache.data.add":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheDataAdd.WithLabelValues().Set(v)
			dbm.blockCacheDataAdd.Collect(mc)
		case "rocksdb.block.cache.data.bytes.insert":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheDataBytesInsert.WithLabelValues().Set(v)
			dbm.blockCacheDataBytesInsert.Collect(mc)
		case "rocksdb.block.cache.bytes.read":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheBytesRead.WithLabelValues().Set(v)
			dbm.blockCacheBytesRead.Collect(mc)
		case "rocksdb.block.cache.bytes.write":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheBytesWrite.WithLabelValues().Set(v)
			dbm.blockCacheBytesWrite.Collect(mc)
		case "rocksdb.bloom.filter.useful":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.bloomFilterUseful.WithLabelValues().Set(v)
			dbm.bloomFilterUseful.Collect(mc)
		case "rocksdb.bloom.filter.full.positive":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.bloomFilterFullPositive.WithLabelValues().Set(v)
			dbm.bloomFilterFullPositive.Collect(mc)
		case "rocksdb.bloom.filter.full.true.positive":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.bloomFilterFullTruePositive.WithLabelValues().Set(v)
			dbm.bloomFilterFullTruePositive.Collect(mc)
		case "rocksdb.bloom.filter.micros":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.bloomFilterMicros.WithLabelValues().Set(v)
			dbm.bloomFilterMicros.Collect(mc)
		case "rocksdb.persistent.cache.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.persistentCacheHit.WithLabelValues().Set(v)
			dbm.persistentCacheHit.Collect(mc)
		case "rocksdb.persistent.cache.miss":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.persistentCacheMiss.WithLabelValues().Set(v)
			dbm.persistentCacheMiss.Collect(mc)
		case "rocksdb.sim.block.cache.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.simBlockCacheHit.WithLabelValues().Set(v)
			dbm.simBlockCacheHit.Collect(mc)
		case "rocksdb.sim.block.cache.miss":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.simBlockCacheMiss.WithLabelValues().Set(v)
			dbm.simBlockCacheMiss.Collect(mc)
		case "rocksdb.memtable.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.memtableHit.WithLabelValues().Set(v)
			dbm.memtableHit.Collect(mc)
		case "rocksdb.memtable.miss":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.memtableMiss.WithLabelValues().Set(v)
			dbm.memtableMiss.Collect(mc)
		case "rocksdb.l0.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.l0Hit.WithLabelValues().Set(v)
			dbm.l0Hit.Collect(mc)
		case "rocksdb.l1.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.l1Hit.WithLabelValues().Set(v)
			dbm.l1Hit.Collect(mc)
		case "rocksdb.l2andup.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.l2andupHit.WithLabelValues().Set(v)
			dbm.l2andupHit.Collect(mc)
		case "rocksdb.compaction.key.drop.new":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactionKeyDropNew.WithLabelValues().Set(v)
			dbm.compactionKeyDropNew.Collect(mc)
		case "rocksdb.compaction.key.drop.obsolete":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactionKeyDropObsolete.WithLabelValues().Set(v)
			dbm.compactionKeyDropObsolete.Collect(mc)
		case "rocksdb.compaction.key.drop.range_del":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactionKeyDropRangeDel.WithLabelValues().Set(v)
			dbm.compactionKeyDropRangeDel.Collect(mc)
		case "rocksdb.compaction.key.drop.user":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactionKeyDropUser.WithLabelValues().Set(v)
			dbm.compactionKeyDropUser.Collect(mc)
		case "rocksdb.compaction.range_del.drop.obsolete":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactionRangeDelDropObsolete.WithLabelValues().Set(v)
			dbm.compactionRangeDelDropObsolete.Collect(mc)
		case "rocksdb.compaction.optimized.del.drop.obsolete":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactionOptimizedDelDropObsolete.WithLabelValues().Set(v)
			dbm.compactionOptimizedDelDropObsolete.Collect(mc)
		case "rocksdb.compaction.cancelled":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactionCancelled.WithLabelValues().Set(v)
			dbm.compactionCancelled.Collect(mc)
		case "rocksdb.number.keys.written":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberKeysWritten.WithLabelValues().Set(v)
			dbm.numberKeysWritten.Collect(mc)
		case "rocksdb.number.keys.read":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberKeysRead.WithLabelValues().Set(v)
			dbm.numberKeysRead.Collect(mc)
		case "rocksdb.number.keys.updated":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberKeysUpdated.WithLabelValues().Set(v)
			dbm.numberKeysUpdated.Collect(mc)
		case "rocksdb.bytes.written":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.bytesWritten.WithLabelValues().Set(v)
			dbm.bytesWritten.Collect(mc)
		case "rocksdb.bytes.read":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.bytesRead.WithLabelValues().Set(v)
			dbm.bytesRead.Collect(mc)
		case "rocksdb.number.db.seek":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberDbSeek.WithLabelValues().Set(v)
			dbm.numberDbSeek.Collect(mc)
		case "rocksdb.number.db.next":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberDbNext.WithLabelValues().Set(v)
			dbm.numberDbNext.Collect(mc)
		case "rocksdb.number.db.prev":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberDbPrev.WithLabelValues().Set(v)
			dbm.numberDbPrev.Collect(mc)
		case "rocksdb.number.db.seek.found":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberDbSeekFound.WithLabelValues().Set(v)
			dbm.numberDbSeekFound.Collect(mc)
		case "rocksdb.number.db.next.found":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberDbNextFound.WithLabelValues().Set(v)
			dbm.numberDbNextFound.Collect(mc)
		case "rocksdb.number.db.prev.found":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberDbPrevFound.WithLabelValues().Set(v)
			dbm.numberDbPrevFound.Collect(mc)
		case "rocksdb.db.iter.bytes.read":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.dbIterBytesRead.WithLabelValues().Set(v)
			dbm.dbIterBytesRead.Collect(mc)
		case "rocksdb.no.file.closes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.noFileCloses.WithLabelValues().Set(v)
			dbm.noFileCloses.Collect(mc)
		case "rocksdb.no.file.opens":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.noFileOpens.WithLabelValues().Set(v)
			dbm.noFileOpens.Collect(mc)
		case "rocksdb.no.file.errors":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.noFileErrors.WithLabelValues().Set(v)
			dbm.noFileErrors.Collect(mc)
		case "rocksdb.l0.slowdown.micros":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.l0SlowdownMicros.WithLabelValues().Set(v)
			dbm.l0SlowdownMicros.Collect(mc)
		case "rocksdb.memtable.compaction.micros":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.memtableCompactionMicros.WithLabelValues().Set(v)
			dbm.memtableCompactionMicros.Collect(mc)
		case "rocksdb.l0.num.files.stall.micros":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.l0NumFilesStallMicros.WithLabelValues().Set(v)
			dbm.l0NumFilesStallMicros.Collect(mc)
		case "rocksdb.stall.micros":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.stallMicros.WithLabelValues().Set(v)
			dbm.stallMicros.Collect(mc)
		case "rocksdb.db.mutex.wait.micros":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.dbMutexWaitMicros.WithLabelValues().Set(v)
			dbm.dbMutexWaitMicros.Collect(mc)
		case "rocksdb.rate.limit.delay.millis":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.rateLimitDelayMillis.WithLabelValues().Set(v)
			dbm.rateLimitDelayMillis.Collect(mc)
		case "rocksdb.num.iterators":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numIterators.WithLabelValues().Set(v)
			dbm.numIterators.Collect(mc)
		case "rocksdb.number.multiget.get":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberMultigetGet.WithLabelValues().Set(v)
			dbm.numberMultigetGet.Collect(mc)
		case "rocksdb.number.multiget.keys.read":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberMultigetKeysRead.WithLabelValues().Set(v)
			dbm.numberMultigetKeysRead.Collect(mc)
		case "rocksdb.number.multiget.bytes.read":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberMultigetBytesRead.WithLabelValues().Set(v)
			dbm.numberMultigetBytesRead.Collect(mc)
		case "rocksdb.number.deletes.filtered":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberDeletesFiltered.WithLabelValues().Set(v)
			dbm.numberDeletesFiltered.Collect(mc)
		case "rocksdb.number.merge.failures":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberMergeFailures.WithLabelValues().Set(v)
			dbm.numberMergeFailures.Collect(mc)
		case "rocksdb.bloom.filter.prefix.checked":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.bloomFilterPrefixChecked.WithLabelValues().Set(v)
			dbm.bloomFilterPrefixChecked.Collect(mc)
		case "rocksdb.bloom.filter.prefix.useful":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.bloomFilterPrefixUseful.WithLabelValues().Set(v)
			dbm.bloomFilterPrefixUseful.Collect(mc)
		case "rocksdb.number.reseeks.iteration":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberReseeksIteration.WithLabelValues().Set(v)
			dbm.numberReseeksIteration.Collect(mc)
		case "rocksdb.getupdatessince.calls":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.getupdatessinceCalls.WithLabelValues().Set(v)
			dbm.getupdatessinceCalls.Collect(mc)
		case "rocksdb.block.cachecompressed.miss":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCachecompressedMiss.WithLabelValues().Set(v)
			dbm.blockCachecompressedMiss.Collect(mc)
		case "rocksdb.block.cachecompressed.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCachecompressedHit.WithLabelValues().Set(v)
			dbm.blockCachecompressedHit.Collect(mc)
		case "rocksdb.block.cachecompressed.add":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCachecompressedAdd.WithLabelValues().Set(v)
			dbm.blockCachecompressedAdd.Collect(mc)
		case "rocksdb.block.cachecompressed.add.failures":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCachecompressedAddFailures.WithLabelValues().Set(v)
			dbm.blockCachecompressedAddFailures.Collect(mc)
		case "rocksdb.wal.synced":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.walSynced.WithLabelValues().Set(v)
			dbm.walSynced.Collect(mc)
		case "rocksdb.wal.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.walBytes.WithLabelValues().Set(v)
			dbm.walBytes.Collect(mc)
		case "rocksdb.write.self":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.writeSelf.WithLabelValues().Set(v)
			dbm.writeSelf.Collect(mc)
		case "rocksdb.write.other":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.writeOther.WithLabelValues().Set(v)
			dbm.writeOther.Collect(mc)
		case "rocksdb.write.timeout":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.writeTimeout.WithLabelValues().Set(v)
			dbm.writeTimeout.Collect(mc)
		case "rocksdb.write.wal":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.writeWal.WithLabelValues().Set(v)
			dbm.writeWal.Collect(mc)
		case "rocksdb.compact.read.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactReadBytes.WithLabelValues().Set(v)
			dbm.compactReadBytes.Collect(mc)
		case "rocksdb.compact.write.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactWriteBytes.WithLabelValues().Set(v)
			dbm.compactWriteBytes.Collect(mc)
		case "rocksdb.flush.write.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.flushWriteBytes.WithLabelValues().Set(v)
			dbm.flushWriteBytes.Collect(mc)
		case "rocksdb.compact.read.marked.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactReadMarkedBytes.WithLabelValues().Set(v)
			dbm.compactReadMarkedBytes.Collect(mc)
		case "rocksdb.compact.read.periodic.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactReadPeriodicBytes.WithLabelValues().Set(v)
			dbm.compactReadPeriodicBytes.Collect(mc)
		case "rocksdb.compact.read.ttl.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactReadTtlBytes.WithLabelValues().Set(v)
			dbm.compactReadTtlBytes.Collect(mc)
		case "rocksdb.compact.write.marked.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactWriteMarkedBytes.WithLabelValues().Set(v)
			dbm.compactWriteMarkedBytes.Collect(mc)
		case "rocksdb.compact.write.periodic.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactWritePeriodicBytes.WithLabelValues().Set(v)
			dbm.compactWritePeriodicBytes.Collect(mc)
		case "rocksdb.compact.write.ttl.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.compactWriteTtlBytes.WithLabelValues().Set(v)
			dbm.compactWriteTtlBytes.Collect(mc)
		case "rocksdb.number.direct.load.table.properties":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberDirectLoadTableProperties.WithLabelValues().Set(v)
			dbm.numberDirectLoadTableProperties.Collect(mc)
		case "rocksdb.number.superversion_acquires":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberSuperversionAcquires.WithLabelValues().Set(v)
			dbm.numberSuperversionAcquires.Collect(mc)
		case "rocksdb.number.superversion_releases":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberSuperversionReleases.WithLabelValues().Set(v)
			dbm.numberSuperversionReleases.Collect(mc)
		case "rocksdb.number.superversion_cleanups":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberSuperversionCleanups.WithLabelValues().Set(v)
			dbm.numberSuperversionCleanups.Collect(mc)
		case "rocksdb.number.block.compressed":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberBlockCompressed.WithLabelValues().Set(v)
			dbm.numberBlockCompressed.Collect(mc)
		case "rocksdb.number.block.decompressed":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberBlockDecompressed.WithLabelValues().Set(v)
			dbm.numberBlockDecompressed.Collect(mc)
		case "rocksdb.number.block.not_compressed":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberBlockNotCompressed.WithLabelValues().Set(v)
			dbm.numberBlockNotCompressed.Collect(mc)
		case "rocksdb.merge.operation.time.nanos":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.mergeOperationTimeNanos.WithLabelValues().Set(v)
			dbm.mergeOperationTimeNanos.Collect(mc)
		case "rocksdb.filter.operation.time.nanos":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.filterOperationTimeNanos.WithLabelValues().Set(v)
			dbm.filterOperationTimeNanos.Collect(mc)
		case "rocksdb.row.cache.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.rowCacheHit.WithLabelValues().Set(v)
			dbm.rowCacheHit.Collect(mc)
		case "rocksdb.row.cache.miss":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.rowCacheMiss.WithLabelValues().Set(v)
			dbm.rowCacheMiss.Collect(mc)
		case "rocksdb.read.amp.estimate.useful.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.readAmpEstimateUsefulBytes.WithLabelValues().Set(v)
			dbm.readAmpEstimateUsefulBytes.Collect(mc)
		case "rocksdb.read.amp.total.read.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.readAmpTotalReadBytes.WithLabelValues().Set(v)
			dbm.readAmpTotalReadBytes.Collect(mc)
		case "rocksdb.number.rate_limiter.drains":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberRateLimiterDrains.WithLabelValues().Set(v)
			dbm.numberRateLimiterDrains.Collect(mc)
		case "rocksdb.number.iter.skip":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberIterSkip.WithLabelValues().Set(v)
			dbm.numberIterSkip.Collect(mc)
		case "rocksdb.blobdb.num.put":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbNumPut.WithLabelValues().Set(v)
			dbm.blobdbNumPut.Collect(mc)
		case "rocksdb.blobdb.num.write":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbNumWrite.WithLabelValues().Set(v)
			dbm.blobdbNumWrite.Collect(mc)
		case "rocksdb.blobdb.num.get":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbNumGet.WithLabelValues().Set(v)
			dbm.blobdbNumGet.Collect(mc)
		case "rocksdb.blobdb.num.multiget":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbNumMultiget.WithLabelValues().Set(v)
			dbm.blobdbNumMultiget.Collect(mc)
		case "rocksdb.blobdb.num.seek":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbNumSeek.WithLabelValues().Set(v)
			dbm.blobdbNumSeek.Collect(mc)
		case "rocksdb.blobdb.num.next":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbNumNext.WithLabelValues().Set(v)
			dbm.blobdbNumNext.Collect(mc)
		case "rocksdb.blobdb.num.prev":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbNumPrev.WithLabelValues().Set(v)
			dbm.blobdbNumPrev.Collect(mc)
		case "rocksdb.blobdb.num.keys.written":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbNumKeysWritten.WithLabelValues().Set(v)
			dbm.blobdbNumKeysWritten.Collect(mc)
		case "rocksdb.blobdb.num.keys.read":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbNumKeysRead.WithLabelValues().Set(v)
			dbm.blobdbNumKeysRead.Collect(mc)
		case "rocksdb.blobdb.bytes.written":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbBytesWritten.WithLabelValues().Set(v)
			dbm.blobdbBytesWritten.Collect(mc)
		case "rocksdb.blobdb.bytes.read":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbBytesRead.WithLabelValues().Set(v)
			dbm.blobdbBytesRead.Collect(mc)
		case "rocksdb.blobdb.write.inlined":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbWriteInlined.WithLabelValues().Set(v)
			dbm.blobdbWriteInlined.Collect(mc)
		case "rocksdb.blobdb.write.inlined.ttl":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbWriteInlinedTtl.WithLabelValues().Set(v)
			dbm.blobdbWriteInlinedTtl.Collect(mc)
		case "rocksdb.blobdb.write.blob":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbWriteBlob.WithLabelValues().Set(v)
			dbm.blobdbWriteBlob.Collect(mc)
		case "rocksdb.blobdb.write.blob.ttl":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbWriteBlobTtl.WithLabelValues().Set(v)
			dbm.blobdbWriteBlobTtl.Collect(mc)
		case "rocksdb.blobdb.blob.file.bytes.written":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbBlobFileBytesWritten.WithLabelValues().Set(v)
			dbm.blobdbBlobFileBytesWritten.Collect(mc)
		case "rocksdb.blobdb.blob.file.bytes.read":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbBlobFileBytesRead.WithLabelValues().Set(v)
			dbm.blobdbBlobFileBytesRead.Collect(mc)
		case "rocksdb.blobdb.blob.file.synced":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbBlobFileSynced.WithLabelValues().Set(v)
			dbm.blobdbBlobFileSynced.Collect(mc)
		case "rocksdb.blobdb.blob.index.expired.count":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbBlobIndexExpiredCount.WithLabelValues().Set(v)
			dbm.blobdbBlobIndexExpiredCount.Collect(mc)
		case "rocksdb.blobdb.blob.index.expired.size":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbBlobIndexExpiredSize.WithLabelValues().Set(v)
			dbm.blobdbBlobIndexExpiredSize.Collect(mc)
		case "rocksdb.blobdb.blob.index.evicted.count":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbBlobIndexEvictedCount.WithLabelValues().Set(v)
			dbm.blobdbBlobIndexEvictedCount.Collect(mc)
		case "rocksdb.blobdb.blob.index.evicted.size":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbBlobIndexEvictedSize.WithLabelValues().Set(v)
			dbm.blobdbBlobIndexEvictedSize.Collect(mc)
		case "rocksdb.blobdb.gc.num.files":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbGcNumFiles.WithLabelValues().Set(v)
			dbm.blobdbGcNumFiles.Collect(mc)
		case "rocksdb.blobdb.gc.num.new.files":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbGcNumNewFiles.WithLabelValues().Set(v)
			dbm.blobdbGcNumNewFiles.Collect(mc)
		case "rocksdb.blobdb.gc.failures":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbGcFailures.WithLabelValues().Set(v)
			dbm.blobdbGcFailures.Collect(mc)
		case "rocksdb.blobdb.gc.num.keys.overwritten":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbGcNumKeysOverwritten.WithLabelValues().Set(v)
			dbm.blobdbGcNumKeysOverwritten.Collect(mc)
		case "rocksdb.blobdb.gc.num.keys.expired":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbGcNumKeysExpired.WithLabelValues().Set(v)
			dbm.blobdbGcNumKeysExpired.Collect(mc)
		case "rocksdb.blobdb.gc.num.keys.relocated":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbGcNumKeysRelocated.WithLabelValues().Set(v)
			dbm.blobdbGcNumKeysRelocated.Collect(mc)
		case "rocksdb.blobdb.gc.bytes.overwritten":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbGcBytesOverwritten.WithLabelValues().Set(v)
			dbm.blobdbGcBytesOverwritten.Collect(mc)
		case "rocksdb.blobdb.gc.bytes.expired":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbGcBytesExpired.WithLabelValues().Set(v)
			dbm.blobdbGcBytesExpired.Collect(mc)
		case "rocksdb.blobdb.gc.bytes.relocated":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbGcBytesRelocated.WithLabelValues().Set(v)
			dbm.blobdbGcBytesRelocated.Collect(mc)
		case "rocksdb.blobdb.fifo.num.files.evicted":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbFifoNumFilesEvicted.WithLabelValues().Set(v)
			dbm.blobdbFifoNumFilesEvicted.Collect(mc)
		case "rocksdb.blobdb.fifo.num.keys.evicted":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbFifoNumKeysEvicted.WithLabelValues().Set(v)
			dbm.blobdbFifoNumKeysEvicted.Collect(mc)
		case "rocksdb.blobdb.fifo.bytes.evicted":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blobdbFifoBytesEvicted.WithLabelValues().Set(v)
			dbm.blobdbFifoBytesEvicted.Collect(mc)
		case "rocksdb.txn.overhead.mutex.prepare":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.txnOverheadMutexPrepare.WithLabelValues().Set(v)
			dbm.txnOverheadMutexPrepare.Collect(mc)
		case "rocksdb.txn.overhead.mutex.old.commit.map":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.txnOverheadMutexOldCommitMap.WithLabelValues().Set(v)
			dbm.txnOverheadMutexOldCommitMap.Collect(mc)
		case "rocksdb.txn.overhead.duplicate.key":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.txnOverheadDuplicateKey.WithLabelValues().Set(v)
			dbm.txnOverheadDuplicateKey.Collect(mc)
		case "rocksdb.txn.overhead.mutex.snapshot":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.txnOverheadMutexSnapshot.WithLabelValues().Set(v)
			dbm.txnOverheadMutexSnapshot.Collect(mc)
		case "rocksdb.txn.get.tryagain":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.txnGetTryagain.WithLabelValues().Set(v)
			dbm.txnGetTryagain.Collect(mc)
		case "rocksdb.number.multiget.keys.found":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numberMultigetKeysFound.WithLabelValues().Set(v)
			dbm.numberMultigetKeysFound.Collect(mc)
		case "rocksdb.num.iterator.created":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numIteratorCreated.WithLabelValues().Set(v)
			dbm.numIteratorCreated.Collect(mc)
		case "rocksdb.num.iterator.deleted":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.numIteratorDeleted.WithLabelValues().Set(v)
			dbm.numIteratorDeleted.Collect(mc)
		case "rocksdb.block.cache.compression.dict.miss":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheCompressionDictMiss.WithLabelValues().Set(v)
			dbm.blockCacheCompressionDictMiss.Collect(mc)
		case "rocksdb.block.cache.compression.dict.hit":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheCompressionDictHit.WithLabelValues().Set(v)
			dbm.blockCacheCompressionDictHit.Collect(mc)
		case "rocksdb.block.cache.compression.dict.add":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheCompressionDictAdd.WithLabelValues().Set(v)
			dbm.blockCacheCompressionDictAdd.Collect(mc)
		case "rocksdb.block.cache.compression.dict.bytes.insert":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheCompressionDictBytesInsert.WithLabelValues().Set(v)
			dbm.blockCacheCompressionDictBytesInsert.Collect(mc)
		case "rocksdb.block.cache.compression.dict.bytes.evict":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheCompressionDictBytesEvict.WithLabelValues().Set(v)
			dbm.blockCacheCompressionDictBytesEvict.Collect(mc)
		case "rocksdb.block.cache.add.redundant":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheAddRedundant.WithLabelValues().Set(v)
			dbm.blockCacheAddRedundant.Collect(mc)
		case "rocksdb.block.cache.index.add.redundant":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheIndexAddRedundant.WithLabelValues().Set(v)
			dbm.blockCacheIndexAddRedundant.Collect(mc)
		case "rocksdb.block.cache.filter.add.redundant":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheFilterAddRedundant.WithLabelValues().Set(v)
			dbm.blockCacheFilterAddRedundant.Collect(mc)
		case "rocksdb.block.cache.data.add.redundant":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheDataAddRedundant.WithLabelValues().Set(v)
			dbm.blockCacheDataAddRedundant.Collect(mc)
		case "rocksdb.block.cache.compression.dict.add.redundant":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.blockCacheCompressionDictAddRedundant.WithLabelValues().Set(v)
			dbm.blockCacheCompressionDictAddRedundant.Collect(mc)
		case "rocksdb.files.marked.trash":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.filesMarkedTrash.WithLabelValues().Set(v)
			dbm.filesMarkedTrash.Collect(mc)
		case "rocksdb.files.deleted.immediately":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.filesDeletedImmediately.WithLabelValues().Set(v)
			dbm.filesDeletedImmediately.Collect(mc)
		case "rocksdb.error.handler.bg.errro.count":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.errorHandlerBgErrroCount.WithLabelValues().Set(v)
			dbm.errorHandlerBgErrroCount.Collect(mc)
		case "rocksdb.error.handler.bg.io.errro.count":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.errorHandlerBgIoErrroCount.WithLabelValues().Set(v)
			dbm.errorHandlerBgIoErrroCount.Collect(mc)
		case "rocksdb.error.handler.bg.retryable.io.errro.count":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.errorHandlerBgRetryableIoErrroCount.WithLabelValues().Set(v)
			dbm.errorHandlerBgRetryableIoErrroCount.Collect(mc)
		case "rocksdb.error.handler.autoresume.count":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.errorHandlerAutoresumeCount.WithLabelValues().Set(v)
			dbm.errorHandlerAutoresumeCount.Collect(mc)
		case "rocksdb.error.handler.autoresume.retry.total.count":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.errorHandlerAutoresumeRetryTotalCount.WithLabelValues().Set(v)
			dbm.errorHandlerAutoresumeRetryTotalCount.Collect(mc)
		case "rocksdb.error.handler.autoresume.success.count":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.errorHandlerAutoresumeSuccessCount.WithLabelValues().Set(v)
			dbm.errorHandlerAutoresumeSuccessCount.Collect(mc)
		case "rocksdb.memtable.payload.bytes.at.flush":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.memtablePayloadBytesAtFlush.WithLabelValues().Set(v)
			dbm.memtablePayloadBytesAtFlush.Collect(mc)
		case "rocksdb.memtable.garbage.bytes.at.flush":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.memtableGarbageBytesAtFlush.WithLabelValues().Set(v)
			dbm.memtableGarbageBytesAtFlush.Collect(mc)
		case "rocksdb.secondary.cache.hits":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.secondaryCacheHits.WithLabelValues().Set(v)
			dbm.secondaryCacheHits.Collect(mc)
		case "rocksdb.verify_checksum.read.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.verifyChecksumReadBytes.WithLabelValues().Set(v)
			dbm.verifyChecksumReadBytes.Collect(mc)
		case "rocksdb.backup.read.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.backupReadBytes.WithLabelValues().Set(v)
			dbm.backupReadBytes.Collect(mc)
		case "rocksdb.backup.write.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.backupWriteBytes.WithLabelValues().Set(v)
			dbm.backupWriteBytes.Collect(mc)
		case "rocksdb.remote.compact.read.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.remoteCompactReadBytes.WithLabelValues().Set(v)
			dbm.remoteCompactReadBytes.Collect(mc)
		case "rocksdb.remote.compact.write.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.remoteCompactWriteBytes.WithLabelValues().Set(v)
			dbm.remoteCompactWriteBytes.Collect(mc)
		case "rocksdb.hot.file.read.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.hotFileReadBytes.WithLabelValues().Set(v)
			dbm.hotFileReadBytes.Collect(mc)
		case "rocksdb.warm.file.read.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.warmFileReadBytes.WithLabelValues().Set(v)
			dbm.warmFileReadBytes.Collect(mc)
		case "rocksdb.cold.file.read.bytes":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.coldFileReadBytes.WithLabelValues().Set(v)
			dbm.coldFileReadBytes.Collect(mc)
		case "rocksdb.hot.file.read.count":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.hotFileReadCount.WithLabelValues().Set(v)
			dbm.hotFileReadCount.Collect(mc)
		case "rocksdb.warm.file.read.count":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.warmFileReadCount.WithLabelValues().Set(v)
			dbm.warmFileReadCount.Collect(mc)
		case "rocksdb.cold.file.read.count":
			v, _ := strconv.ParseFloat(fs[3], 64)
			dbm.coldFileReadCount.WithLabelValues().Set(v)
			dbm.coldFileReadCount.Collect(mc)

		case "rocksdb.db.get.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.dbGetMicros, count, sum, quantiles)
		case "rocksdb.db.write.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.dbWriteMicros, count, sum, quantiles)
		case "rocksdb.compaction.times.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.compactionTimesMicros, count, sum, quantiles)
		case "rocksdb.compaction.times.cpu_micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.compactionTimesCpuMicros, count, sum, quantiles)
		case "rocksdb.subcompaction.setup.times.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.subcompactionSetupTimesMicros, count, sum, quantiles)
		case "rocksdb.table.sync.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.tableSyncMicros, count, sum, quantiles)
		case "rocksdb.compaction.outfile.sync.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.compactionOutfileSyncMicros, count, sum, quantiles)
		case "rocksdb.wal.file.sync.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.walFileSyncMicros, count, sum, quantiles)
		case "rocksdb.manifest.file.sync.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.manifestFileSyncMicros, count, sum, quantiles)
		case "rocksdb.table.open.io.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.tableOpenIoMicros, count, sum, quantiles)
		case "rocksdb.db.multiget.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.dbMultigetMicros, count, sum, quantiles)
		case "rocksdb.read.block.compaction.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.readBlockCompactionMicros, count, sum, quantiles)
		case "rocksdb.read.block.get.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.readBlockGetMicros, count, sum, quantiles)
		case "rocksdb.write.raw.block.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.writeRawBlockMicros, count, sum, quantiles)
		case "rocksdb.l0.slowdown.count":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.l0SlowdownCount, count, sum, quantiles)
		case "rocksdb.memtable.compaction.count":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.memtableCompactionCount, count, sum, quantiles)
		case "rocksdb.num.files.stall.count":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.numFilesStallCount, count, sum, quantiles)
		case "rocksdb.hard.rate.limit.delay.count":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.hardRateLimitDelayCount, count, sum, quantiles)
		case "rocksdb.soft.rate.limit.delay.count":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.softRateLimitDelayCount, count, sum, quantiles)
		case "rocksdb.numfiles.in.singlecompaction":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.numfilesInSinglecompaction, count, sum, quantiles)
		case "rocksdb.db.seek.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.dbSeekMicros, count, sum, quantiles)
		case "rocksdb.db.write.stall":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.dbWriteStall, count, sum, quantiles)
		case "rocksdb.sst.read.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.sstReadMicros, count, sum, quantiles)
		case "rocksdb.num.subcompactions.scheduled":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.numSubcompactionsScheduled, count, sum, quantiles)
		case "rocksdb.bytes.per.read":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.bytesPerRead, count, sum, quantiles)
		case "rocksdb.bytes.per.write":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.bytesPerWrite, count, sum, quantiles)
		case "rocksdb.bytes.per.multiget":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.bytesPerMultiget, count, sum, quantiles)
		case "rocksdb.bytes.compressed":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.bytesCompressed, count, sum, quantiles)
		case "rocksdb.bytes.decompressed":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.bytesDecompressed, count, sum, quantiles)
		case "rocksdb.compression.times.nanos":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.compressionTimesNanos, count, sum, quantiles)
		case "rocksdb.decompression.times.nanos":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.decompressionTimesNanos, count, sum, quantiles)
		case "rocksdb.read.num.merge_operands":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.readNumMergeOperands, count, sum, quantiles)
		case "rocksdb.blobdb.key.size":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbKeySize, count, sum, quantiles)
		case "rocksdb.blobdb.value.size":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbValueSize, count, sum, quantiles)
		case "rocksdb.blobdb.write.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbWriteMicros, count, sum, quantiles)
		case "rocksdb.blobdb.get.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbGetMicros, count, sum, quantiles)
		case "rocksdb.blobdb.multiget.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbMultigetMicros, count, sum, quantiles)
		case "rocksdb.blobdb.seek.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbSeekMicros, count, sum, quantiles)
		case "rocksdb.blobdb.next.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbNextMicros, count, sum, quantiles)
		case "rocksdb.blobdb.prev.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbPrevMicros, count, sum, quantiles)
		case "rocksdb.blobdb.blob.file.write.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbBlobFileWriteMicros, count, sum, quantiles)
		case "rocksdb.blobdb.blob.file.read.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbBlobFileReadMicros, count, sum, quantiles)
		case "rocksdb.blobdb.blob.file.sync.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbBlobFileSyncMicros, count, sum, quantiles)
		case "rocksdb.blobdb.gc.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbGcMicros, count, sum, quantiles)
		case "rocksdb.blobdb.compression.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbCompressionMicros, count, sum, quantiles)
		case "rocksdb.blobdb.decompression.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.blobdbDecompressionMicros, count, sum, quantiles)
		case "rocksdb.db.flush.micros":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.dbFlushMicros, count, sum, quantiles)
		case "rocksdb.sst.batch.size":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.sstBatchSize, count, sum, quantiles)
		case "rocksdb.num.index.and.filter.blocks.read.per.level":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.numIndexAndFilterBlocksReadPerLevel, count, sum, quantiles)
		case "rocksdb.num.data.blocks.read.per.level":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.numDataBlocksReadPerLevel, count, sum, quantiles)
		case "rocksdb.num.sst.read.per.level":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.numSstReadPerLevel, count, sum, quantiles)
		case "rocksdb.error.handler.autoresume.retry.count":
			count, sum, quantiles := parseSummaryStats(fs)
			mc <- prometheus.MustNewConstSummary(dbm.errorHandlerAutoresumeRetryCount, count, sum, quantiles)
		}
	}
}

func parseSummaryStats(fs []string) (count uint64, sum float64, quantiles map[float64]float64) {
	count, _ = strconv.ParseUint(fs[15], 10, 64)
	sum, _ = strconv.ParseFloat(fs[18], 64)

	quantiles = make(map[float64]float64)
	quantiles[0.5], _ = strconv.ParseFloat(fs[3], 64)
	quantiles[0.95], _ = strconv.ParseFloat(fs[6], 64)
	quantiles[0.99], _ = strconv.ParseFloat(fs[9], 64)
	quantiles[1], _ = strconv.ParseFloat(fs[12], 64)

	return
}

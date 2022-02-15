//go:build rocksdb
// +build rocksdb

package db

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/spf13/viper"
	"github.com/tecbot/gorocksdb"
)

func init() {
	dbCreator := func(name string, dir string) (DB, error) {
		return NewRocksDB(name, dir)
	}
	registerDBCreator(RocksDBBackend, dbCreator, false)
}

// RocksDB is a RocksDB backend.
type RocksDB struct {
	db     *gorocksdb.DB
	ro     *gorocksdb.ReadOptions
	wo     *gorocksdb.WriteOptions
	woSync *gorocksdb.WriteOptions
}

var _ DB = (*RocksDB)(nil)

const (
	blockSize    = "block_size"
	blockCache   = "block_cache"
	statistics   = "statistics"
	maxOpenFiles = "max_open_files"
	mmapRead     = "allow_mmap_reads"
	mmapWrite    = "allow_mmap_writes"
)

func NewRocksDB(name string, dir string) (*RocksDB, error) {
	// default rocksdb option, good enough for most cases, including heavy workloads.
	// 1GB table cache, 512MB write buffer(may use 50% more on heavy workloads).
	// compression: snappy as default, need to -lsnappy to enable.
	params := parseOptParams(viper.GetString(FlagRocksdbOpts))

	bbto := gorocksdb.NewDefaultBlockBasedTableOptions()
	if v, ok := params[blockSize]; ok {
		size, err := toBytes(v)
		if err != nil {
			panic(fmt.Sprintf("Invalid options parameter %s: %s", blockSize, err))
		}
		bbto.SetBlockSize(int(size))
	}

	blockCacheSize := 4096 * 1024 * 1024

	bbto.SetBlockCache(gorocksdb.NewLRUCache(blockCacheSize))
	if v, ok := params[blockCache]; ok {
		cache, err := toBytes(v)
		if err != nil {
			panic(fmt.Sprintf("Invalid options parameter %s: %s", blockCache, err))
		}
		bbto.SetBlockCache(gorocksdb.NewLRUCache(cache))
	}
	bbto.SetFilterPolicy(gorocksdb.NewBloomFilter(10))

	opts := gorocksdb.NewDefaultOptions()
	opts.SetBlockBasedTableFactory(bbto)
	opts.SetCreateIfMissing(true)
	opts.IncreaseParallelism(runtime.NumCPU())

	opts.EnableStatistics()
	if v, ok := params[statistics]; ok {
		enable, err := strconv.ParseBool(v)
		if err != nil {
			panic(fmt.Sprintf("Invalid options parameter %s: %s", statistics, err))
		}
		if enable {
			opts.EnableStatistics()
		}
	}

	opts.SetMaxOpenFiles(-1)
	if v, ok := params[maxOpenFiles]; ok {
		maxOpenFiles, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("Invalid options parameter %s: %s", maxOpenFiles, err))
		}
		opts.SetMaxOpenFiles(maxOpenFiles)
	}

	opts.SetAllowMmapReads(false)
	if v, ok := params[mmapRead]; ok {
		enable, err := strconv.ParseBool(v)
		if err != nil {
			panic(fmt.Sprintf("Invalid options parameter %s: %s", mmapRead, err))
		}
		opts.SetAllowMmapReads(enable)
	}

	if v, ok := params[mmapWrite]; ok {
		enable, err := strconv.ParseBool(v)
		if err != nil {
			panic(fmt.Sprintf("Invalid options parameter %s: %s", mmapWrite, err))
		}
		if enable {
			opts.SetAllowMmapWrites(enable)
		}
	}

	opts.OptimizeForPointLookup(blockCacheSize)

	// 1.5GB maximum memory use for writebuffer.
	opts.OptimizeLevelStyleCompaction(512 * 1024 * 1024)
	return NewRocksDBWithOptions(name, dir, opts)
}

func NewRocksDBWithOptions(name string, dir string, opts *gorocksdb.Options) (*RocksDB, error) {
	dbPath := filepath.Join(dir, name+".db")
	db, err := gorocksdb.OpenDb(opts, dbPath)
	if err != nil {
		return nil, err
	}
	ro := gorocksdb.NewDefaultReadOptions()
	wo := gorocksdb.NewDefaultWriteOptions()
	woSync := gorocksdb.NewDefaultWriteOptions()
	woSync.SetSync(true)
	database := &RocksDB{
		db:     db,
		ro:     ro,
		wo:     wo,
		woSync: woSync,
	}
	return database, nil
}

// Get implements DB.
func (db *RocksDB) Get(key []byte) ([]byte, error) {
	key = nonNilBytes(key)
	res, err := db.db.Get(db.ro, key)
	if err != nil {
		return nil, err
	}
	return moveSliceToBytes(res), nil
}

func (db *RocksDB) GetUnsafeValue(key []byte, processor UnsafeValueProcessor) (interface{}, error) {
	key = nonNilBytes(key)
	res, err := db.db.Get(db.ro, key)
	if err != nil {
		return nil, err
	}
	defer res.Free()
	if !res.Exists() {
		return processor(nil)
	}
	return processor(res.Data())
}

// Has implements DB.
func (db *RocksDB) Has(key []byte) (bool, error) {
	bytes, err := db.Get(key)
	if err != nil {
		return false, err
	}
	return bytes != nil, nil
}

// Set implements DB.
func (db *RocksDB) Set(key []byte, value []byte) error {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	err := db.db.Put(db.wo, key, value)
	if err != nil {
		return err
	}
	return nil
}

// SetSync implements DB.
func (db *RocksDB) SetSync(key []byte, value []byte) error {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	err := db.db.Put(db.woSync, key, value)
	if err != nil {
		return err
	}
	return nil
}

// Delete implements DB.
func (db *RocksDB) Delete(key []byte) error {
	key = nonNilBytes(key)
	err := db.db.Delete(db.wo, key)
	if err != nil {
		return err
	}
	return nil
}

// DeleteSync implements DB.
func (db *RocksDB) DeleteSync(key []byte) error {
	key = nonNilBytes(key)
	err := db.db.Delete(db.woSync, key)
	if err != nil {
		return nil
	}
	return nil
}

func (db *RocksDB) DB() *gorocksdb.DB {
	return db.db
}

// Close implements DB.
func (db *RocksDB) Close() error {
	db.ro.Destroy()
	db.wo.Destroy()
	db.woSync.Destroy()
	db.db.Close()
	return nil
}

// Print implements DB.
func (db *RocksDB) Print() error {
	itr, err := db.Iterator(nil, nil)
	if err != nil {
		return err
	}
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		key := itr.Key()
		value := itr.Value()
		fmt.Printf("[%X]:\t[%X]\n", key, value)
	}
	return nil
}

// Stats implements DB.
func (db *RocksDB) Stats() map[string]string {
	keys := []string{"rocksdb.stats"}
	stats := make(map[string]string, len(keys))
	for _, key := range keys {
		stats[key] = db.db.GetProperty(key)
	}
	return stats
}

// NewBatch implements DB.
func (db *RocksDB) NewBatch() Batch {
	return newRocksDBBatch(db)
}

// Iterator implements DB.
func (db *RocksDB) Iterator(start, end []byte) (Iterator, error) {
	itr := db.db.NewIterator(db.ro)
	return newRocksDBIterator(itr, start, end, false), nil
}

// ReverseIterator implements DB.
func (db *RocksDB) ReverseIterator(start, end []byte) (Iterator, error) {
	itr := db.db.NewIterator(db.ro)
	return newRocksDBIterator(itr, start, end, true), nil
}

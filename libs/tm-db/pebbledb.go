package db

import (
	"fmt"
	"github.com/cockroachdb/pebble/bloom"
	"path"

	"github.com/cockroachdb/pebble"
)

// PebbleDB implements DB.
type PebbleDB struct {
	db *pebble.DB
}

var _ DB = (*PebbleDB)(nil)

func init() {
	dbCreator := func(name string, dir string) (DB, error) {
		return NewPebbleDB(name, dir)
	}
	registerDBCreator(PebbleBackend, dbCreator, false)
}

// NewPebbleDB creates a *PebbleDB.
func NewPebbleDB(name, dir string) (*PebbleDB, error) {
	cache := pebble.NewCache(1 << 28)
	defer cache.Unref()
	opts := &pebble.Options{
		Cache:                       cache,
		L0CompactionThreshold:       4,
		L0StopWritesThreshold:       1000,
		LBaseMaxBytes:               64 << 20, // 64 MB
		Levels:                      make([]pebble.LevelOptions, 7),
		MaxConcurrentCompactions:    3,
		MaxOpenFiles:                16384,
		MemTableSize:                64 << 20,
		MemTableStopWritesThreshold: 4,
	}

	for i := 0; i < len(opts.Levels); i++ {
		l := &opts.Levels[i]
		l.BlockSize = 32 << 10       // 32 KB
		l.IndexBlockSize = 256 << 10 // 256 KB
		l.FilterPolicy = bloom.FilterPolicy(10)
		l.FilterType = pebble.TableFilter
		if i > 0 {
			l.TargetFileSize = opts.Levels[i-1].TargetFileSize * 2
		}
		l.EnsureDefaults()
	}
	opts.Levels[6].FilterPolicy = nil
	opts.FlushSplitBytes = opts.Levels[0].TargetFileSize

	opts.EnsureDefaults()
	return NewPebbleDBWithOptions(name, dir, opts)
}

func NewPebbleDBWithOptions(name, dir string, opts *pebble.Options) (*PebbleDB, error) {
	dbPath := path.Join(dir, name+".db")
	db, err := pebble.Open(dbPath, opts)
	if err != nil {
		return nil, err
	}
	return &PebbleDB{db: db}, nil
}

// Get implements DB.
func (db *PebbleDB) Get(key []byte) ([]byte, error) {
	key = nonNilBytes(key)
	res, closer, err := db.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	if err = closer.Close(); err != nil {
		return nil, err
	}
	return res, nil
}

// GetUnsafeValue implements DB.
func (db *PebbleDB) GetUnsafeValue(key []byte, processor UnsafeValueProcessor) (interface{}, error) {
	key = nonNilBytes(key)
	v, err := db.Get(key)
	if err != nil {
		return nil, err
	}
	return processor(v)
}

// Has implements DB.
func (db *PebbleDB) Has(key []byte) (bool, error) {
	value, err := db.Get(key)
	if err != nil {
		return false, err
	}
	return value != nil, nil
}

// Set implements DB.
func (db *PebbleDB) Set(key []byte, value []byte) error {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	return db.db.Set(key, value, nil)
}

// SetSync implements DB.
func (db *PebbleDB) SetSync(key []byte, value []byte) error {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	return db.db.Set(key, value, &pebble.WriteOptions{Sync: true})
}

// Delete implements DB.
func (db *PebbleDB) Delete(key []byte) error {
	key = nonNilBytes(key)
	return db.db.Delete(key, nil)
}

// DeleteSync implements DB.
func (db *PebbleDB) DeleteSync(key []byte) error {
	key = nonNilBytes(key)
	return db.db.Delete(key, &pebble.WriteOptions{Sync: true})
}

// DB implements DB.
func (db *PebbleDB) DB() *pebble.DB {
	return db.db
}

// Close implements DB.
func (db *PebbleDB) Close() error {
	return db.db.Close()
}

// Print implements DB.
func (db *PebbleDB) Print() error {
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
func (db *PebbleDB) Stats() map[string]string {
	return nil
}

// Iterator implements DB.
func (db *PebbleDB) Iterator(start []byte, end []byte) (Iterator, error) {
	iter := db.db.NewIter(&pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	})
	return newPebbleIterator(iter, start, end, false), nil
}

// ReverseIterator implements DB.
func (db *PebbleDB) ReverseIterator(start, end []byte) (Iterator, error) {
	iter := db.db.NewIter(&pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	})
	return newPebbleIterator(iter, start, end, true), nil
}

//NewBatch new
func (db *PebbleDB) NewBatch() Batch {
	batch := &pebbleBatch{
		batch: db.db.NewBatch(),
	}
	return batch
}

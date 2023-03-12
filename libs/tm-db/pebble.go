package db

import (
	"bytes"
	"path/filepath"
	"runtime"

	"github.com/cockroachdb/pebble/bloom"

	"github.com/cockroachdb/pebble"
)

func init() {
	dbCreator := func(name string, dir string) (DB, error) {
		return NewPebble(name, dir)
	}
	registerDBCreator(PebbleBackend, dbCreator, false)
}

type Pebble struct {
	db *pebble.DB
}

func NewPebble(name string, dir string) (*Pebble, error) {
	filedir := filepath.Join(dir, name+".db")
	opts := &pebble.Options{
		Cache:                       pebble.NewCache(int64(256 * 1024 * 1024)),
		MaxOpenFiles:                1024,
		MaxConcurrentCompactions:    func() int { return runtime.NumCPU() },
		MemTableSize:                256 * 1024 * 1024 / 2 / 2,
		MemTableStopWritesThreshold: 2,
		Levels: []pebble.LevelOptions{
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
		},
	}
	db, err := pebble.Open(filedir, opts)
	if err != nil {
		return nil, err
	}
	return &Pebble{db: db}, nil
}

func (p Pebble) Get(key []byte) ([]byte, error) {
	key = nonNilBytes(key)
	res, closer, err := p.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	ret := make([]byte, len(res))
	copy(ret, res)
	closer.Close()
	return ret, nil
}

func (p Pebble) GetUnsafeValue(key []byte, processor UnsafeValueProcessor) (interface{}, error) {
	key = nonNilBytes(key)
	res, closer, err := p.db.Get(key)
	if err != nil {
		if err == pebble.ErrNotFound {
			return processor(nil)
		}
		return nil, err
	}
	ret, err := processor(res)
	closer.Close()
	return ret, err
}

func (p Pebble) Has(key []byte) (bool, error) {
	_, closer, err := p.db.Get(key)
	if err == pebble.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	closer.Close()
	return true, nil
}

func (p Pebble) Set(key []byte, value []byte) error {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	return p.db.Set(key, value, pebble.NoSync)
}

func (p Pebble) SetSync(key []byte, value []byte) error {
	key = nonNilBytes(key)
	value = nonNilBytes(value)
	return p.db.Set(key, value, pebble.Sync)
}

func (p Pebble) Delete(key []byte) error {
	key = nonNilBytes(key)
	return p.db.Delete(key, nil)
}

func (p Pebble) DeleteSync(key []byte) error {
	key = nonNilBytes(key)
	return p.db.Delete(key, pebble.Sync)
}

func (p Pebble) Close() error {
	return p.db.Close()
}

func (p Pebble) Iterator(start, end []byte) (Iterator, error) {
	var iter *pebble.Iterator
	iter = p.db.NewIter(&pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	})
	return newPebbleIterator(iter, start, end, false), nil
}

func (p Pebble) ReverseIterator(start, end []byte) (Iterator, error) {
	var iter *pebble.Iterator
	iter = p.db.NewIter(&pebble.IterOptions{
		LowerBound: start,
		UpperBound: end,
	})
	return newPebbleIterator(iter, start, end, true), nil
}

func (p Pebble) NewBatch() Batch {
	batch := p.db.NewBatch()
	return pebbleBatch{batch}
}

func (p Pebble) Print() error {
	return nil
}

func (p Pebble) Stats() map[string]string {
	return map[string]string{}
}

func (p Pebble) Compact() error {
	return p.db.Compact(nil, bytes.Repeat([]byte{0xff}, 32), true)
}

func upperBound(prefix []byte) (limit []byte) {
	for i := len(prefix) - 1; i >= 0; i-- {
		c := prefix[i]
		if c == 0xff {
			continue
		}
		limit = make([]byte, i+1)
		copy(limit, prefix)
		limit[i] = c + 1
		break
	}
	return limit
}

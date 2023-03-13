package mpt

import (
	"bytes"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/go-amino"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/ethdb"
)

func (store *AsyncKeyValueStore) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	store.mtx.RLock()
	defer store.mtx.RUnlock()

	var (
		pr     = string(prefix)
		st     = string(append(prefix, start...))
		keys   = make([]string, 0, len(store.preCommit.data))
		values = make([]preCommitValue, 0, len(store.preCommit.data))
	)
	// Collect the keys from the memory database corresponding to the given prefix
	// and start
	for key := range store.preCommit.data {
		if !strings.HasPrefix(key, pr) {
			continue
		}
		if key >= st {
			keys = append(keys, key)
		}
	}
	// Sort the items and retrieve the associated values
	sort.Strings(keys)
	for _, key := range keys {
		values = append(values, store.preCommit.data[key])
	}

	dbIter := store.KeyValueStore.NewIterator(prefix, start)

	return &asyncdbIterator{
		dbIter: dbIter,
		keys:   keys,
		values: values,
	}
}

type asyncdbIterator struct {
	dbIter ethdb.Iterator
	keys   []string
	values []preCommitValue
}

func (a *asyncdbIterator) Next() bool {
	if len(a.keys) <= 0 {
		return a.dbIter.Next()
	}

	memkey := amino.StrToBytes(a.keys[0])
	dbKey := a.dbIter.Key()
	switch bytes.Compare(memkey, dbKey) {
	case -1:
		if a.memNext() {
			if a.values[0].deleted {
				return a.Next()
			}
		}
		return a.dbIter.Next()
	case 0:
		if a.memNext() {
			if a.values[0].deleted {
				return a.Next()
			}
			a.dbIter.Next()
			return true
		}
		return a.dbIter.Next()
	case 1:
		if a.dbIter.Next() {
			return true
		}
		if a.memNext() {
			if a.values[0].deleted {
				return a.Next()
			}
			return true
		}
		return false
	}
	return false
}

func (a *asyncdbIterator) memNext() bool {
	a.keys = a.keys[1:]
	a.values = a.values[1:]
	return len(a.keys) > 0
}

func (a *asyncdbIterator) Key() []byte {
	if len(a.keys) <= 0 {
		return a.dbIter.Key()
	}
	memkey := amino.StrToBytes(a.keys[0])
	dbKey := a.dbIter.Key()
	switch bytes.Compare(memkey, dbKey) {
	case -1:
		if a.values[0].deleted {
			a.memNext()
			return a.Key()
		}
		return common.CopyBytes(memkey)
	case 0:
		if a.values[0].deleted {
			a.memNext()
			a.dbIter.Next()
			return a.Key()
		}
		return common.CopyBytes(memkey)
	case 1:
		return dbKey
	}
	return nil
}

func (a *asyncdbIterator) Value() []byte {
	if len(a.keys) <= 0 {
		return a.dbIter.Value()
	}
	memkey := amino.StrToBytes(a.keys[0])
	dbKey := a.dbIter.Key()
	switch bytes.Compare(memkey, dbKey) {
	case -1:
		if a.values[0].deleted {
			a.memNext()
			return a.Value()
		}
		return common.CopyBytes(a.values[0].value)
	case 0:
		if a.values[0].deleted {
			a.memNext()
			a.dbIter.Next()
			return a.Value()
		}
		return common.CopyBytes(a.values[0].value)
	case 1:
		return a.dbIter.Value()
	}
	return nil
}

func (a *asyncdbIterator) Error() error {
	return a.dbIter.Error()
}

func (a *asyncdbIterator) Release() {
	a.dbIter.Release()
}

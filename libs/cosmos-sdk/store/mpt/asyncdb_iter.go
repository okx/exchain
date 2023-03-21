package mpt

import (
	"bytes"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/tendermint/go-amino"
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

	dbIter := store.KeyValueStore.NewIterator(prefix, start)
	if len(keys) == 0 {
		return dbIter
	}

	// Sort the items and retrieve the associated values
	sort.Strings(keys)
	for _, key := range keys {
		values = append(values, store.preCommit.data[key])
	}

	return &asyncdbIterator{
		dbIter: dbIter,
		keys:   keys,
		values: values,

		memMoveNext: true,
		dbMoveNext:  dbIter.Next(),
		dbKey:       dbIter.Key(),
	}
}

type asyncdbIterator struct {
	dbIter ethdb.Iterator
	keys   []string
	values []preCommitValue

	dbKey []byte
	key   []byte
	value []byte

	dbMoveNext  bool
	memMoveNext bool
	moveMemOrDb byte // 0 for mem, 1 for db, 2 for both
	init        bool
}

func (a *asyncdbIterator) Next() bool {
	if len(a.keys) == 0 {
		return a.dbIter.Next()
	}

	if !a.init {
		a.init = true
	} else {
		if a.moveMemOrDb == 0 {
			a.memMoveNext = a.memNext()
		} else if a.moveMemOrDb == 1 {
			a.dbMoveNext = a.dbIter.Next()
			a.dbKey = a.dbIter.Key()
		} else if a.moveMemOrDb == 2 {
			a.memMoveNext = a.memNext()
			a.dbMoveNext = a.dbIter.Next()
			a.dbKey = a.dbIter.Key()
		}
	}

	if !a.memMoveNext && !a.dbMoveNext {
		a.key = nil
		a.value = nil
		return false
	}

	if !a.dbMoveNext {
		a.moveMemOrDb = 0
		if a.values[0].deleted {
			return a.Next()
		}
		a.key = amino.StrToBytes(a.keys[0])
		a.value = a.values[0].value
		return true
	}

	if !a.memMoveNext {
		a.key = a.dbKey
		a.value = a.dbIter.Value()
		return true
	}

	memkey := amino.StrToBytes(a.keys[0])
	switch bytes.Compare(memkey, a.dbKey) {
	case -1:
		a.moveMemOrDb = 0
		if a.values[0].deleted {
			return a.Next()
		}
		a.key = memkey
		a.value = a.values[0].value
		return true
	case 0:
		a.moveMemOrDb = 2
		if a.values[0].deleted {
			return a.Next()
		}
		a.key = memkey
		a.value = a.values[0].value
		return true
	case 1:
		a.moveMemOrDb = 1
		a.key = a.dbKey
		a.value = a.dbIter.Value()
		return true
	}
	return false
}

func (a *asyncdbIterator) memNext() bool {
	a.keys = a.keys[1:]
	a.values = a.values[1:]
	return len(a.keys) > 0
}

func (a *asyncdbIterator) Key() []byte {
	if len(a.keys) == 0 {
		return a.dbIter.Key()
	}
	return a.key
}

func (a *asyncdbIterator) Value() []byte {
	if len(a.keys) == 0 {
		return a.dbIter.Value()
	}
	return a.value
}

func (a *asyncdbIterator) Error() error {
	return a.dbIter.Error()
}

func (a *asyncdbIterator) Release() {
	a.dbIter.Release()
}

package mpt

import (
	"bytes"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/tendermint/go-amino"
)

var iterMapPool = &sync.Pool{
	New: func() interface{} {
		return make(map[string]int)
	},
}

func (store *AsyncKeyValueStore) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	atomic.AddInt64(&store.iterNum, 1)
	if atomic.LoadInt64(&store.waitCommit) == 0 {
		return store.KeyValueStore.NewIterator(prefix, start)
	}

	var (
		pr    = string(prefix)
		st    = string(append(prefix, start...))
		count int
	)

	store.mtx.RLock()
	defer store.mtx.RUnlock()

	curPtr := store.getPreCommitPtr()
	for curPtr.Next() != nil {
		curPtr = curPtr.Next()
		switch op := curPtr.Value.(*commitTask).op.(type) {
		case *singleOp:
			count++
		case multiOp:
			count += len(op)
		}
	}

	var (
		ops = make([]singleOp, 0, count)
	)

	if atomic.LoadInt64(&store.waitCommit) == 0 {
		return store.KeyValueStore.NewIterator(prefix, start)
	}

	var m = iterMapPool.Get().(map[string]int)
	defer iterMapPool.Put(m)
	for k := range m {
		delete(m, k)
	}

	var singleOpHandler = func(op *singleOp, pr, st string, m map[string]int, ops *[]singleOp) {
		strKey := amino.BytesToStr(op.key)
		if !strings.HasPrefix(strKey, pr) {
			return
		}
		if strKey >= st {
			if index, ok := m[strKey]; ok {
				(*ops)[index] = *op
			} else {
				*ops = append(*ops, *op)
				m[strKey] = len(*ops) - 1
			}
		}
	}

	// Collect the keys from the memory database corresponding to the given prefix
	// and start
	curPtr = store.getPreCommitPtr()
	for curPtr.Next() != nil {
		curPtr = curPtr.Next()
		switch op := curPtr.Value.(*commitTask).op.(type) {
		case *singleOp:
			singleOpHandler(op, pr, st, m, &ops)
		case multiOp:
			for _, o := range op {
				singleOpHandler(&o, pr, st, m, &ops)
				if atomic.LoadInt64(&store.waitCommit) == 0 {
					return store.KeyValueStore.NewIterator(prefix, start)
				}
			}
		}
		if atomic.LoadInt64(&store.waitCommit) == 0 {
			return store.KeyValueStore.NewIterator(prefix, start)
		}
	}

	dbIter := store.KeyValueStore.NewIterator(prefix, start)
	if len(ops) == 0 {
		return dbIter
	}

	// Sort the items and retrieve the associated values
	sort.Slice(ops, func(i, j int) bool {
		return bytes.Compare(ops[i].key, ops[j].key) == -1
	})

	return &asyncdbIterator{
		dbIter: dbIter,
		ops:    ops,

		memMoveNext: true,
		dbMoveNext:  dbIter.Next(),
		dbKey:       dbIter.Key(),
	}
}

type asyncdbIterator struct {
	dbIter ethdb.Iterator
	ops    []singleOp

	dbKey []byte
	key   []byte
	value []byte

	dbMoveNext  bool
	memMoveNext bool
	moveMemOrDb byte // 0 for mem, 1 for db, 2 for both
	init        bool
}

func (a *asyncdbIterator) Next() bool {
	if len(a.ops) == 0 {
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
		if a.ops[0].delete {
			return a.Next()
		}
		a.key = a.ops[0].key
		a.value = a.ops[0].value
		return true
	}

	if !a.memMoveNext {
		a.key = a.dbKey
		a.value = a.dbIter.Value()
		return true
	}

	memkey := a.ops[0].key
	switch bytes.Compare(memkey, a.dbKey) {
	case -1:
		a.moveMemOrDb = 0
		if a.ops[0].delete {
			return a.Next()
		}
		a.key = memkey
		a.value = a.ops[0].value
		return true
	case 0:
		a.moveMemOrDb = 2
		if a.ops[0].delete {
			return a.Next()
		}
		a.key = memkey
		a.value = a.ops[0].value
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
	a.ops = a.ops[1:]
	return len(a.ops) > 0
}

func (a *asyncdbIterator) Key() []byte {
	if len(a.ops) == 0 {
		return a.dbIter.Key()
	}
	return a.key
}

func (a *asyncdbIterator) Value() []byte {
	if len(a.ops) == 0 {
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

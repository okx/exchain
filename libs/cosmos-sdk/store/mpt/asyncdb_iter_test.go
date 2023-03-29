package mpt

import (
	"bytes"
	"container/list"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
)

// this mock type should be used only for iterator test.
// it will not commit the changes to the db.
type mockIterTestAsycnDB struct {
	*AsyncKeyValueStore
}

func (store *mockIterTestAsycnDB) Put(key []byte, value []byte) error {
	key, value = common.CopyBytes(key), common.CopyBytes(value)
	task := &commitTask{
		op: &singleOp{
			key:   key,
			value: value,
		},
	}
	store.mtx.Lock()
	store.preCommitTail = store.preCommitList.PushBack(task)
	_ = store.preCommit.Put(key, value)
	store.mtx.Unlock()
	atomic.AddInt64(&store.waitCommit, 1)
	return nil
}

func (store *mockIterTestAsycnDB) Delete(key []byte) error {
	key = common.CopyBytes(key)
	task := &commitTask{
		op: &singleOp{
			key:    key,
			delete: true,
		},
	}
	store.mtx.Lock()
	store.preCommitTail = store.preCommitList.PushBack(task)
	_ = store.preCommit.Delete(key)
	store.mtx.Unlock()
	atomic.AddInt64(&store.waitCommit, 1)
	return nil
}

func TestAsyncdbIterator(t *testing.T) {
	memDb := memorydb.New()
	asyncDb := &mockIterTestAsycnDB{
		NewAsyncKeyValueStoreWithOptions(memDb, AsyncKeyValueStoreOptions{
			DisableAutoPrune: true,
		}),
	}

	iter := asyncDb.NewIterator(nil, nil)
	require.False(t, iter.Next())
	require.Nil(t, iter.Key())
	require.Nil(t, iter.Value())
	require.NoError(t, iter.Error())

	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("key%2d", i))
		value := []byte(fmt.Sprintf("value%2d", i))
		if i%2 == 0 {
			require.NoError(t, asyncDb.Put(key, value))
		} else {
			require.NoError(t, memDb.Put(key, value))
		}
	}

	iter = asyncDb.NewIterator(nil, nil)

	for i := 0; i < 100; i++ {
		require.True(t, iter.Next())
		key := []byte(fmt.Sprintf("key%2d", i))
		value := []byte(fmt.Sprintf("value%2d", i))
		require.Equal(t, string(key), string(iter.Key()))
		require.Equal(t, string(value), string(iter.Value()))
	}
	require.False(t, iter.Next())
	require.Nil(t, iter.Key())
	require.Nil(t, iter.Value())
	require.NoError(t, iter.Error())

	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("key%2d", i))
		value := []byte(fmt.Sprintf("value%2d", i))
		require.NoError(t, asyncDb.Delete(key))
		require.NoError(t, memDb.Put(key, value))
	}
	iter = asyncDb.NewIterator(nil, nil)
	require.False(t, iter.Next())
	require.Nil(t, iter.Key())
	require.Nil(t, iter.Value())
	require.NoError(t, iter.Error())

	asyncDb = &mockIterTestAsycnDB{
		NewAsyncKeyValueStoreWithOptions(memDb, AsyncKeyValueStoreOptions{
			DisableAutoPrune: true,
		}),
	}
	iter = asyncDb.NewIterator(nil, nil)
	for i := 0; i < 100; i++ {
		require.True(t, iter.Next())
		key := []byte(fmt.Sprintf("key%2d", i))
		value := []byte(fmt.Sprintf("value%2d", i))
		require.Equal(t, string(key), string(iter.Key()))
		require.Equal(t, string(value), string(iter.Value()))
	}
	require.False(t, iter.Next())
	require.Nil(t, iter.Key())
	require.Nil(t, iter.Value())
	require.NoError(t, iter.Error())

	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("key%2d", i))
		value := []byte(fmt.Sprintf("value%2d", i))
		require.NoError(t, asyncDb.Put(key, value))
		require.NoError(t, memDb.Delete(key))
	}

	require.Equal(t, 100, asyncDb.preCommit.Len())
	require.Equal(t, 0, memDb.Len())

	iter = asyncDb.NewIterator(nil, nil)
	for i := 0; i < 100; i++ {
		require.True(t, iter.Next())
		key := []byte(fmt.Sprintf("key%2d", i))
		value := []byte(fmt.Sprintf("value%2d", i))
		require.Equal(t, string(key), string(iter.Key()))
		require.Equal(t, string(value), string(iter.Value()))
	}
	require.False(t, iter.Next())
	require.Nil(t, iter.Key())
	require.Nil(t, iter.Value())
	require.NoError(t, iter.Error())

	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("key%2d", i))
		value := []byte(fmt.Sprintf("value%2d", i))
		newValue := []byte(fmt.Sprintf("value%2d", i*2))
		if i%2 != 0 {
			require.NoError(t, asyncDb.Put(key, newValue))
		}
		if i%4 != 0 {
			require.NoError(t, asyncDb.Delete(key))
		}
		require.NoError(t, memDb.Put(key, value))
	}
	asyncDb.Put([]byte("z"), []byte("z"))

	iter = asyncDb.NewIterator(nil, nil)
	for i := 0; i < 100; i++ {
		if i%4 != 0 {
			continue
		}
		require.True(t, iter.Next())
		key := []byte(fmt.Sprintf("key%2d", i))
		value := []byte(fmt.Sprintf("value%2d", i))
		newValue := []byte(fmt.Sprintf("value%2d", i*2))
		require.Equal(t, string(key), string(iter.Key()))
		if i%2 != 0 {
			require.Equal(t, string(newValue), string(iter.Value()))
		} else {
			require.Equal(t, string(value), string(iter.Value()))
		}
	}
	require.True(t, iter.Next())
	require.Equal(t, []byte("z"), iter.Key())
	require.Equal(t, []byte("z"), iter.Value())
	require.False(t, iter.Next())
	require.Nil(t, iter.Key())
	require.Nil(t, iter.Value())
	require.NoError(t, iter.Error())

	memDb = memorydb.New()
	asyncDb = &mockIterTestAsycnDB{
		NewAsyncKeyValueStoreWithOptions(memDb, AsyncKeyValueStoreOptions{
			DisableAutoPrune: true,
		}),
	}

	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("key%2d", i))
		value := []byte(fmt.Sprintf("value%2d", i))
		if i%2 == 0 {
			require.NoError(t, asyncDb.Put(key, value))
		} else {
			require.NoError(t, memDb.Put(key, value))
		}
	}

	iter = asyncDb.NewIterator([]byte("key"), nil)

	for i := 0; i < 100; i++ {
		require.True(t, iter.Next())
		key := []byte(fmt.Sprintf("key%2d", i))
		value := []byte(fmt.Sprintf("value%2d", i))
		require.Equal(t, string(key), string(iter.Key()))
		require.Equal(t, string(value), string(iter.Value()))
	}
	require.False(t, iter.Next())
	require.Nil(t, iter.Key())
	require.Nil(t, iter.Value())
	require.NoError(t, iter.Error())

	iter = asyncDb.NewIterator([]byte("ke"), []byte("y"))

	for i := 0; i < 100; i++ {
		require.True(t, iter.Next())
		key := []byte(fmt.Sprintf("key%2d", i))
		value := []byte(fmt.Sprintf("value%2d", i))
		require.Equal(t, string(key), string(iter.Key()))
		require.Equal(t, string(value), string(iter.Value()))
	}
	require.False(t, iter.Next())
	require.Nil(t, iter.Key())
	require.Nil(t, iter.Value())
	require.NoError(t, iter.Error())

	iter = asyncDb.NewIterator([]byte("z"), []byte(""))
	require.False(t, iter.Next())
	require.Nil(t, iter.Key())
	require.Nil(t, iter.Value())
	require.NoError(t, iter.Error())
}

func BenchmarkAsyncdbIterator(b *testing.B) {
	m := make(map[string][]byte, 10_0000)
	l := list.New()

	for i := 0; i < 10; i++ {
		prefix := bytes.Repeat([]byte{byte(i)}, 32)
		ops := make(multiOp, 0, 1_0000)
		for i := 0; i < 1_0000; i++ {
			key := append(prefix, []byte(fmt.Sprintf("key%d", i))...)
			value := []byte(fmt.Sprintf("value%d", i))
			m[string(key)] = value
			ops = append(ops, singleOp{
				key:   key,
				value: value,
			})
		}
		task := &commitTask{
			op: ops,
		}
		l.PushBack(task)
	}

	ptr := l.Back().Prev().Prev()

	var pool = &sync.Pool{
		New: func() interface{} {
			return make(map[string]int)
		},
	}

	var pool2 = &sync.Pool{
		New: func() interface{} {
			ret := make([]singleOp, 0)
			return &ret
		},
	}

	var (
		prefix = bytes.Repeat([]byte{byte(5)}, 32)
		pr     = string(prefix)
		st     = string(append(prefix, []byte{}...))
	)

	isStEqPr := st == pr

	b.ResetTimer()

	b.Run("1", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var (
				keys   = make([]string, 0, 10_0000)
				values = make([][]byte, 0, 10_0000)
			)
			for key := range m {
				if !strings.HasPrefix(key, pr) {
					continue
				}
				if isStEqPr || key >= st {
					keys = append(keys, key)
				}
			}
			sort.Strings(keys)
			for _, key := range keys {
				values = append(values, m[key])
			}
		}
	})

	b.Run("2", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var count int

			curPtr := ptr
			for curPtr.Next() != nil {
				curPtr = curPtr.Next()
				switch op := curPtr.Value.(*commitTask).op.(type) {
				case *singleOp:
					count++
				case multiOp:
					count += len(op)
				}
			}
			count = 2_0000

			var ops = make([]singleOp, 0, count)
			//var opsPtr = pool2.Get().(*[]singleOp)
			//defer pool2.Put(opsPtr)
			//var ops = *opsPtr
			//if cap(ops) < count {
			//	pool2.Put(opsPtr)
			//	ops = make([]singleOp, 0, count)
			//	opsPtr = &ops
			//}
			ops = ops[:0]

			//var m = make(map[string]int, count)
			var m = pool.Get().(map[string]int)
			defer pool.Put(m)
			for k := range m {
				delete(m, k)
			}

			var singleOpHandler = func(op *singleOp, pr, st string, m map[string]int, ops *[]singleOp) {
				strKey := amino.BytesToStr(op.key)
				if !strings.HasPrefix(strKey, pr) {
					return
				}
				if isStEqPr || strKey >= st {
					if index, ok := m[strKey]; ok {
						(*ops)[index] = *op
					} else {
						*ops = append(*ops, *op)
						m[strKey] = len(*ops) - 1
					}
				}
			}

			curPtr = ptr
			for curPtr.Next() != nil {
				curPtr = curPtr.Next()
				switch op := curPtr.Value.(*commitTask).op.(type) {
				case *singleOp:
					singleOpHandler(op, pr, st, m, &ops)
				case multiOp:
					for _, o := range op {
						singleOpHandler(&o, pr, st, m, &ops)
					}
				}
			}

			sort.Slice(ops, func(i, j int) bool {
				return bytes.Compare(ops[i].key, ops[j].key) == -1
			})
		}
	})

	_ = pool
	_ = pool2
}

func BenchmarkAtomic(b *testing.B) {
	var num = int64(rand.Int31n(100))
	var res *int

	b.Run("1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if num == 200 {
				old := res
				res = new(int)
				*res = *old + 1
			}
		}
	})

	b.Run("2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if atomic.LoadInt64(&num) == 200 {
				old := res
				res = new(int)
				*res = *old + 1
			}
		}
	})

	_ = res
}

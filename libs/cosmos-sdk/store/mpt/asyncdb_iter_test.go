package mpt

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/stretchr/testify/require"
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
		NewAsyncKeyValueStore(memDb, true),
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
		NewAsyncKeyValueStore(memDb, true),
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
}

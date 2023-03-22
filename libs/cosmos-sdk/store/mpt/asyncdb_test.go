package mpt

import (
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/stretchr/testify/require"
)

func TestAsyncDB(t *testing.T) {
	memDb := memorydb.New()

	asyncDb := NewAsyncKeyValueStoreWithOptions(memDb, AsyncKeyValueStoreOptions{
		DisableAutoPrune: true,
		SyncPrune:        true,
	})

	t.Logf("asyncDB started")

	require.NoError(t, asyncDb.Put([]byte("key1"), []byte("value1")))
	require.NoError(t, asyncDb.Put([]byte("key2"), []byte("value2")))
	require.NoError(t, asyncDb.Put([]byte("key3"), []byte("value3")))

	ok, err := asyncDb.Has([]byte("key1"))
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = asyncDb.Has([]byte("key2"))
	require.NoError(t, err)
	require.True(t, ok)

	ok, err = asyncDb.Has([]byte("key3"))
	require.NoError(t, err)
	require.True(t, ok)

	require.NoError(t, asyncDb.Delete([]byte("key1")))

	ok, err = asyncDb.Has([]byte("key1"))
	require.NoError(t, err)
	require.False(t, ok)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	asyncDb.ActionAfterWriteDone(func() { wg.Done() }, true)

	wg.Wait()
	require.Equal(t, 2, memDb.Len())

	require.EqualValues(t, 5, asyncDb.waitPrune)
	asyncDb.Prune()
	require.EqualValues(t, 0, asyncDb.waitPrune)
	require.Equal(t, 0, asyncDb.preCommit.Len())

	err = asyncDb.Close()
	require.NoError(t, err)

	t.Logf("asyncDB closed")
}

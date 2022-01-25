package mempool

import (
	"crypto/rand"
	"testing"

	"github.com/okex/exchain/libs/tendermint/abci/example/kvstore"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/stretchr/testify/require"
)

func TestCacheAdd(t *testing.T) {
	size := 10
	cache := newLruCache(size)
	for i := 0; i < 10; i++ {
		tx := genTestTx(i)
		require.False(t, cache.Contains(tx))
		cache.Add(tx, nil)
		require.True(t, cache.Contains(tx))
	}
	require.Equal(t, size, cache.lru.Len())

	tx := genTestTx(5)
	data := cache.Get(tx)
	_, ok := data.(*WrappedTx)
	require.False(t, ok)
	cache.Add(tx, &WrappedTx{Payload: tx})
	data = cache.Get(tx)
	_, ok = data.(*WrappedTx)
	require.True(t, ok)

	require.True(t, cache.Contains(genTestTx(0)))
	tx = genTestTx(size + 1)
	cache.Add(tx, nil)
	require.True(t, cache.Contains(tx))
	// tx-0 evicted
	require.False(t, cache.Contains(genTestTx(0)))

	cache.Reset()
	require.Equal(t, 0, cache.lru.Len())
}

func genTestTx(i int) types.Tx {
	return append([]byte("test-tx-"), byte(i))
}

func TestCacheRemove(t *testing.T) {
	cache := newLruCache(100)
	numTxs := 10
	txs := make([][]byte, numTxs)
	for i := 0; i < numTxs; i++ {
		// probability of collision is 2**-256
		txBytes := make([]byte, 32)
		rand.Read(txBytes) // nolint: gosec
		txs[i] = txBytes
		cache.Add(txBytes, nil)
		require.Equal(t, i+1, cache.lru.Len())
	}
	for i := 0; i < numTxs; i++ {
		cache.Remove(txs[i])
		require.Equal(t, numTxs-(i+1), cache.lru.Len())
	}
}

func TestCacheAfterUpdate(t *testing.T) {
	app := kvstore.NewApplication()
	cc := proxy.NewLocalClientCreator(app)
	mempool, cleanup := newMempoolWithApp(cc)
	defer cleanup()

	// reAddIndices & txsInCache can have elements > numTxsToCreate
	// also assumes max index is 255 for convenience
	// txs in cache also checks order of elements
	tests := []struct {
		numTxsToCreate int
		updateIndices  []int
		reAddIndices   []int
		txsInCache     []int
	}{
		{1, []int{}, []int{1}, []int{1, 0}},    // adding new txs works
		{2, []int{1}, []int{}, []int{1, 0}},    // update doesn't remove tx from cache
		{2, []int{2}, []int{}, []int{2, 1, 0}}, // update adds new tx to cache
		{2, []int{1}, []int{1}, []int{1, 0}},   // re-adding after update doesn't make dupe
	}
	for tcIndex, tc := range tests {
		for i := 0; i < tc.numTxsToCreate; i++ {
			tx := types.Tx{byte(i)}
			err := mempool.CheckTx(tx, nil, TxInfo{})
			require.NoError(t, err)
			require.True(t, mempool.cache.Contains(tx))
		}

		updateTxs := []types.Tx{}
		for _, v := range tc.updateIndices {
			tx := types.Tx{byte(v)}
			updateTxs = append(updateTxs, tx)
		}
		mempool.Update(int64(tcIndex), updateTxs, abciResponses(len(updateTxs), abci.CodeTypeOK), nil, nil)

		for _, v := range tc.reAddIndices {
			tx := types.Tx{byte(v)}
			_ = mempool.CheckTx(tx, nil, TxInfo{})
		}

		for _, v := range tc.txsInCache {
			tx := types.Tx{byte(v)}
			require.True(t, mempool.cache.Contains(tx))
		}

		cache := mempool.cache.(*lruCache)
		require.Equal(t, len(tc.txsInCache), cache.lru.Len())
		mempool.Flush()
	}
}

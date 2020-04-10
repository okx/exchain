package cache

import (
	"testing"

	"github.com/okex/okchain/x/backend/types"
	"github.com/okex/okchain/x/common"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	cache := NewCache()
	require.Equal(t, 0, len(cache.Transactions))
	require.Equal(t, 2000, cap(cache.Transactions))
	require.Equal(t, 0, len(cache.LatestTicker))
	require.Equal(t, 0, len(cache.ProductsBuf))
	require.Equal(t, 200, cap(cache.ProductsBuf))

	txs := []*types.Transaction{
		{"hash1", types.TxTypeTransfer, "addr1", common.TestToken, types.TxSideFrom, "10.0", "0.1" + common.NativeToken, 100},
		{"hash2", types.TxTypeOrderNew, "addr1", types.TestTokenPair, types.TxSideBuy, "10.0", "0.1" + common.NativeToken, 300},
		{"hash3", types.TxTypeOrderCancel, "addr1", types.TestTokenPair, types.TxSideSell, "10.0", "0.1" + common.NativeToken, 200},
		{"hash4", types.TxTypeTransfer, "addr2", common.TestToken, types.TxSideTo, "10.0", "0.1" + common.NativeToken, 100},
	}

	for _, tx := range txs {
		cache.AddTransaction(tx)
	}

	require.Equal(t, txs, cache.GetTransactions())

	cache.Flush()
	require.Equal(t, 0, len(cache.GetTransactions()))

}

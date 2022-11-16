package global

import (
	"github.com/stretchr/testify/require"
	"math"
	"math/big"
	"math/rand"
	"testing"
)

func TestCacheBlockEvmTxGasUsed(t *testing.T) {
	data := make(map[int64]*big.Int)
	for i := 0; i < 1000; i++ {
		data[int64(i)] = big.NewInt(rand.Int63())
		SetBlockEvmTxGasUsed(int64(i), data[int64(i)])
	}

	for i := math.MaxInt; i > math.MaxInt-1000; i-- {
		data[int64(i)] = big.NewInt(rand.Int63())
		SetBlockEvmTxGasUsed(int64(i), data[int64(i)])
	}

	for k, v := range data {
		cached := GetBlockEvmTxGasUsed(k)
		require.NotNil(t, cached)
		cmp := v.Cmp(cached)
		require.Equal(t, 0, cmp)
	}
}

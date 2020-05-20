package token

import (
	"github.com/okex/okchain/x/token/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCache_AddFeeDetail(t *testing.T) {
	cache := NewCache()
	cache.reset()
}

func TestCache_Clone(t *testing.T) {
	cache := NewCache()

	feeDetail := types.FeeDetail{
		Address:   "xxxxx",
		Fee:       "1000okt",
		FeeType:   "transfer",
		Timestamp: 0,
	}
	cache.addFeeDetail(&feeDetail)

	cloneCache := cache.Clone()
	require.EqualValues(t, len(cache.FeeDetails), len(cloneCache.FeeDetails))

	feeDetail.Address = "yyyyy"
	cache.addFeeDetail(&feeDetail)

	require.EqualValues(t, len(cache.FeeDetails), len(cloneCache.FeeDetails) + 1)
}
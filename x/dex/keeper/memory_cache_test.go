package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	cache := NewCache()
	tokenPair := GetBuiltInTokenPair()

	// AddTokenPair successful
	cache.AddTokenPair(tokenPair)

	// GetTokenPair successful
	product := cache.genTokenPairKey(tokenPair)
	getTokenPair, ok := cache.GetTokenPair(product)
	require.Equal(t, tokenPair, getTokenPair)
	require.Equal(t, true, ok)

	// GetTokenPair failed
	getTokenPair, ok = cache.GetTokenPair(TestProductNotExist)
	require.Nil(t, getTokenPair)
	require.Equal(t, false, ok)
}

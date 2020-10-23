package keeper

import (
	"github.com/okex/okexchain/x/farm/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSatisfyWhiteListAdmittance(t *testing.T) {
	ctx, k := GetKeeper(t)
	quoteSymbol := types.DefaultParams().QuoteSymbol

	baseSymbol := "xxb"
	pool1 := types.FarmPool{
		Name:         "pool1",
		LockedSymbol: baseSymbol,
	}
	k.SetFarmPool(ctx, pool1)
	err := k.satisfyWhiteListAdmittance(ctx, pool1)
	require.NotNil(t, err)
	require.Equal(t, types.CodeTokenNotExist, err.Code())

	// test add LPT
	pool2 := types.FarmPool{
		Name:         "pool2",
		LockedSymbol: "ammswap_" + baseSymbol + "_" + quoteSymbol,
	}
	k.SetFarmPool(ctx, pool2)
	err = k.satisfyWhiteListAdmittance(ctx, pool2)
	require.NotNil(t, err)
	require.Equal(t, types.CodeTokenNotExist, err.Code())

	base1Symbol := "okb"
	pool3 := types.FarmPool{
		Name:         "pool3",
		LockedSymbol: "ammswap_" + baseSymbol + "_" + base1Symbol,
	}
	err = k.satisfyWhiteListAdmittance(ctx, pool3)
	require.NotNil(t, err)
	require.Equal(t, types.CodeTokenNotExist, err.Code())

	SetSwapTokenPair(ctx, k.Keeper, baseSymbol, quoteSymbol)

	err = k.satisfyWhiteListAdmittance(ctx, pool1)
	require.Nil(t, err)

	err = k.satisfyWhiteListAdmittance(ctx, pool2)
	require.Nil(t, err)

	err = k.satisfyWhiteListAdmittance(ctx, pool3)
	require.NotNil(t, err)
	require.Equal(t, types.CodeTokenNotExist, err.Code())

	SetSwapTokenPair(ctx, k.Keeper, base1Symbol, quoteSymbol)
	err = k.satisfyWhiteListAdmittance(ctx, pool3)
	require.Nil(t, err)
}

func TestReadWriteWhiteList(t *testing.T) {
	ctx, k := GetKeeper(t)

	poolName := "pool"
	require.False(t, k.isPoolNameExistedInWhiteList(ctx, poolName))
	k.SetWhitelist(ctx, poolName)
	require.True(t, k.isPoolNameExistedInWhiteList(ctx, poolName))
	k.DeleteWhiteList(ctx, poolName)
	require.False(t, k.isPoolNameExistedInWhiteList(ctx, poolName))
}
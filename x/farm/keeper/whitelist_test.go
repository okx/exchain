package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
	"github.com/stretchr/testify/require"
)

func TestSatisfyWhiteListAdmittance(t *testing.T) {
	ctx, k := GetKeeper(t)
	quoteSymbol := types.DefaultParams().QuoteSymbol

	baseSymbol := "xxb"
	pool1 := types.FarmPool{
		Name:            "pool1",
		MinLockedAmount: sdk.NewDecCoin(baseSymbol, sdk.ZeroInt()),
	}
	k.SetFarmPool(ctx, pool1)
	err := k.satisfyWhiteListAdmittance(ctx, pool1.MinLockedAmount.Denom)
	require.NotNil(t, err)
	require.Equal(t, types.CodeTokenNotExist, err.Code())

	// test add LPT
	pool2 := types.FarmPool{
		Name:            "pool2",
		MinLockedAmount: sdk.NewDecCoin("ammswap_" + baseSymbol + "_" + quoteSymbol, sdk.ZeroInt()),
	}
	k.SetFarmPool(ctx, pool2)
	err = k.satisfyWhiteListAdmittance(ctx, pool2.MinLockedAmount.Denom)
	require.NotNil(t, err)
	require.Equal(t, types.CodeTokenNotExist, err.Code())

	base1Symbol := "okb"
	pool3 := types.FarmPool{
		Name:            "pool3",
		MinLockedAmount:  sdk.NewDecCoin("ammswap_" + baseSymbol + "_" + base1Symbol, sdk.ZeroInt()),
	}
	err = k.satisfyWhiteListAdmittance(ctx, pool3.MinLockedAmount.Denom)
	require.NotNil(t, err)
	require.Equal(t, types.CodeTokenNotExist, err.Code())

	SetSwapTokenPair(ctx, k.Keeper, baseSymbol, quoteSymbol)

	err = k.satisfyWhiteListAdmittance(ctx, pool1.MinLockedAmount.Denom)
	require.Nil(t, err)

	err = k.satisfyWhiteListAdmittance(ctx, pool2.MinLockedAmount.Denom)
	require.Nil(t, err)

	err = k.satisfyWhiteListAdmittance(ctx, pool3.MinLockedAmount.Denom)
	require.NotNil(t, err)
	require.Equal(t, types.CodeTokenNotExist, err.Code())

	SetSwapTokenPair(ctx, k.Keeper, base1Symbol, quoteSymbol)
	err = k.satisfyWhiteListAdmittance(ctx, pool3.MinLockedAmount.Denom)
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
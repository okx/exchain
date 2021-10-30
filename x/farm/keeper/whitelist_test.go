package keeper

import (
	"testing"

	swaptypes "github.com/okex/exchain/x/ammswap/types"
	"github.com/okex/exchain/x/farm/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestSatisfyWhiteListAdmittance(t *testing.T) {
	ctx, k := GetKeeper(t)
	quoteSymbol := types.DefaultParams().QuoteSymbol

	baseSymbol := "xxb"
	pool1 := types.FarmPool{
		Name:          "pool1",
		MinLockAmount: sdk.NewDecCoinFromDec(baseSymbol, sdk.ZeroDec()),
	}
	k.SetFarmPool(ctx, pool1)
	err := k.satisfyWhiteListAdmittance(ctx, pool1)
	require.NotNil(t, err)

	// test add LPT
	lockSymbol := swaptypes.PoolTokenPrefix + baseSymbol + "_" + quoteSymbol
	pool2 := types.FarmPool{
		Name:          "pool2",
		MinLockAmount: sdk.NewDecCoinFromDec(lockSymbol, sdk.ZeroDec()),
	}
	k.SetFarmPool(ctx, pool2)
	err = k.satisfyWhiteListAdmittance(ctx, pool2)
	require.NotNil(t, err)

	base1Symbol := "okb"
	lockSymbol = swaptypes.PoolTokenPrefix + baseSymbol + "_" + base1Symbol
	pool3 := types.FarmPool{
		Name:          "pool3",
		MinLockAmount: sdk.NewDecCoinFromDec(lockSymbol, sdk.ZeroDec()),
	}
	err = k.satisfyWhiteListAdmittance(ctx, pool3)
	require.NotNil(t, err)

	SetSwapTokenPair(ctx, k.Keeper, baseSymbol, quoteSymbol)

	err = k.satisfyWhiteListAdmittance(ctx, pool1)
	require.Nil(t, err)

	err = k.satisfyWhiteListAdmittance(ctx, pool2)
	require.Nil(t, err)

	err = k.satisfyWhiteListAdmittance(ctx, pool3)
	require.NotNil(t, err)

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

package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/exchain/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestAddShares(t *testing.T) {
	ctx, _, keeper := CreateTestInput(t, false, 0)

	shares := []sdk.Dec{
		sdk.OneDec(),
		sdk.OneDec(),
		sdk.OneDec().MulInt64(100),
	}

	keeper.Keeper.SetShares(ctx, addrDels[0], addrVals[0], shares[0])
	keeper.Keeper.SetShares(ctx, addrDels[1], addrVals[0], shares[1])
	keeper.Keeper.SetShares(ctx, addrDels[2], addrVals[0], shares[2])

	var sharesExportedSlice []types.SharesExported
	keeper.IterateShares(ctx,
		func(_ int64, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares types.Shares) (stop bool) {
			sharesExportedSlice = append(sharesExportedSlice, types.NewSharesExported(delAddr, valAddr, shares))
			return false
		})

	require.True(t, len(sharesExportedSlice) == 3)

	shares1, found := keeper.GetShares(ctx, addrDels[0], addrVals[0])
	require.True(t, found)
	require.True(t, shares1.Equal(sdk.OneDec()))

	shares2, found := keeper.GetShares(ctx, addrDels[0], addrVals[2])
	require.False(t, found)
	require.True(t, shares2.IsNil())

	keeper.DeleteShares(ctx, addrVals[0], addrDels[0])
	keeper.DeleteShares(ctx, addrVals[1], addrDels[0])
	_, found = keeper.GetShares(ctx, addrDels[0], addrVals[0])
	require.False(t, found)
}

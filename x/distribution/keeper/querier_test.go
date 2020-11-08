package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/distribution/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQueryParams(t *testing.T) {
	ctx, _, k, _, _ := CreateTestInputDefault(t, false, 1000)
	querior := NewQuerier(k)
	commnutyTax, err := querior(ctx, []string{types.QueryParams, types.ParamCommunityTax}, abci.RequestQuery{})
	require.NoError(t, err)

	var taxData sdk.Dec
	_ = amino.UnmarshalJSON(commnutyTax, &taxData)
	require.Equal(t, sdk.NewDecWithPrec(2,2), taxData)

	enabled, err := querior(ctx, []string{types.QueryParams, types.ParamWithdrawAddrEnabled}, abci.RequestQuery{})
	require.True(t, err == nil)
	var enableData bool
	err1 := amino.UnmarshalJSON(enabled, &enableData)
	require.NoError(t, err1)
	require.Equal(t, true, enableData)

	_, err = querior(ctx, []string{"unknown"}, abci.RequestQuery{})
	require.Error(t, err)
	_, err = querior(ctx, []string{types.QueryParams, "unknown"}, abci.RequestQuery{})
	require.Error(t, err)
}

func TestQueryValidatorCommission(t *testing.T) {
	ctx, _, k, _, _ := CreateTestInputDefault(t, false, 1000)
	querior := NewQuerier(k)
	k.SetValidatorAccumulatedCommission(ctx,valOpAddr1, NewTestDecCoins(15,1))

	bz,err := amino.MarshalJSON(types.NewQueryValidatorCommissionParams(valOpAddr1))
	require.NoError(t,err)
	commission, err := querior(ctx, []string{types.QueryValidatorCommission}, abci.RequestQuery{Data: bz})
	require.NoError(t, err)

	var data sdk.SysCoins
	err = amino.UnmarshalJSON(commission, &data)
	require.NoError(t, err)
	require.Equal(t, NewTestDecCoins(15,1), data)
}

func TestQueryDelegatorWithdrawAddress(t *testing.T) {
	ctx, _, k, _, _ := CreateTestInputDefault(t, false, 1000)
	querior := NewQuerier(k)
	require.NoError(t,k.SetWithdrawAddr(ctx, valAccAddr1, valAccAddr2))

	bz, err := amino.MarshalJSON(types.NewQueryDelegatorWithdrawAddrParams(valAccAddr1))
	require.NoError(t,err)
	withdrawAddr, err := querior(ctx, []string{types.QueryWithdrawAddr}, abci.RequestQuery{Data: bz})
	require.NoError(t, err)

	var data sdk.AccAddress
	err = amino.UnmarshalJSON(withdrawAddr, &data)
	require.NoError(t, err)
	require.Equal(t, valAccAddr2, data)
}

func TestQueryCommunityPool(t *testing.T) {
	ctx, _, k, _, _ := CreateTestInputDefault(t, false, 1000)
	querior := NewQuerier(k)
	feePool := k.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(NewTestDecCoins(123,2))
	k.SetFeePool(ctx, feePool)

	communityPool, err := querior(ctx, []string{types.QueryCommunityPool}, abci.RequestQuery{})
	require.NoError(t, err)

	var data sdk.SysCoins
	err1 := amino.UnmarshalJSON(communityPool, &data)
	require.NoError(t, err1)
	require.Equal(t, NewTestDecCoins(123,2), data)
}
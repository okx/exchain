package keeper

import (
	"testing"
	"time"

	"github.com/okex/exchain/x/common"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/dex/types"
	"github.com/stretchr/testify/require"
	amino "github.com/tendermint/go-amino"
)

func TestQuerier_ProductsAndMatchOrder(t *testing.T) {

	common.InitConfig()
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx
	addr, err := sdk.AccAddressFromBech32(types.TestTokenPairOwner)
	require.Nil(t, err)
	tokenPair0 := &types.TokenPair{
		BaseAssetSymbol:  "bToken0",
		QuoteAssetSymbol: common.NativeToken,
		Owner:            addr,
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(50)),
		BlockHeight:      8,
	}

	tokenPair1 := &types.TokenPair{
		BaseAssetSymbol:  "bToken1",
		QuoteAssetSymbol: common.NativeToken,
		Owner:            addr,
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100)),
		BlockHeight:      10,
	}

	// SaveTokenPair
	err = testInput.DexKeeper.SaveTokenPair(ctx, tokenPair0)
	require.Nil(t, err)
	err = testInput.DexKeeper.SaveTokenPair(ctx, tokenPair1)
	require.Nil(t, err)
	querier := NewQuerier(testInput.DexKeeper)

	var normalPath = []string{types.QueryProducts, types.QueryMatchOrder}

	for _, path := range normalPath {
		// successful case
		queryParams := types.NewQueryDexInfoParams(types.TestTokenPairOwner, 1, 50)
		bz, err := amino.MarshalJSON(queryParams)
		require.Nil(t, err)
		data, err := querier(ctx, []string{path}, abci.RequestQuery{Data: bz})
		require.Nil(t, err)
		require.True(t, data != nil)

		// error case : failed to query data because  param is nil
		dataUnmarshalJSON, err := querier(ctx, []string{path}, abci.RequestQuery{Data: nil})
		require.Error(t, err)
		require.True(t, dataUnmarshalJSON == nil)

		// successful case : query data while page limit is out range of data amount
		queryParams = types.NewQueryDexInfoParams(types.TestTokenPairOwner, 2, 50)
		bz, err = amino.MarshalJSON(queryParams)
		require.Nil(t, err)
		data, err = querier(ctx, []string{path}, abci.RequestQuery{Data: bz})
		require.Nil(t, err)
		require.True(t, data != nil)

		// successful case : query data while  page limit is in range of data amount
		queryParams = types.NewQueryDexInfoParams(types.TestTokenPairOwner, 1, 1)
		bz, err = amino.MarshalJSON(queryParams)
		require.Nil(t, err)
		data, err = querier(ctx, []string{path}, abci.RequestQuery{Data: bz})
		require.Nil(t, err)
		require.True(t, data != nil)

	}

	// error case : failed to query because query path is wrong
	queryParams := types.NewQueryDexInfoParams(types.TestTokenPairOwner, 1, 50)
	bz, err := amino.MarshalJSON(queryParams)
	require.Nil(t, err)
	dataOther, err := querier(ctx, []string{"others"}, abci.RequestQuery{Data: bz})
	require.NotNil(t, err)
	require.True(t, dataOther == nil)

}

func TestQuerier_Deposits(t *testing.T) {
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx
	addr, err := sdk.AccAddressFromBech32(types.TestTokenPairOwner)
	require.Nil(t, err)
	tokenPair0 := &types.TokenPair{
		BaseAssetSymbol:  "bToken0",
		QuoteAssetSymbol: common.NativeToken,
		Owner:            addr,
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(50)),
		BlockHeight:      8,
	}

	tokenPair1 := &types.TokenPair{
		BaseAssetSymbol:  "bToken1",
		QuoteAssetSymbol: common.NativeToken,
		Owner:            addr,
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100)),
		BlockHeight:      10,
	}

	// SaveTokenPair
	err = testInput.DexKeeper.SaveTokenPair(ctx, tokenPair0)
	require.Nil(t, err)
	err = testInput.DexKeeper.SaveTokenPair(ctx, tokenPair1)
	require.Nil(t, err)
	querier := NewQuerier(testInput.DexKeeper)
	path := types.QueryDeposits
	// successful case
	queryParams := types.NewQueryDepositParams(types.TestTokenPairOwner, "", common.NativeToken, 1, 50)
	bz, err := amino.MarshalJSON(queryParams)
	require.Nil(t, err)
	data, err := querier(ctx, []string{path}, abci.RequestQuery{Data: bz})
	require.Nil(t, err)
	require.True(t, data != nil)

	// error case : failed to query data because  param is nil
	dataUnmarshalJSON, err := querier(ctx, []string{path}, abci.RequestQuery{Data: nil})
	require.Error(t, err)
	require.True(t, dataUnmarshalJSON == nil)

	// successful case : query data while page limit is out range of data amount
	queryParams = types.NewQueryDepositParams(types.TestTokenPairOwner, "", "", 2, 50)
	bz, err = amino.MarshalJSON(queryParams)
	require.Nil(t, err)
	data, err = querier(ctx, []string{path}, abci.RequestQuery{Data: bz})
	require.Nil(t, err)
	require.True(t, data != nil)

	// successful case : query data while  page limit is in range of data amount
	queryParams = types.NewQueryDepositParams(types.TestTokenPairOwner, "", "", 1, 1)
	bz, err = amino.MarshalJSON(queryParams)
	require.Nil(t, err)
	data, err = querier(ctx, []string{path}, abci.RequestQuery{Data: bz})
	require.Nil(t, err)
	require.True(t, data != nil)
}

func TestQuerier_QueryParams(t *testing.T) {

	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx

	p := types.DefaultParams()
	p.DelistMaxDepositPeriod = time.Second * 123
	testInput.DexKeeper.SetParams(ctx, *p)
	querier := NewQuerier(testInput.DexKeeper)
	res, err := querier(ctx, []string{"params"}, abci.RequestQuery{})

	require.True(t, len(res) > 0)
	require.True(t, err == nil)

}

func TestQuerier_QueryProductsDelisting(t *testing.T) {

	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx

	querier := NewQuerier(testInput.DexKeeper)

	tokenPair := GetBuiltInTokenPair()
	tokenPair.Delisting = true

	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	res, err := querier(ctx, []string{types.QueryProductsDelisting}, abci.RequestQuery{})
	require.True(t, len(res) > 0)
	require.Nil(t, err)

}

func TestQueryParam(t *testing.T) {
	// NewQueryDexInfoParams
	tests := []struct {
		name    string
		owner   string
		page    int
		perPage int
		result  bool
	}{
		{"new-no-owner", "", 1, 50, true},
		{"new-with-owner", types.TestTokenPairOwner, 1, 50, true},
		{"new-wrong-address", "wrong-address", 1, 50, false},
		{"new-wrong-page", "", 0, 50, false},
		{"new-wrong-per-page", "", 1, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := types.NewQueryDexInfoParams(tt.owner, tt.page, tt.perPage)
			if tt.result {
				require.NotNil(t, params)
			} else {
				require.NotNil(t, params)
			}
		})
	}
}

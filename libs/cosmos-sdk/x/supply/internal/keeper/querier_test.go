package keeper_test

import (
	"fmt"
	"testing"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/stretchr/testify/require"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	keep "github.com/okex/exchain/libs/cosmos-sdk/x/supply/internal/keeper"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/internal/types"
)

func TestNewQuerier(t *testing.T) {
	app, ctx := createTestApp(false)
	keeper := app.SupplyKeeper
	cdc := app.Codec()

	supplyCoins := sdk.NewCoins(
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100)),
		sdk.NewCoin("photon", sdk.NewInt(50)),
		sdk.NewCoin("atom", sdk.NewInt(2000)),
		sdk.NewCoin("btc", sdk.NewInt(21000000)),
	)

	keeper.SetSupply(ctx, types.NewSupply(supplyCoins))

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	querier := keep.NewQuerier(keeper)

	bz, err := querier(ctx, []string{"other"}, query)
	require.Error(t, err)
	require.Nil(t, bz)

	queryTotalSupplyParams := types.NewQueryTotalSupplyParams(1, 20)
	bz, errRes := cdc.MarshalJSON(queryTotalSupplyParams)
	require.Nil(t, errRes)

	query.Path = fmt.Sprintf("/custom/supply/%s", types.QueryTotalSupply)
	query.Data = bz

	_, err = querier(ctx, []string{types.QueryTotalSupply}, query)
	require.Nil(t, err)

	querySupplyParams := types.NewQuerySupplyOfParams(sdk.DefaultBondDenom)
	bz, errRes = cdc.MarshalJSON(querySupplyParams)
	require.Nil(t, errRes)

	query.Path = fmt.Sprintf("/custom/supply/%s", types.QuerySupplyOf)
	query.Data = bz

	_, err = querier(ctx, []string{types.QuerySupplyOf}, query)
	require.Nil(t, err)
}

func TestQuerySupply(t *testing.T) {
	app, ctx := createTestApp(false)
	keeper := app.SupplyKeeper
	cdc := app.Codec()

	supplyCoins := sdk.NewCoins(
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100)),
		sdk.NewCoin("photon", sdk.NewInt(50)),
		sdk.NewCoin("atom", sdk.NewInt(2000)),
		sdk.NewCoin("btc", sdk.NewInt(21000000)),
	)

	querier := keep.NewQuerier(keeper)

	keeper.SetSupply(ctx, types.NewSupply(supplyCoins))

	queryTotalSupplyParams := types.NewQueryTotalSupplyParams(1, 10)
	bz, errRes := cdc.MarshalJSON(queryTotalSupplyParams)
	require.Nil(t, errRes)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	query.Path = fmt.Sprintf("/custom/supply/%s", types.QueryTotalSupply)
	query.Data = bz

	res, err := querier(ctx, []string{types.QueryTotalSupply}, query)
	require.Nil(t, err)

	var totalCoins sdk.Coins
	errRes = cdc.UnmarshalJSON(res, &totalCoins)
	require.Nil(t, errRes)
	require.Equal(t, supplyCoins, totalCoins)

	querySupplyParams := types.NewQuerySupplyOfParams(sdk.DefaultBondDenom)
	bz, errRes = cdc.MarshalJSON(querySupplyParams)
	require.Nil(t, errRes)

	query.Path = fmt.Sprintf("/custom/supply/%s", types.QuerySupplyOf)
	query.Data = bz

	res, err = querier(ctx, []string{types.QuerySupplyOf}, query)
	require.Nil(t, err)

	var supply sdk.Dec
	errRes = supply.UnmarshalJSON(res)
	require.Nil(t, errRes)
	require.True(sdk.DecEq(t, sdk.NewDec(100), supply))

}

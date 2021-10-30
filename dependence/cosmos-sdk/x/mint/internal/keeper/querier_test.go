package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keep "github.com/okex/exchain/dependence/cosmos-sdk/x/mint/internal/keeper"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/mint/internal/types"

	abci "github.com/okex/exchain/dependence/tendermint/abci/types"
)

func TestNewQuerier(t *testing.T) {
	app, ctx := createTestApp(true)
	querier := keep.NewQuerier(app.MintKeeper)

	query := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	_, err := querier(ctx, []string{types.QueryParameters}, query)
	require.NoError(t, err)

	_, err = querier(ctx, []string{"foo"}, query)
	require.Error(t, err)
}

func TestQueryParams(t *testing.T) {
	app, ctx := createTestApp(true)
	querier := keep.NewQuerier(app.MintKeeper)

	var params types.Params

	res, sdkErr := querier(ctx, []string{types.QueryParameters}, abci.RequestQuery{})
	require.NoError(t, sdkErr)

	err := app.Codec().UnmarshalJSON(res, &params)
	require.NoError(t, err)

	expected := app.MintKeeper.GetParams(ctx)
	require.Equal(t, expected.MintDenom, params.MintDenom)
	require.Equal(t, expected.BlocksPerYear, params.BlocksPerYear)
	require.Equal(t, expected.DeflationRate, params.DeflationRate)
	require.Equal(t, expected.DeflationEpoch, params.DeflationEpoch)
	require.Equal(t, expected.FarmProportion, params.FarmProportion)
}

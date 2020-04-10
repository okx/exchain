package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestQuerier(t *testing.T) {
	ctx, keeper := testPrepare(t)
	requestQueryDemo := abci.RequestQuery{}
	// no app upgrade config
	_, err := queryUpgradeConfig(ctx, requestQueryDemo, keeper)
	require.Error(t, err)

	// set app upgrade config
	require.NoError(t, keeper.SetAppUpgradeConfig(ctx, 1, 1, 1024, "software1"))

	queryFunc := NewQuerier(keeper)
	_, err = queryFunc(ctx, []string{QueryUpgradeConfig}, requestQueryDemo)
	require.NoError(t, err)
	_, err = queryFunc(ctx, []string{QueryUpgradeVersion}, requestQueryDemo)
	require.NoError(t, err)
	_, err = queryFunc(ctx, []string{QueryUpgradeFailedVersion}, requestQueryDemo)
	require.NoError(t, err)
	// get a error case

	_, err = queryFunc(ctx, []string{"helloworld"}, requestQueryDemo)
	require.Error(t, err)

}

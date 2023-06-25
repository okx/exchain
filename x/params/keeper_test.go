package params

import (
	"github.com/okex/exchain/x/params/types"
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/store"
	storetypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmdb "github.com/okex/exchain/libs/tm-db"
	"github.com/stretchr/testify/suite"
)

type KeeperSuite struct {
	suite.Suite
	ms           storetypes.CommitMultiStore
	paramsKeeper Keeper
}

func (suite *KeeperSuite) SetupTest() {
	db := tmdb.NewMemDB()
	storeKey := sdk.NewKVStoreKey(StoreKey)
	tstoreKey := sdk.NewTransientStoreKey(TStoreKey)

	suite.ms = store.NewCommitMultiStore(tmdb.NewMemDB())
	suite.ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	suite.ms.MountStoreWithDB(tstoreKey, sdk.StoreTypeTransient, db)
	err := suite.ms.LoadLatestVersion()
	suite.NoError(err)

	suite.paramsKeeper = NewKeeper(ModuleCdc, storeKey, tstoreKey, log.NewNopLogger())
}

func (suite *KeeperSuite) Context(height int64) sdk.Context {
	return sdk.NewContext(suite.ms, abci.Header{Height: height}, false, log.NewNopLogger())
}

func TestKeeper(t *testing.T) {
	suite.Run(t, new(KeeperSuite))
}

func (suite *KeeperSuite) TestGetGasConfig() {
	sub := "params"
	tests := []struct {
		changes []types.ParamChange
		fncheck func(res storetypes.GasConfig)
	}{
		{
			changes: []types.ParamChange{{Subspace: sub, Key: storetypes.GasHasDesc, Value: "\"100\""}},
			fncheck: func(res storetypes.GasConfig) {
				gs := storetypes.GetDefaultGasConfig()
				gs.HasCost = 100
				suite.Equal(*gs, res)
			},
		},
		{
			changes: []types.ParamChange{{Subspace: sub, Key: storetypes.GasDeleteDesc, Value: "\"10\""}},
			fncheck: func(res storetypes.GasConfig) {
				gs := storetypes.GetDefaultGasConfig()
				gs.HasCost = 100
				gs.DeleteCost = 10
				suite.Equal(*gs, res)
			},
		},
		{
			changes: []types.ParamChange{{Subspace: sub, Key: storetypes.GasReadCostFlatDesc, Value: "\"10\""}},
			fncheck: func(res storetypes.GasConfig) {
				gs := storetypes.GetDefaultGasConfig()
				gs.HasCost = 100
				gs.DeleteCost = 10
				gs.ReadCostFlat = 10
				suite.Equal(*gs, res)
			},
		},
		{
			changes: []types.ParamChange{{Subspace: sub, Key: storetypes.GasReadPerByteDesc, Value: "\"10\""}},
			fncheck: func(res storetypes.GasConfig) {
				gs := storetypes.GetDefaultGasConfig()
				gs.HasCost = 100
				gs.DeleteCost = 10
				gs.ReadCostFlat = 10
				gs.ReadCostPerByte = 10
				suite.Equal(*gs, res)
			},
		},
		{
			changes: []types.ParamChange{{Subspace: sub, Key: storetypes.GasWriteCostFlatDesc, Value: "\"10\""}},
			fncheck: func(res storetypes.GasConfig) {
				gs := storetypes.GetDefaultGasConfig()
				gs.HasCost = 100
				gs.DeleteCost = 10
				gs.ReadCostFlat = 10
				gs.ReadCostPerByte = 10
				gs.WriteCostFlat = 10
				suite.Equal(*gs, res)
			},
		},
		{
			changes: []types.ParamChange{{Subspace: sub, Key: storetypes.GasWritePerByteDesc, Value: "\"10\""}},
			fncheck: func(res storetypes.GasConfig) {
				gs := storetypes.GetDefaultGasConfig()
				gs.HasCost = 100
				gs.DeleteCost = 10
				gs.ReadCostFlat = 10
				gs.ReadCostPerByte = 10
				gs.WriteCostFlat = 10
				gs.WriteCostPerByte = 10
				suite.Equal(*gs, res)
			},
		},
		{
			changes: []types.ParamChange{{Subspace: sub, Key: storetypes.GasIterNextCostFlatDesc, Value: "\"10\""}},
			fncheck: func(res storetypes.GasConfig) {
				gs := storetypes.GetDefaultGasConfig()
				gs.HasCost = 100
				gs.DeleteCost = 10
				gs.ReadCostFlat = 10
				gs.ReadCostPerByte = 10
				gs.WriteCostFlat = 10
				gs.WriteCostPerByte = 10
				gs.IterNextCostFlat = 10
				suite.Equal(*gs, res)
			},
		},
	}

	for _, tt := range tests {
		ctx := suite.Context(0)

		changeParams(ctx, &suite.paramsKeeper, types.NewParameterChangeProposal("hello", "word", tt.changes, 1))

		res := suite.paramsKeeper.GetGasConfig(ctx)
		tt.fncheck(*res)
	}
}

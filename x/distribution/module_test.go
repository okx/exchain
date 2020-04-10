package distribution

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okchain/x/distribution/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAppModule(t *testing.T) {
	//valOpAddrs, valConsPks, valConsAddrs := keeper.GetAddrs()
	ctx, _, k, _, supplyKeeper := CreateTestInputDefault(t, false, 1000)

	module := NewAppModule(k, supplyKeeper)
	require.EqualValues(t, ModuleName, module.AppModuleBasic.Name())
	require.EqualValues(t, ModuleName, module.Name())
	require.EqualValues(t, RouterKey, module.Route())
	require.EqualValues(t, QuerierRoute, module.QuerierRoute())

	cdc := codec.New()
	module.RegisterCodec(cdc)

	msg := module.DefaultGenesis()
	require.Nil(t, module.ValidateGenesis(msg))
	require.NotNil(t, module.ValidateGenesis([]byte{}))

	module.InitGenesis(ctx, msg)
	exportMsg := module.ExportGenesis(ctx)

	var gs GenesisState
	types.ModuleCdc.MustUnmarshalJSON(exportMsg, &gs)
	require.Equal(t, msg, exportMsg)

	// for coverage
	module.BeginBlock(ctx, abci.RequestBeginBlock{})
	module.EndBlock(ctx, abci.RequestEndBlock{})
	module.GetQueryCmd(cdc)
	module.GetTxCmd(cdc)
	module.NewQuerierHandler()
	module.NewHandler()
}

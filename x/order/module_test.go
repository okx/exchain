package order

import (
	"testing"

	"github.com/okex/okchain/x/common/version"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okchain/x/order/keeper"
	"github.com/okex/okchain/x/order/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAppModule(t *testing.T) {
	testInput := keeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	module := NewAppModule(version.CurrentProtocolVersion, testInput.OrderKeeper, testInput.SupplyKeeper)

	require.EqualValues(t, ModuleName, module.Name())
	require.EqualValues(t, RouterKey, module.Route())
	require.EqualValues(t, QuerierRoute, module.QuerierRoute())

	cdc := codec.New()
	module.RegisterCodec(cdc)

	msg := module.DefaultGenesis()
	require.Nil(t, module.ValidateGenesis(msg))
	require.NotNil(t, module.ValidateGenesis([]byte{}))

	module.InitGenesis(ctx, msg)
	params := keeper.GetParams(ctx)
	require.EqualValues(t, 259200, params.OrderExpireBlocks)
	exportMsg := module.ExportGenesis(ctx)

	var gs GenesisState
	types.ModuleCdc.MustUnmarshalJSON(exportMsg, &gs)
	require.EqualValues(t, msg, types.ModuleCdc.MustMarshalJSON(gs))

	// for coverage
	module.BeginBlock(ctx, abci.RequestBeginBlock{})
	module.EndBlock(ctx, abci.RequestEndBlock{})
	module.GetQueryCmd(cdc)
	module.GetTxCmd(cdc)
	module.NewQuerierHandler()
	module.NewHandler()
}

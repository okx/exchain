package farm

import (
	"testing"

	cliLcd "github.com/okex/exchain/libs/cosmos-sdk/client/lcd"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/farm/keeper"
	"github.com/okex/exchain/x/farm/types"
	"github.com/stretchr/testify/require"
)

func TestAppModule(t *testing.T) {
	ctx, mk := keeper.GetKeeper(t)
	keeper := mk.Keeper
	var blockHeight int64 = 10
	ctx = ctx.WithBlockHeight(blockHeight)
	BeginBlocker(ctx, abci.RequestBeginBlock{Header: abci.Header{Height: blockHeight}}, keeper)
	module := NewAppModule(keeper)

	require.EqualValues(t, ModuleName, module.Name())
	require.EqualValues(t, RouterKey, module.Route())

	cdc := types.ModuleCdc
	//module.RegisterCodec(cdc)

	msg := module.DefaultGenesis()
	require.Nil(t, module.ValidateGenesis(msg))
	require.NotNil(t, module.ValidateGenesis([]byte{}))

	module.InitGenesis(ctx, msg)
	params := keeper.GetParams(ctx)
	require.EqualValues(t, types.DefaultParams().String(), params.String())
	exportMsg := module.ExportGenesis(ctx)

	var genesis, exportedGenesis types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(exportMsg, &exportedGenesis)
	types.ModuleCdc.MustUnmarshalJSON(msg, &genesis)
	//gs.Pools = types.FarmPools{}
	require.EqualValues(t, genesis, exportedGenesis)

	// for coverage
	module.BeginBlock(ctx, abci.RequestBeginBlock{})
	module.EndBlock(ctx, abci.RequestEndBlock{})
	module.GetQueryCmd(cdc)
	module.GetTxCmd(cdc)
	module.NewQuerierHandler()
	module.NewHandler()
	rs := cliLcd.NewRestServer(cdc, nil)
	module.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
	//module.RegisterInvariants(nil)
	module.RegisterCodec(codec.New())
	module.QuerierRoute()
	module.Name()
}

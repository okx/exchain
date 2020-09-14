package ammswap

import (
	"encoding/json"
	"testing"

	cliLcd "github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAppModule(t *testing.T) {
	mapp, _ := getMockApp(t, 1)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	module := NewAppModule(keeper)

	require.EqualValues(t, ModuleName, module.Name())
	require.EqualValues(t, RouterKey, module.Route())
	require.EqualValues(t, QuerierRoute, module.QuerierRoute())

	cdc := ModuleCdc
	//module.RegisterCodec(cdc)

	msg := module.DefaultGenesis()
	require.Nil(t, module.ValidateGenesis(msg))
	require.NotNil(t, module.ValidateGenesis([]byte{}))

	module.InitGenesis(ctx, msg)
	params := keeper.GetParams(ctx)
	require.EqualValues(t, types.DefaultParams().FeeRate, params.FeeRate)
	exportMsg := module.ExportGenesis(ctx)

	var gs GenesisState
	ModuleCdc.MustUnmarshalJSON(exportMsg, &gs)
	require.EqualValues(t, msg, json.RawMessage(ModuleCdc.MustMarshalJSON(gs)))

	// for coverage
	module.BeginBlock(ctx, abci.RequestBeginBlock{})
	module.EndBlock(ctx, abci.RequestEndBlock{})
	module.GetQueryCmd(cdc)
	module.GetTxCmd(cdc)
	module.NewQuerierHandler()
	module.NewHandler()
	rs := cliLcd.NewRestServer(cdc, nil)
	module.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
	module.RegisterInvariants(nil)
	module.RegisterCodec(codec.New())
}

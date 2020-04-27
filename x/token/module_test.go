package token

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"

	cliLcd "github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/okex/okchain/x/common/version"
	"github.com/okex/okchain/x/token/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAppModule_InitGenesis(t *testing.T) {
	app, tokenKeeper, _ := getMockDexAppEx(t, 0)
	module := NewAppModule(version.ProtocolVersionV0, tokenKeeper, app.supplyKeeper)
	ctx := app.NewContext(true, abci.Header{})
	gs := defaultGenesisState()
	gs.Tokens = nil
	gsJSON := types.ModuleCdc.MustMarshalJSON(gs)

	err := module.ValidateGenesis(gsJSON)
	require.NoError(t, err)

	vu := module.InitGenesis(ctx, gsJSON)
	params := tokenKeeper.GetParams(ctx)
	require.Equal(t, gs.Params, params)
	require.Equal(t, vu, []abci.ValidatorUpdate{})

	export := module.ExportGenesis(ctx)
	require.EqualValues(t, gsJSON, []byte(export))

	require.EqualValues(t, types.ModuleName, module.Name())
	require.EqualValues(t, types.ModuleName, module.AppModuleBasic.Name())
	require.EqualValues(t, types.RouterKey, module.Route())
	require.EqualValues(t, types.QuerierRoute, module.QuerierRoute())
	module.NewHandler()
	module.GetQueryCmd(app.Cdc)
	module.GetTxCmd(app.Cdc)
	module.NewQuerierHandler()
	rs := cliLcd.NewRestServer(app.Cdc, nil)
	module.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
	module.BeginBlock(ctx, abci.RequestBeginBlock{})
	module.EndBlock(ctx, abci.RequestEndBlock{})
	module.DefaultGenesis()
	module.RegisterCodec(codec.New())

	gsJSON = []byte("[[],{}]")
	err = module.ValidateGenesis(gsJSON)
	require.NotNil(t, err)
}

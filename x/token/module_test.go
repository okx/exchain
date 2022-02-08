package token

import (
	"github.com/okex/exchain/x/common"
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"

	cliLcd "github.com/okex/exchain/libs/cosmos-sdk/client/lcd"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/common/version"
	"github.com/okex/exchain/x/token/types"
	"github.com/stretchr/testify/require"
)

func TestAppModule_InitGenesis(t *testing.T) {
	common.InitConfig()
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

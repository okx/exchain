package gov

import (
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	cliLcd "github.com/okex/exchain/libs/cosmos-sdk/client/lcd"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/x/gov/client"
	"github.com/okex/exchain/x/gov/client/rest"
	"github.com/okex/exchain/x/gov/keeper"
	"github.com/okex/exchain/x/gov/types"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

func TestAppModule_BeginBlock(t *testing.T) {

}

func getCmdSubmitProposal(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{}
}

func proposalRESTHandler(cliCtx context.CLIContext) rest.ProposalRESTHandler {
	return rest.ProposalRESTHandler{}
}

func TestNewAppModuleBasic(t *testing.T) {
	ctx, _, gk, _, crisisKeeper := keeper.CreateTestInput(t, false, 1000)

	moduleBasic := NewAppModuleBasic(client.ProposalHandler{
		CLIHandler:  getCmdSubmitProposal,
		RESTHandler: proposalRESTHandler,
	})

	require.Equal(t, types.ModuleName, moduleBasic.Name())

	cdc := codec.New()
	moduleBasic.RegisterCodec(cdc)
	bz, err := cdc.MarshalBinaryBare(types.MsgSubmitProposal{})
	require.NotNil(t, bz)
	require.Nil(t, err)

	jsonMsg := moduleBasic.DefaultGenesis()
	err = moduleBasic.ValidateGenesis(jsonMsg)
	require.Nil(t, err)
	err = moduleBasic.ValidateGenesis(jsonMsg[:len(jsonMsg)-1])
	require.NotNil(t, err)

	rs := cliLcd.NewRestServer(cdc, nil)
	moduleBasic.RegisterRESTRoutes(rs.CliCtx, rs.Mux)

	// todo: check diff after GetTxCmd
	moduleBasic.GetTxCmd(cdc)

	// todo: check diff after GetQueryCmd
	moduleBasic.GetQueryCmd(cdc)

	appModule := NewAppModule(gk, gk.SupplyKeeper())
	require.Equal(t, types.ModuleName, appModule.Name())

	// todo: check diff after RegisterInvariants
	appModule.RegisterInvariants(&crisisKeeper)

	require.Equal(t, types.RouterKey, appModule.Route())

	require.IsType(t, NewHandler(gk), appModule.NewHandler())

	require.Equal(t, types.QuerierRoute, appModule.QuerierRoute())

	require.IsType(t, NewQuerier(gk), appModule.NewQuerierHandler())

	require.Equal(t, []abci.ValidatorUpdate{}, appModule.InitGenesis(ctx, jsonMsg))

	require.Equal(t, jsonMsg, appModule.ExportGenesis(ctx))

	appModule.BeginBlock(ctx, abci.RequestBeginBlock{})

	appModule.EndBlock(ctx, abci.RequestEndBlock{})
}

package distribution

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okexchain/x/distribution/keeper"
	"github.com/okex/okexchain/x/distribution/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestA(t *testing.T){
	_,err:=sdk.AccAddressFromBech32("okexchain1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk")
	fmt.Println(err)
}

func TestAppModule(t *testing.T) {
	ctx, _, k, _, supplyKeeper := keeper.CreateTestInputDefault(t, false, 1000)

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
	require.NotPanics(t, func() {
		types.ModuleCdc.MustUnmarshalJSON(exportMsg, &gs)
	})

	// for coverage
	module.BeginBlock(ctx, abci.RequestBeginBlock{})
	module.EndBlock(ctx, abci.RequestEndBlock{})
	module.GetQueryCmd(cdc)
	module.GetTxCmd(cdc)
	module.NewQuerierHandler()
	module.NewHandler()
}

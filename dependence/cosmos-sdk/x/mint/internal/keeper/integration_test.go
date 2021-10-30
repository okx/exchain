package keeper_test

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/exchain/dependence/cosmos-sdk/simapp"
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/mint/internal/types"
)

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.MintKeeper.SetParams(ctx, types.DefaultParams())
	app.MintKeeper.SetMinter(ctx, types.InitialMinterCustom())

	return app, ctx
}

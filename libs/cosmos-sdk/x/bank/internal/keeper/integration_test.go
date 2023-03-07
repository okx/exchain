package keeper_test

import (
	abci "github.com/okx/exchain/libs/tendermint/abci/types"

	"github.com/okx/exchain/libs/cosmos-sdk/simapp"
	sdk "github.com/okx/exchain/libs/cosmos-sdk/types"
	"github.com/okx/exchain/libs/cosmos-sdk/x/auth"
)

func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})

	app.AccountKeeper.SetParams(ctx, auth.DefaultParams())
	app.BankKeeper.SetSendEnabled(ctx, true)

	return app, ctx
}

package keeper_test

import (
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/okex/exchain/ibc-3rd/cosmos-v443/simapp"
	sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"
	authtypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/x/auth/types"
)

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	return app, ctx
}

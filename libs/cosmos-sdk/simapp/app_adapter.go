package simapp

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/abci/types"
)

type AppAdapter interface {
	App
	Upgrade(ctx sdk.Context, req *types.UpgradeReq) (*types.UpgradeResp, error)
}

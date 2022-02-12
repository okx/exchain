package client

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/02-client/keeper"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/exported"
)

// BeginBlocker updates an existing localhost client with the latest block height.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	_, found := k.GetClientState(ctx, exported.Localhost)
	if !found {
		return
	}

	// update the localhost client with the latest block height
	if err := k.UpdateClient(ctx, exported.Localhost, nil); err != nil {
		panic(err)
	}
}

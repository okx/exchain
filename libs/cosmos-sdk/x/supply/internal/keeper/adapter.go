package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/internal/types"
)

var (
	_ types.CM40AccountKeeper = Keeper{}
)

func (k Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return k.bk.SendCoins(ctx, fromAddr, toAddr, amt)
}

func (k Keeper) NewAccount(ctx sdk.Context, acc authtypes.Account) authtypes.Account {
	return k.ak.NewAccount(ctx, acc)
}

func (k Keeper) SetAccount(ctx sdk.Context, acc authtypes.Account) {
	k.ak.SetAccount(ctx, acc)
}

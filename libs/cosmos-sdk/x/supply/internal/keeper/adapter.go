package keeper

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

func (k Keeper) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return k.bk.SendCoins(ctx,fromAddr,toAddr,amt)
}

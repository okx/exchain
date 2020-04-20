package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

//IsTokenExits checkout token is exit
func (k Keeper) IsTokenExits(ctx sdk.Context, token string) error {
	isExist := k.tokenKeeper.TokenExist(ctx, token)
	if !isExist {
		return sdk.ErrInternal("Failed: token not exits")
	}

	return nil

}

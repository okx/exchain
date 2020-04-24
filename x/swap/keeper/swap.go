package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/swap/types"
)

//IsTokenExist checkout token is exit
func (k Keeper) IsTokenExist(ctx sdk.Context, token string) error {
	isExist := k.tokenKeeper.TokenExist(ctx, token)
	if !isExist {
		return sdk.ErrInternal("Failed: token not exits")
	}

	t := k.tokenKeeper.GetTokenInfo(ctx, token)
	if t.Type == sdk.NewInt(types.GenerateTokenType) {
		return sdk.ErrInternal("Failed to create exchange with poll token")
	}
	return nil

}

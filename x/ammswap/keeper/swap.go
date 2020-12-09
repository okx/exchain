package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/ammswap/types"
)

// IsTokenExist check token is exist
func (k Keeper) IsTokenExist(ctx sdk.Context, token string) error {
	isExist := k.tokenKeeper.TokenExist(ctx, token)
	if !isExist {
		return types.ErrTokenNotExist(types.DefaultCodespace, "Failed: token does not exist")
	}

	t := k.tokenKeeper.GetTokenInfo(ctx, token)
	if t.Type == types.GenerateTokenType {
		return types.ErrInvalidCoins(types.DefaultCodespace, "Failed to create exchange with pool token")
	}
	return nil

}

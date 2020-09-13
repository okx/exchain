package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/dex/types"
)

type StreamKeeper interface {
	OnAddNewTokenPair(ctx sdk.Context, tokenPair *types.TokenPair)
	OnTokenPairUpdated(ctx sdk.Context)
}

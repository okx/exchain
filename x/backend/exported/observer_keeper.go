package exported

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	swaptypes "github.com/okex/okexchain/x/ammswap/types"
)

type BackendKeeper interface {
	OnSwapToken(ctx sdk.Context, address sdk.AccAddress, swapTokenPair swaptypes.SwapTokenPair, sellAmount sdk.SysCoin, buyAmount sdk.SysCoin)
	OnSwapCreateExchange(ctx sdk.Context, swapTokenPair swaptypes.SwapTokenPair)
}

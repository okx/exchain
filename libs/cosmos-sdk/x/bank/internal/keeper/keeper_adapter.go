package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank/internal/types"
)

type CM40ViewKeeper interface {
	IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	BlockedAddr(addr sdk.AccAddress) bool
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

func (adapter BaseKeeper) BlockedAddr(addr sdk.AccAddress) bool {
	return adapter.BlacklistedAddr(addr)
}

func (adapter BaseKeeper) IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error {
	if !adapter.GetSendEnabled(ctx) {
		return sdkerrors.Wrapf(types.ErrSendDisabled, "transfers are currently disabled")
	}
	return nil
}

func (adapter BaseKeeper) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return adapter.GetCoins(ctx, addr)
}

func (adapter BaseKeeper) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	coins := adapter.GetCoins(ctx, addr)
	return sdk.Coin{
		Amount: coins.AmountOf(denom),
		Denom:  denom,
	}
}

func (adapter BaseKeeper) HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool {
	return adapter.GetBalance(ctx, addr, amt.Denom).IsGTE(amt)
}

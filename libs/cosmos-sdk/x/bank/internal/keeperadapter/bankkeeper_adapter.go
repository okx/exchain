package keeperadapter

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank/internal/keeper"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank/internal/types"
)

var (
	_ CM40ViewKeeper = BankKeeperAdapter{}
)

type CM40ViewKeeper interface {
	IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	BlockedAddr(addr sdk.AccAddress) bool
	HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// BankKeeperAdapter is used in wasm module
type BankKeeperAdapter struct {
	keeper.Keeper
}

func NewBankKeeperAdapter(bankKeeper keeper.Keeper) *BankKeeperAdapter {
	return &BankKeeperAdapter{Keeper: bankKeeper}
}

func (adapter BankKeeperAdapter) BlockedAddr(addr sdk.AccAddress) bool {
	return adapter.Keeper.BlacklistedAddr(addr)
}

func (adapter BankKeeperAdapter) IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error {
	if !adapter.Keeper.GetSendEnabled(ctx) {
		return sdkerrors.Wrapf(types.ErrSendDisabled, "transfers are currently disabled")
	}
	// todo weather allow different form okt coin send enable
	return nil
}

func (adapter BankKeeperAdapter) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return adapter.Keeper.GetCoins(ctx, addr)
}

func (adapter BankKeeperAdapter) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	coins := adapter.Keeper.GetCoins(ctx, addr)
	return sdk.Coin{
		Amount: coins.AmountOf(denom),
		Denom:  denom,
	}
}

func (adapter BankKeeperAdapter) HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool {
	return adapter.GetBalance(ctx, addr, amt.Denom).IsGTE(amt)
}

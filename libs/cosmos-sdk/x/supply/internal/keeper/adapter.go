package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/internal/types"
)

var (
	_ types.CM40AccountKeeper = SupplyKeerAdapter{}
	_ types.CM40BankKeeper    = SupplyKeerAdapter{}
)

type SupplyKeerAdapter struct {
	Keeper
}

func NewSupplyKeerAdapter(keeper Keeper) *SupplyKeerAdapter {
	return &SupplyKeerAdapter{Keeper: keeper}
}

func (k SupplyKeerAdapter) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return k.bk.SendCoins(ctx, fromAddr, toAddr, amt)
}

func (k SupplyKeerAdapter) NewAccount(ctx sdk.Context, acc authtypes.Account) authtypes.Account {
	return k.ak.NewAccount(ctx, acc)
}

func (k SupplyKeerAdapter) SetAccount(ctx sdk.Context, acc authtypes.Account) {
	k.ak.SetAccount(ctx, acc)
}

func (k SupplyKeerAdapter) HasBalance(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coin) bool {
	return k.bk.HasBalance(ctx, addr, amt)
}

func (k SupplyKeerAdapter) BlockedAddr(address sdk.AccAddress) bool {
	return k.bk.BlockedAddr(address)
}

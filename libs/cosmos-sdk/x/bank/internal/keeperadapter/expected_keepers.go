package keeperadapter

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	supplyexported "github.com/okx/okbchain/libs/cosmos-sdk/x/supply/exported"
)

type SupplyKeeper interface {
	GetSupply(ctx sdk.Context) (supply supplyexported.SupplyI)
}

type ViewBankKeeper interface {
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetSendEnabled(ctx sdk.Context) bool
}

type MsgServerBankKeeper interface {
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetSendEnabled(ctx sdk.Context) bool
	BlacklistedAddr(addr sdk.AccAddress) bool
}

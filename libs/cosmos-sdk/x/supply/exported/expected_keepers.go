package exported

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetModuleAccount(ctx sdk.Context, moduleName string) ModuleAccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error
	AddFeeBlockPool(amt sdk.Coins)
	GetFeeFromBlockPool() sdk.Coins
	ResetFeeBlockPool()
}

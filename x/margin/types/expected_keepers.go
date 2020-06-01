package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	dex "github.com/okex/okchain/x/dex/types"
)

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	WithKeyTable(table params.KeyTable) params.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps params.ParamSet)
	SetParamSet(ctx sdk.Context, ps params.ParamSet)
}

// TokenKeeper : expected token keeper
type TokenKeeper interface {
	// Token balance
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.DecCoins
}

// SupplyKeeper : expected supply keeper
type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string,
		amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress,
		amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) sdk.Error
	GetModuleAccount(ctx sdk.Context, moduleName string) exported.ModuleAccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
}

// DexKeeper : expected dex keeper
type DexKeeper interface {
	// TokenPair
	GetTokenPair(ctx sdk.Context, product string) *dex.TokenPair
}

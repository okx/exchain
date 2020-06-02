package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
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
	SendCoinsFromAccountToModule(ctx sdk.Context, from sdk.AccAddress, amount sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, to sdk.AccAddress, amount sdk.Coins) sdk.Error
	SendCoinsFromModuleToModule(ctx sdk.Context, recipientModule string, coins sdk.Coins) sdk.Error
}

// DexKeeper : expected dex keeper
type DexKeeper interface {
	// TokenPair
	GetTokenPair(ctx sdk.Context, product string) *dex.TokenPair
}

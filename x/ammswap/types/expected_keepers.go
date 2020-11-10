package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	token "github.com/okex/okexchain/x/token/types"
)

// ParamSubspace defines the expected Subspace interface
type ParamSubspace interface {
	WithKeyTable(table params.KeyTable) params.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps params.ParamSet)
	SetParamSet(ctx sdk.Context, ps params.ParamSet)
}

// BankKeeper defines the expected bank interface
type BankKeeper interface {
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

// SupplyKeeper defines the expected supply interface
type SupplyKeeper interface {
	GetSupplyByDenom(ctx sdk.Context, denom string) sdk.Dec
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string,
		recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress,
		recipientModule string, amt sdk.Coins) sdk.Error
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
}

// TokenKeeper defines the expected token interface
type TokenKeeper interface {
	GetTokenInfo(ctx sdk.Context, symbol string) token.Token
	NewToken(ctx sdk.Context, token token.Token)
	UpdateToken(ctx sdk.Context, token token.Token)
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.SysCoins
	TokenExist(ctx sdk.Context, symbol string) bool
	GetTokensInfo(ctx sdk.Context) (tokens []token.Token)
}

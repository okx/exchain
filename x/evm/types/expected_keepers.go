package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	exported2 "github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
)

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
	GetAllAccounts(ctx sdk.Context) (accounts []authexported.Account)
	IterateAccounts(ctx sdk.Context, cb func(account authexported.Account) bool)
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
	SetAccount(ctx sdk.Context, account authexported.Account)
	RemoveAccount(ctx sdk.Context, account authexported.Account)
	SetObserverKeeper(observer auth.ObserverI)
}

type SupplyKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error

	GetModuleAccount(ctx sdk.Context, moduleName string) exported2.ModuleAccountI
	GetFeeFromBlockPool() sdk.Coins
	AddCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) error
}

type Subspace interface {
	GetParamSet(ctx sdk.Context, ps params.ParamSet)
	SetParamSet(ctx sdk.Context, ps params.ParamSet)
}

type BankKeeper interface {
	BlacklistedAddr(addr sdk.AccAddress) bool
}

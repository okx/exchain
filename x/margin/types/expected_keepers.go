package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/okex/okchain/x/order/types"

	dex "github.com/okex/okchain/x/dex/types"
	token "github.com/okex/okchain/x/token/types"
)

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	WithKeyTable(table params.KeyTable) params.Subspace
	Get(ctx sdk.Context, key []byte, ptr interface{})
	GetParamSet(ctx sdk.Context, ps params.ParamSet)
	SetParamSet(ctx sdk.Context, ps params.ParamSet)
}

/*
When a module wishes to interact with another module, it is good practice to define what it will use
as an interface so the module cannot use things that are not permitted.
TODO: Create interfaces of what you expect the other keepers to have to be able to use this module.
type BankKeeper interface {
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, error)
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}
*/

// TokenKeeper : expected token keeper
type TokenKeeper interface {
	// Token balance
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.DecCoins
	LockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins, lockCoinsType int) error
	UnlockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins, lockCoinsType int) error
	BalanceAccount(ctx sdk.Context, addr sdk.AccAddress, outputCoins sdk.DecCoins, inputCoins sdk.DecCoins) error
	SendCoinsFromAccountToAccount(ctx sdk.Context, from, to sdk.AccAddress, amt sdk.DecCoins) error
	// Fee detail
	AddFeeDetail(ctx sdk.Context, from string, fee sdk.DecCoins, feeType string)
	GetAllLockedCoins(ctx sdk.Context) (locks []token.AccCoins)
	IterateLockedFees(ctx sdk.Context, cb func(acc sdk.AccAddress, coins sdk.DecCoins) (stop bool))
	GetCoinsInfo(ctx sdk.Context, addr sdk.AccAddress) (coinsInfo token.CoinsInfo)
}

// SupplyKeeper : expected supply keeper
type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string,
		amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress,
		amt sdk.Coins) sdk.Error
	GetModuleAccount(ctx sdk.Context, moduleName string) exported.ModuleAccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
}

// DexKeeper : expected dex keeper
type DexKeeper interface {
	// TokenPair
	GetTokenPair(ctx sdk.Context, product string) *dex.TokenPair
	GetTokenPairs(ctx sdk.Context) []*dex.TokenPair
	SortProducts(ctx sdk.Context, products []string)
	SaveTokenPair(ctx sdk.Context, tokenPair *dex.TokenPair) error
	UpdateTokenPair(ctx sdk.Context, product string, tokenPair *dex.TokenPair)
	GetTokenPairsFromStore(ctx sdk.Context) []*dex.TokenPair
	CheckTokenPairUnderDexDelist(ctx sdk.Context, product string) (isDelisting bool, err error)
	LockTokenPair(ctx sdk.Context, product string, lock *types.ProductLock)
	UnlockTokenPair(ctx sdk.Context, product string)
	IsTokenPairLocked(product string) bool
	GetLockedProductsCopy() *types.ProductLockMap
	IsAnyProductLocked() bool
}

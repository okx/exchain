package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/okex/okchain/x/dex/types"
	ordertypes "github.com/okex/okchain/x/order/types"
	"github.com/okex/okchain/x/params"
)

// SupplyKeeper defines the expected supply Keeper
type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress,
		recipientModule string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string,
		recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	GetModuleAccount(ctx sdk.Context, moduleName string) exported.ModuleAccountI
	GetModuleAddress(moduleName string) sdk.AccAddress
	MintCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) sdk.Error
}

// TokenKeeper defines the expected token Keeper
type TokenKeeper interface {
	TokenExist(ctx sdk.Context, symbol string) bool
}

// IKeeper defines the expected dex Keeper
type IKeeper interface {
	GetTokenPair(ctx sdk.Context, product string) *types.TokenPair
	GetTokenPairs(ctx sdk.Context) []*types.TokenPair
	GetUserTokenPairs(ctx sdk.Context, owner sdk.AccAddress) []*types.TokenPair
	GetTokenPairsOrdered(ctx sdk.Context) types.TokenPairs
	GetTokenPairsFromStore(ctx sdk.Context) (tokenPairs []*types.TokenPair)
	SaveTokenPair(ctx sdk.Context, tokenPair *types.TokenPair) error
	DeleteTokenPairByName(ctx sdk.Context, owner sdk.AccAddress, tokenPairName string)
	Deposit(ctx sdk.Context, product string, from sdk.AccAddress, amount sdk.DecCoin) sdk.Error
	Withdraw(ctx sdk.Context, product string, to sdk.AccAddress, amount sdk.DecCoin) sdk.Error
	GetSupplyKeeper() SupplyKeeper
	GetTokenKeeper() TokenKeeper
	GetParamSubspace() params.Subspace
	GetParams(ctx sdk.Context) (params types.Params)
	SetParams(ctx sdk.Context, params types.Params)
	GetFeeCollector() string
	GetCDC() *codec.Codec
	TransferOwnership(ctx sdk.Context, product string, from sdk.AccAddress, to sdk.AccAddress) sdk.Error
	LockTokenPair(ctx sdk.Context, product string, lock *ordertypes.ProductLock)
	LoadProductLocks(ctx sdk.Context) *ordertypes.ProductLockMap
	SetWithdrawInfo(ctx sdk.Context, withdrawInfo types.WithdrawInfo)
	SetWithdrawCompleteTimeAddress(ctx sdk.Context, completeTime time.Time, addr sdk.AccAddress)
	IterateWithdrawAddress(ctx sdk.Context, currentTime time.Time, fn func(index int64, key []byte) (stop bool))
	CompleteWithdraw(ctx sdk.Context, addr sdk.AccAddress) error
	IterateWithdrawInfo(ctx sdk.Context, fn func(index int64, withdrawInfo types.WithdrawInfo) (stop bool))
	DeleteWithdrawCompleteTimeAddress(ctx sdk.Context, timestamp time.Time, delAddr sdk.AccAddress)
	GetMaxTokenPairID(ctx sdk.Context) (tokenPairMaxID uint64)
	SetMaxTokenPairID(ctx sdk.Context, tokenPairMaxID uint64)
}

// StakingKeeper defines the expected staking Keeper (noalias)
type StakingKeeper interface {
	IsValidator(ctx sdk.Context, addr sdk.AccAddress) bool
}

// BankKeeper defines the expected bank Keeper
type BankKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// GovKeeper defines the expected gov Keeper
type GovKeeper interface {
	RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time)
}

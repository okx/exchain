package keeper

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	authexported "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/params"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/supply/exported"
	"github.com/okx/okbchain/libs/ibc-go/modules/apps/transfer/types"
	clienttypes "github.com/okx/okbchain/libs/ibc-go/modules/core/02-client/types"
	tmbytes "github.com/okx/okbchain/libs/tendermint/libs/bytes"
	govtypes "github.com/okx/okbchain/x/gov/types"
)

// GovKeeper defines the expected gov Keeper
type GovKeeper interface {
	GetDepositParams(ctx sdk.Context) govtypes.DepositParams
	GetVotingParams(ctx sdk.Context) govtypes.VotingParams
}

// AccountKeeper defines the expected account keeper interface
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
	SetAccount(ctx sdk.Context, acc authexported.Account)
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authexported.Account
}

type SupplyKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	GetModuleAccount(ctx sdk.Context, moduleName string) exported.ModuleAccountI
}

type Subspace interface {
	GetParamSet(ctx sdk.Context, ps params.ParamSet)
	SetParamSet(ctx sdk.Context, ps params.ParamSet)
}

type BankKeeper interface {
	BlacklistedAddr(addr sdk.AccAddress) bool
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
}

type TransferKeeper interface {
	SendTransfer(
		ctx sdk.Context,
		sourcePort,
		sourceChannel string,
		token sdk.CoinAdapter,
		sender sdk.AccAddress,
		receiver string,
		timeoutHeight clienttypes.Height,
		timeoutTimestamp uint64,
	) error
	DenomPathFromHash(ctx sdk.Context, denom string) (string, error)
	GetDenomTrace(ctx sdk.Context, denomTraceHash tmbytes.HexBytes) (types.DenomTrace, bool)
}

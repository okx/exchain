package types // noalias

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
	govtypes "github.com/okex/exchain/x/gov/types"
)

// StakingKeeper defines the expected staking keeper
type StakingKeeper interface {
	StakingTokenSupply(ctx sdk.Context) sdk.Dec
	BondedRatio(ctx sdk.Context) sdk.Dec
}

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, exported.ModuleAccountI)

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
}

// GovKeeper defines the expected gov Keeper
type GovKeeper interface {
	GetDepositParams(ctx sdk.Context) govtypes.DepositParams
	GetVotingParams(ctx sdk.Context) govtypes.VotingParams
}

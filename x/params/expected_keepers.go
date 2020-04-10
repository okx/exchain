package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper shows the expected action of bank keeper
type BankKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// StakingKeeper shows the expected action of staking keeper
type StakingKeeper interface {
	IsValidator(ctx sdk.Context, addr sdk.AccAddress) bool
}

// GovKeeper shows the expected action of gov keeper
type GovKeeper interface {
	InsertWaitingProposalQueue(ctx sdk.Context, blockHeight, proposalID uint64)
	RemoveFromWaitingProposalQueue(ctx sdk.Context, blockHeight, proposalID uint64)
}
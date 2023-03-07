package keeper

import (
	govtypes "github.com/okx/exchain/x/gov/types"

	sdk "github.com/okx/exchain/libs/cosmos-sdk/types"
	"time"
)

// GovKeeper defines the expected gov Keeper
type GovKeeper interface {
	RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time)
	GetDepositParams(ctx sdk.Context) govtypes.DepositParams
	GetVotingParams(ctx sdk.Context) govtypes.VotingParams
}

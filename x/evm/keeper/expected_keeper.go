package keeper

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	govtypes "github.com/okex/exchain/x/gov/types"
)

// GovKeeper defines the expected gov Keeper
type GovKeeper interface {
	GetDepositParams(ctx sdk.Context) govtypes.DepositParams
	GetVotingParams(ctx sdk.Context) govtypes.VotingParams
}

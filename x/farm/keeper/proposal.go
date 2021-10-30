package keeper

import (
	"fmt"
	"time"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/x/farm/types"
	sdkGov "github.com/okex/exchain/x/gov"
	govKeeper "github.com/okex/exchain/x/gov/keeper"
	govTypes "github.com/okex/exchain/x/gov/types"
)

var _ govKeeper.ProposalHandler = (*Keeper)(nil)

// GetMinDeposit returns min deposit
func (k Keeper) GetMinDeposit(ctx sdk.Context, content sdkGov.Content) (minDeposit sdk.SysCoins) {
	if _, ok := content.(types.ManageWhiteListProposal); ok {
		minDeposit = k.govKeeper.GetDepositParams(ctx).MinDeposit
	}

	return
}

// GetMaxDepositPeriod returns max deposit period
func (k Keeper) GetMaxDepositPeriod(ctx sdk.Context, content sdkGov.Content) (maxDepositPeriod time.Duration) {
	if _, ok := content.(types.ManageWhiteListProposal); ok {
		maxDepositPeriod = k.govKeeper.GetDepositParams(ctx).MaxDepositPeriod
	}

	return
}

// GetVotingPeriod returns voting period
func (k Keeper) GetVotingPeriod(ctx sdk.Context, content sdkGov.Content) (votingPeriod time.Duration) {
	if _, ok := content.(types.ManageWhiteListProposal); ok {
		votingPeriod = k.govKeeper.GetVotingParams(ctx).VotingPeriod
	}

	return
}

// CheckMsgSubmitProposal validates MsgSubmitProposal
func (k Keeper) CheckMsgSubmitProposal(ctx sdk.Context, msg govTypes.MsgSubmitProposal) sdk.Error {
	switch content := msg.Content.(type) {
	case types.ManageWhiteListProposal:
		return k.CheckMsgManageWhiteListProposal(ctx, content)
	default:
		return sdk.ErrUnknownRequest(fmt.Sprintf("unrecognized dex proposal content type: %T", content))
	}
}

// nolint
func (k Keeper) AfterSubmitProposalHandler(_ sdk.Context, _ govTypes.Proposal) {}
func (k Keeper) AfterDepositPeriodPassed(_ sdk.Context, _ govTypes.Proposal)   {}
func (k Keeper) RejectedHandler(_ sdk.Context, _ govTypes.Content)             {}
func (k Keeper) VoteHandler(_ sdk.Context, _ govTypes.Proposal, _ govTypes.Vote) (string, sdk.Error) {
	return "", nil
}

// CheckMsgManageWhiteListProposal checks msg manage white list proposal
func (k Keeper) CheckMsgManageWhiteListProposal(ctx sdk.Context, proposal types.ManageWhiteListProposal) sdk.Error {
	if proposal.IsAdded {
		// add pool name into white list
		// 1. check the existence
		pool, found := k.GetFarmPool(ctx, proposal.PoolName)
		if !found {
			return types.ErrNoFarmPoolFound(proposal.PoolName)
		}
		// 2. check the swap token pair
		if sdkErr := k.satisfyWhiteListAdmittance(ctx, pool); sdkErr != nil {
			return sdkErr
		}

		return nil
	}

	// delete the pool name from the white list
	// 1. check the existence of the pool name in whitelist
	if !k.isPoolNameExistedInWhiteList(ctx, proposal.PoolName) {
		return types.ErrPoolNameNotExistedInWhiteList(proposal.PoolName)
	}

	return nil
}

package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/evm/types"
	sdkGov "github.com/okex/okexchain/x/gov"
	govKeeper "github.com/okex/okexchain/x/gov/keeper"
	govTypes "github.com/okex/okexchain/x/gov/types"
)

var _ govKeeper.ProposalHandler = (*Keeper)(nil)

// GetMinDeposit returns min deposit
func (k Keeper) GetMinDeposit(ctx sdk.Context, content sdkGov.Content) (minDeposit sdk.SysCoins) {
	switch content.(type) {
	case types.ManageContractDeploymentWhitelistProposal, types.ManageContractBlockedListProposal:
		minDeposit = k.govKeeper.GetDepositParams(ctx).MinDeposit
	}

	return
}

// GetMaxDepositPeriod returns max deposit period
func (k Keeper) GetMaxDepositPeriod(ctx sdk.Context, content sdkGov.Content) (maxDepositPeriod time.Duration) {
	switch content.(type) {
	case types.ManageContractDeploymentWhitelistProposal, types.ManageContractBlockedListProposal:
		maxDepositPeriod = k.govKeeper.GetDepositParams(ctx).MaxDepositPeriod
	}

	return
}

// GetVotingPeriod returns voting period
func (k Keeper) GetVotingPeriod(ctx sdk.Context, content sdkGov.Content) (votingPeriod time.Duration) {
	switch content.(type) {
	case types.ManageContractDeploymentWhitelistProposal, types.ManageContractBlockedListProposal:
		votingPeriod = k.govKeeper.GetVotingParams(ctx).VotingPeriod
	}

	return
}

// CheckMsgSubmitProposal validates MsgSubmitProposal
func (k Keeper) CheckMsgSubmitProposal(ctx sdk.Context, msg govTypes.MsgSubmitProposal) sdk.Error {
	switch content := msg.Content.(type) {
	case types.ManageContractDeploymentWhitelistProposal, types.ManageContractBlockedListProposal:
		// whole target address list will be added/deleted to/from the contract deployment whitelist/contract blocked list.
		// It's not necessary to check the existence in CheckMsgSubmitProposal
		return nil
	default:
		return sdk.ErrUnknownRequest(fmt.Sprintf("unrecognized %s proposal content type: %T", types.DefaultCodespace, content))
	}
}

// nolint
func (k Keeper) AfterSubmitProposalHandler(_ sdk.Context, _ govTypes.Proposal) {}
func (k Keeper) AfterDepositPeriodPassed(_ sdk.Context, _ govTypes.Proposal)   {}
func (k Keeper) RejectedHandler(_ sdk.Context, _ govTypes.Content)             {}
func (k Keeper) VoteHandler(_ sdk.Context, _ govTypes.Proposal, _ govTypes.Vote) (string, sdk.Error) {
	return "", nil
}

// CheckMsgManageContractBlockedListProposal checks msg manage contract blocked list proposal
func (k Keeper) CheckMsgManageContractBlockedListProposal(ctx sdk.Context,
	proposal types.ManageContractBlockedListProposal) sdk.Error {
	csdb := types.CreateEmptyCommitStateDB(k.GeneratePureCSDBParams(), ctx)
	if proposal.IsAdded {
		// add contract addr into blocked list
		// 1. check the existence
		if csdb.IsContractInBlockedList(proposal.ContractAddr) {
			return types.ErrContractAlreadyExists(proposal.ContractAddr)
		}

		return nil
	}

	// delete the contract addr from the blocked list
	// 1. check the existence of contract addr in blocked list
	if !csdb.IsContractInBlockedList(proposal.ContractAddr) {
		return types.ErrContractNotExists(proposal.ContractAddr)
	}

	return nil
}

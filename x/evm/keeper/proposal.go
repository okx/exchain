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
	if _, ok := content.(types.ManageContractDeploymentWhitelistProposal); ok {
		minDeposit = k.govKeeper.GetDepositParams(ctx).MinDeposit
	}

	return
}

// GetMaxDepositPeriod returns max deposit period
func (k Keeper) GetMaxDepositPeriod(ctx sdk.Context, content sdkGov.Content) (maxDepositPeriod time.Duration) {
	if _, ok := content.(types.ManageContractDeploymentWhitelistProposal); ok {
		maxDepositPeriod = k.govKeeper.GetDepositParams(ctx).MaxDepositPeriod
	}

	return
}

// GetVotingPeriod returns voting period
func (k Keeper) GetVotingPeriod(ctx sdk.Context, content sdkGov.Content) (votingPeriod time.Duration) {
	if _, ok := content.(types.ManageContractDeploymentWhitelistProposal); ok {
		votingPeriod = k.govKeeper.GetVotingParams(ctx).VotingPeriod
	}

	return
}

// CheckMsgSubmitProposal validates MsgSubmitProposal
func (k Keeper) CheckMsgSubmitProposal(ctx sdk.Context, msg govTypes.MsgSubmitProposal) sdk.Error {
	switch content := msg.Content.(type) {
	case types.ManageContractDeploymentWhitelistProposal:
		return k.CheckMsgManageContractDeploymentWhitelistProposal(ctx, content)
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

// TODO
// CheckMsgManageContractDeploymentWhitelistProposal checks msg manage contract deployment whitelist proposal
func (k Keeper) CheckMsgManageContractDeploymentWhitelistProposal(ctx sdk.Context,
	proposal types.ManageContractDeploymentWhitelistProposal) sdk.Error {
	//if proposal.IsAdded {
	//	// add deployer addr into whitelist
	//	// 1. check the existence
	//	pool, found := k.GetFarmPool(ctx, proposal.PoolName)
	//	if !found {
	//		return types.ErrNoFarmPoolFound(proposal.PoolName)
	//	}
	//	// 2. check the swap token pair
	//	if sdkErr := k.satisfyWhiteListAdmittance(ctx, pool); sdkErr != nil {
	//		return sdkErr
	//	}
	//
	//	return nil
	//}
	//
	//// delete the pool name from the white list
	//// 1. check the existence of the pool name in whitelist
	//if !k.isPoolNameExistedInWhiteList(ctx, proposal.PoolName) {
	//	return types.ErrPoolNameNotExistedInWhiteList(proposal.PoolName)
	//}

	return nil
}

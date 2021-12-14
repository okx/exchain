package keeper

import (
	"fmt"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
	sdkGov "github.com/okex/exchain/x/gov"
	govKeeper "github.com/okex/exchain/x/gov/keeper"
	govTypes "github.com/okex/exchain/x/gov/types"
)

var _ govKeeper.ProposalHandler = (*Keeper)(nil)

// GetMinDeposit returns min deposit
func (k Keeper) GetMinDeposit(ctx sdk.Context, content sdkGov.Content) (minDeposit sdk.SysCoins) {
	switch content.(type) {
	case types.ManageContractDeploymentWhitelistProposal, types.ManageContractBlockedListProposal, types.ManageContractMethodBlockedListProposal:
		minDeposit = k.govKeeper.GetDepositParams(ctx).MinDeposit
	}

	return
}

// GetMaxDepositPeriod returns max deposit period
func (k Keeper) GetMaxDepositPeriod(ctx sdk.Context, content sdkGov.Content) (maxDepositPeriod time.Duration) {
	switch content.(type) {
	case types.ManageContractDeploymentWhitelistProposal, types.ManageContractBlockedListProposal, types.ManageContractMethodBlockedListProposal:
		maxDepositPeriod = k.govKeeper.GetDepositParams(ctx).MaxDepositPeriod
	}

	return
}

// GetVotingPeriod returns voting period
func (k Keeper) GetVotingPeriod(ctx sdk.Context, content sdkGov.Content) (votingPeriod time.Duration) {
	switch content.(type) {
	case types.ManageContractDeploymentWhitelistProposal, types.ManageContractBlockedListProposal, types.ManageContractMethodBlockedListProposal:
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
	case types.ManageContractMethodBlockedListProposal:
		csdb := types.CreateEmptyCommitStateDB(k.GeneratePureCSDBParams(), ctx)
		// can not delete address is not exist
		if !content.IsAdded {
			for i,_ := range content.ContractList {
				bc := csdb.GetContractMethodBlockedByAddress(content.ContractList[i].Address)
				if bc == nil {
					return types.ErrBlockedContractMethodIsNotExist(content.ContractList[i].Address,types.ErrorContractMethodBlockedIsNotExist)
				}
				if err := bc.BlockMethods.DeleteContractMethodMap(content.ContractList[i].BlockMethods);err != nil {
					return types.ErrBlockedContractMethodIsNotExist(content.ContractList[i].Address,err)
				}
			}
		}
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

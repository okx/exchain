package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/dex/types"
	"github.com/okex/okchain/x/gov"
	govTypes "github.com/okex/okchain/x/gov/types"
)

// implement ProposalHandler
func (k Keeper) GetMinDeposit(ctx sdk.Context, content gov.Content) (minDeposit sdk.DecCoins) {
	switch content.(type) {
	case types.DelistProposal:
		minDeposit = k.GetParams(ctx).DelistMinDeposit
	}
	return
}

func (k Keeper) GetMaxDepositPeriod(ctx sdk.Context, content gov.Content) (maxDepositPeriod time.Duration) {
	switch content.(type) {
	case types.DelistProposal:
		maxDepositPeriod = k.GetParams(ctx).DelistMaxDepositPeriod
	}
	return
}

func (k Keeper) GetVotingPeriod(ctx sdk.Context, content gov.Content) (votingPeriod time.Duration) {
	switch content.(type) {
	case types.DelistProposal:
		votingPeriod = k.GetParams(ctx).DelistVotingPeriod
	}
	return
}

// check msg Delist proposal
func (k Keeper) checkMsgDelistProposal(ctx sdk.Context, delistProposal types.DelistProposal, proposer sdk.AccAddress, initialDeposit sdk.DecCoins) sdk.Error {
	// check the proposer of the msg is a validator
	if !k.stakingKeeper.IsValidator(ctx, proposer) {
		return gov.ErrInvalidProposer(types.DefaultCodespace, "failed to submit proposal because the proposer of delist proposal should be a validator")
	}

	// check whether the baseAsset is in the Dex list
	if !k.isTokenPairExisted(ctx, delistProposal.BaseAsset, delistProposal.QuoteAsset) {
		return types.ErrInvalidProduct(fmt.Sprintf("failed to submit proposal because the asset with base asset '%s' and quote asset '%s' didn't exist on the Dex", delistProposal.BaseAsset, delistProposal.QuoteAsset))
	}

	// check the initial deposit
	localMinDeposit := k.GetParams(ctx).DelistMinDeposit.MulDec(sdk.NewDecWithPrec(1, 1))
	err := common.HasSufficientCoins(proposer, initialDeposit, localMinDeposit)

	if err != nil {
		return sdk.NewError(types.DefaultCodespace, types.CodeInvalidAsset, fmt.Sprintf("failed to submit proposal because initial deposit should be more than %s", localMinDeposit.String()))
	}

	// check whether the proposer can afford the initial deposit
	err = common.HasSufficientCoins(proposer, k.bankKeeper.GetCoins(ctx, proposer), initialDeposit)
	if err != nil {
		return sdk.NewError(types.DefaultCodespace, types.CodeInvalidBalanceNotEnough, fmt.Sprintf("failed to submit proposal because proposer %s didn't have enough coins to pay for the initial deposit %s", proposer, initialDeposit))
	}
	return nil
}

func (k Keeper) CheckMsgSubmitProposal(ctx sdk.Context, msg govTypes.MsgSubmitProposal) (sdkErr sdk.Error) {
	switch content := msg.Content.(type) {
	case types.DelistProposal:
		sdkErr = k.checkMsgDelistProposal(ctx, content, msg.Proposer, msg.InitialDeposit)
	default:
		errContent := fmt.Sprintf("unrecognized dex proposal content type: %T", content)
		sdkErr = sdk.ErrUnknownRequest(errContent)
	}
	return
}

// check whether the token pair constituted by baseAsset and quoteAsset exists on the dex
func (k Keeper) isTokenPairExisted(ctx sdk.Context, baseAsset, quoteAsset string) bool {
	tokenPairs := k.GetTokenPairs(ctx)
	tokenPairsLen := len(tokenPairs)
	for i := 0; i < tokenPairsLen; i++ {
		if tokenPairs[i].BaseAssetSymbol == baseAsset && tokenPairs[i].QuoteAssetSymbol == quoteAsset {
			return true
		}
	}
	return false
}

func (k Keeper) AfterSubmitProposalHandler(ctx sdk.Context, proposal govTypes.Proposal) {}

func (k Keeper) VoteHandler(ctx sdk.Context, proposal govTypes.Proposal, vote govTypes.Vote) (string, sdk.Error) {
	switch proposal.Content.(type) {
	case types.DelistProposal:
		delistProposal := proposal.Content.(types.DelistProposal)
		tokenPairName := delistProposal.BaseAsset + "_" + delistProposal.QuoteAsset
		if k.IsTokenPairLocked(tokenPairName) {
			errContent := fmt.Sprintf("the trading pair (%s) is locked, please retry later", tokenPairName)
			return "", sdk.ErrInternal(errContent)
		}
	}
	return "", nil
}

func (k Keeper) RejectedHandler(ctx sdk.Context, content govTypes.Content) {
	switch content := content.(type) {

	case types.DelistProposal:
		tokenPairName := fmt.Sprintf("%s_%s", content.BaseAsset, content.QuoteAsset)
		//update the token info from the store
		tokenPair := k.GetTokenPair(ctx, tokenPairName)
		tokenPair.Delisting = false
		k.UpdateTokenPair(ctx, tokenPairName, tokenPair)
	}
}

func (k Keeper) AfterDepositPeriodPassed(ctx sdk.Context, proposal govTypes.Proposal) {
	switch content := proposal.Content.(type) {
	case types.DelistProposal:
		tokenPairName := fmt.Sprintf("%s_%s", content.BaseAsset, content.QuoteAsset)
		// change the status of the token pair in the store
		tokenPair := k.GetTokenPair(ctx, tokenPairName)
		tokenPair.Delisting = true
		k.UpdateTokenPair(ctx, tokenPairName, tokenPair)
	}
}

func (k Keeper) RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	k.govKeeper.RemoveFromActiveProposalQueue(ctx, proposalID, endTime)
}

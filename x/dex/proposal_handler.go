package dex

import (
	"fmt"
	"github.com/okex/exchain/x/common"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/dex/types"
	govTypes "github.com/okex/exchain/x/gov/types"
)

// NewProposalHandler handles "gov" type message in "dex"
func NewProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		switch c := proposal.Content.(type) {
		case types.DelistProposal:
			return handleDelistProposal(ctx, k, proposal)
		default:
			return common.ErrUnknownProposalType(DefaultCodespace, fmt.Sprintf("%T", c))
		}
	}
}

func handleDelistProposal(ctx sdk.Context, keeper *Keeper, proposal *govTypes.Proposal) (err sdk.Error) {
	p := proposal.Content.(types.DelistProposal)
	logger := ctx.Logger().With("module", types.ModuleName)
	logger.Debug("execute DelistProposal begin")

	tokenPairName := fmt.Sprintf("%s_%s", p.BaseAsset, p.QuoteAsset)
	tokenPair := keeper.GetTokenPair(ctx, tokenPairName)
	if tokenPair == nil {
		return ErrTokenPairNotFound(tokenPairName)
	}
	if keeper.IsTokenPairLocked(ctx, tokenPairName) {
		return types.ErrIsTokenPairLocked(tokenPairName)
	}
	// withdraw
	if tokenPair.Deposits.IsPositive() {
		if err := keeper.Withdraw(ctx, tokenPair.Name(), tokenPair.Owner, tokenPair.Deposits); err != nil {
			return types.ErrWithdrawFailed(err.Error())
		}
	}

	// delete the token pair by its name from store and cache
	keeper.DeleteTokenPairByName(ctx, tokenPair.Owner, tokenPairName)

	// remove the delistProposal from the active proposal queue
	keeper.RemoveFromActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute("token-pair-deleted", tokenPairName),
		))
	return nil
}

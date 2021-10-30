package gov

import (
	"fmt"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/x/gov/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/okex/exchain/x/common/perf"
	"github.com/okex/exchain/x/gov/keeper"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	logger := k.Logger(ctx)

	seq := perf.GetPerf().OnEndBlockEnter(ctx, types.ModuleName)
	defer perf.GetPerf().OnEndBlockExit(ctx, types.ModuleName, seq)

	handleWaitingProposals(ctx, k, logger)
	handleInActiveProposals(ctx, k, logger)
	handleActiveProposals(ctx, k, logger)
}

// handle proposals which is executed on future block height appointed when it is submitted
func handleWaitingProposals(ctx sdk.Context, k keeper.Keeper, logger log.Logger) {
	k.IterateWaitingProposalsQueue(ctx, uint64(ctx.BlockHeight()), func(proposal Proposal) bool {
		handler := k.Router().GetRoute(proposal.ProposalRoute())
		cacheCtx, writeCache := ctx.CacheContext()
		err := handler(cacheCtx, &proposal)
		if err != nil {
			logger.Info(
				fmt.Sprintf("proposal %d (%s) excute failed",
					proposal.ProposalID,
					proposal.GetTitle(),
				),
			)
		} else {
			logger.Info(
				fmt.Sprintf("proposal %d (%s) excute successfully",
					proposal.ProposalID,
					proposal.GetTitle(),
				),
			)
			writeCache()
		}

		return false
	})
}

func handleInActiveProposals(ctx sdk.Context, k keeper.Keeper, logger log.Logger) {
	// delete inactive proposal from store and its deposits
	k.IterateInactiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal Proposal) bool {
		k.DeleteProposal(ctx, proposal.ProposalID)
		k.DistributeDeposits(ctx, proposal.ProposalID)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeInactiveProposal,
				sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
				sdk.NewAttribute(types.AttributeKeyProposalResult, types.AttributeValueProposalDropped),
			),
		)

		logger.Info(
			fmt.Sprintf("proposal %d (%s) didn't meet minimum deposit of %s (had only %s); deleted",
				proposal.ProposalID,
				proposal.GetTitle(),
				k.GetDepositParams(ctx).MinDeposit,
				proposal.TotalDeposit,
			),
		)
		return false
	})
}

func handleActiveProposals(ctx sdk.Context, k keeper.Keeper, logger log.Logger) {
	// fetch active proposals whose voting periods have ended (are passed the block time)
	k.IterateActiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal Proposal) bool {

		status, distribute, tallyResults := keeper.Tally(ctx, k, proposal, true)
		tagValue, logMsg := handleProposalAfterTally(ctx, k, &proposal, distribute, status)
		proposal.FinalTallyResult = tallyResults
		k.SetProposal(ctx, proposal)
		k.RemoveFromActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)
		k.DeleteVotes(ctx, proposal.ProposalID)

		logger.Info(
			fmt.Sprintf("proposal %d (%s) tallied; result: %s",
				proposal.ProposalID, proposal.GetTitle(), logMsg,
			),
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeActiveProposal,
				sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
				sdk.NewAttribute(types.AttributeKeyProposalResult, tagValue),
			),
		)
		return false
	})
}

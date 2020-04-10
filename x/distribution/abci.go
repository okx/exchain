package distribution

import (
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/distribution/keeper"
)

// BeginBlocker set the proposer for determining distribution during endblock
// and distribute rewards for the previous block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	logger := ctx.Logger().With("module", "distr")
	// determine the total power signing the block
	var previousTotalPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
		logger.Debug("Power", sdk.ValAddress(voteInfo.Validator.Address).String(), voteInfo.Validator.Power)
	}

	// TODO this is Tendermint-dependent
	// ref https://github.com/cosmos/cosmos-sdk/issues/3095
	if ctx.BlockHeight() > tmtypes.GetStartBlockHeight()+1 {
		previousProposer := k.GetPreviousProposerConsAddr(ctx)

		/* allocate tokens by okchain custom rule */
		k.AllocateTokens(ctx, previousTotalPower, previousProposer, req.LastCommitInfo.GetVotes())
	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)
}

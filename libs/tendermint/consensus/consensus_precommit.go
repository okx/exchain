package consensus

import (
	"fmt"
	cstypes "github.com/okex/exchain/libs/tendermint/consensus/types"
	"github.com/okex/exchain/libs/tendermint/libs/automation"
	"github.com/okex/exchain/libs/tendermint/types"
)

// Enter: `timeoutPrevote` after any +2/3 prevotes.
// Enter: `timeoutPrecommit` after any +2/3 precommits.
// Enter: +2/3 precomits for block or nil.
// Lock & precommit the ProposalBlock if we have enough prevotes for it (a POL in this round)
// else, unlock an existing lock and precommit nil if +2/3 of prevotes were nil,
// else, precommit nil otherwise.
func (cs *State) enterPrecommit(height int64, round int) {
	logger := cs.Logger.With("height", height, "round", round)

	if cs.Height != height || round < cs.Round || (cs.Round == round && cstypes.RoundStepPrecommit <= cs.Step) {
		logger.Debug(fmt.Sprintf(
			"enterPrecommit(%v/%v): Invalid args. Current step: %v/%v/%v",
			height,
			round,
			cs.Height,
			cs.Round,
			cs.Step))
		return
	}

	cs.initNewHeight()
	cs.trc.Pin("Precommit-%d", round)

	logger.Info(fmt.Sprintf("enterPrecommit(%v/%v). Current: %v/%v/%v", height, round, cs.Height, cs.Round, cs.Step))

	defer func() {
		// Done enterPrecommit:
		cs.updateRoundStep(round, cstypes.RoundStepPrecommit)
		cs.newStep()
	}()

	if automation.PrecommitNil(height, round) {
		cs.signAddVote(types.PrecommitType, nil, types.PartSetHeader{})
		return
	}

	if cs.ProposalBlock == nil {
		if cs.LockedBlock == nil {
			logger.Info("enterPrecommit: +2/3 prevoted for nil.")
		} else {
			logger.Info("enterPrecommit: +2/3 prevoted for nil. Unlocking")
			cs.LockedRound = -1
			cs.LockedBlock = nil
			cs.LockedBlockParts = nil
			cs.eventBus.PublishEventUnlock(cs.RoundStateEvent())
		}
		cs.signAddVote(types.PrecommitType, nil, types.PartSetHeader{})
		return
	}
	// At this point, +2/3 prevoted for a particular block.

	blockID := types.BlockID{Hash: cs.ProposalBlock.Hash(), PartsHeader: cs.ProposalBlockParts.Header()}
	// If we're already locked on that block, precommit it, and update the LockedRound
	if cs.LockedBlock.HashesTo(blockID.Hash) {
		logger.Info("enterPrecommit: +2/3 prevoted locked block. Relocking")
		cs.LockedRound = round
		cs.eventBus.PublishEventRelock(cs.RoundStateEvent())
		cs.signAddVote(types.PrecommitType, blockID.Hash, blockID.PartsHeader)
		return
	}

	if cs.LockedBlock != nil && cs.LockedRound == round {
		// already precommit vote in this round
		return
	}

	// If +2/3 prevoted for proposal block, stage and precommit it
	logger.Info("enterPrecommit: +2/3 prevoted proposal block. Locking", "hash", blockID.Hash)
	// Validate the block.
	if err := cs.blockExec.ValidateBlock(cs.state, cs.ProposalBlock); err != nil {
		panic(fmt.Sprintf("enterPrecommit: +2/3 prevoted for an invalid block: %v", err))
	}
	cs.LockedRound = round
	cs.LockedBlock = cs.ProposalBlock
	cs.LockedBlockParts = cs.ProposalBlockParts
	cs.eventBus.PublishEventLock(cs.RoundStateEvent())
	cs.signAddVote(types.PrecommitType, blockID.Hash, blockID.PartsHeader)
	return

}

// Enter: any +2/3 precommits for next round.
func (cs *State) enterPrecommitWait(height int64, round int) {
	logger := cs.Logger.With("height", height, "round", round)

	if cs.Height != height || round < cs.Round || (cs.Round == round && cs.TriggeredTimeoutPrecommit) {
		logger.Debug(
			fmt.Sprintf(
				"enterPrecommitWait(%v/%v): Invalid args. "+
					"Current state is Height/Round: %v/%v/, TriggeredTimeoutPrecommit:%v",
				height, round, cs.Height, cs.Round, cs.TriggeredTimeoutPrecommit))
		return
	}

	cs.initNewHeight()
	cs.trc.Pin("PrecommitWait-%d", round)

	if !cs.Votes.Precommits(round).HasTwoThirdsAny() {
		panic(fmt.Sprintf("enterPrecommitWait(%v/%v), but Precommits does not have any +2/3 votes", height, round))
	}
	logger.Info(fmt.Sprintf("enterPrecommitWait(%v/%v). Current: %v/%v/%v", height, round, cs.Height, cs.Round, cs.Step))

	defer func() {
		// Done enterPrecommitWait:
		cs.TriggeredTimeoutPrecommit = true
		cs.newStep()
	}()

	// Wait for some more precommits; enterNewRound
	cs.scheduleTimeout(cs.config.Precommit(round), height, round, cstypes.RoundStepPrecommitWait)

}

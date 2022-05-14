package consensus

import (
	"bytes"
	"fmt"
	cstypes "github.com/okex/exchain/libs/tendermint/consensus/types"
	"github.com/okex/exchain/libs/tendermint/types"
)

// Enter (CreateEmptyBlocks): from enterNewRound(height,round)
// Enter (CreateEmptyBlocks, CreateEmptyBlocksInterval > 0 ):
// 		after enterNewRound(height,round), after timeout of CreateEmptyBlocksInterval
// Enter (!CreateEmptyBlocks) : after enterNewRound(height,round), once txs are in the mempool
func (cs *State) enterPropose(height int64, round int) {
	logger := cs.Logger.With("height", height, "round", round)
	if cs.Height != height || round < cs.Round || (cs.Round == round && cstypes.RoundStepPropose <= cs.Step) {
		logger.Debug(fmt.Sprintf(
			"enterPropose(%v/%v): Invalid args. Current step: %v/%v/%v",
			height,
			round,
			cs.Height,
			cs.Round,
			cs.Step))
		return
	}

	cs.initNewHeight()
	isBlockProducer, bpAddr := cs.isBlockProducer()
	cs.trc.Pin("H%d-Propose-%d-%s-%s", height, round, isBlockProducer, bpAddr)

	logger.Info(fmt.Sprintf("enterPropose(%v/%v). Current: %v/%v/%v", height, round, cs.Height, cs.Round, cs.Step))

	defer func() {
		// Done enterPropose:
		cs.updateRoundStep(round, cstypes.RoundStepPropose)
		cs.newStep()

		// If we have the whole proposal + POL, then goto Prevote now.
		// else, we'll enterPrevote when the rest of the proposal is received (in AddProposalBlockPart),
		// or else after timeoutPropose
		if cs.isProposalComplete() {
			cs.enterPrevote(height, cs.Round)
		}
	}()

	// If we don't get the proposal and all block parts quick enough, enterPrevote
	cs.scheduleTimeout(cs.config.Propose(round), height, round, cstypes.RoundStepPropose)

	if isBlockProducer == "y" {
		logger.Info("enterPropose: Our turn to propose",
			"proposer",
			bpAddr,
			"privValidator",
			cs.privValidator)
		cs.decideProposal(height, round)
	} else {
		logger.Info("enterPropose: Not our turn to propose",
			"proposer",
			cs.Validators.GetProposer().Address,
			"privValidator",
			cs.privValidator)
	}
}

func (cs *State) isProposer(address []byte) bool {
	return bytes.Equal(cs.Validators.GetProposer().Address, address)
}

func (cs *State) defaultDecideProposal(height int64, round int) {
	var block *types.Block
	var blockParts *types.PartSet

	// Decide on block
	if cs.ValidBlock != nil {
		// If there is valid block, choose that.
		block, blockParts = cs.ValidBlock, cs.ValidBlockParts
	} else {
		// Create a new proposal block from state/txs from the mempool.
		block, blockParts = cs.createProposalBlock()
		if block == nil {
			return
		}
	}

	// Flush the WAL. Otherwise, we may not recompute the same proposal to sign,
	// and the privValidator will refuse to sign anything.
	cs.wal.FlushAndSync()

	// Make proposal
	propBlockID := types.BlockID{Hash: block.Hash(), PartsHeader: blockParts.Header()}
	proposal := types.NewProposal(height, round, cs.ValidRound, propBlockID)
	if err := cs.privValidator.SignProposal(cs.state.ChainID, proposal); err == nil {

		// send proposal and block parts on internal msg queue
		cs.sendInternalMessage(msgInfo{&ProposalMessage{proposal}, ""})
		for i := 0; i < blockParts.Total(); i++ {
			part := blockParts.GetPart(i)
			cs.sendInternalMessage(msgInfo{&BlockPartMessage{cs.Height, cs.Round, part, cs.Deltas}, ""})
		}
		cs.Logger.Info("Signed proposal", "height", height, "round", round, "proposal", proposal)
		cs.Logger.Debug(fmt.Sprintf("Signed proposal block: %v", block))
	} else if !cs.replayMode {
		cs.Logger.Error("enterPropose: Error signing proposal", "height", height, "round", round, "err", err)
	}
}

// Returns true if the proposal block is complete &&
// (if POLRound was proposed, we have +2/3 prevotes from there).
func (cs *State) isProposalComplete() bool {
	if cs.Proposal == nil || cs.ProposalBlock == nil {
		return false
	}
	// we have the proposal. if there's a POLRound,
	// make sure we have the prevotes from it too
	if cs.Proposal.POLRound < 0 {
		return true
	}
	// if this is false the proposer is lying or we haven't received the POL yet
	return cs.Votes.Prevotes(cs.Proposal.POLRound).HasTwoThirdsMajority()

}

// Create the next block to propose and return it. Returns nil block upon error.
//
// We really only need to return the parts, but the block is returned for
// convenience so we can log the proposal block.
//
// NOTE: keep it side-effect free for clarity.
// CONTRACT: cs.privValidator is not nil.
func (cs *State) createProposalBlock() (block *types.Block, blockParts *types.PartSet) {
	if cs.privValidator == nil {
		panic("entered createProposalBlock with privValidator being nil")
	}

	var commit *types.Commit
	switch {
	case cs.Height == types.GetStartBlockHeight()+1:
		// We're creating a proposal for the first block.
		// The commit is empty, but not nil.
		commit = types.NewCommit(0, 0, types.BlockID{}, nil)
	case cs.LastCommit.HasTwoThirdsMajority():
		// Make the commit from LastCommit
		commit = cs.LastCommit.MakeCommit()
	default: // This shouldn't happen.
		cs.Logger.Error("enterPropose: Cannot propose anything: No commit for the previous block")
		return
	}

	if cs.privValidatorPubKey == nil {
		// If this node is a validator & proposer in the current round, it will
		// miss the opportunity to create a block.
		cs.Logger.Error(fmt.Sprintf("enterPropose: %v", errPubKeyIsNotSet))
		return
	}
	proposerAddr := cs.privValidatorPubKey.Address()

	return cs.blockExec.CreateProposalBlock(cs.Height, cs.state, commit, proposerAddr)
}
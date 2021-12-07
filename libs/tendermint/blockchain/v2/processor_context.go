package v2

import (
	"fmt"

	"github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/types"
)

type processorContext interface {
	applyBlock(blockID types.BlockID, block *types.Block, deltas *types.Deltas) error
	verifyCommit(chainID string, blockID types.BlockID, height int64, commit *types.Commit) error
	saveBlock(block *types.Block, blockParts *types.PartSet, seenCommit *types.Commit)
	saveDeltas(deltas *types.Deltas, height int64)
	tmState() state.State
}

type pContext struct {
	store   blockStore
	dstore  deltaStore
	applier blockApplier
	state   state.State
}

func newProcessorContext(st blockStore, dst deltaStore, ex blockApplier, s state.State) *pContext {
	return &pContext{
		store:   st,
		dstore:  dst,
		applier: ex,
		state:   s,
	}
}

func (pc *pContext) applyBlock(blockID types.BlockID, block *types.Block, deltas *types.Deltas) error {
	newState, _, err := pc.applier.ApplyBlock(pc.state, blockID, block, deltas)
	pc.state = newState
	return err
}

func (pc pContext) tmState() state.State {
	return pc.state
}

func (pc pContext) verifyCommit(chainID string, blockID types.BlockID, height int64, commit *types.Commit) error {
	return pc.state.Validators.VerifyCommit(chainID, blockID, height, commit)
}

func (pc *pContext) saveBlock(block *types.Block, blockParts *types.PartSet, seenCommit *types.Commit) {
	pc.store.SaveBlock(block, blockParts, seenCommit)
}

func (pc *pContext) saveDeltas(deltas *types.Deltas, height int64) {
	pc.dstore.SaveDeltas(deltas, height)
}

type mockPContext struct {
	applicationBL  []int64
	verificationBL []int64
	state          state.State
}

func newMockProcessorContext(
	state state.State,
	verificationBlackList []int64,
	applicationBlackList []int64) *mockPContext {
	return &mockPContext{
		applicationBL:  applicationBlackList,
		verificationBL: verificationBlackList,
		state:          state,
	}
}

func (mpc *mockPContext) applyBlock(blockID types.BlockID, block *types.Block, deltas *types.Deltas) error {
	for _, h := range mpc.applicationBL {
		if h == block.Height {
			return fmt.Errorf("generic application error")
		}
	}
	mpc.state.LastBlockHeight = block.Height
	return nil
}

func (mpc *mockPContext) verifyCommit(chainID string, blockID types.BlockID, height int64, commit *types.Commit) error {
	for _, h := range mpc.verificationBL {
		if h == height {
			return fmt.Errorf("generic verification error")
		}
	}
	return nil
}

func (mpc *mockPContext) saveBlock(block *types.Block, blockParts *types.PartSet, seenCommit *types.Commit) {

}

func (mpc *mockPContext) saveDeltas(deltas *types.Deltas, height int64) {

}

func (mpc *mockPContext) tmState() state.State {
	return mpc.state
}

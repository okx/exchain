package state

import (
	"bytes"
	"fmt"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
)


type prerunContext struct {
	prerunTx bool
	prerunChan chan *executionTask
	prerunResultChan chan *executionTask
	prerunTask *executionTask
	logger log.Logger
}

func newPrerunContex(logger log.Logger) *prerunContext {
	return &prerunContext{
		prerunChan:           make(chan *executionTask, 1),
		prerunResultChan:     make(chan *executionTask, 1),
		logger:         logger,
	}
}

func (pc *prerunContext) checkIndex(block *types.Block) {
	var index int64
	if pc.prerunTask != nil {
		index = pc.prerunTask.index
	}
	pc.logger.Info("Not apply delta", "height", block.Height,
		"block-size", block.Size(), "prerunIndex", index)

}


func (pc *prerunContext) flushPrerunResult() {
	for {
		select {
		case task := <-pc.prerunResultChan:
			task.dump("Flush prerun result")
		default:
			return
		}
	}
}

func (pc *prerunContext) prerunRoutine() {
	pc.prerunTx = true
	for task := range pc.prerunChan {
		task.run()
	}
}

func (pc *prerunContext) dequeueResult() (*ABCIResponses, error) {
	expected := pc.prerunTask
	for context := range pc.prerunResultChan {

		context.dump("Got prerun result")

		if context.stopped {
			continue
		}

		if context.height != expected.block.Height {
			continue
		}

		if context.index != expected.index {
			continue
		}

		if bytes.Equal(context.block.AppHash, expected.block.AppHash) {
			return context.result.res, context.result.err
		} else {
			// todo
			panic("wrong app hash")
		}
	}
	return nil, nil
}

func (pc *prerunContext) stopPrerun(height int64) (index int64) {
	task := pc.prerunTask
	// stop the existing prerun if any
	if task != nil {
		if height > 0 && height != task.block.Height {
			task.dump(fmt.Sprintf(
				"Prerun sanity check failed. block.Height=%d, context.block.Height=%d",
				height,
				task.block.Height))

			// todo
			panic("Prerun sanity check failed")
		}
		task.dump("Stopping prerun")
		task.stop()

		index = task.index
	}
	pc.flushPrerunResult()
	pc.prerunTask = nil
	return index
}


func (pc *prerunContext) notifyPrerun(blockExec *BlockExecutor, block *types.Block) {

	stoppedIndex := pc.stopPrerun(block.Height)
	stoppedIndex++

	pc.prerunTask = newExecutionTask(blockExec, block, stoppedIndex)

	pc.prerunTask.dump("Notify prerun")

	// start a new one
	pc.prerunChan <- pc.prerunTask
}

func (pc *prerunContext) getPrerunResult(block *types.Block, fastSync bool) (res *ABCIResponses, err error) {

	pc.checkIndex(block)

	if fastSync {
		pc.stopPrerun(block.Height)
	}
	// blockExec.prerunContext == nil means:
	// 1. prerunTx disabled
	// 2. we are in fasy-sync: the block comes from BlockPool.AddBlock not State.addProposalBlockPart and no prerun result expected
	if pc.prerunTask != nil {
		res, err = pc.dequeueResult()
		pc.prerunTask = nil
	}
	return
}
package state

import (
	"errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/types"
)

type PreExecBlockResult struct {
	Response *ABCIResponses
	Err      error
}

type InternalMsg struct {
	resChan chan *PreExecBlockResult
}

var (
	RepeatedErr           = errors.New("block can not start over twice")
	CancelErr             = errors.New("block has been canceled")
	NotMatchErr           = errors.New("block has no start record")
	NotProactivelyModeErr = errors.New("not proactively mode")
	GenesisErr            = errors.New("genesis block don't proactively run")
	consensusFailed       bool
)

//single thread process
func (blockExec *BlockExecutor) LoopProcessProactivelyRun() {
	for {
		select {
		case block := <-blockExec.proactiveQueue:
			if res := blockExec.StartProactivelyRun(block); res != nil {
				blockExec.proactiveRes <- res
			}
		}
	}
}

//start a proactively block execution
func (blockExec *BlockExecutor) proactivelyOpen() bool {
	if blockExec.proactivelyFlag {
		return true
	}
	return false
}

//start a proactively block execution
func (blockExec *BlockExecutor) CancelAndStartNewRun(block *types.Block) {
	if !blockExec.proactivelyOpen() {
		return
	}

	// cancel
	// set consensusFailed let unfinished preRun to return
	consensusFailed = true
	// wait preRun return
	blockExec.wg.Wait()
	if len(blockExec.proactiveRes) > 0 {
		if len(blockExec.proactiveRes) != 1 {
			blockExec.logger.Error(" proactiveRes chan length is wrong", "length", len(blockExec.proactiveRes))
		}
		for i := 0; i < len(blockExec.proactiveRes); i++ {
			<-blockExec.proactiveRes
			blockExec.logger.Info(" clean proactiveRes chan")
		}
	}
	// reset consensusFailed and deliverState for next preRun
	consensusFailed = false
	blockExec.ResetDeliverState()
	// start
	blockExec.SetStartBlock(block)
}

func (blockExec *BlockExecutor) SetStartBlock(block *types.Block) {
	if !blockExec.proactivelyOpen() {
		return
	}

	if block != nil {
		blockExec.proactiveQueue <- block
	}
}

//start a proactively block execution
func (blockExec *BlockExecutor) StartProactivelyRun(block *types.Block) *PreExecBlockResult {
	if !blockExec.proactivelyOpen() {
		return nil
	}
	blockExec.wg.Add(1)
	defer blockExec.wg.Done()

	var abciResponses *ABCIResponses
	var err error
	var preBlockRes *PreExecBlockResult
	if blockExec.isAsync {
		abciResponses, err = execBlockOnProxyAppAsync(blockExec.logger, blockExec.proxyApp, block, blockExec.db)
	} else {
		abciResponses, err = execBlockOnProxyApp(blockExec.logger, blockExec.proxyApp, block, blockExec.db)
	}

	if err != nil {
		preBlockRes = &PreExecBlockResult{abciResponses, err}
	} else {
		preBlockRes = &PreExecBlockResult{abciResponses, nil}
	}

	return preBlockRes
}

//return result channel for caller
func (blockExec *BlockExecutor) GetProactivelyRes(reqBlock *types.Block) (chan *PreExecBlockResult, error) {
	if !blockExec.proactivelyOpen() {
		return nil, NotProactivelyModeErr
	}
	return blockExec.proactiveRes, nil
}

//reset deliverState
func (blockExec *BlockExecutor) ResetDeliverState() {
	blockExec.proxyApp.SetOptionSync(abci.RequestSetOption{
		Key: "ResetDeliverState",
	})
}

//set blockExec proactivelyFlag
func (blockExec *BlockExecutor) SetProactivelyFlag(open bool) {
	blockExec.proactivelyFlag = open
	if blockExec.proactivelyFlag {
		go blockExec.LoopProcessProactivelyRun()
	}
}

func IsFirstHeight(block *types.Block) bool {
	if block.Height == 1 {
		return true
	}
	return false
}

package state

import (
	"errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/types"
)

type PreExecBlockResult struct {
	*ABCIResponses
	error
}

type InternalMsg struct {
	cancelFlag bool
	resChan chan *PreExecBlockResult
}

var (
	RepeatedErr     = errors.New("block can not start over twice")
	CancelErr       = errors.New("block has been canceled")
	NotMatchErr     = errors.New("block has no start record")
	consensusFailed bool
)

//start a proactively block execution
func (blockExec *BlockExecutor) StartPreExecBlock(block *types.Block) error {
	if _, ok := blockExec.abciResponse.Load(block); ok {
		// start block twice
		return RepeatedErr
	} else {
		intMsg := &InternalMsg{
			resChan: make(chan *PreExecBlockResult),
		}
		blockExec.abciResponse.Store(block, intMsg)
		go blockExec.DoPreExecBlock(intMsg, block)
		blockExec.lastBlock = block
		return nil
	}
}

//return blockExec.abciResponse num
func (blockExec *BlockExecutor) mapCount() int {
	var count int
	blockExec.abciResponse.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

//gorountine execute block
func (blockExec *BlockExecutor) DoPreExecBlock(channels *InternalMsg, block *types.Block) {
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

	select {
	case channels.resChan <- preBlockRes:
		//only the writer close channel is safe
		close(channels.resChan)
	}
}

//cancel a block already execute
func (blockExec *BlockExecutor) CancelPreExecBlock(block *types.Block) error {
	if channels, ok := blockExec.abciResponse.Load(block); !ok {
		// cancel block not start
		return NotMatchErr
	} else {
		chann := channels.(*InternalMsg)
		// set cancel flag
		chann.cancelFlag = true
		consensusFailed = true
		// read useless result ensure
		<-chann.resChan
		// after execute done reset all
		blockExec.ResetDeliverState()
		blockExec.abciResponse.Delete(blockExec.lastBlock)
		consensusFailed = false
		return nil
	}
}

//return result channel for caller
func (blockExec *BlockExecutor) GetPreExecBlockRes(block *types.Block) (chan *PreExecBlockResult, error) {
	if channels, ok := blockExec.abciResponse.Load(block); !ok {
		// cancel block not start
		return nil, NotMatchErr
	} else {
		chann := channels.(*InternalMsg)
		return chann.resChan, nil
	}
}

//close block channel , clean abciResponse and check abciResponse num
func (blockExec *BlockExecutor) CleanPreExecBlockRes(block *types.Block) {
	if _, ok := blockExec.abciResponse.Load(block); !ok {
		// cancel block not start
		return
	} else {
		blockExec.abciResponse.Delete(block)
		if blockExec.lastBlock == block {
			blockExec.ResetLastBlock()
		}
		if num := blockExec.mapCount(); num != 0 {
			//check sync.Map num, should always be 0
			blockExec.logger.Error("blockExec abciResponse num not 0 ", "num", num)
		}
	}
}

//reset deliverState
func (blockExec *BlockExecutor) ResetDeliverState() {
	blockExec.proxyApp.SetOptionSync(abci.RequestSetOption{
		Key: "ResetDeliverState",
	})

}

//get lastBlock
func (blockExec *BlockExecutor) GetLastBlock() *types.Block {
	return blockExec.lastBlock
}

//reset lastBlock and clean abciResponse
func (blockExec *BlockExecutor) ResetLastBlock() {
	blockExec.abciResponse.Delete(blockExec.lastBlock)
	blockExec.lastBlock = nil
}

//set blockExec proactivelyFlag
func (blockExec *BlockExecutor) SetProactivelyFlag(open bool) {
	blockExec.proactivelyFlag = open
}

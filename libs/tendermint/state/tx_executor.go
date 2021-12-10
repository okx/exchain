package state

import (
	"errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/types"
	"time"

	"sync"
)

type PreExecBlockResult struct {
	Response *ABCIResponses
	State    int
	Break    bool
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
	breakManager *BreakManager
)

func init() {
	breakManager = &BreakManager{}
}

//single thread process
func (blockExec *BlockExecutor) LoopProcessProactivelyRun() error {
	for {
		select {
		case block := <-blockExec.proactiveQueue:
			abicRes := blockExec.StartProactivelyRun(block)
			if res, ok := blockExec.abciAllResponse.Load(block); ok {
				blockRes := res.(*PreExecBlockResult)
				if blockRes.State == 1 {
					breakManager.RemoveBreak(block)
				}
				blockRes.Response = abicRes.Response
				blockRes.Err = abicRes.Err
			} else {
				panic("block not exist!!!")
			}

		}
	}
}

//start a proactively block execution
func (blockExec *BlockExecutor) CancelAndStartNewRun(cancelBlock, startBlock *types.Block) {
	if !blockExec.proactivelyFlag {
		return
	}

	// cancel
	var localCancelBlock *types.Block
	if cancelBlock != nil {
		localCancelBlock = cancelBlock
	} else {
		localCancelBlock = blockExec.lastBlock
	}
	if res, ok := blockExec.abciAllResponse.Load(localCancelBlock); ok {
		result := res.(*PreExecBlockResult)
		if result.State == 1 {
			blockExec.ResetDeliverState()
			if result.Response == nil {
				// block not done
				result.Break = true
			}
		}
	}

	// start
	blockExec.SetStartBlock(startBlock)
}

func (blockExec *BlockExecutor) SetStartBlock(block *types.Block) {
	if !blockExec.proactivelyFlag {
		return
	}

	if block != nil {
		blockExec.lastBlock = block
		blockExec.abciAllResponse.Store(block, &PreExecBlockResult{})
		blockExec.proactiveQueue <- block
	}
}

//start a proactively block execution
func (blockExec *BlockExecutor) StartProactivelyRun(block *types.Block) *PreExecBlockResult {

	if !blockExec.proactivelyFlag {
		return nil
	}

	var abciResponses *ABCIResponses
	var err error
	var preBlockRes *PreExecBlockResult
	if blockExec.isAsync {
		abciResponses, err = execBlockOnProxyAppAsync(blockExec.logger, blockExec.proxyApp, block, blockExec.db)
	} else {
		abciResponses, err = execBlockOnProxyApp(blockExec.logger, blockExec.proxyApp, block, blockExec.db, nil)
	}

	if err != nil {
		preBlockRes = &PreExecBlockResult{abciResponses, 0, false, err}
	} else {
		preBlockRes = &PreExecBlockResult{abciResponses, 0, false, nil}
	}

	return preBlockRes

}

//return result channel for caller
func (blockExec *BlockExecutor) GetProactivelyRes(reqBlock *types.Block) (*PreExecBlockResult, error) {
	if !blockExec.proactivelyFlag {
		return nil, NotProactivelyModeErr
	}
	// 5 ms check
	ticker := time.NewTicker(5 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			if res, ok := blockExec.abciAllResponse.Load(reqBlock); !ok {
				return nil, NotMatchErr
			} else {
				rBlock := res.(*PreExecBlockResult)
				if rBlock.Response != nil {
					blockExec.abciAllResponse.Delete(reqBlock)
					return rBlock, nil
				}
			}
		}
	}
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

func (blockExec *BlockExecutor) SetLastBlockInvalid() error {
	if !blockExec.proactivelyFlag {
		return nil
	}

	if res, ok := blockExec.abciAllResponse.Load(blockExec.lastBlock); !ok {
		return NotMatchErr
	} else {
		rblock := res.(*PreExecBlockResult)
		rblock.State = 1
		blockExec.abciAllResponse.Store(blockExec.lastBlock, rblock)
		return nil
	}

}

type BreakManager struct {
	record sync.Map
}

func (s *BreakManager) SetBreak(block *types.Block) {
	s.record.Store(block, struct{}{})
}

func (s *BreakManager) GetBreak(block *types.Block) bool {
	if _, ok := s.record.Load(block); ok {
		return true
	}
	return false

}

func (s *BreakManager) RemoveBreak(block *types.Block) {
	s.record.Delete(block)
}


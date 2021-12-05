package state

import (
	"errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/types"
	"strings"
	"time"
)

type PreExecBlockResult struct {
	//*Elaped
	*ABCIResponses
	error
}

type InternalMsg struct {
	cancelChan chan struct{}
	resChan    chan *PreExecBlockResult
}

var (
	RepeatedErr = errors.New("block can not start over twice")
	CancelErr   = errors.New("block has been canceled")
	NotMatchErr = errors.New("block has no start record")
)

func GetNowTimeMs() int64 {
	return time.Now().UnixNano() / 1e6
}

func (blockExec *BlockExecutor) StartPreExecBlock(block *types.Block) error {
	if _, ok := blockExec.abciResponse.Load(block); ok {
		// start block twice
		return RepeatedErr
	} else {
		intMsg := &InternalMsg{
			cancelChan: make(chan struct{}),
			resChan:    make(chan *PreExecBlockResult),
		}
		blockExec.abciResponse.Store(block, intMsg)
		go blockExec.DoPreExecBlock(intMsg, block)
		blockExec.lastBlock = block
		return nil
	}
}

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
	case <-channels.cancelChan:
		channels.resChan <- &PreExecBlockResult{nil, CancelErr}
	case channels.resChan <- preBlockRes:

	}
}

func (blockExec *BlockExecutor) CancelPreExecBlock(block *types.Block) error {

	if channels, ok := blockExec.abciResponse.Load(block); !ok {
		// cancel block not start
		return NotMatchErr
	} else {
		chann := channels.(*InternalMsg)
		go func() {
			chann.cancelChan <- struct{}{}
		}()
		return nil
	}
}

func (blockExec *BlockExecutor) GetPreExecBlockRes(block *types.Block) (chan *PreExecBlockResult, error) {
	if channels, ok := blockExec.abciResponse.Load(block); !ok {
		// cancel block not start
		return nil, NotMatchErr
	} else {
		chann := channels.(*InternalMsg)
		return chann.resChan, nil
	}
}

func (blockExec *BlockExecutor) CleanPreExecBlockRes(block *types.Block) {

	if channels, ok := blockExec.abciResponse.Load(block); !ok {
		// cancel block not start
		return
	} else {
		chann := channels.(*InternalMsg)
		close(chann.resChan)
		close(chann.cancelChan)
		blockExec.abciResponse.Delete(block)
		if blockExec.lastBlock == block {
			blockExec.ResetLastBlock()
		}
	}
}

//reset base deliverState
func (blockExec *BlockExecutor) ResetDeliverState() {
	blockExec.proxyApp.SetOptionSync(abci.RequestSetOption{
		Key: "ResetDeliverState",
	})

}

//get lastBlock
func (blockExec *BlockExecutor) GetLastBlock() *types.Block {

	return blockExec.lastBlock
}

//reset lastBlock
func (blockExec *BlockExecutor) ResetLastBlock() {

	blockExec.lastBlock = nil
}

func IsCancelErr(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), CancelErr.Error())
}

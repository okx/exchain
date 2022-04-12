package state

import (
	"sync"
)

type asyncDBContext struct {
	// switch to turn on async save abciResponse and state
	isAsyncSaveDB bool
	// channel to write abciResponse async
	abciResponseQueue chan abciResponse
	/// channe to write state async
	stateQueue chan State
	// channel to feed back height of saved abci response and stat response
	asyncFeedbackQueue chan int64
	// flag to avoid waiting async state result for the first block
	isWaitingLastBlock bool
	//flag to avoid stop twice
	isAsyncQueueStop bool
	//wait group for quiting
	wg sync.WaitGroup
}

type abciResponse struct {
	height    int64
	responses *ABCIResponses
}

func (ctx *BlockExecutor) initAsyncDBContext() {
	ctx.abciResponseQueue = make(chan abciResponse)
	ctx.stateQueue = make(chan State)
	ctx.asyncFeedbackQueue = make(chan int64, 2)

	go ctx.asyncSaveStateRoutine()
	go ctx.asyncSaveABCIRespRoutine()
}

func (blockExec *BlockExecutor) stopAsyncDBContext() {
	if blockExec.isAsyncQueueStop {
		return
	}

	blockExec.wg.Add(2)
	close(blockExec.abciResponseQueue)
	close(blockExec.stateQueue)
	blockExec.wg.Wait()

	blockExec.isAsyncQueueStop = true
}

// ave the abciReponse async
func (blockExec *BlockExecutor) SaveABCIResponsesAsync(height int64, responses *ABCIResponses) {
	blockExec.abciResponseQueue <- abciResponse{height, responses}
}

// save the state async
func (blockExec *BlockExecutor) SaveStateAsync(state State) {
	blockExec.stateQueue <- state
}

// asyncSaveRoutine handle the write state work
func (blockExec *BlockExecutor) asyncSaveStateRoutine() {
	for stateMsg := range blockExec.stateQueue {
		SaveState(blockExec.db, stateMsg)
		blockExec.asyncFeedbackQueue <- stateMsg.LastBlockHeight
	}

	blockExec.wg.Done()
}

// asyncSaveRoutine handle the write abciResponse work
func (blockExec *BlockExecutor) asyncSaveABCIRespRoutine() {
	for abciMsg := range blockExec.abciResponseQueue {
		SaveABCIResponses(blockExec.db, abciMsg.height, abciMsg.responses)
		blockExec.asyncFeedbackQueue <- abciMsg.height
	}
	blockExec.wg.Done()
}

// switch to open async write db feature
func (blockExec *BlockExecutor) SetIsAsyncSaveDB(isAsyncSaveDB bool) {
	blockExec.isAsyncSaveDB = isAsyncSaveDB
}

// wait for the last sate and abciResponse to be saved
func (blockExec *BlockExecutor) tryWaitLastBlockSave(lastHeight int64) {
	if blockExec.isAsyncSaveDB && blockExec.isWaitingLastBlock {
		i := 0
		for r := range blockExec.asyncFeedbackQueue {
			if r != lastHeight {
				panic("Incorrect synced aysnc feed Height")
			}
			if i++; i == 2 {
				break
			}
		}
	}
}

// try to save the abciReponse async
func (blockExec *BlockExecutor) trySaveABCIResponsesAsync(height int64, abciResponses *ABCIResponses) {
	if blockExec.isAsyncSaveDB {
		blockExec.isWaitingLastBlock = true
		blockExec.SaveABCIResponsesAsync(height, abciResponses)
	} else {
		SaveABCIResponses(blockExec.db, height, abciResponses)
	}
}

// try to save the state async
func (blockExec *BlockExecutor) trySaveStateAsync(state State) {
	if blockExec.isAsyncSaveDB {
		blockExec.SaveStateAsync(state)
	} else {
		//Async save state
		SaveState(blockExec.db, state)
	}
}

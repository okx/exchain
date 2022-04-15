package state

import (
	"sync"
	"time"
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

const (
	MAXCHAN_LEN           = 2
	FEEDBACK_LEN          = 2
	QUIT_SIG              = -99
	MAX_WAIT_TIME_SECONDS = 30
)

type abciResponse struct {
	height    int64
	responses *ABCIResponses
}

func (ctx *BlockExecutor) initAsyncDBContext() {
	ctx.abciResponseQueue = make(chan abciResponse, MAXCHAN_LEN)
	ctx.stateQueue = make(chan State, MAXCHAN_LEN)
	ctx.asyncFeedbackQueue = make(chan int64, FEEDBACK_LEN)

	go ctx.asyncSaveStateRoutine()
	go ctx.asyncSaveABCIRespRoutine()
}

func (blockExec *BlockExecutor) stopAsyncDBContext() {
	if blockExec.isAsyncQueueStop {
		return
	}

	blockExec.wg.Add(2)
	blockExec.abciResponseQueue <- abciResponse{height: QUIT_SIG}
	blockExec.stateQueue <- State{LastBlockHeight: QUIT_SIG}

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
		if stateMsg.LastBlockHeight == QUIT_SIG {
			break
		}

		SaveState(blockExec.db, stateMsg)
		blockExec.asyncFeedbackQueue <- stateMsg.LastBlockHeight
	}

	blockExec.wg.Done()
}

// asyncSaveRoutine handle the write abciResponse work
func (blockExec *BlockExecutor) asyncSaveABCIRespRoutine() {
	for abciMsg := range blockExec.abciResponseQueue {
		if abciMsg.height == QUIT_SIG {
			break
		}
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
	timeoutCh := time.After(MAX_WAIT_TIME_SECONDS * time.Second)
	if blockExec.isAsyncSaveDB && blockExec.isWaitingLastBlock {
		i := 0
		for {
			select {
			case r := <-blockExec.asyncFeedbackQueue:
				if r != lastHeight {
					panic("Incorrect synced aysnc feed Height")
				}
				if i++; i == FEEDBACK_LEN {
					return
				}
			case <-timeoutCh:
				// It shouldn't be timeout. something must be wrong here
				panic("Can't get last block aysnc result")
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

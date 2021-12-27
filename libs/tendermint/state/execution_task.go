package state

import (
	"fmt"
	gorid "github.com/okex/exchain/libs/goroutine"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/automation"
	"github.com/okex/exchain/libs/tendermint/trace"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)


type executionResult struct {
	res *ABCIResponses
	err error
}

type executionTask struct {
	height  int64
	block   *types.Block
	stopped bool
	result  *executionResult

	prerunResultChan chan *executionTask
	proxyApp         proxy.AppConnConsensus
	db               dbm.DB
	logger           log.Logger
	index            int64
}

func newExecutionTask(blockExec *BlockExecutor, block *types.Block, index int64) *executionTask {

	return &executionTask{
		height:           block.Height,
		block:            block,
		db:               blockExec.db,
		proxyApp:         blockExec.proxyApp,
		logger:           blockExec.logger,
		prerunResultChan: blockExec.prerunCtx.prerunResultChan,
		index:            index,
	}
}

func (e *executionTask) dump(when string) {

	e.logger.Info(when,
		"gid", gorid.GoRId,
		"stopped", e.stopped,
		"Height", e.block.Height,
		"index", e.index,
		//"AppHash", e.block.AppHash,
	)
}

func (t *executionTask) stop() {
	if t.stopped {
		return
	}

	//reset deliverState
	if t.height != 1 {
		t.proxyApp.SetOptionSync(abci.RequestSetOption{Key: "ResetDeliverState",})
	}
	t.stopped = true
}


func (t *executionTask) run() {
	t.dump("Start prerun")
	trc := trace.NewTracer(fmt.Sprintf("num<%d>, lastRun", t.index))

	if t.height != 1 {
		t.proxyApp.SetOptionSync(abci.RequestSetOption{Key: "ResetDeliverState",})
	}

	abciResponses, err := execBlockOnProxyApp(t)

	if !t.stopped {
		t.result = &executionResult{
			abciResponses, err,
		}
		trace.GetElapsedInfo().AddInfo(trace.Prerun, trc.Format())
	}
	automation.PrerunTimeOut(t.block.Height, int(t.index)-1)
	t.dump("Prerun completed")
	t.prerunResultChan <- t
}

//========================================================
func (blockExec *BlockExecutor) InitPrerun() {
	if blockExec.deltaContext.downloadDelta {
		panic("download delta is not allowed if prerun enabled")
	}
	go blockExec.prerunCtx.prerunRoutine()
}

func (blockExec *BlockExecutor) NotifyPrerun(block *types.Block) {
	blockExec.prerunCtx.notifyPrerun(blockExec, block)
}
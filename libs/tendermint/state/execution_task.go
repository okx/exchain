package state

import (
	"fmt"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/trace"
	"sync/atomic"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

var err_delta_invoked = fmt.Errorf("prerrun stoped because of delta runnging")

type executionResult struct {
	res *ABCIResponses
	err error
}

type executionTask struct {
	height         int64
	index          int64
	block          *types.Block
	stopped        bool
	taskResultChan chan *executionTask
	result         *executionResult
	proxyApp       proxy.AppConnConsensus
	db             dbm.DB
	logger         log.Logger

	fetcher Fetcher
	notifyC  chan error
	// why: atomic is better, if we use mutex ,we have to define more variables
	status int32
}

func newExecutionTask(blockExec *BlockExecutor, block *types.Block, index int64, c Fetcher) *executionTask {

	return &executionTask{
		height:         block.Height,
		block:          block,
		db:             blockExec.db,
		proxyApp:       blockExec.proxyApp,
		logger:         blockExec.logger,
		taskResultChan: blockExec.prerunCtx.taskResultChan,
		index:          index,
		fetcher:       c,
		notifyC:        make(chan error, 1),
	}
}

func (e *executionTask) dump(when string) {

	e.logger.Info(when,
		"stopped", e.stopped,
		"Height", e.block.Height,
		"index", e.index,
		"hash", e.block.Hash(),
		//"AppHash", e.block.AppHash,
	)
}

func (t *executionTask) stop() {
	if t.stopped {
		return
	}

	//reset deliverState
	if t.height != 1 {
		t.proxyApp.SetOptionSync(abci.RequestSetOption{Key: "ResetDeliverState"})
	}
	t.stopped = true
}

func (t *executionTask) run() {
	defer func() {
		if nil != t.notifyC {
			close(t.notifyC)
		}
	}()

	var (
		abciResponses *ABCIResponses
		err           error
	)

	t.dump("Start prerun")
	trc := trace.NewTracer(fmt.Sprintf("num<%d>, lastRun", t.index))

	if t.height != 1 {
		t.proxyApp.SetOptionSync(abci.RequestSetOption{Key: "ResetDeliverState"})
	}
	curStatus := int32(TASK_BEGIN_PRERUN)
	traceHook(CASE_SPECIAL_BEFORE_LOAD_CACHE, t.status, emptyF)
	deltas, _ := t.fetcher.Fetch(t.block.Height)
	traceHook(CASE_SPECIAL_AFTER_LOAD_CACHE, t.status, emptyF)
	if deltas != nil {
		t.logger.Info("currentCacheDelta is not nil", "h", deltas.Height, "currentBlockHash", t.block.Hash())
		if !atomic.CompareAndSwapInt32(&t.status, 0, TASK_BEGIN_DELTA_EXISTS) {
			// case delta running
			traceHook(CASE_PRERRUNDELTA_SITUATION_GET_BEGIN_BLOCK_LOCK_FAILED, t.status, emptyF)
			t.logger.Info("prerun discard,because delta is running")
			return
		}
		traceHook(CASE_PRERRUNDELTA_SITUATION_GET_BEGIN_BLOCK_LOCK_SUCCESS, t.status, emptyF)
		curStatus = TASK_BEGIN_DELTA_EXISTS
		// delta  already downloaded
		execBlockOnProxyAppWithDeltas(t.proxyApp, t.block, t.db)
		resp := ABCIResponses{}
		err = types.Json.Unmarshal(deltas.ABCIRsp(), &resp)
		abciResponses = &resp
	} else {
		t.logger.Info("currentCacheDelta is  nil,so prerun try to execute", "currentBlockHash", t.block.Hash())
		if !atomic.CompareAndSwapInt32(&t.status, 0, TASK_BEGIN_PRERUN) {
			// case: delta get the beginBlock lock
			traceHook(CASE_PRERUN_SITUATION_GET_BEGIN_BLOCK_LOCK_FAILED, t.status, func() {
				// execute again
				execBlockOnProxyAppWithDeltas(t.proxyApp, t.block, t.db)
			})
			return
		}

		traceHook(CASE_PRERUN_SITUATION_GET_BEGIN_BLOCK_LOCK_SUCCESS, t.status, emptyF)

		abciResponses, err = execBlockOnProxyApp(t)
		if nil == err {
			curStatus = TASK_BEGIN_PRERUN
		} else if err == err_delta_invoked {
			// just finish
			traceHook(CASE_PRERRUN_CANCELED_BY_DELTA, t.status, emptyF)
			return
		}
	}
	traceHook(CASE_SPECIAL_PRERUN_BEFORE_RACE_END, 0, emptyF)
	notifyResult(t, abciResponses, err, func() {
		traceHook(CASE_PRERRUN_SITUATION_RACE_END_SUCCESS, t.status, emptyF)
	}, func() {
		traceHook(CASE_PRERRUN_SITUATION_RACE_END_FAIL, t.status, emptyF)
	}, trc, curStatus, TASK_DELTA, TASK_PRERRUN)
}

//========================================================
func (blockExec *BlockExecutor) InitPrerun() {
	if blockExec.deltaContext.downloadDelta {
		go blockExec.prerunCtx.consume()
	}
	go blockExec.prerunCtx.prerunRoutine()
}

func (blockExec *BlockExecutor) NotifyPrerun(block *types.Block) {
	blockExec.prerunCtx.notifyPrerun(blockExec, block)
}

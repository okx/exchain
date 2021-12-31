package state

import (
	"context"
	"fmt"
	gorid "github.com/okex/exchain/libs/goroutine"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/automation"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/trace"
	"github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

// TODO: its better to use interface#Execute or pipeline funcs(preExecute,postExecute to customize and it will be easy to refactor)
type TaskExecutor func(context.Context, *executionTask, chan<- *executionResult)
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

	executor       TaskExecutor
	listener       <-chan interface{}
	listenerCancel func()

	proxyApp proxy.AppConnConsensus
	db       dbm.DB
	logger   log.Logger
}

func newExecutionTask(blockExec *BlockExecutor, block *types.Block, index int64) *executionTask {
	ret := &executionTask{
		height:         block.Height,
		block:          block,
		db:             blockExec.db,
		proxyApp:       blockExec.proxyApp,
		logger:         blockExec.logger,
		taskResultChan: blockExec.prerunCtx.taskResultChan,
		index:          index,
		executor:       blockExec.prerunCtx.taskExecutor,
	}
	ret.listener, ret.listenerCancel = registerListenerBeforePrerun(blockExec, block.Height)
	return ret
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
		t.proxyApp.SetOptionSync(abci.RequestSetOption{Key: "ResetDeliverState"})
	}
	t.stopped = true
}

func (t *executionTask) run() {
	t.dump("Start prerun")
	trc := trace.NewTracer(fmt.Sprintf("num<%d>, lastRun", t.index))

	if t.height != 1 {
		t.proxyApp.SetOptionSync(abci.RequestSetOption{Key: "ResetDeliverState"})
	}

	notifyC := make(chan *executionResult, 1)
	t.executor(context.Background(), t, notifyC)
	ret := <-notifyC
	abciResponses, err := ret.res, ret.err
	//abciResponses, err := t.executor(t)

	if !t.stopped {
		t.result = &executionResult{
			abciResponses, err,
		}
		trace.GetElapsedInfo().AddInfo(trace.Prerun, trc.Format())
	}
	automation.PrerunTimeOut(t.block.Height, int(t.index)-1)
	t.dump("Prerun completed")
	t.taskResultChan <- t
}

func (t *executionTask) listenerDone() *types.Deltas {
	select {
	case v, ok := <-t.listener:
		if !ok {
			return nil
		}
		return v.(*types.Deltas)
	default:
	}
	return nil
}

//========================================================
func (blockExec *BlockExecutor) InitPrerun() {
	//if blockExec.deltaContext.downloadDelta {
	//	panic("download delta is not allowed if prerun enabled")
	//}
	go blockExec.prerunCtx.prerunRoutine()
}

func (blockExec *BlockExecutor) NotifyPrerun(block *types.Block) {
	blockExec.prerunCtx.notifyPrerun(blockExec, block)
}

//============================================
// cases:
//  edge case:
// 			1.deltaDownload the delta => prerrun begin to run => prerun quit without execute execBlockOnProxyAppWithDeltas
//			2.prerun running => prerrun finished ,listener still hang on =>
//				prerun cancel the listener(why:
//				because if delta cant download the specific delta all the time,listener will have memory leak )
// 		 	3.deltaDownloading the delta => prerrun running =>
//		 	listener receive the delta (means deltaDownloaded the data) => quit(execBlockOnProxyAppWithDeltas still running)

func ConcurrentWithPreExecutor(listenerF func(),executor ...TaskExecutor) TaskExecutor {
	return func(ctx context.Context, task *executionTask, results chan<- *executionResult) {
		data := task.listenerDone()
		if data != nil && fillResultByDelta(task, data, results) {
			// edge case 1
			listenerF()
			return
		}
		ConcurrentTaskExecutor(executor...)(ctx,task,results)
	}
}

func ConcurrentTaskExecutor(executors ...TaskExecutor) TaskExecutor {
	return func(ctx context.Context, task *executionTask, results chan<- *executionResult) {
		// FIXME , ITS BETTER TO USE ROUTINE POOL(not mpg)
		for i := 0; i < len(executors); i++ {
			go func(index int) {
				executors[index](ctx, task, results)
			}(i)
		}
	}
}

func RunOnListener(listenerH func()) TaskExecutor {
	return func(ctx context.Context, t *executionTask, results chan<- *executionResult) {
		select {
		case v, ok := <-t.listener:
			if !ok {
				return
			}
			if fillResultByDelta(t, v.(*types.Deltas), results) {
				listenerH()
				// edge case 3
				// TODO when we receive the delta immediately,we should cancel the execBlockOnProxyApp as possibal as we can(by context)
			}
		}
	}
}

func RunOnProxyAppWithPrePost(pre, proxyExecuteH func()) TaskExecutor {
	return func(ctx context.Context, task *executionTask, results chan<- *executionResult) {
		pre()
		RunOnProxyApp(proxyExecuteH)(ctx,task,results)
	}
}

func RunOnProxyApp(proxyExecuteH func()) TaskExecutor {
	return func(ctx context.Context, t *executionTask, results chan<- *executionResult) {
		ret, err := execBlockOnProxyApp(t)
		select {
		case results <- &executionResult{res: ret, err: err}:
			// edge case 2:
			t.listenerCancel()
			proxyExecuteH()
		default:
		}
	}
}

func fillResultByDelta(t *executionTask, data *types.Deltas, ret chan<- *executionResult) bool {
	execBlockOnProxyAppWithDeltas(t.proxyApp, t.block, t.db)
	resp := ABCIResponses{}
	err := types.Json.Unmarshal(data.ABCIRsp(), &resp)
	select {
	case ret <- &executionResult{res: &resp, err: err}:
		return true
	default:
		return false
	}
}

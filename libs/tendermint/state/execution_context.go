package state

import (
	"bytes"
	"fmt"
	"github.com/okex/exchain/libs/queue"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/trace"
	"github.com/okex/exchain/libs/tendermint/types"
	"sync"
	"sync/atomic"
)

var (
	traceHook func(trace, status int32, cb func()) = func(_, _ int32, _ func()) {}
	emptyF    func()
)

func SetTraceHook(f func(trace, status int32, cb func())) {
	traceHook = f
}

type Fetcher interface {
	Fetch(height int64) (*types.Deltas, int64)
}

type DeltaJob struct {
	Delta *types.Deltas
	Cb    func(h int64)
}

type PreRunContextOption func(ctx *prerunContext)

func PreRunContextWithQueue(q queue.Queue) PreRunContextOption {
	return func(ctx *prerunContext) {
		ctx.consumerQ = q
	}
}

func PreRunContextWithFetcher(deltaMap Fetcher) PreRunContextOption {
	return func(ctx *prerunContext) {
		ctx.fetcher = deltaMap
	}
}

func PreRunWithFastSyncCheck(f func() bool) PreRunContextOption {
	return func(ctx *prerunContext) {
		ctx.isFastSync = f
	}
}

type prerunContext struct {
	prerunTx        bool
	taskChan        chan *executionTask
	taskResultChan  chan *executionTask
	prerunTask      *executionTask
	logger          log.Logger
	consumerQ       queue.Queue
	fetcher         Fetcher
	lastPruneHeight int64
	isFastSync      func() bool
}

func newPrerunContex(logger log.Logger, ops ...PreRunContextOption) *prerunContext {
	ret := &prerunContext{
		taskChan:       make(chan *executionTask, 1),
		taskResultChan: make(chan *executionTask, 1),
		logger:         logger,
	}
	for _, opt := range ops {
		opt(ret)
	}

	return ret
}

func (pc *prerunContext) init(f Fetcher, q queue.Queue, e *BlockExecutor) {
	if pc.consumerQ == nil {
		pc.consumerQ = q
	}
	if pc.isFastSync == nil {
		pc.isFastSync = func() bool {
			return e.isFastSync
		}
	}
	if pc.fetcher == nil {
		pc.fetcher = f
	}
	if types.PreRunConsumeDebugEnable {
		SetTraceHook(func(trace, status int32, cb func()) {
			executor := "[special | result]"
			if trace&TRACE_PRERUN_WITH_CACHE >= TRACE_PRERUN_WITH_CACHE {
				executor = "[PRERRUN_WITH_CACHE]"
			} else if trace&TRACE_DELTA >= TRACE_DELTA {
				executor = "[DELTA]"
			} else if trace&TRACE_PRERUN_WITH_NO_CACHE > +TRACE_PRERUN_WITH_NO_CACHE {
				executor = "[TRACE_PRERUN_WITH_NO_CACHE]"
			}
			e.prerunCtx.logger.Info("traceHook", "executor", executor, "trace", trace, "status", status)
		})
	}
}

func (pc *prerunContext) checkIndex(height int64) {
	var index int64
	if pc.prerunTask != nil {
		index = pc.prerunTask.index
	}
	pc.logger.Info("Not apply delta", "height", height, "prerunIndex", index)

}

func (pc *prerunContext) flushPrerunResult() {
	for {
		select {
		case task := <-pc.taskResultChan:
			task.dump("Flush prerun result")
		default:
			return
		}
	}
}

func (pc *prerunContext) prerunRoutine() {
	pc.prerunTx = true
	for task := range pc.taskChan {
		task.run()
	}
}

func (pc *prerunContext) consume() {
	var (
		v interface{}
	)
	pc.logger.Info("prerrun context ,consumer start up")

	for {
		v = pc.consumerQ.Take()
		pc.handleMsg(v)
	}
}

func (pc *prerunContext) handleMsg(msg interface{}) {
	switch v := msg.(type) {
	case *DeltaJob:
		pc.handleDeltaMsg(v)
	default:
		panic("programa error")
	}
}

func (pc *prerunContext) handleDeltaMsg(v *DeltaJob) {
	if pc.isFastSync() {
		return
	}
	pc.logger.Info("receive deltaMsgJob,", "height", v.Delta.Height)
	delta := v.Delta
	traceHook(CASE_SPECIAL_DELTA_BEFORE_BEGIN, 0, emptyF)

	// hold the pointer at first
	curTask := pc.prerunTask
	// in case of  producer is too faster than consumer
	if curTask == nil && delta.Height > pc.lastPruneHeight && delta.Height < pc.lastPruneHeight+10 {
		return
	}

	if curTask == nil {
		pc.logger.Info("currentTask is nil,discard current deltaJob")
		return
	}
	curBlock := curTask.block
	curDb := curTask.db
	app := curTask.proxyApp
	// best effort: try to avoid  running twice
	if curBlock.Height > delta.Height || curTask.stopped {
		// ignore
		return
	} else if curBlock.Height < delta.Height {
		if delta.Height > pc.lastPruneHeight+10 {
			return
		}
		return
	}
	traceHook(CASE_SPECIAL_DELTA_BEFORE_FINAL_STORE, curTask.status, emptyF)

	trc := trace.NewTracer(fmt.Sprintf("num<%d>, lastRun", curTask.index))

	abciResponses := ABCIResponses{}
	err := types.Json.Unmarshal(delta.ABCIRsp(), &abciResponses)
	curStatus := int32(TASK_BEGIN_DELTA)
	if !atomic.CompareAndSwapInt32(&curTask.status, 0, TASK_BEGIN_DELTA) {
		loadStatus := atomic.LoadInt32(&curTask.status)
		if loadStatus&TASK_BEGIN_DELTA_EXISTS >= TASK_BEGIN_DELTA_EXISTS {
			traceHook(CASE_DELTA_SITUATION_BEGIN_BLOCK_FAILED_FOR_DELTA_ALREADY_EXISTS_IN_CACHE_BY_PRERRUN, loadStatus, emptyF)
			pc.logger.Info("currentTask has been execute by prerun_cache,discard", "height", curBlock.Height)
			return
		}
		// edge case 2 ,prerrun already finished the task
		if loadStatus&TASK_PRERRUN >= TASK_PRERRUN {
			pc.logger.Info("currentTask has been executed by prerun,discard", "height", curBlock.Height)
			traceHook(CASE_DELTA_SITUATION_BEGIN_BLOCK_FAILED_FOR_TASK_FINISHED_BY_PRERRUN, loadStatus, emptyF)
			return
		}

		// which means prerrun routine get the lock ,wait it until prerrun finish the beginblock
		// and we have to make sure ,it returns success
		v := <-curTask.notifyC
		if v != nil {
			traceHook(CASE_DELTA_SITUATION_BEGIN_BLOCK_FAILED_AND_NOTIFIED_BY_PRERRUN, loadStatus, emptyF)
			pc.logger.Error("prerrun run the task faield", "err", v)
			return
		}
		traceHook(CASE_DELTA_SITUATION_BEGIN_BLOCK_FAILED_AND_NOTIFIED_BY_PRERRUN, loadStatus, func() { execBlockOnProxyAppWithDeltas(app, curBlock, curDb) })
		curStatus = TASK_BEGIN_PRERUN
	} else {
		pc.logger.Info("handleDelta,executeBeginBlock")
		traceHook(CASE_DELTA_SITUATION_GET_BEGIN_BLOCK_LOCK_SUCCESS, curTask.status, emptyF)
		execBlockOnProxyAppWithDeltas(app, curBlock, curDb)
	}

	traceHook(CASE_SPECIAL_DELTA_BEFORE_RACE_END, 0, emptyF)
	notifyResult(curTask, &abciResponses, err, func() {
		//  case: delta finished before prerrun
		pc.logger.Info("waiting  prerun canceled")
		traceHook(CASE_DELTA_SITUATION_RACE_END_SUCCESS, curTask.status, emptyF)
		for {
			select {
			case _, ok := <-curTask.notifyC:
				if !ok {
					pc.logger.Info("prerun canceled successfully")
					execBlockOnProxyAppWithDeltas(app, curBlock, curDb)
					traceHook(CASE_DELTA_ENTER_CHAN_RECEIVE_WAIT_PRERRUN_CLOSE_NOTIFY, curTask.status, emptyF)
					return
				}
			}
		}
	}, func() {
		traceHook(CASE_DELTA_SITUATION_RACE_END_FAIL, curTask.status, emptyF)
	}, trc, curStatus, TASK_PRERRUN, TASK_DELTA)
}
func notifyResult(curTask *executionTask,
	abciResponses *ABCIResponses,
	err error,
	hook func(),
	elseH func(),
	trc *trace.Tracer,
	currentStatus int32,
	invalidStatus,
	deltaStatus int32) {
	if !atomic.CompareAndSwapInt32(&curTask.status, currentStatus, currentStatus|deltaStatus) {
		return
	}

	currentStatus = curTask.status
	if currentStatus&invalidStatus >= invalidStatus {
		//curTask.status &^= deltaStatus
		panic("programa error")
	}

	curTask.result = &executionResult{res: abciResponses, err: err}
	// finally we can try to cancel the prerun if we have to
	hook()
	select {
	case curTask.taskResultChan <- curTask:
		traceHook(RESULT_TRAEC, currentStatus, emptyF)
		curTask.dump(fmt.Sprintf("curTaskFinished,deltaStatus=%d,currentTaskStatus=%d", deltaStatus,curTask.status))
		trace.GetElapsedInfo().AddInfo(trace.Prerun, trc.Format())
	default:
		panic("programa error")
	}
}

func store(job *DeltaJob, cache *sync.Map) {
	cache.Store(job.Delta.Height, job.Delta)
	if nil != job.Cb {
		job.Cb(job.Delta.Height)
	}
}

func (pc *prerunContext) dequeueResult(b *types.Block) (*ABCIResponses, error) {
	expected := pc.prerunTask
	for context := range pc.taskResultChan {

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
			if !bytes.Equal(b.Hash(), context.block.Hash()) {
				panic("asdkjdkdkdk")
			}
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

	pc.prerunTask = newExecutionTask(blockExec, block, stoppedIndex, pc.fetcher)

	pc.prerunTask.dump("Notify prerun")

	// start a new one
	pc.taskChan <- pc.prerunTask
}

func (pc *prerunContext) getPrerunResult(b *types.Block, height int64, fastSync bool) (res *ABCIResponses, err error) {

	pc.checkIndex(height)
	if fastSync {
		pc.stopPrerun(height)
		return
	}
	// blockExec.prerunContext == nil means:
	// 1. prerunTx disabled
	// 2. we are in fasy-sync: the block comes from BlockPool.AddBlock not State.addProposalBlockPart and no prerun result expected
	if pc.prerunTask != nil {
		res, err = pc.dequeueResult(b)
		pc.prerunTask = nil
	}
	return
}

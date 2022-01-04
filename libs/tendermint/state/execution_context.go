package state

import (
	"bytes"
	"fmt"
	"github.com/okex/exchain/libs/queue"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
	"sync"
	"sync/atomic"
)

const ()

var (
	traceHook func(trace, status int32) = func(_, _ int32) {}
)

func SetTraceHook(f func(trace, status int32)) {
	traceHook = f
}

type PruneCacheJob struct {
	h int64
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

type prerunContext struct {
	prerunTx        bool
	taskChan        chan *executionTask
	taskResultChan  chan *executionTask
	prerunTask      *executionTask
	logger          log.Logger
	consumerQ       queue.Queue
	cache           *sync.Map
	lastPruneHeight int64
}

func newPrerunContex(logger log.Logger, ops ...PreRunContextOption) *prerunContext {
	ret := &prerunContext{
		taskChan:       make(chan *executionTask, 1),
		taskResultChan: make(chan *executionTask, 1),
		logger:         logger,
		cache:          &sync.Map{},
	}

	for _, opt := range ops {
		opt(ret)
	}

	return ret
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

// note: if tendermint network failed all the time ,prerunContext#cache#memory will keep  growing
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
	case PruneCacheJob:
		pc.handlePrune(v)
	case *DeltaJob:
		pc.handleDeltaMsg(v)
	default:
		panic("programa error")
	}
}

func (pc *prerunContext) handleDeltaMsg(v *DeltaJob) {
	delta := v.Delta
	// TODO cmd line

	traceHook(CASE_SPECIAL_DELTA_BEFORE_BEGIN, 0)

	// hold the pointer at first
	curTask := pc.prerunTask
	// in case of  producer is too faster than consumer
	if curTask == nil && delta.Height > pc.lastPruneHeight && delta.Height < pc.lastPruneHeight+10 {
		// cache the delta
		store(v, pc.cache)
		return
	}

	if curTask==nil{
		return
	}
	curBlock := curTask.block
	curDb := curTask.db
	app := curTask.proxyApp

	// best effort: try to avoid to run twice
	if curBlock.Height > delta.Height || curTask.stopped {
		// ignore
		return
	} else if curBlock.Height < delta.Height {
		// cache the delta
		store(v, pc.cache)
		return
	}

	traceHook(CASE_SPECIAL_DELTA_BEFORE_FINAL_STORE, curTask.status)

	// store the delta:  in case of  `cpu edge case`
	// we will delete the height by prune
	store(v, pc.cache)

	abciResponses := ABCIResponses{}
	err := types.Json.Unmarshal(delta.ABCIRsp(), &abciResponses)
	curStatus := int32(TASK_BEGIN_DELTA)
	if !atomic.CompareAndSwapInt32(&curTask.status, 0, TASK_BEGIN_DELTA) {
		loadStatus := atomic.LoadInt32(&curTask.status)
		if loadStatus&TASK_BEGIN_DELTA_EXISTS >= TASK_BEGIN_DELTA_EXISTS {
			traceHook(CASE_DELTA_SITUATION_BEGIN_BLOCK_FAILED_FOR_DELTA_ALREADY_EXISTS_IN_CACHE_BY_PRERRUN, loadStatus)
			return
		}
		// edge case 2 ,prerrun already finished the task
		if loadStatus&TASK_PRERRUN >= TASK_PRERRUN {
			traceHook(CASE_DELTA_SITUATION_BEGIN_BLOCK_FAILED_FOR_TASK_FINISHED_BY_PRERRUN, loadStatus)
			return
		}

		// which means prerrun routine get the lock ,wait it until prerrun finish the beginblock
		// and we have to make sure ,it returns success
		v := <-curTask.notifyC
		if v != nil {
			traceHook(CASE_DELTA_SITUATION_BEGIN_BLOCK_FAILED_AND_NOTIFIED_BY_PRERRUN, loadStatus)
			pc.logger.Error("prerrun run the task faield", "err", v)
			return
		}
		traceHook(CASE_DELTA_SITUATION_BEGIN_BLOCK_FAILED_AND_NOTIFIED_BY_PRERRUN, loadStatus)
		curStatus = TASK_BEGIN_PRERUN
	} else {
		//
		traceHook(CASE_DELTA_SITUATION_GET_BEGIN_BLOCK_LOCK_SUCCESS, curTask.status)
		execBlockOnProxyAppWithDeltas(app, curBlock, curDb)
	}

	traceHook(CASE_SPECIAL_DELTA_BEFORE_RACE_END, 0)
	notifyResult(curTask, &abciResponses, err, func() {
		//  case: delta finished before prerrun
		traceHook(CASE_DELTA_SITUATION_RACE_END_SUCCESS, curTask.status)
		for {
			select {
			case _, ok := <-curTask.notifyC:
				if !ok {
					traceHook(CASE_DELTA_ENTER_CHAN_RECEIVE_WAIT_PRERRUN_CLOSE_NOTIFY, curTask.status)
					return
				}
			}
		}
	}, func() {
		traceHook(CASE_DELTA_SITUATION_RACE_END_FAIL,curTask.status)
	}, curStatus, TASK_PRERRUN, TASK_DELTA)
}
func notifyResult(curTask *executionTask,
	abciResponses *ABCIResponses,
	err error,
	hook func(),
	elseH func(),
	currentStatus int32,
	invalidStatus,
	deltaStatus int32) {
	if !atomic.CompareAndSwapInt32(&curTask.status, currentStatus, currentStatus|deltaStatus) {
		return
	}

	// TODO ,useless code blcok
	currentStatus = curTask.status
	if currentStatus&invalidStatus >= invalidStatus {
		curTask.status &^= deltaStatus
		return
	}

	curTask.result = &executionResult{res: abciResponses, err: err}
	// finally we can try to cancel the prerun if we have to
	hook()
	select {
	case curTask.taskResultChan <- curTask:
		traceHook(RESULT_TRAEC, currentStatus)
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

func (pc *prerunContext) dequeueResult() (*ABCIResponses, error) {
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
			// push prune job
			pc.consumerQ.Push(PruneCacheJob{h: expected.block.Height})
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

	pc.prerunTask = newExecutionTask(blockExec, block, stoppedIndex, pc.cache)

	pc.prerunTask.dump("Notify prerun")

	// start a new one
	pc.taskChan <- pc.prerunTask
}

func (pc *prerunContext) getPrerunResult(height int64, fastSync bool) (res *ABCIResponses, err error) {

	pc.checkIndex(height)

	if fastSync {
		pc.stopPrerun(height)
		return
	}
	// blockExec.prerunContext == nil means:
	// 1. prerunTx disabled
	// 2. we are in fasy-sync: the block comes from BlockPool.AddBlock not State.addProposalBlockPart and no prerun result expected
	if pc.prerunTask != nil {
		res, err = pc.dequeueResult()
		pc.prerunTask = nil
	}
	return
}

func (pc *prerunContext) handlePrune(v PruneCacheJob) {
	if v.h < pc.lastPruneHeight {
		return
	}
	h := v.h
	for i := h; i > pc.lastPruneHeight; i-- {
		pc.cache.Delete(i)
	}
	pc.lastPruneHeight = h
}

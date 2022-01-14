package state

import (
	"bytes"
	"fmt"
	"github.com/okex/exchain/libs/queue"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/trace"
	"github.com/okex/exchain/libs/tendermint/types"
	"sync/atomic"
)

type IAcquire interface {
	acquire(height int64) (*types.Deltas, int64)
}

type DeltaJob struct {
	Delta *types.Deltas
	Cb    func(h int64)
}

const (
	TaskBeginByDelta                 = 1 << 0
	TaskBeginByPrerunWithNoCache     = 1 << 1
	TaskBeginByPrerunWithCacheExists = 1 << 2

	TaskEndByDelta  = 1 << 3
	TaskEndByPrerun = 1 << 4
)

type PreRunContextOption func(ctx *prerunContext)

func PreRunContextWithQueue(q queue.Queue) PreRunContextOption {
	return func(ctx *prerunContext) {
		ctx.consumerQ = q
	}
}

func PreRunContextWithFetcher(deltaMap IAcquire) PreRunContextOption {
	return func(ctx *prerunContext) {
		ctx.acquire = deltaMap
	}
}

func PreRunWithFastSyncCheck(f func() bool) PreRunContextOption {
	return func(ctx *prerunContext) {
		ctx.isFastSync = f
	}
}

type prerunContext struct {
	prerunTx       bool
	taskChan       chan *executionTask
	taskResultChan chan *executionTask
	prerunTask     *executionTask
	logger         log.Logger
	consumerQ      queue.Queue
	acquire        IAcquire
	isFastSync     func() bool
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

func (pc *prerunContext) init(f IAcquire, q queue.Queue, e *BlockExecutor) {
	if pc.consumerQ == nil {
		pc.consumerQ = q
	}
	if pc.isFastSync == nil {
		pc.isFastSync = func() bool {
			return e.isFastSync
		}
	}
	if pc.acquire == nil {
		pc.acquire = f
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
	delta := v.Delta

	// hold the pointer at first
	curTask := pc.prerunTask
	if curTask == nil {
		return
	}
	curBlock := curTask.block
	curDb := curTask.db
	app := curTask.proxyApp
	// best effort: try to avoid  running twice
	if curBlock.Height != delta.Height || curTask.stopped {
		curTask.dump(fmt.Sprintf("delta: discard current delta,deltaHeight=%d", delta.Height))
		// ignore
		return
	}
	if !delta.Validate(curBlock.Height) {
		curTask.dump(fmt.Sprintf("delta, curent delta is invalid,height=%d", v.Delta.Height))
		return
	}

	trc := trace.NewTracer(fmt.Sprintf("num<%d>, lastRun", curTask.index))
	abciResponses := ABCIResponses{}
	err := types.Json.Unmarshal(delta.ABCIRsp(), &abciResponses)

	if !atomic.CompareAndSwapInt32(&curTask.status, 0, TaskBeginByDelta) {
		loadStatus := atomic.LoadInt32(&curTask.status)
		// case1: task executed by prerun_with_cache(maybe not done yet,but we dont care)
		if loadStatus&TaskBeginByPrerunWithCacheExists >= TaskBeginByPrerunWithCacheExists {
			curTask.dump("current task has been execute by prerun_delta_cache,discard")
			return
		}

		// case2: task was finished by prerun(at this time)
		if loadStatus&TaskEndByPrerun >= TaskEndByPrerun {
			curTask.dump("current task has been executed by prerun,discard")
			return
		}
		curTask.dump("delta is waitting prerun to be canceld or finished")

		// which means prerrun routine is running,we will execute again if the task hasnt done yet(delta's priority > prerun's property)
		// before we execute beginBlock ,we have to cancel prerun at first(and we have to wait,because deliverTx will affect deliverState)
		// and we cant  guarantee 'prerun' is done with endBlock(havent notify result yet)  or it is still running deliverTx
		// so we have to use cas instead of using  StoreInt32
		// note: we dont care about cas result(cas is just try to cancel the deliverTx step)
		atomic.CompareAndSwapInt32(&curTask.status, loadStatus, TaskBeginByDelta)
		<-curTask.notifyC
		// which means :
		// prerun is quit(but we dont know the task is done or canceld)
		// note: if prerun is done ,we cant execute beginBlock again ,because if the function `dequeueResult`
		// 		 is executed ,and immediately we call beginBlock again before BlockExecutor#commit `data will reset`
		//		 so we have to check again with current staus
		if atomic.LoadInt32(&curTask.status)&TaskEndByPrerun >= TaskEndByPrerun {
			curTask.dump("current task has been executed by prerun,discard")
			return
		}
		// case3 prerun is canceled,so we can handle it again
	} else {
		// case: 2 blocks ,see executionTask#stop
		defer func() {
			if nil != curTask.notifyC {
				close(curTask.notifyC)
			}
		}()
	}

	// we run here
	// means
	// prerun is quit:
	//			1. prerun donest execute at all
	//			2. prerun canceled
	curTask.dump("start beginBlock by delta ")
	// execute again
	execBlockOnProxyAppWithDeltas(app, curBlock, curDb)

	notifyResult(curTask, &abciResponses, err, TaskBeginByDelta|TaskEndByDelta, trc)
}
func notifyResult(curTask *executionTask,
	abciResponses *ABCIResponses,
	err error,
	lastStatus int32,
	trc *trace.Tracer) {

	// before we end the task ,we update the status at first
	atomic.StoreInt32(&curTask.status, lastStatus)

	if !curTask.stopped {
		curTask.result = &executionResult{res: abciResponses, err: err}
		trace.GetElapsedInfo().AddInfo(trace.Prerun, trc.Format())
	}

	select {
	case curTask.taskResultChan <- curTask:
		curTask.dump(fmt.Sprintf("current task finished, final task status=%d", lastStatus))
	default:
		// edge case : 2 blocks ,we cant let it panic
		// panic("program error")
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

	pc.prerunTask = newExecutionTask(blockExec, block, stoppedIndex, pc.acquire)

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

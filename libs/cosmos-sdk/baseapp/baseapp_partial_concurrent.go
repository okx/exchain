package baseapp

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"sync"
	"time"
)

const (
	maxDeliverTxsConcurrentNum = 4
)

var totalAnteDuration = int64(0)
var totalGasAndMsgsDuration = int64(0)
var totalSerialDuration = int64(0)
var totalSavedTime = int64(0)
var totalWriteTime = int64(0)
var totalDeferGasTime = int64(0)
var totalHandleGasTime = int64(0)
var totalRunMsgsTime = int64(0)
var totalFinishTime = int64(0)

type DeliverTxTask struct {
	tx            sdk.Tx
	index         int
	feeForCollect int64
	anteFailed    bool
	info          *runTxInfo
	from          adb.Address
	fee           sdk.Coins
	//isEvm         bool
	//signCache     sdk.SigCache
	//evmIndex      uint32
	err           error
	//decodeFailed  bool
}

func newDeliverTxTask(tx sdk.Tx, index int) *DeliverTxTask {
	t := &DeliverTxTask{
		tx:    tx,
		index: index,
		info:  &runTxInfo{tx: tx},
	}

	return t
}

type DeliverTxTasksManager struct {
	done       chan int // done for all transactions are executed
	nextSignal chan int // signal for taking a new tx into tasks
	executeSignal chan int // signal for taking a new task from pendingTasks to executingTask
	waitingCount     int
	executeSignalCount int
	mtx           sync.Mutex

	totalCount    int
	curIndex      int
	tasks         sync.Map
	pendingTasks  sync.Map
	executingTask *DeliverTxTask

	txResponses []*abci.ResponseDeliverTx

	app                *BaseApp
	anteDuration       int64
	gasAndMsgsDuration int64
	serialDuration     int64
	writeDuration      int64
	deferGasTime       int64
	handleGasTime      int64
	runMsgsTime        int64
	finishTime         int64
}

func NewDeliverTxTasksManager(app *BaseApp) *DeliverTxTasksManager {
	return &DeliverTxTasksManager{
		app: app,
	}
}

func (dm *DeliverTxTasksManager) deliverTxs(txs [][]byte) {
	dm.done = make(chan int, 1)
	dm.nextSignal = make(chan int, 1)
	dm.executeSignal = make(chan int, 1)
	dm.executeSignalCount = 0
	dm.waitingCount = 0

	dm.totalCount = len(txs)
	dm.curIndex = -1

	dm.tasks = sync.Map{}
	dm.pendingTasks = sync.Map{}
	dm.txResponses = make([]*abci.ResponseDeliverTx, len(txs))
	
	dm.anteDuration = 0
	dm.gasAndMsgsDuration = 0
	dm.serialDuration = 0
	dm.writeDuration = 0
	dm.deferGasTime = 0
	dm.handleGasTime = 0
	dm.runMsgsTime = 0
	dm.finishTime = 0

	go dm.makeTasksRoutine(txs)
	go dm.runTxSerialRoutine()
}

func (dm *DeliverTxTasksManager) makeTasksRoutine(txs [][]byte) {
	taskIndex := 0
	for {
		if taskIndex == dm.totalCount {
			break
		}

		////remaining := taskIndex - (dm.getCurIndex() + 1) //- numTasks - numPending
		//switch {
		////case remaining >= maxDeliverTxsConcurrentNum:
		//case dm.isWaiting(true):
		//	<-dm.nextSignal
		//	dm.executeSignalCount--
		//	if dm.executeSignalCount < 0 {
		//		dm.app.logger.Error("dm.executeSignalCount < 0", "count", dm.executeSignalCount)
		//	}
		//
		//default:
			dm.makeNextTask(txs[taskIndex], taskIndex)
			taskIndex++
			dm.incrementWaitingCount(true)
		//}
	}
}

func (dm *DeliverTxTasksManager) makeNextTask(tx []byte, index int) {
	go dm.runTxPartConcurrent(tx, index)
}

func (dm *DeliverTxTasksManager) runTxPartConcurrent(txByte []byte, index int) {
	start := time.Now()
	// create a new task
	task := dm.makeNewTask(txByte, index)

	defer dm.pushIntoPending(task)

	if task.err != nil {
		return
	}

	mode := runTxModeDeliverPartConcurrent
	task.info.handler = dm.app.getModeHandler(mode) //dm.handler

	// execute ante
	task.info.ctx = dm.app.getContextForTx(mode, task.info.txBytes) // same context for all txs in a block
	//task.fee, task.isEvm, task.signCache = dm.app.getTxFee(task.info.ctx, task.tx)

	task.info.ctx = task.info.ctx.WithCache(sdk.NewCache(dm.app.blockCache, useCache(mode))) // one cache for a tx

	// dm.app.pin(ValTxMsgs, true, mode)
	if err := validateBasicTxMsgs(task.tx.GetMsgs()); err != nil {
		task.err = err
		dm.app.logger.Error("validateBasicTxMsgs failed", "err", err)
		return
	}
	// dm.app.pin(ValTxMsgs, false, mode)

	// dm.app.pin(RunAnte, true, mode)
	if dm.app.anteHandler != nil {
		err := dm.runAnte(task.info, mode)
		if err != nil {
			//dm.app.logger.Error("runAnte failed", "err", err)
			task.anteFailed = true
		}
	}
	// dm.app.pin(RunAnte, false, mode)

	elapsed := time.Since(start).Microseconds()
	dm.anteDuration += elapsed
}

func (dm *DeliverTxTasksManager) makeNewTask(txByte []byte, index int) *DeliverTxTask {
	//dm.app.logger.Info("runTxPartConcurrent", "index", index)
	tx, err := dm.app.txDecoder(txByte)
	task := newDeliverTxTask(tx, index)
	task.info.txBytes = txByte
	if err != nil {
		task.err = err
		//task.decodeFailed = true
		dm.app.logger.Error("tx decode failed"," err", err)
	}

	dm.tasks.Store(task.index, task)
	return task
}

// put task into pendingTasks after execution finished
func (dm *DeliverTxTasksManager) pushIntoPending(task *DeliverTxTask) {
	if task == nil {
		return
	}

	//dm.app.logger.Info("new into pendingTasks", "index", task.index)
	dm.pendingTasks.Store(task.index, task)
	if dm.executingTask == nil && task.index == dm.curIndex+1 {
		dm.executeSignal <- 0
	}
	dm.tasks.Delete(task.index)
}

func (dm *DeliverTxTasksManager) runAnte(info *runTxInfo, mode runTxMode) error {
	//fmt.Printf("runAnte start. gasWanted:%d, startingGas:%d, gasMeter:%d, gasMLimit:%d, bGasMeter:%d\n", info.gasWanted, info.startingGas, info.ctx.GasMeter().GasConsumed(), info.ctx.GasMeter().Limit(), info.ctx.BlockGasMeter().GasConsumed())
	var anteCtx sdk.Context

	// Cache wrap context before AnteHandler call in case it aborts.
	// This is required for both CheckTx and DeliverTx.
	// Ref: https://github.com/cosmos/cosmos-sdk/issues/2772
	//
	// NOTE: Alternatively, we could require that AnteHandler ensures that
	// writes do not happen if aborted/failed.  This may have some
	// performance benefits, but it'll be more difficult to get right.
	anteCtx, info.msCacheAnte = dm.app.cacheTxContext(info.ctx, info.txBytes)
	anteCtx = anteCtx.WithEventManager(sdk.NewEventManager())
	//fmt.Printf("runAnte 1. gasWanted:%d, startingGas:%d, gasMeter:%d, gasMLimit:%d, bGasMeter:%d\n", info.gasWanted, info.startingGas, info.ctx.GasMeter().GasConsumed(), info.ctx.GasMeter().Limit(), info.ctx.BlockGasMeter().GasConsumed())
	newCtx, err := dm.app.anteHandler(anteCtx, info.tx, mode == runTxModeSimulate) // NewAnteHandler
	//fmt.Printf("runAnte 2. gasWanted:%d, startingGas:%d, gasMeter:%d, gasMLimit:%d, bGasMeter:%d\n", info.gasWanted, info.startingGas, info.ctx.GasMeter().GasConsumed(), info.ctx.GasMeter().Limit(), info.ctx.BlockGasMeter().GasConsumed())

	ms := info.ctx.MultiStore()
	info.accountNonce = newCtx.AccountNonce()

	if !newCtx.IsZero() {
		// At this point, newCtx.MultiStore() is cache-wrapped, or something else
		// replaced by the AnteHandler. We want the original multistore, not one
		// which was cache-wrapped for the AnteHandler.
		//
		// Also, in the case of the tx aborting, we need to track gas consumed via
		// the instantiated gas meter in the AnteHandler, so we update the context
		// prior to returning.
		info.ctx = newCtx.WithMultiStore(ms)
	}
	//fmt.Printf("runAnte 3. gasWanted:%d, startingGas:%d, gasMeter:%d, gasMLimit:%d, bGasMeter:%d\n", info.gasWanted, info.startingGas, info.ctx.GasMeter().GasConsumed(), info.ctx.GasMeter().Limit(), info.ctx.BlockGasMeter().GasConsumed())

	// GasMeter expected to be set in AnteHandler
	info.gasWanted = info.ctx.GasMeter().Limit()
	//fmt.Printf("runAnte 4. gasWanted:%d, startingGas:%d, gasMeter:%d, gasMLimit:%d, bGasMeter:%d\n", info.gasWanted, info.startingGas, info.ctx.GasMeter().GasConsumed(), info.ctx.GasMeter().Limit(), info.ctx.BlockGasMeter().GasConsumed())

	if err != nil {
		return err
	}

	return nil
}

func (dm *DeliverTxTasksManager) runTxSerialRoutine() {
	begin := time.Now()
	finished := 0
	for {
		if finished == dm.totalCount {
			break
		}

		if !dm.extractExecutingTask() {
			start := time.Now()
			<-dm.executeSignal
			elapsed := time.Since(start).Microseconds()
			dm.app.logger.Info("time to waiting for extractExecutingTask", "index", dm.curIndex, "us",elapsed)
			dm.anteDuration -= elapsed
			continue
		}

		//dm.app.logger.Info("runTxSerialRoutine", "index", dm.executingTask.index)
		start := time.Now()

		mode := runTxModeDeliverPartConcurrent
		info := dm.executingTask.info
		handler := info.handler

		handleGasFn := func() {
			gasStart := time.Now()

			// dm.app.pin(Refund, true, mode)
			handler.handleDeferRefund(info)
			// dm.app.pin(Refund, false, mode)

			handler.handleDeferGasConsumed(info)

			if r := recover(); r != nil {
				_ = dm.app.runTx_defer_recover(r, info)
				info.msCache = nil //TODO msCache not write
				info.result = nil
			}
			info.gInfo = sdk.GasInfo{GasWanted: info.gasWanted, GasUsed: info.ctx.GasMeter().GasConsumed()}

			gasE := time.Since(gasStart).Microseconds()
			dm.deferGasTime += gasE
		}

		execFinishedFn := func(txRs abci.ResponseDeliverTx) {
			dm.txResponses[dm.executingTask.index] = &txRs
			dm.resetExecutingTask()
			finished++

			elapsed := time.Since(start).Microseconds()
			dm.gasAndMsgsDuration += elapsed
		}

		// execute anteHandler failed
		//var execResult *executeResult
		if dm.executingTask.err != nil {//&& dm.executingTask.decodeFailed {
			//execResult = newExecuteResult(sdkerrors.ResponseDeliverTx(dm.executingTask.err, 0, 0, dm.app.trace), nil, uint32(dm.executingTask.index), uint32(0))
			//dm.txExeResults[dm.executingTask.index] = execResult

			txRs := sdkerrors.ResponseDeliverTx(dm.executingTask.err, 0, 0, dm.app.trace) //execResult.GetResponse()
			execFinishedFn(txRs)
			continue
		}

		gasStart := time.Now()
		err := info.handler.handleGasConsumed(info)
		gasE := time.Since(gasStart).Microseconds()
		dm.handleGasTime += gasE
		if err != nil {
			dm.app.logger.Error("handleGasConsumed failed", "err", err)
			//execResult = newExecuteResult(sdkerrors.ResponseDeliverTx(dm.executingTask.err, 0, 0, dm.app.trace), nil, uint32(dm.executingTask.index), uint32(0))
			//dm.txExeResults[dm.executingTask.index] = execResult

			txRs := sdkerrors.ResponseDeliverTx(dm.executingTask.err, 0, 0, dm.app.trace)
			execFinishedFn(txRs)
			continue
		}

		// todo: if ante failed during concurrently executing, try it again
		if dm.executingTask.anteFailed {
			// dm.app.pin(RunAnte, true, mode)

			if dm.app.anteHandler != nil {
				start := time.Now()
				err := dm.app.runAnte(info, mode)
				elasped := time.Since(start).Microseconds()
				dm.gasAndMsgsDuration -= elasped
				if err != nil {
					dm.app.logger.Error("runAnte failed", "err", err)
					//execResult = newExecuteResult(sdkerrors.ResponseDeliverTx(dm.executingTask.err, 0, 0, dm.app.trace), nil, uint32(dm.executingTask.index), uint32(0))
					//dm.txExeResults[dm.executingTask.index] = execResult

					txRs := sdkerrors.ResponseDeliverTx(dm.executingTask.err, 0, 0, dm.app.trace) //execResult.GetResponse()
					handleGasFn()
					execFinishedFn(txRs)
					// dm.app.pin(RunAnte, false, mode)
					continue
				}
			}
			// dm.app.pin(RunAnte, false, mode)
		}
		// todo: cache is the same for all deliverTx? Maybe it's no need to write cache there.
		wstart := time.Now()
		info.msCacheAnte.Write()
		info.ctx.Cache().Write(true)
		welasped := time.Since(wstart).Microseconds()
		dm.writeDuration += welasped

		// TODO: execute runMsgs etc.
		runMsgStart := time.Now()
		// dm.app.pin(RunMsg, true, mode)
		err = handler.handleRunMsg(info)
		// dm.app.pin(RunMsg, false, mode)
		runMsgE := time.Since(runMsgStart).Microseconds()
		dm.runMsgsTime += runMsgE

		handleGasFn()

		var resp abci.ResponseDeliverTx
		if err != nil {
			//dm.app.logger.Error("handleRunMsg failed", "err", err)
			resp = sdkerrors.ResponseDeliverTx(err, info.gInfo.GasWanted, info.gInfo.GasUsed, dm.app.trace)
		} else {
			resp = abci.ResponseDeliverTx{
				GasWanted: int64(info.gInfo.GasWanted), // TODO: Should type accept unsigned ints?
				GasUsed:   int64(info.gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
				Log:       info.result.Log,
				Data:      info.result.Data,
				Events:    info.result.Events.ToABCIEvents(),
			}
		}
		//execResult = newExecuteResult(resp, info.msCache, uint32(dm.executingTask.index), dm.executingTask.evmIndex)
		//execResult.err = err
		//dm.txExeResults[dm.executingTask.index] = execResult
		//txRs := execResult.GetResponse()
		execFinishedFn(resp)
	}

	// all txs are executed
	if finished == dm.totalCount {
		dm.done <- 0
		close(dm.executeSignal)
		close(dm.nextSignal)
		dm.serialDuration = time.Since(begin).Microseconds()
		totalSerialDuration += dm.serialDuration
	} else {
		dm.app.logger.Error("finished count is not equal to total count", "finished", finished, "total", dm.totalCount)
	}
}

func (dm *DeliverTxTasksManager) extractExecutingTask() bool {
	task, ok := dm.pendingTasks.Load(dm.curIndex+1)
	if ok {
		dm.executingTask = task.(*DeliverTxTask)
		dm.pendingTasks.Delete(dm.executingTask.index)

		dm.incrementWaitingCount(false)
	//} else {
	//	dm.app.logger.Error("extractExecutingTask failed", "index", dm.curIndex+1)
	}
	return ok
}

func (dm *DeliverTxTasksManager) resetExecutingTask() {
	dm.executingTask = nil
}

//func (dm *DeliverTxTasksManager) isWaiting(newTask bool) bool {
//	dm.mtx.Lock()
//	defer dm.mtx.Unlock()
//	if newTask {
//		return dm.waitingCount >= maxDeliverTxsConcurrentNum
//	} else {
//		return dm.waitingCount >= maxDeliverTxsConcurrentNum-1
//	}
//}

func (dm *DeliverTxTasksManager) incrementWaitingCount(increment bool) {
	if increment {
		dm.mtx.Lock()
		dm.waitingCount++
		count := dm.waitingCount
		dm.mtx.Unlock()

		if count >= maxDeliverTxsConcurrentNum{
			<-dm.nextSignal
			dm.executeSignalCount--
			if dm.executeSignalCount < 0 {
				dm.app.logger.Error("dm.executeSignalCount < 0", "count", dm.executeSignalCount)
			}
		}
	} else {
		dm.mtx.Lock()
		dm.curIndex++
		dm.waitingCount--
		count := dm.waitingCount
		dm.mtx.Unlock()

		if count >= maxDeliverTxsConcurrentNum-1 {
			dm.nextSignal <- 0
			dm.executeSignalCount++
			if dm.executeSignalCount > 1 {
				dm.app.logger.Error("dm.executeSignalCount > 1", "count", dm.executeSignalCount)
			}
		}
	}
}

//func (dm *DeliverTxTasksManager) getCurIndex() int {
//	dm.mtx.Lock()
//	defer dm.mtx.Unlock()
//	return dm.curIndex
//}

//-------------------------------------------------------------

func (app *BaseApp) DeliverTxsConcurrent(txs [][]byte) []*abci.ResponseDeliverTx {
	if app.deliverTxsMgr == nil {
		app.deliverTxsMgr = NewDeliverTxTasksManager(app)
	}

	//app.logger.Info("deliverTxs", "txs", len(txs))
	start := time.Now()
	app.deliverTxsMgr.deliverTxs(txs)

	if len(txs) > 0 {
		//waiting for call back
		<-app.deliverTxsMgr.done
		close(app.deliverTxsMgr.done)

		dur := time.Since(start).Microseconds()
		totalAnteDuration += app.deliverTxsMgr.anteDuration
		totalGasAndMsgsDuration += app.deliverTxsMgr.gasAndMsgsDuration
		totalWriteTime += app.deliverTxsMgr.writeDuration
		totalHandleGasTime += app.deliverTxsMgr.handleGasTime
		totalDeferGasTime += app.deliverTxsMgr.deferGasTime
		totalRunMsgsTime += app.deliverTxsMgr.runMsgsTime
		totalSavedTime = totalSavedTime + (app.deliverTxsMgr.anteDuration + app.deliverTxsMgr.gasAndMsgsDuration - app.deliverTxsMgr.serialDuration)
		app.logger.Info("all durations",
			"whole", dur,
			"ante", app.deliverTxsMgr.anteDuration,
			"serial", app.deliverTxsMgr.serialDuration,
			"gasAndMsgs", app.deliverTxsMgr.gasAndMsgsDuration,
			"handleGas", app.deliverTxsMgr.handleGasTime,
			"write", app.deliverTxsMgr.writeDuration,
			"runMsgs", app.deliverTxsMgr.runMsgsTime,
			"deferGas", app.deliverTxsMgr.deferGasTime,
			"serialSum", app.deliverTxsMgr.handleGasTime+app.deliverTxsMgr.writeDuration+app.deliverTxsMgr.runMsgsTime+app.deliverTxsMgr.deferGasTime,
			"handleGasAll", totalHandleGasTime,
			"writeAll", totalWriteTime,
			"runMsgsAll", totalRunMsgsTime,
			"deferGasAll", totalDeferGasTime,
			"serialSumAll", totalHandleGasTime+totalWriteTime+totalRunMsgsTime+totalDeferGasTime,
			"anteAll", totalAnteDuration,
			"gasAndMsgsAll", totalGasAndMsgsDuration,
			"serialAll", totalSerialDuration,
			"totalSavedTime", totalSavedTime,
			"saved", float64(app.deliverTxsMgr.anteDuration) / float64(dur))
	}

	return app.deliverTxsMgr.txResponses
}

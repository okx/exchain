package baseapp

import (
	"encoding/hex"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"sync"
	"time"
)

const (
	maxDeliverTxsConcurrentNum = 3
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
var totalWaitingTime = int64(0)
var totalRerunAnteTime = int64(0)
var totalBasicTime = int64(0)

type DeliverTxTask struct {
	//tx            sdk.Tx
	index         int
	feeForCollect int64
	anteFailed    bool
	info          *runTxInfo
	from          adb.Address
	fee           sdk.Coins
	isEvm         bool
	signCache     sdk.SigCache
	//evmIndex      uint32
	err error
}

func newDeliverTxTask(tx sdk.Tx, index int) *DeliverTxTask {
	t := &DeliverTxTask{
		//tx:    tx,
		index: index,
		info:  &runTxInfo{tx: tx},
	}

	return t
}

type DeliverTxTasksManager struct {
	done               chan int // done for all transactions are executed
	nextSignal         chan int // signal for taking a new tx into tasks
	executeSignal      chan int // signal for taking a new task from pendingTasks to executingTask
	waitingCount       int
	executeSignalCount int
	mtx                sync.Mutex

	totalCount    int
	curIndex      int
	tasks         sync.Map
	pendingTasks  sync.Map
	executingTask *DeliverTxTask
	currTxFee     sdk.Coins

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
	dm.currTxFee = sdk.Coins{}

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
		dm.makeNextTask(txs[taskIndex], taskIndex)
		taskIndex++
		dm.incrementWaitingCount(true)
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
	info := task.info
	info.handler = dm.app.getModeHandler(mode) //dm.handler

	// execute ante
	info.ctx = dm.app.getContextForTx(mode, info.txBytes) // same context for all txs in a block
	task.fee, task.isEvm, task.signCache = dm.app.getTxFee(info.ctx, info.tx)
	info.ctx = info.ctx.WithSigCache(task.signCache)
	info.ctx = info.ctx.WithCache(sdk.NewCache(dm.app.blockCache, useCache(mode))) // one cache for a tx

	if err := validateBasicTxMsgs(info.tx.GetMsgs()); err != nil {
		task.err = err
		dm.app.logger.Error("validateBasicTxMsgs failed", "err", err)
		return
	}

	if dm.app.anteAuthHandler != nil {
		err := dm.runAnte(task, mode) // dm.app.runAnte(task.info, mode)
		if err != nil {
			dm.app.logger.Error("ante failed 1", "err", err)
			task.anteFailed = true
		}
	}

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
		dm.app.logger.Error("tx decode failed", " err", err)
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

func (dm *DeliverTxTasksManager) runAnte(task *DeliverTxTask, mode runTxMode) error {
	info := task.info
	var anteCtx sdk.Context

	// Cache wrap context before AnteHandler call in case it aborts.
	// This is required for both CheckTx and DeliverTx.
	// Ref: https://github.com/cosmos/cosmos-sdk/issues/2772
	//
	// NOTE: Alternatively, we could require that AnteHandler ensures that
	// writes do not happen if aborted/failed.  This may have some
	// performance benefits, but it'll be more difficult to get right.
	anteCtx, info.msCacheAnte = dm.app.cacheTxContext(info.ctx, info.txBytes)	// info.msCacheAnte := ctx.MultiStore().CacheMultiStore(),  anteCtx := ctx.WithMultiStore(info.msCacheAnte)
	anteCtx = anteCtx.WithEventManager(sdk.NewEventManager())
	//if task.isEvm {
	//	info.msCacheAnte.IteratorCache(func(key, value []byte, isDirty bool) bool {
	//		if isDirty {
	//			fmt.Println(task.index, hex.EncodeToString(key), hex.EncodeToString(value))
	//		}
	//		return true
	//	})
	//}

	newCtx, err := dm.app.anteAuthHandler(anteCtx, info.tx, mode == runTxModeSimulate) // NewAnteHandler

	//if task.isEvm {
	//	info.msCacheAnte.IteratorCache(func(key, value []byte, isDirty bool) bool {
	//		if isDirty {
	//			fmt.Println("after", task.index, hex.EncodeToString(key), hex.EncodeToString(value))
	//		}
	//		return true
	//	})
	//}

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
		// todo: CacheMultiStore(info.msCacheAnte) is changed
		info.ctx = newCtx.WithMultiStore(ms)
		dm.updateEvmTxFrom(task)
	}
	// GasMeter expected to be set in AnteHandler
	info.gasWanted = info.ctx.GasMeter().Limit()
	if err != nil {
		return err
	}

	return nil
}

func (dm *DeliverTxTasksManager) updateEvmTxFrom(task *DeliverTxTask) {
	if task.isEvm && dm.app.evmTxFromHandler != nil {
		evmTx, ok := dm.app.evmTxFromHandler(task.info.ctx, task.info.tx)
		if ok {
			task.info.tx = evmTx
		}
	}
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
			dm.app.logger.Info("time to waiting for extractExecutingTask", "index", dm.curIndex, "us", elapsed)
			dm.anteDuration -= elapsed
			totalWaitingTime += elapsed
			continue
		}

		//dm.app.logger.Info("runTxSerialRoutine", "index", dm.executingTask.index)
		start := time.Now()

		info := dm.executingTask.info
		handler := info.handler

		handleGasFn := func() {
			gasStart := time.Now()

			handler.handleDeferRefund(info)

			handler.handleDeferGasConsumed(info)

			if r := recover(); r != nil {
				_ = dm.app.runTx_defer_recover(r, info)
				info.msCache = nil //TODO msCache not write
				info.result = nil
			}
			info.gInfo = sdk.GasInfo{GasWanted: info.gasWanted, GasUsed: info.ctx.GasMeter().GasConsumed()}

			dm.deferGasTime += time.Since(gasStart).Microseconds()
		}

		execFinishedFn := func(txRs abci.ResponseDeliverTx) {
			dm.txResponses[dm.executingTask.index] = &txRs
			dm.resetExecutingTask()
			finished++

			dm.gasAndMsgsDuration += time.Since(start).Microseconds()
		}

		// execute anteHandler failed
		if dm.executingTask.err != nil {
			txRs := sdkerrors.ResponseDeliverTx(dm.executingTask.err, 0, 0, dm.app.trace) //execResult.GetResponse()
			execFinishedFn(txRs)
			continue
		}

		gasStart := time.Now()
		err := info.handler.handleGasConsumed(info)
		dm.handleGasTime += time.Since(gasStart).Microseconds()
		if err != nil {
			dm.app.logger.Error("handleGasConsumed failed", "err", err)

			txRs := sdkerrors.ResponseDeliverTx(err, 0, 0, dm.app.trace)
			execFinishedFn(txRs)
			continue
		}

		// todo: if ante failed during concurrently executing, try it again
		if dm.executingTask.anteFailed && dm.app.anteAuthHandler != nil {
			start := time.Now()
			err := dm.runAnte(dm.executingTask, runTxModeDeliverPartConcurrent) // dm.app.runAnte(info, mode)
			elasped := time.Since(start).Microseconds()
			dm.gasAndMsgsDuration -= elasped
			totalRerunAnteTime += elasped
			if err != nil {
				dm.app.logger.Error("ante failed 2", "err", err)
				//if dm.app.deductFeeHandler != nil {
				//	err = dm.app.deductFeeHandler(info.ctx, info.tx)
				//}

				txRs := sdkerrors.ResponseDeliverTx(err, 0, 0, dm.app.trace) //execResult.GetResponse()
				handleGasFn()
				execFinishedFn(txRs)
				continue
			}
		}
		// update fee
		err = dm.consumeGas(dm.executingTask)
		if err != nil {
			txRs := sdkerrors.ResponseDeliverTx(err, 0, 0, dm.app.trace) //execResult.GetResponse()
			handleGasFn()
			execFinishedFn(txRs)
			continue
		}

		//wstart := time.Now()
		//if global.GetGlobalHeight() == 5810736 {
		//	printLog("msCacheAnte", info.msCacheAnte)
		//}
		//if dm.executingTask.isEvm {
		//	info.msCacheAnte.IteratorCache(func(key, value []byte, isDirty bool) bool {
		//		if isDirty {
		//			fmt.Println(hex.EncodeToString(key), hex.EncodeToString(value))
		//		}
		//		return true
		//	})
		//}
		info.msCacheAnte.Write()
		info.ctx.Cache().Write(true)
		//dm.writeDuration += time.Since(wstart).Microseconds()

		// execute runMsgs
		runMsgStart := time.Now()
		err = handler.handleRunMsg(info)
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
		execFinishedFn(resp)
	}

	// all txs are executed
	if finished == dm.totalCount {
		// update fee collector
		dm.updateFeeCollector()

		dm.done <- 0
		close(dm.executeSignal)
		close(dm.nextSignal)
		dm.serialDuration = time.Since(begin).Microseconds()
		totalSerialDuration += dm.serialDuration
	} else {
		dm.app.logger.Error("finished count is not equal to total count", "finished", finished, "total", dm.totalCount)
	}
}

func (dm *DeliverTxTasksManager) consumeGas(task *DeliverTxTask) error {
	info := task.info
	txType := info.tx.GetType()
	switch txType {
	case sdk.StdTxType:
		if dm.app.deductFeeHandler != nil {
			err := dm.app.deductFeeHandler(info.ctx, info.tx)
			if err != nil {
				dm.app.logger.Error("ante failed 3", "err", err)
				return err
			}
		}
	case sdk.EvmTxType:
		if dm.app.ethGasConsumeHandler != nil {
			if !task.isEvm {
				dm.app.logger.Error("is not an evm tx")
			}
			anteCtx := info.ctx.WithMultiStore(info.msCacheAnte)
			anteCtx = anteCtx.WithEventManager(sdk.NewEventManager())
			newCtx, err := dm.app.ethGasConsumeHandler(anteCtx, info.tx)
			ms := info.ctx.MultiStore()
			if !newCtx.IsZero() {
				info.ctx = newCtx.WithMultiStore(ms)
			}
			// GasMeter expected to be set in AnteHandler
			info.gasWanted = info.ctx.GasMeter().Limit()

			if err != nil {
				dm.app.logger.Error("ante failed 3", "err", err)
				return err
			}
		}
	default:

	}

	dm.calculateFeeForCollector(task.fee, true)

	return nil
}

func (dm *DeliverTxTasksManager) calculateFeeForCollector(fee sdk.Coins, add bool) {
	if add {
		dm.currTxFee = dm.currTxFee.Add(fee...)
	} else {
		dm.currTxFee = dm.currTxFee.Sub(fee)
	}
}

func (dm *DeliverTxTasksManager) updateFeeCollector() {
	dm.app.logger.Info("updateFeeCollector", "now", dm.currTxFee[0].Amount)
	ctx, cache := dm.app.cacheTxContext(dm.app.getContextForTx(runTxModeDeliver, []byte{}), []byte{})
	if err := dm.app.updateFeeCollectorAccHandler(ctx, dm.currTxFee); err != nil {
		panic(err)
	}
	cache.Write()
}

func (dm *DeliverTxTasksManager) extractExecutingTask() bool {
	task, ok := dm.pendingTasks.Load(dm.curIndex + 1)
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

func printLog(msg string, cache sdk.CacheMultiStore) {
	height := global.GetGlobalHeight()
	if height == 5810736 {
		//fmt.Println(msg)
		cache.IteratorCache(func(key, value []byte, isDirty bool) bool {
			if isDirty {
				fmt.Println(msg, hex.EncodeToString(key), hex.EncodeToString(value))
			}
			return true
		})
		fmt.Println("---------------------------------------------------------------------")
	}
}

func (dm *DeliverTxTasksManager) incrementWaitingCount(increment bool) {
	if increment {
		dm.mtx.Lock()
		dm.waitingCount++
		count := dm.waitingCount
		dm.mtx.Unlock()

		if count >= maxDeliverTxsConcurrentNum {
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
			"waitingAll", totalWaitingTime,
			"rerunAnteAll", totalRerunAnteTime,
			"totalSavedTime", totalSavedTime,
			"saved", float64(app.deliverTxsMgr.anteDuration)/float64(dur))
	}

	return app.deliverTxsMgr.txResponses
}

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
	maxDeliverTxsConcurrentNum = 10
)

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
	//pendingSignal chan int // signal for taking a new task from tasks into pendingTasks
	executeSignal chan int // signal for taking a new task from pendingTasks to executingTask
	isWaiting     bool

	totalCount    int
	curIndex      int
	tasks         sync.Map
	pendingTasks  sync.Map
	executingTask *DeliverTxTask

	txResponses []*abci.ResponseDeliverTx

	app *BaseApp
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

	dm.totalCount = len(txs)
	dm.curIndex = -1

	dm.tasks = sync.Map{}
	dm.pendingTasks = sync.Map{}
	dm.txResponses = make([]*abci.ResponseDeliverTx, len(txs))

	go dm.makeTasksRoutine(txs)
	go dm.runTxSerialRoutine()
}

func (dm *DeliverTxTasksManager) makeTasksRoutine(txs [][]byte) {
	taskIndex := 0
	for {
		if taskIndex == dm.totalCount {
			break
		}

		remaining := taskIndex - (dm.curIndex + 1) //- numTasks - numPending
		switch {
		case remaining >= maxDeliverTxsConcurrentNum:
			dm.isWaiting = true
			<-dm.nextSignal

		default:
			dm.isWaiting = false
			dm.makeNextTask(txs[taskIndex], taskIndex)
			taskIndex++
		}
	}
}

func (dm *DeliverTxTasksManager) makeNextTask(tx []byte, index int) {
	go dm.runTxPartConcurrent(tx, index)
}

func (dm *DeliverTxTasksManager) runTxPartConcurrent(txByte []byte, index int) {
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

	dm.app.pin(ValTxMsgs, true, mode)
	if err := validateBasicTxMsgs(task.tx.GetMsgs()); err != nil {
		task.err = err
		dm.app.logger.Error("validateBasicTxMsgs failed", "err", err)
		return
	}
	dm.app.pin(ValTxMsgs, false, mode)

	dm.app.pin(AnteHandler, true, mode)
	if dm.app.anteHandler != nil {
		err := dm.runAnte(task.info, mode)
		if err != nil {
			dm.app.logger.Error("runAnte failed", "err", err)
			task.anteFailed = true
		}
	}
	dm.app.pin(AnteHandler, false, mode)
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

	dm.app.logger.Info("new into pendingTasks", "index", task.index)
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
	finished := 0
	for {
		if finished == dm.totalCount {
			dm.app.logger.Info("break runTxSerialRoutine")
			break
		}

		if !dm.extractExecutingTask() {
			start := time.Now()
			<-dm.executeSignal
			elapsed := time.Since(start).Milliseconds()
			dm.app.logger.Error("time to waiting for extractExecutingTask", "index", dm.curIndex, "ms",elapsed)
			continue
		}
		if dm.isWaiting {
			dm.nextSignal <- 0
		}

		dm.app.logger.Info("runTxSerialRoutine", "index", dm.executingTask.index)

		mode := runTxModeDeliverPartConcurrent
		info := dm.executingTask.info
		handler := info.handler

		handleGasFn := func() {
			dm.app.pin(Refund, true, mode)
			handler.handleDeferRefund(info)
			dm.app.pin(Refund, false, mode)

			handler.handleDeferGasConsumed(info)

			if r := recover(); r != nil {
				_ = dm.app.runTx_defer_recover(r, info)
				info.msCache = nil //TODO msCache not write
				info.result = nil
			}
			info.gInfo = sdk.GasInfo{GasWanted: info.gasWanted, GasUsed: info.ctx.GasMeter().GasConsumed()}
		}

		execFinishedFn := func(txRs abci.ResponseDeliverTx) {
			dm.txResponses[dm.executingTask.index] = &txRs
			dm.resetExecutingTask()
			finished++
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

		dm.app.logger.Info("handleGasConsumed start")
		err := info.handler.handleGasConsumed(info)
		dm.app.logger.Info("handleGasConsumed finished")
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
			dm.app.pin(AnteHandler, true, mode)

			if dm.app.anteHandler != nil {
				// dm.app.logger.Info("rerun Ante", "index", dm.executingTask.index)
				err := dm.app.runAnte(info, mode)
				if err != nil {
					dm.app.logger.Error("runAnte failed", "err", err)
					//execResult = newExecuteResult(sdkerrors.ResponseDeliverTx(dm.executingTask.err, 0, 0, dm.app.trace), nil, uint32(dm.executingTask.index), uint32(0))
					//dm.txExeResults[dm.executingTask.index] = execResult

					txRs := sdkerrors.ResponseDeliverTx(dm.executingTask.err, 0, 0, dm.app.trace) //execResult.GetResponse()
					handleGasFn()
					execFinishedFn(txRs)
					dm.app.pin(AnteHandler, false, mode)
					continue
				}
			}
			dm.app.pin(AnteHandler, false, mode)
		}
		info.msCacheAnte.Write()
		info.ctx.Cache().Write(true)
		dm.app.logger.Info("runAnte succeed")

		// TODO: execute runMsgs etc.
		dm.app.pin(RunMsgs, true, mode)
		err = handler.handleRunMsg(info)
		dm.app.pin(RunMsgs, false, mode)
		dm.app.logger.Info("runMsg succeed")

		handleGasFn()
		dm.app.logger.Info("handleGasFn succeed")

		var resp abci.ResponseDeliverTx
		if err != nil {
			dm.app.logger.Error("handleRunMsg failed", "err", err)
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
		dm.app.logger.Info("execFinishedFn succeed")
	}

	// all txs are executed
	if finished == dm.totalCount {
		dm.app.logger.Info("runTxSerialRoutine finished")
		dm.done <- 0
	}
}

func (dm *DeliverTxTasksManager) extractExecutingTask() bool {
	task, ok := dm.pendingTasks.Load(dm.curIndex+1)
	if ok {
		dm.executingTask = task.(*DeliverTxTask)
		dm.pendingTasks.Delete(dm.executingTask.index)
		dm.curIndex++
		return true
	} else {
		dm.app.logger.Error("extractExecutingTask failed", "index", dm.curIndex+1)
	}
	return false
}

func (dm *DeliverTxTasksManager) resetExecutingTask() {
	dm.executingTask = nil
}

//-------------------------------------------------------------

func (app *BaseApp) DeliverTxsConcurrent(txs [][]byte) []*abci.ResponseDeliverTx {
	if app.deliverTxsMgr == nil {
		app.deliverTxsMgr = NewDeliverTxTasksManager(app)
	}

	app.logger.Info("deliverTxs", "txs", len(txs))
	app.deliverTxsMgr.deliverTxs(txs)

	if len(txs) > 0 {
		//waiting for call back
		<-app.deliverTxsMgr.done
	}

	return app.deliverTxsMgr.txResponses
}

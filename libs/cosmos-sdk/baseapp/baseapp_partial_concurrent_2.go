package baseapp

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/libs/log"

	//"github.com/okex/exchain/libs/tendermint/libs/service"
	"sync"
	"time"
)

//-------------------------------------
type BasicProcessFn func(txByte []byte, index int) *DeliverTxTask
type RunAnteFn func(task *DeliverTxTask) error

type dttRoutine struct {
	//service.BaseService
	done    chan int
	task    *DeliverTxTask
	txByte  chan []byte
	txIndex int
	rerunCh chan int
	index   int8

	logger log.Logger

	basicProFn BasicProcessFn
	runAnteFn  RunAnteFn
}

func newDttRoutine(index int8, basicProcess BasicProcessFn, runAnte RunAnteFn) *dttRoutine {
	dttr := &dttRoutine{
		//done:       make(chan int),
		//txByte:     make(chan []byte),
		//rerunCh:    make(chan int),
		index:      index,
		basicProFn: basicProcess,
		runAnteFn:  runAnte,
	}

	return dttr
}

func (dttr *dttRoutine) setLogger(logger log.Logger) {
	dttr.logger = logger
}

func (dttr *dttRoutine) makeNewTask(txByte []byte, index int) {
	dttr.txIndex = index
	dttr.task = nil
	dttr.logger.Info("makeNewTask 1", "index", dttr.txIndex)
	dttr.txByte <- txByte
}

func (dttr *dttRoutine) OnStart() error {
	dttr.done = make(chan int)
	dttr.txByte = make(chan []byte)
	dttr.rerunCh = make(chan int)
	go dttr.executeTaskRoutine()
	return nil
}

func (dttr *dttRoutine) OnReset() error {
	return nil
}

func (dttr *dttRoutine) stop() {
	dttr.done <- 0
}

func (dttr *dttRoutine) executeTaskRoutine() {
	for {
		select {
		case <-dttr.done:
			close(dttr.txByte)
			close(dttr.rerunCh)
			return
		case tx := <-dttr.txByte:
			//dttr.logger.Info("basicProFn", "index", dttr.txIndex)
			dttr.task = dttr.basicProFn(tx, dttr.txIndex)
			dttr.task.routineIndex = dttr.index
			if dttr.task.err == nil && !dttr.task.needToRerun {
				dttr.runAnteFn(dttr.task)
			} else {
				dttr.logger.Error("DonotRunAnte", "index", dttr.task.index, "needToRerun", dttr.task.needToRerun, "err", dttr.task.err)
			}
		case <-dttr.rerunCh:
			if dttr.task.step == partialConcurrentStepBasicSucceed || dttr.task.step == partialConcurrentStepAnteFailed || dttr.task.step == partialConcurrentStepAnteSucceed {
				dttr.runAnteFn(dttr.task)
			} else {
				dttr.logger.Error("shouldRerunLater", "index", dttr.task.index, "step", dttr.task.step)
				// maybe the task is in other condition, running concurrent execution or running make new task.
				dttr.task.canRerun++
				//dttr.task.needToRerun = false
			}
		}
	}
}

func (dttr *dttRoutine) checkConflict(addr string, index int) bool {
	if dttr.task == nil || dttr.task.index == index || dttr.task.from != addr {
		return false
	}
	if dttr.task.index < index {
		return true
	} else {
		if dttr.task.step != partialConcurrentStepBasicFailed {
			dttr.logger.Error("needToRerun 1", "index", dttr.task.index, "conflicted", index)
			dttr.task.needToRerun = true
		}
		return false
	}
}

func (dttr *dttRoutine) hasExistPrevTask(addr string, index int) bool {
	if dttr.task == nil {
		return false
	}
	// do not deal with the situation that getTxFeeAndFromHandler has not finished
	if dttr.task.index < index && dttr.task.from == addr {
		return true
	}
	return false
}

func (dttr *dttRoutine) couldRerun(index int) {
	if dttr.task == nil || dttr.task.canRerun > 0 || dttr.task.index == index {
		return
	}
	dttr.logger.Error("couldRerun", "index", dttr.task.index, "finished", index)
	dttr.rerunCh <- 0
}

//-------------------------------------

type DTTManager struct {
	done            chan int
	totalCount      int
	txs             [][]byte
	concurrentIndex int
	dttRoutineList  []*dttRoutine //sync.Map	// key: txIndex, value: dttRoutine
	serialIndex     int
	serialTask      *DeliverTxTask
	serialCh        chan *DeliverTxTask

	mtx       sync.Mutex
	currTxFee sdk.Coins

	txResponses []*abci.ResponseDeliverTx
	app         *BaseApp
}

func NewDTTManager(app *BaseApp) *DTTManager {
	dttm := &DTTManager{
		app: app,
	}
	dttm.dttRoutineList = make([]*dttRoutine, 0, maxDeliverTxsConcurrentNum) //sync.Map{}
	for i := 0; i < maxDeliverTxsConcurrentNum; i++ {
		dttr := newDttRoutine(int8(i), dttm.makeNewTask, dttm.runConcurrentAnte)
		dttr.setLogger(dttm.app.logger)
		dttm.dttRoutineList = append(dttm.dttRoutineList, dttr)
		dttm.app.logger.Info("newDttRoutine", "index", i, "list", len(dttm.dttRoutineList))

		//err := dttr.OnStart()
		//if err != nil {
		//	dttm.app.logger.Error("Error starting DttRoutine", "err", err)
		//}
	}

	return dttm
}

func (dttm *DTTManager) deliverTxs(txs [][]byte) {
	dttm.done = make(chan int, 1)

	dttm.totalCount = len(txs)
	dttm.txs = txs
	dttm.currTxFee = sdk.Coins{}
	dttm.serialTask = nil
	dttm.serialIndex = -1
	dttm.serialCh = make(chan *DeliverTxTask)

	dttm.txResponses = make([]*abci.ResponseDeliverTx, len(txs))
	go dttm.serialRoutine()

	//dttm.dttRoutineList = make([]*dttRoutine, 0, maxDeliverTxsConcurrentNum) //sync.Map{}
	for i := 0; i < maxDeliverTxsConcurrentNum; i++ {
		//dttr := newDttRoutine(int8(i), dttm.makeNewTask, dttm.runConcurrentAnte)
		//dttr.setLogger(dttm.app.logger)
		////dttm.dttRoutineList[i] = dttr
		//dttm.dttRoutineList = append(dttm.dttRoutineList, dttr)
		//dttm.app.logger.Info("newDttRoutine", "index", i, "list", len(dttm.dttRoutineList))
		dttr := dttm.dttRoutineList[i]

		err := dttr.OnStart()
		if err != nil {
			dttm.app.logger.Error("Error starting DttRoutine", "err", err)
		}
		//dttm.setConcurrentIndex(i)
		//if dttm.concurrentIndex < i {
			dttm.concurrentIndex = i
		//}
		dttr.makeNewTask(txs[i], i)
	}
}

func (dttm *DTTManager) makeNewTask(txByte []byte, index int) *DeliverTxTask {
	dttm.app.logger.Info("makeNewTask", "index", index)

	// create a new task
	var realTx sdk.Tx
	var err error
	if mem := GetGlobalMempool(); mem != nil {
		realTx, _ = mem.ReapEssentialTx(txByte).(sdk.Tx)
	}
	if realTx == nil {
		realTx, err = dttm.app.txDecoder(txByte)
	}
	task := newDeliverTxTask(realTx, index)
	task.info.txBytes = txByte
	if err != nil {
		task.err = err
		//dm.app.logger.Error("tx decode failed", "err", err)
	}

	if task.err != nil {
		task.step = partialConcurrentStepBasicFailed
		//task.setStep(partialConcurrentStepBasicFailed)
		return task
	}

	task.info.handler = dttm.app.getModeHandler(runTxModeDeliverPartConcurrent)                 //dm.handler
	task.info.ctx = dttm.app.getContextForTx(runTxModeDeliverPartConcurrent, task.info.txBytes) // same context for all txs in a block
	task.setUpdateCount(0)
	task.fee, task.isEvm, task.from = dttm.app.getTxFeeAndFromHandler(task.info.ctx, task.info.tx)

	if err = validateBasicTxMsgs(task.info.tx.GetMsgs()); err != nil {
		task.err = err
		dttm.app.logger.Error("validateBasicTxMsgs failed", "err", err)
		task.step = partialConcurrentStepBasicFailed
	} else {
		//dttm.app.logger.Info("hasExistPrevTask", "index", task.index, "from", task.from)
		task.step = partialConcurrentStepBasicSucceed
		// need to check whether exist running tx who has the same from but smaller txIndex
		count := len(dttm.dttRoutineList)
		for i := 0; i < count; i++ {
			dttr := dttm.dttRoutineList[i]
			if dttr == nil {
				continue
			}
			//dttm.app.logger.Info("hasExistPrevTask 1", "routine", dttr.index, "task", dttr.txIndex)
			task.needToRerun = dttr.hasExistPrevTask(task.from, task.index)
			if task.needToRerun {
				dttm.app.logger.Error("needToRerun 3", "index", task.index, "conflicted", dttr.task.index)
				break
			}
		}
	}

	return task
}

func (dttm *DTTManager) runConcurrentAnte(task *DeliverTxTask) error {
	if dttm.app.anteHandler == nil {
		return fmt.Errorf("anteHandler cannot be nil")
	}

	if task.needToRerun {
		dttm.app.logger.Error("ResetContext", "index", task.index)

		task.step = partialConcurrentStepBasicSucceed
		task.info.ctx = dttm.app.getContextForTx(runTxModeDeliverPartConcurrent, task.info.txBytes) // same context for all txs in a block
		task.setUpdateCount(0)
		task.needToRerun = false
		task.canRerun = 0
	}

	task.step = partialConcurrentStepAnteStart
	dttm.app.logger.Info("RunAnte", "index", task.index, "routine", task.routineIndex, "addr", task.from)

	task.info.ctx = task.info.ctx.WithCache(sdk.NewCache(dttm.app.blockCache, useCache(runTxModeDeliverPartConcurrent))) // one cache for a tx

	dttm.accountUpdated(false, 2, task.from, -1)
	err := dttm.runAnte(task)
	if err != nil {
		dttm.app.logger.Error("ante failed 1", "index", task.index, "err", err)
		//task.anteFailed = true
		task.step = partialConcurrentStepAnteFailed
	} else {
		dttm.app.logger.Info("AnteSucceed", "index", task.index)
		task.step = partialConcurrentStepAnteSucceed
	}
	task.err = err
	count := len(dttm.dttRoutineList)
	for i := 0; i < count; i++ {
		dttr := dttm.dttRoutineList[i]
		if dttr == nil {
			continue
		}
		conflicted := dttr.checkConflict(task.from, task.index)
		if conflicted {
			dttr.logger.Error("needToRerun 2", "index", task.index, "conflicted", dttr.task.index)
			task.needToRerun = true
		}
	}
	if task.canRerun > 0 {
		dttr := dttm.dttRoutineList[task.routineIndex]
		go func() {
			dttr.logger.Error("rerunCh 2", "index", task.index)
			dttr.rerunCh <- 0
		}()
	} else if dttm.serialIndex+1 == task.index {
		if dttm.serialTask == nil {
			dttm.app.logger.Info("AnteFinished 1", "index", task.index)
			dttm.serialCh <- task
		} else {
			dttm.app.logger.Info("AnteFinished 2", "index", task.index)
		}
	} else {
		dttm.app.logger.Info("AnteFinished 3", "index", task.index)
	}

	return nil
}

func (dttm *DTTManager) runAnte(task *DeliverTxTask) error {
	info := task.info
	var anteCtx sdk.Context
	anteCtx, info.msCacheAnte = dttm.app.cacheTxContext(info.ctx, info.txBytes) // info.msCacheAnte := ctx.MultiStore().CacheMultiStore(),  anteCtx := ctx.WithMultiStore(info.msCacheAnte)
	anteCtx = anteCtx.WithEventManager(sdk.NewEventManager())
	//anteCtx = anteCtx.WithAnteTracer(dm.app.anteTracer)

	newCtx, err := dttm.app.anteHandler(anteCtx, info.tx, false) // NewAnteHandler

	ms := info.ctx.MultiStore()
	//info.accountNonce = newCtx.AccountNonce()

	if !newCtx.IsZero() {
		info.ctx = newCtx.WithMultiStore(ms)
	}
	// GasMeter expected to be set in AnteHandler
	info.gasWanted = info.ctx.GasMeter().Limit()
	if err != nil {
		return err
	}

	return nil
}

func (dttm *DTTManager) serialRoutine() {
	for {
		select {
		case task := <-dttm.serialCh:
			// runMsgs etc.
			if task.index == dttm.serialIndex+1 {
				dttm.serialIndex = task.index
				dttm.serialTask = task
				dttm.serialExecution()
				dttm.serialTask = nil
				dttm.app.logger.Info("NextSerialTask", "index", dttm.serialIndex+1)

				if dttm.serialIndex == dttm.totalCount-1 {
					count := len(dttm.dttRoutineList)
					for i := 0; i < count; i++ {
						dttr := dttm.dttRoutineList[i]
						dttr.stop()
					}
					dttm.app.logger.Info("TotalTxFeeForCollector", "fee", dttm.currTxFee)

					dttm.done <- 0
					close(dttm.serialCh)
					return
				}

				// make new task for this routine
				dttr := dttm.dttRoutineList[task.routineIndex]
				//currIndex := dttm.concurrentIndex
				//if currIndex < maxDeliverTxsConcurrentNum {
				//	currIndex = maxDeliverTxsConcurrentNum
				//}
				if dttr != nil && dttm.concurrentIndex < dttm.totalCount-1 {
					dttm.concurrentIndex++
					dttr.makeNewTask(dttm.txs[dttm.concurrentIndex], dttm.concurrentIndex)
				}

				// todo: check whether there are ante-finished task
				count := len(dttm.dttRoutineList)
				var nextTask *DeliverTxTask
				var rerunRoutine *dttRoutine
				for i := 0; i < count; i++ {
					dttr = dttm.dttRoutineList[i]
					if dttr.task == nil {
						continue
					}

					// if exists the next task which has finished the concurrent execution
					if dttr.task.index == dttm.serialIndex+1 {
						dttm.app.logger.Info("WaitNextSerialTask", "index", dttr.task.index, "needToRerun", dttr.task.needToRerun, "step", dttr.task.step)
						if dttr.task.from == task.from {
							//go func() {
							dttr.logger.Error("rerunCh", "index", dttr.task.index)
							dttr.rerunCh <- 0
							//}()
						} else if dttr.task.step == partialConcurrentStepBasicFailed ||
							dttr.task.step == partialConcurrentStepAnteFailed ||
							dttr.task.step == partialConcurrentStepAnteSucceed {
							nextTask = dttr.task
						}
					} else if dttr.task.from == task.from {
						if rerunRoutine == nil {
							rerunRoutine = dttr
						} else if dttr.task.index < rerunRoutine.task.index {
							rerunRoutine = dttr
						}
					}
				}

				if rerunRoutine != nil {
					rerunRoutine.couldRerun(task.index)
				}
				if nextTask != nil {
					go func() {
						dttm.serialCh <- nextTask
					}()
				}
			} else {
				panic(fmt.Sprintf("invalid index for serial execution: expected %x, got %x\n", dttm.serialIndex, task.index))
			}
		}
	}
}

func (dttm *DTTManager) serialExecution() {
	dttm.app.logger.Info("RunStatefulSerialRoutine", "index", dttm.serialTask.index)

	info := dttm.serialTask.info
	handler := info.handler

	handleGasFn := func() {
		gasStart := time.Now()

		dttm.updateFeeCollector()

		//dttm.app.logger.Info("handleDeferRefund", "index", dttm.serialTask.txIndex, "addr", dttm.serialTask.from)
		dttm.accountUpdated(false, 1, dttm.serialTask.from, -1)
		handler.handleDeferRefund(info)

		handler.handleDeferGasConsumed(info)

		if r := recover(); r != nil {
			_ = dttm.app.runTx_defer_recover(r, info)
			info.msCache = nil //TODO msCache not write
			info.result = nil
		}
		info.gInfo = sdk.GasInfo{GasWanted: info.gasWanted, GasUsed: info.ctx.GasMeter().GasConsumed()}

		totalDeferGasTime += time.Since(gasStart).Microseconds()
	}

	execFinishedFn := func(txRs abci.ResponseDeliverTx) {
		dttm.app.logger.Info("SerialFinished", "index", dttm.serialTask.index, "routine", dttm.serialTask.routineIndex)
		dttm.txResponses[dttm.serialTask.index] = &txRs
	}

	// execute anteHandler failed
	if dttm.serialTask.err != nil {
		dttm.app.logger.Error("RunSerialFinished", "index", dttm.serialTask.index, "err", dttm.serialTask.err)
		txRs := sdkerrors.ResponseDeliverTx(dttm.serialTask.err, 0, 0, dttm.app.trace) //execResult.GetResponse()
		execFinishedFn(txRs)
		return
	}

	//dttm.app.logger.Info("WriteAnteCache", "index", dttm.serialTask.txIndex)
	info.msCacheAnte.Write()
	info.ctx.Cache().Write(true)
	dttm.calculateFeeForCollector(dttm.serialTask.fee, true)

	gasStart := time.Now()
	err := info.handler.handleGasConsumed(info)
	//dttm.handleGasTime += time.Since(gasStart).Microseconds()
	totalHandleGasTime += time.Since(gasStart).Microseconds()
	if err != nil {
		dttm.app.logger.Error("handleGasConsumed failed", "err", err)

		txRs := sdkerrors.ResponseDeliverTx(err, 0, 0, dttm.app.trace)
		execFinishedFn(txRs)
		return
	}

	// execute runMsgs
	//dttm.app.logger.Info("handleRunMsg", "index", dttm.serialTask.txIndex, "addr", dttm.serialTask.from)
	dttm.accountUpdated(false, 2, dttm.serialTask.from, -1)
	runMsgStart := time.Now()
	err = handler.handleRunMsg(info)
	totalRunMsgsTime += time.Since(runMsgStart).Microseconds()

	handleGasFn()

	var resp abci.ResponseDeliverTx
	if err != nil {
		//dttm.app.logger.Error("handleRunMsg failed", "err", err)
		resp = sdkerrors.ResponseDeliverTx(err, info.gInfo.GasWanted, info.gInfo.GasUsed, dttm.app.trace)
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

//func (dttm *DTTManager) setConcurrentIndex(index int) {
//	dttm.mtx.Lock()
//	defer dttm.mtx.Unlock()
//
//	dttm.concurrentIndex = index
//}
//
//func (dttm *DTTManager) getConcurrentIndex() int {
//	dttm.mtx.Lock()
//	defer dttm.mtx.Unlock()
//
//	return dttm.concurrentIndex
//}

func (dttm *DTTManager) calculateFeeForCollector(fee sdk.Coins, add bool) {
	dttm.mtx.Lock()
	defer dttm.mtx.Unlock()

	if add {
		dttm.currTxFee = dttm.currTxFee.Add(fee...)
	} else {
		dttm.currTxFee = dttm.currTxFee.Sub(fee)
	}
	//dm.app.logger.Info("CalculateFeeForCollector", "fee", dm.currTxFee)
}

func (dttm *DTTManager) updateFeeCollector() {
	//	dm.app.logger.Info("updateFeeCollector", "now", dm.currTxFee)
	ctx, cache := dttm.app.cacheTxContext(dttm.app.getContextForTx(runTxModeDeliver, []byte{}), []byte{})
	if err := dttm.app.updateFeeCollectorAccHandler(ctx, dttm.currTxFee); err != nil {
		panic(err)
	}
	cache.Write()
}

func (dttm *DTTManager) OnAccountUpdated(acc exported.Account) {
	addr := acc.GetAddress().String()
	if global.GetGlobalHeight() == 5811070 && hex.EncodeToString(acc.GetAddress()) == "34bfa7d438d3b1cb23c3f4557ba5ac6160be4e4c" {
		dttm.app.logger.Error("OnAccountUpdated", "addr", addr)
	}
	//dm.app.logger.Info("OnAccountUpdated", "coins", acc.GetCoins(), "addr", addr)
	waitingIndex := -1
	if dttm.serialTask == nil {
		waitingIndex = dttm.serialIndex + 1
	}
	dttm.accountUpdated(true, 1, addr, waitingIndex)
}

func (dttm *DTTManager) accountUpdated(happened bool, times int8, address string, waitingIndex int) {
	//dttm.mtx.Lock()
	//defer dttm.mtx.Unlock()

	num := len(dttm.dttRoutineList)
	for i := 0; i < num; i++ {
		dttr := dttm.dttRoutineList[i]
		if dttr.task == nil || dttr.task.from != address {
			continue
		}

		task := dttr.task
		count := task.getUpdateCount()
		if happened {
			task.setUpdateCount(count + times)
			// todo: whether should rerun the task
			if task.index != waitingIndex && task.updateCount > 0 && task.needToRerunWhenContextChanged() {
				//task.needToRerun = true
				dttm.app.logger.Error("accountUpdatedToRerun", "index", task.index, "step", task.step)
				dttr.rerunCh <- 0
			}
		} else {
			task.setUpdateCount(count - times)
		}
	}
}

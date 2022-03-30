package baseapp

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"runtime"

	//"github.com/okex/exchain/libs/tendermint/libs/service"
	"sync"
	"time"
)

var totalSerialWaitingCount = 0

//-------------------------------------
type BasicProcessFn func(txByte []byte, index int) *DeliverTxTask
type RunAnteFn func(task *DeliverTxTask) error

//type txBytesWithTxIndex struct {
//	txBytes []byte
//	index int
//}

type dttRoutine struct {
	//service.BaseService
	done    chan int
	task    *DeliverTxTask
	txByte  chan []byte
	txIndex int
	rerunCh chan int
	//runAnteCh chan int
	index   int8
	mtx sync.Mutex

	logger log.Logger

	basicProFn BasicProcessFn
	runAnteFn  RunAnteFn
}

func newDttRoutine(index int8, basicProcess BasicProcessFn, runAnte RunAnteFn) *dttRoutine {
	dttr := &dttRoutine{
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
	//dttr.task = nil
	//dttr.logger.Info("makeNewTask", "index", dttr.txIndex)
	dttr.txByte <- txByte
}

func (dttr *dttRoutine) OnStart() error {
	dttr.done = make(chan int)
	dttr.txByte = make(chan []byte, 2)
	dttr.rerunCh = make(chan int, 5)
	//dttr.runAnteCh = make(chan int, 5)
	go dttr.executeTaskRoutine()
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
			if dttr.task.err == nil {//&& (!dttr.task.needToRerun || dttr.task.canRerun > 0) {
				dttr.runAnteFn(dttr.task)
			//} else {
			//	dttr.logger.Error("DonotRunAnte", "index", dttr.task.index)
			}
		case <-dttr.rerunCh:
			step := dttr.task.getStep()
			if step == partialConcurrentStepFinished ||
				step == partialConcurrentStepSerialExecute ||
				step == partialConcurrentStepBasicFailed {
				dttr.logger.Error("task is empty or finished")
			} else if step == partialConcurrentStepBasicSucceed ||
				step == partialConcurrentStepAnteFailed ||
				step == partialConcurrentStepAnteSucceed {
				//dttr.task.needToRerun = true
				dttr.logger.Error("RerunTask", "index", dttr.task.index, "step", step)
				dttr.runAnteFn(dttr.task)
			} else {
				dttr.logger.Error("shouldRerunLater", "index", dttr.task.index, "step", step)
				// maybe the task is in other condition, running concurrent execution or running make new task.
				dttr.task.canRerun++
				//dttr.task.needToRerun = false
			}

		}
	}
}

func (dttr *dttRoutine) shouldRerun() {
	if len(dttr.rerunCh) > 0 {
		return
	}
	if dttr.task.step == partialConcurrentStepAnteStart || dttr.task.step == partialConcurrentStepAnteFailed || dttr.task.step == partialConcurrentStepAnteSucceed  {
		dttr.logger.Error("shouldRerun", "index", dttr.task.index, "step", dttr.task.step)
		dttr.rerunCh <- 0
	}
}

//func (dttr *dttRoutine) checkConflict(addr string, index int) bool {
//	if dttr.task == nil || dttr.task.index == index || dttr.task.from != addr {
//		return false
//	}
//	if dttr.task.index < index {
//		return true
//	} else {
//		step := dttr.task.getStep()
//		if step != partialConcurrentStepBasicFailed &&
//			step != partialConcurrentStepFinished {
//			dttr.logger.Error("needToRerun 1", "index", dttr.task.index, "conflicted", index)
//			//dttr.task.needToRerun = true
//		}
//		//if dttr.task.step == partialConcurrentStepAnteSucceed || dttr.task.step == partialConcurrentStepAnteFailed {
//		//	return true
//		//}
//		return false
//	}
//}
//
//func (dttr *dttRoutine) hasExistPrevTask(addr string, index int) bool {
//	//step := dttr.task.getStep()
//	if dttr.task == nil ||
//		dttr.task.getStep() == partialConcurrentStepFinished ||
//		dttr.task.getStep() == partialConcurrentStepBasicFailed {
//		return false
//	}
//	// do not deal with the situation that getTxFeeAndFromHandler has not finished
//	if dttr.task.index < index && dttr.task.from == addr {
//		return true
//	}
//	return false
//}
//
//func (dttr *dttRoutine) couldRerun(index int) {
//	if dttr.task == nil || dttr.task.canRerun > 0 || dttr.task.index == index {
//		return
//	}
//	//go func() {
//	dttr.logger.Error("couldRerun", "index", dttr.task.index, "finished", index)
//	//dttr.task.needToRerun = true
//	dttr.rerunCh <- 0
//	//}()
//}

//-------------------------------------

type DTTManager struct {
	done           chan int
	totalCount     int
	txs            [][]byte
	//tasks []*DeliverTxTask
	startFinished  bool
	dttRoutineList []*dttRoutine //sync.Map	// key: txIndex, value: dttRoutine
	serialIndex    int
	serialTask     *DeliverTxTask
	serialCh       chan *DeliverTxTask
	//serialNextCh       chan *DeliverTxTask

	mtx       sync.Mutex
	currTxFee sdk.Coins

	txResponses []*abci.ResponseDeliverTx
	invalidTxs int
	app         *BaseApp
}

func NewDTTManager(app *BaseApp) *DTTManager {
	dttm := &DTTManager{
		app: app,
	}
	dttm.dttRoutineList = make([]*dttRoutine, 0, maxDeliverTxsConcurrentNum) //sync.Map{}
	for i := 0; i < maxDeliverTxsConcurrentNum; i++ {
		dttr := newDttRoutine(int8(i), dttm.concurrentBasic, dttm.runConcurrentAnte)
		dttr.setLogger(dttm.app.logger)
		dttm.dttRoutineList = append(dttm.dttRoutineList, dttr)
		//dttm.app.logger.Info("newDttRoutine", "index", i, "list", len(dttm.dttRoutineList))
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
	dttm.app.logger.Info("TotalTxs", "count", dttm.totalCount)
	totalSerialWaitingCount += dttm.totalCount

	dttm.txs = txs
	//dttm.tasks = make([]*DeliverTxTask, len(txs))
	dttm.currTxFee = sdk.Coins{}
	dttm.serialTask = nil
	dttm.serialIndex = -1
	dttm.serialCh = make(chan *DeliverTxTask, 2)
	//dttm.serialNextCh = make(chan *DeliverTxTask, 1)
	dttm.startFinished = false

	//start := time.Now()
	//dttm.preloadSender(txs)
	//totalPreloadConDuration += time.Since(start).Microseconds()
	//logger.Info("DeliverTxs duration", "preload", totalPreloadConDuration)
	//dttm.app.logger.Error("preloadFinished")

	dttm.txResponses = make([]*abci.ResponseDeliverTx, len(txs))
	dttm.invalidTxs = 0

	go dttm.serialRoutine()
	//go dttm.serialNextRoutine()

	//dttm.dttRoutineList = make([]*dttRoutine, 0, maxDeliverTxsConcurrentNum) //sync.Map{}
	for i := 0; i < maxDeliverTxsConcurrentNum; i++ {
		dttr := dttm.dttRoutineList[i]

		//dttm.app.logger.Info("StartDttRoutine", "index", i, "list", len(dttm.dttRoutineList))
		err := dttr.OnStart()
		if err != nil {
			dttm.app.logger.Error("Error starting DttRoutine", "err", err)
		}
		dttr.makeNewTask(txs[i], i)
		//time.Sleep(1 * time.Millisecond)
	}
	dttm.startFinished = true
}

func (dttm *DTTManager) preloadSender(txs [][]byte) {
	checkStateCtx := dttm.app.checkState.ctx.WithBlockHeight(dttm.app.checkState.ctx.BlockHeight() + 1)

	maxNums := runtime.NumCPU()
	txSize := len(txs)
	if maxNums > txSize {
		maxNums = txSize
	}
	dttm.app.logger.Error("preloadStart", "maxNum", maxNums)

	txJobChan := make(chan []byte)
	//txJobChan := make(chan *txBytesWithTxIndex)
	var wg sync.WaitGroup
	wg.Add(txSize)

	for index := 0; index < maxNums; index++ {
		go func(ch chan []byte, wg *sync.WaitGroup) {
			for txBytes := range ch {
				var realTx sdk.Tx
				var err error
				if mem := GetGlobalMempool(); mem != nil {
					realTx, _ = mem.ReapEssentialTx(txBytes).(sdk.Tx)
				}
				if realTx == nil {
					realTx, err = dttm.app.txDecoder(txBytes)
				}
				//tx, err := dttm.app.txDecoder(tbi.txBytes)
				//task := newDeliverTxTask(realTx, tbi.index)
				//task.info.txBytes = tbi.txBytes
				if err == nil {
					dttm.app.getTxFee(checkStateCtx.WithTxBytes(txBytes), realTx)
					//task.fee, task.isEvm, task.from = dttm.app.getTxFeeAndFromHandler(checkStateCtx.WithTxBytes(tbi.txBytes), realTx)
					//dttm.app.logger.Info("preload", "from", from)
					//task.fee, task.isEvm, task.from = dttm.app.getTxFeeAndFromHandler(checkStateCtx.WithTxBytes(txBytes), task.info.tx)
				//} else {
				//	task.err = err
				//	task.setStep(partialConcurrentStepBasicFailed)
				}
				//dttm.setTask(task)

				wg.Done()
			}
		}(txJobChan, &wg)
	}
	for _, v := range txs {
		txJobChan <- v//&txBytesWithTxIndex{index: k, txBytes: v}
	}

	wg.Wait()
	close(txJobChan)
}

func (dttm *DTTManager) concurrentBasic(txByte []byte, index int) *DeliverTxTask {
	dttm.app.logger.Info("concurrentBasic", "index", index)
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
		task.setStep(partialConcurrentStepBasicFailed)
	} else {
		//dttm.app.logger.Info("hasExistPrevTask", "index", task.index, "from", task.from)
		task.setStep(partialConcurrentStepBasicSucceed)
		if dttm.serialTask == nil && task.index == dttm.serialIndex+1 {
			return task
		}
		// need to check whether exist running tx who has the same from but smaller txIndex
		//count := len(dttm.dttRoutineList)
		//for i := 0; i < count; i++ {
		//	dttr := dttm.dttRoutineList[i]
		//	//dttm.app.logger.Info("hasExistPrevTask 1", "routine", dttr.index, "task", dttr.txIndex)
		//	task.needToRerun = dttr.hasExistPrevTask(task.from, task.index)
		//	if task.needToRerun {
		//		//dttm.app.logger.Error("needToRerun 3", "index", task.index, "conflicted", dttr.task.index)
		//		break
		//	}
		//}
	}
	return task
}

func (dttm *DTTManager) runConcurrentAnte(task *DeliverTxTask) error {
	if dttm.app.anteHandler == nil {
		return fmt.Errorf("anteHandler cannot be nil")
	}

	//if task.needToRerun {
	//	dttm.app.logger.Error("ResetContext", "index", task.index)

		task.info.ctx = dttm.app.getContextForTx(runTxModeDeliverPartConcurrent, task.info.txBytes) // same context for all txs in a block
		task.setUpdateCount(0)
		//task.setStep(partialConcurrentStepBasicSucceed)
		//task.needToRerun = false
		task.canRerun = 0
	//}

	task.setStep(partialConcurrentStepAnteStart)
	//if global.GetGlobalHeight() == 5811244 {
		dttm.app.logger.Info("RunAnte", "index", task.index, "routine", task.routineIndex, "addr", task.from)
	//}

	task.info.ctx = task.info.ctx.WithCache(sdk.NewCache(dttm.app.blockCache, useCache(runTxModeDeliverPartConcurrent))) // one cache for a tx

	dttm.accountUpdated(false, 2, task.from)
	anteStart := time.Now()
	err := dttm.runAnte(task)
	totalAnteDuration += time.Since(anteStart).Microseconds()
	if err != nil {
		dttm.app.logger.Error("anteFailed", "index", task.index, "err", err)
		//task.anteFailed = true
		task.setStep(partialConcurrentStepAnteFailed)
	} else {
		//dttm.app.logger.Info("AnteSucceed", "index", task.index)
		task.setStep(partialConcurrentStepAnteSucceed)
	}
	task.err = err
	//count := len(dttm.dttRoutineList)
	//for i := 0; i < count; i++ {
	//	dttr := dttm.dttRoutineList[i]
	//	if dttr == nil {
	//		continue
	//	}
	//	conflicted := dttr.checkConflict(task.from, task.index)
	//	if conflicted {
	//		dttr.logger.Error("needToRerunFromAnte", "index", task.index, "conflicted", dttr.task.index)
	//		task.needToRerun = true
	//	}
	//}
	if task.canRerun > 0 {
		dttr := dttm.dttRoutineList[task.routineIndex]
		//go func() {
			dttr.logger.Error("rerunChInFromAnte", "index", task.index)
		//dttr.task.needToRerun = true
		dttr.shouldRerun()
		//}()
	} else if dttm.serialIndex+1 == task.index && dttm.serialTask == nil {
			//go func() {
			//dttm.app.logger.Info("ExtractNextSerialFromAnte", "index", task.index)
			dttm.serialCh <- task
			//}()
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

//func (dttm *DTTManager) serialNextRoutine() {
//	for {
//		select {
//		case task := <-dttm.serialNextCh:
//			dttm.serialCh <- task
//		}
//	}
//}

func (dttm *DTTManager) serialRoutine() {
	for {
		select {
		case task := <-dttm.serialCh:
			// runMsgs etc.
			if task.index == dttm.serialIndex+1 {
				dttm.serialIndex = task.index
				dttm.serialTask = task
				dttm.serialTask.setStep(partialConcurrentStepSerialExecute)
				dttm.serialExecution()
				dttm.serialTask = nil
				task.setStep(partialConcurrentStepFinished)
				//if global.GetGlobalHeight() == 5811244 {
				//	dttm.app.logger.Info("NextSerialTask", "index", dttm.serialIndex+1)
				//}

				if dttm.serialIndex == dttm.totalCount-1 {
					//dttm.app.logger.Info("TotalTxFeeForCollector", "fee", dttm.currTxFee)
					count := len(dttm.dttRoutineList)
					for i := 0; i < count; i++ {
						dttr := dttm.dttRoutineList[i]
						dttr.stop()
					}

					dttm.done <- 0
					go func() {
						close(dttm.serialCh)
						//close(dttm.serialNextCh)
					}()
					return
				}

				// make new task for this routine
				dttr := dttm.dttRoutineList[task.routineIndex]
				nextIndex := maxDeliverTxsConcurrentNum + task.index
				if dttr != nil && nextIndex < dttm.totalCount {
					//if !dttm.startFinished {
					//	time.Sleep(maxDeliverTxsConcurrentNum * time.Millisecond)
					//}
					dttr.makeNewTask(dttm.txs[nextIndex], nextIndex)
				}

				// todo: check whether there are ante-finished task
				count := len(dttm.dttRoutineList)
				var nextTask *DeliverTxTask
				var rerunRoutine *dttRoutine
				//getRerun := false
				for i := 0; i < count; i++ {
					//if getRerun {
					//	break
					//}
					dttr = dttm.dttRoutineList[i]
					if dttr.task == nil || dttr.task.index <= task.index || dttr.task.step == partialConcurrentStepFinished || dttr.task.step == partialConcurrentStepBasicFailed {
						continue
					}
					if dttr.task.from == task.from {
						if rerunRoutine == nil {
							rerunRoutine = dttr
						} else if dttr.task.index < rerunRoutine.task.index {
							rerunRoutine = dttr
						}
					} else if dttr.task.index == dttm.serialIndex+1 {
						step := dttr.task.getStep()
						if step == partialConcurrentStepBasicFailed ||
							step == partialConcurrentStepAnteFailed ||
							step == partialConcurrentStepAnteSucceed {
							nextTask = dttr.task
						}
					}

					//// if exists the next task which has finished the concurrent execution
					//if dttr.task.index == dttm.serialIndex+1 {
					//	//dttm.app.logger.Info("WaitNextSerialTask", "index", dttr.task.index, "needToRerun", dttr.task.needToRerun, "step", dttr.task.step)
					//	step := dttr.task.getStep()
					//	if dttr.task.from == task.from {
					//		//go func() {
					//		getRerun = true
					//		dttr.logger.Error("rerunCh", "index", dttr.task.index)
					//		dttr.task.needToRerun = true
					//		dttr.rerunCh <- 0
					//		//}()
					//	} else if dttr.task.needToRerun {
					//		//dttm.app.logger.Info("NeedToWaitRerun", "index", dttr.task.index)
					//		dttr.couldRerun(task.index)
					//	} else if step == partialConcurrentStepBasicFailed ||
					//		step == partialConcurrentStepAnteFailed ||
					//		step == partialConcurrentStepAnteSucceed {
					//		nextTask = dttr.task
					//	}
					//} else if dttr.task.from == task.from {
					//	if rerunRoutine == nil {
					//		rerunRoutine = dttr
					//	} else if dttr.task.index < rerunRoutine.task.index {
					//		rerunRoutine = dttr
					//	}
					//}
				}

				if rerunRoutine != nil {
					//rerunRoutine.couldRerun(task.index)
					dttm.app.logger.Error("rerunRoutine", "index", rerunRoutine.task.index, "serial", task.index)
					rerunRoutine.shouldRerun()
				}
				if nextTask != nil {
					totalSerialWaitingCount--
					go func() {
						//	dttm.app.logger.Info("ExtractNextSerialFromSerial", "index", nextTask.index)
						dttm.serialCh <- nextTask
					}()
					//dttm.serialNextCh <- nextTask
				}
				//} else {
				//	panic(fmt.Sprintf("invalid index for serial execution: expected %x, got %x\n", dttm.serialIndex, task.index))
			}
		}
	}
}

func (dttm *DTTManager) serialExecution() {
	//if global.GetGlobalHeight() == 5811244 {
		dttm.app.logger.Info("RunStatefulSerialRoutine", "index", dttm.serialTask.index)
	//}

	info := dttm.serialTask.info
	handler := info.handler

	handleGasFn := func() {
		gasStart := time.Now()

		dttm.updateFeeCollector()

		//dttm.app.logger.Info("handleDeferRefund", "index", dttm.serialTask.txIndex, "addr", dttm.serialTask.from)
		dttm.accountUpdated(false, 1, dttm.serialTask.from)
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
		if txRs.Code != abci.CodeTypeOK {
			//logger.Debug("Invalid tx", "code", txRs.Code, "log", txRs.Log, "index", dttm.serialTask.index)
			dttm.invalidTxs++
		}
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
	dttm.accountUpdated(false, 2, dttm.serialTask.from)
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
	//if global.GetGlobalHeight() == 5811244 && hex.EncodeToString(acc.GetAddress()) == "4ce08ffc090f5c54013c62efe30d62e6578e738d" {
	//	dttm.app.logger.Error("OnAccountUpdated", "addr", addr)
	//}
	//waitingIndex := -1
	//if dttm.serialTask == nil {
	//	waitingIndex = dttm.serialIndex + 1
	//}
	dttm.accountUpdated(true, 1, addr)
}

func (dttm *DTTManager) accountUpdated(happened bool, times int8, address string) {
	num := len(dttm.dttRoutineList)
	//if global.GetGlobalHeight() == 5811244 && address == "ex1fnsgllqfpaw9gqfuvth7xrtzuetcuuudrhc557" {
	//	dttm.app.logger.Error("OnAccountUpdated", "times", times, "happened", happened)
	//}
	for i := 0; i < num; i++ {
		dttr := dttm.dttRoutineList[i]
		if dttr.task == nil || !dttr.task.needToRerunWhenContextChanged() || dttr.task.from != address {
			continue
		}

		task := dttr.task
		count := task.getUpdateCount()
		if happened {
			task.setUpdateCount(count + times)
			// todo: whether should rerun the task
			if task.getUpdateCount() > 0 {
				//go func() {
				dttm.app.logger.Error("accountUpdatedToRerun", "index", task.index, "step", task.getStep())
				//dttr.task.needToRerun = true
				dttr.shouldRerun()
				//}()
			}
		} else {
			task.setUpdateCount(count - times)
		}
	}
}

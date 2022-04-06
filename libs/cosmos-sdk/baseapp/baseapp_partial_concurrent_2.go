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

type (
	dttRoutineStep uint8
)
const (
	dttRoutineStepNone dttRoutineStep = iota
	dttRoutineStepStart
	//dttRoutineStepBasic
	dttRoutineStepAnteStart
	dttRoutineStepAnteFinished
	//dttRoutineStepNeedRerun
	dttRoutineStepReadyForSerial
	dttRoutineStepSerial
	dttRoutineStepFinished
)

//-------------------------------------
type BasicProcessFn func(txByte []byte, index int) *DeliverTxTask
type RunAnteFn func(task *DeliverTxTask) error

type dttRoutine struct {
	//service.BaseService
	done    chan int8
	task    *DeliverTxTask
	txByte  chan []byte
	txIndex int
	rerunCh chan int8
	//runAnteCh chan int
	index int8
	//mtx   sync.Mutex
	step dttRoutineStep
	needToRerun bool

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
	dttr.step = dttRoutineStepStart
	dttr.txIndex = index
	dttr.needToRerun = false
	//dttr.task = nil
	//dttr.logger.Info("makeNewTask", "index", dttr.txIndex)
	//go func() {
	dttr.txByte <- txByte
	//}()
}

func (dttr *dttRoutine) OnStart() {
	dttr.done = make(chan int8)
	dttr.txByte = make(chan []byte, 3)
	dttr.rerunCh = make(chan int8, 5)
	//dttr.runAnteCh = make(chan int, 5)
	dttr.step = dttRoutineStepNone
	go dttr.executeTaskRoutine()
}

func (dttr *dttRoutine) stop() {
	dttr.step = dttRoutineStepNone
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
			//dttr.step = dttRoutineStepBasic
			dttr.task = dttr.basicProFn(tx, dttr.txIndex)
			dttr.task.routineIndex = dttr.index
			if dttr.task.err == nil {
				dttr.step = dttRoutineStepAnteStart
				dttr.runAnteFn(dttr.task)
				dttr.step = dttRoutineStepAnteFinished
			} else {
				dttr.step = dttRoutineStepReadyForSerial
			}
		case <-dttr.rerunCh:
			//dttr.logger.Error("readRerunCh", "index", dttr.task.index, "step", dttr.step)
			//step := dttr.task.step
			dttr.task.prevTaskIndex = -1
			if dttr.step == dttRoutineStepAnteFinished {
				//dttr.logger.Error("RerunTask", "index", dttr.task.index)
				dttr.needToRerun = false
				dttr.step = dttRoutineStepAnteStart
				dttr.runAnteFn(dttr.task)
				dttr.step = dttRoutineStepAnteFinished
			} else if dttr.step == dttRoutineStepReadyForSerial ||
				dttr.step == dttRoutineStepSerial ||
				dttr.step == dttRoutineStepFinished {
				dttr.needToRerun = false
				dttr.logger.Error("task is empty or finished")
			} else {
				//dttr.logger.Error("shouldRerunLater", "index", dttr.task.index)
				// maybe the task is in other condition, running concurrent execution or running make new task.
				dttr.task.canRerun++
				dttr.needToRerun = false
			}
		}
	}
}

func (dttr *dttRoutine) shouldRerun(fromIndex int) {
	if dttr.step == dttRoutineStepReadyForSerial || dttr.needToRerun == true || (dttr.task.prevTaskIndex >= 0 && dttr.task.prevTaskIndex > fromIndex) {
		//dttr.logger.Error("willnotRerun", "index", dttr.task.index, "prev", dttr.task.prevTaskIndex, "from", fromIndex, "step", dttr.step, "needToRerun", dttr.needToRerun)
		return
	}
	if dttr.step == dttRoutineStepAnteStart || dttr.step == dttRoutineStepAnteFinished {
		//dttr.logger.Error("shouldRerun", "index", dttr.task.index, "from", fromIndex, "step", dttr.step, "needToRerun", dttr.needToRerun)
		//dttr.step = dttRoutineStepNeedRerun
		dttr.needToRerun = true
		//go func() {
			dttr.rerunCh <- 0 // todo: maybe blocked for several milliseconds. why?
			//dttr.logger.Error("sendRerunCh", "index", dttr.task.index)
		//}()
	}
}

func (dttr *dttRoutine) needToRerunWhenContextChanged() bool {
	switch dttr.step {
	case dttRoutineStepNone:
		fallthrough
	case dttRoutineStepStart:
		fallthrough
	case dttRoutineStepReadyForSerial:
		fallthrough
	case dttRoutineStepSerial:
		fallthrough
	case dttRoutineStepFinished:
		return false
	}
	return true
}

func (dttr *dttRoutine) readyForSerialExecution() bool {
	if dttr.task == nil || dttr.needToRerun || dttr.task.canRerun > 0 || dttr.task.prevTaskIndex >= 0 {
		return false
	}
	switch dttr.step {
	case dttRoutineStepNone:
		fallthrough
	case dttRoutineStepStart:
		fallthrough
	//case dttRoutineStepBasic:
	//	fallthrough
	case dttRoutineStepAnteStart:
		fallthrough
	//case dttRoutineStepNeedRerun:
	//	fallthrough
	case dttRoutineStepFinished:
		fallthrough
	case dttRoutineStepSerial:
		return false
	//case dttRoutineStepAnteFinished:
	//	fallthrough
	//case dttRoutineStepReadyForSerial:
	//	return true
	}
	return true
}

//-------------------------------------

type DTTManager struct {
	done       chan int8
	totalCount int
	txs        [][]byte
	//tasks []*DeliverTxTask
	startFinished  bool
	dttRoutineList []*dttRoutine //sync.Map	// key: txIndex, value: dttRoutine
	serialIndex    int
	serialTask     *DeliverTxTask
	serialCh       chan int8 //*DeliverTxTask//
	//serialNextCh       chan *DeliverTxTask

	mtx       sync.Mutex
	currTxFee sdk.Coins

	txResponses   []*abci.ResponseDeliverTx
	invalidTxs    int
	app           *BaseApp
	checkStateCtx sdk.Context
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
	}

	return dttm
}

func (dttm *DTTManager) deliverTxs(txs [][]byte) {
	dttm.done = make(chan int8, 1)

	dttm.totalCount = len(txs)
	dttm.app.logger.Info("TotalTxs", "count", dttm.totalCount)
	totalSerialWaitingCount += dttm.totalCount

	dttm.txs = txs
	//dttm.tasks = make([]*DeliverTxTask, len(txs))
	dttm.currTxFee = sdk.Coins{}
	dttm.serialTask = nil
	dttm.serialIndex = -1
	//dttm.serialCh = make(chan *DeliverTxTask, 2)
	dttm.serialCh = make(chan int8, 3)
	//dttm.serialNextCh = make(chan *DeliverTxTask, 1)
	dttm.startFinished = false

	//dttm.preloadSender(txs)

	dttm.txResponses = make([]*abci.ResponseDeliverTx, len(txs))
	dttm.invalidTxs = 0

	dttm.checkStateCtx = dttm.app.checkState.ctx.WithBlockHeight(dttm.app.checkState.ctx.BlockHeight() + 1)

	go dttm.serialRoutine()
	//go dttm.serialNextRoutine()

	//dttm.dttRoutineList = make([]*dttRoutine, 0, maxDeliverTxsConcurrentNum) //sync.Map{}
	for i := 0; i < maxDeliverTxsConcurrentNum; i++ {
		dttr := dttm.dttRoutineList[i]

		//dttm.app.logger.Info("StartDttRoutine", "index", i, "list", len(dttm.dttRoutineList))
		dttr.OnStart()
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
		txJobChan <- v //&txBytesWithTxIndex{index: k, txBytes: v}
	}

	wg.Wait()
	close(txJobChan)
}

func (dttm *DTTManager) concurrentBasic(txByte []byte, index int) *DeliverTxTask {
	//dttm.app.logger.Info("RunBasic", "index", index)
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
		////dm.app.logger.Error("tx decode failed", "err", err)
		//task.step = partialConcurrentStepBasicFailed
		////task.setStep(partialConcurrentStepBasicFailed)
		return task
	}

	task.info.handler = dttm.app.getModeHandler(runTxModeDeliverPartConcurrent)                 //dm.handler
	//task.info.ctx = dttm.app.getContextForTx(runTxModeDeliverPartConcurrent, task.info.txBytes) // same context for all txs in a block
	//task.resetUpdateCount()
	//task.fee, task.isEvm, task.from = dttm.app.getTxFeeAndFromHandler(task.info.ctx, task.info.tx)
	task.fee, task.isEvm, task.from = dttm.app.getTxFeeAndFromHandler(dttm.checkStateCtx, task.info.tx)

	if err = validateBasicTxMsgs(task.info.tx.GetMsgs()); err != nil {
		task.err = err
		dttm.app.logger.Error("validateBasicTxMsgs failed", "err", err)
		//task.setStep(partialConcurrentStepBasicFailed)
	} else {
		//dttm.app.logger.Info("hasExistPrevTask", "index", task.index, "from", task.from)
		//task.setStep(partialConcurrentStepBasicSucceed)
		if dttm.serialTask == nil && task.index == dttm.serialIndex+1 {
			return task
		}
	}
	return task
}

func (dttm *DTTManager) runConcurrentAnte(task *DeliverTxTask) error {
	if dttm.app.anteHandler == nil {
		return fmt.Errorf("anteHandler cannot be nil")
	}

	count := len(dttm.dttRoutineList)
	for i := 0; i < count; i++ {
		dttr := dttm.dttRoutineList[i]
		if dttr.task == nil ||
			dttr.txIndex != dttr.task.index ||
			dttr.step == dttRoutineStepNone || dttr.step == dttRoutineStepStart ||
			//dttr.step == dttRoutineStepReadyForSerial ||
			dttr.step == dttRoutineStepFinished || dttr.step == dttRoutineStepReadyForSerial ||
			dttr.task.from != task.from {
			continue
		}
		if dttr.task.index < task.index && dttr.task.index > task.prevTaskIndex {
			task.prevTaskIndex = dttr.task.index
			//dttm.app.logger.Error("hasExistPrevTask1", "index", task.index, "prev", task.prevTaskIndex, "prevStep", dttr.step, "from", task.from)
		} else if dttr.task.index > task.index && (dttr.task.prevTaskIndex < 0 || dttr.task.prevTaskIndex < task.index) {
			dttr.task.prevTaskIndex = task.index
			//dttm.app.logger.Error("hasExistPrevTask2", "index", dttr.task.index, "prev", dttr.task.prevTaskIndex, "from", task.from)
		}
	}
	if task.prevTaskIndex >= 0 || task.index <= dttm.serialIndex {//|| dttr.needToRerun {
		return nil
	}

	task.info.ctx = dttm.app.getContextForTx(runTxModeDeliverPartConcurrent, task.info.txBytes) // same context for all txs in a block
	//task.setStep(partialConcurrentStepAnteStart)
	task.resetUpdateCount()
	task.canRerun = 0

	//if global.GetGlobalHeight() == 5811070 {
	//	dttm.app.logger.Info("RunAnte", "index", task.index, "routine", task.routineIndex, "addr", task.from)
	//}

	task.info.ctx = task.info.ctx.WithCache(sdk.NewCache(dttm.app.blockCache, useCache(runTxModeDeliverPartConcurrent))) // one cache for a tx

	//dttm.accountUpdated(false, 2, task.from)
	anteStart := time.Now()
	err := dttm.runAnte(task)
	totalAnteDuration += time.Since(anteStart).Microseconds()
	if err != nil {
		dttm.app.logger.Error("anteFailed", "index", task.index, "err", err)
	}
	task.err = err
	// need to check whether exist running tx who has the same sender but smaller txIndex
	//count := len(dttm.dttRoutineList)
	//for i := 0; i < count; i++ {
	//	dttr := dttm.dttRoutineList[i]
	//	if dttr.task == nil ||
	//		//dttr.txIndex != dttr.task.index ||
	//		//dttr.step == dttRoutineStepNone || dttr.step == dttRoutineStepStart ||
	//		dttr.step == dttRoutineStepReadyForSerial ||
	//		//dttr.step == dttRoutineStepFinished || dttr.step == dttRoutineStepReadyForSerial ||
	//		dttr.task.from != task.from {
	//		continue
	//	}
	//	if dttr.task.index < task.index && dttr.task.index > task.prevTaskIndex {
	//		task.prevTaskIndex = dttr.task.index
	//		dttm.app.logger.Error("hasExistPrevTask", "index", task.index, "prev", task.prevTaskIndex, "prevStep", dttr.step, "from", task.from)
	//	} else if dttr.task.index > task.index && (dttr.task.prevTaskIndex < 0 || dttr.task.prevTaskIndex < task.index) {
	//		dttr.task.prevTaskIndex = task.index
	//		dttm.app.logger.Error("hasExistPrevTask", "index", dttr.task.index, "prev", dttr.task.prevTaskIndex, "from", task.from)
	//	}
	//}
	//dttr := dttm.dttRoutineList[task.routineIndex]
	//dttm.app.logger.Info("AnteFinished", "index", task.index, "step", dttr.step, "prev", task.prevTaskIndex, "canRerun", task.canRerun)
	//if task.prevTaskIndex > 0 || dttr.needToRerun {
	//	return err
	//}
	dttr := dttm.dttRoutineList[task.routineIndex]
	//dttm.app.logger.Info("AnteFinished", "index", task.index, "step", dttr.step, "prev", task.prevTaskIndex, "canRerun", task.canRerun)

	if task.canRerun > 0 {
		//go func() {
		//dttr.logger.Error("rerunChInFromAnte", "index", task.index)
		//dttr.task.needToRerun = true
		dttr.shouldRerun(-1)
		//}()
	} else if dttm.serialIndex+1 == task.index && dttm.serialTask == nil && !dttr.needToRerun && task.prevTaskIndex < 0 {
		//go func() {
		//dttm.app.logger.Info("ExtractNextSerialFromAnte", "index", dttr.task.index, "step", dttr.step, "needToRerun", dttr.needToRerun)
			dttm.serialCh <- task.routineIndex
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

func (dttm *DTTManager) serialRoutine() {
	for {
		select {
		//case task := <-dttm.serialCh:
		case routineIndex := <-dttm.serialCh:
			// runMsgs etc.
			rt := dttm.dttRoutineList[routineIndex]
			task := rt.task
			if task.index == dttm.serialIndex+1 {//&& (rt.step == dttRoutineStepReadyForSerial || rt.step == dttRoutineStepAnteFinished) {
				dttm.serialIndex = task.index
				dttm.serialTask = task
				//dttm.serialTask.setStep(partialConcurrentStepSerialExecute)
				rt.step = dttRoutineStepSerial
				dttm.serialExecution()
				rt.step = dttRoutineStepFinished
				dttm.serialTask = nil
				//task.setStep(partialConcurrentStepFinished)
				//if global.GetGlobalHeight() == 5811070 {
				//	dttm.app.logger.Info("NextSerialTask", "index", dttm.serialIndex+1)
				//}

				if dttm.serialIndex == dttm.totalCount-1 {
					//dttm.app.logger.Info("TotalTxFeeForCollector", "fee", dttm.currTxFee)
					count := len(dttm.dttRoutineList)
					for i := 0; i < count; i++ {
						dttr := dttm.dttRoutineList[i]
						dttr.stop()
					}

					dttm.updateFeeCollector()

					dttm.done <- 0
					//go func() {
						close(dttm.serialCh)
						//close(dttm.serialNextCh)
					//}()
					return
				}

				// make new task for this routine
				nextIndex := maxDeliverTxsConcurrentNum + task.index
				if nextIndex < dttm.totalCount {
					//if !dttm.startFinished {
					//	time.Sleep(maxDeliverTxsConcurrentNum * time.Millisecond)
					//}
					rt.makeNewTask(dttm.txs[nextIndex], nextIndex)
				}

				// check whether there are ante-finished task
				count := len(dttm.dttRoutineList)
				nextTaskRoutine := int8(-1)
				var rerunRoutine *dttRoutine
				for i := 0; i < count; i++ {
					dttr := dttm.dttRoutineList[i]
					if dttr.task == nil || dttr.task.index <= task.index || dttr.step < dttRoutineStepAnteStart {
						continue
					}
					if dttr.task.prevTaskIndex == task.index || dttr.task.from == task.from {
						if dttr.task.prevTaskIndex < task.index {
							dttr.task.prevTaskIndex = task.index
						}
						if rerunRoutine == nil {
							rerunRoutine = dttr
						} else if dttr.task.index < rerunRoutine.task.index {
							rerunRoutine = dttr
						}
					} else if dttr.task.index == dttm.serialIndex+1 {
						if dttr.readyForSerialExecution() {
							nextTaskRoutine = dttr.index
							totalSerialWaitingCount--
							//go func() {
							//dttm.app.logger.Info("ExtractNextSerialFromSerial", "index", dttr.task.index, "step", dttr.step, "needToRerun", dttr.needToRerun)
							dttm.serialCh <- nextTaskRoutine //nextTask//
							//}()
						} else {
							dttm.app.logger.Error("NotReadyForSerial", "index", dttr.task.index, "step", dttr.step, "needToRerun", dttr.needToRerun, "canRerun", dttr.task.canRerun, "prev", dttr.task.prevTaskIndex)
						}
					}
				}

				if rerunRoutine != nil {
					//dttm.app.logger.Error("rerunRoutine", "index", rerunRoutine.task.index, "serial", task.index)
					rerunRoutine.shouldRerun(task.index)
				}
			}
		}
	}
}

func (dttm *DTTManager) serialExecution() {
	//if global.GetGlobalHeight() == 5811070 {
	//	dttm.app.logger.Info("SerialStart", "index", dttm.serialTask.index)
	//}

	info := dttm.serialTask.info
	handler := info.handler

	handleGasFn := func() {
		gasStart := time.Now()

		//dttm.updateFeeCollector()

		//dttm.app.logger.Info("handleDeferRefund", "index", dttm.serialTask.txIndex, "addr", dttm.serialTask.from)
		//dttm.accountUpdated(false, 1, dttm.serialTask.from)
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
		//dttm.app.logger.Info("SerialFinished", "index", dttm.serialTask.index, "routine", dttm.serialTask.routineIndex)
		dttm.txResponses[dttm.serialTask.index] = &txRs
		if txRs.Code != abci.CodeTypeOK {
			//logger.Debug("Invalid tx", "code", txRs.Code, "log", txRs.Log, "index", dttm.serialTask.index)
			dttm.invalidTxs++
		}
	}

	// execute anteHandler failed
	if dttm.serialTask.err != nil {
		//dttm.app.logger.Error("RunSerialFinished", "index", dttm.serialTask.index, "err", dttm.serialTask.err)
		txRs := sdkerrors.ResponseDeliverTx(dttm.serialTask.err, 0, 0, dttm.app.trace) //execResult.GetResponse()
		execFinishedFn(txRs)
		return
	}

	////dttm.app.logger.Info("WriteAnteCache", "index", dttm.serialTask.txIndex)
	////ctx, msCacheAnte := dttm.app.cacheTxContext(dttm.app.getContextForTx(runTxModeDeliver, []byte{}), []byte{})
	dttm.calculateFeeForCollector(dttm.serialTask.fee, true)
	if err := dttm.app.updateFeeCollectorAccHandler(info.ctx, dttm.currTxFee); err != nil {
		panic(err)
	}
	////cache.Write()
	info.msCacheAnte.Write()
	info.ctx.Cache().Write(true)
	//dttm.calculateFeeForCollector(dttm.serialTask.fee, true)

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
	//dttm.accountUpdated(false, 2, dttm.serialTask.from)
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
	//dttm.mtx.Lock()
	//defer dttm.mtx.Unlock()

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

func (dttm *DTTManager) OnAccountUpdated(acc exported.Account, updateState bool) {
	start := time.Now()
	//if global.GetGlobalHeight() == 5811244 && updateState && hex.EncodeToString(acc.GetAddress()) == "4ce08ffc090f5c54013c62efe30d62e6578e738d" {
	//	dttm.app.logger.Error("OnAccountUpdated", "updateState", updateState)
	//}
	if updateState {
		addr := acc.GetAddress().String()
		// called twice on each handleRunMsg
		// addr=ex14h6fzmg37df2yaywr8es2epgvwf38ahpdq367z hex=adf4916d11f352a2748e19f3056428639313f6e1
		//dttm.app.logger.Info("CallAccountUpdated", "addr", addr, "hex", hex.EncodeToString(acc.GetAddress()))
		dttm.accountUpdated(true, 1, addr)
	}
	totalAccountUpdateDuration += time.Since(start).Microseconds()
}

func (dttm *DTTManager) accountUpdated(happened bool, times int8, address string) {
	dttm.mtx.Lock()
	defer dttm.mtx.Unlock()

	num := len(dttm.dttRoutineList)
	//if global.GetGlobalHeight() == 5811070 && address == "ex1fnsgllqfpaw9gqfuvth7xrtzuetcuuudrhc557" {
	//	dttm.app.logger.Info("OnAccountUpdated", "times", times, "happened", happened, "addr", address)
	//}
	for i := 0; i < num; i++ {
		dttr := dttm.dttRoutineList[i]
		if dttr.task == nil || dttr.txIndex != dttr.task.index || !dttr.needToRerunWhenContextChanged() || dttr.task.from != address {
			//if global.GetGlobalHeight() == 5811244 && (dttr.task.index == 3 || address == "ex1fnsgllqfpaw9gqfuvth7xrtzuetcuuudrhc557") {
			//	dttm.app.logger.Error("NoNeedToRerunFromObserver", "index", dttr.task.index, "txIndex", dttr.txIndex, "setp", dttr.step, "from", dttr.task.from)
			//}
			continue
		}

		task := dttr.task
		//count := task.getUpdateCount()
		if task.setUpdateCount(times, happened) {
			dttr.shouldRerun(-1)
			//if dttr.shouldRerun(-1) {
			//	dttm.app.logger.Error("accountUpdatedToRerun", "index", task.index, "step", task.getStep())
			//}
		}
	}
}

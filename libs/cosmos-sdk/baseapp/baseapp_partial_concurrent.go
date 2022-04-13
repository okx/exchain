package baseapp

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/trace"

	"time"
)

var totalSerialWaitingCount = 0

type (
	dttRoutineStep uint8
)

const (
	dttRoutineStepNone dttRoutineStep = iota
	dttRoutineStepStart
	dttRoutineStepAnteStart
	dttRoutineStepAnteFinished
	dttRoutineStepReadyForSerial
	dttRoutineStepSerial
	dttRoutineStepFinished

	maxDeliverTxsConcurrentNum = 4
	keepAliveIntervalMS        = 1
)

var totalAnteDuration = int64(0)
var totalSerialDuration = int64(0)
var totalWriteTime = int64(0)
var totalDeferGasTime = int64(0)
var totalHandleGasTime = int64(0)
var totalRunMsgsTime = int64(0)
var totalWaitingTime = int64(0)
var totalBasicTime = int64(0)
var totalPreloadConDuration = int64(0)
var totalAccountUpdateDuration = int64(0)

type DeliverTxTask struct {
	index        int
	canRerun     int8
	routineIndex int8

	info          *runTxInfo
	fee           sdk.Coins
	from          string
	to            string
	err           error
	prevTaskIndex int // the index of a running tx with a smaller index while its from or to equals to this tx's from
}

func newDeliverTxTask(tx sdk.Tx, index int) *DeliverTxTask {
	t := &DeliverTxTask{
		index:         index,
		info:          &runTxInfo{tx: tx},
		prevTaskIndex: -1,
	}

	return t
}

//-------------------------------------
type BasicProcessFn func(txByte []byte, index int) *DeliverTxTask
type RunAnteFn func(task *DeliverTxTask) error

type dttRoutine struct {
	done        chan int8
	task        *DeliverTxTask
	txByte      chan []byte
	txIndex     int
	rerunCh     chan int8
	index       int8
	step        dttRoutineStep
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
	dttr.txByte <- txByte
}

func (dttr *dttRoutine) OnStart() {
	dttr.done = make(chan int8)
	dttr.txByte = make(chan []byte, 3)
	dttr.rerunCh = make(chan int8, 5)
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
			dttr.task = dttr.basicProFn(tx, dttr.txIndex)
			dttr.task.routineIndex = dttr.index
			if dttr.task.err == nil {
				dttr.runAnteFn(dttr.task)
			} else {
				dttr.step = dttRoutineStepReadyForSerial
			}
		case <-dttr.rerunCh:
			//dttr.logger.Error("readRerunCh", "index", dttr.task.index, "step", dttr.step)
			dttr.task.prevTaskIndex = -1
			if dttr.step == dttRoutineStepAnteFinished {
				//dttr.logger.Error("RerunTask", "index", dttr.task.index)
				dttr.needToRerun = false
				dttr.runAnteFn(dttr.task)
			} else if dttr.step == dttRoutineStepReadyForSerial ||
				dttr.step == dttRoutineStepSerial ||
				dttr.step == dttRoutineStepFinished {
				dttr.needToRerun = false
				//dttr.logger.Error("task is empty or finished")
			} else {
				//dttr.logger.Error("shouldRerunLater", "index", dttr.task.index)
				// maybe the task is in other condition, running concurrent execution or running to make new task.
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
		dttr.needToRerun = true
		dttr.rerunCh <- 0 // todo: maybe blocked for several milliseconds. why?
		//dttr.logger.Error("sendRerunCh", "index", dttr.task.index)
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
	case dttRoutineStepAnteStart:
		fallthrough
	case dttRoutineStepFinished:
		fallthrough
	case dttRoutineStepSerial:
		return false
	}
	return true
}

//-------------------------------------

type DTTManager struct {
	done           chan int8
	totalCount     int
	txs            [][]byte
	dttRoutineList []*dttRoutine // key: txIndex, value: dttRoutine
	serialIndex    int
	serialTask     *DeliverTxTask
	serialCh       chan int8

	//mtx sync.Mutex

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
	dttm.serialTask = nil
	dttm.serialIndex = -1
	dttm.serialCh = make(chan int8, 3)

	dttm.txResponses = make([]*abci.ResponseDeliverTx, len(txs))
	dttm.invalidTxs = 0

	dttm.checkStateCtx = dttm.app.checkState.ctx.WithBlockHeight(dttm.app.checkState.ctx.BlockHeight() + 1)

	go dttm.serialRoutine()

	for i := 0; i < maxDeliverTxsConcurrentNum; i++ {
		dttr := dttm.dttRoutineList[i]
		dttr.OnStart()
		dttr.makeNewTask(txs[i], i)
	}
}

func (dttm *DTTManager) concurrentBasic(txByte []byte, index int) *DeliverTxTask {
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
		return task
	}

	task.info.handler = dttm.app.getModeHandler(runTxModeDeliverPartConcurrent) //dm.handler
	task.fee, _, task.from, task.to = dttm.app.getTxFeeAndFromHandler(dttm.checkStateCtx, task.info.tx)

	if err = validateBasicTxMsgs(task.info.tx.GetMsgs()); err != nil {
		task.err = err
		dttm.app.logger.Error("validateBasicTxMsgs failed", "err", err)
	}
	return task
}

func (dttm *DTTManager) hasConflict(taskA *DeliverTxTask, taskB *DeliverTxTask) bool {
	if taskA.from == taskB.from {
		return true
	}
	if len(taskA.to) == 0 && len(taskB.to) == 0 {
		return false
	}
	if taskA.index < taskB.index && taskA.to == taskB.from {
		//dttm.app.logger.Error("hasConflict", "index", taskA.index, "conflict", taskB.index, "addr", taskA.to)
		return true
	} else if taskA.index > taskB.index && taskA.from == taskB.to {
		//dttm.app.logger.Error("hasConflict", "index", taskA.index, "conflict", taskB.index, "addr", taskB.to)
		return true
	}
	return false
}

func (dttm *DTTManager) runConcurrentAnte(task *DeliverTxTask) error {
	if dttm.app.anteHandler == nil {
		return fmt.Errorf("anteHandler cannot be nil")
	}

	curDttr := dttm.dttRoutineList[task.routineIndex]
	curDttr.step = dttRoutineStepAnteStart
	defer func() {
		curDttr.step = dttRoutineStepAnteFinished
	}()

	count := len(dttm.dttRoutineList)
	for i := 0; i < count; i++ {
		dttr := dttm.dttRoutineList[i]
		if dttr.task == nil ||
			dttr.txIndex != dttr.task.index ||
			dttr.step == dttRoutineStepNone || dttr.step == dttRoutineStepStart ||
			dttr.step == dttRoutineStepFinished || dttr.step == dttRoutineStepReadyForSerial ||
			!dttm.hasConflict(dttr.task, task) {
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
	if task.prevTaskIndex >= 0 || task.index <= dttm.serialIndex { //|| dttr.needToRerun {
		return nil
	}

	dttm.app.logger.Info("RunAnte", "index", task.index)
	task.info.ctx = dttm.app.getContextForTx(runTxModeDeliverPartConcurrent, task.info.txBytes) // same context for all txs in a block
	task.canRerun = 0

	task.info.ctx = task.info.ctx.WithCache(sdk.NewCache(dttm.app.blockCache, useCache(runTxModeDeliverPartConcurrent))) // one cache for a tx

	anteStart := time.Now()
	err := dttm.runAnte(task)
	totalAnteDuration += time.Since(anteStart).Microseconds()
	if err != nil {
		dttm.app.logger.Error("anteFailed", "index", task.index, "err", err)
	}
	task.err = err

	if task.canRerun > 0 {
		//curDttr.logger.Error("rerunChInFromAnte", "index", task.index)
		curDttr.shouldRerun(-1)
	} else if dttm.serialIndex+1 == task.index && !curDttr.needToRerun && task.prevTaskIndex < 0 && dttm.serialTask == nil {
		//dttm.app.logger.Info("ExtractNextSerialFromAnte", "index", curDttr.task.index, "step", curDttr.step, "needToRerun", curDttr.needToRerun)
		dttm.serialCh <- task.routineIndex
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
	keepAliveTicker := time.NewTicker(keepAliveIntervalMS * time.Millisecond)
	nextTaskRoutine := int8(-1)
	for {
		select {
		case routineIndex := <-dttm.serialCh:
			keepAliveTicker.Stop()
			rt := dttm.dttRoutineList[routineIndex]
			task := rt.task
			if task.index != dttm.serialIndex+1 {
				break
			}

			dttm.serialIndex = task.index
			dttm.serialTask = task
			rt.step = dttRoutineStepSerial
			dttm.serialExecution()
			rt.step = dttRoutineStepFinished
			dttm.serialTask = nil

			if dttm.serialIndex == dttm.totalCount-1 {
				//dttm.app.logger.Info("TotalTxFeeForCollector", "fee", dttm.currTxFee)
				count := len(dttm.dttRoutineList)
				for i := 0; i < count; i++ {
					dttr := dttm.dttRoutineList[i]
					dttr.stop()
				}

				dttm.done <- 0
				close(dttm.serialCh)
				return
			}

			// make new task for this routine
			nextIndex := maxDeliverTxsConcurrentNum + task.index
			if nextIndex < dttm.totalCount {
				rt.makeNewTask(dttm.txs[nextIndex], nextIndex)
			}

			// check whether there are ante-finished task
			count := len(dttm.dttRoutineList)
			nextTaskRoutine = -1
			var rerunRoutine *dttRoutine
			for i := 0; i < count; i++ {
				dttr := dttm.dttRoutineList[i]
				if dttr.task == nil || dttr.task.index <= task.index || dttr.step < dttRoutineStepAnteStart {
					continue
				}
				if dttr.task.prevTaskIndex == task.index || dttm.hasConflict(dttr.task, task) { //dttr.task.from == task.from {
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
						//dttm.app.logger.Info("ExtractNextSerialFromSerial", "index", dttr.task.index, "step", dttr.step, "needToRerun", dttr.needToRerun)
						dttm.serialCh <- nextTaskRoutine
					} else {
						//dttm.app.logger.Info("NotReadyForSerial", "index", dttr.task.index, "routine", nextTaskRoutine, "step", dttr.step, "needToRerun", dttr.needToRerun, "canRerun", dttr.task.canRerun, "prev", dttr.task.prevTaskIndex)
						keepAliveTicker.Reset(keepAliveIntervalMS * time.Microsecond)
					}
				}
			}

			if rerunRoutine != nil {
				//dttm.app.logger.Error("rerunRoutine", "index", rerunRoutine.task.index, "serial", task.index)
				rerunRoutine.shouldRerun(task.index)
			}
		case <-keepAliveTicker.C:
			//dttm.app.logger.Error("keepAliveTicker", "routine", nextTaskRoutine)
			if dttm.serialTask == nil && nextTaskRoutine >= 0 {
				//dttm.app.logger.Info("ExtractNextSerialFromTicker", "index", dttm.serialIndex, "routine", nextTaskRoutine)
				dttm.serialCh <- nextTaskRoutine
			}
			keepAliveTicker.Stop()
		}
	}
}

func (dttm *DTTManager) serialExecution() {
	dttm.app.logger.Info("SerialStart", "index", dttm.serialTask.index)
	info := dttm.serialTask.info
	handler := info.handler

	handleGasFn := func() {
		gasStart := time.Now()

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
	info.msCacheAnte.Write()
	info.ctx.Cache().Write(true)

	gasStart := time.Now()
	err := info.handler.handleGasConsumed(info)
	totalHandleGasTime += time.Since(gasStart).Microseconds()
	if err != nil {
		dttm.app.logger.Error("handleGasConsumed failed", "err", err)

		txRs := sdkerrors.ResponseDeliverTx(err, 0, 0, dttm.app.trace)
		execFinishedFn(txRs)
		return
	}

	dttm.app.UpdateFeeForCollector(dttm.serialTask.fee, true)

	// execute runMsgs
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

func (dttm *DTTManager) OnAccountUpdated(acc exported.Account, updateState bool) {
	start := time.Now()
	if updateState {
		addr := hex.EncodeToString(acc.GetAddress())
		dttm.accountUpdated(addr)
	}
	totalAccountUpdateDuration += time.Since(start).Microseconds()
}

func (dttm *DTTManager) accountUpdated(address string) {
	//dttm.mtx.Lock()
	//defer dttm.mtx.Unlock()

	num := len(dttm.dttRoutineList)
	for i := 0; i < num; i++ {
		dttr := dttm.dttRoutineList[i]
		if dttr.task == nil || dttr.txIndex != dttr.task.index || !dttr.needToRerunWhenContextChanged() || dttr.task.from != address {
			continue
		}

		dttr.shouldRerun(-1)
	}
}

//-------------------------------------------------------------

func (app *BaseApp) DeliverTxsConcurrent(txs [][]byte) []*abci.ResponseDeliverTx {
	if app.deliverTxsMgr == nil {
		app.deliverTxsMgr = NewDTTManager(app) //NewDeliverTxTasksManager(app)
	}

	//app.logger.Info("deliverTxs", "txs", len(txs))
	//start := time.Now()
	app.deliverTxsMgr.deliverTxs(txs)

	if len(txs) > 0 {
		//waiting for call back
		<-app.deliverTxsMgr.done
		close(app.deliverTxsMgr.done)
	}
	trace.GetElapsedInfo().AddInfo(trace.InvalidTxs, fmt.Sprintf("%d", app.deliverTxsMgr.invalidTxs))

	return app.deliverTxsMgr.txResponses
}

func (app *BaseApp) OnAccountUpdated(acc exported.Account, updateState bool) {
	if app.deliverTxsMgr != nil {
		app.deliverTxsMgr.OnAccountUpdated(acc, updateState)
	}
}

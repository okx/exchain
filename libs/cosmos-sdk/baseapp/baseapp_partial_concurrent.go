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

type (
	dttRoutineStep uint8
)

const (
	dttRoutineStepNone dttRoutineStep = iota
	dttRoutineStepStart
	dttRoutineStepWaitRerun
	dttRoutineStepAnteStart
	dttRoutineStepAnteFinished
	dttRoutineStepReadyForSerial
	dttRoutineStepSerial
	dttRoutineStepFinished

	maxDeliverTxsConcurrentNum = 4
	keepAliveIntervalMS        = 5
)

type DeliverTxTask struct {
	index        int
	canRerun     int8
	routineIndex int8

	info          *runTxInfo
	fee           sdk.Coins
	isEvm         bool
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
		txIndex:    -1,
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

func (dttr *dttRoutine) start() {
	dttr.done = make(chan int8)
	dttr.txByte = make(chan []byte, 3)
	dttr.rerunCh = make(chan int8, 5)
	dttr.step = dttRoutineStepNone
	go dttr.executeTaskRoutine()
}

func (dttr *dttRoutine) stop() {
	dttr.step = dttRoutineStepNone
	dttr.txIndex = -1
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
			if dttr.task.err == nil { //&& dttr.task.isEvm {
				dttr.runAnteFn(dttr.task)
			} else {
				dttr.step = dttRoutineStepReadyForSerial
			}
		case <-dttr.rerunCh:
			//dttr.logger.Error("readRerunCh", "index", dttr.task.index, "step", dttr.step)
			dttr.task.prevTaskIndex = -1
			if dttr.step == dttRoutineStepAnteFinished {
				dttr.needToRerun = false
				dttr.step = dttRoutineStepWaitRerun
				dttr.runAnteFn(dttr.task)
			} else if dttr.step == dttRoutineStepReadyForSerial ||
				dttr.step == dttRoutineStepSerial ||
				dttr.step == dttRoutineStepFinished {
				dttr.needToRerun = false
			} else {
				// maybe the task is in other condition, running concurrent execution or running to make new task.
				dttr.task.canRerun++
				dttr.needToRerun = false
			}
		}
	}
}

func (dttr *dttRoutine) shouldRerun(fromIndex int, fromAccountUpdate int) {
	if dttr.step == dttRoutineStepReadyForSerial || dttr.needToRerun == true || (dttr.task.prevTaskIndex >= 0 && dttr.task.prevTaskIndex > fromIndex) {
		return
	}
	if dttr.step == dttRoutineStepAnteStart || dttr.step == dttRoutineStepAnteFinished {
		if fromAccountUpdate >= 0 && dttr.task.prevTaskIndex < fromAccountUpdate {
			dttr.task.prevTaskIndex = fromAccountUpdate
		} else {
			dttr.needToRerun = true
			dttr.rerunCh <- 0
		}
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
	case dttRoutineStepWaitRerun:
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
	dttm.totalCount = len(txs)
	dttm.app.logger.Info("TotalTxs", "count", dttm.totalCount)
	dttm.txResponses = make([]*abci.ResponseDeliverTx, dttm.totalCount)
	dttm.invalidTxs = 0
	if dttm.totalCount == 0 {
		return
	}

	dttm.done = make(chan int8, 1)
	dttm.txs = txs
	dttm.serialTask = nil
	dttm.serialIndex = -1
	dttm.serialCh = make(chan int8, 4)

	dttm.checkStateCtx = dttm.app.checkState.ctx
	dttm.checkStateCtx.SetBlockHeight(dttm.app.checkState.ctx.BlockHeight() + 1)

	go dttm.serialRoutine()

	for i := 0; i < maxDeliverTxsConcurrentNum && i < dttm.totalCount; i++ {
		dttr := dttm.dttRoutineList[i]
		dttr.start()
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
	task.fee, task.isEvm, task.from, task.to = dttm.app.getTxFeeAndFromHandler(dttm.checkStateCtx, task.info.tx)

	if err = validateBasicTxMsgs(task.info.tx.GetMsgs()); err != nil {
		task.err = err
	}
	return task
}

func (dttm *DTTManager) hasConflict(taskA *DeliverTxTask, taskB *DeliverTxTask) bool {
	if len(taskA.from) > 0 && taskA.from == taskB.from {
		return true
	}
	if len(taskA.to) == 0 && len(taskB.to) == 0 {
		return false
	}
	if taskA.index < taskB.index && taskA.to == taskB.from {
		return true
	} else if taskA.index > taskB.index && taskA.from == taskB.to {
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
			dttr.task.index == task.index ||
			dttr.step == dttRoutineStepNone || //dttr.step == dttRoutineStepStart ||
			dttr.step == dttRoutineStepFinished || dttr.step == dttRoutineStepReadyForSerial ||
			(dttr.task.index > task.index && dttr.task.prevTaskIndex >= task.index) ||
			(dttr.task.index < task.index && task.prevTaskIndex >= dttr.task.index) {
			continue
		}
		conflict := dttm.hasConflict(dttr.task, task)
		if dttr.task.isEvm && task.isEvm && !conflict {
			continue
		}

		if !dttr.task.isEvm || conflict {
			if dttr.task.index < task.index && dttr.task.index > task.prevTaskIndex {
				task.prevTaskIndex = dttr.task.index
			}
		} else if !task.isEvm || conflict {
			if dttr.task.index > task.index && dttr.task.prevTaskIndex < task.index {
				dttr.task.prevTaskIndex = task.index
			}
		}
	}
	if task.prevTaskIndex < dttm.serialIndex || (task.prevTaskIndex == dttm.serialIndex && dttm.serialTask == nil) {
		task.prevTaskIndex = -1
	} else if task.index <= dttm.serialIndex || task.prevTaskIndex >= 0 {
		return nil
	}

	task.info.ctx = dttm.app.getContextForTx(runTxModeDeliverPartConcurrent, task.info.txBytes) // same context for all txs in a block
	task.canRerun = 0

	task.info.ctx.SetCache(sdk.NewCache(dttm.app.blockCache, useCache(runTxModeDeliverPartConcurrent))) // one cache for a tx

	err := dttm.runAnte(task)
	//if err != nil {
	//	dttm.app.logger.Error("anteFailed", "index", task.index, "err", err)
	//}
	task.err = err

	if task.canRerun > 0 {
		curDttr.shouldRerun(-1, -1)
	} else if dttm.serialIndex+1 == task.index && !curDttr.needToRerun && task.prevTaskIndex < 0 && dttm.serialTask == nil {
		dttm.serialCh <- task.routineIndex
	}

	return nil
}

func (dttm *DTTManager) runAnte(task *DeliverTxTask) error {
	info := task.info
	var anteCtx sdk.Context
	anteCtx, info.msCacheAnte = dttm.app.cacheTxContext(info.ctx, info.txBytes)
	anteCtx.SetEventManager(sdk.NewEventManager())
	//anteCtx = anteCtx.WithAnteTracer(dm.app.anteTracer)

	newCtx, err := dttm.app.anteHandler(anteCtx, info.tx, false) // NewAnteHandler

	ms := info.ctx.MultiStore()

	if !newCtx.IsZero() {
		info.ctx = newCtx
		info.ctx.SetMultiStore(ms)
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
			rt := dttm.dttRoutineList[routineIndex]
			task := rt.task
			if task.index != dttm.serialIndex+1 {
				break
			}
			keepAliveTicker.Stop()
			nextTaskRoutine = -1

			dttm.serialIndex = task.index
			dttm.serialTask = task
			rt.step = dttRoutineStepSerial
			dttm.serialExecution()
			rt.step = dttRoutineStepFinished
			dttm.serialTask = nil

			if dttm.serialIndex == dttm.totalCount-1 {
				count := len(dttm.dttRoutineList)
				for i := 0; i < count && i < dttm.totalCount; i++ {
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
			rerunRoutines := make([]*dttRoutine, 0)
			for i := 0; i < count; i++ {
				dttr := dttm.dttRoutineList[i]
				if dttr.task == nil || dttr.task.index <= task.index || dttr.step < dttRoutineStepAnteStart || dttr.task.prevTaskIndex > task.index {
					continue
				}

				if dttr.task.prevTaskIndex == task.index || !task.isEvm || dttm.hasConflict(dttr.task, task) { //dttr.task.from == task.from {
					if dttr.task.prevTaskIndex < task.index {
						dttr.task.prevTaskIndex = task.index
					}
					rerunRoutines = append(rerunRoutines, dttr)
				} else if dttr.task.index == dttm.serialIndex+1 {
					nextTaskRoutine = dttr.index
					if dttr.readyForSerialExecution() {
						dttm.serialCh <- nextTaskRoutine
					} else {
						keepAliveTicker.Reset(keepAliveIntervalMS * time.Microsecond)
					}
				}
			}

			for _, rerunRoutine := range rerunRoutines {
				rerunRoutine.shouldRerun(task.index, -1)
			}
		case <-keepAliveTicker.C:
			if dttm.serialTask == nil && nextTaskRoutine >= 0 && len(dttm.serialCh) == 0 {
				dttr := dttm.dttRoutineList[nextTaskRoutine]
				if dttr.task.index == dttm.serialIndex+1 && dttr.readyForSerialExecution() {
					keepAliveTicker.Stop()
					dttm.serialCh <- nextTaskRoutine
				}
			}
		}
	}
}

func (dttm *DTTManager) serialExecution() {
	info := dttm.serialTask.info
	handler := info.handler

	execFinishedFn := func(txRs abci.ResponseDeliverTx) {
		dttm.txResponses[dttm.serialTask.index] = &txRs
		if txRs.Code != abci.CodeTypeOK {
			dttm.invalidTxs++
		}
	}

	// execute anteHandler failed
	if dttm.serialTask.err != nil {
		txRs := sdkerrors.ResponseDeliverTx(dttm.serialTask.err, 0, 0, dttm.app.trace)
		execFinishedFn(txRs)
		return
	}

	info.msCacheAnte.Write()
	info.ctx.Cache().Write(true)

	err := info.handler.handleGasConsumed(info)
	if err != nil {
		txRs := sdkerrors.ResponseDeliverTx(err, 0, 0, dttm.app.trace)
		execFinishedFn(txRs)
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = dttm.app.runTx_defer_recover(r, info)
			info.msCache = nil
			info.result = nil
		}
		info.gInfo = sdk.GasInfo{GasWanted: info.gasWanted, GasUsed: info.ctx.GasMeter().GasConsumed()}

		var resp abci.ResponseDeliverTx
		if err != nil {
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
	}()

	defer handler.handleDeferGasConsumed(info)

	defer handler.handleDeferRefund(info)

	dttm.app.UpdateFeeForCollector(dttm.serialTask.fee, true)

	// execute runMsgs
	err = handler.handleRunMsg(info)
	//if err != nil {
	//	dttm.app.logger.Error("RunMsgFailed.", "err", err)
	//}
}

func (dttm *DTTManager) OnAccountUpdated(acc exported.Account, updateState bool) {
	if updateState {
		addr := hex.EncodeToString(acc.GetAddress())
		dttm.accountUpdated(addr)
	}
}

func (dttm *DTTManager) accountUpdated(address string) {
	num := len(dttm.dttRoutineList)
	serialIndex := dttm.serialIndex
	for i := 0; i < num; i++ {
		dttr := dttm.dttRoutineList[i]
		if dttr.task == nil || dttr.txIndex != dttr.task.index || !dttr.needToRerunWhenContextChanged() || dttr.task.from != address {
			continue
		}
		dttr.shouldRerun(-1, serialIndex)
	}
}

//-------------------------------------------------------------

func (app *BaseApp) DeliverTxsConcurrent(txs [][]byte) []*abci.ResponseDeliverTx {
	if app.deliverTxsMgr == nil {
		app.deliverTxsMgr = NewDTTManager(app)
	}

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

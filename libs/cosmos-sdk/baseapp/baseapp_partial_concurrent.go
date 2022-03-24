package baseapp

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"sync"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
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

//var blockHeight = int64(0)
//
//const AssignedBlockHeight = 5810742

type (
	partialConcurrentStep uint8
)

const (
	partialConcurrentStepStart partialConcurrentStep = iota
	partialConcurrentStepBasicFailed
	partialConcurrentStepBasicSucceed
	partialConcurrentStepAnteStart
	partialConcurrentStepAnteFailed
	partialConcurrentStepAnteSucceed
	partialConcurrentStepInRerun
	partialConcurrentStepInPending
	partialConcurrentStepSerialExecute
)

type DeliverTxTask struct {
	//tx            sdk.Tx
	index         int
	feeForCollect int64
	step          partialConcurrentStep
	updateCount   int8
	mtx           sync.Mutex

	info  *runTxInfo
	from  string //sdk.Address//exported.Account
	fee   sdk.Coins
	isEvm bool
	//signCache sdk.SigCache
	err error
	//anteErr        error
	//anteFailed bool
}

func newDeliverTxTask(tx sdk.Tx, index int) *DeliverTxTask {
	t := &DeliverTxTask{
		//tx:    tx,
		index: index,
		info:  &runTxInfo{tx: tx},
	}

	return t
}

func (dtt *DeliverTxTask) setStep(step partialConcurrentStep) {
	//dtt.mtx.Lock()
	//defer dtt.mtx.Unlock()

	if dtt.step == partialConcurrentStepInRerun && step != partialConcurrentStepStart {
		return
	}
	dtt.step = step
}

func (dtt *DeliverTxTask) getStep() partialConcurrentStep {
	//dtt.mtx.Lock()
	//defer dtt.mtx.Unlock()
	return dtt.step
}

func (dtt *DeliverTxTask) needToRerunWhenContextChanged() bool {
	step := dtt.getStep()
	switch step {
	case partialConcurrentStepBasicFailed:
		fallthrough
	case partialConcurrentStepInRerun:
		fallthrough
	case partialConcurrentStepSerialExecute:
		return false
	}
	return true
}

func (dtt *DeliverTxTask) setUpdateCount(count int8) {
	//dtt.mtx.Lock()
	//defer dtt.mtx.Unlock()

	dtt.updateCount = count
}

func (dtt *DeliverTxTask) getUpdateCount() int8 {
	//dtt.mtx.Lock()
	//defer dtt.mtx.Unlock()

	return dtt.updateCount
}

type NeedToRerunFn func(index int)

type sendersMap struct {
	mtx              sync.Mutex
	notFinishedTasks sync.Map // key: address, value: []*DeliverTxTask
	needRerunTasks   sync.Map //[]*DeliverTxTask
	logger           log.Logger
	rerunNotifyFn    NeedToRerunFn
}

func NewSendersMap() *sendersMap {
	sm := &sendersMap{
		notFinishedTasks: sync.Map{},
		needRerunTasks:   sync.Map{}, //[]*DeliverTxTask{},
		//rerunNotifyFn: rerunFn,
	}
	return sm
}

func (sm *sendersMap) setLogger(logger log.Logger) {
	sm.logger = logger
}

func (sm *sendersMap) Push(task *DeliverTxTask) (succeed bool) {
	if task == nil {
		return
	}

	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	address := task.from //.String()
	tasks, ok := sm.notFinishedTasks.Load(address)
	if !ok {
		succeed = true
		tasks = []*DeliverTxTask{task}
	} else {
		tasksList := tasks.([]*DeliverTxTask)
		count := len(tasksList)
		conflict := false
		for i := 0; i < count; i++ {
			if task.index > tasksList[i].index {
				sm.pushToRerunList(task)
				conflict = true
			} else if task.index < tasksList[i].index {
				if tasksList[i].getStep() != partialConcurrentStepInRerun {
					sm.pushToRerunList(tasksList[i])
				}
			} else {
				sm.logger.Error("Push into notFinishedTasks failed.", "index", task.index)
				return
			}
		}

		if !conflict {
			succeed = true
		}
		tasks = append(tasksList, task)
	}
	sm.notFinishedTasks.Store(address, tasks)
	//sm.logger.Info("PushTask", "index", task.index, "addr", address)

	return
}

func (sm *sendersMap) Pop(task *DeliverTxTask) {
	if task == nil {
		return
	}
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	address := task.from //.String()
	tasks, ok := sm.notFinishedTasks.Load(address)
	tasksList := tasks.([]*DeliverTxTask)
	if !ok || len(tasksList) == 0 {
		sm.logger.Error("address does not exist.")
		return
	}

	count := len(tasksList)
	pos := -1
	for i := 0; i < count; i++ {
		if tasksList[i].index == task.index {
			pos = i
			break
		}
	}

	tasksList = append(tasksList[:pos], tasksList[pos+1:]...)
	//if blockHeight == AssignedBlockHeight {
	//	sm.logger.Info("PopTask", "index", task.index, "addr", address)
	//}
	sm.notFinishedTasks.Store(address, tasksList)
	return
}

func (sm *sendersMap) pushToRerunList(task *DeliverTxTask) {
	//if blockHeight == AssignedBlockHeight {
	//}
	_, ok := sm.needRerunTasks.Load(task.index)
	if !ok {
		sm.logger.Error("MoveToRerun", "index", task.index)
		task.setStep(partialConcurrentStepInRerun)
		sm.needRerunTasks.Store(task.index, task)
	}
	//if task.getStep() == partialConcurrentStepInPending {
		sm.rerunNotifyFn(task.index)
	//}
}

func (sm *sendersMap) shouldRerun(task *DeliverTxTask) (rerun bool) {
	if task == nil {
		return
	}

	_, ok := sm.needRerunTasks.Load(task.index)
	if ok {
		rerun = true
	}
	return
}

func (sm *sendersMap) extractNextTask() (*DeliverTxTask, bool) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	minIndex := -1
	var task *DeliverTxTask
	existConflict := false
	sm.needRerunTasks.Range(func(key, value interface{}) bool {
		index := key.(int)
		if minIndex < 0 || index < minIndex {
			task = value.(*DeliverTxTask)
			// check whether exist previous tasks in notFinishedTasks
			address := task.from //.String()
			tmp, ok := sm.notFinishedTasks.Load(address)
			conflict := false
			if ok {
				notFinishedTasks := tmp.([]*DeliverTxTask)
				num := len(notFinishedTasks)
				for i := 0; i < num; i++ {
					if index > notFinishedTasks[i].index {
						sm.logger.Error("RerunTaskConflict", "target", index, "conflict", notFinishedTasks[i].index)
						conflict = true
						existConflict = true
						break
					}
				}
			}

			if !conflict {
				minIndex = index
			}
		}
		sm.logger.Error("NeedRerunTasks", "index", index)
		return true
	})

	if minIndex >= 0 {
		nextTask, ok := sm.needRerunTasks.Load(minIndex)
		if ok {
			//sm.logger.Info("extractNextTask", "index", minIndex)
			sm.needRerunTasks.Delete(minIndex)
			return nextTask.(*DeliverTxTask), existConflict
		}
	}
	return nil, existConflict
}

func (sm *sendersMap) accountUpdated(happened bool, times int8, address string, waitingIndex int) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	tasks, ok := sm.notFinishedTasks.Load(address)
	if !ok {
		return
	}

	tasksList := tasks.([]*DeliverTxTask)
	num := len(tasksList)
	for i := 0; i < num; i++ {
		task := tasksList[i]
		count := task.getUpdateCount()
		if happened {
			tasksList[i].setUpdateCount(count+times)
			// todo: whether should rerun the task
			if task.index != waitingIndex && task.updateCount > 0 && task.needToRerunWhenContextChanged() {
				sm.pushToRerunList(task)
			}
		} else {
			task.setUpdateCount(count-times)
		}
	}
}

func (sm *sendersMap) reset() {
	sm.notFinishedTasks = sync.Map{}
	sm.needRerunTasks = sync.Map{}
}

type DeliverTxTasksManager struct {
	done           chan int // done for all transactions are executed
	nextSignal     chan int // signal for taking a new tx into tasks
	statefulSignal chan int // signal for taking a new task from pendingTasks to statefulTask
	waitingCount   int
	mtx            sync.Mutex

	totalCount    int
	statefulIndex int
	pendingTasks  sync.Map
	statefulTask  *DeliverTxTask
	currTxFee     sdk.Coins
	finished bool

	sendersMap *sendersMap

	txResponses []*abci.ResponseDeliverTx

	app                *BaseApp
	//anteDuration       int64
	//gasAndMsgsDuration int64
	//serialDuration     int64
	////writeDuration      int64
	//deferGasTime       int64
	//handleGasTime      int64
	//runMsgsTime        int64
	//finishTime         int64
}

func NewDeliverTxTasksManager(app *BaseApp) *DeliverTxTasksManager {
	dm := &DeliverTxTasksManager{
		app:        app,
		sendersMap: NewSendersMap(),
	}
	dm.sendersMap.rerunNotifyFn = dm.removeFromPending

	return dm
}

func (dm *DeliverTxTasksManager) OnAccountUpdated(acc exported.Account) {
	addr := acc.GetAddress().String()
	//dm.app.logger.Info("OnAccountUpdated", "coins", acc.GetCoins(), "addr", addr)
	//if !acc.GetCoins().Empty() {
	waitingIndex := int(-1)
	if dm.statefulTask == nil {
		waitingIndex = dm.statefulIndex+1
	}
		dm.sendersMap.accountUpdated(true, 1, addr, waitingIndex)
	//}
}

func (dm *DeliverTxTasksManager) deliverTxs(txs [][]byte) {
	dm.done = make(chan int, 1)
	dm.nextSignal = make(chan int, 1)
	dm.statefulSignal = make(chan int, 1)
	dm.waitingCount = 0

	dm.totalCount = len(txs)
	dm.statefulIndex = -1
	dm.app.logger.Info("TotalTxs", "count", dm.totalCount)

	dm.pendingTasks = sync.Map{}
	dm.statefulTask = nil
	dm.currTxFee = sdk.Coins{}
	dm.finished = false

	dm.sendersMap.reset()
	dm.sendersMap.setLogger(dm.app.logger)

	dm.txResponses = make([]*abci.ResponseDeliverTx, len(txs))

	//dm.anteDuration = 0
	//dm.gasAndMsgsDuration = 0
	//dm.serialDuration = 0
	////dm.writeDuration = 0
	//dm.deferGasTime = 0
	//dm.handleGasTime = 0
	//dm.runMsgsTime = 0
	//dm.finishTime = 0
	//blockHeight = global.GetGlobalHeight()
	go dm.makeTasksRoutine(txs)
	go dm.runStatefulSerialRoutine()
}

func (dm *DeliverTxTasksManager) makeTasksRoutine(txs [][]byte) {
	taskIndex := 0
	for {
		// extract task from sendersMap
		nextTask, existConflict := dm.sendersMap.extractNextTask()
		if nextTask != nil {
			dm.makeNextTask(nil, nextTask.index, nextTask)
		} else if taskIndex < dm.totalCount {
			//dm.app.logger.Info("MakeNextTask", "index", taskIndex, "totalCount", dm.totalCount)
			dm.makeNextTask(txs[taskIndex], taskIndex, nil)
			taskIndex++
			dm.incrementWaitingCount(true)
		} else if existConflict {
			time.Sleep(1 * time.Millisecond)
		} else {
			dm.app.logger.Info("maxDeliverTxsConcurrentNum")
			break
		}
	}
	dm.finished = true
}

func (dm *DeliverTxTasksManager) makeNextTask(tx []byte, index int, task *DeliverTxTask) {
	//if blockHeight == AssignedBlockHeight {
	//dm.app.logger.Info("MakeNextTask", "task", task == nil, "index", index)
	//}
	go dm.runTxPartConcurrent(tx, index, task)
}

func (dm *DeliverTxTasksManager) runTxPartConcurrent(txByte []byte, index int, task *DeliverTxTask) {
	start := time.Now()
	mode := runTxModeDeliverPartConcurrent
	if task == nil {
		// create a new task
		task = dm.makeNewTask(txByte, index)
		//task.step = partialConcurrentStepBasic

		if task.err != nil {
			task.setStep(partialConcurrentStepBasicFailed)
			dm.pushIntoPending(task)
			return
		}

		//info := task.info
		task.info.handler = dm.app.getModeHandler(mode) //dm.handler

		// execute ante
		task.info.ctx = dm.app.getContextForTx(mode, task.info.txBytes) // same context for all txs in a block
		task.setUpdateCount(0)
		task.fee, task.isEvm, task.from = dm.app.getTxFeeAndFromHandler(task.info.ctx, task.info.tx)

		if err := validateBasicTxMsgs(task.info.tx.GetMsgs()); err != nil {
			task.err = err
			dm.app.logger.Error("validateBasicTxMsgs failed", "err", err)
			//dm.sendersMap.Pop(task)
			task.setStep(partialConcurrentStepBasicFailed)
			dm.pushIntoPending(task)
			return
		}

		task.setStep(partialConcurrentStepBasicSucceed)
		if !dm.sendersMap.Push(task) {
			//if blockHeight == AssignedBlockHeight {
			//dm.app.logger.Info("ExitConcurrent", "index", task.index)
			//}
			return
		}
	} else {
		dm.app.logger.Error("ResetContext", "index", task.index)

		task.setStep(partialConcurrentStepStart)
		task.info.ctx = dm.app.getContextForTx(mode, task.info.txBytes) // same context for all txs in a block
		task.setUpdateCount(0)
	}

	if dm.app.anteHandler != nil {
		task.setStep(partialConcurrentStepAnteStart)
		//if blockHeight == AssignedBlockHeight {
		//dm.app.logger.Info("RunAnte", "index", task.index, "addr", task.from)
		//}
		task.info.ctx = task.info.ctx.WithCache(sdk.NewCache(dm.app.blockCache, useCache(mode))) // one cache for a tx

		// todo: will change account. Account updated.
		dm.sendersMap.accountUpdated(false, 2, task.from, -1)
		err := dm.runAnte(task)
		if err != nil {
			dm.app.logger.Error("ante failed 1", "index", task.index, "err", err)
			//task.anteFailed = true
			task.setStep(partialConcurrentStepAnteFailed)
		} else {
			task.setStep(partialConcurrentStepAnteSucceed)
		}
		//dm.app.logger.Info("RunAnteSucceed 1", "index", task.index)
		if !dm.sendersMap.shouldRerun(task) {
			task.err = err

			dm.pushIntoPending(task)
			task.setStep(partialConcurrentStepInPending)
		} else {
			dm.app.logger.Error("NeedToReRunAnte", "index", task.index)
		}
	}
	totalAnteDuration += time.Since(start).Microseconds()
}

func (dm *DeliverTxTasksManager) makeNewTask(txByte []byte, index int) *DeliverTxTask {
	//dm.app.logger.Info("runTxPartConcurrent", "index", index)
	var realTx sdk.Tx
	var err error
	if mem := GetGlobalMempool(); mem != nil {
		realTx, _ = mem.ReapEssentialTx(txByte).(sdk.Tx)
	}
	if realTx == nil {
		realTx, err = dm.app.txDecoder(txByte)
	}
	//tx, err := dm.app.txDecoder(txByte)
	task := newDeliverTxTask(realTx, index)
	task.info.txBytes = txByte
	if err != nil {
		task.err = err
		dm.app.logger.Error("tx decode failed", "err", err)
	}

	//dm.tasks.Store(task.index, task)
	return task
}

// put task into pendingTasks after execution finished
func (dm *DeliverTxTasksManager) pushIntoPending(task *DeliverTxTask) {
	if task == nil {
		return
	}

	dm.mtx.Lock()
	defer dm.mtx.Unlock()

	//dm.app.logger.Info("NewIntoPendingTasks", "index", task.index, "curSerial", dm.statefulIndex+1, "task", dm.statefulTask != nil)
	//task.step = partialConcurrentStepSerialPrepare
	dm.pendingTasks.Store(task.index, task)
	if dm.statefulTask == nil && task.index == dm.statefulIndex+1 {
		dm.statefulSignal <- 0
	}
}

func (dm *DeliverTxTasksManager) removeFromPending(index int) {
	dm.mtx.Lock()
	defer dm.mtx.Unlock()

	task, ok := dm.pendingTasks.LoadAndDelete(index)
	if ok {
		//dm.app.logger.Error("RemoveFromPendingTasks", "index", index)
		if dm.finished {
			go dm.makeNextTask(nil, index, task.(*DeliverTxTask))
		}
	}
}

func (dm *DeliverTxTasksManager) runAnte(task *DeliverTxTask) error {
	info := task.info
	var anteCtx sdk.Context
	anteCtx, info.msCacheAnte = dm.app.cacheTxContext(info.ctx, info.txBytes) // info.msCacheAnte := ctx.MultiStore().CacheMultiStore(),  anteCtx := ctx.WithMultiStore(info.msCacheAnte)
	anteCtx = anteCtx.WithEventManager(sdk.NewEventManager())
	//anteCtx = anteCtx.WithAnteTracer(dm.app.anteTracer)

	newCtx, err := dm.app.anteHandler(anteCtx, info.tx, false) // NewAnteHandler

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

func (dm *DeliverTxTasksManager) runStatefulSerialRoutine() {
	begin := time.Now()
	finished := 0
	for {
		if finished == dm.totalCount {
			break
		}

		if !dm.extractStatefulTask() {
			start := time.Now()
			<-dm.statefulSignal
			elapsed := time.Since(start).Microseconds()
			dm.app.logger.Info("time to waiting for extractStatefulTask", "index", dm.statefulIndex+1, "us", elapsed)
			//dm.anteDuration -= elapsed
			totalWaitingTime += elapsed
			continue
		}

		//if blockHeight == AssignedBlockHeight {
		//dm.app.logger.Info("RunStatefulSerialRoutine", "index", dm.statefulTask.index)
		//}
		//start := time.Now()

		info := dm.statefulTask.info
		handler := info.handler

		handleGasFn := func() {
			gasStart := time.Now()

			//dm.updateFeeCollector()

			// todo: will change account. Account updated.
			//dm.app.logger.Info("handleDeferRefund", "index", dm.statefulTask.index, "addr", dm.statefulTask.from)
			dm.sendersMap.accountUpdated(false, 1, dm.statefulTask.from, -1)
			handler.handleDeferRefund(info)

			handler.handleDeferGasConsumed(info)

			if r := recover(); r != nil {
				_ = dm.app.runTx_defer_recover(r, info)
				info.msCache = nil //TODO msCache not write
				info.result = nil
			}
			info.gInfo = sdk.GasInfo{GasWanted: info.gasWanted, GasUsed: info.ctx.GasMeter().GasConsumed()}

			totalDeferGasTime += time.Since(gasStart).Microseconds()
		}

		execFinishedFn := func(txRs abci.ResponseDeliverTx) {
			//dm.app.logger.Info("SerialFinished", "index", dm.statefulTask.index)
			dm.txResponses[dm.statefulTask.index] = &txRs
			dm.resetStatefulTask()
			finished++

			//dm.gasAndMsgsDuration += time.Since(start).Microseconds()
		}

		// execute anteHandler failed
		//err := dm.statefulTask.err
		//if err == nil {
		//	err = dm.statefulTask.anteErr
		//}
		if dm.statefulTask.err != nil {
			dm.app.logger.Error("RunSerialFinished", "index", dm.statefulTask.index, "err", dm.statefulTask.err)
			txRs := sdkerrors.ResponseDeliverTx(dm.statefulTask.err, 0, 0, dm.app.trace) //execResult.GetResponse()
			execFinishedFn(txRs)
			continue
		}

		//dm.app.logger.Info("WriteAnteCache", "index", dm.statefulTask.index)
		info.msCacheAnte.Write()
		info.ctx.Cache().Write(true)
		dm.calculateFeeForCollector(dm.statefulTask.fee, true)

		gasStart := time.Now()
		err := info.handler.handleGasConsumed(info)
		//dm.handleGasTime += time.Since(gasStart).Microseconds()
		totalHandleGasTime += time.Since(gasStart).Microseconds()
		if err != nil {
			dm.app.logger.Error("handleGasConsumed failed", "err", err)

			txRs := sdkerrors.ResponseDeliverTx(err, 0, 0, dm.app.trace)
			execFinishedFn(txRs)
			continue
		}

		// execute runMsgs
		// todo: will change account. Account updated.
		//dm.app.logger.Info("handleRunMsg", "index", dm.statefulTask.index, "addr", dm.statefulTask.from)
		dm.sendersMap.accountUpdated(false, 2, dm.statefulTask.from, -1)
		runMsgStart := time.Now()
		err = handler.handleRunMsg(info)
		//runMsgE := time.Since(runMsgStart).Microseconds()
		//dm.runMsgsTime += runMsgE
		totalRunMsgsTime += time.Since(runMsgStart).Microseconds()

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
		dm.app.logger.Info("TotalTxFeeForCollector", "fee", dm.currTxFee)

		dm.done <- 0
		close(dm.statefulSignal)
		close(dm.nextSignal)
		//dm.serialDuration = time.Since(begin).Microseconds()
		totalSerialDuration += time.Since(begin).Microseconds()
	} else {
		dm.app.logger.Error("finished count is not equal to total count", "finished", finished, "total", dm.totalCount)
	}
}

func (dm *DeliverTxTasksManager) calculateFeeForCollector(fee sdk.Coins, add bool) {
	dm.mtx.Lock()
	defer dm.mtx.Unlock()

	if add {
		dm.currTxFee = dm.currTxFee.Add(fee...)
	} else {
		dm.currTxFee = dm.currTxFee.Sub(fee)
	}
	//dm.app.logger.Info("CalculateFeeForCollector", "fee", dm.currTxFee)
}

func (dm *DeliverTxTasksManager) updateFeeCollector() {
	//if blockHeight == AssignedBlockHeight {
	//	dm.app.logger.Info("updateFeeCollector", "now", dm.currTxFee)
	//}
	ctx, cache := dm.app.cacheTxContext(dm.app.getContextForTx(runTxModeDeliver, []byte{}), []byte{})
	if err := dm.app.updateFeeCollectorAccHandler(ctx, dm.currTxFee); err != nil {
		panic(err)
	}
	cache.Write()
}

func (dm *DeliverTxTasksManager) extractStatefulTask() bool {
	//dm.app.logger.Info("extractStatefulTask", "index", dm.statefulIndex + 1)
	task, ok := dm.pendingTasks.Load(dm.statefulIndex + 1)
	if ok {
		dm.statefulTask = task.(*DeliverTxTask)
		dm.statefulTask.setStep(partialConcurrentStepSerialExecute)
		dm.pendingTasks.Delete(dm.statefulTask.index)
	}
	return ok
}

func (dm *DeliverTxTasksManager) resetStatefulTask() {
	dm.sendersMap.Pop(dm.statefulTask)
	dm.statefulTask = nil
	dm.incrementWaitingCount(false)
}

func (dm *DeliverTxTasksManager) incrementWaitingCount(increment bool) {
	if increment {
		dm.mtx.Lock()
		dm.waitingCount++
		count := dm.waitingCount
		dm.mtx.Unlock()

		//if blockHeight == AssignedBlockHeight {
		//dm.app.logger.Info("incrementWaitingCount", "count", count)
		//}
		if count == maxDeliverTxsConcurrentNum {
			<-dm.nextSignal
			//dm.statefulSignalCount--
			//if dm.statefulSignalCount < 0 {
			//	dm.app.logger.Error("dm.statefulSignalCount < 0", "count", dm.statefulSignalCount)
			//}
			//} else {
			//	// sleep 10 millisecond in case of the first maxDeliverTxsConcurrentNum txs have the same sender
			//	time.Sleep(1 * time.Millisecond)
		}
	} else {
		dm.mtx.Lock()
		dm.statefulIndex++
		dm.waitingCount--
		count := dm.waitingCount
		dm.mtx.Unlock()

		//if blockHeight == AssignedBlockHeight {
		//dm.app.logger.Info("decrementWaitingCount", "count", count)
		//}
		if count == maxDeliverTxsConcurrentNum-1 {
			dm.nextSignal <- 0
			//dm.statefulSignalCount++
			//if dm.statefulSignalCount > 1 {
			//	dm.app.logger.Error("dm.statefulSignalCount > 1", "count", dm.statefulSignalCount)
			//}
		}
	}
}

//func (dm *DeliverTxTasksManager)OnAccountUpdated(acc exported.Account) {
//	account := acc.GetAddress()
//	dm.app.logger.Info("OnAccountUpdated", hex.EncodeToString(account))
//}

//-------------------------------------------------------------

func (app *BaseApp) DeliverTxsConcurrent(txs [][]byte) []*abci.ResponseDeliverTx {
	if app.deliverTxsMgr == nil {
		app.deliverTxsMgr = NewDeliverTxTasksManager(app)
		// set observer for account's update
		//if app.setAccountObserverFn != nil {
		//	//app.setAccountObserverFn(app.deliverTxsMgr)
		//}
	}

	//app.logger.Info("deliverTxs", "txs", len(txs))
	//start := time.Now()
	app.deliverTxsMgr.deliverTxs(txs)

	if len(txs) > 0 {
		//waiting for call back
		<-app.deliverTxsMgr.done
		close(app.deliverTxsMgr.done)

		//dur := time.Since(start).Microseconds()
		////totalAnteDuration += app.deliverTxsMgr.anteDuration
		////totalGasAndMsgsDuration += app.deliverTxsMgr.gasAndMsgsDuration
		////totalWriteTime += app.deliverTxsMgr.writeDuration
		////totalHandleGasTime += app.deliverTxsMgr.handleGasTime
		////totalDeferGasTime += app.deliverTxsMgr.deferGasTime
		////totalRunMsgsTime += app.deliverTxsMgr.runMsgsTime
		////totalSavedTime = totalSavedTime + (app.deliverTxsMgr.anteDuration + app.deliverTxsMgr.gasAndMsgsDuration - app.deliverTxsMgr.serialDuration)
		//app.logger.Info("all durations",
		//	"whole", dur,
		//	//"ante", app.deliverTxsMgr.anteDuration,
		//	//"serial", app.deliverTxsMgr.serialDuration,
		//	//"gasAndMsgs", app.deliverTxsMgr.gasAndMsgsDuration,
		//	//"handleGas", app.deliverTxsMgr.handleGasTime,
		//	//"write", app.deliverTxsMgr.writeDuration,
		//	//"runMsgs", app.deliverTxsMgr.runMsgsTime,
		//	//"deferGas", app.deliverTxsMgr.deferGasTime,
		//	"serialSum", app.deliverTxsMgr.handleGasTime+app.deliverTxsMgr.runMsgsTime+app.deliverTxsMgr.deferGasTime,
		//	"handleGasAll", totalHandleGasTime,
		//	"writeAll", totalWriteTime,
		//	"runMsgsAll", totalRunMsgsTime,
		//	"deferGasAll", totalDeferGasTime,
		//	"serialSumAll", totalHandleGasTime+totalWriteTime+totalRunMsgsTime+totalDeferGasTime,
		//	"anteAll", totalAnteDuration,
		//	"gasAndMsgsAll", totalGasAndMsgsDuration,
		//	"serialAll", totalSerialDuration,
		//	"waitingAll", totalWaitingTime,
		//	"rerunAnteAll", totalRerunAnteTime,
		//	//"totalSavedTime", totalSavedTime,
		//	//"saved", float64(app.deliverTxsMgr.anteDuration)/float64(dur))
		//)
	}

	return app.deliverTxsMgr.txResponses
}

func (app *BaseApp) OnAccountUpdated(acc exported.Account) {
	if app.deliverTxsMgr != nil {
		app.deliverTxsMgr.OnAccountUpdated(acc)
	}
}

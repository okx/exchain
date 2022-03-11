package baseapp

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/global"
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

type DeliverTxTask struct {
	tx            sdk.Tx
	index         int
	feeForCollect int64
	//anteFailed    bool

	info          *runTxInfo
	from          sdk.Address
	fee           sdk.Coins
	isEvm         bool
	//signCache     sdk.SigCache
	//evmIndex      uint32
	basicVerifyErr error
	anteErr	error
}

func newDeliverTxTask(tx sdk.Tx, index int) *DeliverTxTask {
	t := &DeliverTxTask{
		//tx:    tx,
		index: index,
		info:  &runTxInfo{tx: tx},
	}

	return t
}

type sendersMap struct {
	mtx sync.Mutex
	senders sync.Map
	nextTasks []*DeliverTxTask
	logger      log.Logger
}

func NewSendersMap() *sendersMap {
	sm := &sendersMap{
		senders: sync.Map{},
		//nextTasks: make([]*DeliverTxTask, maxDeliverTxsConcurrentNum),
	}
	return sm
}

func (sm *sendersMap) setLogger(logger log.Logger) {
	sm.logger = logger
}

func (sm *sendersMap) Push(address sdk.Address, task *DeliverTxTask) (alreadyExist bool) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	tasks, ok := sm.senders.Load(address.String())
	if ok {
		alreadyExist = true
		tasks = append(tasks.([]*DeliverTxTask), task)
		sm.logger.Info("Push", "index", task.index, "addr", address)
	} else {
		tasks = []*DeliverTxTask{task}
		sm.senders.Store(address.String(), tasks)
		//sm.logger.Info("New", "index", task.index, "addr", address)
	}

	return
}

func (sm *sendersMap) Pop(address sdk.Address, task *DeliverTxTask) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()

	tasks, ok := sm.senders.Load(address.String())
	tasksList := tasks.([]*DeliverTxTask)
	if ok {
		count := len(tasksList)
		for i:=0;i<count;i++ {
			if tasksList[i] == task {
				tasksList = append(tasksList[:i], tasksList[i+1:]...)
				fmt.Println("Pop", task.index, address.String())
				break
			}
		}
		if len(tasksList) > 0 {
			sm.senders.Store(tasksList, address.String())

			//
			minIndex := tasksList[0].index
			for i:=1;i<count;i++ {
				if tasksList[i].index < minIndex {
					minIndex = tasksList[i].index
				}
			}
			if len(sm.nextTasks) == 0 {
				sm.nextTasks = []*DeliverTxTask{tasksList[minIndex]}
			} else {
				sm.nextTasks = append(sm.nextTasks, tasksList[minIndex])
			}
		} else {
			sm.senders.Delete(address.String())
			fmt.Println("Delete", address.String())
		}
	} else {
		panic("address is not existed.")
	}
}

func (sm *sendersMap) extractNextTask() *DeliverTxTask {
	count := len(sm.nextTasks)
	if count == 0 {
		return nil
	}

	minIndex := sm.nextTasks[0].index
	for i:=1; i<count; i++ {
		if sm.nextTasks[i].index < minIndex {
			minIndex = sm.nextTasks[i].index
		}
	}
	nextTask := sm.nextTasks[minIndex]
	sm.nextTasks = append(sm.nextTasks[:minIndex], sm.nextTasks[minIndex+1:]...)
	return nextTask
}

type DeliverTxTasksManager struct {
	done                chan int // done for all transactions are executed
	nextSignal          chan int // signal for taking a new tx into tasks
	statefulSignal      chan int // signal for taking a new task from pendingTasks to statefulTask
	waitingCount        int
	statefulSignalCount int
	mtx                 sync.Mutex

	totalCount     int
	statefulIndex  int
	tasks          sync.Map
	pendingTasks   sync.Map
	statefulTask   *DeliverTxTask
	currTxFee      sdk.Coins

	sendersMap *sendersMap

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
	dm.statefulSignal = make(chan int, 1)
	dm.statefulSignalCount = 0
	dm.waitingCount = 0

	dm.totalCount = len(txs)
	dm.statefulIndex = -1

	dm.tasks = sync.Map{}
	dm.pendingTasks = sync.Map{}
	dm.statefulTask = nil
	dm.currTxFee = sdk.Coins{}

	dm.sendersMap = NewSendersMap()
	dm.sendersMap.setLogger(dm.app.logger)

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
	go dm.runStatefulSerialRoutine()
}

func (dm *DeliverTxTasksManager) makeTasksRoutine(txs [][]byte) {
	taskIndex := 0
	for {
		if taskIndex == dm.totalCount {
			break
		}

		//todo: extract task from sendersMap
		nextTask := dm.sendersMap.extractNextTask()
		if nextTask != nil {
			dm.makeNextTask(txs[taskIndex], taskIndex, nextTask)
		} else {
			dm.makeNextTask(txs[taskIndex], taskIndex, nil)
			taskIndex++
		}
		dm.incrementWaitingCount(true)
	}
}

func (dm *DeliverTxTasksManager) makeNextTask(tx []byte, index int, task *DeliverTxTask) {
	go dm.runTxPartConcurrent(tx, index, task)
}

func (dm *DeliverTxTasksManager) runTxPartConcurrent(txByte []byte, index int, task *DeliverTxTask) {
	start := time.Now()
	if task == nil {
		// create a new task
		task = dm.makeNewTask(txByte, index)

		if task.basicVerifyErr != nil {
			dm.pushIntoPending(task)
			return
		}

		mode := runTxModeDeliverPartConcurrent
		info := task.info
		info.handler = dm.app.getModeHandler(mode) //dm.handler

		// execute ante
		info.ctx = dm.app.getContextForTx(mode, info.txBytes) // same context for all txs in a block
		var signCache sdk.SigCache
		//task.fee, task.isEvm, signCache = dm.app.getTxFee(info.ctx, info.tx)
		task.tx, task.fee, task.isEvm, task.from, signCache = dm.app.evmTxFromHandler(info.ctx, info.tx)
		info.ctx = info.ctx.WithSigCache(signCache)
		info.ctx = info.ctx.WithCache(sdk.NewCache(dm.app.blockCache, useCache(mode))) // one cache for a tx

		if err := validateBasicTxMsgs(info.tx.GetMsgs()); err != nil {
			task.basicVerifyErr = err
			dm.app.logger.Error("validateBasicTxMsgs failed", "basicVerifyErr", err)
			dm.pushIntoPending(task)
			return
		}

		// todo: check duplicated sender
		if dm.sendersMap.Push(task.from, task) {
			// waiting util
			dm.incrementWaitingCount(false)
			return
		}
	}

	if dm.app.anteAuthHandler != nil {
		//dm.app.logger.Info("runAnte", "index", task.index)
		err := dm.runAnte(task) // dm.app.runAnte(task.info, mode)
		if err != nil {
			dm.app.logger.Error("ante failed 1", "basicVerifyErr", err)
			// todo: should make a judge for the basicVerifyErr. There are some errors don't need to re-run AnteHandler.
			task.anteErr = err
		}
		dm.calculateFeeForCollector(task.fee, true)
	}

	dm.pushIntoPending(task)

	elapsed := time.Since(start).Microseconds()
	dm.anteDuration += elapsed
}

func (dm *DeliverTxTasksManager) makeNewTask(txByte []byte, index int) *DeliverTxTask {
	//dm.app.logger.Info("runTxPartConcurrent", "index", index)
	tx, err := dm.app.txDecoder(txByte)
	task := newDeliverTxTask(tx, index)
	task.info.txBytes = txByte
	if err != nil {
		task.basicVerifyErr = err
		dm.app.logger.Error("tx decode failed", " basicVerifyErr", err)
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
	if dm.statefulTask == nil && task.index == dm.statefulIndex+1 {
		dm.statefulSignal <- 0
	}
	dm.tasks.Delete(task.index)
}

func (dm *DeliverTxTasksManager) runAnte(task *DeliverTxTask) error {
	info := task.info
	var anteCtx sdk.Context

	// Cache wrap context before AnteHandler call in case it aborts.
	// This is required for both CheckTx and DeliverTx.
	// Ref: https://github.com/cosmos/cosmos-sdk/issues/2772
	//
	// NOTE: Alternatively, we could require that AnteHandler ensures that
	// writes do not happen if aborted/failed.  This may have some
	// performance benefits, but it'll be more difficult to get right.
	anteCtx, info.msCacheAnte = dm.app.cacheTxContext(info.ctx, info.txBytes) // info.msCacheAnte := ctx.MultiStore().CacheMultiStore(),  anteCtx := ctx.WithMultiStore(info.msCacheAnte)
	anteCtx = anteCtx.WithEventManager(sdk.NewEventManager())
	//anteCtx = anteCtx.WithAnteTracer(dm.app.anteTracer)
	//if task.isEvm {
	//	info.msCacheAnte.IteratorCache(func(key, value []byte, isDirty bool) bool {
	//		if isDirty {
	//			fmt.Println(task.index, hex.EncodeToString(key), hex.EncodeToString(value))
	//		}
	//		return true
	//	})
	//}

	newCtx, err := dm.app.anteAuthHandler(anteCtx, info.tx, false) // NewAnteHandler

	//if task.isEvm {
	//	info.msCacheAnte.IteratorCache(func(key, value []byte, isDirty bool) bool {
	//		if isDirty {
	//			fmt.Println("after", task.index, hex.EncodeToString(key), hex.EncodeToString(value))
	//		}
	//		return true
	//	})
	//}

	ms := info.ctx.MultiStore()
	//info.accountNonce = newCtx.AccountNonce()

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
		//dm.updateEvmTxFrom(task)
	}
	// GasMeter expected to be set in AnteHandler
	info.gasWanted = info.ctx.GasMeter().Limit()
	if err != nil {
		return err
	}

	info.msCacheAnte.Write()
	info.ctx.Cache().Write(true)

	return nil
}

func (dm *DeliverTxTasksManager) runNonceAndSequenceAnte(task *DeliverTxTask) error {
	info := task.info
	var anteCtx sdk.Context

	// Cache wrap context before AnteHandler call in case it aborts.
	// This is required for both CheckTx and DeliverTx.
	// Ref: https://github.com/cosmos/cosmos-sdk/issues/2772
	//
	// NOTE: Alternatively, we could require that AnteHandler ensures that
	// writes do not happen if aborted/failed.  This may have some
	// performance benefits, but it'll be more difficult to get right.
	anteCtx, _ = dm.app.cacheTxContext(info.ctx, info.txBytes) // info.msCacheAnte := ctx.MultiStore().CacheMultiStore(),  anteCtx := ctx.WithMultiStore(info.msCacheAnte)
	anteCtx = anteCtx.WithEventManager(sdk.NewEventManager())

	var tx sdk.Tx
	if task.isEvm {
		tx = task.tx
	} else {
		tx = info.tx
	}
	newCtx, err := dm.app.nonceSequenceHandler(anteCtx, tx, false) // NewAnteHandler

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
	}
	// GasMeter expected to be set in AnteHandler
	//info.gasWanted = info.ctx.GasMeter().Limit()
	if err != nil {
		return err
	}

	//info.msCacheAnte.Write()
	//info.ctx.Cache().Write(true)

	return nil
}

//func (dm *DeliverTxTasksManager) updateEvmTxFrom(task *DeliverTxTask) {
//	if task.isEvm && dm.app.evmTxFromHandler != nil {
//		evmTx, ok := dm.app.evmTxFromHandler(task.info.ctx, task.info.tx)
//		if ok {
//			task.info.tx = evmTx
//		}
//	}
//}

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
			dm.app.logger.Info("time to waiting for extractStatefulTask", "index", dm.statefulIndex, "us", elapsed)
			dm.anteDuration -= elapsed
			totalWaitingTime += elapsed
			continue
		}

		//dm.app.logger.Info("runStatefulSerialRoutine", "index", dm.statefulTask.index)
		start := time.Now()

		info := dm.statefulTask.info
		handler := info.handler

		handleGasFn := func() {
			gasStart := time.Now()

			dm.updateFeeCollector()
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
			dm.txResponses[dm.statefulTask.index] = &txRs
			dm.resetStatefulTask()
			finished++

			dm.gasAndMsgsDuration += time.Since(start).Microseconds()
		}

		// execute anteHandler failed
		if dm.statefulTask.basicVerifyErr != nil {
			txRs := sdkerrors.ResponseDeliverTx(dm.statefulTask.basicVerifyErr, 0, 0, dm.app.trace) //execResult.GetResponse()
			execFinishedFn(txRs)
			continue
		}

		var tx sdk.Tx
		if dm.statefulTask.isEvm {
			tx = dm.statefulTask.tx
		} else {
			tx = dm.statefulTask.tx
		}
		if dm.app.incrementSeqHandler != nil {
			_, err := dm.app.incrementSeqHandler(info.ctx, tx)
			if err != nil {
				dm.app.logger.Error("incrementSeq failed.")
			}
		}
		if dm.statefulTask.anteErr != nil {
			var err error
			if strings.Contains(dm.statefulTask.anteErr.Error(), "invalid nonce") && dm.app.nonceVerifyHandler != nil {
				err = dm.app.nonceVerifyHandler(info.ctx, tx)
			} else {
				err = errors.New("other ante error")
			}

			if err != nil {
				dm.app.logger.Error("AnteFailed", "err", err)
				txRs := sdkerrors.ResponseDeliverTx(dm.statefulTask.anteErr, 0, 0, dm.app.trace) //execResult.GetResponse()
				handleGasFn()
				execFinishedFn(txRs)
				continue
			}
		}

		gasStart := time.Now()
		basicVerifyErr := info.handler.handleGasConsumed(info)
		dm.handleGasTime += time.Since(gasStart).Microseconds()
		if basicVerifyErr != nil {
			dm.app.logger.Error("handleGasConsumed failed", "basicVerifyErr", basicVerifyErr)

			txRs := sdkerrors.ResponseDeliverTx(basicVerifyErr, 0, 0, dm.app.trace)
			execFinishedFn(txRs)
			continue
		}

		//info.msCacheAnte.Write()
		//info.ctx.Cache().Write(true)

		// execute runMsgs
		runMsgStart := time.Now()
		err := handler.handleRunMsg(info)
		runMsgE := time.Since(runMsgStart).Microseconds()
		dm.runMsgsTime += runMsgE

		handleGasFn()

		var resp abci.ResponseDeliverTx
		if err != nil {
			//dm.app.logger.Error("handleRunMsg failed", "basicVerifyErr", basicVerifyErr)
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
		//// update fee collector
		//dm.updateFeeCollector()

		dm.done <- 0
		close(dm.statefulSignal)
		close(dm.nextSignal)
		dm.serialDuration = time.Since(begin).Microseconds()
		totalSerialDuration += dm.serialDuration
	} else {
		dm.app.logger.Error("finished count is not equal to total count", "finished", finished, "total", dm.totalCount)
	}
}

func (dm *DeliverTxTasksManager) calculateFeeForCollector(fee sdk.Coins, add bool) {
	if add {
		dm.currTxFee = dm.currTxFee.Add(fee...)
	} else {
		dm.currTxFee = dm.currTxFee.Sub(fee)
	}
}

func (dm *DeliverTxTasksManager) updateFeeCollector() {
	//dm.app.logger.Info("updateFeeCollector", "now", dm.currTxFee)
	ctx, cache := dm.app.cacheTxContext(dm.app.getContextForTx(runTxModeDeliver, []byte{}), []byte{})
	if err := dm.app.updateFeeCollectorAccHandler(ctx, dm.currTxFee); err != nil {
		panic(err)
	}
	cache.Write()
}

func (dm *DeliverTxTasksManager) extractStatefulTask() bool {
	task, ok := dm.pendingTasks.Load(dm.statefulIndex + 1)
	if ok {
		dm.statefulTask = task.(*DeliverTxTask)
		dm.pendingTasks.Delete(dm.statefulTask.index)

		dm.incrementWaitingCount(false)
	}
	return ok
}

func (dm *DeliverTxTasksManager) resetStatefulTask() {
	dm.statefulTask = nil
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
			dm.statefulSignalCount--
			if dm.statefulSignalCount < 0 {
				dm.app.logger.Error("dm.statefulSignalCount < 0", "count", dm.statefulSignalCount)
			}
		}
	} else {
		dm.mtx.Lock()
		dm.statefulIndex++
		dm.waitingCount--
		count := dm.waitingCount
		dm.mtx.Unlock()

		if count >= maxDeliverTxsConcurrentNum-1 {
			dm.nextSignal <- 0
			dm.statefulSignalCount++
			if dm.statefulSignalCount > 1 {
				dm.app.logger.Error("dm.statefulSignalCount > 1", "count", dm.statefulSignalCount)
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

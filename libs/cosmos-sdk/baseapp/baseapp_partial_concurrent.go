package baseapp

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/adb"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"sync"
)

const (
	maxDeliverTxsConcurrentNum = 50
	taskIntervalMS             = 1
)

type DeliverTxTask struct {
	tx            sdk.Tx
	index         int
	feeForCollect int64
	anteFailed    bool
	info          *runTxInfo
	from          adb.Address
}

func newDeliverTxTask(tx sdk.Tx, index int) *DeliverTxTask {
	t := &DeliverTxTask{
		tx:    tx,
		index: index,
	}

	return t
}

type DeliverTxTasksManager struct {
	done          chan int // done for all transactions are executed
	nextSignal    chan int // done for
	pendingSignal chan int // done for
	executeSignal chan int
	mtx           sync.Mutex

	totalCount    int
	curIndex      int
	//txs           [][]byte
	tasks         map[int]*DeliverTxTask
	pendingTasks  map[int]*DeliverTxTask
	executingTask *DeliverTxTask

	txExeResults []*executeResult
	txResponses  []*abci.ResponseDeliverTx

	txDecoder sdk.TxDecoder
}

func NewDeliverTxTasksManager(txDecoder sdk.TxDecoder) *DeliverTxTasksManager {
	return &DeliverTxTasksManager{txDecoder: txDecoder}
}

func (dm *DeliverTxTasksManager) deliverTxs(txs [][]byte) {
	dm.done = make(chan int, 1)
	dm.nextSignal = make(chan int, 1)
	dm.pendingSignal = make(chan int, 1)
	dm.executeSignal = make(chan int, 1)

	dm.totalCount = len(txs)
	dm.curIndex = -1

	dm.tasks = make(map[int]*DeliverTxTask, maxDeliverTxsConcurrentNum)
	dm.pendingTasks = make(map[int]*DeliverTxTask, maxDeliverTxsConcurrentNum)
	dm.txExeResults = make([]*executeResult, len(txs))
	dm.txResponses = make([]*abci.ResponseDeliverTx, len(txs))

	go dm.makeTasksRoutine(txs)
	go dm.runTxSerialRoutine()
}

func (dm *DeliverTxTasksManager) makeTasksRoutine(txs [][]byte) {
	taskIndex := 0
	for {
		if taskIndex == dm.totalCount-1 {
			break
		}

		numTasks, numPending := dm.getLen()
		switch {
		case numPending >= maxDeliverTxsConcurrentNum:
			fallthrough
		case numTasks >= maxDeliverTxsConcurrentNum:
			<-dm.nextSignal

		default:
			dm.makeNextTask(txs[taskIndex], taskIndex)
			taskIndex++
		}
	}
}

func (dm *DeliverTxTasksManager) makeNextTask(tx []byte, index int) {
	dm.mtx.Lock()
	defer dm.mtx.Unlock()

	go dm.runTxAnte(tx, index)
}

func (dm *DeliverTxTasksManager) runTxAnte(txByte []byte, index int) {
	// create a new task
	tx, err := dm.txDecoder(txByte)
	if err != nil {
		// TODO: panic??
		return
	}
	task := newDeliverTxTask(tx, index)

	dm.mtx.Lock()
	dm.tasks[task.index] = task
	dm.mtx.Unlock()

	// todo: execute ante

	// ante is executed
	// todo: take consider of this situation pendingTasks.Len == maxDeliverTxsConcurrentNum
	_, numPending := dm.getLen()
	if numPending == maxDeliverTxsConcurrentNum {
		// wait util
		<-dm.pendingSignal
	}

	dm.mtx.Lock()
	defer dm.mtx.Unlock()
	dm.pendingTasks[task.index] = task
	if dm.executingTask == nil && task.index == dm.curIndex+1 {
		dm.executeSignal <- 0
	}
	delete(dm.tasks, task.index)

	if len(dm.tasks) == maxDeliverTxsConcurrentNum-1 {
		dm.nextSignal <- 0
	}
}

func (dm *DeliverTxTasksManager) runTxSerialRoutine() {
	finished := 0
	for {
		numTasks, numPending := dm.getLen()
		if numTasks == 0 && numPending == 0 {
			break
		}

		if !dm.extractExecutingTask() {
			<-dm.executeSignal
			continue
		}

		// TODO: execute runMsgs etc.

		// one task is executed, execute next
		dm.resetExecutingTask()
		finished++
	}

	// all txs are executed
	if finished == dm.totalCount {
		dm.done <- 0
	}
}

func (dm *DeliverTxTasksManager) getLen() (int, int) {
	dm.mtx.Lock()
	defer dm.mtx.Unlock()
	return len(dm.tasks), len(dm.pendingTasks)
}

func (dm *DeliverTxTasksManager) extractExecutingTask() bool {
	dm.mtx.Lock()
	defer dm.mtx.Unlock()
	dm.executingTask = dm.pendingTasks[dm.curIndex+1]
	if dm.executingTask != nil {
		dm.curIndex++
		return true
	}
	return false
}

func (dm *DeliverTxTasksManager) resetExecutingTask() {
	dm.mtx.Lock()
	defer dm.mtx.Unlock()
	dm.executingTask = nil
}

//-------------------------------------------------------------

func (app *BaseApp) DeliverTxsConcurrent(txs [][]byte) []*abci.ResponseDeliverTx {
	if app.deliverTxsMgr == nil {
		app.deliverTxsMgr = NewDeliverTxTasksManager(app.txDecoder)
	}

	app.deliverTxsMgr.deliverTxs(txs)

	if len(txs) > 0 {
		//waiting for call back
		<-app.deliverTxsMgr.done
	}

	return app.deliverTxsMgr.txResponses
}

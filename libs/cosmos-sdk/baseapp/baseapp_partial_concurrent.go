package baseapp

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/trace"
)

const (
	maxDeliverTxsConcurrentNum = 4
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
	updateCount  int8
	needToRerun  bool
	canRerun     int8
	routineIndex int8

	info          *runTxInfo
	from          string //sdk.Address//exported.Account
	to            string
	err           error
	prevTaskIndex int // true: if there exists a not finished tx which has the same sender but smaller index
}

func newDeliverTxTask(tx sdk.Tx, index int) *DeliverTxTask {
	t := &DeliverTxTask{
		index:         index,
		info:          &runTxInfo{tx: tx},
		prevTaskIndex: -1,
	}

	return t
}

func (dtt *DeliverTxTask) setUpdateCount(count int8, add bool) bool {
	if add {
		dtt.updateCount += count
	} else {
		dtt.updateCount -= count
	}
	return dtt.updateCount > 0
}

func (dtt *DeliverTxTask) resetUpdateCount() {
	dtt.updateCount = 0
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

package baseapp

import (
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	gorid "github.com/okex/exchain/libs/goroutine"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"sync"
)

type task interface {
	part1()
	part2()
	id() int
	result() *abci.ResponseDeliverTx
}

type taskImp struct {
	block int64
	idx int
	abciCtx abci.DeliverTxContext

	wg *sync.WaitGroup

	app *BaseApp
	txBytes []byte
	res abci.ResponseDeliverTx

	info *runTxInfo
	
	logger  log.Logger
}

func newTask(id int, txBytes []byte, abciCtx abci.DeliverTxContext, wg *sync.WaitGroup, app *BaseApp) *taskImp {
	t := &taskImp{
		block: app.LastBlockHeight()+1,
		idx: id,
		wg: wg,
		txBytes: txBytes,
		app: app,
		logger: app.logger,
		abciCtx: abciCtx,
	}
	return t
}

func (t *taskImp) setResult(res abci.ResponseDeliverTx) {
	if t.info.finished {
		return
	}

	t.res = res
	t.info.finished = true
}

func (t *taskImp) id() int {
	return t.idx
}

func (t *taskImp) part1() {
	if t.abciCtx != nil && t.abciCtx.Stopped() {
		return
	}

	t.info = &runTxInfo{}
	t.info.txBytes = t.txBytes

	t.logger.Info("Deliver tx part1", "gid", gorid.GoRId, "block", t.block, "txid", t.idx)
	app := t.app
	tx, err := app.txDecoder(t.txBytes)
	if err != nil {
		t.res = sdkerrors.ResponseDeliverTx(err, 0, 0, app.trace)
		return
	}
	t.info.decoded = true

	t.info.tx = tx

	err = app.runtx6_1(t.info, runTxModeDeliver, LatestSimulateTxHeight)
	if err != nil {
		t.logger.Info("Deliver tx part1",
			"gid", gorid.GoRId,
			"block", t.block,
			"txid", t.idx,
			"err", err,
		)
		res := sdkerrors.ResponseDeliverTx(err, t.info.gInfo.GasWanted, t.info.gInfo.GasUsed, app.trace)
		t.setResult(res)
	}
}

func (t *taskImp) part2() {
	t.logger.Info("Deliver tx part2", "gid", gorid.GoRId, "block", t.block, "txid", t.idx)
	defer t.wg.Done()
	if !t.info.decoded {
		return
	}

	err := t.app.runtx6_2(t.info)
	if err != nil {
		t.logger.Info("Deliver tx part2",
			"gid", gorid.GoRId,
			"block", t.block,
			"txid", t.idx,
			"err", err,
			)

		res := sdkerrors.ResponseDeliverTx(err, t.info.gInfo.GasWanted, t.info.gInfo.GasUsed, t.app.trace)
		t.setResult(res)
		return
	}
	t.logger.Info("Deliver tx part2", "info", t.info, )

	if !t.info.finished {

		res := abci.ResponseDeliverTx{
			GasWanted: int64(t.info.gInfo.GasWanted), // TODO: Should type accept unsigned ints?
			GasUsed:   int64(t.info.gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
			Log:       t.info.result1.Log,
			Data:      t.info.result1.Data,
			Events:    t.info.result1.Events.ToABCIEvents(),
		}
		t.setResult(res)
	}
}

func (t *taskImp) result() *abci.ResponseDeliverTx {
	return &t.res
}
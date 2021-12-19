package baseapp

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
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
	finished bool

	info *runTxInfo

	ctx *sdk.Context
	runMsgCtx *sdk.Context
	mode runTxMode
	msCache sdk.CacheMultiStore
	msCacheAnte sdk.CacheMultiStore
	tx sdk.Tx
	accountNonce uint64
	msgs []sdk.Msg
	gasWanted uint64
	startingGas uint64
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

func (t *taskImp) id() int {
	return t.idx
}

func (t *taskImp) part1() {
	if t.abciCtx != nil && t.abciCtx.Stopped() {
		return
	}

	t.logger.Info("Deliver tx part1", "gid", gorid.GoRId, "block", t.block, "txid", t.idx)
	app := t.app
	tx, err := app.txDecoder(t.txBytes)
	if err != nil {
		t.res = sdkerrors.ResponseDeliverTx(err, 0, 0, app.trace)
		t.finished = true
		return
	}

	t.info, err = app.runtx6_1(runTxModeDeliver, t.txBytes, tx, LatestSimulateTxHeight)
	if err != nil {
		t.logger.Info("Deliver tx part1",
			"gid", gorid.GoRId,
			"block", t.block,
			"txid", t.idx,
			"err", err,
		)
		t.res = sdkerrors.ResponseDeliverTx(err, t.info.gInfo.GasWanted, t.info.gInfo.GasUsed, app.trace)
		t.finished = true
	}
}

func (t *taskImp) part2() {
	t.logger.Info("Deliver tx part2", "gid", gorid.GoRId, "block", t.block, "txid", t.idx)
	defer t.wg.Done()


	if t.finished {
		return
	}

	app := t.app
	err := app.runtx6_2(t.info)
	if err != nil {
		t.logger.Info("Deliver tx part2",
			"gid", gorid.GoRId,
			"block", t.block,
			"txid", t.idx,
			"err", err,
			)

		t.res = sdkerrors.ResponseDeliverTx(err, t.info.gInfo.GasWanted, t.info.gInfo.GasUsed, app.trace)
		t.finished = true
		return
	}

	t.res = abci.ResponseDeliverTx{
		GasWanted: int64(t.info.gInfo.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(t.info.gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
		Log:       t.info.result.Log,
		Data:      t.info.result.Data,
		Events:    t.info.result.Events.ToABCIEvents(),
	}
}

func (t *taskImp) result() *abci.ResponseDeliverTx {
	return &t.res
}
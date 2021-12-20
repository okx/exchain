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

	info *runTxInfo
	logger  log.Logger
}

type runTxInfo struct {
	handler modeHandler
	gasWanted uint64
	ctx sdk.Context
	runMsgCtx sdk.Context
	msCache sdk.CacheMultiStore
	msCacheAnte sdk.CacheMultiStore
	accountNonce uint64
	runMsgFinished bool
	startingGas uint64
	gInfo sdk.GasInfo

	result *sdk.Result
	txBytes []byte
	tx sdk.Tx
	finished bool
	decoded bool
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

func  (t *taskImp) genResult()  {
	if t.info.finished {
		return
	}
	t.res = abci.ResponseDeliverTx{
		GasWanted: int64(t.info.gInfo.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(t.info.gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
		Log:       t.info.result.Log,
		Data:      t.info.result.Data,
		Events:    t.info.result.Events.ToABCIEvents(),
	}
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

	err = app.runtx_part1(t.info, runTxModeDeliver, LatestSimulateTxHeight)
	if err != nil {
		t.logger.Info("Deliver tx part1", "gid", gorid.GoRId, "block", t.block, "txid", t.idx, "err", err,)
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

	err := t.app.runtx_part2(t.info)
	if err != nil {
		t.logger.Info("Deliver tx part2", "gid", gorid.GoRId, "block", t.block, "txid", t.idx, "err", err, )

		res := sdkerrors.ResponseDeliverTx(err, t.info.gInfo.GasWanted, t.info.gInfo.GasUsed, t.app.trace)
		t.setResult(res)
		return
	}
	//t.logger.Info("Deliver tx part2", "info", t.info, )
	t.genResult()
}

func (t *taskImp) result() *abci.ResponseDeliverTx {
	return &t.res
}
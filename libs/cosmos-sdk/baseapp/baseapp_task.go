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
	wg *sync.WaitGroup

	app *BaseApp
	txBytes []byte
	res abci.ResponseDeliverTx
	finished bool

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

func newTask(block int64, id int, txBytes []byte, wg *sync.WaitGroup, app *BaseApp) *taskImp {
	t := &taskImp{
		block: block,
		idx: id,
		wg: wg,
		txBytes: txBytes,
		app: app,
		logger: app.logger,
	}
	return t
}

func (t *taskImp) id() int {
	return t.idx
}

func (t *taskImp) part1() {

	t.logger.Info("Deliver tx part1", "gid", gorid.GoRId, "block", t.block, "txid", t.idx)
	app := t.app
	tx, err := app.txDecoder(t.txBytes)
	if err != nil {
		t.res = sdkerrors.ResponseDeliverTx(err, 0, 0, app.trace)
		t.finished = true
		return
	}

	var (
		gInfo  sdk.GasInfo
		result *sdk.Result
	)

	gInfo, result, _, err, t.finished = app.runTxPart1(runTxModeDeliver, t.txBytes, tx, LatestSimulateTxHeight, t) // DeliverTxAsync2
	if err != nil {
		t.res = sdkerrors.ResponseDeliverTx(err, gInfo.GasWanted, gInfo.GasUsed, app.trace)
		t.finished = true
		return
	}

	_ = result

}

func (t *taskImp) part2() {
	t.logger.Info("Deliver tx part2", "gid", gorid.GoRId, "block", t.block, "txid", t.idx)

	defer t.wg.Done()
	if t.finished {
		return
	}

	app := t.app

	gInfo, result, _, err := app.runTxPart2(t)
	if err != nil {
		t.res = sdkerrors.ResponseDeliverTx(err, gInfo.GasWanted, gInfo.GasUsed, app.trace)
		t.finished = true
		return
	}

	t.res = abci.ResponseDeliverTx{
		GasWanted: int64(gInfo.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
		Log:       result.Log,
		Data:      result.Data,
		Events:    result.Events.ToABCIEvents(),
	}
}

func (t *taskImp) result() *abci.ResponseDeliverTx {
	return &t.res
}
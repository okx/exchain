package baseapp

import (
	"fmt"
	"sync"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"runtime/debug"
)


func (app *BaseApp) runTx(mode runTxMode,
	txBytes []byte, tx sdk.Tx, height int64) (gInfo sdk.GasInfo,
	result *sdk.Result, msCacheList sdk.CacheMultiStore, err error) {

	if abci.RunTxByRefactor1 {
		var info *runTxInfo
		info, err = app.runtx_refactor(mode, txBytes, tx, height)
		return info.gInfo, info.result, info.msCacheAnte, err
	} else {
		return app.runtx_org(mode, txBytes, tx, height)
	}
}

func (app *BaseApp) runtx_part1(info *runTxInfo, mode runTxMode, height int64) (err error) {

	mhandler := app.getModeHandler(mode)
	info.handler = mhandler

	err = mhandler.handleStartHeight(info, height)
	if err != nil {
		return err
	}

	err = mhandler.handleGasConsumed(info)
	if err != nil {
		return err
	}

	if err := validateBasicTxMsgs(info.tx.GetMsgs()); err != nil {
		return err
	}

	if app.anteHandler != nil {
		err = app.runAnte(info, mode)
	}

	return err
}


func (app *BaseApp) runtx_part2(info *runTxInfo) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = app.runTx_defer_recover(r, info)
			info.msCache = nil //TODO msCache not write
			info.result = nil
			//app.logger.Info("info.result = nil")
		}
		info.gInfo = sdk.GasInfo{GasWanted: info.gasWanted, GasUsed: info.ctx.GasMeter().GasConsumed()}
	}()

	defer app.runTx_defer_consumegas(info, info.handler.getMode())
	defer app.runTx_defer_refund(info, info.handler.getMode())

	if info.finished {
		return
	}

	err = info.handler.handleRunMsg(info)
	if err == nil && info.result == nil {
		panic("")
	}
	return
}


func (app *BaseApp) runtx_refactor(mode runTxMode, txBytes []byte, tx sdk.Tx, height int64) (info *runTxInfo, err error) {
	info = &runTxInfo{}
	info.handler = app.getModeHandler(mode)
	info.tx = tx
	info.txBytes = txBytes
	mhandler := info.handler

	fmt.Printf("runtx_refactor\n")
	err = mhandler.handleStartHeight(info, height)
	if err != nil {
		return info, err
	}

	err = mhandler.handleGasConsumed(info)
	if err != nil {
		return info, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = app.runTx_defer_recover(r, info)
			info.msCache = nil //TODO msCache not write
			info.result = nil
		}
		info.gInfo = sdk.GasInfo{GasWanted: info.gasWanted, GasUsed: info.ctx.GasMeter().GasConsumed()}
	}()

	defer mhandler.handleDeferGasConsumed(info)

	defer mhandler.handleDeferRefund(info)

	//defer app.runTx_defer_consumegas(info, mode)
	//defer app.runTx_defer_refund(info, mode)

	if err := validateBasicTxMsgs(info.tx.GetMsgs()); err != nil {
		return info, err
	}

	if app.anteHandler != nil {
		err = app.runAnte(info, mode)
		if err != nil {
			return info, err
		}
	}

	err = mhandler.handleRunMsg(info)
	return info, err
}

func (app *BaseApp) dumpResp(res *abci.ResponseDeliverTx, idx int)  {

	app.logger.Info("===========DeliverTx===========",
		"block", app.LastBlockHeight()+1,
		"idx", idx,
		"Data len", len(res.Data),
		"Info", res.Info,
		"GasUsed", res.GasUsed,
		"GasWanted", res.GasWanted,
		"Code", res.Code,
	)

	for i, e := range res.Events {
		app.logger.Info("	Event", "id", i, "type", e.Type)

		if len(e.Attributes) == 0 {
			panic("")
		}
		for j, a := range e.Attributes {
			app.logger.Info("		Attributes", "id", j, "k", string(a.Key), "v", string(a.Value))
		}
	}
}
func (app *BaseApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	fmt.Printf("(app *BaseApp) DeliverTx\n")

	if abci.RunTxByRefactor2 {
		res = app.DeliverTx2Part(req)
	} else {
		res = app.DeliverTxOrg(req)
	}

	app.dumpResp(&res, 0)

	return res
}

func (app *BaseApp) DeliverTx2Part(req abci.RequestDeliverTx) abci.ResponseDeliverTx {

	var wg sync.WaitGroup
	wg.Add(1)
	task := newTask(0, req.Tx, nil, &wg, app)
	task.part1()
	task.part2()
	wg.Wait()
	return *task.result()
}

func (app *BaseApp) DeliverTxOrg(req abci.RequestDeliverTx) abci.ResponseDeliverTx {

	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		return sdkerrors.ResponseDeliverTx(err, 0, 0, app.trace)
	}

	gInfo, result, _, err := app.runTx(runTxModeDeliver, req.Tx, tx, LatestSimulateTxHeight) // DeliverTxConcurrently
	if err != nil {
		return sdkerrors.ResponseDeliverTx(err, gInfo.GasWanted, gInfo.GasUsed, app.trace)
	}

	return abci.ResponseDeliverTx{
		GasWanted: int64(gInfo.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
		Log:       result.Log,
		Data:      result.Data,
		Events:    result.Events.ToABCIEvents(),
	}
}

func (app *BaseApp) DeliverTxConcurrently(txList [][]byte, ctx abci.DeliverTxContext) []*abci.ResponseDeliverTx {
	fmt.Printf("DeliverTxConcurrently\n")
	var wg sync.WaitGroup
	wg.Add(len(txList))
	var taskList []task
	for i, tx := range txList {
		taskList = append(taskList, newTask(i, tx, ctx, &wg, app))
	}

	app.scheduler.start(taskList)
	wg.Wait()

	var results []*abci.ResponseDeliverTx
	for idx, task := range taskList {
		res := task.result()
		app.dumpResp(res, idx)

		if res.Code == 111222 {
			panic("111222")
		}
		results = append(results, res)
	}

	return results
}

// runTx processes a transaction within a given execution mode, encoded transaction
// bytes, and the decoded transaction itself. All state transitions occur through
// a cached Context depending on the mode provided. State only gets persisted
// if all messages get executed successfully and the execution mode is DeliverTx.
// Note, gas execution info is always returned. A reference to a Result is
// returned if the tx does not run out of gas and if all the messages are valid
// and execute successfully. An error is returned otherwise.
func (app *BaseApp) runTx_defer_recover(r interface{}, info *runTxInfo) error {
	var err error
	switch rType := r.(type) {
	// TODO: Use ErrOutOfGas instead of ErrorOutOfGas which would allow us
	// to keep the stracktrace.
	case sdk.ErrorOutOfGas:
		err = sdkerrors.Wrap(
			sdkerrors.ErrOutOfGas, fmt.Sprintf(
				"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
				rType.Descriptor, info.gasWanted, info.ctx.GasMeter().GasConsumed(),
			),
		)

	default:
		err = sdkerrors.Wrap(
			sdkerrors.ErrPanic, fmt.Sprintf(
				"recovered: %v\nstack:\n%v", r, string(debug.Stack()),
			),
		)
	}
	return err
}

func (app *BaseApp) runTx_defer_consumegas(info *runTxInfo, mode runTxMode) {
	app.pin(ConsumeGas, true, mode)
	defer app.pin(ConsumeGas, false, mode)
	if mode == runTxModeDeliver || (mode == runTxModeDeliverInAsync && app.parallelTxManage.isReRun(string(info.txBytes))) {
		info.ctx.BlockGasMeter().ConsumeGas(info.ctx.GasMeter().GasConsumedToLimit(), "block gas meter",)

		if info.ctx.BlockGasMeter().GasConsumed() < info.startingGas {
			panic(sdk.ErrorGasOverflow{Descriptor: "tx gas summation"})
		}
	}
}


func (app *BaseApp) runTx_defer_refund(info *runTxInfo, mode runTxMode){

	if (mode == runTxModeDeliver || mode == runTxModeDeliverInAsync) && app.GasRefundHandler != nil {
		var gasRefundCtx sdk.Context
		if mode == runTxModeDeliver {
			gasRefundCtx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)
		} else if mode == runTxModeDeliverInAsync {
			gasRefundCtx = info.runMsgCtx
			if info.msCache == nil || !info.runMsgFinished { // case: panic when runMsg
				info.msCache = info.msCacheAnte.CacheMultiStore()
				gasRefundCtx = info.ctx.WithMultiStore(info.msCache)
			}
		}
		refundGas, err := app.GasRefundHandler(gasRefundCtx, info.tx)
		if err != nil {
			panic(err)
		}
		info.msCache.Write()
		if mode == runTxModeDeliverInAsync {
			app.parallelTxManage.setRefundFee(string(info.txBytes), refundGas)
		}
	}
}

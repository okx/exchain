package baseapp

import (
	"fmt"
	"runtime/debug"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

type runTxInfo struct {
	handler        modeHandler
	gasWanted      uint64
	ctx            sdk.Context
	runMsgCtx      sdk.Context
	msCache        sdk.CacheMultiStore
	msCacheAnte    sdk.CacheMultiStore
	accountNonce   uint64
	runMsgFinished bool
	startingGas    uint64
	gInfo          sdk.GasInfo

	result  *sdk.Result
	txBytes []byte
	tx      sdk.Tx
}

func (app *BaseApp) runTx(mode runTxMode,
	txBytes []byte, tx sdk.Tx, height int64, from ...string) (gInfo sdk.GasInfo, result *sdk.Result,
	msCacheList sdk.CacheMultiStore, err error) {

	var info *runTxInfo
	info, err = app.runtx(mode, txBytes, tx, height, from...)
	return info.gInfo, info.result, info.msCacheAnte, err
}

func (app *BaseApp) runtx(mode runTxMode, txBytes []byte, tx sdk.Tx, height int64, from ...string) (info *runTxInfo, err error) {
	info = &runTxInfo{}
	err = app.runtxWithInfo(info, mode, txBytes, tx, height, from...)
	return
}
func (app *BaseApp) runtxWithInfo(info *runTxInfo, mode runTxMode, txBytes []byte, tx sdk.Tx, height int64, from ...string) (err error) {
	info.handler = app.getModeHandler(mode)
	info.tx = tx
	info.txBytes = txBytes
	handler := info.handler
	app.pin(ValTxMsgs, true, mode)

	//init info context
	err = handler.handleStartHeight(info, height)
	if err != nil {
		return err
	}
	//info with cache saved in app to load predesessor tx state
	if mode != runTxModeTrace {
		//in trace mode,  info ctx cache was already set to traceBlockCache instead of app.blockCache in app.tracetx()
		//to prevent modifying the deliver state
		//traceBlockCache was created with different root(chainCache) with app.blockCache in app.BeginBlockForTrace()
		info.ctx = info.ctx.WithCache(sdk.NewCache(app.blockCache, useCache(mode)))
	}
	for _, addr := range from {
		// cache from if exist
		if addr != "" {
			info.ctx = info.ctx.WithFrom(addr)
			break
		}
	}

	err = handler.handleGasConsumed(info)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			err = app.runTx_defer_recover(r, info)
			info.msCache = nil //TODO msCache not write
			info.result = nil
		}
		info.gInfo = sdk.GasInfo{GasWanted: info.gasWanted, GasUsed: info.ctx.GasMeter().GasConsumed()}
	}()

	defer handler.handleDeferGasConsumed(info)

	defer func() {
		app.pin(Refund, true, mode)
		defer app.pin(Refund, false, mode)
		handler.handleDeferRefund(info)
	}()

	if err := validateBasicTxMsgs(info.tx.GetMsgs()); err != nil {
		return err
	}
	app.pin(ValTxMsgs, false, mode)

	app.pin(RunAnte, true, mode)
	if app.anteHandler != nil {
		err = app.runAnte(info, mode)
		if err != nil {
			return err
		}
	}
	app.pin(RunAnte, false, mode)

	app.pin(RunMsg, true, mode)
	err = handler.handleRunMsg(info)
	app.pin(RunMsg, false, mode)
	return err
}

func (app *BaseApp) runAnte(info *runTxInfo, mode runTxMode) error {

	var anteCtx sdk.Context

	// Cache wrap context before AnteHandler call in case it aborts.
	// This is required for both CheckTx and DeliverTx.
	// Ref: https://github.com/cosmos/cosmos-sdk/issues/2772
	//
	// NOTE: Alternatively, we could require that AnteHandler ensures that
	// writes do not happen if aborted/failed.  This may have some
	// performance benefits, but it'll be more difficult to get right.

	// 1. CacheTxContext
	app.pin(CacheTxContext, true, mode)
	anteCtx, info.msCacheAnte = app.cacheTxContext(info.ctx, info.txBytes)
	anteCtx = anteCtx.WithEventManager(sdk.NewEventManager())
	app.pin(CacheTxContext, false, mode)

	// 2. AnteChain
	app.pin(AnteChain, true, mode)
	if mode == runTxModeDeliver {
		anteCtx = anteCtx.WithAnteTracer(app.anteTracer)
	}
	newCtx, err := app.anteHandler(anteCtx, info.tx, mode == runTxModeSimulate) // NewAnteHandler
	app.pin(AnteChain, false, mode)


	// 3. AnteOther
	app.pin(AnteOther, true, mode)
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
		info.ctx = newCtx.WithMultiStore(ms)
	}

	// GasMeter expected to be set in AnteHandler
	info.gasWanted = info.ctx.GasMeter().Limit()

	if mode == runTxModeDeliverInAsync {
		app.parallelTxManage.txStatus[string(info.txBytes)].anteErr = err
	}

	if err != nil {
		return err
	}
	app.pin(AnteOther, false, mode)


	// 4. CacheStoreWrite
	if mode != runTxModeDeliverInAsync {
		app.pin(CacheStoreWrite, true, mode)
		info.msCacheAnte.Write()
		info.ctx.Cache().Write(true)
		app.pin(CacheStoreWrite, false, mode)
	}

	return nil
}


func (app *BaseApp) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {

	tx, err := app.txDecoder(req.Tx)
	if err != nil {
		return sdkerrors.ResponseDeliverTx(err, 0, 0, app.trace)
	}

	gInfo, result, _, err := app.runTx(runTxModeDeliver, req.Tx, tx, LatestSimulateTxHeight)
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

func (app *BaseApp) asyncDeliverTx(txWithIndex []byte) {

	txStatus := app.parallelTxManage.txStatus[string(txWithIndex)]
	tx, err := app.txDecoder(getRealTxByte(txWithIndex))
	if err != nil {
		asyncExe := newExecuteResult(sdkerrors.ResponseDeliverTx(err, 0, 0, app.trace), nil, txStatus.indexInBlock, txStatus.evmIndex)
		app.parallelTxManage.workgroup.Push(asyncExe)
		return
	}

	if !txStatus.isEvmTx {
		asyncExe := newExecuteResult(abci.ResponseDeliverTx{}, nil, txStatus.indexInBlock, txStatus.evmIndex)
		app.parallelTxManage.workgroup.Push(asyncExe)
		return
	}

	var resp abci.ResponseDeliverTx
	g, r, m, e := app.runTx(runTxModeDeliverInAsync, txWithIndex, tx, LatestSimulateTxHeight)
	if e != nil {
		resp = sdkerrors.ResponseDeliverTx(e, g.GasWanted, g.GasUsed, app.trace)
	} else {
		resp = abci.ResponseDeliverTx{
			GasWanted: int64(g.GasWanted), // TODO: Should type accept unsigned ints?
			GasUsed:   int64(g.GasUsed),   // TODO: Should type accept unsigned ints?
			Log:       r.Log,
			Data:      r.Data,
			Events:    r.Events.ToABCIEvents(),
		}
	}

	asyncExe := newExecuteResult(resp, m, txStatus.indexInBlock, txStatus.evmIndex)
	asyncExe.err = e
	app.parallelTxManage.workgroup.Push(asyncExe)
}

func useCache(mode runTxMode) bool {
	if !sdk.UseCache {
		return false
	}
	if mode == runTxModeDeliver {
		return true
	}
	return false
}

func (app *BaseApp) newBlockCache() {
	app.blockCache = sdk.NewCache(app.chainCache, useCache(runTxModeDeliver))
	app.deliverState.ctx = app.deliverState.ctx.WithCache(app.blockCache)
}

func (app *BaseApp) commitBlockCache() {
	app.blockCache.Write(true)
	app.chainCache.TryDelete(app.logger, app.deliverState.ctx.BlockHeight())
}

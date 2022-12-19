package baseapp

import (
	"fmt"
	"runtime/debug"

	"github.com/pkg/errors"

	"github.com/okex/exchain/libs/system/trace"

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
	txIndex int

	reusableCacheMultiStore sdk.CacheMultiStore
	overridesBytes          []byte
}

func (info *runTxInfo) GetCacheMultiStore() (sdk.CacheMultiStore, bool) {
	if info.reusableCacheMultiStore == nil {
		return nil, false
	}
	reuse := info.reusableCacheMultiStore
	info.reusableCacheMultiStore = nil
	return reuse, true
}

func (info *runTxInfo) PutCacheMultiStore(cms sdk.CacheMultiStore) {
	info.reusableCacheMultiStore = cms
}

func (app *BaseApp) GetCacheMultiStore(txBytes []byte, height int64) (sdk.CacheMultiStore, bool) {
	if app.reusableCacheMultiStore == nil {
		return nil, false
	}
	reuse := updateCacheMultiStore(app.reusableCacheMultiStore, txBytes, height)
	app.reusableCacheMultiStore = nil
	return reuse, true
}

func (app *BaseApp) PutCacheMultiStore(cms sdk.CacheMultiStore) {
	app.reusableCacheMultiStore = cms
}

func (app *BaseApp) runTxWithIndex(txIndex int, mode runTxMode,
	txBytes []byte, tx sdk.Tx, height int64, from ...string) (info *runTxInfo, err error) {

	info = &runTxInfo{txIndex: txIndex}
	err = app.runtxWithInfo(info, mode, txBytes, tx, height, from...)
	return
}

func (app *BaseApp) runTx(mode runTxMode,
	txBytes []byte, tx sdk.Tx, height int64, from ...string) (info *runTxInfo, err error) {

	info = &runTxInfo{}
	err = app.runtxWithInfo(info, mode, txBytes, tx, height, from...)
	if app.watcherCollector != nil && mode == runTxModeDeliver {
		app.watcherCollector(info.runMsgCtx.GetWatcher())
	}
	return
}

func (app *BaseApp) runtxWithInfo(info *runTxInfo, mode runTxMode, txBytes []byte, tx sdk.Tx, height int64, from ...string) (err error) {
	info.handler = app.getModeHandler(mode)
	info.tx = tx
	info.txBytes = txBytes
	handler := info.handler
	app.pin(trace.ValTxMsgs, true, mode)

	if tx.GetType() != sdk.EvmTxType && mode == runTxModeDeliver {
		// should update the balance of FeeCollector's account when run non-evm tx
		// which uses non-infiniteGasMeter during AnteHandleChain
		app.updateFeeCollectorAccount(false)
	}

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
		if useCache(mode) && tx.GetType() == sdk.EvmTxType {
			info.ctx.SetCache(sdk.NewCache(app.blockCache, true))
		} else {
			info.ctx.SetCache(nil)
		}
	}
	for _, addr := range from {
		// cache from if exist
		if addr != "" {
			info.ctx.SetFrom(addr)
			break
		}
	}

	err = handler.handleGasConsumed(info)
	if err != nil {
		return err
	}

	// There is no need to update BlockGasMeter.GasConsumed and info.gInfo using ctx.GasMeter
	// as gas is not consumed actually when ante failed.
	isAnteSucceed := false
	defer func() {
		if r := recover(); r != nil {
			err = app.runTx_defer_recover(r, info)
			info.msCache = nil //TODO msCache not write
			info.result = nil
		}
		gasUsed := info.ctx.GasMeter().GasConsumed()
		if !isAnteSucceed {
			gasUsed = 0
		}
		info.gInfo = sdk.GasInfo{GasWanted: info.gasWanted, GasUsed: gasUsed}
		if mode == runTxModeDeliver {
			if cms, ok := info.GetCacheMultiStore(); ok {
				app.PutCacheMultiStore(cms)
			}
		}
	}()

	defer func() {
		if isAnteSucceed {
			handler.handleDeferGasConsumed(info)
		}
	}()

	defer func() {
		app.pin(trace.Refund, true, mode)
		defer app.pin(trace.Refund, false, mode)
		handler.handleDeferRefund(info)
	}()

	if err := validateBasicTxMsgs(info.tx.GetMsgs()); err != nil {
		return err
	}
	app.pin(trace.ValTxMsgs, false, mode)

	if mode == runTxModeDeliver {
		if cms, ok := app.GetCacheMultiStore(info.txBytes, info.ctx.BlockHeight()); ok {
			info.PutCacheMultiStore(cms)
		}
	}

	app.pin(trace.RunAnte, true, mode)
	if app.anteHandler != nil {
		err = app.runAnte(info, mode)
		if err != nil {
			return err
		}
	}
	app.pin(trace.RunAnte, false, mode)

	if app.getTxFeeHandler != nil && mode == runTxModeDeliver {
		fee := app.getTxFeeHandler(tx)
		app.UpdateFeeCollector(fee, true)
	}

	isAnteSucceed = true
	app.pin(trace.RunMsg, true, mode)
	err = handler.handleRunMsg(info)
	app.pin(trace.RunMsg, false, mode)
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
	app.pin(trace.CacheTxContext, true, mode)
	if mode == runTxModeDeliver {
		if cms, ok := info.GetCacheMultiStore(); ok {
			anteCtx, info.msCacheAnte = info.ctx, cms
			anteCtx.SetMultiStore(info.msCacheAnte)
		} else {
			anteCtx, info.msCacheAnte = app.cacheTxContext(info.ctx, info.txBytes)
		}
	} else if mode == runTxModeCheck || mode == runTxModeReCheck {
		info.msCacheAnte = app.checkTxCacheMultiStores.GetStore()
		if info.msCacheAnte != nil {
			info.msCacheAnte = updateCacheMultiStore(info.msCacheAnte, info.txBytes, info.ctx.BlockHeight())
			anteCtx = info.ctx
			anteCtx.SetMultiStore(info.msCacheAnte)
		} else {
			anteCtx, info.msCacheAnte = app.cacheTxContext(info.ctx, info.txBytes)
		}
	} else if mode == runTxModeDeliverInAsync {
		anteCtx = info.ctx
		info.msCacheAnte = nil
		msCacheAnte, useCurrentState := app.parallelTxManage.getParentMsByTxIndex(info.txIndex)
		if msCacheAnte == nil {
			return errors.New("Need Skip:txIndex smaller than currentIndex")
		}
		info.ctx.ParaMsg().UseCurrentState = useCurrentState
		info.msCacheAnte = msCacheAnte
		anteCtx.SetMultiStore(info.msCacheAnte)
	} else {
		anteCtx, info.msCacheAnte = app.cacheTxContext(info.ctx, info.txBytes)
	}

	anteCtx.SetEventManager(sdk.NewEventManager())
	app.pin(trace.CacheTxContext, false, mode)

	// 2. AnteChain
	app.pin(trace.AnteChain, true, mode)
	if mode == runTxModeDeliver {
		anteCtx.SetAnteTracer(app.anteTracer)
	}
	newCtx, err := app.anteHandler(anteCtx, info.tx, mode == runTxModeSimulate) // NewAnteHandler
	app.pin(trace.AnteChain, false, mode)

	// 3. AnteOther
	app.pin(trace.AnteOther, true, mode)
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
		info.ctx = newCtx
		info.ctx.SetMultiStore(ms)
	}

	// GasMeter expected to be set in AnteHandler
	info.gasWanted = info.ctx.GasMeter().Limit()

	if mode == runTxModeDeliverInAsync {
		info.ctx.ParaMsg().AnteErr = err
	}

	if err != nil {
		return err
	}
	app.pin(trace.AnteOther, false, mode)

	// 4. CacheStoreWrite
	if mode != runTxModeDeliverInAsync {
		app.pin(trace.CacheStoreWrite, true, mode)
		info.msCacheAnte.Write()
		if mode == runTxModeDeliver {
			info.PutCacheMultiStore(info.msCacheAnte)
			info.msCacheAnte = nil
		} else if mode == runTxModeCheck || mode == runTxModeReCheck {
			app.checkTxCacheMultiStores.PushStore(info.msCacheAnte)
			info.msCacheAnte = nil
		}
		info.ctx.Cache().Write(true)
		app.pin(trace.CacheStoreWrite, false, mode)
	}

	return nil
}

func (app *BaseApp) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {

	var realTx sdk.Tx
	var err error
	if mem := GetGlobalMempool(); mem != nil {
		realTx, _ = mem.ReapEssentialTx(req.Tx).(sdk.Tx)
	}
	if realTx == nil {
		realTx, err = app.txDecoder(req.Tx)
		if err != nil {
			return sdkerrors.ResponseDeliverTx(err, 0, 0, app.trace)
		}
	}

	info, err := app.runTx(runTxModeDeliver, req.Tx, realTx, LatestSimulateTxHeight)
	if err != nil {
		return sdkerrors.ResponseDeliverTx(err, info.gInfo.GasWanted, info.gInfo.GasUsed, app.trace)
	}

	if app.updateGPOHandler != nil {
		app.updateGPOHandler([]sdk.DynamicGasInfo{sdk.NewDynamicGasInfo(realTx.GetGasPrice(), info.gInfo.GasUsed)})
	}

	return abci.ResponseDeliverTx{
		GasWanted: int64(info.gInfo.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(info.gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
		Log:       info.result.Log,
		Data:      info.result.Data,
		Events:    info.result.Events.ToABCIEvents(),
	}
}

func (app *BaseApp) PreDeliverRealTx(tx []byte) abci.TxEssentials {
	var realTx sdk.Tx
	var err error
	if mem := GetGlobalMempool(); mem != nil {
		realTx, _ = mem.ReapEssentialTx(tx).(sdk.Tx)
	}
	if realTx == nil {
		realTx, err = app.txDecoder(tx)
		if err != nil || realTx == nil {
			return nil
		}
	}
	app.blockDataCache.SetTx(tx, realTx)

	if realTx.GetType() == sdk.EvmTxType && app.preDeliverTxHandler != nil {
		ctx := app.deliverState.ctx
		ctx.SetCache(app.chainCache).
			SetMultiStore(app.cms).
			SetGasMeter(sdk.NewInfiniteGasMeter())

		app.preDeliverTxHandler(ctx, realTx, !app.chainCache.IsEnabled())
	}

	return realTx
}

func (app *BaseApp) DeliverRealTx(txes abci.TxEssentials) abci.ResponseDeliverTx {
	var err error
	realTx, _ := txes.(sdk.Tx)
	if realTx == nil {
		realTx, err = app.txDecoder(txes.GetRaw())
		if err != nil {
			return sdkerrors.ResponseDeliverTx(err, 0, 0, app.trace)
		}
	}
	info, err := app.runTx(runTxModeDeliver, realTx.GetRaw(), realTx, LatestSimulateTxHeight)
	if !info.ctx.Cache().IsEnabled() {
		app.blockCache = nil
		app.chainCache = nil
	}
	if err != nil {
		return sdkerrors.ResponseDeliverTx(err, info.gInfo.GasWanted, info.gInfo.GasUsed, app.trace)
	}

	if app.updateGPOHandler != nil {
		app.updateGPOHandler([]sdk.DynamicGasInfo{sdk.NewDynamicGasInfo(realTx.GetGasPrice(), info.gInfo.GasUsed)})
	}

	return abci.ResponseDeliverTx{
		GasWanted: int64(info.gInfo.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(info.gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
		Log:       info.result.Log,
		Data:      info.result.Data,
		Events:    info.result.Events.ToABCIEvents(),
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
				"recovered: %v\n", r,
			),
		)
		app.logger.Info("runTx panic", "recover", r, "stack", string(debug.Stack()))
	}
	return err
}

func (app *BaseApp) asyncDeliverTx(txIndex int) *executeResult {
	pm := app.parallelTxManage
	if app.deliverState == nil { // runTxs already finish
		return nil
	}

	blockHeight := app.deliverState.ctx.BlockHeight()

	txStatus := app.parallelTxManage.extraTxsInfo[txIndex]

	if txStatus.stdTx == nil {
		asyncExe := newExecuteResult(sdkerrors.ResponseDeliverTx(txStatus.decodeErr,
			0, 0, app.trace), nil, uint32(txIndex), nil, blockHeight, sdk.EmptyWatcher{}, nil, app.parallelTxManage, nil)
		return asyncExe
	}

	if !txStatus.isEvm {
		asyncExe := newExecuteResult(abci.ResponseDeliverTx{}, nil, uint32(txIndex), nil,
			blockHeight, sdk.EmptyWatcher{}, nil, app.parallelTxManage, nil)
		return asyncExe
	}

	var resp abci.ResponseDeliverTx
	info, errM := app.runTxWithIndex(txIndex, runTxModeDeliverInAsync, pm.txs[txIndex], txStatus.stdTx, LatestSimulateTxHeight)
	if errM != nil {
		resp = sdkerrors.ResponseDeliverTx(errM, info.gInfo.GasWanted, info.gInfo.GasUsed, app.trace)
	} else {
		resp = abci.ResponseDeliverTx{
			GasWanted: int64(info.gInfo.GasWanted), // TODO: Should type accept unsigned ints?
			GasUsed:   int64(info.gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
			Log:       info.result.Log,
			Data:      info.result.Data,
			Events:    info.result.Events.ToABCIEvents(),
		}
	}

	asyncExe := newExecuteResult(resp, info.msCacheAnte, uint32(txIndex), info.ctx.ParaMsg(),
		blockHeight, info.runMsgCtx.GetWatcher(), info.tx.GetMsgs(), app.parallelTxManage, info.ctx.GetFeeSplitInfo())
	app.parallelTxManage.addMultiCache(info.msCacheAnte, info.msCache)
	return asyncExe
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
	useCache := sdk.UseCache && !app.parallelTxManage.isAsyncDeliverTx
	if app.chainCache == nil {
		app.chainCache = sdk.NewCache(nil, useCache)
	}

	app.blockCache = sdk.NewCache(app.chainCache, useCache)
}

func (app *BaseApp) commitBlockCache() {
	app.blockCache.Write(true)
	app.chainCache.TryDelete(app.logger, app.deliverState.ctx.BlockHeight())
}

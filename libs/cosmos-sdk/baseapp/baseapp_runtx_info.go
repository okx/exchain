package baseapp

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)


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

	result1 *sdk.Result
	txBytes []byte
	tx sdk.Tx
	finished bool
	decoded bool
}

func (app *BaseApp) runTx(mode runTxMode,  // DeliverTxConcurrently
	txBytes []byte, tx sdk.Tx, height int64) (gInfo sdk.GasInfo,
	result *sdk.Result, msCacheList sdk.CacheMultiStore, err error) {

	var info *runTxInfo
	info, err = app.runtx6(mode, txBytes, tx, height)
	return info.gInfo, info.result1, info.msCacheAnte, err
}

func (app *BaseApp) runtx6_1(info *runTxInfo, mode runTxMode, height int64) (err error) {

	mhandler := app.getModeHandler(mode)
	info.handler = mhandler

	fmt.Printf("runtx6-1\n")
	err = mhandler.handleStartHeight(info, height)
	if err != nil {
		return err
	}

	info.startingGas, info.gInfo, err = mhandler.handleGasConsumed(info)
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


func (app *BaseApp) runtx6_2(info *runTxInfo) (err error) {

	defer func() {
		if r := recover(); r != nil {
			err = app.runTx_defer_recover(r, info)
			info.msCache = nil //TODO msCache not write
			info.result1 = nil
			app.logger.Info("info.result = nil")
		}
		info.gInfo = sdk.GasInfo{GasWanted: info.gasWanted, GasUsed: info.ctx.GasMeter().GasConsumed()}
	}()

	defer app.runTx_defer_consumegas(info, info.handler.getMode())
	defer app.runTx_defer_refund(info, info.handler.getMode())

	if info.finished {
		return
	}

	info.result1, err = info.handler.handleRunMsg(info)
	if err == nil && info.result1 == nil {
		panic("")
	}
	return
}


func (app *BaseApp) runtx6(mode runTxMode, txBytes []byte, tx sdk.Tx, height int64) (info *runTxInfo, err error) {
	info = &runTxInfo{}

	info.handler = app.getModeHandler(mode)
	info.tx = tx
	info.txBytes = txBytes
	mhandler := info.handler

	fmt.Printf("runtx6\n")
	err = mhandler.handleStartHeight(info, height)
	if err != nil {
		return info, err
	}

	info.startingGas, info.gInfo, err = mhandler.handleGasConsumed(info)
	if err != nil {
		return info, err
	}

	defer func() {
		if r := recover(); r != nil {
			err = app.runTx_defer_recover(r, info)
			info.msCache = nil //TODO msCache not write
			info.result1 = nil
		}
		info.gInfo = sdk.GasInfo{GasWanted: info.gasWanted, GasUsed: info.ctx.GasMeter().GasConsumed()}
	}()

	defer app.runTx_defer_consumegas(info, mode)
	defer app.runTx_defer_refund(info, mode)

	if err := validateBasicTxMsgs(info.tx.GetMsgs()); err != nil {
		return info, err
	}

	if app.anteHandler != nil {
		err = app.runAnte(info, mode)
		if err != nil {
			return info, err
		}
	}

	info.result1, err = mhandler.handleRunMsg(info)
	return info, err
}


//func (app *BaseApp) runTxModeSimulate(txBytes []byte,
//	tx sdk.Tx, height int64) (gInfo sdk.GasInfo,
//	result *sdk.Result,
//	msCacheList sdk.CacheMultiStore,
//	err error) {
//
//	fmt.Printf("runtx2\n")
//
//	// NOTE: GasWanted should be returned by the AnteHandler. GasUsed is
//	// determined by the GasMeter. We need access to the context to get the gas
//	// meter so we initialize upfront.
//	var gasWanted uint64
//
//	var ctx sdk.Context
//	var runMsgCtx sdk.Context
//	var msCache sdk.CacheMultiStore
//	var msCacheAnte sdk.CacheMultiStore
//	var runMsgFinish bool
//	// simulate tx
//	startHeight := tmtypes.GetStartBlockHeight()
//
//	if height > startHeight && height < app.LastBlockHeight() {
//		ctx, err = app.getContextForSimTx(txBytes, height)
//		if err != nil {
//			return gInfo, result, nil, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
//		}
//	}
//
//
//
//	var startingGas uint64
//
//
//
//
//	defer func() {
//
//		if r := recover(); r != nil {
//			err = app.runTx_defer_recover(r, &ctx, gasWanted)
//			msCacheList = msCacheAnte
//			msCache = nil //TODO msCache not write
//			result = nil
//		}
//		gInfo = sdk.GasInfo{GasWanted: gasWanted, GasUsed: ctx.GasMeter().GasConsumed()}
//	}()
//
//	// If BlockGasMeter() panics it will be caught by the above recover and will
//	// return an error - in any case BlockGasMeter will consume gas past the limit.
//	//
//	// NOTE: This must exist in a separate defer function for the above recovery
//	// to recover from this one.
//	defer app.runTx_defer_consumegas(&ctx, mode, txBytes, startingGas)
//
//	defer func() {
//		msCache = app.runTx_defer_refund(&ctx, &runMsgCtx, mode, tx, txBytes, msCache, msCacheAnte, runMsgFinish)
//	}()
//
//
//	msgs := tx.GetMsgs()
//	if err := validateBasicTxMsgs(msgs); err != nil {
//		return sdk.GasInfo{}, nil, nil, err
//	}
//
//	accountNonce := uint64(0)
//	if app.anteHandler != nil {
//		accountNonce, gasWanted, err = app.runAnte(&ctx, mode, tx, txBytes, msCacheAnte)
//		if err != nil {
//			return gInfo, nil, nil, err
//		}
//	}
//
//	// Create a new Context based off of the existing Context with a cache-wrapped
//	// MultiStore in case message processing fails. At this point, the MultiStore
//	// is doubly cached-wrapped.
//
//
//
//	runMsgCtx, msCache = app.cacheTxContext(ctx, txBytes)
//
//
//	// Attempt to execute all messages and only update state if all messages pass
//	// and we're in DeliverTx. Note, runMsgs will never return a reference to a
//	// Result if any single message fails or does not have a registered Handler.
//
//	result, err = app.runMsgs(runMsgCtx, msgs, mode)
//
//
//
//	runMsgFinish = true
//
//
//
//
//	if err != nil {
//		if sdk.HigherThanMercury(ctx.BlockHeight()) {
//			codeSpace, code, info := sdkerrors.ABCIInfo(err, app.trace)
//			err = sdkerrors.New(codeSpace, abci.CodeTypeNonceInc+code, info)
//		}
//		msCache = nil
//	}
//
//
//	return gInfo, result, nil, err
//}

//func (app *BaseApp) runtx2(mode runTxMode, txBytes []byte,
//	tx sdk.Tx, height int64) (gInfo sdk.GasInfo,
//	result *sdk.Result,
//	msCacheList sdk.CacheMultiStore,
//	err error) {
//
//	fmt.Printf("runtx2\n")
//
//	// NOTE: GasWanted should be returned by the AnteHandler. GasUsed is
//	// determined by the GasMeter. We need access to the context to get the gas
//	// meter so we initialize upfront.
//	var gasWanted uint64
//
//	var ctx sdk.Context
//	var runMsgCtx sdk.Context
//	var msCache sdk.CacheMultiStore
//	var msCacheAnte sdk.CacheMultiStore
//	var runMsgFinish bool
//	// simulate tx
//	startHeight := tmtypes.GetStartBlockHeight()
//	if mode == runTxModeSimulate && height > startHeight && height < app.LastBlockHeight() {
//		ctx, err = app.getContextForSimTx(txBytes, height)
//		if err != nil {
//			return gInfo, result, nil, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
//		}
//	} else if height < startHeight && height != 0 {
//		return gInfo, result, nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
//			fmt.Sprintf("height(%d) should be greater than start block height(%d)", height, startHeight))
//	} else {
//		ctx = app.getContextForTx(mode, txBytes)
//	}
//
//
//	var startingGas uint64
//	if mode == runTxModeDeliver || mode == runTxModeDeliverInAsync {
//		if ctx.BlockGasMeter().IsOutOfGas() {
//			gInfo = sdk.GasInfo{GasUsed: ctx.BlockGasMeter().GasConsumed()}
//			return gInfo, nil, nil,
//				sdkerrors.Wrap(sdkerrors.ErrOutOfGas, "no block gas left to run tx")
//		}
//		startingGas = ctx.BlockGasMeter().GasConsumed()
//	}
//
//
//	defer func() {
//		app.pin(Recover, true, mode)
//		defer app.pin(Recover, false, mode)
//		if r := recover(); r != nil {
//			err = app.runTx_defer_recover(r, &ctx, gasWanted)
//			msCacheList = msCacheAnte
//			msCache = nil //TODO msCache not write
//			result = nil
//		}
//		gInfo = sdk.GasInfo{GasWanted: gasWanted, GasUsed: ctx.GasMeter().GasConsumed()}
//	}()
//
//	// If BlockGasMeter() panics it will be caught by the above recover and will
//	// return an error - in any case BlockGasMeter will consume gas past the limit.
//	//
//	// NOTE: This must exist in a separate defer function for the above recovery
//	// to recover from this one.
//	defer app.runTx_defer_consumegas(&ctx, mode, txBytes, startingGas)
//
//	defer func() {
//		msCache = app.runTx_defer_refund(&ctx, &runMsgCtx, mode, tx, txBytes, msCache, msCacheAnte, runMsgFinish)
//	}()
//
//
//	msgs := tx.GetMsgs()
//	if err := validateBasicTxMsgs(msgs); err != nil {
//		return sdk.GasInfo{}, nil, nil, err
//	}
//
//	accountNonce := uint64(0)
//	if app.anteHandler != nil {
//		accountNonce, gasWanted, err = app.runAnte(&ctx, mode, tx, txBytes, msCacheAnte)
//		if err != nil {
//			return gInfo, nil, nil, err
//		}
//	}
//
//	// Create a new Context based off of the existing Context with a cache-wrapped
//	// MultiStore in case message processing fails. At this point, the MultiStore
//	// is doubly cached-wrapped.
//
//	if mode == runTxModeDeliverInAsync {
//		msCache = msCacheAnte.CacheMultiStore()
//		runMsgCtx = ctx.WithMultiStore(msCache)
//	} else {
//		runMsgCtx, msCache = app.cacheTxContext(ctx, txBytes)
//	}
//
//	// Attempt to execute all messages and only update state if all messages pass
//	// and we're in DeliverTx. Note, runMsgs will never return a reference to a
//	// Result if any single message fails or does not have a registered Handler.
//
//	result, err = app.runMsgs(runMsgCtx, msgs, mode)
//	if err == nil && (mode == runTxModeDeliver) {
//		msCache.Write()
//	}
//
//	runMsgFinish = true
//
//	if mode == runTxModeCheck {
//		exTxInfo := app.GetTxInfo(ctx, tx)
//		exTxInfo.SenderNonce = accountNonce
//
//		data, err := json.Marshal(exTxInfo)
//		if err == nil {
//			result.Data = data
//		}
//	}
//
//	if err != nil {
//		if sdk.HigherThanMercury(ctx.BlockHeight()) {
//			codeSpace, code, info := sdkerrors.ABCIInfo(err, app.trace)
//			err = sdkerrors.New(codeSpace, abci.CodeTypeNonceInc+code, info)
//		}
//		msCache = nil
//	}
//
//	if mode == runTxModeDeliverInAsync {
//		if msCache != nil {
//			msCache.Write()
//		}
//		return gInfo, result, msCacheAnte, err
//	}
//	app.pin(RunMsgs, false, mode)
//	return gInfo, result, nil, err
//}

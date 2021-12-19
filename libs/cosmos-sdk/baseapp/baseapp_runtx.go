package baseapp

import (
	"encoding/json"
	"fmt"
	"sync"

	//"github.com/Workiva/go-datastructures/threadsafe/err"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"runtime/debug"
)

func (app *BaseApp) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	//res := app.DeliverTxCon(req)
	res := app.DeliverTxOrg(req)

	app.logger.Info("===========DeliverTx===========",
		"block", app.LastBlockHeight()+1,
		"Data len", len(res.Data),
		"Info", res.Info,
		"GasUsed", res.GasUsed,
		"GasWanted", res.GasWanted,
		"Code", res.Code,
	)
	for i, e := range res.Events {
		app.logger.Info("	Event", "id", i, "type", e.Type)
		for j, a := range e.Attributes {
			app.logger.Info("		Attributes", "id", j, "k", string(a.Key), "v", string(a.Value))
		}
	}

	return res
}

func (app *BaseApp) DeliverTxCon(req abci.RequestDeliverTx) abci.ResponseDeliverTx {

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

	var wg sync.WaitGroup
	wg.Add(len(txList))
	var taskList []task
	for i, tx := range txList {
		taskList = append(taskList, newTask(i, tx, ctx, &wg, app))
	}

	app.scheduler.start(taskList)
	wg.Wait()

	var results []*abci.ResponseDeliverTx
	for _, task := range taskList {
		results = append(results, task.result())
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
func (app *BaseApp) runTx_defer_recover(r interface{}, ctx *sdk.Context, gasWanted uint64) error {
	var err error
	switch rType := r.(type) {
	// TODO: Use ErrOutOfGas instead of ErrorOutOfGas which would allow us
	// to keep the stracktrace.
	case sdk.ErrorOutOfGas:
		err = sdkerrors.Wrap(
			sdkerrors.ErrOutOfGas, fmt.Sprintf(
				"out of gas in location: %v; gasWanted: %d, gasUsed: %d",
				rType.Descriptor, gasWanted, ctx.GasMeter().GasConsumed(),
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

func (app *BaseApp) runTx_defer_consumegas(ctx *sdk.Context, mode runTxMode, txBytes []byte, startingGas uint64) {
	app.pin(ConsumeGas, true, mode)
	defer app.pin(ConsumeGas, false, mode)
	if mode == runTxModeDeliver || (mode == runTxModeDeliverInAsync && app.parallelTxManage.isReRun(string(txBytes))) {
		ctx.BlockGasMeter().ConsumeGas(ctx.GasMeter().GasConsumedToLimit(), "block gas meter",)

		if ctx.BlockGasMeter().GasConsumed() < startingGas {
			panic(sdk.ErrorGasOverflow{Descriptor: "tx gas summation"})
		}
	}
}


func (app *BaseApp) runTx_defer_refund(ctx *sdk.Context,
	runMsgCtx *sdk.Context,
	mode runTxMode,
	tx sdk.Tx,
	txBytes []byte,
	msCache sdk.CacheMultiStore,
	msCacheAnte sdk.CacheMultiStore,
	runMsgFinish bool,
	) (sdk.CacheMultiStore) {

	app.pin(Refund, true, mode)
	defer app.pin(Refund, false, mode)
	if (mode == runTxModeDeliver || mode == runTxModeDeliverInAsync) && app.GasRefundHandler != nil {
		var gasRefundCtx sdk.Context
		if mode == runTxModeDeliver {
			gasRefundCtx, msCache = app.cacheTxContext(*ctx, txBytes)
		} else if mode == runTxModeDeliverInAsync {
			gasRefundCtx = *runMsgCtx
			if msCache == nil || !runMsgFinish { // case: panic when runMsg
				msCache = msCacheAnte.CacheMultiStore()
				gasRefundCtx = ctx.WithMultiStore(msCache)
			}
		}
		refundGas, err := app.GasRefundHandler(gasRefundCtx, tx)
		if err != nil {
			panic(err)
		}
		msCache.Write()
		if mode == runTxModeDeliverInAsync {
			app.parallelTxManage.setRefundFee(string(txBytes), refundGas)
		}
	}

	return msCache
}

func (app *BaseApp) runTxPart1(mode runTxMode, txBytes []byte, tx sdk.Tx,
	height int64, task *taskImp) (gInfo sdk.GasInfo,
	result *sdk.Result,
	msCacheList sdk.CacheMultiStore,
	err error) {

	app.pin(InitCtx, true, mode)

	// NOTE: GasWanted should be returned by the AnteHandler. GasUsed is
	// determined by the GasMeter. We need access to the context to get the gas
	// meter so we initialize upfront.
	var gasWanted uint64

	var ctx sdk.Context
	var runMsgCtx sdk.Context
	var msCache sdk.CacheMultiStore
	var msCacheAnte sdk.CacheMultiStore
	// simulate tx
	startHeight := tmtypes.GetStartBlockHeight()
	if mode == runTxModeSimulate && height > startHeight && height < app.LastBlockHeight() {
		ctx, err = app.getContextForSimTx(txBytes, height)
		if err != nil {
			return gInfo, result, nil, sdkerrors.Wrap(sdkerrors.ErrInternal, err.Error())
		}
	} else if height < startHeight && height != 0 {

		err = sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("height(%d) should be greater than start block height(%d)",
				height, startHeight))

		return gInfo, result, nil, err

	} else {
		ctx = app.getContextForTx(mode, txBytes)
	}

	ms := ctx.MultiStore()

	// only run the tx if there is block gas remaining
	if (mode == runTxModeDeliver || mode == runTxModeDeliverInAsync) && ctx.BlockGasMeter().IsOutOfGas() {
		gInfo = sdk.GasInfo{GasUsed: ctx.BlockGasMeter().GasConsumed()}
		err = sdkerrors.Wrap(sdkerrors.ErrOutOfGas, "no block gas left to run tx")
		return gInfo, nil, nil, err
	}

	var startingGas uint64
	if mode == runTxModeDeliver || mode == runTxModeDeliverInAsync {
		startingGas = ctx.BlockGasMeter().GasConsumed()
	}

	app.pin(InitCtx, false, mode)

	defer func() {
		app.pin(Recover, true, mode)
		defer app.pin(Recover, false, mode)
		if r := recover(); r != nil {
			err = app.runTx_defer_recover(r, &ctx, gasWanted)
			msCacheList = msCacheAnte
			msCache = nil //TODO msCache not write
			result = nil
		}
		gInfo = sdk.GasInfo{GasWanted: gasWanted, GasUsed: ctx.GasMeter().GasConsumed()}
	}()

	// If BlockGasMeter() panics it will be caught by the above recover and will
	// return an error - in any case BlockGasMeter will consume gas past the limit.
	//
	// NOTE: This must exist in a separate defer function for the above recovery
	// to recover from this one.
	defer app.runTx_defer_consumegas(&ctx, mode, txBytes, startingGas)

	defer func() {
		msCache = app.runTx_defer_refund(&ctx, &runMsgCtx, mode, tx, txBytes, msCache, msCacheAnte, false)
	}()

	app.pin(ValTxMsgs, true, mode)

	msgs := tx.GetMsgs()
	if err := validateBasicTxMsgs(msgs); err != nil {
		return sdk.GasInfo{}, nil, nil, err
	}
	app.pin(ValTxMsgs, false, mode)

	app.pin(AnteHandler, true, mode)

	accountNonce := uint64(0)
	if app.anteHandler != nil {
		var anteCtx sdk.Context

		// Cache wrap context before AnteHandler call in case it aborts.
		// This is required for both CheckTx and DeliverTx.
		// Ref: https://github.com/cosmos/cosmos-sdk/issues/2772
		//
		// NOTE: Alternatively, we could require that AnteHandler ensures that
		// writes do not happen if aborted/failed.  This may have some
		// performance benefits, but it'll be more difficult to get right.
		anteCtx, msCacheAnte = app.cacheTxContext(ctx, txBytes)
		anteCtx = anteCtx.WithEventManager(sdk.NewEventManager())
		newCtx, err := app.anteHandler(anteCtx, tx, mode == runTxModeSimulate)

		accountNonce = newCtx.AccountNonce()
		if !newCtx.IsZero() {
			// At this point, newCtx.MultiStore() is cache-wrapped, or something else
			// replaced by the AnteHandler. We want the original multistore, not one
			// which was cache-wrapped for the AnteHandler.
			//
			// Also, in the case of the tx aborting, we need to track gas consumed via
			// the instantiated gas meter in the AnteHandler, so we update the context
			// prior to returning.
			ctx = newCtx.WithMultiStore(ms)
		}

		// GasMeter expected to be set in AnteHandler
		gasWanted = ctx.GasMeter().Limit()

		if mode == runTxModeDeliverInAsync {
			app.parallelTxManage.txStatus[string(txBytes)].anteErr = err
		}

		if err != nil {
			return gInfo, nil, nil, err
		}

		if mode != runTxModeDeliverInAsync {
			msCacheAnte.Write()
		}
	}
	app.pin(AnteHandler, false, mode)

	app.pin(RunMsgs, true, mode)

	// Create a new Context based off of the existing Context with a cache-wrapped
	// MultiStore in case message processing fails. At this point, the MultiStore
	// is doubly cached-wrapped.
	if mode == runTxModeDeliverInAsync {
		msCache = msCacheAnte.CacheMultiStore()
		runMsgCtx = ctx.WithMultiStore(msCache)
	} else {
		runMsgCtx, msCache = app.cacheTxContext(ctx, txBytes)
	}

	task.ctx = &ctx
	task.runMsgCtx = &runMsgCtx
	task.mode = mode
	task.msCache = msCache
	task.msCacheAnte = msCacheAnte
	task.tx = tx
	task.accountNonce = accountNonce
	task.msgs = msgs
	task.gasWanted = gasWanted
	task.startingGas = startingGas

	return gInfo, nil, nil, nil
}

func (app *BaseApp) runTxPart2(task *taskImp) (gInfo sdk.GasInfo,
	result *sdk.Result,
	msCacheList sdk.CacheMultiStore,
	err error) {

	ctx := task.ctx
	runMsgCtx := task.runMsgCtx
	mode := task.mode
	msCache := task.msCache
	msCacheAnte := task.msCacheAnte
	tx := task.tx
	accountNonce := task.accountNonce
	msgs := task.msgs
	gasWanted := task.gasWanted
	startingGas := task.startingGas
	txBytes := task.txBytes

	var runMsgFinished bool

	defer func() {
		app.pin(Recover, true, mode)
		defer app.pin(Recover, false, mode)
		if r := recover(); r != nil {
			err = app.runTx_defer_recover(r, ctx, gasWanted)
			msCacheList = msCacheAnte
			msCache = nil //TODO msCache not write
			result = nil
		}
		gInfo = sdk.GasInfo{GasWanted: gasWanted, GasUsed: ctx.GasMeter().GasConsumed()}
	}()

	// If BlockGasMeter() panics it will be caught by the above recover and will
	// return an error - in any case BlockGasMeter will consume gas past the limit.
	//
	// NOTE: This must exist in a separate defer function for the above recovery
	// to recover from this one.
	defer app.runTx_defer_consumegas(ctx, mode, txBytes, startingGas)

	defer func() {
		msCache = app.runTx_defer_refund(ctx, runMsgCtx, mode, tx, txBytes, msCache, msCacheAnte, runMsgFinished)
	}()


	// Attempt to execute all messages and only update state if all messages pass
	// and we're in DeliverTx. Note, runMsgs will never return a reference to a
	// Result if any single message fails or does not have a registered Handler.

	result, err = app.runMsgs(*runMsgCtx, msgs, mode)
	if err == nil && (mode == runTxModeDeliver) {
		msCache.Write()
	}

	runMsgFinished = true

	if mode == runTxModeCheck {
		exTxInfo := app.GetTxInfo(*ctx, tx)
		exTxInfo.SenderNonce = accountNonce

		data, err := json.Marshal(exTxInfo)
		if err == nil {
			result.Data = data
		}
	}

	if err != nil {
		if sdk.HigherThanMercury(ctx.BlockHeight()) {
			codeSpace, code, info := sdkerrors.ABCIInfo(err, app.trace)
			err = sdkerrors.New(codeSpace, abci.CodeTypeNonceInc+code, info)
		}
		msCache = nil
	}

	if mode == runTxModeDeliverInAsync {
		if msCache != nil {
			msCache.Write()
		}
		return gInfo, result, msCacheAnte, err
	}
	app.pin(RunMsgs, false, mode)
	return gInfo, result, nil, err
}
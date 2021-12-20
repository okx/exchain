package baseapp

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"runtime/debug"
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

	result *sdk.Result
	txBytes []byte
	tx sdk.Tx
	finished bool
	decoded bool
}

var RunTxByRefactor bool = true
//var RunTxByRefactor bool = false

func (app *BaseApp) runTx(mode runTxMode,
	txBytes []byte, tx sdk.Tx, height int64) (gInfo sdk.GasInfo,
	result *sdk.Result, msCacheList sdk.CacheMultiStore, err error) {

	if RunTxByRefactor {
		var info *runTxInfo
		info, err = app.RunTxByRefactor(mode, txBytes, tx, height)
		return info.gInfo, info.result, info.msCacheAnte, err
	} else {
		return app.runtx_org(mode, txBytes, tx, height)
	}
}

func (app *BaseApp) RunTxByRefactor(mode runTxMode, txBytes []byte, tx sdk.Tx, height int64) (info *runTxInfo, err error) {
	info = &runTxInfo{}
	info.handler = app.getModeHandler(mode)
	info.tx = tx
	info.txBytes = txBytes
	handler := info.handler

	fmt.Printf("runtx_refactor\n")
	err = handler.handleStartHeight(info, height)
	if err != nil {
		return info, err
	}

	err = handler.handleGasConsumed(info)
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

	defer handler.handleDeferGasConsumed(info)
	defer handler.handleDeferRefund(info)

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

	err = handler.handleRunMsg(info)
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
	return app.DeliverTxOrg(req)
}

func (app *BaseApp) DeliverTxOrg(req abci.RequestDeliverTx) abci.ResponseDeliverTx {

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

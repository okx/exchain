package baseapp

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

func (app *BaseApp) PushAnteHandler(ah sdk.AnteHandler) {
	app.anteHandler = ah
}

func (app *BaseApp) GetDeliverStateCtx() sdk.Context {
	return app.deliverState.ctx
}

func (app *BaseApp) TraceTx(data []byte) (*sdk.Result, error) {

	var traceTxRequest sdk.QueryTraceParams
	err := codec.Cdc.UnmarshalJSON(data, &traceTxRequest)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid trace tx input")
	}
	// Begin block
	req := abci.RequestBeginBlock{
		Hash:   traceTxRequest.Block.Hash(),
		Header: tmtypes.TM2PB.Header(&traceTxRequest.Block.Header),
	}
	app.BeginBlock(req)
	//prepare context to deliver tx
	runningMode := runTxModeDeliver
	info := &runTxInfo{}
	info.handler = app.getModeHandler(runningMode)
	info.tx = *traceTxRequest.TraceTx

	var initialTx []byte
	if len(traceTxRequest.Predecessors) == 0 {
		initialTx = traceTxRequest.TxBytes
	} else {
		initialTx = traceTxRequest.PredecessorsBytes[0]
	}
	info.txBytes = initialTx

	info.ctx, err = app.getContextForSimTx(traceTxRequest.TxBytes, traceTxRequest.Block.Height)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "get context failed for trace tx")
	}
	info.ctx = info.ctx.WithCache(sdk.NewCache(app.blockCache, useCache(runningMode)))
	info.ctx.WithIsTraceTx(false)

	//pre deliver prodesessor tx to get the right context
	for index, prodesessor := range traceTxRequest.Predecessors {
		info, err = app.runTxWithInfo(info, runningMode, traceTxRequest.PredecessorsBytes[index], *prodesessor, traceTxRequest.Block.Height)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "run prodesessor failed for trace tx")
		}
	}

	//trace tx
	info.ctx.WithIsTraceTx(true)
	info, err = app.runTxWithInfo(info, runningMode, traceTxRequest.TxBytes, *traceTxRequest.TraceTx, traceTxRequest.Block.Height)
	return info.result, err
}

/*
 * @Author: worm
 * @Description:
 * @Date: 2021-11-08 16:19:11
 * @LastEditors: worm
 * @LastEditTime: 2022-01-13 20:25:05
 * @FilePath: /exchain/libs/cosmos-sdk/baseapp/helpers_okchain.go
 */
package baseapp

import (
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

func (app *BaseApp) TraceTx(txData []byte, targetTx sdk.Tx, txIndex uint32, block *tmtypes.Block) (*sdk.Result, error) {
	//prepare context to deliver tx
	runningMode := runTxModeTrace
	info := &runTxInfo{}
	info.handler = app.getModeHandler(runningMode)

	var initialTx sdk.Tx
	var initialTxBytes []byte
	predesessors := block.Txs[:txIndex]
	if len(predesessors) == 0 {
		initialTx = targetTx
		initialTxBytes = txData
	} else {
		tmp, err := app.txDecoder(predesessors[0])
		if err != nil {
			return nil, sdkerrors.Wrap(err, "invalid prodesessor")
		}
		initialTx = tmp
		initialTxBytes = predesessors[0]
	}
	info.tx = initialTx
	info.txBytes = initialTxBytes

	//begin block
	err := app.beginBlockForTracing(info, txData, block)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "invalid prodesessor")
	}

	info.ctx = info.ctx.WithIsTraceTx(false)

	//pre deliver prodesessor tx to get the right context
	for _, predesessor := range predesessors {
		tx, err := app.txDecoder(predesessor)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "invalid prodesessor")
		}
		info.tx = tx
		info.txBytes = predesessor
		info, err = app.runTxWithInfo(info, runningMode, tx, block.Height)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "run prodesessor failed for trace tx")
		}
	}

	//trace tx
	info.tx = targetTx
	info.txBytes = txData
	info.ctx = info.ctx.WithIsTraceTx(true)
	info, err = app.runTxWithInfo(info, runningMode, targetTx, block.Height)
	return info.result, err
}

func (app *BaseApp) beginBlockForTracing(info *runTxInfo, txData []byte, block *tmtypes.Block) error {

	// Begin block
	req := abci.RequestBeginBlock{
		Hash:   block.Hash(),
		Header: tmtypes.TM2PB.Header(&block.Header),
	}
	//app.setDeliverState(req.Header)
	var err error
	info.ctx, err = app.getContextForSimTx(txData, block.Height)
	if err != nil {
		return sdkerrors.Wrap(err, "get context failed for trace tx")
	}

	//app.newBlockCache()
	// use block cache instead of app.blockCache to save all tx results in one block
	chainCache := sdk.NewChainCache()
	blockCache := sdk.NewCache(chainCache, true)
	info.ctx = info.ctx.WithCache(blockCache)

	// add block gas meter
	var gasMeter sdk.GasMeter
	if maxGas := app.getMaximumBlockGas(); maxGas > 0 {
		gasMeter = sdk.NewGasMeter(maxGas)
	} else {
		gasMeter = sdk.NewInfiniteGasMeter()
	}

	info.ctx = info.ctx.WithBlockGasMeter(gasMeter)

	if app.beginBlocker != nil {
		_ = app.beginBlocker(info.ctx, req)
	}

	// set the signed validators for addition to context in deliverTx
	// No need
	return nil
}

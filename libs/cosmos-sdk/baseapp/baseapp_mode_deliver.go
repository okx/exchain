package baseapp

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func (m *modeHandlerDeliver) handleRunMsg(info *runTxInfo) (err error) {
	app := m.app
	mode := m.mode

	info.runMsgCtx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)
	info.ctx.Cache().Write(nil, info.ctx, false)
	info.result, err = app.runMsgs(info.runMsgCtx, info.tx.GetMsgs(), mode)
	if err == nil {
		info.msCache.Write()
		info.ctx.Cache().Write(nil, info.ctx,true)
	}

	info.runMsgFinished = true
	err = m.checkHigherThanMercury(err, info)
	return
}

func (m *modeHandlerDeliver) handleDeferRefund(info *runTxInfo) {
	app := m.app

	if app.GasRefundHandler == nil {
		return
	}

	var gasRefundCtx sdk.Context
	info.ctx.Cache().Write(nil, info.ctx,false)
	gasRefundCtx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)

	_, err := app.GasRefundHandler(gasRefundCtx, info.tx)
	if err != nil {
		panic(err)
	}
	info.msCache.Write()
	info.ctx.Cache().Write(nil, info.ctx,true)
}

func (m *modeHandlerDeliver) handleDeferGasConsumed(info *runTxInfo) {
	m.setGasConsumed(info)
}

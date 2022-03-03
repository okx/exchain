package baseapp

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func (m *modeHandlerDeliver) handleRunMsg(info *runTxInfo) (err error) {
	app := m.app
	mode := m.mode

	info.runMsgCtx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)
	info.runMsgCtx = info.runMsgCtx.WithCache(sdk.NewCache(info.ctx.Cache(), useCache(mode)))
	info.result, err = app.runMsgs(info.runMsgCtx, info.tx.GetMsgs(), mode)
	if err == nil {
		info.msCache.Write()
		info.runMsgCtx.Cache().Write(true)
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
	gasRefundCtx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)
	gasRefundCtx = gasRefundCtx.WithCache(sdk.NewCache(info.ctx.Cache(), useCache(m.mode)))

	_, err := app.GasRefundHandler(gasRefundCtx, info.tx)
	if err != nil {
		panic(err)
	}
	info.msCache.Write()
	gasRefundCtx.Cache().Write(true)
}

func (m *modeHandlerDeliver) handleDeferGasConsumed(info *runTxInfo) {
	m.setGasConsumed(info)
}

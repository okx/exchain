package baseapp

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func (m *modeHandlerDeliver) handleRunMsg(info *runTxInfo) (err error) {
	app := m.app
	mode := m.mode

	var ok bool
	if info.msCache, ok = info.GetCacheMultiStore(); ok {
		info.runMsgCtx = info.ctx
		info.runMsgCtx.SetMultiStore(info.msCache)
	} else {
		info.ctx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)
	}

	info.ctx.Cache().Write(false)
	info.result, err = app.runMsgs(info.runMsgCtx, info.tx.GetMsgs(), mode)
	if err == nil {
		info.msCache.Write()
		info.ctx.Cache().Write(true)
		info.PutCacheMultiStore(info.msCache)
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
	info.ctx.Cache().Write(false)

	var ok bool
	if info.msCache, ok = info.GetCacheMultiStore(); ok {
		gasRefundCtx = info.ctx
		gasRefundCtx.SetMultiStore(info.msCache)
	} else {
		gasRefundCtx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)
	}

	refund, err := app.GasRefundHandler(gasRefundCtx, info.tx)
	if err != nil {
		panic(err)
	}
	info.msCache.Write()
	info.PutCacheMultiStore(info.msCache)
	info.ctx.Cache().Write(true)

	app.UpdateFeeForCollector(refund, false)
}

func (m *modeHandlerDeliver) handleDeferGasConsumed(info *runTxInfo) {
	m.setGasConsumed(info)
}

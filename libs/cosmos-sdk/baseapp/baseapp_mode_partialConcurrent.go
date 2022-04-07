package baseapp

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

func (m *modeHandlerDeliverPartConcurrent) handleRunMsg(info *runTxInfo) (err error) {
	app := m.app
	mode := m.mode

	info.runMsgCtx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)
	info.ctx.Cache().Write(false)
	info.result, err = app.runMsgs(info.runMsgCtx, info.tx.GetMsgs(), mode)
	if err == nil {
		info.msCache.Write()
		info.ctx.Cache().Write(true)
	}

	info.runMsgFinished = true
	err = m.checkHigherThanMercury(err, info)
	return
}

func (m *modeHandlerDeliverPartConcurrent) handleDeferRefund(info *runTxInfo) {
	app := m.app

	if app.GasRefundHandler == nil {
		return
	}

	var gasRefundCtx sdk.Context
	info.ctx.Cache().Write(false)
	gasRefundCtx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)

	refundGas, err := app.GasRefundHandler(gasRefundCtx, info.tx)
	if err != nil {
		panic(err)
	}
	info.msCache.Write()
	info.ctx.Cache().Write(true)

	info.ctx = info.ctx.UpdateFeeForCollector(refundGas, false)
	app.UpdateFeeForCollector(info.ctx.FeeForCollector(), true)

	app.deliverTxsMgr.calculateFeeForCollector(refundGas, false)

	diff, hasNeg := app.feeForCollector.SafeSub(app.deliverTxsMgr.currTxFee)
	if hasNeg || !diff.IsZero() {
		app.logger.Error("NotEqual.", info.ctx.FeeForCollector(), app.deliverTxsMgr.currTxFee)
	}
}

func (m *modeHandlerDeliverPartConcurrent) handleDeferGasConsumed(info *runTxInfo) {
	m.setGasConsumed(info)
}

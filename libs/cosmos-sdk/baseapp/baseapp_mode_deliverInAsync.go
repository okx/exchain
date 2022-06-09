package baseapp

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

func (m *modeHandlerDeliverInAsync) handleDeferRefund(info *runTxInfo) {
	app := m.app

	if app.GasRefundHandler == nil {
		return
	}
	if info.msCacheAnte == nil {
		return
	}
	var gasRefundCtx sdk.Context
	gasRefundCtx = info.runMsgCtx
	if info.msCache == nil || !info.runMsgFinished { // case: panic when runMsg
		info.msCache = app.parallelTxManage.chainMultiStores.GetStoreWithParent(info.msCacheAnte)
		gasRefundCtx = info.ctx
		gasRefundCtx.SetMultiStore(info.msCache)
	}

	refundGas, err := app.GasRefundHandler(gasRefundCtx, info.tx)
	if err != nil {
		panic(err)
	}
	info.msCache.Write()
	info.ctx.ParaMsg().RefundFee = refundGas
}

func (m *modeHandlerDeliverInAsync) handleDeferGasConsumed(info *runTxInfo) {
}
func (m *modeHandlerDeliverInAsync) handleRunMsg(info *runTxInfo) (err error) {
	app := m.app
	mode := m.mode

	info.msCache = app.parallelTxManage.chainMultiStores.GetStoreWithParent(info.msCacheAnte)
	info.runMsgCtx = info.ctx
	info.runMsgCtx.SetMultiStore(info.msCache)

	info.result, err = app.runMsgs(info.runMsgCtx, info.tx.GetMsgs(), mode)
	info.runMsgFinished = true
	err = m.checkHigherThanMercury(err, info)

	if info.msCache != nil {
		info.msCache.Write()
	}

	return
}

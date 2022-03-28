package baseapp

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

func (m *modeHandlerDeliverInAsync) handleDeferRefund(info *runTxInfo) {
	app := m.app

	if app.GasRefundHandler == nil {
		return
	}
	var gasRefundCtx sdk.Context
	gasRefundCtx = info.runMsgCtx
	if info.msCache == nil || !info.runMsgFinished { // case: panic when runMsg
		info.msCache = info.msCacheAnte.CacheMultiStore()
		gasRefundCtx = info.ctx.WithMultiStore(info.msCache)
	}

	refundGas, err := app.GasRefundHandler(gasRefundCtx, info.tx)
	if err != nil {
		panic(err)
	}
	info.msCache.Write()
	info.ctx.ParaMsg().RefundFee = refundGas
}

func (m *modeHandlerDeliverInAsync) handleDeferGasConsumed(info *runTxInfo) {
	if m.app.parallelTxManage.isReRun(string(info.txBytes)) {
		m.setGasConsumed(info)
	}
}
func (m *modeHandlerDeliverInAsync) handleRunMsg(info *runTxInfo) (err error) {
	app := m.app
	mode := m.mode
	msCacheAnte := info.msCacheAnte

	info.msCache = msCacheAnte.CacheMultiStore()
	info.runMsgCtx = info.ctx.WithMultiStore(info.msCache)

	info.result, err = app.runMsgs(info.runMsgCtx, info.tx.GetMsgs(), mode)
	info.runMsgFinished = true
	err = m.checkHigherThanMercury(err, info)

	if info.msCache != nil {
		info.msCache.Write()
	}

	return
}

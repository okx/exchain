package baseapp

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

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
	if err == nil {
		info.msCache.Write()
	}
	err = m.wrapError(err, info)

	return
}

// ====================================================
// 2. handleGasConsumed
func (m *modeHandlerDeliverInAsync) handleGasConsumed(info *runTxInfo) (err error) {
	m.app.parallelTxManage.blockGasMeterMu.Lock()
	defer m.app.parallelTxManage.blockGasMeterMu.Unlock()
	return m.modeHandlerBase.handleGasConsumed(info)
}

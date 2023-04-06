package baseapp

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func (m *modeHandlerDeliver) handleRunMsg(info *runTxInfo) (err error) {
	app := m.app
	mode := m.mode
	if cms, ok := info.GetCacheMultiStore(); ok {
		info.runMsgCtx, info.msCache = info.ctx, cms
		info.runMsgCtx.SetMultiStore(info.msCache)
	} else {
		info.runMsgCtx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)
	}

	info.ctx.Cache().Write(false)
	info.result, err = app.runMsgs(info.runMsgCtx, info.tx.GetMsgs(), mode)
	if err == nil {
		info.msCache.Write()
		info.ctx.Cache().Write(true)
		info.PutCacheMultiStore(info.msCache)
		info.msCache = nil
	}

	info.runMsgFinished = true
	err = m.checkHigherThanMercury(err, info)
	return
}

type CacheTxContextFunc func(ctx sdk.Context, txBytes []byte) (sdk.Context, sdk.CacheMultiStore)

//this handleGasRefund func is also called by modeHandlerTrace.handleDeferRefund
//in this func, edit any member in BaseApp is prohibited
func handleGasRefund(info *runTxInfo, cacheTxCtxFunc CacheTxContextFunc, gasRefundHandler sdk.GasRefundHandler) sdk.DecCoins {
	var gasRefundCtx sdk.Context
	gasRefundCtx.SetOutOfGas(info.outOfGas)
	info.ctx.Cache().Write(false)
	if cms, ok := info.GetCacheMultiStore(); ok {
		gasRefundCtx, info.msCache = info.ctx, cms
		gasRefundCtx.SetMultiStore(info.msCache)
	} else {
		gasRefundCtx, info.msCache = cacheTxCtxFunc(info.ctx, info.txBytes)
	}

	refund, err := gasRefundHandler(gasRefundCtx, info.tx)
	if err != nil {
		panic(err)
	}
	info.msCache.Write()
	info.PutCacheMultiStore(info.msCache)
	info.msCache = nil
	info.ctx.Cache().Write(true)
	return refund
}
func (m *modeHandlerDeliver) handleDeferRefund(info *runTxInfo) {
	if m.app.GasRefundHandler == nil {
		return
	}
	refund := handleGasRefund(info, m.app.cacheTxContext, m.app.GasRefundHandler)
	m.app.UpdateFeeCollector(refund, false)
	if info.ctx.GetFeeSplitInfo().HasFee {
		m.app.FeeSplitCollector = append(m.app.FeeSplitCollector, info.ctx.GetFeeSplitInfo())
	}

}

func (m *modeHandlerDeliver) handleDeferGasConsumed(info *runTxInfo) {
	m.setGasConsumed(info)
}

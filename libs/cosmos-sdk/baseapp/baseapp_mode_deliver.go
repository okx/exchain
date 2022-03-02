package baseapp

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

func (m *modeHandlerDeliver) handleRunMsg(info *runTxInfo) (err error) {
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

func (a *BaseApp) printMs(extraData string) {
	if a.parallelTxManage.isAsyncDeliverTx {
		return
	}
	cnt := 0
	a.deliverState.ms.IteratorCache(func(key, value []byte, isDirty bool, isdelete bool, s sdk.StoreKey) bool {
		if isDirty {
			cnt++
		}

		//fmt.Println("key---", hex.EncodeToString(key), isDirty)
		return true
	}, nil)
	fmt.Println("PrintMsInfo", extraData, cnt)
}

func (m *modeHandlerDeliver) handleDeferRefund(info *runTxInfo) {
	app := m.app

	if app.GasRefundHandler == nil {
		return
	}

	var gasRefundCtx sdk.Context
	info.ctx.Cache().Write(false)
	gasRefundCtx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)

	_, err := app.GasRefundHandler(gasRefundCtx, info.tx)
	if err != nil {
		panic(err)
	}
	info.msCache.Write()
	info.ctx.Cache().Write(true)
	//app.printMs("Defer")
}

func (m *modeHandlerDeliver) handleDeferGasConsumed(info *runTxInfo) {
	m.setGasConsumed(info)
}

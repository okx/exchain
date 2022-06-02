package baseapp

func (m *modeHandlerTrace) handleStartHeight(info *runTxInfo, height int64) (err error) { return }

func (m *modeHandlerTrace) handleDeferRefund(info *runTxInfo) {
	if m.app.GasRefundHandler == nil {
		return
	}
	handleGasRefund(info, m.app.cacheTxContext, m.app.GasRefundHandler)
}

package baseapp

func (m *modeHandlerTrace) handleStartHeight(info *runTxInfo, height int64) (err error) { return }

func (m *modeHandlerTrace) handleDeferRefund(info *runTxInfo) {
	app := m.app
	if app.GasRefundHandler == nil {
		return
	}
	handleGasRefund(info, app.cacheTxContext, app.GasRefundHandler)
}

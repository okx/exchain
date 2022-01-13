package baseapp

func (m *modeHandlerTrace) handleStartHeight(info *runTxInfo, height int64) error {
	//no need to set info ctx for each tx, because it is set at first (in app.beginBlockForTracing)
	return nil
}

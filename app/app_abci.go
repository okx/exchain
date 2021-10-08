package app

import (
	appconfig "github.com/okex/exchain/app/config"
	"github.com/okex/exchain/x/analyzer"
	"github.com/okex/exchain/x/evm"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/trace"
)


// BeginBlock implements the Application interface
func (app *OKExChainApp) BeginBlock(req abci.RequestBeginBlock) (res abci.ResponseBeginBlock) {

	analyzer.OnAppBeginBlockEnter(app.Logger(), app.LastBlockHeight()+1)
	defer analyzer.OnAppBeginBlockExit()


	// dump app.LastBlockHeight()-1 info for reactor sync mode
	trace.GetElapsedInfo().Dump(app.Logger())
	return app.BaseApp.BeginBlock(req)
}


func (app *OKExChainApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {

	analyzer.OnAppDeliverTxEnter()
	defer analyzer.OnAppDeliverTxExit()

	resp := app.BaseApp.DeliverTx(req)
	if (app.BackendKeeper.Config.EnableBackend || app.StreamKeeper.AnalysisEnable()) && resp.IsOK() {
		app.syncTx(req.Tx)
	}

	if appconfig.GetOecConfig().GetEnableDynamicGp() {
		tx, err := evm.TxDecoder(app.Codec())(req.Tx)
		if err == nil {
			app.blockGasPrice = append(app.blockGasPrice, tx.GetTxInfo(app.GetDeliverStateCtx()).GasPrice)
		}
	}

	return resp
}

// EndBlock implements the Application interface
func (app *OKExChainApp) EndBlock(req abci.RequestEndBlock) (res abci.ResponseEndBlock) {

	analyzer.OnAppEndBlockEnter()
	defer analyzer.OnAppEndBlockExit()

	return app.BaseApp.EndBlock(req)
}



// Commit implements the Application interface
func (app *OKExChainApp) Commit() abci.ResponseCommit {

	//analyzer.OnCommitEnter()
	//defer analyzer.OnCommitExit()
	res := app.BaseApp.Commit()

	return res
}

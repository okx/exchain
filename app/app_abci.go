package app

import (
	"github.com/okex/exchain/x/analyzer"
	"github.com/okex/exchain/x/common/perf"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/trace"
)


// BeginBlock implements the Application interface
func (app *OKExChainApp) BeginBlock(req abci.RequestBeginBlock) (res abci.ResponseBeginBlock) {

	seq := perf.GetPerf().OnAppBeginBlockEnter(app.LastBlockHeight() + 1)
	analyzer.OnAppBeginBlockEnter(app.Logger(), app.LastBlockHeight()+1)
	defer perf.GetPerf().OnAppBeginBlockExit(app.LastBlockHeight()+1, seq)
	defer analyzer.OnAppBeginBlockExit()


	// dump app.LastBlockHeight()-1 info
	trace.GetElapsedInfo().Dump(app.Logger())
	return app.BaseApp.BeginBlock(req)
}

// EndBlock implements the Application interface
func (app *OKExChainApp) EndBlock(req abci.RequestEndBlock) (res abci.ResponseEndBlock) {

	seq := perf.GetPerf().OnAppEndBlockEnter(app.LastBlockHeight() + 1)
	analyzer.OnAppEndBlockEnter()
	defer perf.GetPerf().OnAppEndBlockExit(app.LastBlockHeight()+1, seq)
	defer analyzer.OnAppEndBlockExit()

	return app.BaseApp.EndBlock(req)
}


// Commit implements the Application interface
func (app *OKExChainApp) Commit() abci.ResponseCommit {

	seq := perf.GetPerf().OnCommitEnter(app.LastBlockHeight() + 1)
	analyzer.OnCommitEnter()
	defer perf.GetPerf().OnCommitExit(app.LastBlockHeight()+1, seq, app.Logger())
	defer analyzer.OnCommitExit()
	res := app.BaseApp.Commit()

	return res
}

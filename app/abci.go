package app

import (
	"fmt"

	"github.com/okex/okexchain/x/common/perf"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okexchain/app/protocol"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// DeliverTx implements the Application interface
func (app *OKChainApp) DeliverTx(req abci.RequestDeliverTx) (res abci.ResponseDeliverTx) {
	protocol.GetEngine().GetCurrentProtocol().CheckStopped()

	resp := app.BaseApp.DeliverTx(req)
	if (protocol.GetEngine().GetCurrentProtocol().GetBackendKeeper().Config.EnableBackend ||
		protocol.GetEngine().GetCurrentProtocol().GetStreamKeeper().AnalysisEnable()) && resp.IsOK() {
		app.syncTx(req.Tx)
	}

	return resp
}

// InitChain implements the Application interface
func (app *OKChainApp) InitChain(req abci.RequestInitChain) (res abci.ResponseInitChain) {

	app.log("[ABCI interface] ---> InitChain")
	return app.BaseApp.InitChain(req)
}

// BeginBlock implements the Application interface
func (app *OKChainApp) BeginBlock(req abci.RequestBeginBlock) (res abci.ResponseBeginBlock) {

	protocol.GetEngine().GetCurrentProtocol().CheckStopped()

	app.log("--------- Block[%d], Protocol[v%d] ---------", app.LastBlockHeight()+1,
		protocol.GetEngine().GetCurrentProtocol().GetVersion())

	seq := perf.GetPerf().OnAppBeginBlockEnter(app.LastBlockHeight() + 1)
	defer perf.GetPerf().OnAppBeginBlockExit(app.LastBlockHeight()+1, seq)

	return app.BaseApp.BeginBlock(req)
}

// EndBlock implements the Application interface
func (app *OKChainApp) EndBlock(req abci.RequestEndBlock) (res abci.ResponseEndBlock) {
	protocol.GetEngine().GetCurrentProtocol().CheckStopped()

	seq := perf.GetPerf().OnAppEndBlockEnter(app.LastBlockHeight() + 1)
	defer perf.GetPerf().OnAppEndBlockExit(app.LastBlockHeight()+1, seq)

	return app.BaseApp.EndBlock(req)
}

// Commit implements the Application interface
func (app *OKChainApp) Commit() abci.ResponseCommit {
	protocol.GetEngine().GetCurrentProtocol().CheckStopped()

	seq := perf.GetPerf().OnCommitEnter(app.LastBlockHeight() + 1)
	defer perf.GetPerf().OnCommitExit(app.LastBlockHeight()+1, seq, app.Logger())

	res := app.BaseApp.Commit()
	return res
}

// sync txBytes to backend module
func (app *OKChainApp) syncTx(txBytes []byte) {
	if tx, err := auth.DefaultTxDecoder(protocol.GetEngine().GetCurrentProtocol().GetCodec())(txBytes); err == nil {
		if stdTx, ok := tx.(auth.StdTx); ok {
			txHash := fmt.Sprintf("%X", tmhash.Sum(txBytes))
			app.Logger().Debug(fmt.Sprintf("[Sync Tx(%s) to backend module]", txHash))
			ctx := app.GetState(baseapp.RunTxModeDeliver()).Context()
			protocol.GetEngine().GetCurrentProtocol().GetBackendKeeper().SyncTx(ctx, &stdTx, txHash,
				ctx.BlockHeader().Time.Unix())
			protocol.GetEngine().GetCurrentProtocol().GetStreamKeeper().SyncTx(ctx, &stdTx, txHash,
				ctx.BlockHeader().Time.Unix())
		}
	}
}

// log format
func (app *OKChainApp) log(format string, a ...interface{}) {
	app.Logger().Info(fmt.Sprintf(format, a...))
}

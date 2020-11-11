package app

import (
	"fmt"
	"github.com/okex/okexchain/app/protocol"
	"github.com/okex/okexchain/x/common/perf"
	"github.com/okex/okexchain/x/common/version"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/x/auth"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

const appName = "OKExChainApp"

// OKExChainApp extends BaseApp(ABCI application)
type OKExChainApp struct {
	*baseapp.BaseApp
}

// NewOKExChainApp returns a reference to an initialized OKExChainApp.
func NewOKExChainApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp)) *OKExChainApp {
	bApp := baseapp.NewBaseApp(appName, logger, db, nil, baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	// set app version
	bApp.SetAppVersion(version.Version)
	// set protocol version
	bApp.ProtocolVersion = int32(version.CurrentProtocolVersion)

	app := &OKExChainApp{
		BaseApp: bApp,
	}
	// set hook function postEndBlocker
	bApp.PostEndBlocker = app.postEndBlocker

	// add new protocol based on new version
	if err := protocol.GetEngine().FillProtocol(app, logger, 0); err != nil {
		panic(err)
	}

	// recover the main store
	app.recoverLocalEnv(loadLatest)

	// load the status of current protocol from the store
	isLoaded, current := protocol.GetEngine().LoadCurrentProtocol(app.GetCommitMultiStore().GetKVStore(
		protocol.GetMainStoreKey()))
	if !isLoaded {
		cmn.Exit(fmt.Sprintf("Your software doesn't support the required protocol (version %d)!", current))
	}
	logger.Debug(fmt.Sprintf("launch app with version: %v", current))

	// set txDecoder for BaseApp
	app.SetTxDecoder(auth.DefaultTxDecoder(protocol.GetEngine().GetCurrentProtocol().GetCodec()))

	// enable perf
	perf.GetPerf().EnableCheck()

	return app
}

// LoadHeight loads data on a particular height
func (app *OKExChainApp) LoadHeight(height int64) error {
	//return app.LoadVersion(height, app.keys[bam.MainStoreKey])
	return app.LoadVersion(height, protocol.GetMainStoreKey())
}

// hook function for BaseApp's EndBlock(upgrade)
func (app *OKExChainApp) postEndBlocker(res *abci.ResponseEndBlock) {

}

func (app *OKExChainApp) recoverLocalEnv(loadLatest bool) {
	// the current field in AppProtocolEngine is 0
	// on the beginning for the running of NewOKExChainApp()

	// it will mount protocolv0.GetKVStoreKeysMap()
	app.MountKVStores(protocol.GetKVStoreKeysMap())
	// it will mount protocolv0.GetTransientStoreKeysMap()
	app.MountTransientStores(protocol.GetTransientStoreKeysMap())
	if loadLatest {
		if err := app.LoadLatestVersion(protocol.GetMainStoreKey()); err != nil {
			cmn.Exit(err.Error())
		}
	} else {
		if err := app.GetCommitMultiStore().LoadVersion(0); err != nil {
			cmn.Exit(err.Error())
		}
	}
}

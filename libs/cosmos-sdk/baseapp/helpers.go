package baseapp

import (
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"regexp"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var isAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func (app *BaseApp) Check(tx sdk.Tx) (sdk.GasInfo, *sdk.Result, error) {
	info, e := app.runTx(runTxModeCheck, nil, tx, LatestSimulateTxHeight)
	return info.gInfo, info.result, e
}

func (app *BaseApp) Simulate(txBytes []byte, tx sdk.Tx, height int64, overridesBytes []byte, from ...string) (sdk.GasInfo, *sdk.Result, error) {
	info := &runTxInfo{
		overridesBytes: overridesBytes,
	}
	e := app.runtxWithInfo(info, runTxModeSimulate, txBytes, tx, height, from...)
	return info.gInfo, info.result, e
}

func (app *BaseApp) MempoolSimulate(txBytes []byte, tx sdk.Tx, height int64, overridesBytes []byte, from ...string) (sdk.GasInfo, *sdk.Result, error) {
	info := &runTxInfo{
		overridesBytes:  overridesBytes,
		mempoolSimulate: true,
	}
	if cfg.DynamicConfig.GetEnableMempoolSimGuFactor() {
		info.mempoolSimulate = false // the mempoolSimulate is false is need gu factor
	}
	e := app.runtxWithInfo(info, runTxModeSimulate, txBytes, tx, height, from...)
	return info.gInfo, info.result, e
}

func (app *BaseApp) Deliver(tx sdk.Tx) (sdk.GasInfo, *sdk.Result, error) {
	info, e := app.runTx(runTxModeDeliver, nil, tx, LatestSimulateTxHeight)
	return info.gInfo, info.result, e
}

// Context with current {check, deliver}State of the app used by tests.
func (app *BaseApp) NewContext(isCheckTx bool, header abci.Header) (ctx sdk.Context) {
	if isCheckTx {
		ctx = sdk.NewContext(app.checkState.ms, header, true, app.logger)
		ctx.SetMinGasPrices(app.minGasPrices)
		return
	}

	ctx = sdk.NewContext(app.deliverState.ms, header, false, app.logger)
	return
}

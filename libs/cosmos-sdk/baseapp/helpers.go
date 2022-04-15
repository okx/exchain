package baseapp

import (
	"regexp"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var isAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

func (app *BaseApp) Check(tx sdk.Tx) (sdk.GasInfo, *sdk.Result, error) {
	info, e := app.runTx(runTxModeCheck, nil, tx, LatestSimulateTxHeight)
	return info.gInfo, info.result, e
}

func (app *BaseApp) Simulate(txBytes []byte, tx sdk.Tx, height int64, from ...string) (sdk.GasInfo, *sdk.Result, error) {
	info, e := app.runTx(runTxModeSimulate, txBytes, tx, height, from...)
	return info.gInfo, info.result, e
}

func (app *BaseApp) Deliver(tx sdk.Tx) (sdk.GasInfo, *sdk.Result, error) {
	info, e := app.runTx(runTxModeDeliver, nil, tx, LatestSimulateTxHeight)
	return info.gInfo, info.result, e
}

// Context with current {check, deliver}State of the app used by tests.
func (app *BaseApp) NewContext(isCheckTx bool, header abci.Header) sdk.Context {
	if isCheckTx {
		return sdk.NewContext(app.checkState.ms, header, true, app.logger).
			WithMinGasPrices(app.minGasPrices)
	}

	return sdk.NewContext(app.deliverState.ms, header, false, app.logger)
}

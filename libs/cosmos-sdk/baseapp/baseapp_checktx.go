package baseapp

import (
	"encoding/json"
	"fmt"
	"sync/atomic"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/global"
)

// CheckTx implements the ABCI interface and executes a tx in CheckTx mode. In
// CheckTx mode, messages are not executed. This means messages are only validated
// and only the AnteHandler is executed. State is persisted to the BaseApp's
// internal CheckTx state if the AnteHandler passes. Otherwise, the ResponseCheckTx
// will contain releveant error information. Regardless of tx execution outcome,
// the ResponseCheckTx will contain relevant gas execution context.
func (app *BaseApp) CheckTx(req abci.RequestCheckTx) abci.ResponseCheckTx {
	tx, err := app.txDecoder(req.Tx, global.GetGlobalHeight())
	if err != nil {
		return sdkerrors.ResponseCheckTx(err, 0, 0, app.trace)
	}

	var mode runTxMode

	switch {
	case req.Type == abci.CheckTxType_New:
		mode = runTxModeCheck
		atomic.AddInt64(&app.checkTxNum, 1)
	case req.Type == abci.CheckTxType_Recheck:
		mode = runTxModeReCheck

	case req.Type == abci.CheckTxType_WrappedCheck:
		mode = runTxModeWrappedCheck
		atomic.AddInt64(&app.wrappedCheckTxNum, 1)

	default:
		panic(fmt.Sprintf("unknown RequestCheckTx type: %s", req.Type))
	}

	if abci.GetDisableCheckTx() {
		var ctx sdk.Context
		ctx = app.getContextForTx(mode, req.Tx)
		exTxInfo := app.GetTxInfo(ctx, tx)
		data, _ := json.Marshal(exTxInfo)

		return abci.ResponseCheckTx{
			Data: data,
		}
	}

	gInfo, result, _, err := app.runTx(mode, req.Tx, tx, LatestSimulateTxHeight, req.From)
	if err != nil {
		return sdkerrors.ResponseCheckTx(err, gInfo.GasWanted, gInfo.GasUsed, app.trace)
	}

	return abci.ResponseCheckTx{
		GasWanted: int64(gInfo.GasWanted), // TODO: Should type accept unsigned ints?
		GasUsed:   int64(gInfo.GasUsed),   // TODO: Should type accept unsigned ints?
		Log:       result.Log,
		Data:      result.Data,
		Events:    result.Events.ToABCIEvents(),
	}
}

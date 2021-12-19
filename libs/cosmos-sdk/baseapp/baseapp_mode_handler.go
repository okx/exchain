package baseapp

import (
	"fmt"
	"encoding/json"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

type modeHandler interface {
	getMode() runTxMode
	handleStartHeight(info *runTxInfo, height int64) error
	handleGasConsumed(info *runTxInfo)(startingGas uint64, gInfo sdk.GasInfo, err error)
	handleRunMsg(info *runTxInfo) (result *sdk.Result, err error)
}

func (app *BaseApp) getModeHandler(mode runTxMode) modeHandler {
	var h modeHandler
	switch mode {
	case runTxModeCheck:
		h = &modeHandlerCheck {&modeHandlerBase {mode: mode, app: app,}}
	case runTxModeReCheck:
		h = &modeHandlerRecheck {&modeHandlerBase {mode: mode, app: app,}}
	case runTxModeDeliver:
		h = &modeHandlerDeliver {&modeHandlerBase {mode: mode, app: app,}}
	case runTxModeSimulate:
		h = &modeHandlerSimulate {&modeHandlerBase {mode: mode, app: app,}}
	case runTxModeDeliverInAsync:
		h = &modeHandlerDeliverInAsync {&modeHandlerBase {mode: mode, app: app,}}
	default:
		h = &modeHandlerBase {mode: mode, app: app,}
	}

	return h
}

type modeHandlerBase struct {
	mode runTxMode
	app *BaseApp
}


type modeHandlerDeliverInAsync struct {
	*modeHandlerBase
}

type modeHandlerDeliver struct {
	*modeHandlerBase
}
type modeHandlerCheck struct {
	*modeHandlerBase
}

type modeHandlerRecheck struct {
	*modeHandlerBase
}

type modeHandlerSimulate struct {
	*modeHandlerBase
}

func (m *modeHandlerBase) getMode() runTxMode {
	return m.mode
}

// ====================================================
// handleStartHeight
func (m *modeHandlerSimulate) handleStartHeight(info *runTxInfo, height int64) error {
	app := m.app
	startHeight := tmtypes.GetStartBlockHeight()

	var err error
	if height > startHeight && height < app.LastBlockHeight() {
		info.ctx, err = app.getContextForSimTx(info.txBytes, height)
	}

	return err
}

func (m *modeHandlerBase) handleStartHeight(info *runTxInfo, height int64) error {
	app := m.app
	startHeight := tmtypes.GetStartBlockHeight()

	if height < startHeight && height != 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("height(%d) should be greater than start block height(%d)", height, startHeight))
	} else {
		info.ctx = app.getContextForTx(m.mode, info.txBytes)
	}

	return nil
}

// ====================================================
// handleGasConsumed
func (m *modeHandlerBase) handleGasConsumed(info *runTxInfo) (startingGas uint64, gInfo sdk.GasInfo, err error) {

	if info.ctx.BlockGasMeter().IsOutOfGas() {
		gInfo = sdk.GasInfo{GasUsed: info.ctx.BlockGasMeter().GasConsumed()}
		return startingGas, gInfo, sdkerrors.Wrap(sdkerrors.ErrOutOfGas, "no block gas left to run tx")
	}
	startingGas = info.ctx.BlockGasMeter().GasConsumed()

	return startingGas, gInfo, nil
}

// noop
func (m *modeHandlerRecheck) handleGasConsumed(*runTxInfo) (startingGas uint64, gInfo sdk.GasInfo, err error){return}
func (m *modeHandlerCheck) handleGasConsumed(*runTxInfo) (startingGas uint64, gInfo sdk.GasInfo, err error){return}
func (m *modeHandlerSimulate) handleGasConsumed(*runTxInfo) (startingGas uint64, gInfo sdk.GasInfo, err error){return}

//==========================================================================
// handleRunMsg
func (m *modeHandlerBase) handleRunMsg(info *runTxInfo) (result *sdk.Result, err error) {
	app := m.app
	mode := m.mode
	msCacheAnte := info.msCacheAnte
	msCache := info.msCache

	if mode == runTxModeDeliverInAsync {
		info.msCache = msCacheAnte.CacheMultiStore()
		info.runMsgCtx = info.ctx.WithMultiStore(msCache)
	} else {
		info.runMsgCtx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)
	}
	msCache = info.msCache

	result, err = app.runMsgs(info.runMsgCtx, info.tx.GetMsgs(), mode)
	if err == nil && (mode == runTxModeDeliver) {
		msCache.Write()
	}

	info.runMsgFinished = true

	if mode == runTxModeCheck {
		exTxInfo := app.GetTxInfo(info.ctx, info.tx)
		exTxInfo.SenderNonce = info.accountNonce

		data, err := json.Marshal(exTxInfo)
		if err == nil {
			result.Data = data
		}
	}

	if err != nil {
		if sdk.HigherThanMercury(info.ctx.BlockHeight()) {
			codeSpace, code, info := sdkerrors.ABCIInfo(err, app.trace)
			err = sdkerrors.New(codeSpace, abci.CodeTypeNonceInc+code, info)
		}
		msCache = nil
	}

	if mode == runTxModeDeliverInAsync {
		if msCache != nil {
			msCache.Write()
		}
	}

	return
}

//=============================




//func (m *modeHandlerBase) handleGasConsumed3(height int64, txBytes []byte) error {
//	app := m.app
//
//
//	return nil
//}
//=============================
package baseapp

import (
	"encoding/json"
	"fmt"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"

	//"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

type modeHandler interface {
	getMode() runTxMode

	handleStartHeight(info *runTxInfo, height int64) error
	handleGasConsumed(info *runTxInfo) error
	handleRunMsg(info *runTxInfo) error
	handleDeferRefund(info *runTxInfo)
	handleDeferGasConsumed(info *runTxInfo)
}

func (app *BaseApp) getModeHandler(mode runTxMode) modeHandler {
	var h modeHandler
	switch mode {
	case runTxModeCheck, runTxModeWrappedCheck:
		h = &modeHandlerCheck{&modeHandlerBase{mode: mode, app: app}}
	case runTxModeReCheck:
		h = &modeHandlerRecheck{&modeHandlerBase{mode: mode, app: app}}
	case runTxModeTrace:
		h = &modeHandlerTrace{&modeHandlerDeliver{&modeHandlerBase{mode: mode, app: app}}}
	case runTxModeDeliver:
		h = &modeHandlerDeliver{&modeHandlerBase{mode: mode, app: app}}
	case runTxModeSimulate:
		h = &modeHandlerSimulate{&modeHandlerBase{mode: mode, app: app}}
	case runTxModeDeliverInAsync:
		h = &modeHandlerDeliverInAsync{&modeHandlerBase{mode: mode, app: app}}
	case runTxModeDeliverPartConcurrent:
		h = &modeHandlerDeliverPartConcurrent{&modeHandlerBase{mode: mode,app: app}}
	default:
		h = &modeHandlerBase{mode: mode, app: app}
	}

	return h
}

type modeHandlerBase struct {
	mode runTxMode
	app  *BaseApp
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

//modeHandlerTrace derived from modeHandlerDeliver
type modeHandlerTrace struct {
	*modeHandlerDeliver
}

type modeHandlerDeliverPartConcurrent struct {
	*modeHandlerBase
}

func (m *modeHandlerBase) getMode() runTxMode {
	return m.mode
}

// ====================================================
// 1. handleStartHeight
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
// 2. handleGasConsumed
func (m *modeHandlerBase) handleGasConsumed(info *runTxInfo) (err error) {

	if info.ctx.BlockGasMeter().IsOutOfGas() {
		info.gInfo = sdk.GasInfo{GasUsed: info.ctx.BlockGasMeter().GasConsumed()}
		err = sdkerrors.Wrap(sdkerrors.ErrOutOfGas, "no block gas left to run tx")
	} else {
		info.startingGas = info.ctx.BlockGasMeter().GasConsumed()
	}

	return err
}

// noop
func (m *modeHandlerRecheck) handleGasConsumed(*runTxInfo) (err error)  { return }
func (m *modeHandlerCheck) handleGasConsumed(*runTxInfo) (err error)    { return }
func (m *modeHandlerSimulate) handleGasConsumed(*runTxInfo) (err error) { return }

//==========================================================================
// 3. handleRunMsg

// modeHandlerBase.handleRunMsg derived by:
// (m *modeHandlerRecheck)
// (m *modeHandlerCheck)
// (m *modeHandlerSimulate)
func (m *modeHandlerBase) handleRunMsg(info *runTxInfo) (err error) {
	app := m.app
	mode := m.mode

	info.runMsgCtx, info.msCache = app.cacheTxContext(info.ctx, info.txBytes)
	info.result, err = app.runMsgs(info.runMsgCtx, info.tx.GetMsgs(), mode)
	info.runMsgFinished = true

	m.handleRunMsg4CheckMode(info)
	err = m.checkHigherThanMercury(err, info)
	return
}

//=============================
// 4. handleDeferGasConsumed
func (m *modeHandlerBase) handleDeferGasConsumed(*runTxInfo) {}

//====================================================================
// 5. handleDeferRefund
func (m *modeHandlerBase) handleDeferRefund(*runTxInfo) {}

//===========================================================================================
// other members
func (m *modeHandlerBase) setGasConsumed(info *runTxInfo) {
	info.ctx.BlockGasMeter().ConsumeGas(info.ctx.GasMeter().GasConsumedToLimit(), "block gas meter")
	if info.ctx.BlockGasMeter().GasConsumed() < info.startingGas {
		panic(sdk.ErrorGasOverflow{Descriptor: "tx gas summation"})
	}
}

func (m *modeHandlerBase) checkHigherThanMercury(err error, info *runTxInfo) error {

	if err != nil {
		if tmtypes.HigherThanMercury(info.ctx.BlockHeight()) {
			codeSpace, code, info := sdkerrors.ABCIInfo(err, m.app.trace)
			err = sdkerrors.New(codeSpace, abci.CodeTypeNonceInc+code, info)
		}
		info.msCache = nil
	}
	return err
}

func (m *modeHandlerBase) handleRunMsg4CheckMode(info *runTxInfo) {
	if m.mode != runTxModeCheck && m.mode != runTxModeWrappedCheck {
		return
	}

	exTxInfo := m.app.GetTxInfo(info.ctx, info.tx)
	exTxInfo.SenderNonce = info.accountNonce

	data, err := json.Marshal(exTxInfo)
	if err == nil {
		info.result.Data = data
	}
}

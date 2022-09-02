package baseapp

import (
	"fmt"

	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

func (m *modeHandlerSimulate) handleStartHeight(info *runTxInfo, height int64) error {
	app := m.app
	startHeight := tmtypes.GetStartBlockHeight()

	var err error
	lastHeight := app.LastBlockHeight()
	if height > startHeight && height <= lastHeight {
		info.ctx, err = app.getContextForSimTx(info.txBytes, height)
	} else if height < startHeight && height != 0 {
		err = sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("height(%d) should be greater than start block height(%d)", height, startHeight))
	} else if height > lastHeight {
		err = sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("height(%d) should be less than or equal to latest block height(%d)", height, lastHeight))
	} else {
		info.ctx = app.getContextForTx(m.mode, info.txBytes)
	}
	if info.overridesBytes != nil {
		info.ctx.SetOverrideBytes(info.overridesBytes)
	}
	return err
}

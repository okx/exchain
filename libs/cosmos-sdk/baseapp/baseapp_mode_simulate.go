package baseapp

import (
	"fmt"

	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	tmtypes "github.com/okx/okbchain/libs/tendermint/types"
)

func (m *modeHandlerSimulate) handleStartHeight(info *runTxInfo, height int64) error {
	app := m.app
	startHeight := tmtypes.GetStartBlockHeight()

	var err error
	lastHeight := app.LastBlockHeight()
	if height == 0 {
		height = lastHeight
	}
	if height <= startHeight {
		err = sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("height(%d) should be greater than start block height(%d)", height, startHeight))
	} else if height > startHeight && height <= lastHeight {
		info.ctx, err = app.getContextForSimTx(info.txBytes, height)
	} else {
		err = sdkerrors.Wrap(sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("height(%d) should be less than or equal to latest block height(%d)", height, lastHeight))
	}
	if info.overridesBytes != nil {
		info.ctx.SetOverrideBytes(info.overridesBytes)
	}
	return err
}

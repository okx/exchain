package ante

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	evmtypes "github.com/okx/okbchain/x/evm/types"
)

// AccountBlockedVerificationDecorator check whether signer is blocked.
type AccountBlockedVerificationDecorator struct {
	evmKeeper EVMKeeper
}

// NewAccountBlockedVerificationDecorator creates a new AccountBlockedVerificationDecorator instance
func NewAccountBlockedVerificationDecorator(evmKeeper EVMKeeper) AccountBlockedVerificationDecorator {
	return AccountBlockedVerificationDecorator{
		evmKeeper: evmKeeper,
	}
}

// AnteHandle check wether signer of tx(contains cosmos-tx and eth-tx) is blocked.
func (abvd AccountBlockedVerificationDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// simulate means 'eth_call' or 'eth_estimateGas', when it means 'eth_estimateGas' we can not 'VerifySig'.so skip here
	if simulate {
		return next(ctx, tx, simulate)
	}
	pinAnte(ctx.AnteTracer(), "AccountBlockedVerificationDecorator")

	var signers []sdk.AccAddress
	if ethTx, ok := tx.(*evmtypes.MsgEthereumTx); ok {
		signers = ethTx.GetSigners()
	} else {
		signers = tx.GetSigners()
	}

	currentGasMeter := ctx.GasMeter()
	infGasMeter := sdk.GetReusableInfiniteGasMeter()
	ctx.SetGasMeter(infGasMeter)

	for _, signer := range signers {
		//TODO it may be optimizate by cache blockedAddressList
		if ok := abvd.evmKeeper.IsAddressBlocked(ctx, signer); ok {
			ctx.SetGasMeter(currentGasMeter)
			sdk.ReturnInfiniteGasMeter(infGasMeter)
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "address: %s has been blocked", signer.String())
		}
	}
	ctx.SetGasMeter(currentGasMeter)
	sdk.ReturnInfiniteGasMeter(infGasMeter)
	return next(ctx, tx, simulate)
}

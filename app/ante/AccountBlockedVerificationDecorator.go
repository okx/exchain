package ante

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"sync"
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

var resuableGasMeterPool = &sync.Pool{
	New: func() interface{} {
		return sdk.NewReusableInfiniteGasMeter()
	},
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
		signers = []sdk.AccAddress{ethTx.AccountAddress()}
	} else {
		signers = tx.GetSigners()
	}

	currentGasMeter := ctx.GasMeter()
	infGasMeter := resuableGasMeterPool.Get().(sdk.ReusableGasMeter)
	infGasMeter.Reset()
	ctx.SetGasMeter(infGasMeter)

	for _, signer := range signers {
		//TODO it may be optimizate by cache blockedAddressList
		if ok := abvd.evmKeeper.IsAddressBlocked(ctx, signer); ok {
			ctx.SetGasMeter(currentGasMeter)
			resuableGasMeterPool.Put(infGasMeter)
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "address: %s has been blocked", signer.String())
		}
	}
	ctx.SetGasMeter(currentGasMeter)
	resuableGasMeterPool.Put(infGasMeter)
	return next(ctx, tx, simulate)
}

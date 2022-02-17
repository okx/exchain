package ante

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authante "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ante"
	evmtypes "github.com/okex/exchain/x/evm/types"
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
	signers, err := getSigners(tx)
	if err != nil {
		return ctx, err
	}
	currentGasMeter := ctx.GasMeter()
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	for _, signer := range signers {
		//TODO it may be optimizate by cache blockedAddressList
		if ok := abvd.evmKeeper.IsAddressBlocked(ctx, signer); ok {
			ctx = ctx.WithGasMeter(currentGasMeter)
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "address: %s has been blocked", signer.String())
		}
	}
	ctx = ctx.WithGasMeter(currentGasMeter)
	return next(ctx, tx, simulate)
}

// getSigners get signers of tx(contains cosmos-tx and eth-tx.
func getSigners(tx sdk.Tx) ([]sdk.AccAddress, error) {
	signers := make([]sdk.AccAddress, 0)
	switch tx.(type) {
	case auth.StdTx:
		sigTx, ok := tx.(authante.SigVerifiableTx)
		if !ok {
			return signers, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "invalid transaction type")
		}
		signers = append(signers, sigTx.GetSigners()...)
	case evmtypes.MsgEthereumTx:
		msgEthTx, ok := tx.(evmtypes.MsgEthereumTx)
		if !ok {
			return signers, sdkerrors.Wrapf(sdkerrors.ErrTxDecode, "invalid transaction type: %T", tx)
		}
		signers = append(signers, msgEthTx.GetSigners()...)

	default:
		return signers, sdkerrors.Wrapf(sdkerrors.ErrTxDecode, "invalid transaction type: %T", tx)
	}
	return signers, nil
}

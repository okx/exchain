package ante

import (
	"fmt"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// CheckedMempoolFeeDecorator used for checked tx to check mempool fee
type CheckedMempoolFeeDecorator struct {
	evmkeeper EVMKeeper
}

// NewCheckedMempoolFeeDecorator create a new CheckedMempoolFeeDecorator
func NewCheckedMempoolFeeDecorator(e EVMKeeper) CheckedMempoolFeeDecorator {
	return CheckedMempoolFeeDecorator{
		evmkeeper: e,
	}
}

// nolint
// the Ante logic for this decorator
// a copy from EthMempoolFeeDecorator adjust for the CheckedTx and OriginTx
func (cmfd CheckedMempoolFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if !ctx.IsCheckTx() {
		return next(ctx, tx, simulate)
	}

	var minFees sdk.Dec // after type asset this struct is non-empty
	var fee sdk.DecCoin
	minGasPrices := ctx.MinGasPrices()
	evmDenom := sdk.DefaultBondDenom
	if len(tx.GetTxCarriedData()) > 0 {
		if t, ok := tx.(evmtypes.MsgEthereumCheckedTx); ok {
			// TODO: add the FeeTx for the MsgEthereumChedkedTx to conver the fee
			minFees = minGasPrices.AmountOf(evmDenom).MulInt64(int64(t.Data.GasLimit))
		} else {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
		}
	} else {

	}

	var hasEnoughFees bool
	if fee.Amount.GTE(minFees) {
		hasEnoughFees = true
	}
	if !ctx.MinGasPrices().IsZero() && !hasEnoughFees {
		return ctx, sdkerrors.Wrap(
			sdkerrors.ErrInsufficientFee,
			fmt.Sprintf("insufficient fee, got: %q required: %q", fee, sdk.NewDecCoinFromDec(evmDenom, minFees)),
		)
	}

	return
}

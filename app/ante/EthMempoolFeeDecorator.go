package ante

import (
	"fmt"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	evmtypes "github.com/okx/okbchain/x/evm/types"
	"math/big"
	"sync"
)

// EthMempoolFeeDecorator validates that sufficient fees have been provided that
// meet a minimum threshold defined by the proposer (for mempool purposes during CheckTx).
type EthMempoolFeeDecorator struct {
	evmKeeper EVMKeeper
}

// NewEthMempoolFeeDecorator creates a new EthMempoolFeeDecorator
func NewEthMempoolFeeDecorator(ek EVMKeeper) EthMempoolFeeDecorator {
	return EthMempoolFeeDecorator{
		evmKeeper: ek,
	}
}

var feeIntsPool = &sync.Pool{
	New: func() interface{} {
		return &[2]big.Int{}
	},
}

// AnteHandle verifies that enough fees have been provided by the
// Ethereum transaction that meet the minimum threshold set by the block
// proposer.
//
// NOTE: This should only be run during a CheckTx mode.
func (emfd EthMempoolFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {

	// simulate means 'eth_call' or 'eth_estimateGas', when it means 'eth_estimateGas' we can not 'VerifySig'.so skip here
	if !ctx.IsCheckTx() || simulate {
		return next(ctx, tx, simulate)
	}

	msgEthTx, ok := tx.(*evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
	}

	const evmDenom = sdk.DefaultBondDenom

	feeInts := feeIntsPool.Get().(*[2]big.Int)

	// fee = gas price * gas limit
	fee := sdk.NewDecCoinFromDec(evmDenom, sdk.NewDecWithBigIntAndPrec(msgEthTx.CalcFee(&feeInts[0]), sdk.Precision))

	minGasPrices := ctx.MinGasPrices()
	// minFees := minGasPrices.AmountOf(evmDenom).MulInt64(int64(msgEthTx.Data.GasLimit))
	var minFees = sdk.Dec{&feeInts[1]}
	minGasPrices.AmountOf(evmDenom).MulInt64To(int64(msgEthTx.Data.GasLimit), &minFees)

	// check that fee provided is greater than the minimum defined by the validator node
	// NOTE: we only check if the evm denom tokens are present in min gas prices. It is up to the
	// sender if they want to send additional fees in other denominations.
	var hasEnoughFees bool
	if fee.Amount.GTE(minFees) {
		hasEnoughFees = true
	}

	// reject transaction if minimum gas price is not zero and the transaction does not
	// meet the minimum fee
	if !ctx.MinGasPrices().IsZero() && !hasEnoughFees {
		err = sdkerrors.Wrap(
			sdkerrors.ErrInsufficientFee,
			fmt.Sprintf("insufficient fee, got: %q required: %q", fee, sdk.NewDecCoinFromDec(evmDenom, minFees)),
		)
		feeIntsPool.Put(feeInts)
		return ctx, err
	}
	feeIntsPool.Put(feeInts)

	return next(ctx, tx, simulate)
}

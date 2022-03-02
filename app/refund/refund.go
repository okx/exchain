package refund

import (
	"math/big"

	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/ante"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/keeper"

	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/refund"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

func NewGasRefundHandler(ak auth.AccountKeeper, sk types.SupplyKeeper) sdk.GasRefundHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx,
	) (refundFee sdk.Coins, err error) {
		var gasRefundHandler sdk.GasRefundHandler

		if tx.GetType() == sdk.EvmTxType {
			gasRefundHandler = NewGasRefundDecorator(ak, sk)
		} else {
			return nil, nil
		}
		return gasRefundHandler(ctx, tx)
	}
}

type Handler struct {
	ak           keeper.AccountKeeper
	supplyKeeper types.SupplyKeeper
}

func (handler Handler) GasRefund(ctx sdk.Context, tx sdk.Tx) (refundGasFee sdk.Coins, err error) {

	currentGasMeter := ctx.GasMeter()
	TempGasMeter := sdk.NewInfiniteGasMeter()
	ctx = ctx.WithGasMeter(TempGasMeter)

	defer func() {
		ctx = ctx.WithGasMeter(currentGasMeter)
	}()

	gasLimit := currentGasMeter.Limit()
	gasUsed := currentGasMeter.GasConsumed()

	if gasUsed >= gasLimit {
		return nil, nil
	}

	feeTx, ok := tx.(ante.FeeTx)
	if !ok {
		return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer(ctx)
	feePayerAcc := handler.ak.GetAccount(ctx, feePayer)
	if feePayerAcc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	gas := feeTx.GetGas()
	fees := feeTx.GetFee()
	gasFees := caculateRefundFees(gasUsed, gas, fees)
	//fmt.Println("refund-", ethcommon.BytesToAddress(feePayerAcc.GetAddress()), gasUsed)
	err = refund.RefundFees(handler.supplyKeeper, ctx, feePayerAcc.GetAddress(), gasFees)
	if err != nil {
		return nil, err
	}

	return gasFees, nil
}

func NewGasRefundDecorator(ak auth.AccountKeeper, sk types.SupplyKeeper) sdk.GasRefundHandler {
	chandler := Handler{
		ak:           ak,
		supplyKeeper: sk,
	}

	return func(ctx sdk.Context, tx sdk.Tx) (refund sdk.Coins, err error) {
		return chandler.GasRefund(ctx, tx)
	}
}

func caculateRefundFees(gasUsed uint64, gas uint64, fees sdk.DecCoins) sdk.Coins {

	refundFees := make(sdk.Coins, len(fees))
	for i, fee := range fees {
		gasPrice := new(big.Int).Div(fee.Amount.BigInt(), new(big.Int).SetUint64(gas))
		gasConsumed := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gasUsed))
		gasCost := sdk.NewCoin(fee.Denom, sdk.NewDecFromBigIntWithPrec(gasConsumed, sdk.Precision))
		gasRefund := fee.Sub(gasCost)

		refundFees[i] = gasRefund
	}
	return refundFees
}

// CaculateRefundFees provides the way to calculate the refunded gas with gasUsed, fees and gasPrice,
// as refunded gas = fees - gasPrice * gasUsed
func CaculateRefundFees(gasUsed uint64, fees sdk.DecCoins, gasPrice *big.Int) sdk.Coins {
	gas := new(big.Int).Div(fees[0].Amount.BigInt(), gasPrice).Uint64()
	return caculateRefundFees(gasUsed, gas, fees)
}

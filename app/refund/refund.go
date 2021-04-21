package refund

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"

	"github.com/cosmos/cosmos-sdk/x/auth/refund"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

func NewGasRefundHandler(ak auth.AccountKeeper, sk types.SupplyKeeper) sdk.GasRefundHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx,
	) (err error) {
		var gasRefundHandler sdk.GasRefundHandler
		switch tx.(type) {
		case evmtypes.MsgEthereumTx:
			gasRefundHandler = NewGasRefundDecorator(ak, sk)
		default:
			return nil
		}
		return gasRefundHandler(ctx, tx)
	}
}

type Handler struct {
	ak           keeper.AccountKeeper
	supplyKeeper types.SupplyKeeper
}

func (handler Handler) GasRefund(ctx sdk.Context, tx sdk.Tx) (err error) {

	currentGasMeter := ctx.GasMeter()
	TempGasMeter := sdk.NewInfiniteGasMeter()
	ctx = ctx.WithGasMeter(TempGasMeter)

	defer func() {
		ctx = ctx.WithGasMeter(currentGasMeter)
	}()

	gasLimit := currentGasMeter.Limit()
	gasUsed := currentGasMeter.GasConsumed()

	if gasUsed >= gasLimit {
		return nil
	}

	feeTx, ok := tx.(ante.FeeTx)
	if !ok {
		return sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feePayer := feeTx.FeePayer()
	feePayerAcc := handler.ak.GetAccount(ctx, feePayer)
	if feePayerAcc == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	gas := feeTx.GetGas()
	fees := feeTx.GetFee()
	gasFees := make(sdk.Coins, len(fees))

	for i, fee := range fees {
		gasPrice := new(big.Int).Div(fee.Amount.BigInt(), new(big.Int).SetUint64(gas))
		gasConsumed := new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gasUsed))
		gasCost := sdk.NewCoin(fee.Denom, sdk.NewDecFromBigIntWithPrec(gasConsumed, sdk.Precision))
		gasRefund := fee.Sub(gasCost)

		gasFees[i] = gasRefund
	}

	err = refund.RefundFees(handler.supplyKeeper, ctx, feePayerAcc.GetAddress(), gasFees)
	if err != nil {
		return err
	}

	return nil
}

func NewGasRefundDecorator(ak auth.AccountKeeper, sk types.SupplyKeeper) sdk.GasRefundHandler {
	chandler := Handler{
		ak:           ak,
		supplyKeeper: sk,
	}

	return func(ctx sdk.Context, tx sdk.Tx) (err error) {
		return chandler.GasRefund(ctx, tx)
	}
}

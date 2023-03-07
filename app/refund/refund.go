package refund

import (
	"math/big"
	"sync"

	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/ante"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/keeper"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/innertx"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/types"
)

func NewGasRefundHandler(ak auth.AccountKeeper, sk types.SupplyKeeper, ik innertx.InnerTxKeeper) sdk.GasRefundHandler {
	evmGasRefundHandler := NewGasRefundDecorator(ak, sk, ik)

	return func(
		ctx sdk.Context, tx sdk.Tx,
	) (refundFee sdk.Coins, err error) {
		var gasRefundHandler sdk.GasRefundHandler

		if tx.GetType() == sdk.EvmTxType {
			gasRefundHandler = evmGasRefundHandler
		} else {
			return nil, nil
		}
		return gasRefundHandler(ctx, tx)
	}
}

type Handler struct {
	ak           keeper.AccountKeeper
	supplyKeeper types.SupplyKeeper
	ik           innertx.InnerTxKeeper
}

func (handler Handler) GasRefund(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, error) {
	return gasRefund(handler.ik, handler.ak, handler.supplyKeeper, ctx, tx)
}

type accountKeeperInterface interface {
	SetAccount(ctx sdk.Context, acc exported.Account)
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account
}

func gasRefund(ik innertx.InnerTxKeeper, ak accountKeeperInterface, sk types.SupplyKeeper, ctx sdk.Context, tx sdk.Tx) (refundGasFee sdk.Coins, err error) {
	currentGasMeter := ctx.GasMeter()
	ctx.SetGasMeter(sdk.NewInfiniteGasMeter())

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
	feePayerAcc := ak.GetAccount(ctx, feePayer)
	if feePayerAcc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "fee payer address: %s does not exist", feePayer)
	}

	gas := feeTx.GetGas()
	fees := feeTx.GetFee()
	gasFees := calculateRefundFees(gasUsed, gas, fees)
	newCoins := feePayerAcc.GetCoins().Add(gasFees...)

	// set coins and record innertx
	err = feePayerAcc.SetCoins(newCoins)
	if !ctx.IsCheckTx() {
		fromAddr := sk.GetModuleAddress(types.FeeCollectorName)
		ik.UpdateInnerTx(ctx.TxBytes(), ctx.BlockHeight(), innertx.CosmosDepth, fromAddr, feePayerAcc.GetAddress(), innertx.CosmosCallType, innertx.SendCallName, gasFees, err)
	}
	if err != nil {
		return nil, err
	}
	ak.SetAccount(ctx, feePayerAcc)

	return gasFees, nil
}

func NewGasRefundDecorator(ak auth.AccountKeeper, sk types.SupplyKeeper, ik innertx.InnerTxKeeper) sdk.GasRefundHandler {
	chandler := Handler{
		ak:           ak,
		supplyKeeper: sk,
		ik:           ik,
	}
	return chandler.GasRefund
}

var bigIntsPool = &sync.Pool{
	New: func() interface{} {
		return &[2]big.Int{}
	},
}

func calculateRefundFees(gasUsed uint64, gas uint64, fees sdk.DecCoins) sdk.Coins {
	bitInts := bigIntsPool.Get().(*[2]big.Int)
	defer bigIntsPool.Put(bitInts)

	refundFees := make(sdk.Coins, len(fees))
	for i, fee := range fees {
		gasPrice := bitInts[0].SetUint64(gas)
		gasPrice = gasPrice.Div(fee.Amount.Int, gasPrice)

		gasConsumed := bitInts[1].SetUint64(gasUsed)
		gasConsumed = gasConsumed.Mul(gasPrice, gasConsumed)

		gasCost := sdk.NewDecCoinFromDec(fee.Denom, sdk.NewDecWithBigIntAndPrec(gasConsumed, sdk.Precision))
		gasRefund := fee.Sub(gasCost)

		refundFees[i] = gasRefund
	}
	return refundFees
}

// CalculateRefundFees provides the way to calculate the refunded gas with gasUsed, fees and gasPrice,
// as refunded gas = fees - gasPrice * gasUsed
func CalculateRefundFees(gasUsed uint64, fees sdk.DecCoins, gasPrice *big.Int) sdk.Coins {
	gas := new(big.Int).Div(fees[0].Amount.BigInt(), gasPrice).Uint64()
	return calculateRefundFees(gasUsed, gas, fees)
}

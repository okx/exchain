package ante

import (
	"github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"math/big"
)

type AccountAggregateValidateDecorator struct {
	ak        auth.AccountKeeper
	sk        types.SupplyKeeper
	evmKeeper EVMKeeper
}

func NewAccountAggregateValidateDecorator(ak auth.AccountKeeper, sk types.SupplyKeeper, ek EVMKeeper) AccountAggregateValidateDecorator {
	return AccountAggregateValidateDecorator{
		ak:        ak,
		sk:        sk,
		evmKeeper: ek,
	}
}

// AnteHandle validates the signature and returns sender address
func (aavd AccountAggregateValidateDecorator) AnteHandle1(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	oldGasMeter := ctx.GasMeter()
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	msgEthTx, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
	}

	evmParams := aavd.evmKeeper.GetParams(ctx)
	if msgEthTx.GetGas() > evmParams.MaxGasLimitPerTx {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrTxTooLarge, "too large gas limit, it must be less than %d", evmParams.MaxGasLimitPerTx)
	}

	if ctx.From() != "" {
		msgEthTx.SetFrom(ctx.From())
	}

	address := msgEthTx.From()
	if address.Empty() {
		panic("sender address cannot be empty")
	}

	if ok := aavd.evmKeeper.IsContractInBlockedList(ctx, address); ok && evmParams.EnableContractBlockedList {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "address: %s has been blocked", address.String())
	}

	acc := aavd.ak.GetAccount(ctx, address)
	if acc == nil {
		return ctx, sdkerrors.Wrapf(
			sdkerrors.ErrUnknownAddress,
			"account %s (%s) is nil", common.BytesToAddress(address.Bytes()), address,
		)
	}

	// validate sender has enough funds to pay for gas cost
	balance := acc.GetCoins().AmountOf(sdk.DefaultBondDenom)
	if balance.BigInt().Cmp(msgEthTx.Cost()) < 0 {
		return ctx, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds,
			"sender balance < tx gas cost (%s%s < %s%s)", balance.String(), sdk.DefaultBondDenom, sdk.NewDecFromBigIntWithPrec(msgEthTx.Cost(), sdk.Precision).String(), sdk.DefaultBondDenom,
		)
	}

	seq := acc.GetSequence()
	if msgEthTx.Data.AccountNonce != seq {
		return ctx, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidSequence,
			"invalid nonce; got %d, expected %d", msgEthTx.Data.AccountNonce, seq,
		)
	}
	seq++
	if err := acc.SetSequence(seq); err != nil {
		return ctx, err
	}

	gasLimit := msgEthTx.GetGas()
	if gasLimit == 0 {
		return ctx, sdkerrors.Wrapf(
			sdkerrors.ErrOutOfGas, "invalid gas can not be 0")
	}
	// Charge sender for gas up to limit

	// Cost calculates the fees paid to validators based on gas limit and price
	cost := new(big.Int).Mul(msgEthTx.Data.Price, new(big.Int).SetUint64(gasLimit))
	feeAmt := sdk.NewCoins(
		sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewDecFromBigIntWithPrec(cost, sdk.Precision)), // int2dec
	)
	if !feeAmt.IsValid() {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", feeAmt)
	}

	oldCoins := acc.GetCoins()
	newCoins, hasNeg := oldCoins.SafeSub(feeAmt)
	if hasNeg {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", oldCoins, feeAmt)
	}
	if err := acc.SetCoins(newCoins); err != nil {
		return ctx, err
	}

	feeAcc := aavd.sk.GetModuleAccount(ctx, types.FeeCollectorName)
	if feeAcc == nil {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", types.FeeCollectorName)
	}
	feeCoin := feeAcc.GetCoins()
	feeNewCoin := feeCoin.Add(feeAmt...)
	if feeNewCoin.IsAnyNegative() {
		return ctx, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds, "insufficient account funds; %s < %s", feeCoin, feeAmt,
		)
	}
	if err := feeAcc.SetCoins(feeNewCoin); err != nil {
		return ctx, err
	}

	aavd.ak.SetAccount(ctx, acc)
	aavd.ak.SetAccount(ctx, feeAcc)

	ctx = ctx.WithGasMeter(oldGasMeter)
	return next(ctx, tx, simulate)
}

func (aavd AccountAggregateValidateDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	pinAnte(ctx.AnteTracer(), "AccountAggregateValidateDecorator")
	//oldGasMeter := ctx.GasMeter()
	//ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	msgEthTx, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
	}

	// simulate means 'eth_call' or 'eth_estimateGas', when it's 'eth_estimateGas' we set the sender from ctx.
	if ctx.From() != "" {
		msgEthTx.SetFrom(ctx.From())
	}
	address := msgEthTx.From()
	if address.Empty() {
		panic("sender address cannot be empty")
	}

	acc := aavd.ak.GetAccount(ctx, address)
	if acc == nil {
		return ctx, sdkerrors.Wrapf(
			sdkerrors.ErrUnknownAddress,
			"account %s (%s) is nil", common.BytesToAddress(address.Bytes()), address,
		)
	}

	seq := acc.GetSequence()
	if msgEthTx.Data.AccountNonce != seq {
		return ctx, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidSequence,
			"invalid nonce; got %d, expected %d", msgEthTx.Data.AccountNonce, seq,
		)
	}
	seq++
	if err := acc.SetSequence(seq); err != nil {
		panic(err)
	}
	gasLimit := msgEthTx.GetGas()
	gas, err := ethcore.IntrinsicGas(msgEthTx.Data.Payload, []ethtypes.AccessTuple{}, msgEthTx.To() == nil, true, false)
	if err != nil {
		return ctx, sdkerrors.Wrap(err, "failed to compute intrinsic gas cost")
	}

	// intrinsic gas verification during CheckTx
	if ctx.IsCheckTx() && gasLimit < gas {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrOutOfGas, "intrinsic gas too low: %d < %d", gasLimit, gas)
	}

	// Charge sender for gas up to limit
	if gasLimit != 0 {
		// Cost calculates the fees paid to validators based on gas limit and price
		cost := new(big.Int).Mul(msgEthTx.Data.Price, new(big.Int).SetUint64(gasLimit))

		evmDenom := sdk.DefaultBondDenom

		feeAmt := sdk.NewCoins(
			sdk.NewCoin(evmDenom, sdk.NewDecFromBigIntWithPrec(cost, sdk.Precision)), // int2dec
		)

		if !feeAmt.IsValid() {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", feeAmt)
		}

		oldCoins := acc.GetCoins()
		newCoins, hasNeg := oldCoins.SafeSub(feeAmt)
		if hasNeg {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
				"insufficient funds to pay for fees; %s < %s", oldCoins, feeAmt)
		}
		if err := acc.SetCoins(newCoins); err != nil {
			return ctx, err
		}
		aavd.ak.SetAccount(ctx, acc)

		recipientAcc := aavd.sk.GetModuleAccount(ctx, types.FeeCollectorName)
		if recipientAcc == nil {
			aavd.ak.Logger(ctx).Error("AccountAggregateValidateDecorator", "getfeeacc", "err")
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "module account %s does not exist", types.FeeCollectorName)
		}
		feeCoin := recipientAcc.GetCoins()
		feeNewCoin := feeCoin.Add(feeAmt...)
		if !feeNewCoin.IsValid() || feeCoin.IsAnyNegative() {
			return ctx, sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, feeCoin.String())
		}
		if err := recipientAcc.SetCoins(feeNewCoin); err != nil {
			return ctx, err
		}
		aavd.ak.SetAccount(ctx, recipientAcc)
	}
	newCtx = auth.SetGasMeter(simulate, ctx, gasLimit)
	return next(newCtx, tx, simulate)
}

package ante

import (
	"github.com/ethereum/go-ethereum/common"
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

func (aavd AccountAggregateValidateDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	pinAnte(ctx.AnteTracer(), "AccountAggregateValidateDecorator")
	//oldGasMeter := ctx.GasMeter()
	//ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	msgEthTx, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
	}

	pinAnte(ctx.AnteTracer(), "AccountAggregateValidateDecorator-getParams")
	evmParams := aavd.evmKeeper.GetParams(ctx)
	if msgEthTx.GetGas() > evmParams.MaxGasLimitPerTx {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrTxTooLarge, "too large gas limit, it must be less than %d", evmParams.MaxGasLimitPerTx)
	}
	pinAnte(ctx.AnteTracer(), "AccountAggregateValidateDecorator-getFrom")

	// simulate means 'eth_call' or 'eth_estimateGas', when it's 'eth_estimateGas' we set the sender from ctx.
	if ctx.From() != "" {
		msgEthTx.SetFrom(ctx.From())
	}
	address := msgEthTx.From()
	if address.Empty() {
		panic("sender address cannot be empty")
	}

	pinAnte(ctx.AnteTracer(), "AccountAggregateValidateDecorator-isContractBlocked")
	if evmParams.EnableContractBlockedList {
		if ok := aavd.evmKeeper.IsContractInBlockedList(ctx, address); ok {
			return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "address: %s has been blocked", address.String())
		}
	}
	pinAnte(ctx.AnteTracer(), "AccountAggregateValidateDecorator-Set/GetAccount")

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

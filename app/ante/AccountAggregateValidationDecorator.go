package ante

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
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
	if err := feeAcc.SetCoins(feeCoin); err != nil {
		return ctx, err
	}

	aavd.ak.SetAccount(ctx, acc)
	aavd.ak.SetAccount(ctx, feeAcc)

	ctx = ctx.WithGasMeter(oldGasMeter)
	return next(ctx, tx, simulate)
}

func (aavd AccountAggregateValidateDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	oldGasMeter := ctx.GasMeter()
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
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

	// on InitChain make sure account number == 0
	if ctx.BlockHeight() == 0 && acc.GetAccountNumber() != 0 {
		return ctx, sdkerrors.Wrapf(
			sdkerrors.ErrInvalidSequence,
			"invalid account number for height zero (got %d)", acc.GetAccountNumber(),
		)
	}

	evmDenom := sdk.DefaultBondDenom

	// validate sender has enough funds to pay for gas cost
	balance := acc.GetCoins().AmountOf(evmDenom)
	if balance.BigInt().Cmp(msgEthTx.Cost()) < 0 {
		return ctx, sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds,
			"sender balance < tx gas cost (%s%s < %s%s)", balance.String(), evmDenom, sdk.NewDecFromBigIntWithPrec(msgEthTx.Cost(), sdk.Precision).String(), evmDenom,
		)
	}

	seq := acc.GetSequence()
	// if multiple transactions are submitted in succession with increasing nonces,
	// all will be rejected except the first, since the first needs to be included in a block
	// before the sequence increments
	if ctx.IsCheckTx() {
		ctx = ctx.WithAccountNonce(seq)
		// will be checkTx and RecheckTx mode
		if ctx.IsReCheckTx() {
			// recheckTx mode

			// sequence must strictly increasing
			if msgEthTx.Data.AccountNonce != seq {
				return ctx, sdkerrors.Wrapf(
					sdkerrors.ErrInvalidSequence,
					"invalid nonce; got %d, expected %d", msgEthTx.Data.AccountNonce, seq,
				)
			}
		} else {
			if baseapp.IsMempoolEnablePendingPool() {
				if msgEthTx.Data.AccountNonce < seq {
					return ctx, sdkerrors.Wrapf(
						sdkerrors.ErrInvalidSequence,
						"invalid nonce; got %d, expected %d",
						msgEthTx.Data.AccountNonce, seq,
					)
				}
			} else {
				// checkTx mode
				checkTxModeNonce := seq

				if !baseapp.IsMempoolEnableRecheck() {
					// if is enable recheck, the sequence of checkState will increase after commit(), so we do not need
					// to add pending txs len in the mempool.
					// but, if disable recheck, we will not increase sequence of checkState (even in force recheck case, we
					// will also reset checkState), so we will need to add pending txs len to get the right nonce
					gPool := baseapp.GetGlobalMempool()
					if gPool != nil {
						cnt := gPool.GetUserPendingTxsCnt(evmtypes.EthAddressStringer(common.BytesToAddress(address.Bytes())).String())
						checkTxModeNonce = seq + uint64(cnt)
					}
				}

				if baseapp.IsMempoolEnableSort() {
					if msgEthTx.Data.AccountNonce < seq || msgEthTx.Data.AccountNonce > checkTxModeNonce {
						return ctx, sdkerrors.Wrapf(
							sdkerrors.ErrInvalidSequence,
							"invalid nonce; got %d, expected in the range of [%d, %d]",
							msgEthTx.Data.AccountNonce, seq, checkTxModeNonce,
						)
					}
				} else {
					if msgEthTx.Data.AccountNonce != checkTxModeNonce {
						return ctx, sdkerrors.Wrapf(
							sdkerrors.ErrInvalidSequence,
							"invalid nonce; got %d, expected %d",
							msgEthTx.Data.AccountNonce, checkTxModeNonce,
						)
					}
				}
			}
		}
	} else {
		// only deliverTx mode
		if msgEthTx.Data.AccountNonce != seq {
			return ctx, sdkerrors.Wrapf(
				sdkerrors.ErrInvalidSequence,
				"invalid nonce; got %d, expected %d", msgEthTx.Data.AccountNonce, seq,
			)
		}
	}

	ctx = ctx.WithGasMeter(oldGasMeter)
	return next(ctx, tx, simulate)
}

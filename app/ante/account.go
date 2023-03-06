package ante

import (
	"bytes"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/innertx"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/types"
	evmtypes "github.com/okx/okbchain/x/evm/types"
)

type accountKeeperInterface interface {
	SetAccount(ctx sdk.Context, acc exported.Account)
}

type AccountAnteDecorator struct {
	ak        auth.AccountKeeper
	sk        types.SupplyKeeper
	evmKeeper EVMKeeper
}

// NewAccountVerificationDecorator creates a new AccountVerificationDecorator
func NewAccountAnteDecorator(ak auth.AccountKeeper, ek EVMKeeper, sk types.SupplyKeeper) AccountAnteDecorator {
	return AccountAnteDecorator{
		ak:        ak,
		sk:        sk,
		evmKeeper: ek,
	}
}

func accountVerification(ctx *sdk.Context, acc exported.Account, tx *evmtypes.MsgEthereumTx) error {
	if ctx.BlockHeight() == 0 && acc.GetAccountNumber() != 0 {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidSequence,
			"invalid account number for height zero (got %d)", acc.GetAccountNumber(),
		)
	}

	const evmDenom = sdk.DefaultBondDenom

	feeInts := feeIntsPool.Get().(*[2]big.Int)
	defer feeIntsPool.Put(feeInts)

	// validate sender has enough funds to pay for gas cost
	balance := acc.GetCoins().AmountOf(evmDenom)
	if balance.Int.Cmp(tx.CalcCostTo(&feeInts[0])) < 0 {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds,
			"sender balance < tx gas cost (%s%s < %s%s)", balance.String(), evmDenom, sdk.NewDecFromBigIntWithPrec(tx.Cost(), sdk.Precision).String(), evmDenom,
		)
	}
	return nil
}

func nonceVerificationInCheckTx(ctx sdk.Context, seq uint64, msgEthTx *evmtypes.MsgEthereumTx, isReCheckTx bool) error {
	if isReCheckTx {
		// recheckTx mode
		// sequence must strictly increasing
		if msgEthTx.Data.AccountNonce != seq {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidSequence,
				"invalid nonce; got %d, expected %d", msgEthTx.Data.AccountNonce, seq,
			)
		}
	} else {
		if baseapp.IsMempoolEnablePendingPool() {
			if msgEthTx.Data.AccountNonce < seq {
				return sdkerrors.Wrapf(
					sdkerrors.ErrInvalidSequence,
					"invalid nonce; got %d, expected %d", msgEthTx.Data.AccountNonce, seq,
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
					addr := msgEthTx.GetSender(ctx)
					if pendingNonce, ok := gPool.GetPendingNonce(addr); ok {
						checkTxModeNonce = pendingNonce + 1
					}
				}
			}

			if baseapp.IsMempoolEnableSort() {
				if msgEthTx.Data.AccountNonce < seq || msgEthTx.Data.AccountNonce > checkTxModeNonce {
					accNonceStr := strconv.FormatUint(msgEthTx.Data.AccountNonce, 10)
					seqStr := strconv.FormatUint(seq, 10)
					checkTxModeNonceStr := strconv.FormatUint(checkTxModeNonce, 10)

					errStr := strings.Join([]string{
						"invalid nonce; got ", accNonceStr,
						", expected in the range of [", seqStr, ", ", checkTxModeNonceStr, "]"},
						"")

					return sdkerrors.WrapNoStack(sdkerrors.ErrInvalidSequence, errStr)
				}
			} else {
				if msgEthTx.Data.AccountNonce != checkTxModeNonce {
					accNonceStr := strconv.FormatUint(msgEthTx.Data.AccountNonce, 10)
					checkTxModeNonceStr := strconv.FormatUint(checkTxModeNonce, 10)

					errStr := strings.Join([]string{
						"invalid nonce; got ", accNonceStr, ", expected ", checkTxModeNonceStr},
						"")

					return sdkerrors.WrapNoStack(sdkerrors.ErrInvalidSequence, errStr)
				}
			}
		}
	}
	return nil
}

func nonceVerification(ctx sdk.Context, acc exported.Account, msgEthTx *evmtypes.MsgEthereumTx) (sdk.Context, error) {
	seq := acc.GetSequence()
	// if multiple transactions are submitted in succession with increasing nonces,
	// all will be rejected except the first, since the first needs to be included in a block
	// before the sequence increments
	if ctx.IsCheckTx() {
		ctx.SetAccountNonce(seq)
		// will be checkTx and RecheckTx mode
		err := nonceVerificationInCheckTx(ctx, seq, msgEthTx, ctx.IsReCheckTx())
		if err != nil {
			return ctx, err
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
	return ctx, nil
}

func ethGasConsume(ik innertx.InnerTxKeeper, ak accountKeeperInterface, sk types.SupplyKeeper, ctx *sdk.Context, acc exported.Account, accGetGas sdk.Gas, msgEthTx *evmtypes.MsgEthereumTx, simulate bool) error {
	gasLimit := msgEthTx.GetGas()
	gas, err := ethcore.IntrinsicGas(msgEthTx.Data.Payload, []ethtypes.AccessTuple{}, msgEthTx.To() == nil, true, false)
	if err != nil {
		return sdkerrors.Wrap(err, "failed to compute intrinsic gas cost")
	}

	// intrinsic gas verification during CheckTx
	if ctx.IsCheckTx() && gasLimit < gas {
		return sdkerrors.Wrapf(sdkerrors.ErrOutOfGas, "intrinsic gas too low: %d < %d", gasLimit, gas)
	}

	// Charge sender for gas up to limit
	if gasLimit != 0 {
		feeInts := feeIntsPool.Get().(*[2]big.Int)
		defer feeIntsPool.Put(feeInts)
		// Cost calculates the fees paid to validators based on gas limit and price
		cost := (&feeInts[0]).SetUint64(gasLimit)
		cost = cost.Mul(msgEthTx.Data.Price, cost)

		const evmDenom = sdk.DefaultBondDenom

		feeAmt := sdk.NewDecCoinsFromDec(evmDenom, sdk.NewDecWithBigIntAndPrec(cost, sdk.Precision))

		ctx.UpdateFromAccountCache(acc, accGetGas)

		err = deductFees(ik, ak, sk, *ctx, acc, feeAmt)
		if err != nil {
			return err
		}
	}

	// Set gas meter after ante handler to ignore gaskv costs
	auth.SetGasMeter(simulate, ctx, gasLimit)
	return nil
}

func deductFees(ik innertx.InnerTxKeeper, ak accountKeeperInterface, sk types.SupplyKeeper, ctx sdk.Context, acc exported.Account, fees sdk.Coins) error {
	blockTime := ctx.BlockTime()
	coins := acc.GetCoins()

	if !fees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid fee amount: %s", fees)
	}

	// verify the account has enough funds to pay for fees
	balance, hasNeg := coins.SafeSub(fees)
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", coins, fees)
	}

	// Validate the account has enough "spendable" coins as this will cover cases
	// such as vesting accounts.
	spendableCoins := acc.SpendableCoins(blockTime)
	if _, hasNeg := spendableCoins.SafeSub(fees); hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to pay for fees; %s < %s", spendableCoins, fees)
	}

	// set coins and record innertx
	err := acc.SetCoins(balance)
	if !ctx.IsCheckTx() {
		toAcc := sk.GetModuleAddress(types.FeeCollectorName)
		ik.UpdateInnerTx(ctx.TxBytes(), ctx.BlockHeight(), innertx.CosmosDepth, acc.GetAddress(), toAcc, innertx.CosmosCallType, innertx.SendCallName, fees, err)
	}
	if err != nil {
		return err
	}
	ak.SetAccount(ctx, acc)

	return nil
}

func incrementSeq(ctx sdk.Context, msgEthTx *evmtypes.MsgEthereumTx, accAddress sdk.AccAddress, ak auth.AccountKeeper, acc exported.Account) {
	if ctx.IsCheckTx() && !ctx.IsReCheckTx() && !baseapp.IsMempoolEnableRecheck() && !ctx.IsTraceTx() {
		return
	}

	// get and set account must be called with an infinite gas meter in order to prevent
	// additional gas from being deducted.
	infGasMeter := sdk.GetReusableInfiniteGasMeter()
	defer sdk.ReturnInfiniteGasMeter(infGasMeter)
	ctx.SetGasMeter(infGasMeter)

	// increment sequence of all signers
	// eth tx only has one signer
	if accAddress.Empty() {
		accAddress = msgEthTx.AccountAddress()
	}
	var sacc exported.Account
	if acc != nil && bytes.Equal(accAddress, acc.GetAddress()) {
		// because we use infinite gas meter, we can don't care about the gas
		sacc = acc
	} else {
		sacc = ak.GetAccount(ctx, accAddress)
	}
	seq := sacc.GetSequence()
	if !baseapp.IsMempoolEnablePendingPool() {
		seq++
	} else if msgEthTx.Data.AccountNonce == seq {
		seq++
	}
	if err := sacc.SetSequence(seq); err != nil {
		panic(err)
	}
	ak.SetAccount(ctx, sacc)

	return
}

func (avd AccountAnteDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	msgEthTx, ok := tx.(*evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
	}

	var acc exported.Account
	var getAccGasUsed sdk.Gas

	address := msgEthTx.AccountAddress()
	if address.Empty() && ctx.From() != "" {
		msgEthTx.SetFrom(ctx.From())
		address = msgEthTx.AccountAddress()
	}

	if !simulate {
		if address.Empty() {
			panic("sender address cannot be empty")
		}
		if ctx.IsCheckTx() {
			acc = avd.ak.GetAccount(ctx, address)
			if acc == nil {
				acc = avd.ak.NewAccountWithAddress(ctx, address)
				avd.ak.SetAccount(ctx, acc)
			}
			// on InitChain make sure account number == 0
			err = accountVerification(&ctx, acc, msgEthTx)
			if err != nil {
				return ctx, err
			}
		}

		acc, getAccGasUsed = getAccount(&avd.ak, &ctx, address, acc)
		if acc == nil {
			return ctx, sdkerrors.Wrapf(
				sdkerrors.ErrUnknownAddress,
				"account %s (%s) is nil", common.BytesToAddress(address.Bytes()), address,
			)
		}

		// account would not be updated
		ctx, err = nonceVerification(ctx, acc, msgEthTx)
		if err != nil {
			return ctx, err
		}

		// consume gas for compatible
		ctx.GasMeter().ConsumeGas(getAccGasUsed, "get account")

		ctx.EnableAccountCache()
		// account would be updated
		err = ethGasConsume(avd.evmKeeper, avd.ak, avd.sk, &ctx, acc, getAccGasUsed, msgEthTx, simulate)
		acc = nil
		acc, _ = ctx.GetFromAccountCacheData().(exported.Account)
		ctx.DisableAccountCache()
		if err != nil {
			return ctx, err
		}
	}

	incrementSeq(ctx, msgEthTx, address, avd.ak, acc)

	return next(ctx, tx, simulate)
}

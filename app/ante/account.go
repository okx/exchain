package ante

import (
	"bytes"
	"math/big"

	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"

	"github.com/ethereum/go-ethereum/common"
	ethcore "github.com/ethereum/go-ethereum/core"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

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

func accountVertification(ctx *sdk.Context, acc exported.Account, tx evmtypes.MsgEthereumTx) error {
	if ctx.BlockHeight() == 0 && acc.GetAccountNumber() != 0 {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidSequence,
			"invalid account number for height zero (got %d)", acc.GetAccountNumber(),
		)
	}

	evmDenom := sdk.DefaultBondDenom

	// validate sender has enough funds to pay for gas cost
	balance := acc.GetCoins().AmountOf(evmDenom)
	if balance.BigInt().Cmp(tx.Cost()) < 0 {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInsufficientFunds,
			"sender balance < tx gas cost (%s%s < %s%s)", balance.String(), evmDenom, sdk.NewDecFromBigIntWithPrec(tx.Cost(), sdk.Precision).String(), evmDenom,
		)
	}
	return nil
}

func nounceVertification(ctx sdk.Context, acc exported.Account, msgEthTx evmtypes.MsgEthereumTx) (sdk.Context, error) {
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
						cnt := gPool.GetUserPendingTxsCnt(evmtypes.EthAddressStringer(common.BytesToAddress(msgEthTx.From().Bytes())).String())
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
	return ctx, nil
}

func ethGasConsume(ctx sdk.Context, acc exported.Account, accGetGas sdk.Gas, msgEthTx evmtypes.MsgEthereumTx, simulate bool, sk types.SupplyKeeper) (sdk.Context, error) {
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

		if ctx.AccountCache() != nil {
			cache := ctx.AccountCache()
			cache.FromAcc = acc
			cache.FromAccGettedGas = accGetGas
			ctx.SetAccountCache(cache)
		} else {
			ctx.SetAccountCache(&sdk.AccountCache{FromAcc: acc, FromAccGettedGas: accGetGas})
		}

		err = auth.DeductFees(sk, ctx, acc, feeAmt)
		if err != nil {
			return ctx, err
		}
	}

	// Set gas meter after ante handler to ignore gaskv costs
	ctx = auth.SetGasMeter(simulate, ctx, gasLimit)
	return ctx, nil
}

func incrementSeq(ctx sdk.Context, msgEthTx evmtypes.MsgEthereumTx, ak auth.AccountKeeper, acc exported.Account) {
	if ctx.IsCheckTx() && !ctx.IsReCheckTx() && !baseapp.IsMempoolEnableRecheck() && !ctx.IsTraceTx() {
		return
	}

	// get and set account must be called with an infinite gas meter in order to prevent
	// additional gas from being deducted.
	ctx = ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	// increment sequence of all signers
	for _, addr := range msgEthTx.GetSigners() {
		var sacc exported.Account
		if bytes.Equal(addr, acc.GetAddress().Bytes()) {
			sacc = getAccount(&ak, &ctx, addr, sacc)
		} else {
			sacc = ak.GetAccount(ctx, addr)
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
	}
	return
}

func (avd AccountAnteDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	msgEthTx, ok := tx.(evmtypes.MsgEthereumTx)
	if !ok {
		return ctx, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid transaction type: %T", tx)
	}
	address := msgEthTx.From()
	if address.Empty() {
		panic("sender address cannot be empty")
	}

	var acc exported.Account

	if !simulate {
		if ctx.IsCheckTx() {
			acc = avd.ak.GetAccount(ctx, address)
			if acc == nil {
				acc = avd.ak.NewAccountWithAddress(ctx, address)
				avd.ak.SetAccount(ctx, acc)
			}
			// on InitChain make sure account number == 0
			err = accountVertification(&ctx, acc, msgEthTx)
			if err != nil {
				return ctx, err
			}
		}

		acc = getAccount(&avd.ak, &ctx, address, acc)
		if acc == nil {
			return ctx, sdkerrors.Wrapf(
				sdkerrors.ErrUnknownAddress,
				"account %s (%s) is nil", common.BytesToAddress(address.Bytes()), address,
			)
		}

		// account would not be updated
		ctx, err = nounceVertification(ctx, acc, msgEthTx)
		if err != nil {
			return ctx, err
		}

		// consume gas for compatible
		gas := getAccountGas(&avd.ak, acc)
		ctx.GasMeter().ConsumeGas(gas, "get account")
		// fmt.Printf("gas used: %d; changed: %d\n", ctx.GasMeter().GasConsumed(), gas)

		ctx.SetAccountCache(&sdk.AccountCache{})
		// account would be updated
		ctx, err = ethGasConsume(ctx, acc, gas, msgEthTx, simulate, avd.sk)
		if err != nil {
			return ctx, err
		}
		acc = nil // account has be updated
		if ctx.AccountCache() != nil && ctx.AccountCache().FromAcc != nil {
			acc = ctx.AccountCache().FromAcc.(exported.Account)
		}
	}

	incrementSeq(ctx, msgEthTx, avd.ak, acc)

	return next(ctx, tx, simulate)
}

package ante

import (
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"strconv"
	"strings"
)

func GetCheckTxNonceFromMempool(addr string) uint64 {
	if !baseapp.IsMempoolEnableRecheck() {
		// if is enable recheck, the sequence of checkState will increase after commit(), so we do not need
		// to add pending txs len in the mempool.
		// but, if disable recheck, we will not increase sequence of checkState (even in force recheck case, we
		// will also reset checkState), so we will need to add pending txs len to get the right nonce
		gPool := baseapp.GetGlobalMempool()
		if gPool != nil {
			if pendingNonce, ok := gPool.GetPendingNonce(addr); ok {
				return pendingNonce + 1
			}
		}
	}
	return 0
}

func nonceVerification(ctx sdk.Context, seq uint64, txNonce uint64, addr string, simulate bool) error {
	if simulate || //
		(txNonce == 0) || // no wrapCMtx no need verify
		!ctx.IsCheckTx() || // deliverTx mode no need check
		ctx.IsReCheckTx() { // recheckTx mode sequence must strictly increasing, get nonce from account
		return nil
	}
	// will be checkTx mode
	err := nonceVerificationInCheckTx(seq, txNonce, addr)
	if err != nil {
		return err
	}
	return nil
}

func nonceVerificationInCheckTx(seq uint64, txNonce uint64, addr string) error {
	if baseapp.IsMempoolEnablePendingPool() {
		if txNonce < seq {
			return sdkerrors.Wrapf(
				sdkerrors.ErrInvalidSequence,
				"cmtx invalid nonce; got %d, expected %d", txNonce, seq,
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
				if pendingNonce, ok := gPool.GetPendingNonce(addr); ok {
					checkTxModeNonce = pendingNonce + 1
				}
			}
		}

		if baseapp.IsMempoolEnableSort() { // this is only for replace the same nonce tx
			if txNonce < seq || txNonce > checkTxModeNonce {
				accNonceStr := strconv.FormatUint(txNonce, 10)
				seqStr := strconv.FormatUint(seq, 10)
				checkTxModeNonceStr := strconv.FormatUint(checkTxModeNonce, 10)

				errStr := strings.Join([]string{
					"cmtx invalid nonce; got ", accNonceStr,
					", expected in the range of [", seqStr, ", ", checkTxModeNonceStr, "]"},
					"")

				return sdkerrors.WrapNoStack(sdkerrors.ErrInvalidSequence, errStr)
			}
		} else {
			if txNonce != checkTxModeNonce {
				accNonceStr := strconv.FormatUint(txNonce, 10)
				checkTxModeNonceStr := strconv.FormatUint(checkTxModeNonce, 10)

				errStr := strings.Join([]string{
					"cmtx invalid nonce; got ", accNonceStr, ", expected ", checkTxModeNonceStr},
					"")

				return sdkerrors.WrapNoStack(sdkerrors.ErrInvalidSequence, errStr)
			}
		}
	}
	return nil
}

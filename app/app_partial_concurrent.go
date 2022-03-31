package app

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authante "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ante"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

//type SetAccountObserver func(o keeper.ObserverI) ()

// getTxFeeAndFromHandler get tx fee and from
func getTxFeeAndFromHandler(ak auth.AccountKeeper) sdk.GetTxFeeAndFromHandler {
	return func(ctx sdk.Context, tx sdk.Tx) (fee sdk.Coins, isEvm bool, from string) {
		if evmTx, ok := tx.(*evmtypes.MsgEthereumTx); ok {
			isEvm = true
			_ = evmTx.VerifySig(evmTx.ChainID(), ctx.BlockHeight())
			//from = evmTx.From
			fee = evmTx.GetFee()
			feePayer := evmTx.FeePayer(ctx)//.AccountAddress()
			//from = ak.GetAccount(ctx, feePayer)
			feePayerAcc := ak.GetAccount(ctx, feePayer)
			from = feePayerAcc.GetAddress().String()//hex.EncodeToString(feePayerAcc.GetAddress())
		} else if feeTx, ok := tx.(authante.FeeTx); ok {
			fee = feeTx.GetFee()
			feePayer := feeTx.FeePayer(ctx)
			//from = ak.GetAccount(ctx, feePayer)
			feePayerAcc := ak.GetAccount(ctx, feePayer)
			from = feePayerAcc.GetAddress().String()// ex17xpfvakm2amg962yls6f84z3kell8c5lcs49z2
			//hex.EncodeToString(feePayerAcc.GetAddress())// f1829676db577682e944fc3493d451b67ff3e29f
		}

		return
	}
}

//func setAccountObserver(ak auth.AccountKeeper) sdk.SetAccountObserver {
//	return func(funco keeper.ObserverI) {
//		ak.SetObserverKeeper(o)
//	}
//}

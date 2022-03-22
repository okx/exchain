package app

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authante "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ante"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/keeper"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

type SetAccountObserver func(o keeper.ObserverI) ()

// getTxFeeAndFromHandler get tx fee and from
func getTxFeeAndFromHandler(ak auth.AccountKeeper) sdk.GetTxFeeAndFromHandler {
	return func(ctx sdk.Context, tx sdk.Tx) (fee sdk.Coins, isEvm bool, from sdk.Address) {
		if evmTx, ok := tx.(*evmtypes.MsgEthereumTx); ok {
			isEvm = true
			_ = evmTx.VerifySig(evmTx.ChainID(), ctx.BlockHeight())
			from = evmTx.AccountAddress()
		}
		if feeTx, ok := tx.(authante.FeeTx); ok {
			fee = feeTx.GetFee()
			feePayer := feeTx.FeePayer(ctx)
			feePayerAcc := ak.GetAccount(ctx, feePayer)
			from = feePayerAcc.GetAddress()
			//hex.EncodeToString(feePayerAcc.GetAddress())
		}

		return
	}
}

func setAccountObserver(ak auth.AccountKeeper) SetAccountObserver {
	//return func(o keeper.ObserverI) {
	//	ak.SetObserverKeeper(o)
	//}
	return func(o keeper.ObserverI) {
		ak.SetObserverKeeper(o)
	}
}

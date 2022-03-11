package app

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authante "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ante"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// feeCollectorHandler set or get the value of feeCollectorAcc
func updateFeeCollectorHandler(bk bank.Keeper, sk supply.Keeper) sdk.UpdateFeeCollectorAccHandler {
	return func(ctx sdk.Context, balance sdk.Coins) error {
		//feeCollectorAcc := sk.GetModuleAddress(auth.FeeCollectorName)
		//fmt.Println("FeeCollectorAcc", feeCollectorAcc)
		return bk.SetCoins(ctx, sk.GetModuleAddress(auth.FeeCollectorName), balance)
	}
}

// evmTxFeeHandler get tx fee for evm tx
func evmTxFeeHandler() sdk.GetTxFeeHandler {
	return func(ctx sdk.Context, tx sdk.Tx) (fee sdk.Coins, isEvm bool, signCache sdk.SigCache) {
		if evmTx, ok := tx.(evmtypes.MsgEthereumTx); ok {
			isEvm = true
			signCache, _ = evmTx.VerifySig(evmTx.ChainID(), ctx.BlockHeight(), ctx.TxBytes(), ctx.SigCache())
		}
		if feeTx, ok := tx.(authante.FeeTx); ok {
			fee = feeTx.GetFee()
		}

		return
	}
}

// fixLogForParallelTxHandler fix log for parallel tx
func fixLogForParallelTxHandler(ek *evm.Keeper) sdk.LogFix {
	return func(execResults [][]string) (logs [][]byte) {
		return ek.FixLog(execResults)
	}
}


// evmTxFromHandler get tx fee for evm tx
func evmTxFromHandler(ak auth.AccountKeeper) sdk.EvmTxFromHandler {
	return func(ctx sdk.Context, tx sdk.Tx) (evmTx sdk.Tx, fee sdk.Coins, isEvm bool, from sdk.Address, signCache sdk.SigCache) {
		if evmTxTmp, ok := tx.(evmtypes.MsgEthereumTx); ok {
			isEvm = true
			signCache, _ = evmTxTmp.VerifySig(evmTxTmp.ChainID(), ctx.BlockHeight(), ctx.TxBytes(), ctx.SigCache())
			evmTxTmp.SetFromUseSigCache(signCache)
			from = evmTxTmp.From()
			evmTx = evmTxTmp
		}
		if feeTx, ok := tx.(authante.FeeTx); ok {
			fee = feeTx.GetFee()
			feePayer := feeTx.FeePayer(ctx)
			feePayerAcc := ak.GetAccount(ctx, feePayer)
			from = feePayerAcc.GetAddress()
		}

		return
	}
	//return func(ctx sdk.Context, tx sdk.Tx) (sdk.Tx, bool) {
	//	if ctx.SigCache() != nil {
	//		if evmTx, ok := tx.(evmtypes.MsgEthereumTx); ok {
	//			evmTx.SetFromUseSigCache(ctx.SigCache())
	//			//log.Printf("evmTxFromHandler from: %s\n", hex.EncodeToString(evmTx.From().Bytes()))
	//			return evmTx, true
	//		}
	//	}
	//
	//	return evmtypes.MsgEthereumTx{}, false
	//}
}
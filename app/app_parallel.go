package app

import (
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authante "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ante"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/pkg/errors"
)

// feeCollectorHandler set or get the value of feeCollectorAcc
func updateFeeCollectorHandler(bk bank.Keeper, sk supply.Keeper) sdk.UpdateFeeCollectorAccHandler {
	return func(ctx sdk.Context, balance sdk.Coins) error {
		return bk.SetCoins(ctx, sk.GetModuleAddress(auth.FeeCollectorName), balance)
	}
}

// evmTxFeeHandler get tx fee for evm tx
func evmTxFeeHandler() sdk.GetTxFeeHandler {
	return func(ctx sdk.Context, tx sdk.Tx) (fee sdk.Coins, isEvm bool) {
		if evmTx, ok := tx.(*evmtypes.MsgEthereumTx); ok {
			isEvm = true
			if evmTx.BaseTx.From == "" && ctx.From() != "" {
				evmTx.BaseTx.From = ctx.From()
			} else {
				chainIDEpoch, err := ethermint.ParseChainID(ctx.ChainID())
				if err == nil {
					_ = evmTx.VerifySig(chainIDEpoch, ctx.BlockHeight())
				}
			}
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

func evmTxVerifySigHandler() sdk.TxVerifySigHandler {
	return func(ctx sdk.Context, tx sdk.Tx) error {
		if evmTx, ok := tx.(*evmtypes.MsgEthereumTx); ok {
			if evmTx.BaseTx.From == "" && ctx.From() != "" {
				evmTx.BaseTx.From = ctx.From()
				return nil
			}
			chainIDEpoch, err := ethermint.ParseChainID(ctx.ChainID())
			if err != nil {
				return err
			}
			err = evmTx.VerifySig(chainIDEpoch, ctx.BlockHeight())
			if err != nil {
				return err
			}
			return nil
		}
		return errors.New("tx type is not evm tx")
	}
}

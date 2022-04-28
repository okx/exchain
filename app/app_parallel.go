package app

import (
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// feeCollectorHandler set or get the value of feeCollectorAcc
func updateFeeCollectorHandler(bk bank.Keeper, sk supply.Keeper) sdk.UpdateFeeCollectorAccHandler {
	return func(ctx sdk.Context, balance sdk.Coins) error {
		return bk.SetCoins(ctx, sk.GetModuleAccount(ctx, auth.FeeCollectorName).GetAddress(), balance)
	}
}

// fixLogForParallelTxHandler fix log for parallel tx
func fixLogForParallelTxHandler(ek *evm.Keeper) sdk.LogFix {
	return func(logIndex []int, anteErrs []error) (logs [][]byte) {
		return ek.FixLog(logIndex, anteErrs)
	}
}

func preDeliverTxHandler(ak auth.AccountKeeper) sdk.PreDeliverTxHandler {
	return func(ctx sdk.Context, tx sdk.Tx, onlyVerifySig bool) {
		if evmTx, ok := tx.(*evmtypes.MsgEthereumTx); ok {
			if evmTx.BaseTx.From == "" {
				if ctx.From() != "" {
					evmTx.BaseTx.From = ctx.From()
				}
			}
			if evmTx.BaseTx.From == "" {
				_ = evmTxVerifySigHandler(ctx.ChainID(), ctx.BlockHeight(), evmTx)
			}

			if onlyVerifySig {
				return
			}

			if from := evmTx.AccountAddress(); from != nil {
				ak.LoadAccount(ctx, from)
			}
			if to := evmTx.Data.Recipient; to != nil {
				ak.LoadAccount(ctx, to.Bytes())
			}
		}
	}
}

func evmTxVerifySigHandler(chainID string, blockHeight int64, evmTx *evmtypes.MsgEthereumTx) error {
	chainIDEpoch, err := ethermint.ParseChainID(chainID)
	if err != nil {
		return err
	}
	err = evmTx.VerifySig(chainIDEpoch, blockHeight)
	if err != nil {
		return err
	}
	return nil
}

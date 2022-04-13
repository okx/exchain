package watcher

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	evm "github.com/okex/exchain/x/evm/keeper"
)

func NewWatcherHandler(evmKeeper *evm.Keeper) sdk.WatcherHandler {
	return func(
		tx sdk.Tx, evmResultData interface{}, gasUsed uint64,
	) (err error) {
		var watcherHandler sdk.WatcherHandler

		if tx.GetType() == sdk.EvmTxType {
			watcherHandler = NewHandler(evmKeeper)
		} else {
			return nil
		}
		return watcherHandler(tx, evmResultData, gasUsed)
	}
}

type Handler struct {
	EvmKeeper *evm.Keeper
}

func NewHandler(evmKeeper *evm.Keeper) sdk.WatcherHandler {
	handler := Handler{
		EvmKeeper: evmKeeper,
	}

	return func(tx sdk.Tx, evmResultData interface{}, gasUsed uint64) (err error) {
		return handler.SaveEvmTxAndReceipt(tx, evmResultData, gasUsed)
	}
}

func (handler Handler) SaveEvmTxAndReceipt(tx sdk.Tx, evmResultData interface{}, gasUsed uint64) error {
	if evmResultData == nil {
		txHash := ethcmn.BytesToHash(tx.TxHash())
		return handler.saveEvmTxAndFailedReceipt(tx, txHash, gasUsed)
	}
	return handler.saveEvmTxAndSuccessReceipt(tx, evmResultData, gasUsed)
}

func (handler Handler) saveEvmTxAndSuccessReceipt(tx sdk.Tx, evmResultData interface{}, gasUsed uint64) error {
	return handler.EvmKeeper.Watcher.SaveTxAndSuccessReceipt(tx, handler.EvmKeeper.TxIndexInBlock, evmResultData, gasUsed)
}

func (handler Handler) saveEvmTxAndFailedReceipt(tx sdk.Tx, txHash ethcmn.Hash, gasUsed uint64) error {
	return handler.EvmKeeper.Watcher.SaveTxAndFailedReceipt(tx, handler.EvmKeeper.TxIndexInBlock, txHash, gasUsed)
}

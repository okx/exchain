package watcher

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	evm "github.com/okex/exchain/x/evm/keeper"
)

func NewWatcherHandler(evmKeeper evm.Keeper) sdk.WatcherHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, evmResultData interface{},
	) (err error) {
		var watcherHandler sdk.WatcherHandler

		if tx.GetType() == sdk.EvmTxType {
			watcherHandler = NewHandler(evmKeeper)
		} else {
			return nil
		}
		return watcherHandler(ctx, tx, evmResultData)
	}
}

type Handler struct {
	EvmKeeper evm.Keeper
}

func NewHandler(evmKeeper evm.Keeper) sdk.WatcherHandler {
	handler := Handler{
		EvmKeeper: evmKeeper,
	}

	return func(ctx sdk.Context, tx sdk.Tx, evmResultData interface{}) (err error) {
		return handler.SaveEvmTxAndReceipt(ctx, tx, evmResultData)
	}
}

func (handler Handler) SaveEvmTxAndReceipt(ctx sdk.Context, tx sdk.Tx, evmResultData interface{}) error {
	return nil
}

func (handler Handler) saveEvmTxAndSuccessReceipt(ctx sdk.Context, tx sdk.Tx) {
}

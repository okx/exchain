package keeper

import (
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerror "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/x/evm/types"
)

var (
	_ types.EvmHooks = MultiEvmHooks{}
	_ types.EvmHooks = LogProcessEvmHook{}
)

// MultiEvmHooks combine multiple evm hooks, all hook functions are run in array sequence
type MultiEvmHooks []types.EvmHooks

// NewMultiEvmHooks combine multiple evm hooks
func NewMultiEvmHooks(hooks ...types.EvmHooks) MultiEvmHooks {
	return hooks
}

// PostTxProcessing delegate the call to underlying hooks
func (mh MultiEvmHooks) PostTxProcessing(ctx sdk.Context, from common.Address, to *common.Address, receipt *ethtypes.Receipt) error {
	for i := range mh {
		if err := mh[i].PostTxProcessing(ctx, from, to, receipt); err != nil {
			return sdkerror.Wrapf(err, "EVM hook %T failed", mh[i])
		}
	}
	return nil
}

// LogProcessEvmHook is an evm hook that convert specific contract logs into native module calls
type LogProcessEvmHook struct {
	handlers map[common.Hash]types.EvmLogHandler
}

func NewLogProcessEvmHook(handlers ...types.EvmLogHandler) LogProcessEvmHook {
	handlerMap := make(map[common.Hash]types.EvmLogHandler)
	for _, h := range handlers {
		handlerMap[h.EventID()] = h
	}
	return LogProcessEvmHook{handlerMap}
}

// PostTxProcessing delegate the call to underlying hooks
func (lh LogProcessEvmHook) PostTxProcessing(ctx sdk.Context, from common.Address, to *common.Address, receipt *ethtypes.Receipt) error {
	for _, log := range receipt.Logs {
		if len(log.Topics) == 0 {
			continue
		}
		if handler, ok := lh.handlers[log.Topics[0]]; ok {
			if err := handler.Handle(ctx, log.Address, log.Data); err != nil {
				return err
			}
		}
	}
	return nil
}

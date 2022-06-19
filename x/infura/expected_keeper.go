package infura

import evm "github.com/okex/exchain/x/evm/watcher"

type EvmKeeper interface {
	SetObserverKeeper(keeper evm.InfuraKeeper)
}

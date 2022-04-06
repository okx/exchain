package evm

import (
	ethvm "github.com/ethereum/go-ethereum/core/vm"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
)

func init() {
	server.TrapSignal(func() {
		ethvm.CloseDB()
	})
}

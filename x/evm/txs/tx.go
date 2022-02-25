package txs

import "github.com/okex/exchain/x/evm/types"

type Tx interface {
	Prepare(msg *types.MsgEthereumTx) error
	Transition() error
	Finalize() error
}

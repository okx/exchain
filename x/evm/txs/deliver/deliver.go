package deliver

import "github.com/okex/exchain/x/evm/txs/base"

type tx = base.Tx

func NewTx(config base.Config) *tx {
	return base.NewTx(config)
}

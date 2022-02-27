package check

import (
	"github.com/okex/exchain/x/evm/txs/base"
)

type tx struct {
	baseTx *base.Tx
}

func NewTx(config base.Config) *tx {
	return &tx{
		baseTx: base.NewTx(config),
	}
}

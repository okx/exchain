package check

import (
	"github.com/okex/exchain/x/evm/txs/base"
)

type Tx struct {
	*base.Tx
}

func NewTx(config base.Config) *Tx {
	return &Tx{
		Tx: base.NewTx(config),
	}
}

package check

import (
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/types"
)

type tx struct {
	baseTx base.Tx
}

func NewTx(config base.Config) *tx {
	return &tx{
		base.Tx{Config: config},
	}
}

func (tx *tx) Prepare(msg *types.MsgEthereumTx) (err error) {
	return tx.baseTx.Prepare(msg)
}

func (tx *tx) Transition() error {
	//TODO implement me
	panic("implement me")
}

func (tx *tx) Finalize() error {
	//TODO implement me
	panic("implement me")
}

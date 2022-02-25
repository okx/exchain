package check

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/types"
)

type tx struct {
	baseTx *base.Tx
}

func NewTx(config base.Config) *tx {
	return &tx{
		baseTx: base.NewTx(config),
	}
}

// Exec simulated tx do not submit state db.
func (tx *tx) Exec(msg *types.MsgEthereumTx) (result *sdk.Result, err error) {
	err = tx.baseTx.Prepare(msg)
	if err != nil {
		return
	}

	err = tx.baseTx.Transition()
	if err != nil {
		return tx.baseTx.DecorateError(err)
	}

	result = tx.baseTx.Emit(msg)

	return nil, nil

}

func (tx *tx) Finalize() error {

	//TODO implement me
	panic("implement me")
}

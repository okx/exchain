package check

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/types"
)

type tx struct {
	baseTx *base.Tx
}

func (t tx) Prepare(msg *types.MsgEthereumTx) (err error) {
	//TODO implement me
	panic("implement me")
}

func (t tx) SaveTx(msg *types.MsgEthereumTx) {
	//TODO implement me
	panic("implement me")
}

func (t tx) Transition() (result *base.Result, err error) {
	//TODO implement me
	panic("implement me")
}

func (t tx) DecorateResultError(result *base.Result, err error) (*base.Result, error) {
	//TODO implement me
	panic("implement me")
}

func (t tx) EmitEvent() *sdk.Result {
	//TODO implement me
	panic("implement me")
}

func NewTx(config base.Config) *tx {
	return &tx{
		baseTx: base.NewTx(config),
	}
}

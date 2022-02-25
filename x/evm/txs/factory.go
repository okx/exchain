package txs

import (
	"fmt"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/txs/check"
)

type factory struct {
	base.Config
}

func NewFactory(config base.Config) *factory {
	return &factory{config}
}

func (factory *factory) CreateTx() (Tx, error) {
	if factory == nil {
		return nil, fmt.Errorf("evm txs factory not inited")
	}
	if factory.Ctx.IsCheckTx() {
		return check.NewTx(factory.Config), nil
	}

	return nil, fmt.Errorf("unkown evm txs")
}

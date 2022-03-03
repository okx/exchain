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

// AnalyzeStart check Tx do not analyze start
func (t *Tx) AnalyzeStart(tag string) {}

// AnalyzeStop check Tx do not analyze stop
func (t *Tx) AnalyzeStop(tag string) {}

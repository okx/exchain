package tracetxlog

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/txs/check"
)

// tx trace tx log depends on check tx
type tx struct {
	*check.Tx
}

func NewTx(config base.Config) *tx {
	return &tx{
		Tx: check.NewTx(config),
	}
}

// DecorateResult trace log tx need modify the result to log, and swallow error
func (t *tx) DecorateResult(inResult *base.Result, inErr error) (result *sdk.Result, err error) {
	if inResult == nil || inResult.ExecResult == nil || inResult.ExecResult.Result == nil {
		return nil, fmt.Errorf("result is nil")
	}
	inResult.ExecResult.Result.Data = inResult.ExecResult.TraceLogs

	return inResult.ExecResult.Result, nil
}

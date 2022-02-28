package txs

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/types"
)

type Tx interface {
	// Prepare convert msg to tx
	Prepare(msg *types.MsgEthereumTx) (err error)

	// SaveTx since the txCount is used by the stateDB, and a simulated tx is run only on the node it's submitted to,
	// then this will cause the txCount/stateDB of the node that ran the simulated tx to be different with the
	// other nodes, causing a consensus error
	SaveTx(msg *types.MsgEthereumTx)

	// Transition execute evm tx
	Transition() (result *base.Result, err error)

	// DecorateResultError when
	DecorateResultError(result *base.Result, err error) (*base.Result, error)

	// EmitEvent emit event
	EmitEvent() *sdk.Result
}

// TransitionEvmTx execute evm transition template
func TransitionEvmTx(tx Tx, msg *types.MsgEthereumTx) (result *sdk.Result, err error) {
	// Prepare convert msg to state transition
	err = tx.Prepare(msg)
	if err != nil {
		return nil, err
	}

	// save tx
	tx.SaveTx(msg)

	// execute transition, the result
	var baseResult *base.Result
	baseResult, err = tx.Transition()
	if err != nil {
		baseResult, err = tx.DecorateResultError(baseResult, err)
		return baseResult.ExecResult.Result, err
	}

	result = tx.EmitEvent()

	return
}

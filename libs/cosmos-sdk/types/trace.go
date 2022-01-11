package types

import (
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

type QueryTraceParams struct {
	// msgEthereumTx for the requested transaction
	TraceTx *Tx
	TxBytes tmtypes.Tx
	// the predecessor transactions included in the same block
	Predecessors      []*Tx
	PredecessorsBytes []tmtypes.Tx
	Block             *tmtypes.Block
}

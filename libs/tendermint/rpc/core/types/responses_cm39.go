package coretypes

import (
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/bytes"
	"github.com/okex/exchain/libs/tendermint/types"
)

// Query abci msg
type CM39ResultABCIQuery struct {
	Response abci.CM39ResponseQuery `json:"response"`
}

// Result of querying for a tx
type CM39ResultTx struct {
	Hash     bytes.HexBytes             `json:"hash"`
	Height   int64                      `json:"height"`
	Index    uint32                     `json:"index"`
	TxResult abci.CM39ResponseDeliverTx `json:"tx_result"`
	Tx       types.Tx                   `json:"tx"`
	Proof    types.TxProof              `json:"proof,omitempty"`
}

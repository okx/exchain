package client

import (
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/libs/bytes"
	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
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

/////////

func ConvCM39ToCM4(c *CM39ResultTx, ret *coretypes.ResultTx) {
	ret.Hash = c.Hash
	ret.Height = c.Height
	ret.Index = c.Index
	ret.TxResult = abci.ResponseDeliverTx{
		Code:      c.TxResult.Code,
		Data:      c.TxResult.Data,
		Log:       c.TxResult.Log,
		Info:      c.TxResult.Info,
		GasWanted: c.TxResult.GasWanted,
		GasUsed:   c.TxResult.GasUsed,
		Events:    c.TxResult.Events,
	}
	ret.Tx = c.Tx
	ret.Proof = c.Proof
}

func ConvTCM39ResultABCIQuery2CM4(c *CM39ResultABCIQuery, ret *coretypes.ResultABCIQuery) {
	ret.Response = abci.ResponseQuery{
		Code:      c.Response.Code,
		Log:       c.Response.Log,
		Info:      c.Response.Info,
		Index:     c.Response.Index,
		Key:       c.Response.Key,
		Value:     c.Response.Value,
		Proof:     c.Response.Proof,
		Height:    c.Response.Height,
		Codespace: c.Response.Codespace,
	}
}

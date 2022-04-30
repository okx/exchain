package client

import (
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
)

func ConvCM39ToCM4(c *coretypes.CM39ResultTx, ret *coretypes.ResultTx) {
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

func ConvTCM39ResultABCIQuery2CM4(c *coretypes.CM39ResultABCIQuery, ret *coretypes.ResultABCIQuery) {
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

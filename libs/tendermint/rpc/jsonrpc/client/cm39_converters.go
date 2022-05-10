package client

import (
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	merkle "github.com/okex/exchain/libs/tendermint/crypto/merkle"
	"github.com/okex/exchain/libs/tendermint/libs/bytes"
	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	"github.com/okex/exchain/libs/tendermint/types"
)

type CM39ResultBroadcastTxCommit struct {
	CheckTx   abci.ResponseCheckTx  `json:"check_tx"`
	DeliverTx CM39ResponseDeliverTx `json:"deliver_tx"`
	Hash      bytes.HexBytes        `json:"hash"`
	Height    int64                 `json:"height"`
}

type CM39ResponseQuery struct {
	Code uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// bytes data = 2; // use "value" instead.
	Log       string        `protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	Info      string        `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	Index     int64         `protobuf:"varint,5,opt,name=index,proto3" json:"index,omitempty"`
	Key       []byte        `protobuf:"bytes,6,opt,name=key,proto3" json:"key,omitempty"`
	Value     []byte        `protobuf:"bytes,7,opt,name=value,proto3" json:"value,omitempty"`
	Proof     *merkle.Proof `protobuf:"bytes,8,opt,name=proof,proto3" json:"proof,omitempty"`
	Height    int64         `protobuf:"varint,9,opt,name=height,proto3" json:"height,omitempty"`
	Codespace string        `protobuf:"bytes,10,opt,name=codespace,proto3" json:"codespace,omitempty"`
}

type CM39ResponseDeliverTx struct {
	Code      uint32       `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Data      []byte       `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Log       string       `protobuf:"bytes,3,opt,name=log,proto3" json:"log,omitempty"`
	Info      string       `protobuf:"bytes,4,opt,name=info,proto3" json:"info,omitempty"`
	GasWanted int64        `protobuf:"varint,5,opt,name=gas_wanted,json=gasWanted,proto3" json:"gas_wanted,omitempty"`
	GasUsed   int64        `protobuf:"varint,6,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	Events    []abci.Event `protobuf:"bytes,7,rep,name=events,proto3" json:"events,omitempty"`
	Codespace string       `protobuf:"bytes,8,opt,name=codespace,proto3" json:"codespace,omitempty"`
}

// Query abci msg
type CM39ResultABCIQuery struct {
	Response CM39ResponseQuery `json:"response"`
}

// Result of querying for a tx
type CM39ResultTx struct {
	Hash     bytes.HexBytes        `json:"hash"`
	Height   int64                 `json:"height"`
	Index    uint32                `json:"index"`
	TxResult CM39ResponseDeliverTx `json:"tx_result"`
	Tx       types.Tx              `json:"tx"`
	Proof    types.TxProof         `json:"proof,omitempty"`
}

/////////

func ConvTCM39BroadcastCommitTx2CM4(c *CM39ResultBroadcastTxCommit, ret *coretypes.ResultBroadcastTxCommit) {
	ret.CheckTx = c.CheckTx
	ret.DeliverTx = abci.ResponseDeliverTx{
		Code:      c.DeliverTx.Code,
		Data:      c.DeliverTx.Data,
		Log:       c.DeliverTx.Log,
		Info:      c.DeliverTx.Info,
		GasWanted: c.DeliverTx.GasWanted,
		GasUsed:   c.DeliverTx.GasUsed,
		Events:    c.DeliverTx.Events,
		Codespace: c.DeliverTx.Codespace,
	}
	ret.Hash = c.Hash
	ret.Height = c.Height
}

func ConvTCM392CM4(c *CM39ResultTx, ret *coretypes.ResultTx) {
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
		Codespace: c.TxResult.Codespace,
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

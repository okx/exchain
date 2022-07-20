package types

import (
	"github.com/ethereum/go-ethereum/common"
)

type QueryTraceTx struct {
	TxHash      common.Hash `json:"tx"`
	ConfigBytes []byte      `json:"config"`
}

type SimulateData struct {
	TxBytes        []byte `json:"tx"`
	OverridesBytes []byte `json:"overrides"`
}
type QueryTraceBlock struct {
	Height      int64  `json:"height"`
	ConfigBytes []byte `json:"config"`
}
type QueryTraceTxResult struct {
	TxIndex int         `json:"tx_index"`
	TxHash  common.Hash `json:"tx_hash"`
	Result  []byte      `json:"result"`
	Error   string      `json:"error"`
}

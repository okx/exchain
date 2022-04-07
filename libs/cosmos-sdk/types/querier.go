package types

import "github.com/ethereum/go-ethereum/common"

type SimulateData struct {
	TxBytes        []byte `json:"tx"`
	OverridesBytes []byte `json:"overrides"`
}

type QueryTraceTx struct {
	TxHash      common.Hash `json:"tx"`
	ConfigBytes []byte      `json:"config"`
}

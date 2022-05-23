package types

import (
	"github.com/ethereum/go-ethereum/common"
)

type QueryTraceTx struct {
	TxHash      common.Hash `json:"tx"`
	ConfigBytes []byte      `json:"config"`
}

type QuerySimulateData struct {
	TxBytes        []byte `json:"tx"`
	OverridesBytes []byte `json:"overrides"`
}

package types

import "github.com/ethereum/go-ethereum/common"

type QueryTraceTx struct {
	TxHash      common.Hash `json:"tx"`
	ConfigBytes []byte      `json:"config"`
}

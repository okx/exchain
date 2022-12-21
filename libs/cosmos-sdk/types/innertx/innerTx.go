package innertx

import (
	"math/big"
)

const (
	CosmosCallType = "cosmos"
	CosmosDepth    = 0

	SendCallName       = "send"
	DelegateCallName   = "delegate"
	MultiCallName      = "multi-send"
	UndelegateCallName = "undelegate"
	EvmCallName        = "call"
	EvmCreateName      = "create"

	IsAvailable = false
)

var BIG0 = big.NewInt(0)

type InnerTxKeeper interface {
	InitInnerBlock(...interface{})
	UpdateInnerTx(...interface{})
	UpdateWasmInnerTx(...interface{})
}

func AddDefaultInnerTx(...interface{}) interface{} {
	return nil
}

func UpdateDefaultInnerTx(...interface{}) {
}

func ParseInnerTxAndContract(...interface{}) (interface{}, interface{}) {
	return nil, nil
}

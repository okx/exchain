package innertx

import (
	"math/big"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

const (
	COSMOS_CALL_TYPE = "cosmos"
	COSMOS_DEPTH     = 0

	SEND_CALL_NAME       = "send"
	DELEGATE_CALL_NAME   = "delegate"
	MULTI_CALL_NAME      = "multi-send"
	UNDELEGATE_CALL_NAME = "undelegate"
	EVM_CALL_NAME        = "call"
	EVM_CREATE_NAME      = "create"
)

var BIG0 = big.NewInt(0)

type InnerTxKeeper interface {
	InitInnerBlock(hash string)
	UpdateInnerTx(txBytes []byte, dept int64, from, to sdk.AccAddress, callType, name string, valueWei sdk.Coins, err error)
}

func AddDefaultInnerTx(...interface{}) interface{} {
	return nil
}

func UpdateDefaultInnerTx(...interface{}) {
}

func ParseInnerTxAndContract(...interface{}) (interface{}, interface{}) {
	return nil, nil
}

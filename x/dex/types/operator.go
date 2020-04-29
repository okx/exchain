package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//OperatorAddress sdk.ValAddress `json:"operator_address" yaml:"operator_address"`

type DEXOperator struct {
	Address            sdk.AccAddress `json:"address"`
	HandlingFeeAddress sdk.AccAddress `json:"handling_fee_address"`
	Website            string         `json:"website"`
	InitHeight         int64          `json:"init_height"`
	TxHash             string         `json:"tx_hash"`
}
package types

import (
	"fmt"
	"strings"

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

// nolint
func (o DEXOperator) String() string {
	return fmt.Sprintf(`DEXOperator :
  Address:              %s
  Handling Fee Address: %s
  Website:              %s
  Init Height:          %d
  TxHash:               %s`,
		o.Address, o.HandlingFeeAddress, o.Website,
		o.InitHeight, o.TxHash,
	)
}

type DEXOperators []DEXOperator

func (o DEXOperators) String() string {
	out := ""
	for _, v := range o {
		out += fmt.Sprintf(`DEXOperator :
  Address:              %s
  Handling Fee Address: %s
  Website:              %s
  Init Height:          %d
  TxHash:               %s`,
			v.Address, v.HandlingFeeAddress, v.Website,
			v.InitHeight, v.TxHash,
		)
	}

	return strings.TrimSpace(out)
}

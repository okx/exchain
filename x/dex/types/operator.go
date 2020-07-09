package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OperatorAddress sdk.ValAddress `json:"operator_address" yaml:"operator_address"`
type DEXOperator struct {
	Address            sdk.AccAddress `json:"address"`
	HandlingFeeAddress sdk.AccAddress `json:"handling_fee_address"`
	Website            string         `json:"website"`
	InitHeight         int64          `json:"init_height"`
	TxHash             string         `json:"tx_hash"`
}

// DEXOperatorInfo includes handling_fee
type DEXOperatorInfo struct {
	Address            sdk.AccAddress `json:"address"`
	HandlingFeeAddress sdk.AccAddress `json:"handling_fee_address"`
	Website            string         `json:"website"`
	InitHeight         int64          `json:"init_height"`
	TxHash             string         `json:"tx_hash"`
	//HandlingFees       string         `json:"handling_fees"`
}

// nolint
func NewDEXOperatorInfo(operator DEXOperator) DEXOperatorInfo {
	info := DEXOperatorInfo{}
	info.Address = operator.Address
	info.HandlingFeeAddress = operator.HandlingFeeAddress
	info.Website = operator.Website
	info.InitHeight = operator.InitHeight
	info.TxHash = operator.TxHash

	return info
}

// nolint
func (o DEXOperatorInfo) String() string {
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

// nolint
type DEXOperatorInfos []DEXOperatorInfo

// nolint
func (o DEXOperatorInfos) String() string {
	out := ""
	for _, v := range o {
		out += v.String()
	}
	return strings.TrimSpace(out)
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
		out += v.String()
	}

	return strings.TrimSpace(out)
}

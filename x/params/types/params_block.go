package types

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params/subspace"
)

const (
	QueryBlockConfig   = "blockconfig"
	MaxGasUsedPerBlock = "MaxGasUsedPerBlock"
)

// BlockConfig is the struct of the parameters in this module
type BlockConfig struct {
	MaxGasUsedPerBlock int64 `json:"maxGasUsedPerBlock"`
}

func NewDefaultBlockConfig() *BlockConfig {
	return &BlockConfig{
		MaxGasUsedPerBlock: sdk.DefaultMaxGasUsedPerBlock,
	}
}

func (p BlockConfig) String() string {
	return fmt.Sprintf(`
MaxGasUsedPerBlock: %d,
`, p.MaxGasUsedPerBlock)
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *BlockConfig) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{[]byte(MaxGasUsedPerBlock), &p.MaxGasUsedPerBlock, ValidateInt64("maxGasUsedPerBlock")},
	}
}

func ValidateInt64(param string) subspace.ValueValidatorFn {
	return func(i interface{}) error {
		v, ok := i.(int64)
		if !ok {
			return fmt.Errorf("invalid parameter type: %T", i)
		}

		if v < -1 {
			return fmt.Errorf("%s must be equal or greater than -1: %d", param, v)
		}

		return nil
	}
}

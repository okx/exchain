package types

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params/subspace"
	"github.com/okex/exchain/x/common"

	stypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdkparams "github.com/okex/exchain/libs/cosmos-sdk/x/params"
)

// GasConfigKeyTable returns the key declaration for parameters
func GasConfigKeyTable() sdkparams.KeyTable {
	return sdkparams.NewKeyTable().RegisterParamSet(&GasConfig{})
}

// GasConfig is the struct of the parameters in this module
type GasConfig struct {
	ReadCostFlat    uint64 `json:"read_cost_flat"`
	ReadCostPerByte uint64 `json:"Read_cost_perByte"`
}

func (p GasConfig) String() string {
	return fmt.Sprintf(`
ReadCostFlat: %d,
ReadCostPerByte:       %d,
`, p.ReadCostFlat, p.ReadCostPerByte)
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *GasConfig) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{[]byte(stypes.GasReadCostFlatDesc), &p.ReadCostFlat, common.ValidateUint64Positive("gas read cost flat")},
		{[]byte(stypes.GasReadPerByteDesc), &p.ReadCostPerByte, common.ValidateUint64Positive("gas read per byte")},
	}
}

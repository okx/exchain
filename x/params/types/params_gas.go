package types

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params/subspace"
	"github.com/okex/exchain/x/common"

	stypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
)

const (
	QueryGasConfig = "gasconfig"
)

// GasConfig is the struct of the parameters in this module
type GasConfig struct {
	stypes.GasConfig
}

func (p GasConfig) String() string {
	return fmt.Sprintf(`
HasCost: %d,
DeleteCost: %d,
ReadCostFlat: %d,
ReadCostPerByte: %d,
WriteCostFlat: %d,
WriteCostPerByte: %d,
IterNextCostFlat: %d,
`, p.HasCost, p.DeleteCost, p.ReadCostFlat, p.ReadCostPerByte, p.WriteCostFlat, p.WriteCostPerByte, p.IterNextCostFlat)
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *GasConfig) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{[]byte(stypes.GasHasDesc), &p.HasCost, common.ValidateUint64Positive("gas has")},
		{[]byte(stypes.GasDeleteDesc), &p.DeleteCost, common.ValidateUint64Positive("gas delete")},
		{[]byte(stypes.GasReadCostFlatDesc), &p.ReadCostFlat, common.ValidateUint64Positive("gas read cost flat")},
		{[]byte(stypes.GasReadPerByteDesc), &p.ReadCostPerByte, common.ValidateUint64Positive("gas read per byte")},
		{[]byte(stypes.GasWriteCostFlatDesc), &p.WriteCostFlat, common.ValidateUint64Positive("gas write cost flat")},
		{[]byte(stypes.GasWritePerByteDesc), &p.WriteCostPerByte, common.ValidateUint64Positive("gas write cost per byte")},
		{[]byte(stypes.GasIterNextCostFlatDesc), &p.IterNextCostFlat, common.ValidateUint64Positive("gas iter next cost flat")},
	}
}

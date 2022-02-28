package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	simtypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/types/simulation"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/simulation"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/staking/types"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyMaxValidators),
			func(r *rand.Rand) string {
				return fmt.Sprintf("%d", genMaxValidators(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyUnbondingTime),
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", genUnbondingTime(r))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyHistoricalEntries),
			func(r *rand.Rand) string {
				return fmt.Sprintf("%d", getHistEntries(r))
			},
		),
	}
}

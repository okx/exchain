package simulation

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/application/transfer/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/simulation"
	gogotypes "github.com/gogo/protobuf/types"
	"math/rand"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(r *rand.Rand) []simulation.ParamChange{
	return []simulation.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, string(types.KeySendEnabled),
			func(r *rand.Rand) string {
				sendEnabled := RadomEnabled(r)
				return string(types.ModuleCdc.MustMarshalJSON(&gogotypes.BoolValue{Value: sendEnabled}))
			},
		),
		simulation.NewSimParamChange(types.ModuleName, string(types.KeyReceiveEnabled),
			func(r *rand.Rand) string {
				receiveEnabled := RadomEnabled(r)
				return string(types.ModuleCdc.MustMarshalJSON(&gogotypes.BoolValue{Value: receiveEnabled}))
			},
		),
	}
}

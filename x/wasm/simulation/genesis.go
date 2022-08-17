package simulation

import (
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"

	"github.com/okex/exchain/x/wasm/types"
)

// RandomizeGenState generates a random GenesisState for wasm
func RandomizedGenState(simstate *module.SimulationState) {
	params := RandomParams(simstate.Rand)
	wasmGenesis := types.GenesisState{
		Params:    params,
		Codes:     nil,
		Contracts: nil,
		Sequences: []types.Sequence{
			{IDKey: types.KeyLastCodeID, Value: simstate.Rand.Uint64()},
			{IDKey: types.KeyLastInstanceID, Value: simstate.Rand.Uint64()},
		},
		GenMsgs: nil,
	}

	_, err := simstate.Cdc.MarshalJSON(&wasmGenesis)
	if err != nil {
		panic(err)
	}

	simstate.GenState[types.ModuleName] = simstate.Cdc.MustMarshalJSON(&wasmGenesis)
}

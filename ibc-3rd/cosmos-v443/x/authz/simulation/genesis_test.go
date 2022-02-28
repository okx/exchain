package simulation_test

import (
	"encoding/json"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/ibc-3rd/cosmos-v443/simapp"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/types/module"
	simtypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/types/simulation"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/authz"
	"github.com/okex/exchain/ibc-3rd/cosmos-v443/x/authz/simulation"
)

func TestRandomizedGenState(t *testing.T) {
	app := simapp.Setup(false)

	s := rand.NewSource(1)
	r := rand.New(s)

	simState := module.SimulationState{
		AppParams:    make(simtypes.AppParams),
		Cdc:          app.AppCodec(),
		Rand:         r,
		NumBonded:    3,
		Accounts:     simtypes.RandomAccounts(r, 3),
		InitialStake: 1000,
		GenState:     make(map[string]json.RawMessage),
	}

	simulation.RandomizedGenState(&simState)
	var authzGenesis authz.GenesisState
	simState.Cdc.MustUnmarshalJSON(simState.GenState[authz.ModuleName], &authzGenesis)

	require.Len(t, authzGenesis.Authorization, len(simState.Accounts)-1)
}

package simulation

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/simulation"
	"math/rand"
)

// GenClientGenesis returns the default client genesis state.
func GenClientGenesis(_ *rand.Rand, _ []simulation.Account) types.GenesisState {
	return types.DefaultGenesisState()
}

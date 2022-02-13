package simulation

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/03-connection/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/simulation"
	"github.com/okex/exchain/libs/tendermint/libs/rand"
)

// GenConnectionGenesis returns the default connection genesis state.
func GenConnectionGenesis(_ *rand.Rand, _ []simulation.Account) types.GenesisState {
	return types.DefaultGenesisState()
}

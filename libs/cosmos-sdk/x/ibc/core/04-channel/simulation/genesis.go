package simulation

import (
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/simulation"
	"github.com/okex/exchain/libs/tendermint/libs/rand"
)

// GenChannelGenesis returns the default channel genesis state.
func GenChannelGenesis(_ *rand.Rand, _ []simulation.Account) types.GenesisState {
	return types.DefaultGenesisState()
}

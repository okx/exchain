package simulation

import (
	"encoding/json"
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	clienttypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/02-client/types"
	connectiontypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/03-connection/types"
	channeltypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/04-channel/types"
	clientsims "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/02-client/simulation"
	connectionsims "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/03-connection/simulation"
	channelsims "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/04-channel/simulation"
	host "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/24-host"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/types"
	"math/rand"
)

// DONTCOVER


// Simulation parameter constants
const (
	clientGenesis     = "client_genesis"
	connectionGenesis = "connection_genesis"
	channelGenesis    = "channel_genesis"
)

// RandomizedGenState generates a random GenesisState for evidence
func RandomizedGenState(simState *module.SimulationState) {
	var (
		clientGenesisState     clienttypes.GenesisState
		connectionGenesisState connectiontypes.GenesisState
		channelGenesisState    channeltypes.GenesisState
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, clientGenesis, &clientGenesisState, simState.Rand,
		func(r *rand.Rand) { clientGenesisState = clientsims.GenClientGenesis(r, simState.Accounts) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, connectionGenesis, &connectionGenesisState, simState.Rand,
		func(r *rand.Rand) { connectionGenesisState = connectionsims.GenConnectionGenesis(r, simState.Accounts) },
	)

	simState.AppParams.GetOrGenerate(
		simState.Cdc, channelGenesis, &channelGenesisState, simState.Rand,
		func(r *rand.Rand) { channelGenesisState = channelsims.GenChannelGenesis(r, simState.Accounts) },
	)

	ibcGenesis := types.GenesisState{
		ClientGenesis:     clientGenesisState,
		ConnectionGenesis: connectionGenesisState,
		ChannelGenesis:    channelGenesisState,
	}

	bz, err := json.MarshalIndent(&ibcGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated %s parameters:\n%s\n", host.ModuleName, bz)
	simState.GenState[host.ModuleName] = simState.Cdc.MustMarshalJSON(&ibcGenesis)
}

package staking

import (
	"encoding/json"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
	"github.com/okex/exchain/libs/ibc-go/testing/simapp/adapter"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/staking"
	"github.com/okex/exchain/x/staking/types"
)

type StakingModule struct {
	staking.AppModule

	keeper       staking.Keeper
	accKeeper    types.AccountKeeper
	supplyKeeper types.SupplyKeeper
}

func TNewStakingModule(keeper staking.Keeper, accKeeper types.AccountKeeper,
	supplyKeeper types.SupplyKeeper) StakingModule {

	s := staking.NewAppModule(keeper, accKeeper, supplyKeeper)

	return StakingModule{
		s,
		keeper,
		accKeeper,
		supplyKeeper,
	}
}

func (s StakingModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState staking.GenesisState
	adapter.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	genesisState.Params = types.DefaultParams()
	genesisState.Params.UnbondingTime = types2.DefaultUnbondingTime

	return staking.InitGenesis(ctx, s.keeper, s.accKeeper, s.supplyKeeper, genesisState)
}

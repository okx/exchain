package ica

import (
	"encoding/json"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"

	"github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	ica "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts"
	controllerkeeper "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/controller/keeper"
	hostkeeper "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/host/keeper"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

type TestICAModuleBaisc struct {
	ica.AppModuleBasic
}

func (b TestICAModuleBaisc) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(types.DefaultGenesis())
}

type TestICAModule struct {
	ica.AppModule
	ck *controllerkeeper.Keeper
	hk *hostkeeper.Keeper
}

func NewTestICAModule(cdc *codec.CodecProxy, ck *controllerkeeper.Keeper, hk *hostkeeper.Keeper) *TestICAModule {
	return &TestICAModule{
		AppModule: ica.NewAppModule(cdc, ck, hk),
		ck:        ck,
		hk:        hk,
	}
}

func (am TestICAModule) InitGenesis(s sdk.Context, message json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(message, &genesisState)

	if am.ck != nil {
		controllerkeeper.InitGenesis(s, *am.ck, genesisState.ControllerGenesisState)
	}

	if am.hk != nil {
		hostkeeper.InitGenesis(s, *am.hk, genesisState.HostGenesisState)
	}

	return nil
}

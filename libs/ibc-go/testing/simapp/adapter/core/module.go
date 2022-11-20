package core

import (
	"encoding/json"
	"fmt"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	ibc "github.com/okex/exchain/libs/ibc-go/modules/core"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	"github.com/okex/exchain/libs/ibc-go/modules/core/keeper"
	"github.com/okex/exchain/libs/ibc-go/modules/core/types"
	"github.com/okex/exchain/libs/ibc-go/testing/simapp/adapter"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

type CoreModule struct {
	ibc.AppModule

	tkeeper *keeper.Keeper

	// create localhost by default
	tcreateLocalhost bool
}

// NewAppModule creates a new AppModule object
func NewIBCCOreAppModule(k *ibc.Keeper) *CoreModule {
	a := ibc.NewAppModule(k)
	ret := &CoreModule{
		AppModule:        a,
		tkeeper:          k.V2Keeper,
		tcreateLocalhost: false,
	}
	return ret
}

// DefaultGenesis returns default genesis state as raw bytes for the ibc
// module.
func (CoreModule) DefaultGenesis() json.RawMessage {
	return adapter.ModuleCdc.MustMarshalJSON(ibc.DefaultGenesisState())
}

// InitGenesis performs genesis initialization for the ibc module. It returns
// no validator updates.
//func (am CoreModule) InitGenesis(ctx sdk.Context, cdc Corec.JSONMarshaler, bz json.RawMessage) []abci.ValidatorUpdate {
func (am CoreModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return am.initGenesis(ctx, data)
}

func (am CoreModule) initGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var gs types.GenesisState
	err := adapter.ModuleCdc.UnmarshalJSON(data, &gs)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal %s genesis state: %s", host.ModuleName, err))
	}
	gs.Params.EnableIbc = true
	ibc.InitGenesis(ctx, *am.tkeeper, am.tcreateLocalhost, &gs)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the ibc
// module.
func (am CoreModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return am.exportGenesis(ctx)
}

func (am CoreModule) exportGenesis(ctx sdk.Context) json.RawMessage {
	return adapter.ModuleCdc.MustMarshalJSON(ibc.ExportGenesis(ctx, *am.tkeeper))
}

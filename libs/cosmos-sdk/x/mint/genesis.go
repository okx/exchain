package mint

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// InitGenesis new mint genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetMinter(ctx, data.Minter)
	keeper.SetParams(ctx, data.Params)
	if data.Treasures != nil {
		keeper.SetTreasures(ctx, data.Treasures)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	minter := keeper.GetMinterCustom(ctx)
	params := keeper.GetParams(ctx)
	genesisState := NewGenesisState(minter, params, keeper.GetOriginalMintedPerBlock())
	treasures := keeper.GetTreasures(ctx)
	if treasures != nil {
		genesisState.Treasures = treasures
	}
	return genesisState
}

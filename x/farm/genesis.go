package farm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/keeper"
	"github.com/okex/okexchain/x/farm/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k keeper.Keeper /* TODO: Define what keepers the module needs */, data types.GenesisState) {
	// TODO: Define logic for when you would like to initialize a new genesis
	////////////////////////////////////////////////////////////
	// TODO: demo for test. remove it later
	tPool := types.NewFarmPool(
		"pool-airtoken1-eth",
		"locked_token_symbol",
		types.YieldedTokenInfos{
			types.NewYieldedTokenInfo(
				sdk.NewDecCoinFromDec("btc", sdk.OneDec()),
				1024,
				sdk.OneDec(),
			)},
		sdk.NewDecCoinFromDec("btc", sdk.OneDec()),
		sdk.Coins{sdk.NewDecCoinFromDec("btc", sdk.OneDec())},
		2048,
		sdk.OneDec(),
	)
	k.SetFarmPool(ctx, tPool)
	////////////////////////////////////////////////////////////
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) (data types.GenesisState) {
	// TODO: Define logic for exporting state
	return types.NewGenesisState()
}

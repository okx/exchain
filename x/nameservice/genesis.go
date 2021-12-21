package nameservice

import (
	// abci "github.com/tendermint/tendermint/abci/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/nameservice/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	for _, record := range data.WhoisRecords {
		keeper.SetWhois(ctx, record.Value, record)
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k Keeper) types.GenesisState {
	var records []types.Whois
	iterator := k.GetNamesIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {

		name := string(iterator.Key())
		whois, _ := k.GetWhois(ctx, name)
		records = append(records, whois)
	}
	return types.GenesisState{WhoisRecords: records}
}

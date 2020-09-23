package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

// NewQuerier creates a new querier for farm clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, k)
			// TODO: Put the modules query routes
		default:
			return nil, sdk.ErrUnknownRequest("unknown farm query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {

	return []byte{}, nil
}

// TODO: Add the modules query functions
// They will be similar to the above one: queryParams()

package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okexchain/x/token"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/farm/types"
)

// Keeper of the farm store
type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	paramspace   types.ParamSubspace
	supplyKeeper supply.Keeper
	tokenKeeper  token.Keeper
}

// NewKeeper creates a farm keeper
func NewKeeper(supplyKeeper supply.Keeper, tokenKeeper token.Keeper, paramspace types.ParamSubspace, key sdk.StoreKey,
	cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		paramspace:   paramspace.WithKeyTable(types.ParamKeyTable()),
		supplyKeeper: supplyKeeper,
		tokenKeeper:  tokenKeeper,
	}
}

func (k Keeper) StoreKey() sdk.StoreKey {
	return k.storeKey
}

func (k Keeper) SupplyKeeper() supply.Keeper {
	return k.supplyKeeper
}

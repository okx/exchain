package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okexchain/x/token"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	swap "github.com/okex/okexchain/x/ammswap/keeper"
	"github.com/okex/okexchain/x/farm/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the farm store
type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	feeCollectorName string // name of the FeeCollector ModuleAccount
	paramSubspace    types.ParamSubspace
	supplyKeeper     supply.Keeper
	tokenKeeper      token.Keeper
	swapKeeper       swap.Keeper
}

// NewKeeper creates a farm keeper
func NewKeeper(feeCollectorName string, supplyKeeper supply.Keeper,
	tokenKeeper token.Keeper,
	swapKeeper swap.Keeper,
	paramSubspace types.ParamSubspace, key sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		storeKey:         key,
		cdc:              cdc,
		feeCollectorName: feeCollectorName,
		paramSubspace:    paramSubspace.WithKeyTable(types.ParamKeyTable()),
		supplyKeeper:     supplyKeeper,
		tokenKeeper:      tokenKeeper,
		swapKeeper:       swapKeeper,
	}
}

func (k Keeper) StoreKey() sdk.StoreKey {
	return k.storeKey
}

func (k Keeper) SupplyKeeper() supply.Keeper {
	return k.supplyKeeper
}

func (k Keeper) TokenKeeper() token.Keeper {
	return k.tokenKeeper
}

// GetFeeCollector returns feeCollectorName
func (k Keeper) GetFeeCollector() string {
	return k.feeCollectorName
}

// Logger returns a module-specific logger.
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

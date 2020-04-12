package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/debug/types"
)

// keeper of debug module
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec
	// for test for sending fee
	feePoolModuleName string
	stakingKeeper     types.StakingKeeper
	supplyKeeper      types.SupplyKeeper
	tokenKeeper       types.TokenKeeper
	orderKeeper       types.OrderKeeper
	StopFunc          func()
}

func NewDebugKeeper(cdc *codec.Codec, storeKey sdk.StoreKey,
	orderKeeper types.OrderKeeper,
	stakingKeeper types.StakingKeeper,
	tokenKeeper types.TokenKeeper,
	supplyKeeper types.SupplyKeeper, feePoolModuleName string, stop func()) Keeper {
	return Keeper{
		storeKey:          storeKey,
		cdc:               cdc,
		feePoolModuleName: feePoolModuleName,
		orderKeeper:       orderKeeper,
		stakingKeeper:     stakingKeeper,
		supplyKeeper:      supplyKeeper,
		tokenKeeper:       tokenKeeper,
		StopFunc:          stop,
	}
}

func (k *Keeper) GetCDC() *codec.Codec {
	return k.cdc
}

func (k *Keeper) DumpStore(ctx sdk.Context, m string) {
	logger := ctx.Logger().With("module", "debug")

	logger.Error(fmt.Sprintf("--------------------------------------------------------------------------"))
	logger.Error(fmt.Sprintf("---------- Dump <%s> KV Store at BlockHeight <%d> started ----------",
		m, ctx.BlockHeight()))
	defer logger.Error(fmt.Sprintf("--------------------------------------------------------------------------"))
	defer logger.Error(fmt.Sprintf("---------- Dump <%s> KV Store at BlockHeight <%d> finished ----------",
		m, ctx.BlockHeight()))

	if m == "order" {
		k.orderKeeper.DumpStore(ctx)
	}
}

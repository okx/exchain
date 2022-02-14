package keeper

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	capabilitykeeper "github.com/okex/exchain/libs/cosmos-sdk/x/capability/keeper"
	clientkeeper "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/02-client/keeper"
	clienttypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/02-client/types"
	connectionkeeper "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/03-connection/keeper"
	channelkeeper "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/04-channel/keeper"
	portkeeper "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/05-port/keeper"
	porttypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/05-port/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/types"
	paramtypes "github.com/okex/exchain/libs/cosmos-sdk/x/params"
)

var _ types.QueryServer = (*Keeper)(nil)

// Keeper defines each ICS keeper for IBC
type Keeper struct {
	// implements gRPC QueryServer interface
	types.QueryServer

	cdc codec.Codec

	ClientKeeper     clientkeeper.Keeper
	ConnectionKeeper connectionkeeper.Keeper
	ChannelKeeper    channelkeeper.Keeper
	PortKeeper       portkeeper.Keeper
	Router           *porttypes.Router
}

// NewKeeper creates a new ibc Keeper
func NewKeeper(
	cdc codec.Codec, key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper clienttypes.StakingKeeper, scopedKeeper capabilitykeeper.ScopedKeeper,
) *Keeper {
	clientKeeper := clientkeeper.NewKeeper(cdc, key, paramSpace, stakingKeeper)
	connectionKeeper := connectionkeeper.NewKeeper(cdc, key, clientKeeper)
	portKeeper := portkeeper.NewKeeper(scopedKeeper)
	channelKeeper := channelkeeper.NewKeeper(cdc, key, clientKeeper, connectionKeeper, portKeeper, scopedKeeper)

	return &Keeper{
		cdc:              cdc,
		ClientKeeper:     clientKeeper,
		ConnectionKeeper: connectionKeeper,
		ChannelKeeper:    channelKeeper,
		PortKeeper:       portKeeper,
	}
}

// Codec returns the IBC module codec.
func (k Keeper) Codec() codec.Codec {
	return k.cdc
}

// SetRouter sets the Router in IBC Keeper and seals it. The method panics if
// there is an existing router that's already sealed.
func (k *Keeper) SetRouter(rtr *porttypes.Router) {
	if k.Router != nil && k.Router.Sealed() {
		panic("cannot reset a sealed router")
	}
	k.Router = rtr
	k.Router.Seal()
}

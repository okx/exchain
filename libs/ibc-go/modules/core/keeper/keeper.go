package keeper

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	capabilitykeeper "github.com/okex/exchain/libs/cosmos-sdk/x/capability/keeper"
	paramtypes "github.com/okex/exchain/libs/cosmos-sdk/x/params"
	clientkeeper "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/keeper"
	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	connectionkeeper "github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/keeper"
	connectiontypes "github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/types"
	channelkeeper "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/keeper"
	portkeeper "github.com/okex/exchain/libs/ibc-go/modules/core/05-port/keeper"
	porttypes "github.com/okex/exchain/libs/ibc-go/modules/core/05-port/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/types"
)

var _ types.QueryServer = (*Keeper)(nil)

// Keeper defines each ICS keeper for IBC
type Keeper struct {
	// implements gRPC QueryServer interface
	types.QueryServer

	cdc *codec.MarshalProxy

	ClientKeeper     clientkeeper.Keeper
	ConnectionKeeper connectionkeeper.Keeper
	ChannelKeeper    channelkeeper.Keeper
	PortKeeper       portkeeper.Keeper
	Router           *porttypes.Router
}

// NewKeeper creates a new ibc Keeper
func NewKeeper(
	proxy *codec.MarshalProxy,
	key sdk.StoreKey, paramSpace paramtypes.Subspace,
	stakingKeeper clienttypes.StakingKeeper, scopedKeeper *capabilitykeeper.ScopedKeeper,
	registry types2.InterfaceRegistry,
) *Keeper {
	//mm := codec.NewProtoCodec(registry)
	//proxy:=codec.NewMarshalProxy(mm,cdcc)
	if !paramSpace.HasKeyTable() {
		keyTable := clienttypes.ParamKeyTable()
		keyTable.RegisterParamSet(&connectiontypes.Params{})
		paramSpace = paramSpace.WithKeyTable(keyTable)
	}
	clientKeeper := clientkeeper.NewKeeper(proxy, key, paramSpace, stakingKeeper)
	connectionKeeper := connectionkeeper.NewKeeper(proxy, key, paramSpace, clientKeeper)
	portKeeper := portkeeper.NewKeeper(scopedKeeper)
	channelKeeper := channelkeeper.NewKeeper(proxy, key, clientKeeper, connectionKeeper, portKeeper, scopedKeeper)

	return &Keeper{
		cdc:              proxy,
		ClientKeeper:     clientKeeper,
		ConnectionKeeper: connectionKeeper,
		ChannelKeeper:    channelKeeper,
		PortKeeper:       portKeeper,
	}
}

// Codec returns the IBC module codec.
func (k Keeper) Codec() *codec.MarshalProxy {
	return k.cdc
}

// SetRouter sets the Router in IBC Keeper and seals it. The method panics if
// there is an existing router that's already sealed.
func (k *Keeper) SetRouter(rtr *porttypes.Router) {
	if k.Router != nil && k.Router.Sealed() {
		panic("cannot reset a sealed router")
	}

	k.PortKeeper.Router = rtr
	k.Router = rtr
	k.Router.Seal()
}

package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	supplyexported "github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	ibcexported "github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

// AccountKeeper defines the expected account keeper
type AccountKeeper interface {
	NewAccount(ctx sdk.Context, acc authtypes.Account) authtypes.Account
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) authtypes.Account
	SetAccount(ctx sdk.Context, acc authtypes.Account)
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI
	GetModuleAddress(name string) sdk.AccAddress
}

// ICS4Wrapper defines the expected ICS4Wrapper for middleware
type ICS4Wrapper interface {
	SendPacket(ctx sdk.Context, channelCap *capabilitytypes.Capability, packet ibcexported.PacketI) error
	GetAppVersion(ctx sdk.Context, portID, channelID string) (string, bool)
}

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channeltypes.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	GetConnection(ctx sdk.Context, connectionID string) (ibcexported.ConnectionI, error)
}

// PortKeeper defines the expected IBC port keeper
type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability
	IsBound(ctx sdk.Context, portID string) bool
}

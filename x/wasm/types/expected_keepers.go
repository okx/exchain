package types

import (
	"context"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	auth "github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/okex/exchain/libs/cosmos-sdk/x/staking/types"
	connectiontypes "github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/types"
	channeltypes "github.com/okex/exchain/libs/ibc-go/modules/core/04-channel/types"
	ibcexported "github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

// BankViewKeeper defines a subset of methods implemented by the cosmos-sdk bank keeper
type BankViewKeeper interface {
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

// BankKeeper defines a subset of methods implemented by the cosmos-sdk bank keeper
type BankKeeper interface {
	BankViewKeeper
	//	Burner
	IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error
	BlockedAddr(addr sdk.AccAddress) bool
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error
	GetSendEnabled(ctx sdk.Context) bool
}

// AccountKeeper defines a subset of methods implemented by the cosmos-sdk account keeper
type AccountKeeper interface {
	// Return a new account with the next account number and the specified address. Does not save the new account to the store.
	NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) auth.Account
	// Retrieve an account from the store.
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) auth.Account
	// Set an account in the store.
	SetAccount(ctx sdk.Context, acc auth.Account, updateState ...bool)
}

// DistributionKeeper defines a subset of methods implemented by the cosmos-sdk distribution keeper
type DistributionKeeper interface {
	DelegationRewards(c context.Context, req *types.QueryDelegationRewardsParams) (*sdk.DecCoins, error)
}

// StakingKeeper defines a subset of methods implemented by the cosmos-sdk staking keeper
type StakingKeeper interface {
	// BondDenom - Bondable coin denomination
	BondDenom(ctx sdk.Context) (res string)
	// GetValidator get a single validator
	GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, found bool)
	// GetBondedValidatorsByPower get the current group of bonded validators sorted by power-rank
	GetBondedValidatorsByPower(ctx sdk.Context) []stakingtypes.Validator
	// GetAllDelegatorDelegations return all delegations for a delegator
	GetAllDelegatorDelegations(ctx sdk.Context, delegator sdk.AccAddress) []stakingtypes.Delegation
	// GetDelegation return a specific delegation
	GetDelegation(ctx sdk.Context,
		delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, found bool)
	// HasReceivingRedelegation check if validator is receiving a redelegation
	HasReceivingRedelegation(ctx sdk.Context,
		delAddr sdk.AccAddress, valDstAddr sdk.ValAddress) bool
}

// ChannelKeeper defines the expected IBC channel keeper
type ChannelKeeper interface {
	GetChannel(ctx sdk.Context, srcPort, srcChan string) (channel channeltypes.Channel, found bool)
	GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool)
	SendPacket(ctx sdk.Context, channelCap *capabilitytypes.Capability, packet ibcexported.PacketI) error
	ChanCloseInit(ctx sdk.Context, portID, channelID string, chanCap *capabilitytypes.Capability) error
	GetAllChannels(ctx sdk.Context) (channels []channeltypes.IdentifiedChannel)
	IterateChannels(ctx sdk.Context, cb func(channeltypes.IdentifiedChannel) bool)
	SetChannel(ctx sdk.Context, portID, channelID string, channel channeltypes.Channel)
}

// ClientKeeper defines the expected IBC client keeper
type ClientKeeper interface {
	GetClientConsensusState(ctx sdk.Context, clientID string) (connection ibcexported.ConsensusState, found bool)
}

// ConnectionKeeper defines the expected IBC connection keeper
type ConnectionKeeper interface {
	GetConnection(ctx sdk.Context, connectionID string) (connection connectiontypes.ConnectionEnd, found bool)
}

// PortKeeper defines the expected IBC port keeper
type PortKeeper interface {
	BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability
}

type CapabilityKeeper interface {
	GetCapability(ctx sdk.Context, name string) (*capabilitytypes.Capability, bool)
	ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error
	AuthenticateCapability(ctx sdk.Context, capability *capabilitytypes.Capability, name string) bool
}

// ICS20TransferPortSource is a subset of the ibc transfer keeper.
type ICS20TransferPortSource interface {
	GetPort(ctx sdk.Context) string
}

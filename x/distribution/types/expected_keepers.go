package types

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	supplyexported "github.com/okx/okbchain/libs/cosmos-sdk/x/supply/exported"
	"github.com/okx/okbchain/x/common"
	govtypes "github.com/okx/okbchain/x/gov/types"
	stakingexported "github.com/okx/okbchain/x/staking/exported"
)

// StakingKeeper expected staking keeper (noalias)
type StakingKeeper interface {
	// iterate through validators by operator address, execute func for each validator
	IterateValidators(sdk.Context,
		func(index int64, validator stakingexported.ValidatorI) (stop bool))

	// iterate through bonded validators by operator address, execute func for each validator
	IterateBondedValidatorsByPower(sdk.Context,
		func(index int64, validator stakingexported.ValidatorI) (stop bool))

	// iterate through the consensus validator set of the last block by operator address
	// execute func for each validator
	IterateLastValidators(sdk.Context,
		func(index int64, validator stakingexported.ValidatorI) (stop bool))

	// get a particular validator by operator address
	Validator(sdk.Context, sdk.ValAddress) stakingexported.ValidatorI
	// get a particular validator by consensus address
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) stakingexported.ValidatorI

	Jail(sdk.Context, sdk.ConsAddress)   // jail a validator
	Unjail(sdk.Context, sdk.ConsAddress) // unjail a validator

	// MaxValidators returns the maximum amount of bonded validators
	MaxValidators(sdk.Context) uint16

	GetLastTotalPower(ctx sdk.Context) sdk.Int
	GetLastValidatorPower(ctx sdk.Context, valAddr sdk.ValAddress) int64

	Delegator(ctx sdk.Context, delAddr sdk.AccAddress) stakingexported.DelegatorI

	IsValidator(ctx sdk.Context, addr sdk.AccAddress) bool

	ParamsConsensusType(ctx sdk.Context) (consensusType common.ConsensusType)
}

// StakingHooks event hooks for staking validator object (noalias)
type StakingHooks interface {
	// Must be called when a validator is created
	AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress)
	// Must be called when a validator is deleted
	AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress)
	// Must be called when a delegation is created
	BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddrs []sdk.ValAddress)
	// Must be called when a delegation's shares are modified
	BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress)
	AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress)
	BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec)
}

// SupplyKeeper defines the expected supply Keeper (noalias)
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)

	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule string, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// GovKeeper defines the expected gov Keeper
type GovKeeper interface {
	GetDepositParams(ctx sdk.Context) govtypes.DepositParams
	GetVotingParams(ctx sdk.Context) govtypes.VotingParams
	CheckMsgSubmitProposal(ctx sdk.Context, msg govtypes.MsgSubmitProposal) sdk.Error
}

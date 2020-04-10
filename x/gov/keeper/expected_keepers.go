package keeper

import (
	stakingexported "github.com/okex/okchain/x/staking/exported"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BankKeeper defines expected bank keeper
type BankKeeper interface {
	GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins

	// TODO remove once governance doesn't require use of accounts
	SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SetSendEnabled(ctx sdk.Context, enabled bool)
	SubtractCoins(ctx sdk.Context, addr sdk.AccAddress, amt sdk.Coins) (sdk.Coins, sdk.Error)
}

// StakingKeeper defines expected staking keeper (Validator and Delegator sets)
type StakingKeeper interface {
	// iterate through bonded validators by operator address, execute func for each validator
	// gov use it for getting votes of validator
	IterateBondedValidatorsByPower(sdk.Context,
		func(index int64, validator stakingexported.ValidatorI) (stop bool))

	// gov use it for getting votes of delegator which has been voted to validator
	Delegator(ctx sdk.Context, delAddr sdk.AccAddress) stakingexported.DelegatorI
}

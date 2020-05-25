// DONTCOVER
// nolint
package v0_36

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	v034staking "github.com/okex/okchain/x/staking/legacy/v0_34"
)

// Migrate accepts exported genesis state from v0.34 and migrates it to v0.36 genesis state
// All entries are identical except for validator slashing events which now include the period.
func Migrate(oldGenState v034staking.GenesisState) GenesisState {
	return NewGenesisState(
		oldGenState.Params,
		oldGenState.LastTotalPower,
		oldGenState.LastValidatorPowers,
		migrateValidators(oldGenState.Validators),
		oldGenState.Delegations,
		oldGenState.UnbondingDelegations,
		oldGenState.Redelegations,
		oldGenState.Exported,
	)
}
func migrateValidators(oldValidators v034staking.Validators) Validators {
	validators := make(Validators, len(oldValidators))

	for i, val := range oldValidators {
		bechConsPubKey, err := sdk.Bech32ifyConsPub(val.ConsPubKey)
		if err != nil {
			panic(err)
		}
		validators[i] = ValidatorExported{
			OperatorAddress:         val.OperatorAddress,
			ConsPubKey:              bechConsPubKey,
			Jailed:                  val.Jailed,
			Status:                  val.Status,
			DelegatorShares:         val.DelegatorShares,
			Description:             val.Description,
			UnbondingHeight:         val.UnbondingHeight,
			UnbondingCompletionTime: val.UnbondingCompletionTime,
			MinSelfDelegation:       sdk.MustNewDecFromStr("0.001"),
		}
	}

	return validators
}

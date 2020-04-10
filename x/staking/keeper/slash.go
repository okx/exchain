package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Jail sents a validator to jail
func (k Keeper) Jail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator := k.mustGetValidatorByConsAddr(ctx, consAddr)
	k.jailValidator(ctx, validator)
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("validator %s jailed", consAddr))
	// TODO Return event(s), blocked on https://github.com/tendermint/tendermint/pull/1803
}

// Unjail discharges a validator by unjailing
func (k Keeper) Unjail(ctx sdk.Context, consAddr sdk.ConsAddress) {
	validator := k.mustGetValidatorByConsAddr(ctx, consAddr)
	k.unjailValidator(ctx, validator)
	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("validator %s unjailed", consAddr))
	// TODO Return event(s), blocked on https://github.com/tendermint/tendermint/pull/1803
}

package staking

import (
	"fmt"

	"github.com/okx/okbchain/x/common"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/x/staking/keeper"
	"github.com/okx/okbchain/x/staking/types"
)

// BeginBlocker will persist the current header and validator set as a historical entry
// and prune the oldest entry based on the HistoricalEntries parameter
func BeginBlocker(ctx sdk.Context, k Keeper) {
	k.TrackHistoricalInfo(ctx)
}

// EndBlocker is called every block, update validator set
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	// calculate validator set changes
	var validatorUpdates []abci.ValidatorUpdate
	if k.ParamsConsensusType(ctx) == common.PoA {
		validatorUpdates = PoAValidatorsUpdate(ctx, k)
	} else {
		validatorUpdates = DPoSValidatorsUpdate(ctx, k)
	}

	// Unbond all mature validators from the unbonding queue.
	k.UnbondAllMatureValidatorQueue(ctx)

	k.IterateKeysBeforeCurrentTime(ctx, ctx.BlockHeader().Time,
		func(index int64, key []byte) (stop bool) {
			oldTime, delAddr := types.SplitCompleteTimeWithAddrKey(key)
			k.DeleteAddrByTimeKey(ctx, oldTime, delAddr)

			quantity, err := k.CompleteUndelegation(ctx, delAddr)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("complete withdraw failed: %s", err))
			} else {
				ctx.EventManager().EmitEvent(
					sdk.NewEvent(
						types.EventTypeCompleteUnbonding,
						sdk.NewAttribute(types.AttributeKeyDelegator, delAddr.String()),
						sdk.NewAttribute(sdk.AttributeKeyAmount, quantity.String()),
					),
				)
			}
			return false
		})

	return validatorUpdates
}

func PoAValidatorsUpdate(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	// calculate validator set changes
	validatorUpdates := make([]abci.ValidatorUpdate, 0)
	proposedValidators := k.GetProposeValidators(ctx)
	if len(proposedValidators) > 0 {
		validatorUpdates = k.PoAApplyAndReturnValidatorSetUpdates(ctx, proposedValidators)
		k.DeleteProposeValidators(ctx)
	} else if k.IsKickedOut(ctx) {
		// if there are some validators to kick out in an epoch
		validatorUpdates = k.KickOutAndReturnValidatorSetUpdates(ctx)
		k.DeleteAbandonedValidatorAddrs(ctx)
	}
	// update epoch
	if k.IsEndOfEpoch(ctx) {
		oldEpoch, newEpoch := k.GetEpoch(ctx), k.ParamsEpoch(ctx)
		if oldEpoch != newEpoch {
			k.SetEpoch(ctx, newEpoch)
		}
		k.SetTheEndOfLastEpoch(ctx)
	}
	return validatorUpdates
}

func DPoSValidatorsUpdate(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	// calculate validator set changes
	validatorUpdates := make([]abci.ValidatorUpdate, 0)
	if k.IsEndOfEpoch(ctx) {
		oldEpoch, newEpoch := k.GetEpoch(ctx), k.ParamsEpoch(ctx)
		if oldEpoch != newEpoch {
			k.SetEpoch(ctx, newEpoch)
		}
		k.SetTheEndOfLastEpoch(ctx)
		//ctx.Logger().Debug("validatorUpdates epoch", "old", oldEpoch, "new", newEpoch)
		//ctx.Logger().Debug(fmt.Sprintf("old epoch end blockHeight: %d", lastEpochEndHeight))

		validatorUpdates = k.ApplyAndReturnValidatorSetUpdates(ctx)
		// dont forget to delete in case that some validator need to kick out when an epoch ends
		k.DeleteAbandonedValidatorAddrs(ctx)
	} else if k.IsKickedOut(ctx) {
		// if there are some validators to kick out in an epoch
		validatorUpdates = k.KickOutAndReturnValidatorSetUpdates(ctx)
		k.DeleteAbandonedValidatorAddrs(ctx)
	}
	return validatorUpdates
}

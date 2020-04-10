package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/staking/exported"
)

// endblock at specific block height
func tally(ctx sdk.Context, versionProtocol uint64, k Keeper, threshold sdk.Dec) (passes bool) {

	totalVotingPower := sdk.ZeroDec()
	signalsVotingPower := sdk.ZeroDec()

	// computing voting power
	k.IterateBondedValidatorsByPower(ctx, func(index int64, validator exported.ValidatorI) (stop bool) {
		totalVotingPower = totalVotingPower.Add(sdk.NewDec(validator.GetConsensusPower()))
		valAcc := validator.GetConsAddr().String()
		if ok := k.GetSignal(ctx, versionProtocol, valAcc); ok {
			signalsVotingPower = signalsVotingPower.Add(sdk.NewDec(validator.GetConsensusPower()))
		}
		return false
	})

	ctx.Logger().Info("Tally Start", "SignalsVotingPower", signalsVotingPower.String(),
		"TotalVotingPower", totalVotingPower.String(),
		"SignalsVotingPower/TotalVotingPower", signalsVotingPower.Quo(totalVotingPower).String(),
		"Threshold", threshold.String())
	// If more than TH of validator update, do activate new protocol
	if signalsVotingPower.Quo(totalVotingPower).GTE(threshold) {
		return true
	}
	return false
}

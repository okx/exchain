package keeper

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/distribution/types"
	stakingexported "github.com/okex/exchain/x/staking/exported"
	"reflect"
)

// HandleChangeDistributionTypeProposal is a handler for executing a passed change distribution type proposal
func HandleChangeDistributionTypeProposal(ctx sdk.Context, k Keeper, p types.ChangeDistributionTypeProposal) error {
	logger := k.Logger(ctx)

	//1.check if it's the same
	if k.GetDistributionType(ctx) == p.Type {
		logger.Debug(fmt.Sprintf("do nothing, same distribution type, %d", p.Type))
		return nil
	}

	//2. if on chain, iteration validators and init val which has not outstanding
	if p.Type == types.DistributionTypeOnChain && !k.CheckInitExistedValidatorFlag(ctx) {
		k.SetInitExistedValidatorFlag(ctx, true)
		k.stakingKeeper.IterateValidators(ctx, func(index int64, validator stakingexported.ValidatorI) (stop bool) {
			if validator != nil {
				k.initExistedValidatorForDistrProposal(ctx, validator)
			}
			return false
		})
	}
	//3. set it
	k.SetDistributionType(ctx, p.Type)

	return nil
}

// HandleWithdrawRewardEnabledProposal is a handler for executing a passed set withdraw reward enabled proposal
func HandleWithdrawRewardEnabledProposal(ctx sdk.Context, k Keeper, p types.WithdrawRewardEnabledProposal) error {
	logger := k.Logger(ctx)
	logger.Debug(fmt.Sprintf("set withdraw reward enabled:%t", p.Enabled))
	k.SetWithdrawRewardEnabled(ctx, p.Enabled)
	return nil
}

// HandleRewardTruncatePrecisionProposal is a handler for executing a passed reward truncate precision proposal
func HandleRewardTruncatePrecisionProposal(ctx sdk.Context, k Keeper, p types.RewardTruncatePrecisionProposal) error {
	logger := k.Logger(ctx)
	logger.Debug(fmt.Sprintf("set reward truncate retain precision :%d", p.Precision))
	k.SetRewardTruncatePrecision(ctx, p.Precision)
	return nil
}

func HandleExtendProposal(ctx sdk.Context, k Keeper, p types.DistrExtendProposal) error {
	f := reflect.ValueOf(&k).MethodByName(types.InvokeExtendProposalName)
	result := f.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(p.Method), reflect.ValueOf(p.Params)})
	err := result[0].Interface()
	a, _ := err.(error)
	return a
}

func (k Keeper) InvokeExtendProposal(ctx sdk.Context, method string, params string) error {
	testExtend, err := types.NewTestExtend(params)
	if err != nil {
		//TODO zhujianguo
		return nil
	}

	//TODO the new change
	_ = testExtend

	return nil
}

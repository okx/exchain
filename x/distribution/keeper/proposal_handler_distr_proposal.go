package keeper

import (
	"fmt"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/distribution/types"
)

// HandleChangeDistributionTypeProposal is a handler for executing a passed change distribution type proposal
func HandleChangeDistributionTypeProposal(ctx sdk.Context, k Keeper, p types.ChangeDistributionTypeProposal) error {
	logger := k.Logger(ctx)

	//1.check if it's the same
	if k.GetDistributionType(ctx) == p.Type {
		logger.Debug(fmt.Sprintf("do nothing, same distribution type, %d", p.Type))
		return nil
	}

	//2. set it
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

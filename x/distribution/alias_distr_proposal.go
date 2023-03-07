// nolint
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/okx/exchain/x/distribution/types
// ALIASGEN: github.com/okx/exchain/x/distribution/client
package distribution

import (
	"github.com/okx/exchain/x/distribution/client"
	"github.com/okx/exchain/x/distribution/types"
)

var (
	NewMsgWithdrawDelegatorReward          = types.NewMsgWithdrawDelegatorReward
	CommunityPoolSpendProposalHandler      = client.CommunityPoolSpendProposalHandler
	ChangeDistributionTypeProposalHandler  = client.ChangeDistributionTypeProposalHandler
	WithdrawRewardEnabledProposalHandler   = client.WithdrawRewardEnabledProposalHandler
	RewardTruncatePrecisionProposalHandler = client.RewardTruncatePrecisionProposalHandler
	NewMsgWithdrawDelegatorAllRewards      = types.NewMsgWithdrawDelegatorAllRewards
)

// nolint
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/okex/exchain/x/distribution/types
// ALIASGEN: github.com/okex/exchain/x/distribution/client
package distribution

import (
	"github.com/okx/okbchain/x/distribution/client"
	"github.com/okx/okbchain/x/distribution/types"
)

var (
	NewMsgWithdrawDelegatorReward          = types.NewMsgWithdrawDelegatorReward
	CommunityPoolSpendProposalHandler      = client.CommunityPoolSpendProposalHandler
	ChangeDistributionTypeProposalHandler  = client.ChangeDistributionTypeProposalHandler
	WithdrawRewardEnabledProposalHandler   = client.WithdrawRewardEnabledProposalHandler
	RewardTruncatePrecisionProposalHandler = client.RewardTruncatePrecisionProposalHandler
	NewMsgWithdrawDelegatorAllRewards      = types.NewMsgWithdrawDelegatorAllRewards
)

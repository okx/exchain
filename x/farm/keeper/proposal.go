package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkGov "github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/okex/okexchain/x/farm/types"
	govKeeper "github.com/okex/okexchain/x/gov/keeper"
	govTypes "github.com/okex/okexchain/x/gov/types"
	"time"
)

var _ govKeeper.ProposalHandler = (*Keeper)(nil)

// GetMinDeposit returns min deposit
func (k Keeper) GetMinDeposit(ctx sdk.Context, content sdkGov.Content) sdk.DecCoins {
	var minDeposit sdk.DecCoins
	if _, ok := content.(types.ManageWhiteListProposal); ok {
		minDeposit = k.GetParams(ctx).ManageWhiteListMinDeposit
	}
	return minDeposit
}

func (k Keeper) GetMaxDepositPeriod(ctx sdk.Context, content sdkGov.Content) time.Duration {
	panic("implement me")
}

func (k Keeper) GetVotingPeriod(ctx sdk.Context, content sdkGov.Content) time.Duration {
	panic("implement me")
}

func (k Keeper) CheckMsgSubmitProposal(ctx sdk.Context, msg govTypes.MsgSubmitProposal) sdk.Error {
	panic("implement me")
}

func (k Keeper) AfterSubmitProposalHandler(ctx sdk.Context, proposal govTypes.Proposal) {
	panic("implement me")
}

func (k Keeper) VoteHandler(ctx sdk.Context, proposal govTypes.Proposal, vote govTypes.Vote) (string, sdk.Error) {
	panic("implement me")
}

func (k Keeper) AfterDepositPeriodPassed(ctx sdk.Context, proposal govTypes.Proposal) {
	panic("implement me")
}

func (k Keeper) RejectedHandler(ctx sdk.Context, content govTypes.Content) {
	panic("implement me")
}

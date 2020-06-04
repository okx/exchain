package token

import (
	"fmt"
	"time"

	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/gov"
	"github.com/okex/okchain/x/token/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govTypes "github.com/okex/okchain/x/gov/types"
)

// GetMinDeposit implements ProposalHandler interface
func (k Keeper) GetMinDeposit(ctx sdk.Context, content gov.Content) (minDeposit sdk.DecCoins) {
	switch content.(type) {
	case types.CertifiedTokenProposal:
		minDeposit = k.GetParams(ctx).CertifiedTokenMinDeposit
	}

	return
}

// GetMaxDepositPeriod implements ProposalHandler interface
func (k Keeper) GetMaxDepositPeriod(ctx sdk.Context, content gov.Content) (maxDepositPeriod time.Duration) {
	switch content.(type) {
	case types.CertifiedTokenProposal:
		maxDepositPeriod = k.GetParams(ctx).CertifiedTokenMaxDepositPeriod
	}

	return
}

// GetVotingPeriod implements ProposalHandler interface
func (k Keeper) GetVotingPeriod(ctx sdk.Context, content gov.Content) (votingPeriod time.Duration) {
	switch content.(type) {
	case types.CertifiedTokenProposal:
		votingPeriod = k.GetParams(ctx).CertifiedTokenVotingPeriod
	}

	return
}

// CheckMsgSubmitProposal implements ProposalHandler interface
func (k Keeper) CheckMsgSubmitProposal(ctx sdk.Context, msg gov.MsgSubmitProposal) sdk.Error {
	// check initial deposit more than or equal to ratio of MinDeposit
	initDeposit := k.GetParams(ctx).CertifiedTokenMinDeposit.MulDec(sdk.NewDecWithPrec(1, 1))
	if err := common.HasSufficientCoins(msg.Proposer, msg.InitialDeposit, initDeposit); err != nil {
		return sdk.NewError(types.DefaultCodespace, CodeInvalidAsset, err.Error())
	}
	// check proposer has sufficient coins
	if err := common.HasSufficientCoins(msg.Proposer, k.bankKeeper.GetCoins(ctx, msg.Proposer), msg.InitialDeposit); err != nil {
		return sdk.NewError(types.DefaultCodespace, CodeInvalidAsset, err.Error())
	}

	proposal := msg.Content.(types.CertifiedTokenProposal)
	if k.TokenExist(ctx, proposal.Token.Symbol) {
		return sdk.NewError(types.DefaultCodespace, types.CodeInvalidToken, fmt.Sprintf("%s already exists", proposal.Token.Symbol))
	}

	return nil
}

// nolint
func (Keeper) VoteHandler(_ sdk.Context, _ gov.Proposal, _ gov.Vote) (string, sdk.Error) {
	return "", nil
}
func (Keeper) AfterSubmitProposalHandler(_ sdk.Context, _ gov.Proposal) {}
func (Keeper) AfterDepositPeriodPassed(_ sdk.Context, _ gov.Proposal)   {}
func (Keeper) RejectedHandler(_ sdk.Context, _ gov.Content)             {}

// NewCertifiedTokenProposalHandler handles "gov" type message in "token"
func NewCertifiedTokenProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) (err sdk.Error) {
		switch c := proposal.Content.(type) {
		case types.CertifiedTokenProposal:
			return handleCertifiedTokenProposal(ctx, k, proposal)
		default:
			errMsg := fmt.Sprintf("unrecognized token proposal content type: %s", c)
			return sdk.ErrUnknownRequest(errMsg)
		}
	}
}

func handleCertifiedTokenProposal(ctx sdk.Context, keeper *Keeper, proposal *govTypes.Proposal) (err sdk.Error) {
	logger := ctx.Logger().With("module", types.ModuleName)
	logger.Debug("execute CertifiedTokenProposal begin")
	p := proposal.Content.(types.CertifiedTokenProposal)

	keeper.SetCertifiedToken(ctx, proposal.ProposalID, p.Token)

	return nil
}

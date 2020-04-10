package keeper

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/common/proto"
	"github.com/okex/okchain/x/gov"
	govTypes "github.com/okex/okchain/x/gov/types"
	"github.com/okex/okchain/x/token"
	"github.com/okex/okchain/x/upgrade/types"
)

// implement ProposalHandler interface
func (k Keeper) GetMinDeposit(ctx sdk.Context, content gov.Content) (minDeposit sdk.DecCoins) {
	switch content.(type) {
	case types.AppUpgradeProposal:
		minDeposit = k.GetParams(ctx).AppUpgradeMinDeposit
	}

	return
}

func (k Keeper) GetMaxDepositPeriod(ctx sdk.Context, content gov.Content) (maxDepositPeriod time.Duration) {
	switch content.(type) {
	case types.AppUpgradeProposal:
		maxDepositPeriod = k.GetParams(ctx).AppUpgradeMaxDepositPeriod
	}

	return
}

func (k Keeper) GetVotingPeriod(ctx sdk.Context, content gov.Content) (votingPeriod time.Duration) {
	switch content.(type) {
	case types.AppUpgradeProposal:
		votingPeriod = k.GetParams(ctx).AppUpgradeVotingPeriod
	}

	return
}

func (k Keeper) CheckMsgSubmitProposal(ctx sdk.Context, msg govTypes.MsgSubmitProposal) sdk.Error {
	// check message sender is current validator
	if !k.stakingKeeper.IsValidator(ctx, msg.Proposer) {
		return gov.ErrInvalidProposer(types.DefaultCodespace, fmt.Sprintf("proposer of App Upgrade Proposal must be validator"))
	}
	// check initial deposit more than or equal to ratio of MinDeposit
	initDeposit := k.GetParams(ctx).AppUpgradeMinDeposit.MulDec(sdk.NewDecWithPrec(1, 1))
	err := common.HasSufficientCoins(msg.Proposer, msg.InitialDeposit, initDeposit)
	if err != nil {
		return sdk.NewError(types.DefaultCodespace, token.CodeInvalidAsset, fmt.Sprintf("%s", err.Error()))
	}
	// check proposer has sufficient coins
	err = common.HasSufficientCoins(msg.Proposer, k.bankKeeper.GetCoins(ctx, msg.Proposer), msg.InitialDeposit)
	if err != nil {
		return sdk.NewError(types.DefaultCodespace, token.CodeInvalidAsset, fmt.Sprintf("%s", err.Error()))
	}

	upgradeProposal := msg.Content.(types.AppUpgradeProposal)
	if !k.protocolKeeper.IsValidVersion(ctx, upgradeProposal.ProtocolDefinition.Version) {
		return types.ErrInvalidVersion(types.DefaultCodespace, upgradeProposal.ProtocolDefinition.Version)
	}

	if uint64(ctx.BlockHeight()) > upgradeProposal.ProtocolDefinition.Height {
		return types.ErrInvalidSwitchHeight(types.DefaultCodespace, uint64(ctx.BlockHeight()), upgradeProposal.ProtocolDefinition.Height)
	}

	if _, ok := k.protocolKeeper.GetUpgradeConfig(ctx); ok {
		return types.ErrSwitchPeriodInProcess(types.DefaultCodespace)
	}
	return nil
}

func (k Keeper) AfterSubmitProposalHandler(ctx sdk.Context, proposal govTypes.Proposal) {}

func (k Keeper) VoteHandler(ctx sdk.Context, proposal govTypes.Proposal, vote govTypes.Vote) (string, sdk.Error) {
	return "", nil
}

func (k Keeper) AfterDepositPeriodPassed(ctx sdk.Context, proposal govTypes.Proposal) {}

func (k Keeper) RejectedHandler(ctx sdk.Context, content govTypes.Content) {}

func NewAppUpgradeProposalHandler(k *Keeper) govTypes.Handler {
	return func(ctx sdk.Context, proposal *govTypes.Proposal) sdk.Error {
		switch c := proposal.Content.(type) {
		case types.AppUpgradeProposal:
			return handleAppUpgradeProposal(ctx, k, proposal)

		default:
			errMsg := fmt.Sprintf("unrecognized param proposal content type: %s", c.ProposalType())
			return sdk.ErrUnknownRequest(errMsg)
		}
	}
}

func handleAppUpgradeProposal(ctx sdk.Context, k *Keeper, proposal *govTypes.Proposal) sdk.Error {
	logger := ctx.Logger().With("module", types.ModuleName)
	logger.Info("Begin to Execute AppUpgradeProposal")
	upgradeProposal := proposal.Content.(types.AppUpgradeProposal)
	if _, found := k.protocolKeeper.GetUpgradeConfig(ctx); found {
		logger.Error("Execute AppUpgradeProposal Failure", "info",
			fmt.Sprintf("App Upgrade Switch Period is in process."))
		return nil
	}

	if !k.protocolKeeper.IsValidVersion(ctx, upgradeProposal.ProtocolDefinition.Version) {
		logger.Error("Execute AppUpgradeProposal Failure", "info",
			fmt.Sprintf("version [%d] in AppUpgradeProposal is NOT valid", upgradeProposal.ProtocolDefinition.Version))
		return nil
	}

	if uint64(ctx.BlockHeight())+1 >= upgradeProposal.ProtocolDefinition.Height {
		logger.Error("Execute AppUpgradeProposal Failure", "info",
			fmt.Sprintf("switch height [%d] in AppUpgradeProposal must be more than current block height",
				upgradeProposal.ProtocolDefinition.Height))
		return nil
	}

	k.protocolKeeper.SetUpgradeConfig(ctx, proto.NewAppUpgradeConfig(proposal.ProposalID, upgradeProposal.ProtocolDefinition))

	logger.Info("Execute AppUpgradeProposal Success")

	return nil
}

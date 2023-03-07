package params

import (
	"fmt"
	"math"
	"time"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	sdkparams "github.com/okx/okbchain/libs/cosmos-sdk/x/params"
	"github.com/okx/okbchain/x/common"
	govtypes "github.com/okx/okbchain/x/gov/types"
	"github.com/okx/okbchain/x/params/types"
)

// NewParamChangeProposalHandler returns the rollback function of the param proposal handler
func NewParamChangeProposalHandler(k *Keeper) govtypes.Handler {
	return func(ctx sdk.Context, proposal *govtypes.Proposal) sdk.Error {
		switch c := proposal.Content.(type) {
		case types.ParameterChangeProposal:
			return handleParameterChangeProposal(ctx, k, proposal)
		default:
			return common.ErrUnknownProposalType(DefaultCodespace, fmt.Sprintf("%T", c))
		}
	}
}

func handleParameterChangeProposal(ctx sdk.Context, k *Keeper, proposal *govtypes.Proposal) sdk.Error {
	logger := ctx.Logger().With("module", ModuleName)
	logger.Info("Execute ParameterProposal begin")
	paramProposal := proposal.Content.(types.ParameterChangeProposal)
	curHeight := uint64(ctx.BlockHeight())
	if paramProposal.Height > curHeight {
		k.gk.InsertWaitingProposalQueue(ctx, paramProposal.Height, proposal.ProposalID)
		return nil
	}

	defer k.gk.RemoveFromWaitingProposalQueue(ctx, paramProposal.Height, proposal.ProposalID)
	return changeParams(ctx, k, paramProposal)
}

func changeParams(ctx sdk.Context, k *Keeper, paramProposal types.ParameterChangeProposal) sdk.Error {
	defer k.signalUpdate()
	for _, c := range paramProposal.Changes {
		ss, ok := k.GetSubspace(c.Subspace)
		if !ok {
			return sdkerrors.Wrap(sdkparams.ErrUnknownSubspace, c.Subspace)
		}

		err := ss.Update(ctx, []byte(c.Key), []byte(c.Value))
		if err != nil {
			return sdkerrors.Wrap(sdkparams.ErrSettingParameter, err.Error())
		}
	}
	return nil
}

func (k *Keeper) RegisterSignal(handler func()) {
	k.signals = append(k.signals, handler)
}
func (k *Keeper) signalUpdate() {
	for i, _ := range k.signals {
		k.signals[i]()
	}
}

func checkDenom(paramProposal types.ParameterChangeProposal) sdk.Error {
	for _, c := range paramProposal.Changes {
		if c.Subspace == "evm" && c.Key == "EVMDenom" {
			return sdkerrors.Wrap(sdkparams.ErrSettingParameter, "evm denom can not be reset")
		}
		if c.Subspace == "staking" && c.Key == "BondDenom" {
			return sdkerrors.Wrap(sdkparams.ErrSettingParameter, "staking bond denom can not be reset")
		}
	}
	return nil
}

// GetMinDeposit implements ProposalHandler interface
func (keeper Keeper) GetMinDeposit(ctx sdk.Context, content govtypes.Content) (minDeposit sdk.SysCoins) {
	switch content.(type) {
	case types.ParameterChangeProposal, types.UpgradeProposal:
		minDeposit = keeper.GetParams(ctx).MinDeposit
	}

	return
}

// GetMaxDepositPeriod implements ProposalHandler interface
func (keeper Keeper) GetMaxDepositPeriod(ctx sdk.Context, content govtypes.Content) (maxDepositPeriod time.Duration) {
	switch content.(type) {
	case types.ParameterChangeProposal, types.UpgradeProposal:
		maxDepositPeriod = keeper.GetParams(ctx).MaxDepositPeriod
	}

	return
}

// GetVotingPeriod implements ProposalHandler interface
func (keeper Keeper) GetVotingPeriod(ctx sdk.Context, content govtypes.Content) (votingPeriod time.Duration) {
	switch content.(type) {
	case types.ParameterChangeProposal, types.UpgradeProposal:
		votingPeriod = keeper.GetParams(ctx).VotingPeriod
	}

	return
}

// CheckMsgSubmitProposal implements ProposalHandler interface
func (keeper Keeper) CheckMsgSubmitProposal(ctx sdk.Context, msg govtypes.MsgSubmitProposal) sdk.Error {
	switch proposal := msg.Content.(type) {
	case types.ParameterChangeProposal:
		return keeper.checkSubmitParamsChangeProposal(ctx, msg.Proposer, msg.InitialDeposit, proposal)
	case types.UpgradeProposal:
		return keeper.checkSubmitUpgradeProposal(ctx, msg.Proposer, msg.InitialDeposit, proposal)
	default:
		return common.ErrUnknownProposalType(DefaultCodespace, fmt.Sprintf("%T", proposal))
	}

}

func (keeper Keeper) checkSubmitParamsChangeProposal(ctx sdk.Context, proposer sdk.AccAddress, initialDeposit sdk.SysCoins, paramsChangeProposal types.ParameterChangeProposal) sdk.Error {
	if err := keeper.proposalCommonCheck(ctx, true, proposer, initialDeposit); err != nil {
		return err
	}

	curHeight := uint64(ctx.BlockHeight())
	maxHeight := keeper.GetParams(ctx).MaxBlockHeight
	if maxHeight == 0 {
		maxHeight = math.MaxInt64 - paramsChangeProposal.Height
	}
	if paramsChangeProposal.Height < curHeight || paramsChangeProposal.Height > curHeight+maxHeight {
		return govtypes.ErrInvalidHeight(paramsChangeProposal.Height, curHeight, maxHeight)
	}

	// run simulation with cache context
	cacheCtx, _ := ctx.CacheContext()
	return changeParams(cacheCtx, &keeper, paramsChangeProposal)
}

func (keeper Keeper) checkSubmitUpgradeProposal(ctx sdk.Context, proposer sdk.AccAddress, initialDeposit sdk.SysCoins, proposal types.UpgradeProposal) sdk.Error {
	if err := keeper.proposalCommonCheck(ctx, true, proposer, initialDeposit); err != nil {
		return err
	}

	if err := checkUpgradeValidEffectiveHeight(ctx, &keeper, proposal.ExpectHeight); err != nil {
		return err
	}

	if keeper.isUpgradeExist(ctx, proposal.Name) {
		keeper.Logger(ctx).Error("upgrade has been exist", "name", proposal.Name)
		return sdk.ErrInternal(fmt.Sprintf("upgrade proposal name '%s' has been exist", proposal.Name))
	}
	return nil
}

func (keeper Keeper) proposalCommonCheck(ctx sdk.Context, checkIsValidator bool, proposer sdk.AccAddress, initialDeposit sdk.SysCoins) sdk.Error {
	// check message sender is current validator
	if checkIsValidator && !keeper.sk.IsValidator(ctx, proposer) {
		return govtypes.ErrInvalidProposer()
	}
	// check initial deposit more than or equal to ratio of MinDeposit
	initDeposit := keeper.GetParams(ctx).MinDeposit.MulDec(sdk.NewDecWithPrec(1, 1))
	if err := common.HasSufficientCoins(proposer, initialDeposit, initDeposit); err != nil {
		return sdk.ErrInvalidCoins(fmt.Sprintf("InitialDeposit must not be less than %s", initDeposit.String()))
	}
	// check proposer has sufficient coins
	if err := common.HasSufficientCoins(proposer, keeper.ck.GetCoins(ctx, proposer), initialDeposit); err != nil {
		return sdk.ErrInvalidCoins(err.Error())
	}

	return nil
}

// nolint
func (keeper Keeper) AfterSubmitProposalHandler(ctx sdk.Context, proposal govtypes.Proposal) {
	switch content := proposal.Content.(type) {
	case types.UpgradeProposal:
		// must be no error in the normal situation, for the error comes from upgrade name has been exist,
		// which has checked in CheckMsgSubmitProposal.
		_ = storePreparingUpgrade(ctx, &keeper, content)

	}
}

func (keeper Keeper) VoteHandler(ctx sdk.Context, proposal govtypes.Proposal, vote govtypes.Vote) (string, sdk.Error) {
	switch content := proposal.Content.(type) {
	case types.UpgradeProposal:
		return checkUpgradeVote(ctx, proposal.ProposalID, content, vote)
	}
	return "", nil
}
func (keeper Keeper) AfterDepositPeriodPassed(ctx sdk.Context, proposal govtypes.Proposal) {}
func (keeper Keeper) RejectedHandler(ctx sdk.Context, content govtypes.Content)            {}

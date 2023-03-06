package params

import (
	"fmt"
	"math"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	"github.com/okx/okbchain/x/common"
	govtypes "github.com/okx/okbchain/x/gov/types"
	"github.com/okx/okbchain/x/params/types"
)

func NewUpgradeProposalHandler(k *Keeper) govtypes.Handler {
	return func(ctx sdk.Context, proposal *govtypes.Proposal) sdk.Error {
		switch c := proposal.Content.(type) {
		case types.UpgradeProposal:
			return handleUpgradeProposal(ctx, k, proposal.ProposalID, c)
		default:
			return common.ErrUnknownProposalType(DefaultCodespace, fmt.Sprintf("%T", c))
		}
	}
}

func handleUpgradeProposal(ctx sdk.Context, k *Keeper, proposalID uint64, proposal types.UpgradeProposal) sdk.Error {
	curHeight := uint64(ctx.BlockHeight())
	confirmHeight, err := getUpgradeProposalConfirmHeight(curHeight, proposal)
	if err != nil {
		return err
	}
	effectiveHeight := confirmHeight + 1

	if curHeight < confirmHeight {
		k.gk.InsertWaitingProposalQueue(ctx, confirmHeight, proposalID)
		_ = storeWaitingUpgrade(ctx, k, proposal, effectiveHeight) // ignore error
		return nil
	}

	// proposal will be confirmed right now, check if ready.
	cbs, ready := k.queryReadyForUpgrade(proposal.Name)
	if !ready {
		// if no module claims that has ready for this upgrade,
		// that probably means program's version is too low.
		// To avoid status machine broken, we panic.
		errMsg := fmt.Sprintf("there's a upgrade proposal named '%s' has been take effective, "+
			"and the upgrade is incompatible, but your binary seems not ready for this upgrade. "+
			"To avoid state machine broken, the program is panic. "+
			"Using the latest version binary and re-run it to avoid this panic.", proposal.Name)
		k.Logger(ctx).Error(errMsg)
		panic(errMsg)
	}

	storedInfo, err := storeEffectiveUpgrade(ctx, k, proposal, effectiveHeight)
	if err != nil {
		return err
	}

	for _, cb := range cbs {
		if cb != nil {
			cb(storedInfo)
		}
	}
	return nil
}

func getUpgradeProposalConfirmHeight(currentHeight uint64, proposal types.UpgradeProposal) (uint64, sdk.Error) {
	// confirm height is the height proposal is confirmed.
	// confirmed is not become effective. Becoming effective will happen at
	// the next block of confirm block. see `storeEffectiveUpgrade` and `IsUpgradeEffective`
	confirmHeight := proposal.ExpectHeight - 1
	if proposal.ExpectHeight == 0 {
		// if height is not specified, this upgrade will become effective
		// at the next block of the block which the proposal is passed
		// (i.e. become effective at next block).
		confirmHeight = currentHeight
	}

	if confirmHeight < currentHeight {
		// if it's too late to make the proposal become effective at the height which we expected,
		// refuse to effective this proposal
		return 0, sdkerrors.New(DefaultCodespace, types.BaseParamsError,
			fmt.Sprintf("current height '%d' has exceed "+
				"the expect height '%d' of upgrade proposal '%s'",
				currentHeight, proposal.ExpectHeight, proposal.Name))
	}
	return confirmHeight, nil
}

func storePreparingUpgrade(ctx sdk.Context, k *Keeper, upgrade types.UpgradeProposal) sdk.Error {
	info := types.UpgradeInfo{
		Name:         upgrade.Name,
		ExpectHeight: upgrade.ExpectHeight,
		Config:       upgrade.Config,

		EffectiveHeight: 0,
		Status:          types.UpgradeStatusPreparing,
	}

	return k.writeUpgradeInfo(ctx, info, false)
}

func storeWaitingUpgrade(ctx sdk.Context, k *Keeper, upgrade types.UpgradeProposal, effectiveHeight uint64) error {
	info := types.UpgradeInfo{
		Name:         upgrade.Name,
		ExpectHeight: upgrade.ExpectHeight,
		Config:       upgrade.Config,

		EffectiveHeight: effectiveHeight,
		Status:          types.UpgradeStatusWaitingEffective,
	}

	return k.writeUpgradeInfo(ctx, info, true)
}

func storeEffectiveUpgrade(ctx sdk.Context, k *Keeper, upgrade types.UpgradeProposal, effectiveHeight uint64) (types.UpgradeInfo, sdk.Error) {
	info := types.UpgradeInfo{
		Name:         upgrade.Name,
		ExpectHeight: upgrade.ExpectHeight,
		Config:       upgrade.Config,

		EffectiveHeight: effectiveHeight,
		Status:          types.UpgradeStatusEffective,
	}

	return info, k.writeUpgradeInfo(ctx, info, true)
}

// a upgrade valid effective height must be:
//  1. zero, or
//  2. bigger than current height and not too far away from current height
func checkUpgradeValidEffectiveHeight(ctx sdk.Context, k *Keeper, effectiveHeight uint64) sdk.Error {
	if effectiveHeight == 0 {
		return nil
	}

	curHeight := uint64(ctx.BlockHeight())

	maxHeight := k.GetParams(ctx).MaxBlockHeight
	if maxHeight == 0 {
		maxHeight = math.MaxInt64 - effectiveHeight
	}

	if effectiveHeight <= curHeight || effectiveHeight-curHeight > maxHeight {
		return govtypes.ErrInvalidHeight(effectiveHeight, curHeight, maxHeight)
	}
	return nil
}

func checkUpgradeVote(ctx sdk.Context, proposalID uint64, proposal types.UpgradeProposal, _ govtypes.Vote) (string, sdk.Error) {
	curHeight := uint64(ctx.BlockHeight())

	if proposal.ExpectHeight != 0 && proposal.ExpectHeight <= curHeight {
		return "", sdkerrors.New(DefaultCodespace, types.BaseParamsError,
			fmt.Sprintf("can not vote： current height '%d' has exceed "+
				"the expect height '%d' of upgrade proposal '%s'(proposal id '%d')",
				curHeight, proposal.ExpectHeight, proposal.Name, proposalID))
	}

	return "", nil
}

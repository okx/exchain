package params

import (
	"encoding/json"
	"fmt"
	"github.com/okex/exchain/x/params/types"
	"math"
	"time"

	"github.com/okex/exchain/libs/cosmos-sdk/store/prefix"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	sdkparams "github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"github.com/okex/exchain/x/common"
	govtypes "github.com/okex/exchain/x/gov/types"
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

func handleUpgradeProposal(ctx sdk.Context, k *Keeper, proposalID uint64, proposal types.UpgradeProposal) sdk.Error {
	curHeight := uint64(ctx.BlockHeight())
	confirmHeight := getUpgradeProposalConfirmHeight(curHeight, proposal)

	if curHeight < confirmHeight {
		k.gk.InsertWaitingProposalQueue(ctx, confirmHeight, proposalID)
	}
	return confirmUpgrade(ctx, k, proposal, confirmHeight+1)
}

func confirmUpgrade(ctx sdk.Context, k *Keeper, proposal types.UpgradeProposal, effectiveHeight uint64) sdk.Error {
	key := []byte(proposal.Name)

	store := getUpgradeStore(ctx, k)
	if store.Has(key) {
		k.Logger(ctx).Error("upgrade proposal name has been exist", "proposal name", proposal.Name)
		return sdk.ErrInternal(fmt.Sprintf("upgrade proposal name '%s' has been exist", proposal.Name))
	}

	proposal.UpgradeInfo.EffectiveHeight = effectiveHeight
	data, err := json.Marshal(proposal.UpgradeInfo)
	if err != nil {
		k.Logger(ctx).Error("marshal upgrade proposal error", "upgrade info", proposal.UpgradeInfo, "error", err)
		return sdk.ErrInternal(err.Error())
	}
	store.Set(key, data)
	return nil
}

func getUpgradeProposalConfirmHeight(currentHeight uint64, proposal types.UpgradeProposal) uint64 {
	// confirm height is the height proposal is confirmed.
	// confirmed is not become effective. Becoming effective will happen at
	// the next block of confirm block. see last argument of `confirmUpgrade` and `IsUpgradeEffective`
	confirmHeight := proposal.ExpectHeight - 1

	if proposal.ExpectHeight == 0 {
		// if height is not specified, this upgrade will become effective
		// at the next block of the block which the proposal is passed
		// (i.e. become effective at next block).
		confirmHeight = currentHeight
	} else if proposal.ExpectHeight <= currentHeight {
		// if it's too late to make the proposal become effective at the height
		// which we expected, make the proposal take effect at next block.
		confirmHeight = currentHeight
	}

	return confirmHeight
}

func getUpgradeStore(ctx sdk.Context, k *Keeper) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(k.storeKey), []byte("upgrade"))
}

func getUpgradeInfo(ctx sdk.Context, k *Keeper, name string) (bool, types.UpgradeInfo, error) {
	store := getUpgradeStore(ctx, k)
	data := store.Get([]byte(name))
	if len(data) == 0 {
		k.Logger(ctx).Info("upgrade not exist", "name", name)
		return false, types.UpgradeInfo{}, nil
	}

	var info types.UpgradeInfo
	if err := json.Unmarshal(data, &info); err != nil {
		k.Logger(ctx).Error("unmarshal upgrade proposal error", "error", err, "name", name, "data", data)
		return false, info, err
	}

	return true, info, nil
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
	if err := keeper.checkValidEffectiveHeight(ctx, paramsChangeProposal.Height); err != nil {
		return err
	}

	// run simulation with cache context
	cacheCtx, _ := ctx.CacheContext()
	return changeParams(cacheCtx, &keeper, paramsChangeProposal)
}

func (keeper Keeper) checkSubmitUpgradeProposal(ctx sdk.Context, proposer sdk.AccAddress, initialDeposit sdk.SysCoins, proposal types.UpgradeProposal) sdk.Error {
	if err := keeper.proposalCommonCheck(ctx, true, proposer, initialDeposit); err != nil {
		return err
	}

	if proposal.ExpectHeight != 0 {
		if err := keeper.checkValidEffectiveHeight(ctx, proposal.ExpectHeight); err != nil {
			return err
		}
	}

	// run simulation with cache context
	cacheCtx, _ := ctx.CacheContext()
	return confirmUpgrade(cacheCtx, &keeper, proposal, proposal.ExpectHeight)
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

// a valid effective height must be:
//  1. bigger than current height
//  2. not too far away from current height
func (keeper Keeper) checkValidEffectiveHeight(ctx sdk.Context, effectiveHeight uint64) sdk.Error {
	curHeight := uint64(ctx.BlockHeight())

	maxHeight := keeper.GetParams(ctx).MaxBlockHeight
	if maxHeight == 0 {
		maxHeight = math.MaxInt64 - effectiveHeight
	}

	if effectiveHeight < curHeight || effectiveHeight-curHeight > maxHeight {
		return govtypes.ErrInvalidHeight(effectiveHeight, curHeight, maxHeight)
	}
	return nil
}

// nolint
func (keeper Keeper) AfterSubmitProposalHandler(ctx sdk.Context, proposal govtypes.Proposal) {}
func (keeper Keeper) VoteHandler(ctx sdk.Context, proposal govtypes.Proposal, vote govtypes.Vote) (string, sdk.Error) {
	return "", nil
}
func (keeper Keeper) AfterDepositPeriodPassed(ctx sdk.Context, proposal govtypes.Proposal) {}
func (keeper Keeper) RejectedHandler(ctx sdk.Context, content govtypes.Content)            {}

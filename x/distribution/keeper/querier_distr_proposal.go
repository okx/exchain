package keeper

import (
	"encoding/json"

	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	comm "github.com/okx/okbchain/x/common"
	"github.com/okx/okbchain/x/distribution/types"
)

func queryDelegationRewards(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	if !k.CheckDistributionProposalValid(ctx) {
		return nil, types.ErrCodeNotSupportDistributionProposal()
	}

	var params types.QueryDelegationRewardsParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, comm.ErrUnMarshalJSONFailed(err.Error())
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()

	val := k.stakingKeeper.Validator(ctx, params.ValidatorAddress)
	if val == nil {
		return nil, sdkerrors.Wrap(types.ErrCodeEmptyValidatorDistInfo(), params.ValidatorAddress.String())
	}

	del := k.stakingKeeper.Delegator(ctx, params.DelegatorAddress)
	if del == nil {
		return nil, types.ErrCodeEmptyDelegationDistInfo()
	}

	found := false
	for _, valAddr := range del.GetShareAddedValidatorAddresses() {
		if valAddr.Equals(params.ValidatorAddress) {
			found = true
		}
	}
	if !found {
		return nil, sdkerrors.Wrap(types.ErrCodeEmptyDelegationVoteValidator(), params.ValidatorAddress.String())
	}

	logger := k.Logger(ctx)
	if !k.HasDelegatorStartingInfo(ctx, val.GetOperator(), params.DelegatorAddress) {
		if del.GetLastAddedShares().IsZero() {
			return nil, sdkerrors.Wrap(types.ErrCodeZeroDelegationShares(), params.DelegatorAddress.String())
		}
		k.initExistedDelegationStartInfo(ctx, val, del)
	}

	endingPeriod := k.incrementValidatorPeriod(ctx, val)
	rewards := k.calculateDelegationRewards(ctx, val, params.DelegatorAddress, endingPeriod)
	if rewards == nil {
		rewards = sdk.DecCoins{}
	}

	logger.Debug("queryDelegationRewards", "Validator", val.GetOperator(),
		"Delegator", params.DelegatorAddress, "Reward", rewards)

	bz, err := codec.MarshalJSONIndent(k.cdc, rewards)
	if err != nil {
		return nil, comm.ErrMarshalJSONFailed(err.Error())
	}

	return bz, nil
}

func queryDelegatorTotalRewards(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	if !k.CheckDistributionProposalValid(ctx) {
		return nil, types.ErrCodeNotSupportDistributionProposal()
	}

	var params types.QueryDelegatorParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, comm.ErrUnMarshalJSONFailed(err.Error())
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()

	del := k.stakingKeeper.Delegator(ctx, params.DelegatorAddress)
	if del == nil {
		return nil, types.ErrCodeEmptyDelegationDistInfo()
	}

	total := sdk.DecCoins{}
	var delRewards []types.DelegationDelegatorReward

	for _, valAddr := range del.GetShareAddedValidatorAddresses() {
		val := k.stakingKeeper.Validator(ctx, valAddr)
		if val == nil {
			continue
		}

		logger := k.Logger(ctx)
		if !k.HasDelegatorStartingInfo(ctx, val.GetOperator(), params.DelegatorAddress) {
			if del.GetLastAddedShares().IsZero() {
				return nil, sdkerrors.Wrap(types.ErrCodeZeroDelegationShares(), params.DelegatorAddress.String())
			}
			k.initExistedDelegationStartInfo(ctx, val, del)
		}

		endingPeriod := k.incrementValidatorPeriod(ctx, val)
		delReward := k.calculateDelegationRewards(ctx, val, params.DelegatorAddress, endingPeriod)
		if delReward == nil {
			delReward = sdk.DecCoins{}
		}
		delRewards = append(delRewards, types.NewDelegationDelegatorReward(valAddr, delReward))
		total = total.Add(delReward...)
		logger.Debug("queryDelegatorTotalRewards", "Validator", val.GetOperator(),
			"Delegator", params.DelegatorAddress, "Reward", delReward)
	}

	totalRewards := types.NewQueryDelegatorTotalRewardsResponse(delRewards, total)

	bz, err := json.Marshal(totalRewards)
	if err != nil {
		return nil, comm.ErrMarshalJSONFailed(err.Error())
	}

	return bz, nil
}

func queryValidatorOutstandingRewards(ctx sdk.Context, path []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	if !k.CheckDistributionProposalValid(ctx) {
		return nil, types.ErrCodeNotSupportDistributionProposal()
	}

	var params types.QueryValidatorOutstandingRewardsParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, comm.ErrUnMarshalJSONFailed(err.Error())
	}

	rewards := k.GetValidatorOutstandingRewards(ctx, params.ValidatorAddress)
	if rewards == nil {
		rewards = sdk.DecCoins{}
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, rewards)
	if err != nil {
		return nil, comm.ErrMarshalJSONFailed(err.Error())
	}

	return bz, nil
}

func queryDelegatorValidators(ctx sdk.Context, _ []string, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryDelegatorParams
	err := k.cdc.UnmarshalJSON(req.Data, &params)
	if err != nil {
		return nil, comm.ErrUnMarshalJSONFailed(err.Error())
	}

	// cache-wrap context as to not persist state changes during querying
	ctx, _ = ctx.CacheContext()

	delegator := k.stakingKeeper.Delegator(ctx, params.DelegatorAddress)

	if delegator == nil {
		return nil, types.ErrCodeEmptyDelegationDistInfo()
	}

	bz, err := codec.MarshalJSONIndent(k.cdc, delegator.GetShareAddedValidatorAddresses())
	if err != nil {
		return nil, comm.ErrMarshalJSONFailed(err.Error())
	}

	return bz, nil
}

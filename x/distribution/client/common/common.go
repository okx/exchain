package common

import (
	"encoding/json"
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"github.com/okex/exchain/x/distribution/types"
)

// QueryParams actually queries distribution params
func QueryParams(cliCtx context.CLIContext, queryRoute string) (params types.Params, err error) {
	route := fmt.Sprintf("custom/%s/params/%s", queryRoute, types.ParamCommunityTax)
	var communityTax sdk.Dec
	var withdrawAddrEnabled bool
	bytes, _, err := cliCtx.QueryWithData(route, []byte{})
	if err != nil {
		return
	}
	cliCtx.Codec.MustUnmarshalJSON(bytes, &communityTax)

	route = fmt.Sprintf("custom/%s/params/%s", queryRoute, types.ParamWithdrawAddrEnabled)
	bytes, _, err = cliCtx.QueryWithData(route, []byte{})
	if err != nil {
		return
	}
	cliCtx.Codec.MustUnmarshalJSON(bytes, &withdrawAddrEnabled)

	var distributionType uint32
	route = fmt.Sprintf("custom/%s/params/%s", queryRoute, types.ParamDistributionType)
	bytes, _, err = cliCtx.QueryWithData(route, []byte{})
	if err == nil {
		cliCtx.Codec.MustUnmarshalJSON(bytes, &distributionType)
	} else if !ignoreError(err.Error()) {
		return
	}

	var withdrawRewardEnabled bool
	route = fmt.Sprintf("custom/%s/params/%s", queryRoute, types.ParamWithdrawRewardEnabled)
	bytes, _, err = cliCtx.QueryWithData(route, []byte{})
	if err == nil {
		cliCtx.Codec.MustUnmarshalJSON(bytes, &withdrawRewardEnabled)
	} else if !ignoreError(err.Error()) {
		return
	}

	var rewardTruncatePrecision int64
	route = fmt.Sprintf("custom/%s/params/%s", queryRoute, types.ParamRewardTruncatePrecision)
	bytes, _, err = cliCtx.QueryWithData(route, []byte{})
	if err == nil {
		cliCtx.Codec.MustUnmarshalJSON(bytes, &rewardTruncatePrecision)
	} else if !ignoreError(err.Error()) {
		return
	}

	return types.NewParams(communityTax, withdrawAddrEnabled, distributionType, withdrawRewardEnabled, rewardTruncatePrecision), nil
}

func ignoreError(err string) bool {
	type ParamsError struct {
		Code uint32 `json:"code"`
	}
	var paramsError ParamsError
	jsonErr := json.Unmarshal([]byte(err), &paramsError)
	if jsonErr != nil {
		return false
	}
	return paramsError.Code == types.CodeUnknownDistributionParamType
}

// QueryValidatorCommission returns a validator's commission.
func QueryValidatorCommission(cliCtx context.CLIContext, queryRoute string, validatorAddr sdk.ValAddress) (
	[]byte, error) {
	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryValidatorCommission),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryValidatorCommissionParams(validatorAddr)),
	)
	return res, err
}

// WithdrawValidatorRewardsAndCommission builds a two-message message slice to be
// used to withdraw both validation's commission and self-delegation reward.
func WithdrawValidatorRewardsAndCommission(validatorAddr sdk.ValAddress) ([]sdk.Msg, error) {
	commissionMsg := types.NewMsgWithdrawValidatorCommission(validatorAddr)
	if err := commissionMsg.ValidateBasic(); err != nil {
		return nil, err
	}

	return []sdk.Msg{commissionMsg}, nil
}

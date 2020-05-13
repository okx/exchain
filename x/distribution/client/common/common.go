package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/distribution/types"
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

	return types.NewParams(communityTax, withdrawAddrEnabled), err
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

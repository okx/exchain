package common

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/x/distribution/types"
	staking "github.com/okx/okbchain/x/staking/types"
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
	if err != nil {
		return
	}
	cliCtx.Codec.MustUnmarshalJSON(bytes, &distributionType)

	var withdrawRewardEnabled bool
	route = fmt.Sprintf("custom/%s/params/%s", queryRoute, types.ParamWithdrawRewardEnabled)
	bytes, _, err = cliCtx.QueryWithData(route, []byte{})
	if err != nil {
		return
	}
	cliCtx.Codec.MustUnmarshalJSON(bytes, &withdrawRewardEnabled)

	var rewardTruncatePrecision int64
	route = fmt.Sprintf("custom/%s/params/%s", queryRoute, types.ParamRewardTruncatePrecision)
	bytes, _, err = cliCtx.QueryWithData(route, []byte{})
	if err != nil {
		return
	}
	cliCtx.Codec.MustUnmarshalJSON(bytes, &rewardTruncatePrecision)

	return types.NewParams(communityTax, withdrawAddrEnabled, distributionType, withdrawRewardEnabled, rewardTruncatePrecision), nil
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

type WrapError struct {
	Code      uint32 `json:"code"`
	Log       string `json:"log"`
	Codespace string `json:"codespace"`
	Changed   bool   `json:"-"`
	RawError  error  `json:"-"`
}

func (e *WrapError) Error() string {
	data, jsonErr := json.Marshal(e)
	if jsonErr != nil {
		fmt.Fprintf(os.Stderr, "Trans wrap error, marshal err=%v\n", jsonErr)
		return ""
	}
	return string(data)
}

func (e *WrapError) setLog(log string) {
	e.Log = log
}

func (e *WrapError) Trans(code uint32, newLog string) {
	if e.Code == abci.CodeTypeNonceInc+code {
		e.setLog(fmt.Sprintf("%s;%s", newLog, e.Log))
		e.Changed = true
	}
}

func NewWrapError(err error) *WrapError {
	if err == nil {
		return nil
	}
	var wrapErr WrapError
	jsonErr := json.Unmarshal([]byte(err.Error()), &wrapErr)
	if jsonErr != nil {
		fmt.Fprintf(os.Stderr, "Trans wrap error, unmarshal err=%v\n", jsonErr)
		return nil
	}
	wrapErr.RawError = err
	return &wrapErr
}

func IsValidator(cliCtx context.CLIContext, cdc *codec.Codec, valAddress sdk.ValAddress) bool {
	resKVs, _, err := cliCtx.QuerySubspace(staking.ValidatorsKey, staking.StoreKey)
	if err != nil {
		return false
	}

	for _, kv := range resKVs {
		if staking.MustUnmarshalValidator(cdc, kv.Value).GetOperator().Equals(valAddress) {
			return true
		}
	}

	return false
}

func IsDelegator(cliCtx context.CLIContext, cdc *codec.Codec, delAddr sdk.AccAddress) bool {
	delegator := staking.NewDelegator(delAddr)
	resp, _, err := cliCtx.QueryStore(staking.GetDelegatorKey(delAddr), staking.StoreKey)
	if err != nil {
		return false
	}
	if len(resp) == 0 {
		return false
	}
	cdc.MustUnmarshalBinaryLengthPrefixed(resp, &delegator)

	if delegator.Tokens.IsZero() {
		return false
	}

	return true
}

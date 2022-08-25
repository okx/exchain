package common

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/x/gov/types"
)

// QueryParams actually queries gov params
func QueryParams(cliCtx context.CLIContext, paramType string) (types.CM45Params, int64, error) {
	route := fmt.Sprintf("custom/gov/%s/%s", types.QueryParams, paramType)
	var height int64
	vp := types.DefaultVotingParams()
	tp := types.DefaultTallyParams()
	dp := types.DefaultDepositParams()
	switch paramType {
	case types.ParamVoting:
		var voting types.VotingParams
		bytes, h, err := cliCtx.Query(route)
		if err != nil {
			return types.NewCM45Params(vp.ToCM45VotingParams(), tp, dp.ToCM45DepositParams()), 0, err
		}
		cliCtx.Codec.MustUnmarshalJSON(bytes, &voting)
		vp = voting
		height = h
	case types.ParamTallying:
		var tallying types.TallyParams
		bytes, h, err := cliCtx.Query(route)
		if err != nil {
			return types.NewCM45Params(vp.ToCM45VotingParams(), tp, dp.ToCM45DepositParams()), 0, err
		}
		cliCtx.Codec.MustUnmarshalJSON(bytes, &tallying)
		tp = tallying
		height = h
	case types.ParamDeposit:
		var deposit types.DepositParams
		bytes, h, err := cliCtx.Query(route)
		if err != nil {
			return types.NewCM45Params(vp.ToCM45VotingParams(), tp, dp.ToCM45DepositParams()), 0, err
		}
		cliCtx.Codec.MustUnmarshalJSON(bytes, &deposit)
		dp = deposit
		height = h
	default:
		return types.NewCM45Params(vp.ToCM45VotingParams(), tp, dp.ToCM45DepositParams()), 0,
			sdkerrors.Wrapf(sdkerrors.ErrInvalidType, "%s is not a valid param type", paramType)
	}
	return types.NewCM45Params(vp.ToCM45VotingParams(), tp, dp.ToCM45DepositParams()), height, nil
}

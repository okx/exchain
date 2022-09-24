package evm

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/staking"
	"github.com/okex/exchain/x/staking/client/cli"
	"github.com/okex/exchain/x/staking/types"
)

var (
	//cdc          = appCodec.MakeCodec(app.ModuleBasics)
	//interfaceReg = appCodec.MakeIBC(app.ModuleBasics)
	//protoCdc     = codec.NewProtoCodec(interfaceReg)
	storeName = staking.StoreKey
)

func QueryDelegator(cliCtx context.CLIContext, addr string) ([]byte, error) {
	delAddr, err := sdk.AccAddressFromBech32(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid address：%s", addr)
	}

	cdc := cliCtx.Codec
	delegator, undelegation := types.NewDelegator(delAddr), types.DefaultUndelegation()
	resp, _, err := cliCtx.QueryStore(types.GetDelegatorKey(delAddr), storeName)
	if err != nil {
		return nil, err
	}
	if len(resp) != 0 {
		cdc.MustUnmarshalBinaryLengthPrefixed(resp, &delegator)
	}

	// query for the undelegation info
	bytes, err := cdc.MarshalJSON(types.NewQueryDelegatorParams(delAddr))
	if err != nil {
		return nil, err
	}

	route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryUnbondingDelegation)
	res, _, err := cliCtx.QueryWithData(route, bytes)
	// if err!= nil , we treat it as there's no undelegation of the delegator
	if err == nil {
		if err := cdc.UnmarshalJSON(res, &undelegation); err != nil {
			return nil, err
		}
	}
	result := convertToDelegatorResp(delegator, undelegation)
	return cdc.MarshalJSON(result)
}

func convertToDelegatorResp(delegator types.Delegator, undelegation types.UndelegationInfo,
) cli.DelegatorResponse {
	return cli.DelegatorResponse{
		delegator.DelegatorAddress,
		delegator.ValidatorAddresses,
		delegator.Shares,
		delegator.Tokens,
		undelegation.Quantity,
		undelegation.CompletionTime,
		delegator.IsProxy,
		delegator.TotalDelegatedTokens,
		delegator.ProxyAddress,
	}
}

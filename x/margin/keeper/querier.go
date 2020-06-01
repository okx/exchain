package keeper

import (
	"fmt"

	"github.com/okex/okchain/x/common"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/margin/types"
)

// NewQuerier creates a new querier for margin clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryParameters:
			return queryParams(ctx, k)
		case types.QueryMarginAccount:
			return queryMarginAccount(ctx, path, req, k)

		default:
			return nil, sdk.ErrUnknownRequest("unknown dex query endpoint")
		}
	}
}

func queryMarginAccount(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	addr, err := sdk.AccAddressFromBech32(path[1])
	if err != nil {
		return nil, sdk.ErrInvalidAddress(fmt.Sprintf("invalid address：%s", path[1]))
	}

	//marginAcc := types.GetMarginAccount(addr.String())
	marginDeposit := keeper.GetAccounts(ctx, addr)
	res, err := common.JSONMarshalV2(marginDeposit)
	if err != nil {
		return nil, sdk.ErrInternal(err.Error())
	}
	return res, nil
}

func queryParams(ctx sdk.Context, k Keeper) (res []byte, err sdk.Error) {
	params := k.GetParams(ctx)

	res, errUnmarshal := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if errUnmarshal != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal result to JSON", errUnmarshal.Error()))
	}

	return res, nil
}

// TODO: Add the modules query functions
// They will be similar to the above one: queryParams()

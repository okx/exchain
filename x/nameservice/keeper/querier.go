/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/21 5:28 上午
# @File : querier.go
# @Description :
# @Attention :
*/
package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/nameservice/types"
)

// NewQuerier creates a new querier for nameservice clients.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		// this line is used by starport scaffolding # 2
		case types.QueryResolveName:
			return resolveName(ctx, path[1:], k)
		case types.QueryListWhois:
			return listWhois(ctx, k)
		case types.QueryGetWhois:
			return getWhois(ctx, path[1:], k)
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "unknown nameservice query endpoint")
		}
	}
}

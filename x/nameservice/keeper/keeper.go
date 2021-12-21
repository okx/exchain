/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/21 4:54 上午
# @File : keeper.go
# @Description :
# @Attention :
*/
package keeper

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	store "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/nameservice/types"
)

// Keeper of the nameservice store
type Keeper struct {
	CoinKeeper bank.Keeper
	storeKey   store.StoreKey
	cdc        *codec.Codec
	// paramspace types.ParamSubspace
}

// NewKeeper creates a nameservice keeper
func NewKeeper(coinKeeper bank.Keeper, cdc *codec.Codec, key store.StoreKey) Keeper {
	keeper := Keeper{
		CoinKeeper: coinKeeper,
		storeKey:   key,
		cdc:        cdc,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

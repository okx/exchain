/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/21 5:09 上午
# @File : alias.go
# @Description :
# @Attention :
*/
package nameservice

import (
	"github.com/okex/exchain/x/nameservice/keeper"
	"github.com/okex/exchain/x/nameservice/types"
)

type (
	// nolint
	Keeper = keeper.Keeper
)
const (
	StoreKey            = types.StoreKey
	DefaultParamspace   = types.DefaultParamspace
)
var (
	ModuleName = types.ModuleName
	NewKeeper  = keeper.NewKeeper
	RouterKey         = types.RouterKey
)

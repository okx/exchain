/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/19 8:44 上午
# @File : errors.go
# @Description :
# @Attention :
*/
package types

import (
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

var (
	ErrNameDoesNotExist = sdkerrors.Register(ModuleName, 1, "name does not exist")
)

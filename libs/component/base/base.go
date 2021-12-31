/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/31 2:57 下午
# @File : base.go
# @Description :
# @Attention :
*/
package base

import (
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/libs/service"
)

// used for DI
type IComponent interface {
	service.Service
}

type BaseComponent struct {
	service.BaseService
}

func NewBaseComponent(logger log.Logger, name string, com IComponent) *BaseComponent {
	ret := &BaseComponent{
		BaseService: *service.NewBaseService(logger, name, com),
	}
	return ret
}

/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/31 2:59 下午
# @File : component.go
# @Description :
# @Attention :
*/
package listener

import "github.com/okex/exchain/libs/component/base"


// What: producer-consumer listener component
// Why: tendermint event plugin cant satisfy what i want (remove immediatelly after called)
// Performance: just alloc chan memory and map put/delete/get  ,thats all
// support repeated registration
type IListenerComponent interface {
	base.IComponent

	RegisterListener(topic ...string) <-chan interface{}
	NotifyListener(data interface{}, listenerIds ...string)
	CancelAsync(listenerIds ...string)
}


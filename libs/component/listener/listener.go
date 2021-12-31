/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/31 3:00 下午
# @File : impl.go
# @Description :
# @Attention :
*/
package listener

import (
	"github.com/okex/exchain/libs/component/base"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)

var (
	_ IListenerComponent = (*DefaultListenerComponent)(nil)
)

type DefaultListenerComponent struct {
	base.BaseComponent

	pubsub *PubSub
}

func DefaultNewListenerComponent(logger log.Logger, name string) *DefaultListenerComponent {
	const c int = 256
	return NewDefaultListenerComponent(c, logger, name)
}

func NewDefaultListenerComponent(cap int, logger log.Logger, name string) *DefaultListenerComponent {
	ret := &DefaultListenerComponent{}
	ret.BaseComponent = *base.NewBaseComponent(logger, name, ret)
	ret.pubsub = New(cap)
	return ret
}

func (l *DefaultListenerComponent) RegisterListener(topic ...string) <-chan interface{} {
	return l.pubsub.SubOnce(topic...)
}

func (l *DefaultListenerComponent) NotifyListener(data interface{}, listenerIds ...string) {
	l.pubsub.Pub(data, listenerIds...)
}

func (l *DefaultListenerComponent) CancelAsync(listenerIds ...string) {
	l.pubsub.Pub(nil, listenerIds...)
}

func (l *DefaultListenerComponent) OnStart() error {
	go l.pubsub.start()
	return nil
}

func (l *DefaultListenerComponent) OnStop() {
	l.pubsub.Stop()
}

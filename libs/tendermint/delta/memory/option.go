/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2022/1/2 11:05 上午
# @File : ops.go
# @Description :
# @Attention :
*/
package memory

import "time"

type MemoryBrokerOption func(m *MemoryBroker)

func MemoryBrokerWithTTL(t time.Duration) MemoryBrokerOption {
	return func(m *MemoryBroker) {
		m.ttl = t
	}
}

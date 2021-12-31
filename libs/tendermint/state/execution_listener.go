/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/31 4:33 下午
# @File : execution_listener.go
# @Description :
# @Attention :
*/
package state

import "strconv"

var closedCh = make(chan interface{}, 1)

func init() {
	close(closedCh)
}

// TODO, HASH IS BETTER
func registerListenerBeforePrerun(ec *BlockExecutor, h int64) (<-chan interface{}, func()) {
	if ec.deltaContext.downloadDelta && nil != ec.listener {
		p, _ := ec.deltaContext.dataMap.acquire(h)
		if p != nil {
			ret := make(chan interface{}, 1)
			ret <- p
			return ret, func() {}
		}
		return ec.listener.RegisterListener(strconv.Itoa(int(h))), func() {
			ec.listener.CancelAsync(strconv.Itoa(int(h)))
		}
	}
	return closedCh, func() {}
}

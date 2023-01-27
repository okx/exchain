package global

import "sync/atomic"

var record int32

func SetGlobalRecord(flag int32) {
	atomic.StoreInt32(&record, flag)
}

func Record() bool {
	return atomic.LoadInt32(&record) == 1
}

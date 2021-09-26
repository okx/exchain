package pkg

import (
	"math"
	"time"
)

type operateInfo struct {
	count    int64 `json:"count"`
	timeCost int64 `json:"timeCost"`
	lastCall int64 `json:"lastCall"`
	avg      int64 `json:"avg"`
	min      int64 `json:"min"`
	max      int64 `json:"max"`
}

func newOperateInfo() *operateInfo {
	tmp := &operateInfo{
		lastCall: time.Now().Unix(),
		min:      math.MaxInt64,
	}
	return tmp
}

func (s *operateInfo) StartOper() {
	s.lastCall = time.Now().Unix()
}

func (s *operateInfo) StopOper() {
	callTime := time.Now().Unix() - s.lastCall
	if callTime > s.max {
		s.max = callTime
	}
	if callTime < s.min {
		s.min = callTime
	}
	s.timeCost += callTime
	s.count++
	s.avg = s.timeCost / s.count
}

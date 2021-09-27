package pkg

import (
	"math"
)

type operateInfo struct {
	Count    int64 `json:"count"`
	TimeCost int64 `json:"timeCost"`
	LastCall int64 `json:"lastCall"`
	Avg      int64 `json:"avg"`
	Min      int64 `json:"min"`
	Max      int64 `json:"max"`
}

func newOperateInfo() *operateInfo {
	tmp := &operateInfo{
		LastCall: GetNowTimeNs(),
		Min:      math.MaxInt64,
	}
	return tmp
}

func (s *operateInfo) StartOper() {
	s.LastCall = GetNowTimeNs()
}

func (s *operateInfo) StopOper() {

	callTime := GetNowTimeNs() - s.LastCall

	if callTime > s.Max {
		s.Max = callTime
	}
	if callTime < s.Min {
		s.Min = callTime
	}
	s.TimeCost += callTime
	s.Count++
	s.Avg = s.TimeCost / s.Count

}

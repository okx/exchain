package trace

type operateInfo struct {
	TimeCost int64 `json:"timeCost"`
	LastCall int64 `json:"lastCall"`
	started  bool
}

func newOperateInfo() *operateInfo {
	tmp := &operateInfo{
		LastCall: getNowTimeMs(),
	}
	return tmp
}

func (s *operateInfo) StartOper() {
	if s.started {
		panic("wrong state")
	}
	s.started = true
	s.LastCall = getNowTimeMs()
}

func (s *operateInfo) StopOper() {
	if !s.started {
		panic("wrong state")
	}
	s.started = false
	callTime := getNowTimeMs() - s.LastCall
	s.TimeCost += callTime
}

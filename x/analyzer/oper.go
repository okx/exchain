package analyzer

type operateInfo struct {
	Count    int64 `json:"count"`
	TimeCost int64 `json:"timeCost"`
	LastCall int64 `json:"lastCall"`
}

func newOperateInfo() *operateInfo {
	tmp := &operateInfo{
		LastCall: GetNowTimeMs(),
	}
	return tmp
}

func (s *operateInfo) StartOper() {
	s.LastCall = GetNowTimeMs()
}

func (s *operateInfo) StopOper() {
	callTime := GetNowTimeMs() - s.LastCall
	s.TimeCost += callTime
	s.Count++
}

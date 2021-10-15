package analyzer

type txLog struct {
	startTime int64
	AllCost   int64
	Record    map[string]*operateInfo
}

func newTxLog() *txLog {
	tmp := &txLog{
		startTime: GetNowTimeMs(),
		Record:    make(map[string]*operateInfo),
	}

	return tmp
}

func (s *txLog) StartTxLog(oper string) {
	if _, ok := s.Record[oper]; !ok {
		s.Record[oper] = newOperateInfo()
	}
	s.Record[oper].StartOper()
}

func (s *txLog) StopTxLog(oper string) {
	if _, ok := s.Record[oper]; !ok {
		return
	}
	s.Record[oper].StopOper()
}

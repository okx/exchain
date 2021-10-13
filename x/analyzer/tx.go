package analyzer

import (
	"fmt"
)

type txLog struct {
	startTime int64
	AllCost   int64
	Record map[string]*operateInfo
}

func newTxLog() *txLog {
	tmp := &txLog{
		startTime: GetNowTimeMs(),
		Record:    make(map[string]*operateInfo),
	}

	return tmp
}

func (s *txLog) StartTxLog(oper string) error {

	if _, ok := s.Record[oper]; !ok {
		s.Record[oper] = newOperateInfo()
	}
	s.Record[oper].StartOper()
	return nil
}

func (s *txLog) StopTxLog(oper string) error {

	if _, ok := s.Record[oper]; !ok {
		return fmt.Errorf("%s oper not found", oper)
	}

	s.Record[oper].StopOper()

	return nil
}

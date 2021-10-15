package analyzer

import (
	"fmt"
	"sync"
)

type txLog struct {
	startTime int64
	AllCost   int64
	Record    sync.Map
}

func newTxLog() *txLog {
	tmp := &txLog{
		startTime: GetNowTimeMs(),
	}

	return tmp
}

func (s *txLog) StartTxLog(oper string) error {
	if v, ok := s.Record.Load(oper); !ok {
		newOper := newOperateInfo()
		s.Record.Store(oper, newOper)
		newOper.StartOper()
	}else{
		oper := v.(*operateInfo)
		oper.StartOper()
	}

	return nil
}

func (s *txLog) StopTxLog(oper string) error {
	if v, ok := s.Record.Load(oper); !ok {
		return fmt.Errorf("%s oper not found", oper)
	}else{
		oper := v.(*operateInfo)
		oper.StopOper()
	}

	return nil
}

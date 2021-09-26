package pkg

import (
	"fmt"
	"time"
)

type txLog struct {
	startTime int64
	evmCost   int64
	record    map[string]*txBase
}

func newTxLog(module []string) *txLog {
	tmp := &txLog{
		startTime: time.Now().Unix(),
		record:    make(map[string]*txBase),
	}
	for _, v := range module {
		tmp.record[v] = newTxBase()
	}
	return tmp
}

func (s *txLog) StartTxLog(module, oper string) error {
	if _, ok := s.record[module]; !ok {
		return fmt.Errorf("%s module not found", module)
	}
	s.record[module].StartCost(oper)
	return nil
}

func (s *txLog) StopTxLog(module, oper string) error {
	if _, ok := s.record[module]; !ok {
		return fmt.Errorf("%s module not found", module)
	}
	s.record[module].StopCost(oper)
	//统计evm 耗时
	if v, ok := s.record[COMMIT_STATE_DB]; ok {
		s.evmCost = time.Now().Unix() - s.startTime - v.AllCost()
	}
	return nil
}

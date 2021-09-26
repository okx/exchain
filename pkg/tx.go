package pkg

import (
	"fmt"
)

type txLog struct {
	startTime int64
	EvmCost   int64
	Record    map[string]*txBase
}

func newTxLog(module []string) *txLog {
	tmp := &txLog{
		startTime: GetNowTimeMs(),
		Record:    make(map[string]*txBase),
	}
	for _, v := range module {
		tmp.Record[v] = newTxBase()
	}
	return tmp
}

func (s *txLog) StartTxLog(module, oper string) error {
	if _, ok := s.Record[module]; !ok {
		return fmt.Errorf("%s module not found", module)
	}
	s.Record[module].StartCost(oper)
	return nil
}

func (s *txLog) StopTxLog(module, oper string) error {
	if _, ok := s.Record[module]; !ok {
		return fmt.Errorf("%s module not found", module)
	}
	s.Record[module].StopCost(oper)
	//统计evm 耗时
	if v, ok := s.Record[COMMIT_STATE_DB]; ok {
		s.EvmCost = GetNowTimeMs() - s.startTime - v.AllCost()
	}
	return nil
}

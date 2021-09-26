package pkg

import (
	"errors"
	"fmt"
	"sync"
)

type txBase struct {
	lock   sync.RWMutex
	record map[string]*operateInfo
}

func newTxBase() *txBase {
	return &txBase{
		record: make(map[string]*operateInfo),
	}
}

func (s *txBase) StartCost(oper string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if v, ok := s.record[oper]; !ok {
		s.record[oper] = newOperateInfo()
	} else {
		v.StartOper()
	}
}

func (s *txBase) StopCost(oper string) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	if _, ok := s.record[oper]; !ok {
		return errors.New(fmt.Sprintf("%s not Start", oper))
	}
	s.record[oper].StopOper()
	return nil
}

func (s *txBase) Format() map[string]*operateInfo {
	return s.record
}

func (s *txBase) AllCost() int64 {
	var res int64
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, v := range s.record {
		res += v.timeCost
	}
	return res
}

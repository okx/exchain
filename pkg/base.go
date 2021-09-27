package pkg

import (
	"errors"
	"fmt"
	"sync"
)

type txBase struct {
	lock   sync.RWMutex
	Record map[string]*operateInfo
}

func newTxBase() *txBase {
	return &txBase{
		Record: make(map[string]*operateInfo),
	}
}

func (s *txBase) StartCost(oper string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.Record[oper]; !ok {

		s.Record[oper] = newOperateInfo()
	}

	s.Record[oper].StartOper()

}

func (s *txBase) StopCost(oper string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.Record[oper]; !ok {
		return errors.New(fmt.Sprintf("%s not Start", oper))
	}
	s.Record[oper].StopOper()
	return nil
}

func (s *txBase) AllCost() int64 {
	var res int64
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, v := range s.Record {
		res += v.TimeCost
	}
	return res
}

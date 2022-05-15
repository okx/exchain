package trace

import (
	"sync"
	"time"
)


func getNowTimeMs() int64 {
	return time.Now().UnixNano() / 1e6
}

type DbRecord struct {
	lock sync.RWMutex
	oper map[string]int
}

func newDbRecord() *DbRecord {
	return &DbRecord{
		oper: make(map[string]int),
	}
}

func (s *DbRecord) GetOperType(oper string) int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if _, ok := s.oper[oper]; !ok {
		return -1
	}
	return s.oper[oper]
}

func (s *DbRecord) AddOperType(oper string, value int) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.oper[oper] = value
}

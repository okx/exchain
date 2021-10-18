package analyzer

import (
	"runtime"
	"strings"
	"sync"
	"time"
)

func RunFuncName() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	names := strings.Split(f.Name(), ".")
	return names[len(names)-1]
}

func GetNowTimeMs() int64 {
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

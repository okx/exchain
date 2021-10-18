package cache

import "github.com/ethereum/go-ethereum/common"

type Monitor struct {
	cacheData map[common.Hash][]byte
}

func NewMonitor() *Monitor {
	return &Monitor{
		cacheData: make(map[common.Hash][]byte),
	}
}

func (m *Monitor) SetState(key common.Hash, value []byte) {
	m.cacheData[key] = value
}

func (m *Monitor) Empty() {
	m.cacheData = make(map[common.Hash][]byte)
}

func (m *Monitor) Iterator(f func(key common.Hash, value []byte)) {
	if f == nil {
		return
	}
	for k, v := range m.cacheData {
		f(k, v)
	}
}

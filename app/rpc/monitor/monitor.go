package monitor

import (
	"fmt"
	"github.com/go-kit/kit/metrics"
	"github.com/tendermint/tendermint/libs/log"
	"time"
)

type Monitor struct {
	method        string
	logger        log.Logger
	lastTimestamp int64
}

func GetMonitor(method string, logger log.Logger) *Monitor {
	return &Monitor{
		method: method,
		logger: logger,
	}
}

func (m *Monitor) OnBegin(metrics map[string]metrics.Counter) {
	m.lastTimestamp = time.Now().UnixNano()

	if metrics == nil {
		return
	}
	if _, ok := metrics[m.method]; ok {
		metrics[m.method].Add(1)
	}
}

func (m *Monitor) OnEnd(args ...interface{}) {
	now := time.Now().UnixNano()
	m.logger.Debug(fmt.Sprintf("RPC: Method<%s>, Interval<%dms>, Params<%v>", m.method, (now-m.lastTimestamp)/1000, args))
}

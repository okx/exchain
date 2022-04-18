package watcher

import (
	"github.com/okex/exchain/libs/tendermint/libs/log"
	stdlog "log"
	"testing"
	"time"
)

type mylogger struct {
}

func (m mylogger) Debug(msg string, keyvals ...interface{}) {}

func (m mylogger) Info(msg string, keyvals ...interface{}) {}

func (m mylogger) Error(msg string, keyvals ...interface{}) {
	stdlog.Println(msg, keyvals)
}

func (m mylogger) With(keyvals ...interface{}) log.Logger { return m }

func TestWatcher_dispatchJob(t *testing.T) {
	var logger mylogger
	watcher := &Watcher{
		log: logger,
		sw:  true,
	}
	go watcher.jobRoutine()
	time.Sleep(time.Microsecond)

	const JobChanBuffer = 15
	for i := 0; i < JobChanBuffer*2; i++ {
		index := i
		watcher.dispatchJob(func() {
			t.Logf("fired %v \n", index)
		})
	}
	time.Sleep(time.Millisecond)
	t.Log("test the error log")
}

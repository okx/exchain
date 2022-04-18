package watcher

import (
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/stretchr/testify/suite"
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

type WatcherTestSuite struct {
	suite.Suite
	watcher Watcher
}

func (suite *WatcherTestSuite) SetupTest() {
	var logger mylogger
	suite.watcher = Watcher{
		log: logger,
		sw:  true,
	}
	go suite.watcher.jobRoutine()
	suite.watcher.txRoutine()
	time.Sleep(time.Millisecond)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(WatcherTestSuite))
}

func (suite *WatcherTestSuite) TestWatcher_dispatchJob_ErrorMsg() {
	const JobChanBuffer = 15
	for i := 0; i < JobChanBuffer*2; i++ {
		index := i
		suite.watcher.dispatchJob(func() {
			suite.T().Logf("fired %v \n", index)
			time.Sleep(10 * time.Microsecond)
		})
	}
	time.Sleep(time.Millisecond)
	suite.T().Log("test the error log: watch dispatch job too busy.")
}

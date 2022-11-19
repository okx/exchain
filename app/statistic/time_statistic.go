package statistic

import (
	"fmt"
	"time"

	"github.com/okex/exchain/libs/tendermint/libs/log"
)

type TimeStatistic struct {
	blockTimeCh chan<- blockTimeInfo
	exitCh      chan<- struct{}
}

type btchType uint8

const (
	BlockBeginTime btchType = iota
	BlockEndTime
)

type blockTimeInfo struct {
	tType  btchType
	height int64
	time   time.Time
}

var now func() time.Time = time.Now

func NewTimeStatistic(appStartTime time.Time, logger log.Logger) *TimeStatistic {
	exitCh := make(chan struct{})
	blockTimeCh := make(chan blockTimeInfo)

	go timeStatisticLoop(appStartTime, blockTimeCh, exitCh, logger)

	return &TimeStatistic{blockTimeCh, exitCh}
}

func (ts *TimeStatistic) BeginBlock(height int64) {
	ts.blockTimeCh <- blockTimeInfo{BlockBeginTime, height, now()}
}

func (ts *TimeStatistic) EndBlock(height int64) {
	ts.blockTimeCh <- blockTimeInfo{BlockEndTime, height, now()}
}

func (ts *TimeStatistic) Exit() {
	close(ts.exitCh)
}

func timeStatisticLoop(appStartTime time.Time, blockTimeCh <-chan blockTimeInfo, exitCh <-chan struct{}, logger log.Logger) {
	blockBeginTimes := make(map[int64]time.Time)
	blockTotalTime := time.Duration(0)

	for {
		select {
		case <-exitCh:
			logger.Info("Time Statistic Loop exit")
			return
		case bt := <-blockTimeCh:
			switch bt.tType {
			case BlockBeginTime:
				blockBeginTimes[bt.height] = bt.time
			case BlockEndTime:
				beginTime, ok := blockBeginTimes[bt.height]
				if !ok {
					logger.Error(fmt.Sprintf("no block height %d when receiving a BlockEndTime message", bt.height))
					continue
				}

				blockTime := bt.time.Sub(beginTime)
				appStartedTime := now().Sub(appStartTime)
				blockTotalTime += blockTime
				ratio := float64(blockTotalTime.Nanoseconds()) / float64(appStartedTime.Nanoseconds()) * 100
				logger.Info("Block Execution Time Statistic", "Current Height", bt.height, "Current Block Time", blockTime, "Total Block Time", blockTotalTime, "Node Started Time", appStartedTime, "Total Block Time Ratio", fmt.Sprintf("%.2f%%", ratio))
			}
		}
	}
}

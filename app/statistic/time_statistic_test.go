package statistic

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/okex/exchain/libs/tendermint/libs/log"
)

type mockNow struct {
	nowTime time.Time
	lock    sync.Mutex
}

func newMockNow() *mockNow {
	return &mockNow{
		nowTime: time.Now(),
	}
}

func (mn *mockNow) add(d time.Duration) {
	mn.lock.Lock()
	defer mn.lock.Unlock()

	mn.nowTime = mn.nowTime.Add(d)
}

func (mn *mockNow) setNow(t time.Time) {
	mn.lock.Lock()
	defer mn.lock.Unlock()

	mn.nowTime = t
}

func (mn *mockNow) now() time.Time {
	mn.lock.Lock()
	defer mn.lock.Unlock()

	return mn.nowTime
}

func TestBlockTimeStatisticOnce(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := log.NewTMLogger(buf)

	oldNow := now
	defer func() { now = oldNow }()

	mnow := newMockNow()
	now = mnow.now

	const HEIGHT = 1
	const START_MORE_DELTA = time.Second
	const BLOCK_DELTA = time.Second
	const TOTAL_START_DELTA = START_MORE_DELTA + BLOCK_DELTA
	const RATIO = (float64(BLOCK_DELTA) / float64(TOTAL_START_DELTA)) * 100

	s := NewTimeStatistic(now(), logger)
	defer s.Exit()
	mnow.add(START_MORE_DELTA)

	s.BeginBlock(1)
	mnow.add(BLOCK_DELTA)
	s.EndBlock(1)
	time.Sleep(time.Second)

	expectOutput := fmt.Sprintf("Block Execution Time Statistic. CurrentHeight=%d CurrentBlockTime=%v TotalBlockTime=%v NodeStartedTime=%v TotalBlockTimeRatio=%.2f%%", HEIGHT, BLOCK_DELTA, BLOCK_DELTA, BLOCK_DELTA+START_MORE_DELTA, RATIO)

	if !strings.Contains(buf.String(), expectOutput) {
		t.Fatalf("expect time statistic output \"%s\" but got \"%s\"\n", expectOutput, buf.String())
	}
}

func TestBlockTimeStatisticMulti(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	logger := log.NewTMLogger(buf)

	oldNow := now
	defer func() { now = oldNow }()

	mnow := newMockNow()
	now = mnow.now

	const HEIGHT = 1
	const START_MORE_DELTA = time.Second
	const BLOCK_DELTA = time.Second
	const BLOCK_COUNT = 3

	s := NewTimeStatistic(now(), logger)
	defer s.Exit()
	mnow.add(START_MORE_DELTA)

	for i := int64(1); i < BLOCK_COUNT+1; i++ {
		s.BeginBlock(i)
		mnow.add(BLOCK_DELTA)
		s.EndBlock(i)
		time.Sleep(time.Second)

		totalBlockTime := time.Duration(i * BLOCK_DELTA.Nanoseconds())
		nodeStarted := totalBlockTime + START_MORE_DELTA
		ratio := (float64(totalBlockTime) / float64(nodeStarted)) * 100
		expectOutput := fmt.Sprintf("Block Execution Time Statistic. CurrentHeight=%d CurrentBlockTime=%v TotalBlockTime=%v NodeStartedTime=%v TotalBlockTimeRatio=%.2f%%", i, BLOCK_DELTA, totalBlockTime, nodeStarted, ratio)
		if !strings.Contains(buf.String(), expectOutput) {
			t.Fatalf("loop %d: expect time statistic output \"%s\" but got \"%s\"\n", i, expectOutput, buf.String())
		}
	}

}

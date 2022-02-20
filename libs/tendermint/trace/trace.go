package trace

import (
	"fmt"
	"time"

	"github.com/okex/exchain/libs/tendermint/libs/log"
)

const (
	GasUsed     = "GasUsed"
	Produce     = "Produce"
	RunTx       = "RunTx"
	Height      = "Height"
	Tx          = "Tx"
	BlockSize   = "BlockSize"
	Elapsed     = "Elapsed"
	CommitRound = "CommitRound"
	Round       = "Round"
	Evm         = "Evm"
	Iavl        = "Iavl"
	FlatKV      = "FlatKV"
	WtxRatio    = "WtxRatio"
	DeliverTxs  = "DeliverTxs"
	AnteHandler = "AnthHandler"

	Delta = "Delta"

	Abci       = "abci"
	InvalidTxs = "InvalidTxs"
	SaveResp   = "saveResp"
	Persist    = "persist"
	SaveState  = "saveState"

	ApplyBlock = "ApplyBlock"
	Consensus  = "Consensus"

	MempoolCheckTxCnt = "checkTxCnt"
	MempoolTxsCnt     = "mempoolTxsCnt"

	Prerun = "Prerun"
)

type IElapsedTimeInfos interface {
	AddInfo(key string, info string)
	Dump(logger log.Logger)
	SetElapsedTime(elapsedTime int64)
	GetElapsedTime() int64
}

func SetInfoObject(e IElapsedTimeInfos) {
	if e != nil {
		elapsedInfo = e
	}
}

var elapsedInfo IElapsedTimeInfos = &EmptyTimeInfo{}

func GetElapsedInfo() IElapsedTimeInfos {
	return elapsedInfo
}

type Tracer struct {
	name             string
	startTime        int64
	lastPin          string
	lastPinStartTime int64
	pins             []string
	pinMap           map[string]int64
	intervals        []int64
	elapsedTime      int64
	ignoredTags       string
	ignoreOverallElapsed bool
}

func NewTracer(name string) *Tracer {
	t := &Tracer{
		startTime: time.Now().UnixNano(),
		name:      name,
		pinMap:    make(map[string]int64),
	}
	return t
}

func (t *Tracer) Pin(format string, args ...interface{}) {
	t.pinByFormat(fmt.Sprintf(format, args...))
}

func (t *Tracer) pinByFormat(tag string) {
	if len(tag) == 0 {
		//panic("invalid tag")
		return
	}

	if len(t.pins) > 100 {
		// 100 pins limitation
		return
	}

	now := time.Now().UnixNano()

	if len(t.lastPin) > 0 {
		t.pins = append(t.pins, t.lastPin)
		t.intervals = append(t.intervals, (now-t.lastPinStartTime)/1e6)
	}
	t.lastPinStartTime = now
	t.lastPin = tag
}

func (t *Tracer) Format() string {
	if len(t.pins) == 0 {
		now := time.Now().UnixNano()
		t.elapsedTime = (now - t.startTime) / 1e6
		return fmt.Sprintf("%s<%dms>",
			t.name,
			t.elapsedTime,
		)
	}

	t.Pin("_")

	now := time.Now().UnixNano()
	t.elapsedTime = (now - t.startTime) / 1e6
	info := fmt.Sprintf("%s<%dms>",
		t.name,
		t.elapsedTime,
		)

	for i := range t.pins {
		info += fmt.Sprintf(", %s<%dms>", t.pins[i], t.intervals[i])
	}
	return info
}

func (t *Tracer) SetIgnoredTag(tag string) {
	t.ignoredTags = tag
}

func (t *Tracer) SetIgnoreOverallElapsed() {
	t.ignoreOverallElapsed = true
}



func (t *Tracer) RepeatingPin(format string, args ...interface{}) {
	t.repeatingPinByFormat(fmt.Sprintf(format, args...))
}

func (t *Tracer) repeatingPinByFormat(tag string) {
	if len(tag) == 0 {
		//panic("invalid tag")
		return
	}

	if len(t.pinMap) > 100 {
		// 100 pins limitation
		return
	}

	now := time.Now().UnixNano()

	if len(t.lastPin) > 0 {
		t.pinMap[t.lastPin] += (now-t.lastPinStartTime)/1e6
	}
	t.lastPinStartTime = now
	t.lastPin = tag
}

func (t *Tracer) FormatRepeatingPins() string {
	if len(t.pinMap) == 0 {
		now := time.Now().UnixNano()
		t.elapsedTime = (now - t.startTime) / 1e6
		return fmt.Sprintf("%s<%dms>",
			t.name,
			t.elapsedTime,
		)
	}

	t.RepeatingPin("_")

	var info, comma string

	for tag, interval := range t.pinMap {
		if tag == t.ignoredTags {
			continue
		}
		info += fmt.Sprintf("%s%s<%dms>", comma, tag, interval)
		comma = ", "
	}
	return info
}

func (t *Tracer) GetElapsedTime() int64 {
	return t.elapsedTime
}

func (t *Tracer) Reset() {
	t.startTime = time.Now().UnixNano()
	t.lastPin = ""
	t.lastPinStartTime = 0
	t.pins = nil
	t.intervals = nil
}

type EmptyTimeInfo struct {
}

func (e *EmptyTimeInfo) AddInfo(key string, info string) {
}

func (e *EmptyTimeInfo) Dump(logger log.Logger) {
}

func (e *EmptyTimeInfo) SetElapsedTime(elapsedTime int64) {
}

func (e *EmptyTimeInfo) GetElapsedTime() int64 {
	return 0
}

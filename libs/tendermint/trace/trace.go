package trace

import (
	"fmt"
	"time"

	"github.com/okex/exchain/libs/tendermint/libs/log"
)

const (
	GasUsed          = "GasUsed"
	Produce          = "Produce"
	RunTx            = "RunTx"
	Height           = "Height"
	Tx               = "Tx"
	BlockSize        = "BlockSize"
	Elapsed          = "Elapsed"
	CommitRound      = "CommitRound"
	Round            = "Round"
	Evm              = "Evm"
	Iavl             = "Iavl"
	FlatKV           = "FlatKV"
	WtxRatio         = "WtxRatio"
	DeliverTxs       = "DeliverTxs"
	EvmHandlerDetail = "EvmHandlerDetail"
	RunAnteDetail    = "RunAnteDetail"
	AnteChainDetail  = "AnteChainDetail"

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
	startTime        time.Time
	lastPin          string
	lastPinStartTime time.Time
	pins             []string
	intervals        []time.Duration
	elapsedTime      time.Duration

	pinMap map[string]time.Duration
}

func NewTracer(name string) *Tracer {
	t := &Tracer{
		startTime: time.Now(),
		name:      name,
		pinMap:    make(map[string]time.Duration),
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

	now := time.Now()

	if len(t.lastPin) > 0 {
		t.pins = append(t.pins, t.lastPin)
		t.intervals = append(t.intervals, now.Sub(t.lastPinStartTime))
	}
	t.lastPinStartTime = now
	t.lastPin = tag
}

func (t *Tracer) Format() string {
	if len(t.pins) == 0 {
		now := time.Now()
		t.elapsedTime = now.Sub(t.startTime)
		return fmt.Sprintf("%s<%dms>",
			t.name,
			t.elapsedTime,
		)
	}

	t.Pin("_")

	now := time.Now()
	t.elapsedTime = now.Sub(t.startTime)
	info := fmt.Sprintf("%s<%dms>",
		t.name,
		t.elapsedTime,
	)

	for i := range t.pins {
		info += fmt.Sprintf(", %s<%dms>", t.pins[i], t.intervals[i].Milliseconds())
	}
	return info
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

	now := time.Now()

	if len(t.lastPin) > 0 {
		t.pinMap[t.lastPin] += now.Sub(t.lastPinStartTime)
	}
	t.lastPinStartTime = now
	t.lastPin = tag
}

func (t *Tracer) FormatRepeatingPins(ignoredTags string) string {
	var info, comma string

	if len(t.pinMap) == 0 {
		return info
	}

	t.RepeatingPin("_")

	for tag, interval := range t.pinMap {
		if tag == ignoredTags {
			continue
		}
		info += fmt.Sprintf("%s%s<%dms>", comma, tag, interval.Milliseconds())
		comma = ", "
	}
	return info
}

func (t *Tracer) GetElapsedTime() int64 {
	return t.elapsedTime.Milliseconds()
}

func (t *Tracer) Reset() {
	t.startTime = time.Now()
	t.lastPin = ""
	t.lastPinStartTime = time.Date(2018, 1, 1, 1, 1, 1, 1, time.Local)
	t.pins = nil
	t.intervals = nil
	t.pinMap = make(map[string]time.Duration)
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

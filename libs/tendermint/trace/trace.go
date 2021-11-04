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
	Elapsed     = "Elapsed"
	CommitRound = "CommitRound"
	Round       = "Round"
	Evm         = "Evm"
	Iavl        = "Iavl"
	DeliverTxs  = "DeliverTxs"


	Abci       = "abci"
	SaveResp   = "saveResp"
	Persist    = "persist"
	SaveState  = "saveState"

	ApplyBlock = "ApplyBlock"
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

func NewTracer(name string) *Tracer {
	t := &Tracer{
		startTime: time.Now().UnixNano(),
		name: name,
	}
	return t
}

type Tracer struct {
	name             string
	startTime        int64
	lastPin          string
	lastPinStartTime int64
	pins             []string
	intervals        []int64
	elapsedTime      int64
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
		return ""
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

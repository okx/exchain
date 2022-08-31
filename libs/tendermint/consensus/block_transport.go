package consensus

import (
	"fmt"
	"github.com/okex/exchain/libs/system/trace"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"sync"
	"time"
)

type BlockTransport struct {
	height                 int64
	recvProposal           time.Time
	firstPart              time.Time
	droppedDue2NotExpected int
	droppedDue2NotAdded    int
	droppedDue2Error       int
	droppedDue2WrongHeight int
	totalParts             int
	Logger                 log.Logger

	bpStatMtx       sync.RWMutex
	bpSend          int
	bpNOTransByData int
	bpNOTransByACK  int

	recvPrevote  bool
	firstPrevote time.Time
}

func (bt *BlockTransport) onProposal(height int64) {
	if bt.height == height || bt.height == 0 {
		bt.recvProposal = time.Now()
		bt.height = height
	}
}

func (bt *BlockTransport) reset(height int64) {
	bt.height = height
	bt.droppedDue2NotExpected = 0
	bt.droppedDue2NotAdded = 0
	bt.droppedDue2Error = 0
	bt.droppedDue2WrongHeight = 0
	bt.totalParts = 0
	bt.bpNOTransByData = 0
	bt.bpNOTransByACK = 0
	bt.bpSend = 0
	bt.recvPrevote = false
}

func (bt *BlockTransport) on1stPart(height int64) {
	if bt.height == height || bt.height == 0 {
		bt.firstPart = time.Now()
		bt.height = height
	}
}

func (bt *BlockTransport) onRecvBlock(height int64) {
	if bt.height == height {
		//totalElapsed := time.Now().Sub(bt.recvProposal)
		//trace.GetElapsedInfo().AddInfo(trace.RecvBlock, fmt.Sprintf("<%dms>", totalElapsed.Milliseconds()))
		first2LastPartElapsed := time.Now().Sub(bt.firstPart)
		trace.GetElapsedInfo().AddInfo(trace.First2LastPart, fmt.Sprintf("%dms", first2LastPartElapsed.Milliseconds()))
	}
}

// blockpart send times
func (bt *BlockTransport) onBPSend() {
	bt.bpStatMtx.Lock()
	bt.bpSend++
	bt.bpStatMtx.Unlock()
}

// blockpart-ack receive times, specific blockpart won't send  to the peer from where the ack fired
func (bt *BlockTransport) onBPACKHit() {
	bt.bpStatMtx.Lock()
	bt.bpNOTransByACK++
	bt.bpStatMtx.Unlock()
}

// blockpart data receive times, specific blockpart won't send to the peer from where the data fired
func (bt *BlockTransport) onBPDataHit() {
	bt.bpStatMtx.Lock()
	bt.bpNOTransByData++
	bt.bpStatMtx.Unlock()
}

//prevote time
func (bt *BlockTransport) On1stPrevote(height int64) {
	if (bt.height == height || bt.height == 0) && !bt.recvPrevote {
		bt.recvPrevote = true
		bt.firstPrevote = time.Now()
	}
}

func (bt *BlockTransport) on23AnyPrevote(height int64) {
	if bt.height == height {
		first2AnyElapsed := time.Now().Sub(bt.firstPrevote)
		trace.GetElapsedInfo().AddInfo(trace.FirstTo23AnyPrevote, fmt.Sprintf("%dms", first2AnyElapsed.Milliseconds()))
	}
}

func (bt *BlockTransport) on23MajPrevote(height int64) {
	if bt.height == height {
		first2MajElapsed := time.Now().Sub(bt.firstPrevote)
		trace.GetElapsedInfo().AddInfo(trace.FirstTo23MajPrevote, fmt.Sprintf("%dms", first2MajElapsed.Milliseconds()))
	}
}

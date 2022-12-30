package consensus

import (
	"fmt"
	"github.com/okex/exchain/libs/system/trace"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/p2p"
	"github.com/okex/exchain/libs/tendermint/types"
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

	lastBPSender        string
	lastPrevoteSender   string
	lastPrecommitSender string
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
	bt.lastBPSender = "-"
	bt.lastPrecommitSender = "-"
	bt.lastPrevoteSender = "-"
}

func (bt *BlockTransport) on1stPart(height int64) {
	if bt.height == height || bt.height == 0 {
		bt.firstPart = time.Now()
		bt.height = height
	}
}

func (bt *BlockTransport) onRecvBlock(height int64, peerID p2p.ID, selfAddr types.Address) {
	if bt.height == height {
		//totalElapsed := time.Now().Sub(bt.recvProposal)
		//trace.GetElapsedInfo().AddInfo(trace.RecvBlock, fmt.Sprintf("<%dms>", totalElapsed.Milliseconds()))
		first2LastPartElapsed := time.Now().Sub(bt.firstPart)
		trace.GetElapsedInfo().AddInfo(trace.First2LastPart, fmt.Sprintf("%dms", first2LastPartElapsed.Milliseconds()))
		var peerAddress types.Address
		if peerID == "" {
			peerAddress = selfAddr
		} else if v, ok := pID2Pubkey.Load(peerID); ok {
			peerAddress = v.(types.Address)
		}

		if peerAddress != nil {
			bt.lastBPSender = peerAddress.String()[:6]
		}
	}
}

func (bt *BlockTransport) onLastPrevote(height int64, peerID p2p.ID, voteAddr types.Address, selfAddr types.Address) {
	if bt.height != height || voteAddr == nil {
		return
	}

	var peerAddress types.Address
	if peerID == "" {
		peerAddress = selfAddr
	} else if v, ok := pID2Pubkey.Load(peerID); ok {
		peerAddress = v.(types.Address)
	}

	if peerAddress != nil {
		bt.lastPrevoteSender = fmt.Sprintf("%s<%s>", peerAddress.String()[:6], voteAddr.String()[:6])
	}

}

func (bt *BlockTransport) onLastPrecommit(height int64, peerID p2p.ID, voteAddr types.Address, selfAddr types.Address) {
	if bt.height != height || voteAddr == nil {
		return
	}

	var peerAddress types.Address
	if peerID == "" {
		peerAddress = selfAddr
	} else if v, ok := pID2Pubkey.Load(peerID); ok {
		peerAddress = v.(types.Address)
	}

	if peerAddress != nil {
		bt.lastPrecommitSender = fmt.Sprintf("%s<%s>", peerAddress.String()[:6], voteAddr.String()[:6])
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

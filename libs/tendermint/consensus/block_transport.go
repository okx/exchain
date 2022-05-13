package consensus

import (
	"fmt"
	"github.com/okex/exchain/libs/system/trace"
	"time"
)

type BlockTransport struct {
	height int64
	recvProposal time.Time
	firstPart time.Time
	totalElapsed time.Duration
	first2LastPartElapsed time.Duration
}

func (bt *BlockTransport) onProposal(height int64)  {
	if bt.height+1 == height || bt.height == 0 {
		bt.recvProposal = time.Now()
		bt.height = height
	} else {
		//panic("invalid height")
	}
}

func (bt *BlockTransport) on1stPart(height int64)  {
	if bt.height+1 == height || bt.height == height || bt.height == 0 {
		bt.firstPart = time.Now()
		bt.height = height
	} else {
		//panic("invalid height")
	}
}

func (bt *BlockTransport) onRecvBlock(height int64)  {
	if bt.height == height {
		bt.totalElapsed = time.Now().Sub(bt.recvProposal)
		bt.first2LastPartElapsed = time.Now().Sub(bt.firstPart)
		trace.GetElapsedInfo().AddInfo(trace.RecvBlock,
			fmt.Sprintf("%d<%dms>", height, bt.totalElapsed.Milliseconds()))
		trace.GetElapsedInfo().AddInfo(trace.First2LastPart,
			fmt.Sprintf("%d<%dms>", height, bt.first2LastPartElapsed.Milliseconds()))
	} else {
		//panic("invalid height")
	}
}

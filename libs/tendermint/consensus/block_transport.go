package consensus

import (
	"fmt"
	"github.com/okex/exchain/libs/system/trace"
	"time"
)

type BlockTransport struct {
	height int64
	recvProposal time.Time
	elapsed time.Duration
}

func (bt *BlockTransport) onProposal(height int64)  {
	if bt.height+1 == height || bt.height == 0 {
		bt.recvProposal = time.Now()
		bt.height = height
	} else {
		//panic("invalid height")
	}
}

func (bt *BlockTransport) onRecvBlock(height int64)  {
	if bt.height == height {
		bt.elapsed = time.Now().Sub(bt.recvProposal)
		trace.GetElapsedInfo().AddInfo(trace.RecvBlock,
			fmt.Sprintf("%d<%dms>", height, bt.elapsed.Milliseconds()))
	} else {
		//panic("invalid height")
	}
}

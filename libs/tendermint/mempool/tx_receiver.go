package mempool

import (
	"context"
	"fmt"
	pb "github.com/okex/exchain/libs/tendermint/proto/mempool"
	"google.golang.org/protobuf/types/known/emptypb"
)

type txReceiverServer struct {
	pb.UnimplementedMempoolTxReceiverServer
	memR    *Reactor
	Port    int
	Started int64
}

func newTxReceiverServer(memR *Reactor) *txReceiverServer {
	return &txReceiverServer{memR: memR}
}

func (s *txReceiverServer) Enabled() bool {
	if s == nil || s.Port == 0 || atomic.LoadInt64(&s.Started) == 0 {
		return false
	}
	return true
}

func (s *txReceiverServer) Receive(ctx context.Context, req *pb.TxRequest) (*emptypb.Empty, error) {
	if len(req.Tx) > 0 {
		var txjob txJob
		txjob.tx = req.Tx
		if req.PeerId != 0 {
			txjob.info.SenderID = uint16(req.PeerId)
		} else {
			return nil, errZeroId
		}
		select {
		case s.memR.txCh <- txjob:
		case <-ctx.Done():
		}
	} else {
		s.memR.Logger.Error("txReceiverServer.Receive empty tx")
		return nil, errEmpty
	}

	return empty, nil
}

var empty = &emptypb.Empty{}

var errEmpty = fmt.Errorf("empty tx")
var errZeroId = fmt.Errorf("peerId is 0")

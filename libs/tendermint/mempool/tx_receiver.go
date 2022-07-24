package mempool

import (
	"context"
	pb "github.com/okex/exchain/libs/tendermint/proto/mempool"
	"google.golang.org/protobuf/types/known/emptypb"
)

type txReceiverServer struct {
	pb.UnimplementedMempoolTxReceiverServer
	memR *Reactor
	Port int
}

func newTxReceiverServer(memR *Reactor) *txReceiverServer {
	return &txReceiverServer{memR: memR}
}

func (s *txReceiverServer) Receive(ctx context.Context, req *pb.TxRequest) (*emptypb.Empty, error) {
	if len(req.Tx) > 0 {
		var txjob txJob
		txjob.tx = req.Tx
		if req.PeerId != 0 {
			txjob.info.SenderID = uint16(req.PeerId)
		}
		select {
		case s.memR.txCh <- txjob:
		case <-ctx.Done():
		}
	}

	return empty, nil
}

var empty = &emptypb.Empty{}

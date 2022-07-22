package mempool

import (
	"context"
	pb "github.com/okex/exchain/libs/tendermint/proto/mempool"
	"google.golang.org/protobuf/types/known/emptypb"
)

type txReceiverServer struct {
	pb.UnimplementedMempoolTxReceiverServer
	memR *Reactor
}

func NewTxReceiverServer(memR *Reactor) pb.MempoolTxReceiverServer {
	return &txReceiverServer{memR: memR}
}

func (s *txReceiverServer) Receive(ctx context.Context, req *pb.TxRequest) (*emptypb.Empty, error) {
	if len(req.Tx) > 0 {
		var txjob txJob
		txjob.tx = req.Tx
		if req.PeerId != 0 {
			txjob.info.SenderID = uint16(req.PeerId)
		}
		s.memR.txCh <- txjob
	}

	return nil, nil
}

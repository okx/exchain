package mempool

import (
	"context"
	"fmt"
	pb "github.com/okex/exchain/libs/tendermint/proto/mempool"
	"github.com/okex/exchain/libs/tendermint/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync/atomic"
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

var num1, num2 uint64

func (s *txReceiverServer) Receive(ctx context.Context, req *pb.TxsRequest) (*emptypb.Empty, error) {
	for _, tx := range req.Txs {
		txInfo := TxInfo{
			SenderID: uint16(req.PeerId),
		}

		if atomic.AddUint64(&num1, 1)%1000 == 0 {
			fmt.Println("continue receiving...", "size", s.memR.mempool.Size())
		}
		if err := s.memR.mempool.CheckTx(tx, nil, txInfo); err != nil && err != ErrTxInCache {
			fmt.Println("checkTx error", err)
			return nil, err
		} else if err == nil {
			atomic.AddUint64(&num2, 1)
		}
	}
	return empty, nil
}

func (s *txReceiverServer) ReceiveSentry(ctx context.Context, req *pb.SentryTxs) (*emptypb.Empty, error) {
	for _, stx := range req.Txs {
		types.SignatureCache().Add(types.Tx(stx.Tx).Hash(s.memR.mempool.Height()), stx.From)
		s.memR.txMap.Add(stx.TxIndex, stx.Tx)
	}
	s.memR.mempool.mustNotifyTxsAvailable()
	return empty, nil
}

func (s *txReceiverServer) TxIndices(ctx context.Context, req *pb.IndicesRequest) (*pb.IndicesResponse, error) {
	indices := s.memR.mempool.ReapTxIndicesMaxBytesMaxGas(req.MaxBytes, req.MaxGas)
	return &pb.IndicesResponse{Indices: indices}, nil
}

var empty = &emptypb.Empty{}

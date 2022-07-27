package mempool

import (
	"context"
	"encoding/hex"
	"fmt"
	pb "github.com/okex/exchain/libs/tendermint/proto/mempool"
	"github.com/okex/exchain/libs/tendermint/types"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
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

var num1, num2, num3 uint64
var dupTx = make(map[string]int64)
var dupMtx sync.Mutex

func (s *txReceiverServer) Receive(ctx context.Context, req *pb.TxsRequest) (*emptypb.Empty, error) {
	for _, tx := range req.Txs {
		txInfo := TxInfo{
			SenderID: uint16(req.PeerId),
		}

		if atomic.AddUint64(&num1, 1)%1000 == 0 {
			dupMtx.Lock()
			fmt.Println("mempool size", s.memR.mempool.Size(), "batch size", len(req.Txs), atomic.LoadUint64(&num1), atomic.LoadUint64(&num2), atomic.LoadUint64(&num3), "dup len", len(dupTx))
			dupMtx.Unlock()
		}

		if err := s.memR.mempool.CheckTx(tx, nil, txInfo); err != nil && err != ErrTxInCache {
			fmt.Println("checkTx error", err)
			return nil, err
		} else if err == nil {
			atomic.AddUint64(&num2, 1)
		} else if err == ErrTxInCache {
			dupMtx.Lock()
			dupTx[hex.EncodeToString(tx)]++
			dupMtx.Unlock()
			atomic.AddUint64(&num3, 1)
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

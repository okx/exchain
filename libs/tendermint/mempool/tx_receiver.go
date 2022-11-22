package mempool

import (
	"context"
	"encoding/hex"
	"fmt"
	gogocodec "github.com/gogo/protobuf/codec"
	gogoproto "github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/p2p"
	pb "github.com/okex/exchain/libs/tendermint/proto/mempool"
	"github.com/okex/exchain/libs/tendermint/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding"
	"google.golang.org/grpc/encoding/proto"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func init() {
	encoding.RegisterCodec(newGogoCodec())
}

type txReceiverClient struct {
	Client pb.MempoolTxReceiverClient
	Conn   *grpc.ClientConn
	ID     uint16
}

type txReceiverServer struct {
	pb.UnimplementedMempoolTxReceiverServer
	memR    *Reactor
	Port    int
	Started int64
	Logger  log.Logger
}

func newTxReceiverServer(memR *Reactor) *txReceiverServer {
	return &txReceiverServer{memR: memR}
}

func (s *txReceiverServer) SetLogger(l log.Logger) {
	s.Logger = l
}

func (s *txReceiverServer) Enabled() bool {
	if s == nil || s.Port == 0 || atomic.LoadInt64(&s.Started) == 0 {
		return false
	}
	return true
}

func (s *txReceiverServer) CheckTx(_ context.Context, req *pb.TxRequest) (*pb.Empty, error) {
	return s.checkTx(req, nil)
}

func (s *txReceiverServer) CheckTxAsync(_ context.Context, req *pb.TxRequest) (*pb.Empty, error) {
	return s.checkTx(req, s.memR.getTxJobChannel(req.From))
}

func (s *txReceiverServer) checkTx(req *pb.TxRequest, ch chan<- txJob) (*pb.Empty, error) {
	if req == nil {
		return nil, errEmpty
	}
	if len(req.Tx) > 0 {
		if req.PeerId == 0 {
			return nil, errZeroId
		}

		var info = TxInfo{
			SenderID: uint16(req.PeerId),
		}

		if ch == nil {
			s.memR.checkTx(req.Tx, info)
		} else {
			ch <- txJob{
				tx:   req.Tx,
				info: info,
			}
		}

		return nil, nil
	} else {
		s.memR.Logger.Error("txReceiverServer.Receive empty tx")
		return nil, errEmpty
	}
}

func (s *txReceiverServer) CheckTxStream(stream pb.MempoolTxReceiver_CheckTxStreamServer) error {
	txCh := make(chan txJob)

	go func(s *txReceiverServer, ch chan txJob) {
		for job := range ch {
			s.memR.checkTx(job.tx, job.info)
		}
	}(s, txCh)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			close(txCh)
			return stream.SendAndClose(empty)
		}
		if err != nil {
			close(txCh)
			return err
		}
		_, err = s.checkTx(req, txCh)
		if err != nil {
			close(txCh)
			return err
		}
	}
}

func (s *txReceiverServer) CheckTxStreamAsync(stream pb.MempoolTxReceiver_CheckTxStreamAsyncServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(empty)
		}
		if err != nil {
			return err
		}
		_, err = s.checkTx(req, s.memR.getTxJobChannel(req.From))
		if err != nil {
			return err
		}
	}
}

var empty = &pb.Empty{}

var errEmpty = fmt.Errorf("empty tx")
var errZeroId = fmt.Errorf("peerId is 0")

type txReceiver struct {
	Server     *txReceiverServer
	Clients    map[uint16]txReceiverClient
	ClientsMtx sync.RWMutex
	Logger     log.Logger

	s    *grpc.Server
	memR *Reactor
}

func newTxReceiver(memR *Reactor) *txReceiver {
	return &txReceiver{
		Server:  newTxReceiverServer(memR),
		Clients: make(map[uint16]txReceiverClient),
		memR:    memR,
	}
}

func (r *txReceiver) SetLogger(l log.Logger) {
	r.Logger = l
	r.Server.SetLogger(l)
}

func (r *txReceiver) AddClient(id uint16, client txReceiverClient) int {
	var receiverCount int
	r.ClientsMtx.Lock()
	r.Clients[id] = client
	receiverCount = len(r.Clients)
	r.ClientsMtx.Unlock()
	return receiverCount
}

func (r *txReceiver) GetClient(id uint16) (client txReceiverClient, ok bool) {
	r.ClientsMtx.RLock()
	client, ok = r.Clients[id]
	r.ClientsMtx.RUnlock()
	return
}

func (r *txReceiver) ReceiveTxReceiverInfo(src p2p.Peer, bz []byte) {
	var info pb.ReceiverInfo
	err := info.Unmarshal(bz)
	if err != nil {
		r.Logger.Error("receiveTxReceiverInfo:unmarshal", "error", err)
		return
	}

	addr := src.SocketAddr().IP.String() + ":" + strconv.FormatInt(info.Port, 10)
	r.Logger.Info("receiveTxReceiverInfo:pre dial", "addr", addr)

	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		//grpc.WithKeepaliveParams(keepalive.ClientParameters{
		//	Time:                30 * time.Second,
		//	Timeout:             10 * time.Second,
		//	PermitWithoutStream: true,
		//}),
	)
	if err != nil {
		r.Logger.Error("receiveTxReceiverInfo:dial", "error", err)
		return
	} else {
		client := pb.NewMempoolTxReceiverClient(conn)
		var receiverCount = r.AddClient(r.memR.ids.GetForPeer(src), txReceiverClient{client, conn, uint16(info.YourId)})
		r.Logger.Info("receiveTxReceiverInfo:success", "peer", src, "yourID", info.YourId, "port", info.Port, "clientCount", receiverCount)
	}
}

func (r *txReceiver) SendTxReceiverInfo(peer p2p.Peer) {
	if !r.Server.Enabled() {
		return
	}
	port := int64(r.Server.Port)
	if r.memR.config.TxReceiverExternalPort != 0 {
		port = int64(r.memR.config.TxReceiverExternalPort)
	}
	var info pb.ReceiverInfo
	info.Port = port
	info.YourId = uint32(r.memR.ids.GetForPeer(peer))
	bz, err := info.Marshal()
	if err != nil {
		r.Logger.Error("sendTxReceiverInfo:marshal", "error", err)
		return
	}

	var retry = 0

	for {
		if !r.memR.IsRunning() || !peer.IsRunning() {
			r.Logger.Error("sendTxReceiverInfo:peer is not running", "peer", peer)
			return
		}
		if retry == 10 {
			r.Logger.Error("sendTxReceiverInfo:try", "times", retry, "peer", peer)
			return
		}
		// make sure the peer is up to date
		_, ok := peer.Get(types.PeerStateKey).(PeerState)
		if !ok {
			// Peer does not have a state yet. We set it in the consensus reactor, but
			// when we add peer in Switch, the order we call reactors#AddPeer is
			// different every time due to us using a map. Sometimes other reactors
			// will be initialized before the consensus reactor. We should wait a few
			// milliseconds and retry.
			time.Sleep(peerCatchupSleepIntervalMS * time.Millisecond)
			continue
		}
		ok = peer.Send(TxReceiverChannel, bz)
		if !ok {
			retry++
			continue
		}
		r.Logger.Info("sendTxReceiverInfo:success", "peer", peer, "peerID", info.YourId, "port", info.Port)
		return
	}
}

func (r *txReceiver) RemovePeer(peer p2p.Peer) {
	r.Logger.Info("pre removePeer", "peer", peer)

	peerID := r.memR.ids.GetForPeer(peer)
	r.ClientsMtx.Lock()
	if c, ok := r.Clients[peerID]; ok {
		var count int
		delete(r.Clients, peerID)
		count = len(r.Clients)
		r.ClientsMtx.Unlock()
		r.Logger.Info("Removing peer from tx receiver", "peer", peer.ID(), "peerID", peerID, "clientCountAfterRemove", count)
		if err := c.Conn.Close(); err != nil {
			r.Logger.Error("Failed to close tx receiver connection", "peer", peer.ID(), "peerID", peerID, "err", err)
		}
	} else {
		r.ClientsMtx.Unlock()
	}
}

func (r *txReceiver) Start(configPort string) {
	configPort = strings.ToLower(configPort)
	if configPort == "off" {
		return
	} else if configPort == "auto" {
		configPort = "0"
	}

	if port, err := strconv.Atoi(configPort); err == nil {
		lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
		if err != nil {
			r.Logger.Error("Failed to start tx receiver:Listen", "err", err)
		} else {
			var options = []grpc.ServerOption{
				//grpc.KeepaliveParams(keepalive.ServerParameters{
				//	Time:    30 * time.Second,
				//	Timeout: 10 * time.Second,
				//}),
				//grpc.ForceServerCodec(newGogoCodec()),
			}
			r.s = grpc.NewServer(options...)
			pb.RegisterMempoolTxReceiverServer(r.s, r.Server)
			r.Server.Port = lis.Addr().(*net.TCPAddr).Port
			r.Logger.Info("Tx receiver listening on port", "port", r.Server.Port)
			atomic.StoreInt64(&r.Server.Started, 1)
			go func() {
				if err := r.s.Serve(lis); err != nil {
					atomic.StoreInt64(&r.Server.Started, 0)
					r.Logger.Error("Failed to start tx receiver:Serve", "err", err)
				}
			}()
		}
	}
}

func (r *txReceiver) CheckTx(peerID uint16, memTx *mempoolTx) bool {
	client, ok := r.GetClient(peerID)
	if ok {
		r.Logger.Debug("CheckTx:Receive", "tx", hex.EncodeToString(memTx.realTx.TxHash()))
		_, err := client.Client.CheckTx(context.Background(), &pb.TxRequest{Tx: memTx.tx, PeerId: uint32(client.ID), From: memTx.from})
		if err != nil {
			r.Logger.Error("CheckTx:Receive", "err", err)
			return false
		}
		return true
	}
	return false
}

func (r *txReceiver) CheckTxByStream(memTx *mempoolTx, peerID uint16, peer p2p.Peer, client *txReceiverClient, streamp *pb.MempoolTxReceiver_CheckTxStreamClient) bool {
	if streamp == nil || client == nil {
		return false
	}

	var err error
	var stream = *streamp

	if stream == nil {
		if client.Client == nil {
			clientV, ok := r.GetClient(peerID)
			if !ok {
				return false
			}
			*client = clientV
		}
		r.Logger.Debug("CheckTxStream:Send", "tx", hex.EncodeToString(memTx.realTx.TxHash()))
		stream, err = client.Client.CheckTxStream(context.Background())
		if err != nil {
			r.Logger.Error("Error CheckTxs", "err", err)
			return false
		}
		*streamp = stream
	}

	err = stream.Send(&pb.TxRequest{Tx: memTx.tx, PeerId: uint32(client.ID), From: memTx.from})
	if err != nil {
		r.Logger.Error("Error Send", "err", err)
		_, err = stream.CloseAndRecv()
		if err != nil {
			r.Logger.Error("Error closing checktxs stream", "peer", peer, "err", err)
		}
		*streamp = nil
		return false
	} else {
		return true
	}
}

func (r *txReceiver) CheckTxAsyncByStream(memTx *mempoolTx, peerID uint16, peer p2p.Peer, client *txReceiverClient, streamp *pb.MempoolTxReceiver_CheckTxStreamAsyncClient) bool {
	if streamp == nil || client == nil {
		return false
	}

	var err error
	var stream = *streamp

	if stream == nil {
		if client.Client == nil {
			clientV, ok := r.GetClient(peerID)
			if !ok {
				return false
			}
			*client = clientV
		}
		r.Logger.Debug("CheckTxStreamAsync:Send", "tx", hex.EncodeToString(memTx.realTx.TxHash()))
		stream, err = client.Client.CheckTxStreamAsync(context.Background())
		if err != nil {
			r.Logger.Error("Error CheckTxs", "err", err)
			return false
		}
		*streamp = stream
	}

	err = stream.Send(&pb.TxRequest{Tx: memTx.tx, PeerId: uint32(client.ID), From: memTx.from})
	if err != nil {
		r.Logger.Error("Error Send", "err", err)
		_, err = stream.CloseAndRecv()
		if err != nil {
			r.Logger.Error("Error closing checktxs stream", "peer", peer, "err", err)
		}
		*streamp = nil
		return false
	} else {
		return true
	}
}

func (r *txReceiver) Stop() {
	if r.s != nil {
		atomic.StoreInt64(&r.Server.Started, 0)
		r.s.Stop()
	}
}

type gogoCodec struct {
	gogocodec.Codec
}

func (_ *gogoCodec) Name() string {
	return proto.Name
}

type marshaler interface {
	MarshalTo(data []byte) (n int, err error)
}

func (c *gogoCodec) Unmarshal(data []byte, v interface{}) error {
	if req, ok := v.(*pb.TxRequest); ok {
		req.Reset()
		return req.Unmarshal(data)
	}
	return c.Codec.Unmarshal(data, v)
}

func (c *gogoCodec) Marshal(v interface{}) ([]byte, error) {
	if req, ok := v.(*pb.TxRequest); ok {
		n := req.Size()
		buf := make([]byte, n)
		_, err := req.MarshalTo(buf)
		if err != nil {
			return nil, err
		}
		return buf, nil
	}

	if m, ok := v.(marshaler); ok {
		n, ok := getSize(v)
		if !ok {
			return gogoproto.Marshal(v.(gogoproto.Message))
		}
		buf := make([]byte, n)
		_, err := m.MarshalTo(buf)
		if err != nil {
			return nil, err
		}
		return buf, nil
	}
	return gogoproto.Marshal(v.(gogoproto.Message))
}

func getSize(v interface{}) (int, bool) {
	if sz, ok := v.(interface {
		Size() (n int)
	}); ok {
		return sz.Size(), true
	} else if sz, ok := v.(interface {
		ProtoSize() (n int)
	}); ok {
		return sz.ProtoSize(), true
	} else {
		return 0, false
	}
}

//func (c *gogoCodec) Marshal(v interface{}) ([]byte, error) {
//	fmt.Println("Marshal")
//	return c.Codec.Marshal(v)
//}

func newGogoCodec() encoding.Codec {
	return &gogoCodec{
		Codec: gogocodec.New(4096),
	}
}

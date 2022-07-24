package mempool

import (
	"bytes"
	"fmt"
	"math"
	"reflect"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"

	lru "github.com/hashicorp/golang-lru"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/libs/clist"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/p2p"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/tendermint/go-amino"
)

const (
	MempoolChannel = byte(0x30)

	aminoOverheadForTxMessage = 8

	peerCatchupSleepIntervalMS = 100 // If peer is behind, sleep this amount

	// UnknownPeerID is the peer ID to use when running CheckTx when there is
	// no peer (e.g. RPC)
	UnknownPeerID uint16 = 0

	maxActiveIDs = math.MaxUint16
)

// Reactor handles mempool tx broadcasting amongst peers.
// It maintains a map from peer ID to counter, to prevent gossiping txs to the
// peers you received it from.
type Reactor struct {
	p2p.BaseReactor
	config            *cfg.MempoolConfig
	mempool           *CListMempool
	ids               *mempoolIDs
	nodeKey           *p2p.NodeKey
	nodeKeyWhitelist  map[string]struct{}
	enableWtx         bool
	enableBatchTx     bool
	isSentryNode      bool
	sentryPartner     string
	sentryPartnerPeer p2p.Peer
	responseChan      chan *TxIndicesMessage
	txMap             *lru.Cache // txIndex -> tx
}

func (memR *Reactor) SetNodeKey(key *p2p.NodeKey) {
	memR.nodeKey = key
}

type mempoolIDs struct {
	mtx       sync.RWMutex
	peerMap   map[p2p.ID]uint16
	nextID    uint16              // assumes that a node will never have over 65536 active peers
	activeIDs map[uint16]struct{} // used to check if a given peerID key is used, the value doesn't matter
}

// Reserve searches for the next unused ID and assigns it to the
// peer.
func (ids *mempoolIDs) ReserveForPeer(peer p2p.Peer) {
	ids.mtx.Lock()
	defer ids.mtx.Unlock()

	curID := ids.nextPeerID()
	ids.peerMap[peer.ID()] = curID
	ids.activeIDs[curID] = struct{}{}
}

// nextPeerID returns the next unused peer ID to use.
// This assumes that ids's mutex is already locked.
func (ids *mempoolIDs) nextPeerID() uint16 {
	if len(ids.activeIDs) == maxActiveIDs {
		panic(fmt.Sprintf("node has maximum %d active IDs and wanted to get one more", maxActiveIDs))
	}

	_, idExists := ids.activeIDs[ids.nextID]
	for idExists {
		ids.nextID++
		_, idExists = ids.activeIDs[ids.nextID]
	}
	curID := ids.nextID
	ids.nextID++
	return curID
}

// Reclaim returns the ID reserved for the peer back to unused pool.
func (ids *mempoolIDs) Reclaim(peer p2p.Peer) {
	ids.mtx.Lock()
	defer ids.mtx.Unlock()

	removedID, ok := ids.peerMap[peer.ID()]
	if ok {
		delete(ids.activeIDs, removedID)
		delete(ids.peerMap, peer.ID())
	}
}

// GetForPeer returns an ID reserved for the peer.
func (ids *mempoolIDs) GetForPeer(peer p2p.Peer) uint16 {
	ids.mtx.RLock()
	defer ids.mtx.RUnlock()

	return ids.peerMap[peer.ID()]
}

func newMempoolIDs() *mempoolIDs {
	return &mempoolIDs{
		peerMap:   make(map[p2p.ID]uint16),
		activeIDs: map[uint16]struct{}{0: {}},
		nextID:    1, // reserve unknownPeerID(0) for mempoolReactor.BroadcastTx
	}
}

// NewReactor returns a new Reactor with the given config and mempool.
func NewReactor(config *cfg.MempoolConfig, mempool *CListMempool) *Reactor {
	cache, err := lru.New(cfg.DynamicConfig.GetMempoolSize())
	if err != nil {
		panic(err)
	}
	memR := &Reactor{
		config:           config,
		mempool:          mempool,
		ids:              newMempoolIDs(),
		nodeKeyWhitelist: make(map[string]struct{}),
		enableWtx:        cfg.DynamicConfig.GetEnableWtx(),
		isSentryNode:     cfg.DynamicConfig.IsSentryNode(),
		sentryPartner:    cfg.DynamicConfig.GetSentryPartner(),
		responseChan:     make(chan *TxIndicesMessage, 10),
		txMap:            cache,
	}

	if memR.sentryPartner != "" && !memR.isSentryNode {
		mempool.SetSentryMempool(memR)
	}

	for _, nodeKey := range config.GetNodeKeyWhitelist() {
		memR.nodeKeyWhitelist[nodeKey] = struct{}{}
	}
	memR.BaseReactor = *p2p.NewBaseReactor("Mempool", memR)
	memR.press()
	return memR
}

func (memR *Reactor) Update(height int64, txs types.Txs) {
	for _, tx := range txs {
		types.SignatureCache().Remove(tx.Hash(height))
	}
}

func (memR *Reactor) ReapTxs(maxBytes, maxGas int64) []types.Tx {
	msg := &FetchMessage{
		MaxBytes: maxBytes,
		MaxGas:   maxGas,
	}
	msgBz := memR.encodeMsg(msg)

	for len(memR.responseChan) != 0 {
		<-memR.responseChan
		memR.Logger.Error("error: responseChan not empty")
		fmt.Println("error: responseChan not empty")
	}

	start := time.Now()
	success := memR.sentryPartnerPeer.Send(MempoolChannel, msgBz)
	if !success {
		memR.Logger.Error("fetch txs from sentry node failed", "peer id", memR.sentryPartnerPeer.ID())
		return nil
	}

	var txs []types.Tx
	ticker := time.After(time.Second)
	for {
		select {
		case msg := <-memR.responseChan:
			txs = make([]types.Tx, 0, len(msg.Indices))
			for _, index := range msg.Indices {
				v, ok := memR.txMap.Get(index)
				if ok {
					txs = append(txs, v.([]byte))
				}
			}
			fmt.Println("fetch tx indices time cost", time.Since(start))
			return txs

		case <-ticker:
			memR.Logger.Error("wait for txs from sentry node timeout", "peer id", memR.sentryPartnerPeer.ID())
			fmt.Println("wait for txs from sentry node timeout", "peer id", memR.sentryPartnerPeer.ID())
			return txs
		}
	}
}

// InitPeer implements Reactor by creating a state for the peer.
func (memR *Reactor) InitPeer(peer p2p.Peer) p2p.Peer {
	memR.ids.ReserveForPeer(peer)
	return peer
}

// SetLogger sets the Logger on the reactor and the underlying mempool.
func (memR *Reactor) SetLogger(l log.Logger) {
	memR.Logger = l
	memR.mempool.SetLogger(l)
}

// OnStart implements p2p.BaseReactor.
func (memR *Reactor) OnStart() error {
	if !memR.config.Broadcast {
		memR.Logger.Info("Tx broadcasting is disabled")
	}
	return nil
}

// GetChannels implements Reactor.
// It returns the list of channels for this reactor.
func (memR *Reactor) GetChannels() []*p2p.ChannelDescriptor {
	return []*p2p.ChannelDescriptor{
		{
			ID:       MempoolChannel,
			Priority: 5,
		},
	}
}

// AddPeer implements Reactor.
// It starts a broadcast routine ensuring all txs are forwarded to the given peer.
func (memR *Reactor) AddPeer(peer p2p.Peer) {
	fmt.Println("debug add peer", peer.ID())
	if string(peer.ID()) == memR.sentryPartner {
		memR.sentryPartnerPeer = peer
	}
	go memR.broadcastTxRoutine(peer)
}

// RemovePeer implements Reactor.
func (memR *Reactor) RemovePeer(peer p2p.Peer, reason interface{}) {
	memR.ids.Reclaim(peer)
	// broadcast routine checks if peer is gone and returns
}

// txMessageDecodePool is a sync.Pool of *TxMessage.
// memR.decodeMsg will call txMessageDeocdePool.Get, and memR.Receive will reset the Msg after use, then call txMessageDeocdePool.Put.
var txMessageDeocdePool = &sync.Pool{
	New: func() interface{} {
		return &TxMessage{}
	},
}

var logParamsPool = &sync.Pool{
	New: func() interface{} {
		return &[6]interface{}{}
	},
}

func (memR *Reactor) logReceive(peer p2p.Peer, chID byte, msg Message) {
	logParams := logParamsPool.Get().(*[6]interface{})

	logParams[0] = "src"
	logParams[1] = peer
	logParams[2] = "chId"
	logParams[3] = chID
	logParams[4] = "msg"
	logParams[5] = msg

	memR.Logger.Debug("Receive", logParams[:]...)

	logParamsPool.Put(logParams)
}

var txIDStringerPool = &sync.Pool{
	New: func() interface{} {
		return &txIDStringer{}
	},
}

func (memR *Reactor) logCheckTxError(tx []byte, height int64, err error) {
	logParams := logParamsPool.Get().(*[6]interface{})
	txStr := txIDStringerPool.Get().(*txIDStringer)
	txStr.tx = tx
	txStr.height = height

	logParams[0] = "tx"
	logParams[1] = txStr
	logParams[2] = "err"
	logParams[3] = err

	memR.Logger.Info("Could not check tx", logParams[:4]...)

	txIDStringerPool.Put(txStr)
	logParamsPool.Put(logParams)
}

// Receive implements Reactor.
// It adds any received transactions to the mempool.
func (memR *Reactor) Receive(chID byte, src p2p.Peer, msgBytes []byte) {
	if memR.mempool.config.Sealed {
		return
	}
	msg, err := memR.decodeMsg(msgBytes)
	if err != nil {
		memR.Logger.Error("Error decoding message", "src", src, "chId", chID, "msg", msg, "err", err, "bytes", msgBytes)
		memR.Switch.StopPeerForError(src, err)
		return
	}
	memR.logReceive(src, chID, msg)

	txInfo := TxInfo{SenderID: memR.ids.GetForPeer(src)}
	if src != nil {
		txInfo.SenderP2PID = src.ID()
	}
	var tx types.Tx

	//TODO: simplify the code
	switch msg := msg.(type) {
	case *FetchMessage:
		indices := memR.mempool.ReapTxIndicesMaxBytesMaxGas(msg.MaxBytes, msg.MaxGas)
		success := src.Send(MempoolChannel, memR.encodeMsg(&TxIndicesMessage{
			Indices: indices,
		}))
		if !success {
			memR.Logger.Error("response txs to validator failed", "peerID", src.ID())
		}

	case *StxMessage:
		if string(src.ID()) == memR.sentryPartner {
			for _, stx := range msg.Stx {
				types.SignatureCache().Add(types.Tx(stx.Tx).Hash(memR.mempool.Height()), stx.From)
				memR.txMap.Add(stx.Index, stx.Tx)
			}
			memR.mempool.mustNotifyTxsAvailable()
		}
		return

	case *TxIndicesMessage:
		if !memR.isSentryNode && string(src.ID()) == memR.sentryPartner {
			select {
			case memR.responseChan <- msg:
			default:
				memR.Logger.Error("debug responseChan full")
			}
		}
		return
	case *TxsMessage:
		if !memR.isSentryNode && memR.sentryPartner != "" {
			return
		}
		for _, tx := range msg.Txs {
			err = memR.mempool.CheckTx(tx, nil, txInfo)
			if err != nil {
				memR.Logger.Info("Could not check tx", "tx", txID(tx, memR.mempool.Height()), "err", err)
			}
		}
		return

	case *TxMessage:
		tx = msg.Tx
		msg.Tx = nil
		txMessageDeocdePool.Put(msg)

	case *WtxMessage:
		tx = msg.Wtx.Payload
		if err := msg.Wtx.verify(memR.nodeKeyWhitelist); err != nil {
			memR.Logger.Error("wtx.verify", "error", err, "txhash",
				common.BytesToHash(types.Tx(msg.Wtx.Payload).Hash(memR.mempool.Height())),
			)
		} else {
			txInfo.wtx = msg.Wtx
			txInfo.checkType = abci.CheckTxType_WrappedCheck
		}
		// broadcasting happens from go routines per peer
	default:
		memR.Logger.Error(fmt.Sprintf("Unknown message type %v", reflect.TypeOf(msg)))
		return
	}

	err = memR.mempool.CheckTx(tx, nil, txInfo)
	if err != nil {
		memR.logCheckTxError(tx, memR.mempool.height, err)
	}
}

// PeerState describes the state of a peer.
type PeerState interface {
	GetHeight() int64
}

// Send new mempool txs to peer.
func (memR *Reactor) broadcastTxRoutine(peer p2p.Peer) {
	if !memR.config.Broadcast {
		return
	}

	peerID := memR.ids.GetForPeer(peer)
	var next *clist.CElement
	for {
		// In case of both next.NextWaitChan() and peer.Quit() are variable at the same time
		if !memR.IsRunning() || !peer.IsRunning() {
			return
		}
		// This happens because the CElement we were looking at got garbage
		// collected (removed). That is, .NextWait() returned nil. Go ahead and
		// start from the beginning.
		if next == nil {
			select {
			case <-memR.mempool.TxsWaitChan(): // Wait until a tx is available
				if next = memR.mempool.BroadcastTxsFront(); next == nil {
					continue
				}
			case <-peer.Quit():
				return
			case <-memR.Quit():
				return
			}
		}

		// make sure the peer is up to date
		peerState, ok := peer.Get(types.PeerStateKey).(PeerState)
		if !ok {
			// Peer does not have a state yet. We set it in the consensus reactor, but
			// when we add peer in Switch, the order we call reactors#AddPeer is
			// different every time due to us using a map. Sometimes other reactors
			// will be initialized before the consensus reactor. We should wait a few
			// milliseconds and retry.
			time.Sleep(peerCatchupSleepIntervalMS * time.Millisecond)
			continue
		}
		memTx := next.Value.(*mempoolTx)
		if peerState.GetHeight() < memTx.Height()-1 { // Allow for a lag of 1 block
			time.Sleep(peerCatchupSleepIntervalMS * time.Millisecond)
			continue
		}

		if cfg.DynamicConfig.GetEnableBatchTx() {
			var msg Message
			if memR.isSentryNode && string(peer.ID()) == memR.sentryPartner {
				var sentryTxs []*SentryTx
				sentryTxs, next = memR.mempool.sentryTxs(next)
				msg = StxMessage{
					Stx: sentryTxs,
				}

			} else {
				var txs []types.Tx
				txs, next = memR.mempool.getTxs(next, peerID)
				msg = TxsMessage{
					Txs: txs,
				}
			}

			msgBz := memR.encodeMsg(msg)
			success := peer.Send(MempoolChannel, msgBz)
			if !success {
				time.Sleep(peerCatchupSleepIntervalMS * time.Millisecond)
				continue
			}

		} else {
			// ensure peer hasn't already sent us this tx
			if _, ok := memTx.senders.Load(peerID); !ok {
				var getFromPool bool
				// send memTx
				var msg Message
				if memTx.nodeKey != nil && memTx.signature != nil {
					msg = &WtxMessage{
						Wtx: &WrappedTx{
							Payload:   memTx.tx,
							From:      memTx.from,
							Signature: memTx.signature,
							NodeKey:   memTx.nodeKey,
						},
					}
				} else if memR.enableWtx {
					if wtx, err := memR.wrapTx(memTx.tx, memTx.from); err == nil {
						msg = &WtxMessage{
							Wtx: wtx,
						}
					}
				} else {
					txMsg := txMessageDeocdePool.Get().(*TxMessage)
					txMsg.Tx = memTx.tx
					msg = txMsg
					getFromPool = true
				}

				msgBz := memR.encodeMsg(msg)
				if getFromPool {
					getFromPool = false
					txMessageDeocdePool.Put(msg)
				}

				success := peer.Send(MempoolChannel, msgBz)
				if !success {
					time.Sleep(peerCatchupSleepIntervalMS * time.Millisecond)
					continue
				}
			}
		}

		select {
		case <-next.NextWaitChan():
			// see the start of the for loop for nil check
			next = next.Next()
		case <-peer.Quit():
			return
		case <-memR.Quit():
			return
		}
	}
}

//-----------------------------------------------------------------------------
// Messages

// Message is a message sent or received by the Reactor.
type Message interface{}

func RegisterMessages(cdc *amino.Codec) {
	cdc.RegisterInterface((*Message)(nil), nil)
	cdc.RegisterConcrete(&TxMessage{}, "tendermint/mempool/TxMessage", nil)
	cdc.RegisterConcrete(&WtxMessage{}, "tendermint/mempool/WtxMessage", nil)
	cdc.RegisterConcrete(&StxMessage{}, "tendermint/mempool/StxMessage", nil)
	cdc.RegisterConcrete(&TxsMessage{}, "tendermint/mempool/TxsMessage", nil)
	cdc.RegisterConcrete(&FetchMessage{}, "tendermint/mempool/FetchMessage", nil)
	cdc.RegisterConcrete(&TxIndicesMessage{}, "tendermint/mempool/TxIndicesMessage", nil)

	cdc.RegisterConcreteMarshaller("tendermint/mempool/TxMessage", func(codec *amino.Codec, i interface{}) ([]byte, error) {
		txmp, ok := i.(*TxMessage)
		if ok {
			return txmp.MarshalToAmino(codec)
		}
		txm, ok := i.(TxMessage)
		if ok {
			return txm.MarshalToAmino(codec)
		}
		return nil, fmt.Errorf("%T is not a TxMessage", i)
	})
	cdc.RegisterConcreteUnmarshaller("tendermint/mempool/TxMessage", func(cdc *amino.Codec, bz []byte) (interface{}, int, error) {
		m := &TxMessage{}
		err := m.UnmarshalFromAmino(cdc, bz)
		if err != nil {
			return nil, 0, err
		}
		return m, len(bz), nil
	})
}

// decodeMsg decodes the bz bytes into a Message,
// if err is nil and Message is a TxMessage, you must put Message to txMessageDeocdePool after use.
func (memR *Reactor) decodeMsg(bz []byte) (Message, error) {
	maxMsgSize := calcMaxMsgSize(memR.config.MaxTxBytes)
	l := len(bz)
	if l > maxMsgSize {
		return nil, ErrTxTooLarge{maxMsgSize, l}
	}

	tp := getTxMessageAminoTypePrefix()
	if l >= len(tp) && bytes.Equal(bz[:len(tp)], tp) {
		txmsg := txMessageDeocdePool.Get().(*TxMessage)
		err := txmsg.UnmarshalFromAmino(cdc, bz[len(tp):])
		if err == nil {
			return txmsg, nil
		}
		txmsg.Tx = nil
		txMessageDeocdePool.Put(txmsg)
	}
	var msg Message
	err := cdc.UnmarshalBinaryBare(bz, &msg)
	return msg, err
}

func (memR *Reactor) encodeMsg(msg Message) []byte {
	var ok bool
	var txmp *TxMessage
	var txm TxMessage
	if txmp, ok = msg.(*TxMessage); !ok {
		txmp = nil
		if txm, ok = msg.(TxMessage); ok {
			txmp = &txm
		}
	}
	if txmp != nil {
		buf := &bytes.Buffer{}
		tp := getTxMessageAminoTypePrefix()
		buf.Grow(len(tp) + txmp.AminoSize(cdc))
		// we manually assemble the encoded bytes for performance
		buf.Write(tp)
		err := txmp.MarshalAminoTo(cdc, buf)
		if err == nil {
			return buf.Bytes()
		}
	}
	return cdc.MustMarshalBinaryBare(msg)
}

//-------------------------------------

// TxMessage is a Message containing a transaction.
type TxMessage struct {
	Tx types.Tx
}

func (m TxMessage) AminoSize(_ *amino.Codec) int {
	size := 0
	if len(m.Tx) > 0 {
		size += 1 + amino.ByteSliceSize(m.Tx)
	}
	return size
}

func (m TxMessage) MarshalToAmino(cdc *amino.Codec) ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Grow(m.AminoSize(cdc))
	err := m.MarshalAminoTo(cdc, buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m TxMessage) MarshalAminoTo(_ *amino.Codec, buf *bytes.Buffer) error {
	if len(m.Tx) != 0 {
		const pbKey = byte(1<<3 | amino.Typ3_ByteLength)
		err := amino.EncodeByteSliceWithKeyToBuffer(buf, m.Tx, pbKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *TxMessage) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
	if len(data) == 0 {
		return nil
	}

	if data[0] != 1<<3|byte(amino.Typ3_ByteLength) {
		return fmt.Errorf("error pb type")
	}

	data = data[1:]
	dataLen, n, err := amino.DecodeUvarint(data)
	if err != nil {
		return err
	}
	data = data[n:]
	if len(data) != int(dataLen) {
		return fmt.Errorf("invalid datalen")
	}

	m.Tx = nil
	if dataLen > 0 {
		m.Tx = make([]byte, dataLen)
		copy(m.Tx, data)
	}

	return nil
}

// String returns a string representation of the TxMessage.
func (m *TxMessage) String() string {
	return fmt.Sprintf("[TxMessage %v]", m.Tx)
}

// calcMaxMsgSize returns the max size of TxMessage
// account for amino overhead of TxMessage
func calcMaxMsgSize(maxTxSize int) int {
	return maxTxSize + aminoOverheadForTxMessage
}

// FetchMessage is a Message of info to fetch txs from sentry mempool.
type FetchMessage struct {
	MaxBytes int64
	MaxGas   int64
}

// String returns a string representation of the FetchMessage.
func (m *FetchMessage) String() string {
	return fmt.Sprintf("[FetchMessage %d, %d]", m.MaxBytes, m.MaxGas)
}

// TxIndicesMessage is a Message containing transaction indices.
type TxIndicesMessage struct {
	Indices []uint32
}

// String returns a string representation of the TxIndicesMessage.
func (m *TxIndicesMessage) String() string {
	return fmt.Sprintf("[TxIndicesMessage %v]", m.Indices)
}

// TxsMessage is a Message containing transactions.
type TxsMessage struct {
	Txs     []types.Tx
	Indices []uint32
}

// String returns a string representation of the TxsMessage.
func (m *TxsMessage) String() string {
	return fmt.Sprintf("[TxsMessage %v]", m.Txs)
}

// StxMessage is a Message containing a Stx.
type StxMessage struct {
	Stx []*SentryTx
}

// String returns a string representation of the StxMessage.
func (m *StxMessage) String() string {
	return fmt.Sprintf("[StxMessage %v]", m.Stx)
}

// SentryTx is a result of verify signature from sentry node.
type SentryTx struct {
	Tx    []byte
	From  string
	Index uint32
}

// WtxMessage is a Message containing a transaction.
type WtxMessage struct {
	Wtx *WrappedTx
}

// String returns a string representation of the WtxMessage.
func (m *WtxMessage) String() string {
	return fmt.Sprintf("[WtxMessage %v]", m.Wtx)
}

type WrappedTx struct {
	Payload   []byte `json:"payload"`   // std tx or evm tx
	From      string `json:"from"`      // from address of evm tx or ""
	Signature []byte `json:"signature"` // signature for payload
	NodeKey   []byte `json:"nodeKey"`   // pub key of the node who signs the tx
}

func (wtx *WrappedTx) GetPayload() []byte {
	if wtx != nil {
		return wtx.Payload
	}
	return nil
}

func (wtx *WrappedTx) GetSignature() []byte {
	if wtx != nil {
		return wtx.Signature
	}
	return nil
}

func (wtx *WrappedTx) GetNodeKey() []byte {
	if wtx != nil {
		return wtx.NodeKey
	}
	return nil
}

func (wtx *WrappedTx) GetFrom() string {
	if wtx != nil {
		return wtx.From
	}
	return ""
}

func (w *WrappedTx) verify(whitelist map[string]struct{}) error {
	pub := p2p.BytesToPubKey(w.NodeKey)
	if _, ok := whitelist[string(p2p.PubKeyToID(pub))]; !ok {
		return fmt.Errorf("node key [%s] not in whitelist", p2p.PubKeyToID(pub))
	}
	if !pub.VerifyBytes(append(w.Payload, w.From...), w.Signature) {
		return fmt.Errorf("invalid signature of wtx")
	}
	return nil
}

func (memR *Reactor) wrapTx(tx types.Tx, from string) (*WrappedTx, error) {
	wtx := &WrappedTx{
		Payload: tx,
		From:    from,
		NodeKey: memR.nodeKey.PubKey().Bytes(),
	}
	sig, err := memR.nodeKey.PrivKey.Sign(append(wtx.Payload, from...))
	if err != nil {
		return nil, err
	}
	wtx.Signature = sig
	return wtx, nil
}

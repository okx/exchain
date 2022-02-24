package mempool

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math"
	"reflect"
	"sync"
	"time"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/libs/clist"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/p2p"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/viper"
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
	config           *cfg.MempoolConfig
	mempool          *CListMempool
	ids              *mempoolIDs
	nodeKey          *p2p.NodeKey
	nodeKeyWhitelist map[string]struct{}
	enableWtx        bool
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
	memR := &Reactor{
		config:           config,
		mempool:          mempool,
		ids:              newMempoolIDs(),
		nodeKeyWhitelist: make(map[string]struct{}),
		enableWtx:        viper.GetBool(abci.FlagEnableWrappedTx),
	}
	for _, nodeKey := range config.GetNodeKeyWhitelist() {
		memR.nodeKeyWhitelist[nodeKey] = struct{}{}
	}
	memR.BaseReactor = *p2p.NewBaseReactor("Mempool", memR)
	memR.press()
	return memR
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
	go memR.broadcastTxRoutine(peer)
}

// RemovePeer implements Reactor.
func (memR *Reactor) RemovePeer(peer p2p.Peer, reason interface{}) {
	memR.ids.Reclaim(peer)
	// broadcast routine checks if peer is gone and returns
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
	memR.Logger.Debug("Receive", "src", src, "chId", chID, "msg", msg)

	txInfo := TxInfo{SenderID: memR.ids.GetForPeer(src)}
	if src != nil {
		txInfo.SenderP2PID = src.ID()
	}
	var tx types.Tx

	switch msg := msg.(type) {
	case *TxMessage:
		tx = msg.Tx

	case *WtxMessage:
		tx = msg.Wtx.Payload
		if err := msg.Wtx.verify(memR.nodeKeyWhitelist); err != nil {
			memR.Logger.Error("wtx.verify", "error", err, "txhash",
				common.BytesToHash(types.Tx(msg.Wtx.Payload).Hash(memR.mempool.height)),
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
		memR.Logger.Info("Could not check tx", "tx", txID(tx, memR.mempool.height), "err", err)
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

		memTx := next.Value.(*mempoolTx)

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
		if peerState.GetHeight() < memTx.Height()-1 { // Allow for a lag of 1 block
			time.Sleep(peerCatchupSleepIntervalMS * time.Millisecond)
			continue
		}

		// ensure peer hasn't already sent us this tx
		if _, ok := memTx.senders.Load(peerID); !ok {
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
				msg = &TxMessage{
					Tx: memTx.tx,
				}
			}

			success := peer.Send(MempoolChannel, cdc.MustMarshalBinaryBare(msg))
			if !success {
				time.Sleep(peerCatchupSleepIntervalMS * time.Millisecond)
				continue
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
}

func (memR *Reactor) decodeMsg(bz []byte) (msg Message, err error) {
	maxMsgSize := calcMaxMsgSize(memR.config.MaxTxBytes)
	if l := len(bz); l > maxMsgSize {
		return msg, ErrTxTooLarge{maxMsgSize, l}
	}
	err = cdc.UnmarshalBinaryBare(bz, &msg)
	return
}

//-------------------------------------

// TxMessage is a Message containing a transaction.
type TxMessage struct {
	Tx types.Tx
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

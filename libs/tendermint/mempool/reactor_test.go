package mempool

import (
	"math/rand"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/okex/exchain/libs/tendermint/crypto/ed25519"

	"github.com/fortytw2/leaktest"
	"github.com/go-kit/kit/log/term"
	"github.com/okex/exchain/libs/tendermint/abci/example/kvstore"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/p2p"
	"github.com/okex/exchain/libs/tendermint/p2p/mock"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

type peerState struct {
	height int64
}

func (ps peerState) GetHeight() int64 {
	return ps.height
}

// mempoolLogger is a TestingLogger which uses a different
// color for each validator ("validator" key must exist).
func mempoolLogger() log.Logger {
	return log.TestingLoggerWithColorFn(func(keyvals ...interface{}) term.FgBgColor {
		for i := 0; i < len(keyvals)-1; i += 2 {
			if keyvals[i] == "validator" {
				return term.FgBgColor{Fg: term.Color(uint8(keyvals[i+1].(int) + 1))}
			}
		}
		return term.FgBgColor{}
	})
}

// connect N mempool reactors through N switches
func makeAndConnectReactors(config *cfg.Config, n int) []*Reactor {
	reactors := make([]*Reactor, n)
	logger := mempoolLogger()
	for i := 0; i < n; i++ {
		app := kvstore.NewApplication()
		cc := proxy.NewLocalClientCreator(app)
		mempool, cleanup := newMempoolWithApp(cc)
		defer cleanup()

		reactors[i] = NewReactor(config.Mempool, mempool) // so we dont start the consensus states
		reactors[i].SetLogger(logger.With("validator", i))
	}

	p2p.MakeConnectedSwitches(config.P2P, n, func(i int, s *p2p.Switch) *p2p.Switch {
		s.AddReactor("MEMPOOL", reactors[i])
		return s

	}, p2p.Connect2Switches)
	return reactors
}

func waitForTxsOnReactors(t *testing.T, txs types.Txs, reactors []*Reactor) {
	// wait for the txs in all mempools
	wg := new(sync.WaitGroup)
	for i, reactor := range reactors {
		wg.Add(1)
		go func(r *Reactor, reactorIndex int) {
			defer wg.Done()
			waitForTxsOnReactor(t, txs, r, reactorIndex)
		}(reactor, i)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	timer := time.After(Timeout)
	select {
	case <-timer:
		t.Fatal("Timed out waiting for txs")
	case <-done:
	}
}

func waitForTxsOnReactor(t *testing.T, txs types.Txs, reactor *Reactor, reactorIndex int) {
	mempool := reactor.mempool
	for mempool.Size() < len(txs) {
		time.Sleep(time.Millisecond * 50)
	}

	reapedTxs := mempool.ReapMaxTxs(len(txs))
	for i, tx := range txs {
		assert.Equalf(t, tx, reapedTxs[i],
			"txs at index %d on reactor %d don't match: %v vs %v", i, reactorIndex, tx, reapedTxs[i])
	}
}

// ensure no txs on reactor after some timeout
func ensureNoTxs(t *testing.T, reactor *Reactor, timeout time.Duration) {
	time.Sleep(timeout) // wait for the txs in all mempools
	assert.Zero(t, reactor.mempool.Size())
}

const (
	NumTxs  = 1000
	Timeout = 120 * time.Second // ridiculously high because CircleCI is slow
)

//TODO fix random failure case
func testReactorBroadcastTxMessage(t *testing.T) {
	config := cfg.TestConfig()
	const N = 4
	reactors := makeAndConnectReactors(config, N)
	defer func() {
		for _, r := range reactors {
			r.Stop()
		}
	}()
	for _, r := range reactors {
		for _, peer := range r.Switch.Peers().List() {
			peer.Set(types.PeerStateKey, peerState{1})
		}
	}

	// send a bunch of txs to the first reactor's mempool
	// and wait for them all to be received in the others
	txs := checkTxs(t, reactors[0].mempool, NumTxs, UnknownPeerID)
	waitForTxsOnReactors(t, txs, reactors)
}

func TestReactorNoBroadcastToSender(t *testing.T) {
	config := cfg.TestConfig()
	const N = 2
	reactors := makeAndConnectReactors(config, N)
	defer func() {
		for _, r := range reactors {
			r.Stop()
		}
	}()

	// send a bunch of txs to the first reactor's mempool, claiming it came from peer
	// ensure peer gets no txs
	checkTxs(t, reactors[0].mempool, NumTxs, 1)
	ensureNoTxs(t, reactors[1], 100*time.Millisecond)
}

func TestBroadcastTxForPeerStopsWhenPeerStops(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	config := cfg.TestConfig()
	const N = 2
	reactors := makeAndConnectReactors(config, N)
	defer func() {
		for _, r := range reactors {
			r.Stop()
		}
	}()

	// stop peer
	sw := reactors[1].Switch
	sw.StopPeerForError(sw.Peers().List()[0], errors.New("some reason"))

	// check that we are not leaking any go-routines
	// i.e. broadcastTxRoutine finishes when peer is stopped
	leaktest.CheckTimeout(t, 10*time.Second)()
}

func TestBroadcastTxForPeerStopsWhenReactorStops(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	config := cfg.TestConfig()
	const N = 2
	reactors := makeAndConnectReactors(config, N)

	// stop reactors
	for _, r := range reactors {
		r.Stop()
	}

	// check that we are not leaking any go-routines
	// i.e. broadcastTxRoutine finishes when reactor is stopped
	leaktest.CheckTimeout(t, 10*time.Second)()
}

func TestMempoolIDsBasic(t *testing.T) {
	ids := newMempoolIDs()

	peer := mock.NewPeer(net.IP{127, 0, 0, 1})

	ids.ReserveForPeer(peer)
	assert.EqualValues(t, 1, ids.GetForPeer(peer))
	ids.Reclaim(peer)

	ids.ReserveForPeer(peer)
	assert.EqualValues(t, 2, ids.GetForPeer(peer))
	ids.Reclaim(peer)
}

func TestMempoolIDsPanicsIfNodeRequestsOvermaxActiveIDs(t *testing.T) {
	if testing.Short() {
		return
	}

	// 0 is already reserved for UnknownPeerID
	ids := newMempoolIDs()

	for i := 0; i < maxActiveIDs-1; i++ {
		peer := mock.NewPeer(net.IP{127, 0, 0, 1})
		ids.ReserveForPeer(peer)
	}

	assert.Panics(t, func() {
		peer := mock.NewPeer(net.IP{127, 0, 0, 1})
		ids.ReserveForPeer(peer)
	})
}

func TestDontExhaustMaxActiveIDs(t *testing.T) {
	config := cfg.TestConfig()
	const N = 1
	reactors := makeAndConnectReactors(config, N)
	defer func() {
		for _, r := range reactors {
			r.Stop()
		}
	}()
	reactor := reactors[0]

	for i := 0; i < maxActiveIDs+1; i++ {
		peer := mock.NewPeer(nil)
		reactor.Receive(MempoolChannel, peer, []byte{0x1, 0x2, 0x3})
		reactor.AddPeer(peer)
	}
}

func TestVerifyWtx(t *testing.T) {
	nodeKey := &p2p.NodeKey{
		PrivKey: ed25519.GenPrivKey(),
	}
	memR := &Reactor{
		nodeKey: nodeKey,
	}

	wtx, err := memR.wrapTx([]byte("test-tx"), "test-from")
	assert.Nil(t, err)

	nodeKeyWhitelist := make(map[string]struct{})
	err = wtx.verify(nodeKeyWhitelist)
	assert.NotNil(t, err)

	nodeKeyWhitelist[string(p2p.PubKeyToID(nodeKey.PubKey()))] = struct{}{}
	err = wtx.verify(nodeKeyWhitelist)
	assert.Nil(t, err)
}

func TestTxMessageAmino(t *testing.T) {
	testcases := []TxMessage{
		{},
		{[]byte{}},
		{[]byte{1, 2, 3, 4, 5, 6, 7}},
		{[]byte{}},
	}

	var typePrefix = make([]byte, 8)
	tpLen, err := cdc.GetTypePrefix(TxMessage{}, typePrefix)
	require.NoError(t, err)
	typePrefix = typePrefix[:tpLen]
	reactor := Reactor{
		config: &cfg.MempoolConfig{
			MaxTxBytes: 1024 * 1024,
		},
	}

	for _, tx := range testcases {
		var m Message
		m = tx
		expectBz, err := cdc.MarshalBinaryBare(m)
		require.NoError(t, err)
		actualBz, err := tx.MarshalToAmino(cdc)
		require.NoError(t, err)

		require.Equal(t, expectBz, append(typePrefix, actualBz...))
		require.Equal(t, len(expectBz), tpLen+tx.AminoSize(cdc))

		actualBz, err = cdc.MarshalBinaryBareWithRegisteredMarshaller(tx)
		require.NoError(t, err)

		require.Equal(t, expectBz, actualBz)
		require.Equal(t, cdc.MustMarshalBinaryBare(m), reactor.encodeMsg(&tx))
		require.Equal(t, cdc.MustMarshalBinaryBare(m), reactor.encodeMsg(tx))

		var expectValue Message
		err = cdc.UnmarshalBinaryBare(expectBz, &expectValue)
		require.NoError(t, err)
		var actualValue Message
		actualValue, err = cdc.UnmarshalBinaryBareWithRegisteredUnmarshaller(expectBz, &actualValue)
		require.Equal(t, expectValue, actualValue)

		actualValue, err = reactor.decodeMsg(expectBz)
		require.NoError(t, err)
		require.Equal(t, expectValue, actualValue)
		actualValue.(*TxMessage).Tx = nil
		txMessageDeocdePool.Put(actualValue)
	}

	// special case
	{
		var bz = []byte{1<<3 | 2, 0}
		bz = append(typePrefix, bz...)
		var expectValue Message
		err = cdc.UnmarshalBinaryBare(bz, &expectValue)
		require.NoError(t, err)
		var actualValue Message
		actualValue, err = cdc.UnmarshalBinaryBareWithRegisteredUnmarshaller(bz, &actualValue)
		require.NoError(t, err)
		require.Equal(t, expectValue, actualValue)

		actualValue, err = reactor.decodeMsg(bz)
		require.NoError(t, err)
		require.Equal(t, expectValue, actualValue)
	}
}

func BenchmarkTxMessageAminoMarshal(b *testing.B) {
	var bz = make([]byte, 256)
	rand.Read(bz)
	reactor := &Reactor{}
	var msg Message
	b.ResetTimer()

	b.Run("amino", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			msg = TxMessage{bz}
			_, err := cdc.MarshalBinaryBare(&msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("marshaller", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			msg = &TxMessage{bz}
			_, err := cdc.MarshalBinaryBareWithRegisteredMarshaller(msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("encodeMsgOld", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			msg = &TxMessage{bz}
			reactor.encodeMsg(msg)
		}
	})
	b.Run("encodeMsg", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			txm := txMessageDeocdePool.Get().(*TxMessage)
			txm.Tx = bz
			msg = txm
			reactor.encodeMsg(msg)
			txMessageDeocdePool.Put(txm)
		}
	})
}

func decodeMsgOld(memR *Reactor, bz []byte) (msg Message, err error) {
	maxMsgSize := calcMaxMsgSize(memR.config.MaxTxBytes)
	if l := len(bz); l > maxMsgSize {
		return msg, ErrTxTooLarge{maxMsgSize, l}
	}
	err = cdc.UnmarshalBinaryBare(bz, &msg)
	return
}

func BenchmarkTxMessageUnmarshal(b *testing.B) {
	txMsg := TxMessage{
		Tx: make([]byte, 512),
	}
	rand.Read(txMsg.Tx)
	bz := cdc.MustMarshalBinaryBare(&txMsg)

	//msg := conn.PacketMsg{
	//	ChannelID: MempoolChannel,
	//	Bytes:     bz,
	//}
	// msgBz := cdc.MustMarshalBinaryBare(&msg)

	//hashMap := make(map[string]struct{})
	var h []byte

	reactor := &Reactor{
		config: &cfg.MempoolConfig{
			MaxTxBytes: 1024 * 1024,
		},
	}

	var msg Message
	var err error

	b.ResetTimer()

	b.Run("decode", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			//var m conn.PacketMsg
			//err := m.UnmarshalFromAmino(nil, msgBz)
			//if err != nil {
			//	b.Fatal(err)
			//}
			msg, err = reactor.decodeMsg(bz)
			if err != nil {
				b.Fatal(err)
			}
			msg.(*TxMessage).Tx = nil
			txMessageDeocdePool.Put(msg)
		}
	})
	b.Run("decode-old", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			//var m conn.PacketMsg
			//err := m.UnmarshalFromAmino(nil, msgBz)
			//if err != nil {
			//	b.Fatal(err)
			//}
			msg, err = decodeMsgOld(reactor, bz)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	b.Run("amino", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var m TxMessage
			err := m.UnmarshalFromAmino(cdc, bz[4:])
			msg = &m
			if err != nil {
				b.Fatal(err)
			}
		}
	})
	//b.Run("hash", func(b *testing.B) {
	//	b.ReportAllocs()
	//	for i := 0; i < b.N; i++ {
	//		var m conn.PacketMsg
	//		err := m.UnmarshalFromAmino(nil, msgBz)
	//		if err != nil {
	//			b.Fatal(err)
	//		}
	//		_ = crypto.Sha256(bz)
	//	}
	//})
	_ = h
	_ = msg
}

func BenchmarkReactorLogReceive(b *testing.B) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "benchmark")
	var options []log.Option
	options = append(options, log.AllowInfoWith("module", "benchmark"))
	logger = log.NewFilter(logger, options...)

	memR := &Reactor{}
	memR.Logger = logger

	chID := byte(10)
	var msg Message = &TxMessage{Tx: make([]byte, 512)}
	var src p2p.Peer

	b.Run("pool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			memR.logReceive(src, chID, msg)
		}
	})

	b.Run("logger", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			memR.Logger.Debug("Receive", "src", src, "chId", chID, "msg", msg)
		}
	})
}

func BenchmarkReactorLogCheckTxError(b *testing.B) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "benchmark")
	var options []log.Option
	options = append(options, log.AllowErrorWith("module", "benchmark"))
	logger = log.NewFilter(logger, options...)

	memR := &Reactor{}
	memR.Logger = logger
	memR.mempool = &CListMempool{height: 123456}

	var msg Message = &TxMessage{Tx: make([]byte, 512)}
	tx := msg.(*TxMessage).Tx
	err := errors.New("error")

	b.Run("pool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			memR.logCheckTxError(tx, memR.mempool.height, err)
		}
	})

	b.Run("logger", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			memR.Logger.Info("Could not check tx", "tx", txIDStringer{tx, memR.mempool.height}, "err", err)
		}
	})
}

func TestDecodeMsg(t *testing.T) {
	memR := NewReactor(cfg.TestMempoolConfig(), nil)
	var originMsg, decMsg Message

	originMsg = &FetchMessage{
		MaxBytes: 2 << 20,
		MaxGas:   -1,
	}
	decMsg, _ = memR.decodeMsg(memR.encodeMsg(originMsg))
	require.Equal(t, originMsg, decMsg)

	originMsg = &StxMessage{
		Stx: []*SentryTx{{Tx: []byte("test"), From: "testFrom"}},
	}
	decMsg, _ = memR.decodeMsg(memR.encodeMsg(originMsg))
	require.Equal(t, originMsg, decMsg)

	originMsg = &TxsMessage{
		Txs: []types.Tx{[]byte("testTx")},
	}
	decMsg, _ = memR.decodeMsg(memR.encodeMsg(originMsg))
	require.Equal(t, originMsg, decMsg)

	originMsg = &TxIndicesMessage{}
	decMsg, _ = memR.decodeMsg(memR.encodeMsg(originMsg))
	require.Equal(t, originMsg, decMsg)

	originMsg = &TxIndicesMessage{
		Indices: []uint32{666, 888},
	}
	decMsg, _ = memR.decodeMsg(memR.encodeMsg(originMsg))
	require.Equal(t, originMsg, decMsg)

}

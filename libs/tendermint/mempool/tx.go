package mempool

import (
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/types"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

// mempoolTx is a transaction that successfully ran
type mempoolTx struct {
	height      int64    // height that this tx had been validated in
	gasWanted   int64    // amount of gas this tx states it will require
	tx          types.Tx //
	realTx      abci.TxEssentials
	nodeKey     []byte
	signature   []byte
	from        string
	senderNonce uint64

	// ids of peers who've sent us this tx (as a map for quick lookups).
	// senders: PeerID -> bool
	senders sync.Map

	// timestamp is the time at which the node first received the transaction from
	// a peer. It is used as a second dimension is prioritizing transactions when
	// two transactions have the same priority.
	timestamp time.Time
}

// Height returns the height for this transaction
func (memTx *mempoolTx) Height() int64 {
	return atomic.LoadInt64(&memTx.height)
}


//--------------------------------------------------------------------------------
type ExTxInfo struct {
	Sender      string   `json:"sender"`
	SenderNonce uint64   `json:"sender_nonce"`
	GasPrice    *big.Int `json:"gas_price"`
	Nonce       uint64   `json:"nonce"`
}
//--------------------------------------------------------------------------------
// mempoolTxList implements a thread-safe list of *mempoolTx objects that can be
// used to build generic transaction indexes in the mempool. It accepts a
// comparator function, less(a, b *mempoolTx) bool, that compares two mempoolTx
// references which is used during Insert in order to determine sorted order. If
// less returns true, a <= b.
type MempoolTxList struct {
	mtx  sync.RWMutex
	txs  []*mempoolTx
	less func(*mempoolTx, *mempoolTx) bool
}

func NewMempoolTxList(less func(*mempoolTx, *mempoolTx) bool) *MempoolTxList {
	return &MempoolTxList{
		txs:  make([]*mempoolTx, 0),
		less: less,
	}
}

// Size returns the number of mempoolTx objects in the list.
func (mtl *MempoolTxList) Size() int {
	mtl.mtx.RLock()
	defer mtl.mtx.RUnlock()

	return len(mtl.txs)
}

// Reset resets the list of transactions to an empty list.
func (mtl *MempoolTxList) Reset() {
	mtl.mtx.Lock()
	defer mtl.mtx.Unlock()

	mtl.txs = make([]*mempoolTx, 0)
}

// Insert inserts a mempoolTx reference into the sorted list based on the list's
// comparator function.
func (mtl *MempoolTxList) Insert(mtx *mempoolTx) {
	mtl.mtx.Lock()
	defer mtl.mtx.Unlock()

	i := sort.Search(len(mtl.txs), func(i int) bool {
		return mtl.less(mtl.txs[i], mtx)
	})

	if i == len(mtl.txs) {
		// insert at the end
		mtl.txs = append(mtl.txs, mtx)
		return
	}

	// Make space for the inserted element by shifting values at the insertion
	// index up one index.
	//
	// NOTE: The call to append does not allocate memory when cap(wtl.txs) > len(wtl.txs).
	mtl.txs = append(mtl.txs[:i+1], mtl.txs[i:]...)
	mtl.txs[i] = mtx
}

// Remove attempts to remove a mempoolTx from the sorted list.
func (mtl *MempoolTxList) Remove(mtx *mempoolTx) {
	mtl.mtx.Lock()
	defer mtl.mtx.Unlock()

	i := sort.Search(len(mtl.txs), func(i int) bool {
		return mtl.less(mtl.txs[i], mtx)
	})

	// Since the list is sorted, we evaluate all elements starting at i. Note, if
	// the element does not exist, we may potentially evaluate the entire remainder
	// of the list. However, a caller should not be expected to call Remove with a
	// non-existing element.
	for i < len(mtl.txs) {
		if mtl.txs[i] == mtx {
			mtl.txs = append(mtl.txs[:i], mtl.txs[i+1:]...)
			return
		}

		i++
	}
}
package watcher

import (
	"math/big"

	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	types2 "github.com/okex/okexchain/x/evm/types"
	"github.com/tendermint/tendermint/abci/types"
)

type Watcher struct {
	store         *WatchStore
	height        uint64
	blockHash     common.Hash
	header        types.Header
	batch         []WatchMessage
	cumulativeGas map[uint64]uint64
	gasUsed       uint64
	blockTxs      []common.Hash
	sw            bool
}

func NewWatcher() *Watcher {
	return &Watcher{store: InstanceOfWatchStore(), sw: viper.GetBool(ViperFlagFastQuery)}
}

func (w Watcher) enabled() bool {
	return w.sw
}

func (w *Watcher) Enable(sw bool) {
	w.sw = sw
}

func (w *Watcher) NewHeight(height uint64, blockHash common.Hash, header types.Header) {
	if !w.enabled() {
		return
	}
	w.batch = []WatchMessage{}
	w.header = header
	w.height = height
	w.blockHash = blockHash
	w.cumulativeGas = make(map[uint64]uint64)
	w.gasUsed = 0
	w.blockTxs = []common.Hash{}
}

func (w *Watcher) SaveEthereumTx(msg types2.MsgEthereumTx, txHash common.Hash, index uint64) {
	if !w.enabled() {
		return
	}
	wMsg := NewMsgEthTx(&msg, txHash, w.blockHash, w.height, index)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
	w.UpdateBlockTxs(txHash)
}

func (w *Watcher) SaveContractCode(addr common.Address, code []byte) {
	if !w.enabled() {
		return
	}
	wMsg := NewMsgCode(addr, code, w.height)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) SaveTransactionReceipt(status uint32, msg types2.MsgEthereumTx, txHash common.Hash, txIndex uint64, data *types2.ResultData, gasUsed uint64) {
	if !w.enabled() {
		return
	}
	w.UpdateCumulativeGas(txIndex, gasUsed)
	wMsg := NewMsgTransactionReceipt(status, &msg, txHash, w.blockHash, txIndex, w.height, data, w.cumulativeGas[txIndex], gasUsed)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) UpdateCumulativeGas(txIndex, gasUsed uint64) {
	if !w.enabled() {
		return
	}
	if len(w.cumulativeGas) == 0 {
		w.cumulativeGas[txIndex] = gasUsed
	} else {
		w.cumulativeGas[txIndex] = w.cumulativeGas[txIndex-1] + gasUsed
	}
	w.gasUsed += gasUsed
}

func (w *Watcher) UpdateBlockTxs(txHash common.Hash) {
	if !w.enabled() {
		return
	}
	w.blockTxs = append(w.blockTxs, txHash)
}

func (w *Watcher) SaveBlock(bloom ethtypes.Bloom) {
	if !w.enabled() {
		return
	}
	wMsg := NewMsgBlock(w.height, bloom, w.blockHash, w.header, uint64(0xffffffff), big.NewInt(int64(w.gasUsed)), w.blockTxs)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}

	wInfo := NewMsgBlockInfo(w.height, w.blockHash)
	if wInfo != nil {
		w.batch = append(w.batch, wInfo)
	}
	w.SaveLatestHeight(w.height)
}

func (w *Watcher) SaveLatestHeight(height uint64) {
	if !w.enabled() {
		return
	}
	wMsg := NewMsgLatestHeight(height)
	if wMsg != nil {
		w.batch = append(w.batch, wMsg)
	}
}

func (w *Watcher) Commit() {
	if !w.enabled() {
		return
	}
	//hold it in temp
	batch := w.batch
	go func() {
		for _, b := range batch {
			w.store.Set([]byte(b.GetKey()), []byte(b.GetValue()))
		}
	}()
}

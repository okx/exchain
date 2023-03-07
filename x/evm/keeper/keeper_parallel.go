package keeper

import (
	"math/big"
	"sync"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"

	abci "github.com/okx/okbchain/libs/tendermint/abci/types"

	"github.com/okx/okbchain/x/evm/types"
)

func (k *Keeper) FixLog(tx []sdk.Tx, logIndex []int, hasEnterEvmTx []bool, anteErrs []error, resp []abci.ResponseDeliverTx) [][]byte {
	txSize := len(logIndex)
	res := make([][]byte, txSize, txSize)
	logSize := uint(0)
	txInBlock := -1
	k.Bloom = new(big.Int)

	for index := 0; index < txSize; index++ {
		if hasEnterEvmTx[index] {
			txInBlock++
		}
		rs, ok := k.LogsManages.Get(logIndex[index])
		if ok && anteErrs[index] == nil && rs.ResultData != nil {
			for _, v := range rs.ResultData.Logs {
				v.Index = logSize
				v.TxIndex = uint(txInBlock)
				logSize++
			}

			k.Bloom = k.Bloom.Or(k.Bloom, rs.ResultData.Bloom.Big())
			data, err := types.EncodeResultData(rs.ResultData)
			if err != nil {
				panic(err)
			}
			res[index] = data
		}
		// save transaction and transactionReceipt to watcher
		k.saveParallelTxResult(tx[index], rs.ResultData, resp[index])
	}

	return res
}

type LogsManager struct {
	cnt     int
	mu      sync.RWMutex
	Results map[int]TxResult
}

func NewLogManager() *LogsManager {
	return &LogsManager{
		mu:      sync.RWMutex{},
		Results: make(map[int]TxResult),
	}
}

func (l *LogsManager) Set(value TxResult) int {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.cnt++
	l.Results[l.cnt] = value
	return l.cnt
}

func (l *LogsManager) Get(index int) (TxResult, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	data, ok := l.Results[index]
	return data, ok
}

func (l *LogsManager) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.Results)
}

func (l *LogsManager) Reset() {
	if l == nil {
		return
	}
	for k := range l.Results {
		delete(l.Results, k)
	}
	l.cnt = 0
}

type TxResult struct {
	ResultData *types.ResultData
}

func (k *Keeper) saveParallelTxResult(tx sdk.Tx, resultData *types.ResultData, resp abci.ResponseDeliverTx) {
	if !k.Watcher.Enabled() {
		return
	}
	k.Watcher.SaveParallelTx(tx, resultData, resp)
}

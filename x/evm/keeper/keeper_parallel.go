package keeper

import (
	"math/big"
	"sync"

	"github.com/okex/exchain/x/evm/types"
)

func (k *Keeper) FixLog(logIndex []int, anteErrs []error) [][]byte {
	txSize := len(logIndex)
	res := make([][]byte, txSize, txSize)
	logSize := uint(0)
	txInBlock := -1
	k.Bloom = new(big.Int)

	for index := 0; index < txSize; index++ {
		rs, ok := k.LogsManages.Get(logIndex[index])
		if !ok || anteErrs[index] != nil {
			continue
		}
		txInBlock++
		if rs.ResultData == nil {
			continue
		}

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
	k.LogsManages.Reset()
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
	l.Results = make(map[int]TxResult)
	l.cnt = 0
}

type TxResult struct {
	ResultData *types.ResultData
	Err        error
}

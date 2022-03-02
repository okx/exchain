package keeper

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/okex/exchain/x/evm/types"
)

func (k *Keeper) FixLog(execResults [][]string) [][]byte {
	res := make([][]byte, len(execResults), len(execResults))
	logSize := uint(0)
	txInBlock := int(-1)
	k.Bloom = new(big.Int)

	for index := 0; index < len(execResults); index++ {
		rs, ok := k.LogsManages.Get(execResults[index][0])
		if !ok || execResults[index][1] != "" {
			fmt.Println("log-err", index, ok, execResults[index][1])
			continue
		}
		txInBlock++
		if rs.ResultData == nil {
			fmt.Println("billll", index)
			continue
		}

		for _, v := range rs.ResultData.Logs {
			v.Index = logSize
			v.TxIndex = uint(txInBlock)
			logSize++
		}
		//
		//ss := sha256.New()
		//ss.Write(rs.ResultData.Bloom.Bytes())
		//sum := ss.Sum(nil)
		//fmt.Println("log", index, hex.EncodeToString(sum))
		k.Bloom = k.Bloom.Or(k.Bloom, rs.ResultData.Bloom.Big())
		data, err := types.EncodeResultData(*rs.ResultData)
		if err != nil {
			panic(err)
		}
		res[index] = data
	}
	k.LogsManages.Reset()
	return res
}

type LogsManager struct {
	mu      sync.RWMutex
	Results map[string]TxResult
}

func NewLogManager() *LogsManager {
	return &LogsManager{
		mu:      sync.RWMutex{},
		Results: make(map[string]TxResult),
	}
}

func (l *LogsManager) Set(txBytes string, value TxResult) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.Results[txBytes]; ok {
		fmt.Println("attention!!!!!!!!", hex.EncodeToString([]byte(txBytes)), value.Err, value.ResultData.TxHash.String())
	}
	l.Results[txBytes] = value
}

func (l *LogsManager) Get(txBytes string) (TxResult, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	data, ok := l.Results[txBytes]
	return data, ok
}

func (l *LogsManager) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.Results)
}

func (l *LogsManager) Reset() {
	l.Results = make(map[string]TxResult)
}

type TxResult struct {
	ResultData *types.ResultData
	Err        error
}

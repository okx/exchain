package app

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/trace"
)

var (
	once         sync.Once
	CUSTOM_PRINT = []string{
		trace.Evm,
		trace.Delta,
		trace.Iavl,
		trace.FlatKV,
		trace.DeliverTxs,
		trace.EvmHandlerDetail,

		trace.RunAnteDetail,
		trace.AnteChainDetail,
		trace.Round,
		trace.CommitRound,
		trace.Produce}

	DefaultElapsedSchemas = fmt.Sprintf("%s=1,%s=1,%s=1,%s=1,%s=1,%s=1,%s=1,%s=1,%s=0,%s=0,%s=0",
		trace.Evm,
		trace.Delta,
		trace.Iavl,
		trace.FlatKV,
		trace.DeliverTxs,
		trace.EvmHandlerDetail,
		trace.RunAnteDetail,
		trace.AnteChainDetail,
		trace.Round,
		trace.CommitRound,
		trace.Produce)
)

const (
	Elapsed = "elapsed"
)

func init() {
	once.Do(func() {
		elapsedInfo := &ElapsedTimeInfos{
			infoMap:   make(map[string]string),
			schemaMap: make(map[string]bool),
		}

		elapsedInfo.decodeElapseParam(DefaultElapsedSchemas)

		trace.SetInfoObject(elapsedInfo)
	})
}

type ElapsedTimeInfos struct {
	mtx         sync.Mutex
	infoMap     map[string]string
	schemaMap   map[string]bool
	initialized bool
	elapsedTime int64
}

func (e *ElapsedTimeInfos) AddInfo(key string, info string) {
	if len(key) == 0 || len(info) == 0 {
		return
	}

	e.mtx.Lock()
	defer e.mtx.Unlock()

	e.infoMap[key] = info
}

func (e *ElapsedTimeInfos) Dump(logger log.Logger) {

	e.mtx.Lock()
	defer e.mtx.Unlock()

	if len(e.infoMap) == 0 {
		return
	}

	if !e.initialized {
		e.decodeElapseParam(viper.GetString(Elapsed))
		e.initialized = true
	}

	var detailInfo string
	for _, k := range CUSTOM_PRINT {
		if v, ok := e.schemaMap[k]; ok {
			if v {
				detailInfo += fmt.Sprintf("%s[%s], ", k, e.infoMap[k])
			}
		}
	}

	info := fmt.Sprintf("%s<%s>, %s<%s>, %s<%s>, %s<%s>, %s<%s>, %s<%s>, %s[%s], %s[%s], %s<%s>, %s<%s>",
		trace.Height, e.infoMap[trace.Height],
		trace.Tx, e.infoMap[trace.Tx],
		trace.BlockSize, e.infoMap[trace.BlockSize],
		trace.GasUsed, e.infoMap[trace.GasUsed],
		trace.WtxRatio, e.infoMap[trace.WtxRatio],
		trace.InvalidTxs, e.infoMap[trace.InvalidTxs],
		trace.RunTx, e.infoMap[trace.RunTx],
		trace.Prerun, e.infoMap[trace.Prerun],
		trace.MempoolCheckTxCnt, e.infoMap[trace.MempoolCheckTxCnt],
		trace.MempoolTxsCnt, e.infoMap[trace.MempoolTxsCnt],
	)

	if len(detailInfo) > 0 {
		detailInfo = strings.TrimRight(detailInfo, ", ")
		info += ", " + detailInfo
	}

	logger.Info(info)
	e.infoMap = make(map[string]string)
}

func (e *ElapsedTimeInfos) decodeElapseParam(elapsed string) {

	// suppose elapsd is like Evm=x,Iavl=x,DeliverTxs=x,DB=x,Round=x,CommitRound=x,Produce=x
	elapsdA := strings.Split(elapsed, ",")
	for _, v := range elapsdA {
		setVal := strings.Split(v, "=")
		if len(setVal) == 2 && setVal[1] == "1" {
			e.schemaMap[setVal[0]] = true
		}
	}
}

func (e *ElapsedTimeInfos) SetElapsedTime(elapsedTime int64) {
	e.elapsedTime = elapsedTime
}

func (e *ElapsedTimeInfos) GetElapsedTime() int64 {
	return e.elapsedTime
}

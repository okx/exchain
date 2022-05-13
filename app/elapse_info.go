package app

import (
	"fmt"
	"github.com/okex/exchain/libs/system/trace"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type SchemaConfig struct {
	schema  string
	enabled int
}

var (
	once         sync.Once
	CUSTOM_PRINT = []SchemaConfig{
		{trace.Evm, 1},
		{trace.Delta, 1},
		{trace.Iavl, 1},
		{trace.DeliverTxs, 1},
		{trace.EvmHandlerDetail, 0},

		{trace.RunAnteDetail, 0},
		{trace.AnteChainDetail, 0},
		{trace.Round, 0},
		{trace.CommitRound, 0},
		{trace.RecvBlock, 1},
		{trace.Produce, 0},
	}

	DefaultElapsedSchemas string
)

const (
	Elapsed = "elapsed"
)

func init() {
	for _, k := range CUSTOM_PRINT {
		DefaultElapsedSchemas += fmt.Sprintf("%s=%d,", k.schema, k.enabled)
	}

	elapsedInfo := &ElapsedTimeInfos{
		infoMap:   make(map[string]string),
		schemaMap: make(map[string]bool),
	}

	elapsedInfo.decodeElapseParam(DefaultElapsedSchemas)
	trace.SetInfoObject(elapsedInfo)

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

func (e *ElapsedTimeInfos) Dump(input interface{}) {

	logger, ok := input.(log.Logger)
	if !ok {
		panic("Invalid input")
	}
	e.mtx.Lock()
	defer e.mtx.Unlock()

	if _, ok := e.infoMap[trace.Height]; !ok {
		return
	}

	if !e.initialized {
		e.decodeElapseParam(viper.GetString(Elapsed))
		e.initialized = true
	}

	var detailInfo string
	for _, k := range CUSTOM_PRINT {
		if v, ok := e.schemaMap[k.schema]; ok {
			if v {
				detailInfo += fmt.Sprintf("%s[%s], ", k.schema, e.infoMap[k.schema])
			}
		}
	}

	info := fmt.Sprintf("%s<%s>, %s<%s>, %s<%s>, %s[%s], %s[%s], %s<%s>, %s<%s>, %s[%s], %s[%s], %s<%s>, %s<%s>, %s<%s>",
		trace.Height, e.infoMap[trace.Height],
		trace.Tx, e.infoMap[trace.Tx],
		trace.BlockSize, e.infoMap[trace.BlockSize],
		trace.BlockCompress, e.infoMap[trace.BlockCompress],
		trace.BlockUncompress, e.infoMap[trace.BlockUncompress],
		trace.GasUsed, e.infoMap[trace.GasUsed],
		trace.InvalidTxs, e.infoMap[trace.InvalidTxs],
		trace.RunTx, e.infoMap[trace.RunTx],
		trace.Prerun, e.infoMap[trace.Prerun],
		trace.MempoolCheckTxCnt, e.infoMap[trace.MempoolCheckTxCnt],
		trace.MempoolTxsCnt, e.infoMap[trace.MempoolTxsCnt],
		trace.SigCacheRatio, e.infoMap[trace.SigCacheRatio],
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

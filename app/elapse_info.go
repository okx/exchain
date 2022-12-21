package app

import (
	"fmt"
	"strings"
	"sync"

	"github.com/okex/exchain/libs/system/trace"
	"github.com/okex/exchain/libs/tendermint/libs/log"

	"github.com/spf13/viper"
)

type SchemaConfig struct {
	schema  string
	enabled int
}

var (
	optionalSchemas = []SchemaConfig{
		{trace.MempoolCheckTxCnt, 0},
		{trace.MempoolCheckTxTime, 0},
		{trace.SigCacheRatio, 0},
		{trace.Evm, 1},
		{trace.Delta, 1},
		{trace.Iavl, 1},
		{trace.DeliverTxs, 1},
		{trace.EvmHandlerDetail, 0},

		{trace.IavlRuntime, 0},
		{trace.RunAnteDetail, 0},
		{trace.AnteChainDetail, 0},
		{trace.Round, 0},
		{trace.CommitRound, 0},
		//{trace.RecvBlock, 1},
		{trace.First2LastPart, 0},
		{trace.BlockParts, 0},
		{trace.BlockPartsP2P, 0},
		{trace.Produce, 0},
		{trace.CompressBlock, 0},
		{trace.UncompressBlock, 0},
	}

	mandatorySchemas = []string{
		trace.Height,
		trace.Tx,
		trace.SimTx,
		trace.BlockSize,
		trace.BTInterval,
		trace.LastBlockTime,
		trace.GasUsed,
		trace.SimGasUsed,
		trace.InvalidTxs,
		trace.LastRun,
		trace.RunTx,
		trace.Prerun,
		trace.MempoolTxsCnt,
		trace.Workload,
		trace.ACOffset,
		trace.PersistDetails,
	}

	DefaultElapsedSchemas string
)

const (
	Elapsed = "elapsed"
)

func init() {
	for _, k := range optionalSchemas {
		DefaultElapsedSchemas += fmt.Sprintf("%s=%d,", k.schema, k.enabled)
	}

	elapsedInfo := &ElapsedTimeInfos{
		infoMap:   make(map[string]string),
		schemaMap: make(map[string]struct{}),
	}

	elapsedInfo.decodeElapseParam(DefaultElapsedSchemas)
	trace.SetInfoObject(elapsedInfo)

}

type ElapsedTimeInfos struct {
	mtx         sync.Mutex
	infoMap     map[string]string
	schemaMap   map[string]struct{}
	initialized bool
	elapsedTime int64
}

func (e *ElapsedTimeInfos) AddInfo(key string, info string) {
	if len(key) == 0 || len(info) == 0 {
		return
	}

	_, ok := e.schemaMap[key]
	if !ok {
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

	var mandatoryInfo string
	for _, key := range mandatorySchemas {
		_, ok := e.infoMap[key]
		if !ok {
			continue
		}
		mandatoryInfo += fmt.Sprintf("%s<%s>, ", key, e.infoMap[key])
	}

	var optionalInfo string
	var comma string
	for _, k := range optionalSchemas {
		if _, found := e.schemaMap[k.schema]; found {
			_, ok := e.infoMap[k.schema]
			if !ok {
				continue
			}
			optionalInfo += fmt.Sprintf("%s%s[%s]", comma, k.schema, e.infoMap[k.schema])
			comma = ", "
		}
	}

	logger.Info(mandatoryInfo + optionalInfo)
	e.infoMap = make(map[string]string)
}

func (e *ElapsedTimeInfos) decodeElapseParam(elapsed string) {
	// elapsed looks like: Evm=x,Iavl=x,DeliverTxs=x,DB=x,Round=x,CommitRound=x,Produce=x,IavlRuntime=x
	elapsedKV := strings.Split(elapsed, ",")
	for _, v := range elapsedKV {
		setVal := strings.Split(v, "=")
		if len(setVal) == 2 && setVal[1] == "1" {
			e.schemaMap[setVal[0]] = struct{}{}
		}
	}

	for _, key := range mandatorySchemas {
		e.schemaMap[key] = struct{}{}
	}
}

func (e *ElapsedTimeInfos) SetElapsedTime(elapsedTime int64) {
	e.elapsedTime = elapsedTime
}

func (e *ElapsedTimeInfos) GetElapsedTime() int64 {
	return e.elapsedTime
}

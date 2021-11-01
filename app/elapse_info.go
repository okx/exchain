package app

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"sync"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/trace"
)

var (
	once         sync.Once
	CUSTOM_PRINT = []string{trace.Evm, "Iavl", "DeliverTxs", trace.Round, trace.CommitRound, trace.Produce}
)

const (
	DefaultElapsedSchemas = "Evm=1,Iavl=1,DeliverTxs=1,Round=0,CommitRound=0,Produce=0"
	Elapsed = "elapsed"
)

func init() {
	once.Do(func() {
		elapsedInfo := &ElapsedTimeInfos{
			infoMap:     make(map[string]string),
			schemaMap: make(map[string]bool),
		}

		elapsedInfo.decodeElapseParam(DefaultElapsedSchemas)

		trace.SetInfoObject(elapsedInfo)
	})
}

type ElapsedTimeInfos struct {
	infoMap         map[string]string
	schemaMap       map[string]bool
	initialized     bool
	elapsedTime     int64
}

func (e *ElapsedTimeInfos) AddInfo(key string, info string) {
	if len(key) == 0 || len(info) == 0 {
		return
	}

	e.infoMap[key] = info
}

func (e *ElapsedTimeInfos) Dump(logger log.Logger) {

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

	info := fmt.Sprintf("%s<%s>, %s<%s>, %s<%s>, %s[%s]",
		trace.Height, e.infoMap[trace.Height],
		trace.Tx, e.infoMap[trace.Tx],
		trace.GasUsed, e.infoMap[trace.GasUsed],
		trace.RunTx, e.infoMap[trace.RunTx],
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

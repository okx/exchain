package app

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"sync"

	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/trace"
)

var (
	once         sync.Once
	CUSTOM_PRINT = []string{trace.Evm, "DB", "DeliverTxs", trace.Round, trace.CommitRound, trace.Produce}
)

const (
	Elapsed = "elapsed"
)

func init() {
	once.Do(func() {
		elapsedInfo := &ElapsedTimeInfos{
			infoMap:     make(map[string]string),
			showFlagMap: make(map[string]struct{}),
		}
		trace.SetInfoObject(elapsedInfo)
	})
}

type ElapsedTimeInfos struct {
	infoMap         map[string]string
	showFlagMap     map[string]struct{}
	showFlagInitail bool
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

	e.decodeElapsParam()

	var detailInfo string
	for _, v := range CUSTOM_PRINT {
		if _, ok := e.showFlagMap[v]; ok {
			detailInfo += fmt.Sprintf("%s[%s], ", v, e.infoMap[v])
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

func (e *ElapsedTimeInfos) decodeElapsParam() {

	if e.showFlagInitail {
		return
	}

	elapsd := viper.GetString(Elapsed)
	// suppose elapsd is like Evm=x,DeliverTx=x,DB=x,Round=x,CommitRound=x,Produce=x
	elapsdA := strings.Split(elapsd, ",")
	for _, v := range elapsdA {
		setVal := strings.Split(v, "=")
		if len(setVal) == 2 && setVal[1] == "1" {
			e.showFlagMap[setVal[0]] = struct{}{}
		}
	}
	e.showFlagInitail = true
}

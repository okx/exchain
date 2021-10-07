package app

import (
	"fmt"
	elapse "github.com/tendermint/tendermint/trace"
	"github.com/tendermint/tendermint/libs/log"

	"sync"
)
var once sync.Once

func init() {
	once.Do(func() {
		elapsedInfo := &ElapsedTimeInfos{
			infoMap: make(map[string]string),
		}
		elapse.SetInfoObject(elapsedInfo)
	})
}

type ElapsedTimeInfos struct {
	infoMap map[string]string
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

	info := fmt.Sprintf("%s<%s>, %s<%s>, %s<%s>, %s[%s], %s[%s]",
		elapse.Height, e.infoMap[elapse.Height],
		elapse.Tx, e.infoMap[elapse.Tx],
		elapse.GasUsed, e.infoMap[elapse.GasUsed],
		elapse.Produce, e.infoMap[elapse.Produce],
		elapse.RunTx, e.infoMap[elapse.RunTx],
		)

	logger.Info(info)
	e.infoMap = make(map[string]string)
}
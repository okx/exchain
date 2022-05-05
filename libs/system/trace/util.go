package trace

import (
	"github.com/spf13/viper"
	"sync"
)

const FlagEnableAnalyzer string = "enable-analyzer"

var (
	openAnalyzer bool
	dynamicConfig IDynamicConfig = MockDynamicConfig{}
	forceAnalyzerTags map[string]struct{}
	status            bool
	once              sync.Once
)

func EnableAnalyzer(flag bool)  {
	status = flag
}

func initForceAnalyzerTags() {
	forceAnalyzerTags = map[string]struct{}{
		RunAnte: {},
		Refund:  {},
		RunMsg:  {},
	}
}

func init() {
	initForceAnalyzerTags()

	dbOper = newDbRecord()
	for _, v := range STATEDB_READ {
		dbOper.AddOperType(v, READ)
	}
	for _, v := range STATEDB_WRITE {
		dbOper.AddOperType(v, WRITE)
	}
	for _, v := range EVM_OPER {
		dbOper.AddOperType(v, EVMALL)
	}
}

func OnAppBeginBlockEnter(height int64) {
	analyzer.reset(height)
	if !dynamicConfig.GetEnableAnalyzer() {
		openAnalyzer = false
		return
	}
	openAnalyzer = true
	lastElapsedTime := GetElapsedInfo().GetElapsedTime()
	if singlePprofDumper != nil && lastElapsedTime > singlePprofDumper.triggerAbciElapsed {
		singlePprofDumper.cpuProfile(height)
	}
}

func skip(oper string) bool {
	if openAnalyzer {
		return false
	}
	_, ok := forceAnalyzerTags[oper]
	return !ok
}

func OnAppDeliverTxEnter() {
	if analyzer != nil {
		analyzer.onAppDeliverTxEnter()
	}
}

func OnCommitDone() {
	if analyzer != nil {
		analyzer.onCommitDone()
	}
}

func StartTxLog(oper string) {
	if !skip(oper) {
		analyzer.startTxLog(oper)
	}
}

func StopTxLog(oper string) {
	if !skip(oper) {
		analyzer.stopTxLog(oper)
	}
}

func SetDynamicConfig(c IDynamicConfig) {
	dynamicConfig = c
}

type IDynamicConfig interface {
	GetEnableAnalyzer() bool
}

type MockDynamicConfig struct {
}

func (c MockDynamicConfig) GetEnableAnalyzer() bool {
	return viper.GetBool(FlagEnableAnalyzer)
}
package analyzer

import (
	"fmt"
	"github.com/okex/exchain/libs/system/trace"

	//"github.com/okex/exchain/libs/system/trace"
	"strconv"
	"strings"
	"sync"

	sm "github.com/okex/exchain/libs/tendermint/state"
	"github.com/spf13/viper"
)

const FlagEnableAnalyzer string = "enable-analyzer"

type insertFuncType func(tag string, elaspe int64)

var (
	singleAnalys *analyer
	openAnalyzer bool

	dynamicConfig IDynamicConfig = MockDynamicConfig{}

	forceAnalyzerTags map[string]struct{}

	isParalleledTxOn *bool
	insertElapse     insertFuncType
	once             sync.Once
)

func SetInsertFunc(f insertFuncType)  {
	once.Do(func() {
		insertElapse = f
	})
}

func initForceAnalyzerTags() {
	forceAnalyzerTags = map[string]struct{}{
		trace.RunAnte: {},
		trace.Refund:  {},
		trace.RunMsg:  {},
	}
}

func init() {
	initForceAnalyzerTags()
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

type analyer struct {
	status         bool
	currentTxIndex int64
	blockHeight    int64
	dbRead         int64
	dbWrite        int64
	evmCost        int64
	txs            []*txLog
}

func init() {
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

func newAnalys(height int64) {
	if isParalleledTxOn == nil {
		isParalleledTxOn = new(bool)
		*isParalleledTxOn =  sm.DeliverTxsExecMode(viper.GetInt(sm.FlagDeliverTxsExecMode)) != sm.DeliverTxsExecModeParallel
	}

	singleAnalys = &analyer{}
	singleAnalys.blockHeight = height
	singleAnalys.status = *isParalleledTxOn
}

func OnAppBeginBlockEnter(height int64) {
	newAnalys(height)
	if !dynamicConfig.GetEnableAnalyzer() {
		openAnalyzer = false
		return
	}
	openAnalyzer = true
	lastElapsedTime := trace.GetElapsedInfo().GetElapsedTime()
	if singlePprofDumper != nil && lastElapsedTime > singlePprofDumper.triggerAbciElapsed {
		singlePprofDumper.cpuProfile(height)
	}
}

func skip(a *analyer, oper string) bool {
	if a != nil {
		if openAnalyzer {
			return false
		}
		_, ok := forceAnalyzerTags[oper]
		return !ok
	} else {
		return true
	}
}

func OnAppDeliverTxEnter() {
	if singleAnalys != nil {
		singleAnalys.onAppDeliverTxEnter()
	}
}

func OnCommitDone() {
	if singleAnalys != nil {
		singleAnalys.onCommitDone()
	}
}

func StartTxLog(oper string) {
	if !skip(singleAnalys, oper) {
		singleAnalys.startTxLog(oper)
	}
}

func StopTxLog(oper string) {
	if !skip(singleAnalys, oper) {
		singleAnalys.stopTxLog(oper)
	}
}

func (s *analyer) onAppDeliverTxEnter() {
	if s.status {
		s.newTxLog()
	}
}

func (s *analyer) onCommitDone() {
	if s.status {
		s.format()
	}
}

func (s *analyer) newTxLog() {
	s.currentTxIndex++
	s.txs = append(s.txs, newTxLog())
}

func (s *analyer) startTxLog(oper string) {
	if s.status {
		if s.currentTxIndex > 0 && int64(len(s.txs)) == s.currentTxIndex {
			s.txs[s.currentTxIndex-1].StartTxLog(oper)
		}
	}
}

func (s *analyer) stopTxLog(oper string) {
	if s.status {
		if s.currentTxIndex > 0 && int64(len(s.txs)) == s.currentTxIndex {
			s.txs[s.currentTxIndex-1].StopTxLog(oper)
		}
	}
}


func (s *analyer) format() {

	evmcore, record := s.genRecord()

	if !openAnalyzer {
		formatNecessaryDeliverTx(record)
		return
	}
	formatDeliverTx(record)
	formatRunAnteDetail(record)
	formatEvmHandlerDetail(record)

	if insertElapse != nil {
		for k, v := range record {
			insertElapse(k, v)
		}
	}

	// evm
	trace.GetElapsedInfo().AddInfo(trace.Evm, fmt.Sprintf(EVM_FORMAT, s.dbRead, s.dbWrite, evmcore-s.dbRead-s.dbWrite))
}

// formatRecord format the record in the format fmt.Sprintf(", %s<%dms>", v, record[v])
func formatRecord(i int, key string, ms int64) string {
	t := strconv.FormatInt(ms, 10)
	b := strings.Builder{}
	b.Grow(2 + len(key) + 1 + len(t) + 3)
	if i != 0 {
		b.WriteString(", ")
	}
	b.WriteString(key)
	b.WriteString("<")
	b.WriteString(t)
	b.WriteString("ms>")
	return b.String()
}

func addInfo(name string, keys []string, record map[string]int64) {
	var strs = make([]string, len(keys))
	length := 0
	for i, v := range keys {
		strs[i] = formatRecord(i, v, record[v])
		length += len(strs[i])
	}
	builder := strings.Builder{}
	builder.Grow(length)
	for _, v := range strs {
		builder.WriteString(v)
	}
	trace.GetElapsedInfo().AddInfo(name, builder.String())
}

func (s *analyer) genRecord() (int64, map[string]int64) {
	var evmcore int64
	var record = make(map[string]int64)
	for _, v := range s.txs {
		for oper, operObj := range v.Record {
			operType := dbOper.GetOperType(oper)
			switch operType {
			case READ:
				s.dbRead += operObj.TimeCost
			case WRITE:
				s.dbWrite += operObj.TimeCost
			case EVMALL:
				evmcore += operObj.TimeCost
			default:
				if _, ok := record[oper]; !ok {
					record[oper] = operObj.TimeCost
				} else {
					record[oper] += operObj.TimeCost
				}
			}
		}
	}

	return evmcore, record
}

func formatNecessaryDeliverTx(record map[string]int64) {
	// deliver txs
	var deliverTxsKeys = []string{
		trace.RunAnte,
		trace.RunMsg,
		trace.Refund,
	}
	addInfo(trace.DeliverTxs, deliverTxsKeys, record)
}

func formatDeliverTx(record map[string]int64) {

	// deliver txs
	var deliverTxsKeys = []string{
		//----- DeliverTx
		//bam.DeliverTx,
		//bam.TxDecoder,
		//bam.RunTx,
		//----- run_tx
		//bam.InitCtx,
		trace.ValTxMsgs,
		trace.RunAnte,
		trace.RunMsg,
		trace.Refund,
		trace.EvmHandler,
	}
	addInfo(trace.DeliverTxs, deliverTxsKeys, record)
}

func formatEvmHandlerDetail(record map[string]int64) {

	// run msg
	var evmHandlerKeys = []string{
		//bam.ConsumeGas,
		//bam.Recover,
		//----- handler
		//bam.EvmHandler,
		//bam.ParseChainID,
		//bam.VerifySig,
		trace.Txhash,
		trace.SaveTx,
		trace.TransitionDb,
		//bam.Bloomfilter,
		//bam.EmitEvents,
		//bam.HandlerDefer,
		//-----
	}
	addInfo(trace.EvmHandlerDetail, evmHandlerKeys, record)
}

func formatRunAnteDetail(record map[string]int64) {

	// ante
	var anteKeys = []string{
		trace.CacheTxContext,
		trace.AnteChain,
		trace.AnteOther,
		trace.CacheStoreWrite,
	}
	addInfo(trace.RunAnteDetail, anteKeys, record)

}

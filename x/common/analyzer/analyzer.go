package analyzer

import (
	"fmt"
	"strconv"
	"strings"

	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sm "github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/trace"
	"github.com/spf13/viper"
)

const FlagEnableAnalyzer string = "enable-analyzer"

var (
	singleAnalys *analyer
	openAnalyzer bool

	dynamicConfig IDynamicConfig = MockDynamicConfig{}

	forceAnalyzerTags map[string]struct{}
)

func initForceAnalyzerTags() {
	forceAnalyzerTags = map[string]struct{}{
		bam.RunAnte: {},
		bam.Refund:  {},
		bam.RunMsg:  {},
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
	if singleAnalys == nil {
		singleAnalys = &analyer{
			status:      !viper.GetBool(sm.FlagParalleledTx),
			blockHeight: height,
		}
	}
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
	singleAnalys = nil
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
	singleAnalys = nil
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

	formatDeliverTx(record)
	formatRunAnteDetail(record)
	formatEvmHandlerDetail(record)
	// evm
	trace.GetElapsedInfo().AddInfo(trace.Evm, fmt.Sprintf(EVM_FORMAT, s.dbRead, s.dbWrite, evmcore-s.dbRead-s.dbWrite))
}

func addInfo(name string, keys []string, record map[string]int64) {
	var strs = make([]string, len(keys))
	length := 0
	for i, v := range keys {
		t := strconv.FormatInt(record[v], 10)
		b := strings.Builder{}
		b.Grow(2 + len(v) + 1 + len(t) + 3)
		b.WriteString(", ")
		b.WriteString(v)
		b.WriteString("<")
		b.WriteString(t)
		b.WriteString("ms>")
		length += b.Len()
		strs[i] = b.String()
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

func formatDeliverTx(record map[string]int64) {

	// deliver txs
	var deliverTxsKeys = []string{
		//----- DeliverTx
		//bam.DeliverTx,
		//bam.TxDecoder,
		//bam.RunTx,
		//----- run_tx
		//bam.InitCtx,
		bam.ValTxMsgs,
		bam.RunAnte,
		bam.RunMsg,
		bam.Refund,
		bam.EvmHandler,
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
		bam.Txhash,
		bam.SaveTx,
		bam.TransitionDb,
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
		bam.CacheTxContext,
		bam.AnteChain,
		bam.AnteOther,
		bam.CacheStoreWrite,
	}
	addInfo(trace.RunAnteDetail, anteKeys, record)

}

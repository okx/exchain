package analyzer

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
	"sync"

	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sm "github.com/okex/exchain/libs/tendermint/state"
	"github.com/okex/exchain/libs/tendermint/trace"
)

const FlagEnableAnalyzer string = "enable-analyzer"

var (
	singleAnalys *analyer
	openAnalyzer bool
	once         sync.Once
)

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

func getOpen() bool {
	once.Do(func() {
		openAnalyzer = viper.GetBool(FlagEnableAnalyzer)
	})
	return openAnalyzer
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
	if !getOpen() {
		return
	}
	newAnalys(height)
	lastElapsedTime := trace.GetElapsedInfo().GetElapsedTime()
	if singlePprofDumper != nil && lastElapsedTime > singlePprofDumper.triggerAbciElapsed {
		singlePprofDumper.cpuProfile(height)
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
	if singleAnalys != nil {
		singleAnalys.startTxLog(oper)
	}
}

func StopTxLog(oper string) {
	if singleAnalys != nil {
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
	var evmcore int64
	var format string
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

	var keys = []string{
		//----- DeliverTx
		//bam.DeliverTx,
		//bam.TxDecoder,
		//bam.RunTx,
		//----- run_tx
		//bam.InitCtx,
		bam.ValTxMsgs,
		bam.AnteHandler,
		bam.RunMsgs,
		bam.Refund,
		//bam.ConsumeGas,
		//bam.Recover,
		//----- handler
		bam.EvmHandler,
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

	for _, v := range keys {
		format += fmt.Sprintf("%s<%dms>, ", v, record[v])

	}
	format = strings.TrimRight(format, ", ")
	trace.GetElapsedInfo().AddInfo(trace.Evm, fmt.Sprintf(EVM_FORMAT, s.dbRead, s.dbWrite, evmcore-s.dbRead-s.dbWrite))

	trace.GetElapsedInfo().AddInfo(trace.DeliverTxs, format)
}

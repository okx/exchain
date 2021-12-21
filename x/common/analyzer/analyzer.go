package analyzer

import (
	"fmt"
	sm "github.com/okex/exchain/libs/tendermint/state"
	"github.com/spf13/viper"
	"strings"

	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	"github.com/okex/exchain/libs/tendermint/trace"
)

var singleAnalys *analyer

type analyer struct {
	status          bool
	currentTxIndex  int64
	blockHeight     int64
	startBeginBlock int64
	beginBlockCost  int64
	startdelliverTx int64
	delliverTxCost  int64
	startEndBlock   int64
	endBlockCost    int64
	startCommit     int64
	commitCost      int64
	dbRead          int64
	dbWrite         int64
	allCost         int64
	evmCost         int64
	txs             []*txLog
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
	singleAnalys.onAppBeginBlockEnter()
	lastElapsedTime := trace.GetElapsedInfo().GetElapsedTime()
	if singlePprofDumper != nil && lastElapsedTime > singlePprofDumper.triggerAbciElapsed {
		singlePprofDumper.cpuProfile(height)
	}
}

func OnAppBeginBlockExit() {
	if singleAnalys != nil {
		singleAnalys.onAppBeginBlockExit()
	}
}

func OnAppDeliverTxEnter() {
	if singleAnalys != nil {
		singleAnalys.onAppDeliverTxEnter()
	}
}

func OnAppDeliverTxExit() {
	if singleAnalys != nil {
		singleAnalys.onAppDeliverTxExit()
	}
}

func OnAppEndBlockEnter() {
	if singleAnalys != nil {
		singleAnalys.onAppEndBlockEnter()
	}
}

func OnAppEndBlockExit() {
	if singleAnalys != nil {
		singleAnalys.onAppEndBlockExit()
	}
}

func OnCommitEnter() {
	if singleAnalys != nil {
		singleAnalys.onCommitEnter()
	}
}

func OnCommitExit() {
	if singleAnalys != nil {
		singleAnalys.onCommitExit()
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

func (s *analyer) onAppBeginBlockEnter() {
	if s.status {
		s.startBeginBlock = GetNowTimeMs()
	}
}

func (s *analyer) onAppBeginBlockExit() {
	if s.status {
		s.beginBlockCost = GetNowTimeMs() - s.startBeginBlock
	}
}

func (s *analyer) onAppDeliverTxEnter() {
	if s.status {
		s.startdelliverTx = GetNowTimeMs()
		s.newTxLog()
	}
}

func (s *analyer) onAppDeliverTxExit() {
	if s.status {
		s.delliverTxCost += GetNowTimeMs() - s.startdelliverTx
	}
}

func (s *analyer) onAppEndBlockEnter() {
	if s.status {
		s.startEndBlock = GetNowTimeMs()
	}
}

func (s *analyer) onAppEndBlockExit() {
	if s.status {
		s.endBlockCost = GetNowTimeMs() - s.startEndBlock
	}
}

func (s *analyer) onCommitEnter() {
	if s.status {
		s.startCommit = GetNowTimeMs()
	}
}

func (s *analyer) onCommitExit() {
	if s.status {
		s.commitCost = GetNowTimeMs() - s.startCommit
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
	s.allCost = s.beginBlockCost + s.delliverTxCost + s.endBlockCost + s.commitCost
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

	trace.GetElapsedInfo().AddInfo("DeliverTxs", format)
}

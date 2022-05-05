package trace

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	analyzer *Analyzer = &Analyzer{}
)

type Analyzer struct {
	status         bool
	currentTxIndex int64
	blockHeight    int64
	dbRead         int64
	dbWrite        int64
	evmCost        int64
	txs            []*txLog
}

func (s *Analyzer) reset (height int64) {
	s.status = status
	s.currentTxIndex = 0
	s.blockHeight = height
	s.dbRead = 0
	s.dbWrite = 0
	s.evmCost = 0
	s.txs = nil
}

func (s *Analyzer) onAppDeliverTxEnter() {
	if s.status {
		s.newTxLog()
	}
}

func (s *Analyzer) onCommitDone() {
	if s.status {
		s.format()
	}
}

func (s *Analyzer) newTxLog() {
	s.currentTxIndex++
	s.txs = append(s.txs, newTxLog())
}

func (s *Analyzer) startTxLog(oper string) {
	if s.status {
		if s.currentTxIndex > 0 && int64(len(s.txs)) == s.currentTxIndex {
			s.txs[s.currentTxIndex-1].StartTxLog(oper)
		}
	}
}

func (s *Analyzer) stopTxLog(oper string) {
	if s.status {
		if s.currentTxIndex > 0 && int64(len(s.txs)) == s.currentTxIndex {
			s.txs[s.currentTxIndex-1].StopTxLog(oper)
		}
	}
}

func (s *Analyzer) format() {

	evmcore, record := s.genRecord()
	for k, v := range record {
		insertElapse(k, v)
	}

	if !openAnalyzer {
		formatNecessaryDeliverTx(record)
		return
	}
	formatDeliverTx(record)
	formatRunAnteDetail(record)
	formatEvmHandlerDetail(record)

	// evm
	GetElapsedInfo().AddInfo(Evm, fmt.Sprintf(EVM_FORMAT, s.dbRead, s.dbWrite, evmcore-s.dbRead-s.dbWrite))
}

func (s *Analyzer) genRecord() (int64, map[string]int64) {
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
		RunAnte,
		RunMsg,
		Refund,
	}
	addInfo(DeliverTxs, deliverTxsKeys, record)
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
		ValTxMsgs,
		RunAnte,
		RunMsg,
		Refund,
		//EvmHandler,
	}
	addInfo(DeliverTxs, deliverTxsKeys, record)
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
		Txhash,
		SaveTx,
		TransitionDb,
		//bam.Bloomfilter,
		//bam.EmitEvents,
		//bam.HandlerDefer,
		//-----
	}
	addInfo(EvmHandlerDetail, evmHandlerKeys, record)
}

func formatRunAnteDetail(record map[string]int64) {
	// ante
	var anteKeys = []string{
		CacheTxContext,
		AnteChain,
		AnteOther,
		CacheStoreWrite,
	}
	addInfo(RunAnteDetail, anteKeys, record)
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
	GetElapsedInfo().AddInfo(name, builder.String())
}
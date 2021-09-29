package analyzer

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/log"
	"time"
)

var singleAnalys *analyer

type analyer struct {
	logger          log.Logger
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
	allCost         int64
	evmCost         int64
	tx              []*txLog
}

func init() {
	dbOper = newDbRecord()
	for _, v := range STATEDB_READ {
		dbOper.AddOperType(v, READ)
	}
	for _, v := range STATEDB_WRITE {
		dbOper.AddOperType(v, WRITE)
	}
}

func newAnalys(log log.Logger, height int64) {
	if singleAnalys == nil {
		singleAnalys = &analyer{
			logger:      log,
			status:      true,
			blockHeight: height,
		}
	}
}

func OnAppBeginBlockEnter(log log.Logger, height int64) {
	newAnalys(log, height)

	singleAnalys.onAppBeginBlockEnter()
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

func StartTxLog(module, oper string) {
	if singleAnalys != nil {
		singleAnalys.startTxLog(module, oper)
	}
}

func StopTxLog(module, oper string) {
	if singleAnalys != nil {
		singleAnalys.stopTxLog(module, oper)
	}
}

func CloseAnalys() {
	singleAnalys.Close()
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
		//format to print analyzer and release current
		s.formatLog()
	}
}

func (s *analyer) newTxLog() {
	s.currentTxIndex++
	s.tx = append(s.tx, newTxLog(module))
}

func (s *analyer) startTxLog(module, oper string) {
	if s.status {
		if s.currentTxIndex > 0 && int64(len(s.tx)) == s.currentTxIndex {
			s.tx[s.currentTxIndex-1].StartTxLog(module, oper)
		}
	}
}

func (s *analyer) stopTxLog(module, oper string) {
	if s.status {
		if s.currentTxIndex > 0 && int64(len(s.tx)) == s.currentTxIndex {
			s.tx[s.currentTxIndex-1].StopTxLog(module, oper)
		}
	}
}

func (s *analyer) Close() {
	s.status = false
}

func (s *analyer) formatLog() {
	var tx_detail, tx_debug string
	var debug bool
	s.allCost = s.beginBlockCost + s.delliverTxCost + s.endBlockCost + s.commitCost
	if s.allCost > 5*int64(time.Millisecond) {
		debug = true
	}
	for index, v := range s.tx {
		s.evmCost += v.EvmCost
		var txRead, txWrite int64

		for _, operMap := range v.Record {
			tx_debug = ""
			for action, oper := range operMap.Record {
				operType, err := dbOper.GetOperType(action)
				if err != nil {
					continue
				}
				if operType == READ {
					txRead += oper.TimeCost
				}
				if operType == WRITE {
					txWrite += oper.TimeCost
				}
				if debug {
					tx_debug += fmt.Sprintf(TX_DEBUG_FORMAT, action, oper.Count, oper.TimeCost)
				}
			}
		}
		tx_detail += fmt.Sprintf(TX_FORMAT, index+1, v.AllCost, txRead, txWrite, v.EvmCost)
		if debug {
			tx_detail += tx_debug
		}
	}

	s.logger.Info(fmt.Sprintf(BLOCK_FORMAT, s.blockHeight, s.allCost, s.evmCost, tx_detail))
}

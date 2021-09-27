package pkg

import (
	"fmt"
	"github.com/tendermint/tendermint/libs/log"
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
	tx              []*txLog
}

func NewAnalys(log log.Logger, height int64) *analyer {
	singleAnalys = &analyer{
		logger:      log,
		status:      true,
		blockHeight: height,
	}
	return singleAnalys
}

func GetCurrentAnalys() *analyer {
	return singleAnalys
}

func OnCommitExit() {
	if singleAnalys != nil {
		singleAnalys.OnCommitExit()
	}
	singleAnalys = nil
}

func CloseAnalys() {
	singleAnalys.Close()
}

func (s *analyer) OnAppBeginBlockEnter() {
	if s.status {
		s.startBeginBlock = GetNowTimeNs()
	}
}

func (s *analyer) OnAppBeginBlockExit() {
	if s.status {
		s.beginBlockCost = GetNowTimeNs() - s.startBeginBlock
	}
}

func (s *analyer) OnAppDeliverTxEnter() {
	if s.status {
		s.startdelliverTx = GetNowTimeNs()
		s.newTxLog()
	}
}

func (s *analyer) OnAppDeliverTxExit() {
	if s.status {
		s.delliverTxCost = GetNowTimeNs() - s.startdelliverTx
	}
}

func (s *analyer) OnAppEndBlockEnter() {
	if s.status {
		s.startEndBlock = GetNowTimeNs()
	}
}

func (s *analyer) OnAppEndBlockExit() {
	if s.status {
		s.endBlockCost = GetNowTimeNs() - s.startEndBlock
	}
}

func (s *analyer) OnCommitEnter() {
	if s.status {
		s.startCommit = GetNowTimeNs()
	}
}

func (s *analyer) OnCommitExit() {
	if s.status {
		s.commitCost = GetNowTimeNs() - s.startCommit
		//format to print log and release current
		s.formatLog()
	}
}

func (s *analyer) newTxLog() {
	s.currentTxIndex++
	s.tx = append(s.tx, newTxLog(module))
}

func (s *analyer) StartTxLog(module, oper string) {
	if s.status {
		if s.currentTxIndex > 0 && int64(len(s.tx)) == s.currentTxIndex {
			s.tx[s.currentTxIndex-1].StartTxLog(module, oper)
		}
	}
}

func (s *analyer) StopTxLog(module, oper string) {
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
	s.allCost = s.beginBlockCost + s.delliverTxCost + s.endBlockCost + s.commitCost
	var tx_info string
	for index, v := range s.tx {
		tx_info += fmt.Sprintf(TX_FORMAT, index+1, v.EvmCost)
		var tx_detail string
		for module, operMap := range v.Record {
			tx_detail += fmt.Sprintf("moduleName: %s", module)
			for action, oper := range operMap.Record {
				tx_detail += fmt.Sprintf(TX_DETAIL, action, oper.Count, oper.TimeCost, oper.Min, oper.Max, oper.Avg)
			}
		}
		tx_info += tx_detail
	}

	s.logger.Info(fmt.Sprintf(BLOCK_FORMAT, s.blockHeight, s.allCost, tx_info))
}

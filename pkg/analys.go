package pkg

import "time"

var singleAnalys *analyer

type analyer struct {
	status          bool
	currentTxIndex  int
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

func GetCurrentAnalys() *analyer {
	if singleAnalys != nil {
		return singleAnalys
	}
	return nil
}

func NewAnalys(height int64) *analyer {
	return &analyer{
		blockHeight: height,
	}
}

func (s *analyer) StartBeginBlock() {
	s.startBeginBlock = time.Now().Unix()
}

func (s *analyer) StopBeginBlock() {
	s.beginBlockCost = time.Now().Unix() - s.startBeginBlock
}

func (s *analyer) StartDelliverTx() {
	s.startdelliverTx = time.Now().Unix()
	s.newTxLog()
}

func (s *analyer) StopDelliverTx() {
	s.beginBlockCost = time.Now().Unix() - s.startdelliverTx
}

func (s *analyer) StartEndBlock() {
	s.startEndBlock = time.Now().Unix()
}

func (s *analyer) StopEndBlock() {
	s.endBlockCost = time.Now().Unix() - s.startEndBlock
}

func (s *analyer) StartCommitBlock() {
	s.startCommit = time.Now().Unix()
}

func (s *analyer) StopCommitBlock() {
	s.commitCost = time.Now().Unix() - s.startCommit
}

func (s *analyer) newTxLog() {
	s.currentTxIndex++
	s.tx = append(s.tx, newTxLog(module))
}

func (s *analyer) StartTxLog(module, oper string) {
	s.tx[s.currentTxIndex].StartTxLog(module, oper)
}

func (s *analyer) StopTxLog(module, oper string) {
	s.tx[s.currentTxIndex].StopTxLog(module, oper)
}

func (s *analyer) Stop() {
	s.allCost = s.beginBlockCost + s.delliverTxCost + s.endBlockCost + s.commitCost
	//print log
}

func (s *analyer) FormatLog() {
	// here to print the logs
}

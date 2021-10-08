package analyzer

var preSingleAnalys, singleAnalys *analyer

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

func newAnalys(height int64) {
	if singleAnalys == nil {
		singleAnalys = &analyer{
			status:      true,
			blockHeight: height,
		}
	}
}

func OnAppBeginBlockEnter(height int64) {
	newAnalys(height)
	singleAnalys.onAppBeginBlockEnter()
}

func OnAppBeginBlockExit() {
	if singleAnalys != nil {
		singleAnalys.onAppBeginBlockExit()
	}
}

func EvmCost() int64 {
	if preSingleAnalys != nil {
		return preSingleAnalys.EvmCost()
	}
	return 0
}

func DbReadCost() int64 {
	if preSingleAnalys != nil {
		return preSingleAnalys.DbReadCost()
	}
	return 0
}

func DbWriteCost() int64 {
	if preSingleAnalys != nil {
		return preSingleAnalys.DbWriteCost()
	}
	return 0
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
		preSingleAnalys = singleAnalys
	}
	singleAnalys = nil
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

func (s *analyer) EvmCost() int64 {
	return s.evmCost
}

func (s *analyer) DbReadCost() int64 {
	return s.dbRead
}

func (s *analyer) DbWriteCost() int64 {
	return s.dbWrite
}

func (s *analyer) format() {
	s.allCost = s.beginBlockCost + s.delliverTxCost + s.endBlockCost + s.commitCost
	for _, v := range s.tx {
		s.evmCost += v.EvmCost
		for _, operMap := range v.Record {
			for action, oper := range operMap.Record {
				operType, err := dbOper.GetOperType(action)
				if err != nil {
					continue
				}
				if operType == READ {
					s.dbRead += oper.TimeCost
				}
				if operType == WRITE {
					s.dbWrite += oper.TimeCost
				}
			}
		}
	}
}

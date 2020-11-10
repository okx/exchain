package orm

import (
	"github.com/okex/okexchain/x/backend/types"
)

// AddSwapInfo insert into swap token pairs details
func (orm *ORM) AddSwapInfo(swapInfos []*types.SwapInfo) (addedCnt int, err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)
	cnt := 0

	for _, swapInfo := range swapInfos {
		if swapInfo != nil {
			ret := tx.Create(swapInfo)
			if ret.Error != nil {
				return cnt, ret.Error
			} else {
				cnt++
			}
		}
	}

	tx.Commit()
	return cnt, nil
}

// nolint
func (orm *ORM) GetSwapInfo(startTime int64) []types.SwapInfo {
	var swapInfos []types.SwapInfo
	query := orm.db.Model(types.SwapInfo{}).Where("timestamp >= ?", startTime)

	query.Order("timestamp asc").Find(&swapInfos)
	return swapInfos
}

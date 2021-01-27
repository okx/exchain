package orm

import "github.com/okex/okexchain/x/backend/types"

// AddClaimInfo insert farm claimed coins into db
func (orm *ORM) AddClaimInfo(claimInfos []*types.ClaimInfo) (addedCnt int, err error) {
	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)
	cnt := 0

	for _, claimInfo := range claimInfos {
		if claimInfo != nil {
			ret := tx.Create(claimInfo)
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
func (orm *ORM) GetAccountClaimInfos(address string) []types.ClaimInfo {
	var claimInfos []types.ClaimInfo
	query := orm.db.Model(types.ClaimInfo{}).Where("address = ?", address)

	query.Order("timestamp asc").Find(&claimInfos)
	return claimInfos
}

func (orm *ORM) GetAccountClaimedByPool(address string, poolName string) []types.ClaimInfo {
	var claimInfos []types.ClaimInfo
	query := orm.db.Model(types.ClaimInfo{}).Where("address = ? and pool_name = ?", address, poolName)

	query.Order("timestamp asc").Find(&claimInfos)
	return claimInfos
}

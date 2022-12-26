package mysqldb

import (
	"github.com/okex/exchain/x/evm/statistics/orm/model"
)

func (mdb *mysqlDB) InsertClaim(claim model.Claim) {
	if *claim.Height <= mdb.latestSavedHeight {
		return
	}
	tx := mdb.db.Table("claim").Create(&claim)
	if tx.Error != nil {
		panic(tx.Error)
	}
}

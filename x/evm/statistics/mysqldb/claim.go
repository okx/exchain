package mysqldb

import (
	"github.com/okex/exchain/x/evm/statistics/orm/model"
	"time"
)

const (
	// ExpireDate date --date='2023-02-01 00:00:00' +%s
	TermThreshold = 1675209600
)

func (mdb *mysqlDB) InsertClaim(claim model.Claim) {
	if *claim.Height <= mdb.latestSavedHeight {
		return
	}
	if (*claim.BlockTime).Add(time.Duration(*(claim.Term))*time.Duration(24)*time.Hour).Unix() > TermThreshold {
		return
	}
	tx := mdb.db.Table("claim").Create(&claim)
	if tx.Error != nil {
		panic(tx.Error)
	}
}

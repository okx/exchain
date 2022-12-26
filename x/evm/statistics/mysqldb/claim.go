package mysqldb

import (
	"github.com/okex/exchain/x/evm/statistics/orm/model"
)

func (mdb *mysqlDB) InsertClaim(claim model.Claim) {
	if len(mdb.claimBatch) > batchSize {
		tx := mdb.db.Table("claim").CreateInBatches(mdb.claimBatch, len(mdb.claimBatch))
		if tx.Error != nil {
			panic(tx.Error)
		}

		mdb.claimBatch = nil
	}
	mdb.claimBatch = append(mdb.claimBatch, claim)
}

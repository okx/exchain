package mysqldb

import (
	"github.com/okex/exchain/x/evm/statistics/orm/model"
)

func (mdb *mysqlDB) InsertClaim(claim model.Claim) {
	if len(mdb.claimBatch) >= batchSize {
		tx := mdb.db.CreateInBatches(mdb.claimBatch, len(mdb.claimBatch))
		if tx.Error != nil {
			panic(tx.Error)
		}
		mdb.claimBatch = mdb.claimBatch[:0]
	} else {
		mdb.claimBatch = append(mdb.claimBatch, claim)
	}
}

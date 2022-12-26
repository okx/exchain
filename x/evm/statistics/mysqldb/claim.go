package mysqldb

import (
	"github.com/okex/exchain/x/evm/statistics/orm/model"
	"log"
)

func (mdb *mysqlDB) InsertClaim(claim model.Claim) {
	log.Println("insert claim to db")
	if len(mdb.claimBatch) >= batchSize {
		tx := mdb.db.CreateInBatches(mdb.claimBatch, len(mdb.claimBatch))
		log.Printf("insert claim %v\n", len(mdb.claimBatch))
		if tx.Error != nil {
			panic(tx.Error)
		}
		mdb.claimBatch = mdb.claimBatch[:0]
	} else {
		mdb.claimBatch = append(mdb.claimBatch, claim)
	}
}

package mysqldb

import (
	"github.com/okex/exchain/x/evm/statistics/orm/model"
	"log"
	"runtime/debug"
)

func (mdb *mysqlDB) InsertClaim(claim model.Claim) {
	if *claim.Useraddr == "0xacf041fc5a59978016e3b6c339b61a65762d10e2" {
		debug.PrintStack()
	}
	mdb.claimBatch = append(mdb.claimBatch, claim)
	if len(mdb.claimBatch) >= batchSize {
		//		tx := mdb.db.Table("claim").CreateInBatches(mdb.claimBatch, len(mdb.claimBatch))
		//		if tx.Error != nil {
		//			panic(tx.Error)
		//		}
		//
		//		mdb.claimBatch = nil
		for i, v := range mdb.claimBatch {
			log.Printf("%v %v\n", i, v.Useraddr)
			if *v.Useraddr == "0xacf041fc5a59978016e3b6c339b61a65762d10e2" {
				log.Printf("---giskook %v %v \n", i, v)
			}
			tx := mdb.db.Table("claim").Create(&v)
			if tx.Error != nil {
				panic(tx.Error)
			}
		}
		mdb.claimBatch = make([]model.Claim, 0)
	}
}

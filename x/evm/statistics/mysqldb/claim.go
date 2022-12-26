package mysqldb

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/x/evm/statistics/orm/model"
	"log"
)

func (mdb *mysqlDB) InsertClaim(claim model.Claim) {
	if mdb.claimSavedHeight == 0 {
		mdb.claimSavedHeight = global.GetGlobalHeight()
	}
	if *claim.Height >= mdb.claimSavedHeight+2 && len(mdb.claimBatch) > 0 {
		for _, v := range mdb.claimBatch {
			if *v.Useraddr == "0xacf041fc5a59978016e3b6c339b61a65762d10e2" {
				log.Println(v)
			}
			tx := mdb.db.Table("claim").Create(&v)
			if tx.Error != nil {
				panic(tx.Error)
			}
		}
		//		tx := mdb.db.Table("claim").CreateInBatches(mdb.claimBatch, len(mdb.claimBatch))
		//		if tx.Error != nil {
		//			panic(tx.Error)
		//		}
		height := *mdb.claimBatch[0].Height
		for _, v := range mdb.claimBatch {
			if *v.Height != height {
				panic(fmt.Sprintf("%v", height))
			}
		}
		log.Printf("insert claim height %v batch %v\n", mdb.claimSavedHeight+1, len(mdb.claimBatch))

		mdb.claimBatch = nil
		mdb.claimSavedHeight++
	}
	mdb.claimBatch = append(mdb.claimBatch, claim)
}

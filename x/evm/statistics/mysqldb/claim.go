package mysqldb

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/x/evm/statistics/orm/model"
	"log"
)

func (mdb *mysqlDB) InsertClaim(claim model.Claim) {
	if *claim.Useraddr == "0xacf041fc5a59978016e3b6c339b61a65762d10e2" {
		log.Printf("giskook -------giskook %v\n", claim)
	}
	if mdb.claimSavedHeight == 0 {
		mdb.claimSavedHeight = global.GetGlobalHeight()
	}
	if *claim.Height >= mdb.claimSavedHeight+2 && len(mdb.claimBatch) > 0 {
		tx := mdb.db.Table("claim").CreateInBatches(mdb.claimBatch, len(mdb.claimBatch))
		if tx.Error != nil {
			panic(tx.Error)
		}
		height := *mdb.claimBatch[0].Height
		for _, v := range mdb.claimBatch {
			if *v.Height != height {
				panic(fmt.Sprintf("%v", height))
			}
		}

		mdb.claimBatch = mdb.claimBatch[:0]
		mdb.claimSavedHeight++
	}
	mdb.claimBatch = append(mdb.claimBatch, claim)
}

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
		log.Printf("giskook height %v %v\n", *claim.Height, mdb.claimSavedHeight)
	}
	if *claim.Height == mdb.claimSavedHeight+2 {
		tx := mdb.db.CreateInBatches(mdb.claimBatch, len(mdb.claimBatch))
		if tx.Error != nil {
			panic(tx.Error)
		}
		height := *claim.Height
		for _, v := range mdb.claimBatch {
			if *v.Height != height-1 {
				panic(fmt.Sprintf("%v", height))
			}
		}
		log.Printf("insert claim %v %v\n", *claim.Height, len(mdb.claimBatch))

		mdb.claimBatch = mdb.claimBatch[:0]
		mdb.claimSavedHeight++
	}
	mdb.claimBatch = append(mdb.claimBatch, claim)
}

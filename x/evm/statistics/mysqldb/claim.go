package mysqldb

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/x/evm/statistics/orm/model"
	"log"
)

func (mdb *mysqlDB) InsertClaim(claim model.Claim) {
	log.Printf("giskook height %v %v\n", *claim.Height, global.GetGlobalHeight())
	if global.GetGlobalHeight()+2 == *claim.Height {
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
	}
	mdb.claimBatch = append(mdb.claimBatch, claim)
}

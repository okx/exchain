package mysqldb

import (
	"github.com/okex/exchain/x/evm/statistics/orm/model"
	"log"
)

func (mdb *mysqlDB) InsertReward(reward model.Reward) {
	if len(mdb.rewardBatch) >= batchSize {
		tx := mdb.db.CreateInBatches(mdb.rewardBatch, len(mdb.rewardBatch))
		if tx.Error != nil {
			panic(tx.Error)
		}
		log.Printf("insert reward %v\n", len(mdb.rewardBatch))
		mdb.rewardBatch = mdb.rewardBatch[:0]
	} else {
		mdb.rewardBatch = append(mdb.rewardBatch, reward)
	}
}

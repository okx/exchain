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
		var rewards model.Reward
		if rows, err := tx.Rows(); err == nil {
			if rows.Next() {
				e := tx.ScanRows(rows, &rewards)
				if e != nil {
					panic(e)
				}
				log.Println(rewards)
			}
		} else {
			panic(err)
		}

		mdb.rewardBatch = mdb.rewardBatch[:0]
	} else {
		mdb.rewardBatch = append(mdb.rewardBatch, reward)
		//mdb.db.Model(model.Claim{}).Where("useraddr=? and reward=0", reward.Useraddr)
	}
}

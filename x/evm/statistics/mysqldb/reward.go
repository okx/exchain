package mysqldb

import (
	"github.com/okex/exchain/x/evm/statistics/orm/model"
)

func (mdb *mysqlDB) InsertReward(reward model.Reward) {
	if len(mdb.rewardBatch) >= batchSize {
		tx := mdb.db.CreateInBatches(mdb.rewardBatch, len(mdb.rewardBatch))
		if tx.Error != nil {
			panic(tx.Error)
		}
		var claims []model.Claim
		useraddr := ""
		//for _, _ = range mdb.rewardBatch {
		useraddr = "0xacf041fc5a59978016e3b6c339b61a65762d10e2"
		mdb.db.Table("claim").Where("useraddr=? and reward=0", useraddr).Find(&claims)
		//		if len(claims) > 1 || len(claims) == 0 {
		//			panic(fmt.Sprintf("error claims %v %v", len(claims), claims))
		//		}
		mdb.db.Table("claim").Where("useraddr=? and reward=0", useraddr).Update("reward", 1)
		//}

		mdb.rewardBatch = mdb.rewardBatch[:0]
	} else {
		mdb.rewardBatch = append(mdb.rewardBatch, reward)
		//mdb.db.Model(model.Claim{}).Where("useraddr=? and reward=0", reward.Useraddr)
	}
}

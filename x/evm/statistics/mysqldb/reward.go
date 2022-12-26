package mysqldb

import (
	"fmt"
	"github.com/okex/exchain/x/evm/statistics/orm/model"
)

func (mdb *mysqlDB) InsertReward(reward model.Reward) {
	if *reward.Height <= mdb.latestSavedHeight {
		return
	}
	tx := mdb.db.Table("reward").Create(&reward)
	if tx.Error != nil {
		panic(tx.Error)
	}
	var dbReward model.Reward
	mdb.db.Table("reward").Where("useraddr=?", *reward.Useraddr).Last(&dbReward)

	userAddr := *dbReward.Useraddr
	var claims []model.Claim
	tx = mdb.db.Table("claim").Where("useraddr=? and reward=0", userAddr).Find(&claims)
	if tx.Error != nil {
		panic(tx.Error)
	}
	if len(claims) != 1 {
		panic(fmt.Sprintf("useraddr %v dup or empty %v", userAddr, len(claims)))
	}
	var r int64 = 1
	tx = mdb.db.Table("claim").Model(&claims[0]).Updates(&model.Claim{Reward: &r, RewardID: &dbReward.ID})
	if tx.Error != nil {
		panic(tx.Error)
	}
}

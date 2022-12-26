package mysqldb

import (
	"fmt"
	"github.com/okex/exchain/x/evm/statistics/orm/model"
	"log"
	"runtime/debug"
)

func (mdb *mysqlDB) InsertReward(reward model.Reward) {
	tx := mdb.db.Table("reward").Create(&reward)
	if tx.Error != nil {
		panic(tx.Error)
	}
	var dbReward model.Reward
	tx.Last(&dbReward)
	log.Printf("%v %v %v\n", dbReward.ID, *dbReward.Txhash, *dbReward.Useraddr)
	debug.PrintStack()

	userAddr := "0xacf041fc5a59978016e3b6c339b61a65762d10e2"
	//useraddr := *dbReward.Useraddr
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

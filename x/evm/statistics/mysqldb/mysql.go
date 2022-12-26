package mysqldb

import (
	"fmt"
	"github.com/okex/exchain/x/evm/statistics/orm/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	MysqlConfig = "okc:okcpassword@(localhost:3306)/xen_stats?charset=utf8mb4&parseTime=True&loc=Local"
	batchSize   = 100
)

var db *mysqlDB

func init() {
	db = &mysqlDB{}
}

func GetInstance() *mysqlDB {
	return db
}

type mysqlDB struct {
	db               *gorm.DB
	claimBatch       []model.Claim
	rewardBatch      []model.Reward
	claimSavedHeight int64
}

func (mdb *mysqlDB) Init() {
	var err error
	mdb.db, err = gorm.Open(mysql.Open(MysqlConfig))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}
}

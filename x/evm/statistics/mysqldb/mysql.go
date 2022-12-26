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
	db                *gorm.DB
	claimBatch        []model.Claim
	rewardBatch       []model.Reward
	latestSavedHeight int64
}

func (mdb *mysqlDB) Init() {
	var err error
	mdb.db, err = gorm.Open(mysql.Open(MysqlConfig))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}
}

func (mdb *mysqlDB) GetMaxHeight(table string) int64 {
	var claim model.Claim
	tx := mdb.db.Table(table).Last(&claim)
	if tx.Error != nil && tx.Error.Error() == "record not found" {
		return 0
	}
	if tx.Error != nil {
		panic(tx.Error)
	}
	return *claim.Height
}

func (mdb *mysqlDB) DeleteHeight(table string, height int64) {
	tx := mdb.db.Table(table).Where("height=?", height).Delete(&model.Claim{})
	if tx.Error != nil {
		panic(tx.Error)
	}
}

func (mdb *mysqlDB) GetLatestHeightAndDeleteHeight() {
	claimHeight := mdb.GetMaxHeight("claim")
	rewardHeight := mdb.GetMaxHeight("reward")
	mdb.DeleteHeight("claim", claimHeight)
	mdb.DeleteHeight("reward", rewardHeight)
	if rewardHeight > claimHeight {
		panic("reward should not greater than claim")
	}
	if claimHeight > 0 {
		mdb.latestSavedHeight = claimHeight - 1
	}
}

func (mdb *mysqlDB) GetLatestSavedHeight() int64 {
	return mdb.latestSavedHeight
}

func (mdb *mysqlDB) DeleteFromHeight(height int64) {
	tx := mdb.db.Table("claim").Where("height>=?", height).Delete(&model.Claim{})
	if tx.Error != nil {
		panic(tx.Error)
	}
	tx = mdb.db.Table("reward").Where("height>=?", height).Delete(&model.Reward{})
	if tx.Error != nil {
		panic(tx.Error)
	}
}

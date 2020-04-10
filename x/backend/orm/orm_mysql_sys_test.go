package orm

import (
	"testing"

	"github.com/okex/okchain/x/backend/types"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/token"
)

type DangrousORM struct {
	*ORM
}

func (orm *DangrousORM) CleanupDataInTestEvn() (err error) {

	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)

	r := tx.Delete(&types.Deal{})
	r = tx.Delete(&types.Order{})
	r = tx.Delete(&token.FeeDetail{})
	r = tx.Delete(&types.Transaction{})
	r = tx.Delete(&types.MatchResult{})

	if r.Error == nil {
		tx.Commit()
	} else {
		return r.Error
	}
	return nil
}

func NewMysqlORM() (orm *ORM, e error) {
	engineInfo := OrmEngineInfo{
		EngineType: EngineTypeMysql,
		ConnectStr: "okdexer:okdex123!@tcp(127.0.0.1:13306)/okdex?charset=utf8mb4&parseTime=True",
	}
	mysqlOrm, e := New(false, &engineInfo, nil)

	dorm := DangrousORM{mysqlOrm}
	dorm.CleanupDataInTestEvn()

	return mysqlOrm, e
}

func TestMysql_ORMDeals(t *testing.T) {
	common.SkipSysTestChecker(t)
	orm, _ := NewMysqlORM()
	testORMDeals(t, orm)
}

func TestMysql_FeeDetails(t *testing.T) {
	common.SkipSysTestChecker(t)
	orm, _ := NewMysqlORM()
	testORMFeeDetails(t, orm)
}

func TestMysql_Orders(t *testing.T) {
	common.SkipSysTestChecker(t)
	orm, _ := NewMysqlORM()
	testORMOrders(t, orm)
}

func TestMysql_Transactions(t *testing.T) {
	common.SkipSysTestChecker(t)
	orm, _ := NewMysqlORM()
	testORMTransactions(t, orm)
}

func TestNewORM_BatchInsert(t *testing.T) {
	common.SkipSysTestChecker(t)
	orm, _ := NewMysqlORM()
	testORMBatchInsert(t, orm)
}

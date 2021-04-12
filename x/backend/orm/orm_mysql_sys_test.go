package orm

import (
	"testing"

	"github.com/okex/exchain/x/backend/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/token"
)

type DangrousORM struct {
	*ORM
}

func (orm *DangrousORM) CleanupDataInTestEvn() (err error) {

	orm.singleEntryLock.Lock()
	defer orm.singleEntryLock.Unlock()

	tx := orm.db.Begin()
	defer orm.deferRollbackTx(tx, err)
	dealDB := tx.Delete(&types.Deal{})
	orderDB := tx.Delete(&types.Order{})
	feeDB := tx.Delete(&token.FeeDetail{})
	txDB := tx.Delete(&types.Transaction{})
	matchDB := tx.Delete(&types.MatchResult{})

	if err = types.NewErrorsMerged(dealDB.Error, orderDB.Error, feeDB.Error, txDB.Error, matchDB.Error); err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func NewMysqlORM() (orm *ORM, e error) {
	engineInfo := OrmEngineInfo{
		EngineType: EngineTypeMysql,
		ConnectStr: "okdexer:okdex123!@tcp(127.0.0.1:13306)/okdex?charset=utf8mb4&parseTime=True",
	}
	mysqlOrm, e := New(false, &engineInfo, nil)

	dorm := DangrousORM{mysqlOrm}
	if err := dorm.CleanupDataInTestEvn(); err != nil {
		return nil, err
	}

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

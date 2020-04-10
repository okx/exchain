package orm

import (
	"fmt"
	"log"
	"os"
	"time"
)

// MockSqlite3ORM create sqlite db for test, return orm
func MockSqlite3ORM() (*ORM, string) {
	dbDir := "/tmp"
	dbName := fmt.Sprintf("testdb_%010d.db", time.Now().Unix())

	orm, err := NewSqlite3ORM(false, dbDir, dbName, nil)
	if err != nil {
		fmt.Print("failed to create sqlite orm")
	}
	return orm, dbDir + "/" + dbName

}

// NewSqlite3ORM create sqlite db, return orm
func NewSqlite3ORM(enableLog bool, baseDir string, dbName string, logger *log.Logger) (orm *ORM, e error) {
	engineInfo := OrmEngineInfo{
		EngineType: EngineTypeSqlite,
		ConnectStr: baseDir + string(os.PathSeparator) + dbName,
	}
	return New(enableLog, &engineInfo, nil)
}

// CommitKlines insert into klines for test
func (orm *ORM) CommitKlines(klines ...[]interface{}) {
	tx := orm.db.Begin()
	for _, kline := range klines {
		for _, k := range kline {
			tx.Create(k)
		}
	}
	tx.Commit()
}

// DeleteDB remove the sqlite db
func DeleteDB(dbPath string) {
	if err := os.Remove(dbPath); err != nil {
		fmt.Print("failed to remove " + dbPath)
	}
}

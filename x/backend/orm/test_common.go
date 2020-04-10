package orm

import (
	"fmt"
	"log"
	"os"
	"time"
)

func MockSqlite3ORM() (*ORM, string) {
	dbDir := "/tmp"
	dbName := fmt.Sprintf("testdb_%010d.db", time.Now().Unix())

	orm, err := NewSqlite3ORM(false, dbDir, dbName, nil)
	if err != nil {
		fmt.Print("error: sqlite3 orm")
	}
	return orm, dbDir + "/" + dbName

}

func NewSqlite3ORM(enableLog bool, baseDir string, dbName string, logger *log.Logger) (orm *ORM, e error) {
	engineInfo := OrmEngineInfo{
		EngineType: EngineTypeSqlite,
		ConnectStr: baseDir + string(os.PathSeparator) + dbName,
	}
	return New(enableLog, &engineInfo, nil)
}

func (orm *ORM) MockCommitKlines(klines ...[]interface{}) {
	tx := orm.db.Begin()
	for _, kline := range klines {
		for _, k := range kline {
			tx.Create(k)
		}
	}
	tx.Commit()
}

func DeleteDB(dbPath string) {
	os.Remove(dbPath)
}

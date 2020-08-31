package analyservice

import (
	"fmt"

	"github.com/okex/okchain/x/backend/orm"

	"github.com/okex/okchain/x/backend"
)

func NewMysqlORM(url string) *backend.ORM {
	engineInfo := backend.OrmEngineInfo{
		EngineType: orm.EngineTypeMysql,
		ConnectStr: url + "?charset=utf8mb4&parseTime=True",
	}
	mysqlOrm, err := backend.NewORM(false, &engineInfo, nil)
	if err != nil {
		fmt.Println("error: ", err)
	}
	return mysqlOrm
}

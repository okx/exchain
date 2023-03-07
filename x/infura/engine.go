package infura

import (
	"errors"
	"fmt"

	"github.com/okx/okbchain/libs/tendermint/libs/log"
	"github.com/okx/okbchain/x/infura/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const batchSize = 1000

func newStreamEngine(cfg *types.Config, logger log.Logger) (types.IStreamEngine, error) {
	if cfg.MysqlUrl == "" {
		return nil, errors.New("infura.mysql-url is empty")
	}
	return newMySQLEngine(cfg.MysqlUrl, cfg.MysqlUser, cfg.MysqlPass, cfg.MysqlDB, logger)
}

type MySQLEngine struct {
	db     *gorm.DB
	logger log.Logger
}

func newMySQLEngine(url, user, pass, dbName string, l log.Logger) (types.IStreamEngine, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, url, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&types.TransactionReceipt{}, &types.TransactionLog{},
		&types.LogTopic{}, &types.Block{}, &types.Transaction{}, &types.ContractCode{})
	return &MySQLEngine{
		db:     db,
		logger: l,
	}, nil
}

func (e *MySQLEngine) Write(streamData types.IStreamData) bool {
	e.logger.Debug("Begin MySqlEngine write")
	data := streamData.ConvertEngineData()
	trx := e.db.Begin()
	// write TransactionReceipts
	for i := 0; i < len(data.TransactionReceipts); i += batchSize {
		end := i + batchSize
		if end > len(data.TransactionReceipts) {
			end = len(data.TransactionReceipts)
		}
		ret := trx.CreateInBatches(data.TransactionReceipts[i:end], len(data.TransactionReceipts[i:end]))
		if ret.Error != nil {
			return e.rollbackWithError(trx, ret.Error)
		}
	}

	// write Block
	ret := trx.Omit("Transactions").Create(data.Block)
	if ret.Error != nil {
		return e.rollbackWithError(trx, ret.Error)
	}

	// write Transactions
	for i := 0; i < len(data.Block.Transactions); i += batchSize {
		end := i + batchSize
		if end > len(data.Block.Transactions) {
			end = len(data.Block.Transactions)
		}
		ret := trx.CreateInBatches(data.Block.Transactions[i:end], len(data.Block.Transactions[i:end]))
		if ret.Error != nil {
			return e.rollbackWithError(trx, ret.Error)
		}
	}

	// write contract code
	for _, code := range data.ContractCodes {
		ret := trx.Create(code)
		if ret.Error != nil {
			return e.rollbackWithError(trx, ret.Error)
		}
	}

	trx.Commit()
	e.logger.Debug("End MySqlEngine write")
	return true
}

func (e *MySQLEngine) rollbackWithError(trx *gorm.DB, err error) bool {
	trx.Rollback()
	e.logger.Error(err.Error())
	return false
}

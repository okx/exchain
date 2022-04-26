package types

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	evm "github.com/okex/exchain/x/evm/watcher"
	"gorm.io/gorm"
)

type EngineData struct {
	TransactionReceipts []*TransactionReceipt
	//TransactionLogs     []*TransactionLog
	//LogTopics           []*LogTopic
}

type StreamData struct {
	TransactionReceipts []evm.TransactionReceipt
}

func (sd StreamData) ConvertEngineData() EngineData {
	transactionReceipts := make([]*TransactionReceipt, len(sd.TransactionReceipts))
	for i, t := range sd.TransactionReceipts {
		// convert  TransactionReceipt
		receipt := &TransactionReceipt{
			Status:            uint64(t.Status),
			CumulativeGasUsed: uint64(t.CumulativeGasUsed),
			TransactionHash:   t.TransactionHash,
			GasUsed:           uint64(t.GasUsed),
			BlockHash:         t.BlockHash,
			BlockNumber:       int64(t.BlockNumber),
			TransactionIndex:  uint64(t.TransactionIndex),
			From:              t.From,
		}
		if t.ContractAddress != nil {
			receipt.ContractAddress = t.ContractAddress.String()
		}
		if t.To != nil {
			receipt.To = t.To.String()
		}

		// convert  TransactionLog
		transactionLogs := make([]TransactionLog, len(t.Logs))
		for i, l := range t.Logs {
			log := TransactionLog{
				Address:          l.Address.String(),
				Data:             hexutil.Encode(l.Data),
				TransactionHash:  t.TransactionHash,
				TransactionIndex: receipt.TransactionIndex,
				LogIndex:         uint64(l.Index),
				BlockNumber:      receipt.BlockNumber,
				BlockHash:        t.BlockHash,
			}

			// convert  LogTopic
			logTopics := make([]LogTopic, len(l.Topics))
			for i, topic := range l.Topics {
				logTopics[i] = LogTopic{
					//TransactionHash: t.TransactionHash,
					//LogIndex:        uint64(l.Index),
					Topic: topic.String(),
				}
			}
			log.Topics = logTopics

			transactionLogs[i] = log
		}

		receipt.Logs = transactionLogs
		transactionReceipts[i] = receipt
	}

	return EngineData{
		TransactionReceipts: transactionReceipts,
		//TransactionLogs:     transactionLogs,
		//LogTopics:           logTopics,
	}
}

type TransactionReceipt struct {
	gorm.Model
	Status            uint64 `gorm:"type:tinyint(4)"`
	CumulativeGasUsed uint64 `gorm:"type:int(11)"`
	TransactionHash   string `gorm:"type:varchar(66);index:unique_hash,unique;not null"`
	ContractAddress   string `gorm:"type:varchar(42)"`
	GasUsed           uint64 `gorm:"type:int(11)"`
	BlockHash         string `gorm:"type:varchar(66)"`
	BlockNumber       int64
	TransactionIndex  uint64 `gorm:"type:int(11)"`
	From              string `gorm:"type:varchar(42)"`
	To                string `gorm:"type:varchar(42)"`
	Logs              []TransactionLog
}

type TransactionLog struct {
	gorm.Model
	Address              string `gorm:"type:varchar(42);index;not null"`
	Data                 string `gorm:"type:varchar(256)"`
	TransactionHash      string `gorm:"type:varchar(66);index;not null"`
	TransactionIndex     uint64 `gorm:"type:int(11)"`
	LogIndex             uint64 `gorm:"type:int(11)"`
	BlockHash            string `gorm:"type:varchar(66);index;not null"`
	BlockNumber          int64  `gorm:"index;not null"`
	TransactionReceiptID uint
	Topics               []LogTopic
}

type LogTopic struct {
	gorm.Model
	//TransactionHash string `gorm:"type:varchar(66);index;not null"`
	//LogIndex        uint64 `gorm:"index;not null"`
	Topic            string `gorm:"type:varchar(66);index;not null"`
	TransactionLogID uint
}

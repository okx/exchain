package types

import (
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	evm "github.com/okx/okbchain/x/evm/watcher"
	"gorm.io/gorm"
)

type EngineData struct {
	TransactionReceipts []*TransactionReceipt
	Block               *Block
	ContractCodes       []*ContractCode
}

type StreamData struct {
	TransactionReceipts []evm.TransactionReceipt
	Block               evm.Block
	Transactions        []evm.Transaction
	ContractCodes       map[string][]byte
}

func (sd StreamData) ConvertEngineData() EngineData {
	return EngineData{
		TransactionReceipts: convertTransactionReceipts(sd.TransactionReceipts),
		Block:               convertBlocks(sd.Block, sd.Transactions),
		ContractCodes:       convertContractCodes(sd.ContractCodes, int64(sd.Block.Number)),
	}
}

func convertTransactionReceipts(trs []evm.TransactionReceipt) []*TransactionReceipt {
	transactionReceipts := make([]*TransactionReceipt, len(trs))
	for i, t := range trs {
		// convert  TransactionReceipt
		receipt := &TransactionReceipt{
			Status:            uint64(t.Status),
			CumulativeGasUsed: uint64(t.CumulativeGasUsed),
			TransactionHash:   t.GetHash(),
			GasUsed:           uint64(t.GasUsed),
			BlockHash:         t.GetBlockHash(),
			BlockNumber:       int64(t.BlockNumber),
			TransactionIndex:  uint64(t.TransactionIndex),
			From:              t.GetFrom(),
		}
		if t.ContractAddress != nil {
			receipt.ContractAddress = t.ContractAddress.String()
		}
		to := t.GetTo()
		if to != nil {
			receipt.To = to.String()
		}

		// convert  TransactionLog
		transactionLogs := make([]TransactionLog, len(t.Logs))
		for i, l := range t.Logs {
			log := TransactionLog{
				Address:          l.Address.String(),
				Data:             hexutil.Encode(l.Data),
				TransactionHash:  receipt.TransactionHash,
				TransactionIndex: receipt.TransactionIndex,
				LogIndex:         uint64(l.Index),
				BlockNumber:      receipt.BlockNumber,
				BlockHash:        receipt.BlockHash,
			}

			// convert  LogTopic
			logTopics := make([]LogTopic, len(l.Topics))
			for i, topic := range l.Topics {
				logTopics[i] = LogTopic{
					Topic: topic.String(),
				}
			}
			log.Topics = logTopics

			transactionLogs[i] = log
		}

		receipt.Logs = transactionLogs
		transactionReceipts[i] = receipt
	}
	return transactionReceipts
}

func convertBlocks(evmBlock evm.Block, evmTransactions []evm.Transaction) *Block {
	block := &Block{
		Number:           int64(evmBlock.Number),
		Hash:             evmBlock.Hash.String(),
		ParentHash:       evmBlock.ParentHash.String(),
		TransactionsRoot: evmBlock.TransactionsRoot.String(),
		StateRoot:        evmBlock.StateRoot.String(),
		Miner:            evmBlock.Miner.String(),
		Size:             uint64(evmBlock.Size),
		GasLimit:         uint64(evmBlock.GasLimit),
		GasUsed:          evmBlock.GasUsed.ToInt().Uint64(),
		Timestamp:        uint64(evmBlock.Timestamp),
	}

	transactions := make([]*Transaction, len(evmTransactions))
	for i, t := range evmTransactions {
		tx := &Transaction{
			BlockHash:   t.BlockHash.String(),
			BlockNumber: t.BlockNumber.ToInt().Int64(),
			From:        t.From.String(),
			Gas:         uint64(t.Gas),
			GasPrice:    t.GasPrice.String(),
			Hash:        t.Hash.String(),
			Input:       t.Input.String(),
			Nonce:       uint64(t.Nonce),
			Index:       uint64(*t.TransactionIndex),
			Value:       t.Value.String(),
			V:           t.V.String(),
			R:           t.R.String(),
			S:           t.S.String(),
		}
		if t.To != nil {
			tx.To = t.To.String()
		}
		transactions[i] = tx
	}
	block.Transactions = transactions
	return block
}

func convertContractCodes(codes map[string][]byte, height int64) []*ContractCode {
	contractCodes := make([]*ContractCode, 0, len(codes))
	for k, v := range codes {
		contractCodes = append(contractCodes, &ContractCode{
			Address:     k,
			Code:        hexutil.Encode(v),
			BlockNumber: height,
		})
	}
	return contractCodes
}

type TransactionReceipt struct {
	gorm.Model
	Status            uint64 `gorm:"type:tinyint(4)"`
	CumulativeGasUsed uint64 `gorm:"type:int(11)"`
	TransactionHash   string `gorm:"type:varchar(66);index;not null"`
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
	Data                 string `gorm:"type:text"`
	TransactionHash      string `gorm:"type:varchar(66)"`
	TransactionIndex     uint64 `gorm:"type:int(11)"`
	LogIndex             uint64 `gorm:"type:int(11)"`
	BlockHash            string `gorm:"type:varchar(66);index;not null"`
	BlockNumber          int64  `gorm:"index;not null"`
	TransactionReceiptID uint
	Topics               []LogTopic
}

type LogTopic struct {
	gorm.Model
	Topic            string `gorm:"type:varchar(66)"`
	TransactionLogID uint
}

type Block struct {
	Number           int64  `gorm:"primaryKey"`
	Hash             string `gorm:"type:varchar(66);index;not null"`
	ParentHash       string `gorm:"type:varchar(66)"`
	TransactionsRoot string `gorm:"type:varchar(66)"`
	StateRoot        string `gorm:"type:varchar(66)"`
	Miner            string `gorm:"type:varchar(42)"`
	Size             uint64 `gorm:"type:int(11)"`
	GasLimit         uint64
	GasUsed          uint64
	Timestamp        uint64 `gorm:"type:int(11)"`
	Transactions     []*Transaction
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}

type Transaction struct {
	gorm.Model
	BlockHash   string `gorm:"type:varchar(66)"`
	BlockNumber int64
	From        string `gorm:"type:varchar(42)"`
	Gas         uint64 `gorm:"type:int(11)"`
	GasPrice    string `gorm:"type:varchar(66)"`
	Hash        string `gorm:"type:varchar(66)"`
	Input       string `gorm:"type:text"`
	Nonce       uint64 `gorm:"type:int(11)"`
	To          string `gorm:"type:varchar(42)"`
	Index       uint64 `gorm:"type:int(11)"`
	Value       string `gorm:"type:varchar(255)"`
	V           string `gorm:"type:varchar(255)"`
	R           string `gorm:"type:varchar(255)"`
	S           string `gorm:"type:varchar(255)"`
}

type ContractCode struct {
	gorm.Model
	Address     string `gorm:"type:varchar(42);index:unique_address,unique;not null"`
	Code        string
	BlockNumber int64
}

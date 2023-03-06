package watcher

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	prototypes "github.com/okx/okbchain/x/evm/watcher/proto"
)

func transactionToProto(tr *Transaction) *prototypes.Transaction {
	var to []byte
	if tr.To != nil {
		to = tr.To.Bytes()
	}
	return &prototypes.Transaction{
		BlockHash:        tr.BlockHash.Bytes(),
		BlockNumber:      tr.BlockNumber.String(),
		From:             tr.From.Bytes(),
		Gas:              uint64(tr.Gas),
		GasPrice:         tr.GasPrice.String(),
		Hash:             tr.Hash.Bytes(),
		Input:            tr.Input,
		Nonce:            uint64(tr.Nonce),
		To:               to,
		TransactionIndex: uint64(*tr.TransactionIndex),
		Value:            tr.Value.String(),
		V:                tr.V.String(),
		R:                tr.R.String(),
		S:                tr.S.String(),
	}
}

func protoToTransaction(tr *prototypes.Transaction) *Transaction {
	blockHash := common.BytesToHash(tr.BlockHash)
	blockNum := hexutil.MustDecodeBig(tr.BlockNumber)
	gasPrice := hexutil.MustDecodeBig(tr.GasPrice)
	var to *common.Address
	if len(tr.To) > 0 {
		addr := common.BytesToAddress(tr.To)
		to = &addr
	}
	index := hexutil.Uint64(tr.TransactionIndex)
	value := hexutil.MustDecodeBig(tr.Value)
	v := hexutil.MustDecodeBig(tr.V)
	r := hexutil.MustDecodeBig(tr.R)
	s := hexutil.MustDecodeBig(tr.S)
	return &Transaction{
		BlockHash:        &blockHash,
		BlockNumber:      (*hexutil.Big)(blockNum),
		From:             common.BytesToAddress(tr.From),
		Gas:              hexutil.Uint64(tr.Gas),
		GasPrice:         (*hexutil.Big)(gasPrice),
		Hash:             common.BytesToHash(tr.Hash),
		Input:            tr.Input,
		Nonce:            hexutil.Uint64(tr.Nonce),
		To:               to,
		TransactionIndex: &index,
		Value:            (*hexutil.Big)(value),
		V:                (*hexutil.Big)(v),
		R:                (*hexutil.Big)(r),
		S:                (*hexutil.Big)(s),
	}
}

func receiptToProto(tr *TransactionReceipt) *prototypes.TransactionReceipt {
	logs := make([]*prototypes.Log, len(tr.Logs))
	for i, log := range tr.Logs {
		topics := make([][]byte, len(log.Topics))
		for j, topic := range log.Topics {
			topics[j] = topic.Bytes()
		}
		logs[i] = &prototypes.Log{
			Address:     log.Address.Bytes(),
			Topics:      topics,
			Data:        log.Data,
			BlockNumber: log.BlockNumber,
			TxHash:      log.TxHash.Bytes(),
			TxIndex:     uint64(log.TxIndex),
			BlockHash:   log.BlockHash.Bytes(),
			Index:       uint64(log.Index),
			Removed:     log.Removed,
		}
	}
	var contractAddr []byte
	if tr.ContractAddress != nil {
		contractAddr = tr.ContractAddress.Bytes()
	}
	var to []byte
	if tr.To != nil {
		to = tr.To.Bytes()
	}
	return &prototypes.TransactionReceipt{
		Status:            uint64(tr.Status),
		CumulativeGasUsed: uint64(tr.CumulativeGasUsed),
		LogsBloom:         tr.LogsBloom.Bytes(),
		Logs:              logs,
		TransactionHash:   tr.TransactionHash,
		ContractAddress:   contractAddr,
		GasUsed:           uint64(tr.GasUsed),
		BlockHash:         tr.BlockHash,
		BlockNumber:       uint64(tr.BlockNumber),
		TransactionIndex:  uint64(tr.TransactionIndex),
		From:              tr.From,
		To:                to,
	}
}

func protoToReceipt(tr *prototypes.TransactionReceipt) *TransactionReceipt {
	logs := make([]*ethtypes.Log, len(tr.Logs))
	for i, log := range tr.Logs {
		topics := make([]common.Hash, len(log.Topics))
		for j, topic := range log.Topics {
			topics[j] = common.BytesToHash(topic)
		}
		logs[i] = &ethtypes.Log{
			Address:     common.BytesToAddress(log.Address),
			Topics:      topics,
			Data:        log.Data,
			BlockNumber: log.BlockNumber,
			TxHash:      common.BytesToHash(log.TxHash),
			TxIndex:     uint(log.TxIndex),
			BlockHash:   common.BytesToHash(log.BlockHash),
			Index:       uint(log.Index),
			Removed:     log.Removed,
		}
	}
	var contractAddr *common.Address
	if len(tr.ContractAddress) > 0 {
		addr := common.BytesToAddress(tr.ContractAddress)
		contractAddr = &addr
	}
	var to *common.Address
	if len(tr.To) > 0 {
		addr := common.BytesToAddress(tr.To)
		to = &addr
	}
	return &TransactionReceipt{
		Status:            hexutil.Uint64(tr.Status),
		CumulativeGasUsed: hexutil.Uint64(tr.CumulativeGasUsed),
		LogsBloom:         ethtypes.BytesToBloom(tr.LogsBloom),
		Logs:              logs,
		TransactionHash:   tr.TransactionHash,
		ContractAddress:   contractAddr,
		GasUsed:           hexutil.Uint64(tr.GasUsed),
		BlockHash:         tr.BlockHash,
		BlockNumber:       hexutil.Uint64(tr.BlockNumber),
		TransactionIndex:  hexutil.Uint64(tr.TransactionIndex),
		From:              tr.From,
		To:                to,
	}
}

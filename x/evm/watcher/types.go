package watcher

import (
	"encoding/json"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	rpctypes "github.com/okex/okexchain/app/rpc/types"
	"github.com/okex/okexchain/x/evm/types"
	"github.com/status-im/keycard-go/hexutils"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	prefixTx           = "0x1"
	prefixBlock        = "0x2"
	prefixReceipt      = "0x3"
	prefixCode         = "0x4"
	prefixBlockInfo    = "0x5"
	prefixLatestHeight = "0x6"

	KeyLatestHeight = "LatestHeight"

	TransactionSuccess = 1
	TransactionFailed  = 0
)

type WatchMessage interface {
	GetKey() string
	GetValue() string
}

type MsgEthTx struct {
	Key       string
	JsonEthTx string
}

func NewMsgEthTx(tx *types.MsgEthereumTx, txHash, blockHash common.Hash, height, index uint64) *MsgEthTx {
	ethTx, e := rpctypes.NewTransaction(tx, txHash, blockHash, height, index)
	if e != nil {
		return nil
	}
	jsTx, e := json.Marshal(ethTx)
	if e != nil {
		return nil
	}
	msg := MsgEthTx{
		Key:       txHash.String(),
		JsonEthTx: string(jsTx),
	}
	return &msg
}

func (m MsgEthTx) GetKey() string {
	return prefixTx + m.Key
}

func (m MsgEthTx) GetValue() string {
	return m.JsonEthTx
}

type MsgCode struct {
	Key  string
	Code string
}

type CodeInfo struct {
	Height uint64 `json:"height"`
	Code   string `json:"code"`
}

func NewMsgCode(contractAddr common.Address, code []byte, height uint64) *MsgCode {
	codeInfo := CodeInfo{
		Height: height,
		Code:   hexutils.BytesToHex(code),
	}
	jsCode, e := json.Marshal(codeInfo)
	if e != nil {
		return nil
	}
	return &MsgCode{
		Key:  contractAddr.String(),
		Code: string(jsCode),
	}
}

func (m MsgCode) GetKey() string {
	return prefixCode + m.Key
}

func (m MsgCode) GetValue() string {
	return m.Code
}

type MsgTransactionReceipt struct {
	txHash  string
	receipt string
}

type TransactionReceipt struct {
	Status            hexutil.Uint64  `json:"status"`
	CumulativeGasUsed hexutil.Uint64  `json:"cumulativeGasUsed"`
	LogsBloom         ethtypes.Bloom  `json:"logsBloom"`
	Logs              []*ethtypes.Log `json:"logs"`
	TransactionHash   string          `json:"transactionHash"`
	ContractAddress   *common.Address `json:"contractAddress"`
	GasUsed           hexutil.Uint64  `json:"gasUsed"`
	BlockHash         string          `json:"blockHash"`
	BlockNumber       hexutil.Uint64  `json:"blockNumber"`
	TransactionIndex  hexutil.Uint64  `json:"transactionIndex"`
	From              string          `json:"from"`
	To                *common.Address `json:"to"`
}

func NewMsgTransactionReceipt(status uint32, tx *types.MsgEthereumTx, txHash, blockHash common.Hash, txIndex, height uint64, data *types.ResultData, cumulativeGas, GasUsed uint64) *MsgTransactionReceipt {

	tr := TransactionReceipt{
		Status:            hexutil.Uint64(status),
		CumulativeGasUsed: hexutil.Uint64(cumulativeGas),
		LogsBloom:         data.Bloom,
		Logs:              data.Logs,
		TransactionHash:   txHash.String(),
		ContractAddress:   &data.ContractAddress,
		GasUsed:           hexutil.Uint64(GasUsed),
		BlockHash:         blockHash.String(),
		BlockNumber:       hexutil.Uint64(height),
		TransactionIndex:  hexutil.Uint64(txIndex),
		From:              common.BytesToAddress(tx.From().Bytes()).Hex(),
		To:                tx.To(),
	}

	//contract address will set to 0x0000000000000000000000000000000000000000 if contract deploy failed
	if tr.ContractAddress != nil && tr.ContractAddress.String() == "0x0000000000000000000000000000000000000000" {
		//set to nil to keep sycn with ethereum rpc
		tr.ContractAddress = nil
	}
	jsTr, e := json.Marshal(tr)
	if e != nil {
		return nil
	}
	return &MsgTransactionReceipt{txHash: txHash.String(), receipt: string(jsTr)}
}

func (m MsgTransactionReceipt) GetKey() string {
	return prefixReceipt + m.txHash
}

func (m MsgTransactionReceipt) GetValue() string {
	return m.receipt
}

type MsgBlock struct {
	blockHash string
	block     string
}

type EthBlock struct {
	Number           hexutil.Uint64 `json:"number"`
	Hash             common.Hash    `json:"hash"`
	ParentHash       common.Hash    `json:"parentHash"`
	Nonce            uint64         `json:"nonce"`
	Sha3Uncles       common.Hash    `json:"sha3Uncles"`
	LogsBloom        ethtypes.Bloom `json:"logsBloom"`
	TransactionsRoot common.Hash    `json:"transactionsRoot"`
	StateRoot        common.Hash    `json:"stateRoot"`
	Miner            common.Address `json:"miner"`
	MixHash          common.Hash    `json:"mixHash"`
	Difficulty       hexutil.Uint64 `json:"difficulty"`
	TotalDifficulty  hexutil.Uint64 `json:"totalDifficulty"`
	ExtraData        hexutil.Bytes  `json:"extraData"`
	Size             hexutil.Uint64 `json:"size"`
	GasLimit         hexutil.Uint64 `json:"gasLimit"`
	GasUsed          *hexutil.Big   `json:"gasUsed"`
	Timestamp        hexutil.Uint64 `json:"timestamp"`
	Uncles           []string       `json:"uncles"`
	ReceiptsRoot     common.Hash    `json:"receiptsRoot"`
	Transactions     interface{}    `json:"transactions"`
}

func NewMsgBlock(height uint64, blockBloom ethtypes.Bloom, blockHash common.Hash, header abci.Header, gasLimit uint64, gasUsed *big.Int, txs interface{}) *MsgBlock {
	b := EthBlock{
		Number:           hexutil.Uint64(height),
		Hash:             blockHash,
		ParentHash:       common.BytesToHash(header.LastBlockId.Hash),
		Nonce:            0,
		Sha3Uncles:       common.Hash{},
		LogsBloom:        blockBloom,
		TransactionsRoot: common.BytesToHash(header.DataHash),
		StateRoot:        common.BytesToHash(header.AppHash),
		Miner:            common.BytesToAddress(header.ProposerAddress),
		MixHash:          common.Hash{},
		Difficulty:       0,
		TotalDifficulty:  0,
		ExtraData:        nil,
		Size:             hexutil.Uint64(header.Size()),
		GasLimit:         hexutil.Uint64(gasLimit),
		GasUsed:          (*hexutil.Big)(gasUsed),
		Timestamp:        hexutil.Uint64(header.Time.Unix()),
		Uncles:           []string{},
		ReceiptsRoot:     common.Hash{},
		Transactions:     txs,
	}
	jsBlock, e := json.Marshal(b)
	if e != nil {
		return nil
	}
	return &MsgBlock{blockHash: blockHash.String(), block: string(jsBlock)}
}

func (m MsgBlock) GetKey() string {
	return prefixBlock + m.blockHash
}

func (m MsgBlock) GetValue() string {
	return m.block
}

type MsgBlockInfo struct {
	height string
	hash   string
}

func NewMsgBlockInfo(height uint64, blockHash common.Hash) *MsgBlockInfo {
	return &MsgBlockInfo{
		height: strconv.Itoa(int(height)),
		hash:   blockHash.String(),
	}
}

func (b MsgBlockInfo) GetKey() string {
	return prefixBlockInfo + b.height
}

func (b MsgBlockInfo) GetValue() string {
	return b.hash
}

type MsgLatestHeight struct {
	height string
}

func NewMsgLatestHeight(height uint64) *MsgLatestHeight {
	return &MsgLatestHeight{
		height: strconv.Itoa(int(height)),
	}
}

func (b MsgLatestHeight) GetKey() string {
	return prefixLatestHeight + KeyLatestHeight
}

func (b MsgLatestHeight) GetValue() string {
	return b.height
}

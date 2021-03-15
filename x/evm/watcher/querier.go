package watcher

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

type Querier struct {
	store *WatchStore
}

func NewQuerier() *Querier {
	return &Querier{store: InstanceOfWatchStore()}
}

func (q Querier) GetTransactionReceipt(hash common.Hash) (*TransactionReceipt, error) {
	var receipt TransactionReceipt
	b, e := q.store.Get([]byte(prefixReceipt + hash.String()))
	if e != nil {
		return nil, e
	}
	e = json.Unmarshal(b, &receipt)
	if e != nil {
		return nil, e
	}
	return &receipt, nil
}

func (q Querier) GetBlockByHash(hash common.Hash) (*EthBlock, error) {
	var block EthBlock
	b, e := q.store.Get([]byte(prefixBlock + hash.String()))
	if e != nil {
		return nil, e
	}
	e = json.Unmarshal(b, &block)
	if e != nil {
		return nil, e
	}
	return &block, nil
}

func (q Querier) GetBlockByNumber(number uint64) (*EthBlock, error) {
	hash, e := q.store.Get([]byte(prefixBlockInfo + strconv.Itoa(int(number))))
	if e != nil {
		return nil, e
	}
	return q.GetBlockByHash(common.HexToHash(string(hash)))
}

func (q Querier) GetCode(contractAddr common.Address, height uint64) ([]byte, error) {
	var codeInfo CodeInfo
	info, e := q.store.Get([]byte(prefixCode + contractAddr.String()))
	if e != nil {
		return nil, e
	}
	e = json.Unmarshal(info, &codeInfo)
	if e != nil {
		return nil, e
	}
	if height < codeInfo.Height && height > 0 {
		return nil, errors.New("the target height has not deployed this contract")
	}
	return hex.DecodeString(codeInfo.Code)
}

func (q Querier) GetLatestBlockNumber() (uint64, error) {
	height, e := q.store.Get([]byte(prefixLatestHeight + KeyLatestHeight))
	if e != nil {
		return 0, e
	}
	h, e := strconv.Atoi(string(height))
	return uint64(h), e
}

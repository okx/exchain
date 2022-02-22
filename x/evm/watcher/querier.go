package watcher

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	lru "github.com/hashicorp/golang-lru"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"github.com/okex/exchain/app/rpc/namespaces/eth/state"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	"github.com/okex/exchain/app/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

const MsgFunctionDisable = "fast query function has been disabled"

var errNotFound = errors.New("leveldb: not found")

type Querier struct {
	store *WatchStore
	sw    bool
	lru   *lru.Cache
}

func (q Querier) enabled() bool {
	return q.sw
}

func (q *Querier) Enable(sw bool) {
	q.sw = sw
}

func NewQuerier() *Querier {
	lru, e := lru.New(GetWatchLruSize())
	if e != nil {
		panic(errors.New("Failed to init LRU Cause " + e.Error()))
	}
	return &Querier{store: InstanceOfWatchStore(), sw: IsWatcherEnabled(), lru: lru}
}

func (q Querier) GetTransactionReceipt(hash common.Hash) (*TransactionReceipt, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	var receipt TransactionReceipt
	b, e := q.store.Get(append(prefixReceipt, hash.Bytes()...))
	if e != nil {
		return nil, e
	}
	if b == nil {
		return nil, errNotFound
	}
	e = json.Unmarshal(b, &receipt)
	if e != nil {
		return nil, e
	}
	if receipt.Logs == nil {
		receipt.Logs = []*ethtypes.Log{}
	}
	return &receipt, nil
}

func (q Querier) GetBlockByHash(hash common.Hash, fullTx bool) (*EthBlock, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	var block EthBlock
	b, e := q.store.Get(append(prefixBlock, hash.Bytes()...))
	if e != nil {
		return nil, e
	}
	if b == nil {
		return nil, errNotFound
	}

	e = json.Unmarshal(b, &block)
	if e != nil {
		return nil, e
	}
	if fullTx && block.Transactions != nil {
		txsHash := block.Transactions.([]interface{})
		txList := []rpctypes.Transaction{}
		if len(txsHash) == 0 {
			block.TransactionsRoot = ethtypes.EmptyRootHash
		}
		for _, tx := range txsHash {
			transaction, e := q.GetTransactionByHash(common.HexToHash(tx.(string)))
			if e == nil && transaction != nil {
				txList = append(txList, *transaction)
			}
		}
		block.Transactions = txList
	}
	block.UncleHash = ethtypes.EmptyUncleHash
	block.ReceiptsRoot = ethtypes.EmptyRootHash

	return &block, nil
}

func (q Querier) GetBlockHashByNumber(number uint64) (common.Hash, error) {
	if !q.enabled() {
		return common.Hash{}, errors.New(MsgFunctionDisable)
	}
	var height = number
	var err error
	if height == 0 {
		height, err = q.GetLatestBlockNumber()
		if err != nil {
			return common.Hash{}, err
		}
	}
	hash, e := q.store.Get(append(prefixBlockInfo, []byte(strconv.Itoa(int(height)))...))
	if e != nil {
		return common.Hash{}, e
	}
	if hash == nil {
		return common.Hash{}, errNotFound
	}
	return common.HexToHash(string(hash)), e
}

func (q Querier) GetBlockByNumber(number uint64, fullTx bool) (*EthBlock, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	var height = number
	var err error
	if height == 0 {
		height, err = q.GetLatestBlockNumber()
		if err != nil {
			return nil, err
		}
	}
	hash, e := q.store.Get(append(prefixBlockInfo, []byte(strconv.Itoa(int(height)))...))
	if e != nil {
		return nil, e
	}
	if hash == nil {
		return nil, errNotFound
	}

	return q.GetBlockByHash(common.HexToHash(string(hash)), fullTx)
}

func (q Querier) GetCode(contractAddr common.Address, height uint64) ([]byte, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	var codeInfo CodeInfo
	info, e := q.store.Get(append(prefixCode, contractAddr.Bytes()...))
	if e != nil {
		return nil, e
	}
	if info == nil {
		return nil, errNotFound
	}

	e = json.Unmarshal(info, &codeInfo)
	if e != nil {
		return nil, e
	}
	if height < codeInfo.Height && height > 0 {
		return nil, errors.New("the target height has not deploy this contract yet")
	}
	return hex.DecodeString(codeInfo.Code)
}

func (q Querier) GetCodeByHash(codeHash []byte) ([]byte, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	cacheCode, ok := q.lru.Get(common.BytesToHash(codeHash))
	if ok {
		data, ok := cacheCode.([]byte)
		if ok {
			return data, nil
		}
	}
	code, e := q.store.Get(append(prefixCodeHash, codeHash...))
	if e != nil {
		return nil, e
	}
	if code == nil {
		return nil, errNotFound
	}
	q.lru.Add(common.BytesToHash(codeHash), code)
	return code, nil
}

func (q Querier) GetLatestBlockNumber() (uint64, error) {
	if !q.enabled() {
		return 0, errors.New(MsgFunctionDisable)
	}
	height, e := q.store.Get(append(prefixLatestHeight, KeyLatestHeight...))
	if e != nil {
		return 0, e
	}
	if height == nil {
		return 0, errNotFound
	}
	h, e := strconv.Atoi(string(height))
	return uint64(h), e
}

func (q Querier) GetTransactionByHash(hash common.Hash) (*rpctypes.Transaction, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	var tx rpctypes.Transaction
	transaction, e := q.store.Get(append(prefixTx, hash.Bytes()...))
	if e != nil {
		return nil, e
	}
	if transaction == nil {
		return nil, errNotFound
	}
	e = json.Unmarshal(transaction, &tx)
	if e != nil {
		return nil, e
	}
	return &tx, nil
}

func (q Querier) GetTransactionByBlockNumberAndIndex(number uint64, idx uint) (*rpctypes.Transaction, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	block, e := q.GetBlockByNumber(number, true)
	if e != nil {
		return nil, e
	}
	return q.getTransactionByBlockAndIndex(block, idx)
}

func (q Querier) GetTransactionByBlockHashAndIndex(hash common.Hash, idx uint) (*rpctypes.Transaction, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	block, e := q.GetBlockByHash(hash, true)
	if e != nil {
		return nil, e
	}
	return q.getTransactionByBlockAndIndex(block, idx)
}

func (q Querier) getTransactionByBlockAndIndex(block *EthBlock, idx uint) (*rpctypes.Transaction, error) {
	if block.Transactions == nil {
		return nil, errors.New("no such transaction in target block")
	}
	txs := block.Transactions.([]rpctypes.Transaction)

	for _, tx := range txs {
		rawTx := tx
		if idx == uint(*rawTx.TransactionIndex) {
			return &rawTx, nil
		}
	}
	return nil, errors.New("no such transaction in target block")
}

func (q Querier) GetTransactionsByBlockNumber(number, offset, limit uint64) ([]*rpctypes.Transaction, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	block, err := q.GetBlockByNumber(number, true)
	if err != nil {
		return nil, err
	}
	if block.Transactions == nil {
		return nil, errors.New("no such transaction in target block")
	}

	rawTxs := block.Transactions.([]rpctypes.Transaction)
	var txs []*rpctypes.Transaction
	for idx := offset; idx < offset+limit && int(idx) < len(rawTxs); idx++ {
		rawTx := rawTxs[idx]
		txs = append(txs, &rawTx)
	}

	return txs, nil
}

func (q Querier) MustGetAccount(addr sdk.AccAddress) (*types.EthAccount, error) {
	acc, e := q.GetAccount(addr)
	//todo delete account from rdb if we get Account from H db successfully
	if e != nil {
		acc, e = q.GetAccountFromRdb(addr)
	} else {
		q.DeleteAccountFromRdb(addr)
	}
	return acc, e
}

func (q Querier) GetAccount(addr sdk.AccAddress) (*types.EthAccount, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	var acc types.EthAccount
	b, e := q.store.Get([]byte(GetMsgAccountKey(addr.Bytes())))
	if e != nil {
		return nil, e
	}
	if b == nil {
		return nil, errNotFound
	}
	e = json.Unmarshal(b, &acc)
	if e != nil {
		return nil, e
	}
	return &acc, nil
}

func (q Querier) GetAccountFromRdb(addr sdk.AccAddress) (*types.EthAccount, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	var acc types.EthAccount
	key := append(prefixRpcDb, GetMsgAccountKey(addr.Bytes())...)

	b, e := q.store.Get(key)
	if e != nil {
		return nil, e
	}
	if b == nil {
		return nil, errNotFound
	}
	e = json.Unmarshal(b, &acc)
	if e != nil {
		return nil, e
	}
	return &acc, nil
}

func (q Querier) DeleteAccountFromRdb(addr sdk.AccAddress) {
	if !q.enabled() {
		return
	}
	q.store.Delete(append(prefixRpcDb, GetMsgAccountKey(addr.Bytes())...))
}

func (q Querier) MustGetState(addr common.Address, key []byte) ([]byte, error) {
	orgKey := GetMsgStateKey(addr, key)
	realKey := common.BytesToHash(orgKey)
	data := state.GetStateFromLru(realKey)
	if data != nil {
		return data, nil
	}
	b, e := q.GetState(orgKey)
	if e != nil {
		b, e = q.GetStateFromRdb(orgKey)
	} else {
		q.DeleteStateFromRdb(addr, key)
	}
	if e == nil {
		state.SetStateToLru(realKey, b)
	}
	return b, e
}

func (q Querier) GetState(key []byte) ([]byte, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	b, e := q.store.Get(key)
	if e != nil {
		return nil, e
	}
	if b == nil {
		return nil, errNotFound
	}
	return b, nil
}

func (q Querier) GetStateFromRdb(key []byte) ([]byte, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	b, e := q.store.Get(append(prefixRpcDb, key...))
	if e != nil {
		return nil, e
	}
	if b == nil {
		return nil, errNotFound
	}

	return b, nil
}

func (q Querier) DeleteStateFromRdb(addr common.Address, key []byte) {
	if !q.enabled() {
		return
	}
	q.store.Delete(append(prefixRpcDb, GetMsgStateKey(addr, key)...))
}

func (q Querier) GetParams() (*evmtypes.Params, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	params := q.store.GetEvmParams()
	return &params, nil
}

func (q Querier) HasContractBlockedList(key []byte) bool {
	if !q.enabled() {
		return false
	}
	return q.store.Has(append(prefixBlackList, key...))
}
func (q Querier) GetContractMethodBlockedList(key []byte) ([]byte, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	return q.store.Get(append(prefixBlackList, key...))
}

func (q Querier) HasContractDeploymentWhitelist(key []byte) bool {
	if !q.enabled() {
		return false
	}
	return q.store.Has(append(prefixWhiteList, key...))
}

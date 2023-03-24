package watcher

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/gogo/protobuf/proto"
	lru "github.com/hashicorp/golang-lru"

	"github.com/okx/okbchain/app/rpc/namespaces/eth/state"
	"github.com/okx/okbchain/app/types"
	clientcontext "github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	evmtypes "github.com/okx/okbchain/x/evm/types"
	prototypes "github.com/okx/okbchain/x/evm/watcher/proto"
)

const MsgFunctionDisable = "fast query function has been disabled"

var errNotFound = errors.New("leveldb: not found")
var errDisable = errors.New(MsgFunctionDisable)

const hashPrefixKeyLen = 33

var hashPrefixKeyPool = &sync.Pool{
	New: func() interface{} {
		return &[hashPrefixKeyLen]byte{}
	},
}

func getHashPrefixKey(prefix []byte, hash []byte) ([]byte, error) {
	if len(prefix)+len(hash) > hashPrefixKeyLen {
		return nil, errors.New("invalid prefix or hash len")
	}
	key := hashPrefixKeyPool.Get().(*[hashPrefixKeyLen]byte)
	copy(key[:], prefix)
	copy(key[len(prefix):], hash)
	return key[:len(prefix)+len(hash)], nil
}

func putHashPrefixKey(key []byte) {
	hashPrefixKeyPool.Put((*[hashPrefixKeyLen]byte)(key[:hashPrefixKeyLen]))
}

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
		return nil, errDisable
	}
	var protoReceipt prototypes.TransactionReceipt
	b, e := q.store.Get(append(prefixReceipt, hash.Bytes()...))
	if e != nil {
		return nil, e
	}
	if b == nil {
		return nil, errNotFound
	}
	e = proto.Unmarshal(b, &protoReceipt)
	if e != nil {
		return nil, e
	}
	receipt := protoToReceipt(&protoReceipt)
	return receipt, nil
}

func (q Querier) GetTransactionResponse(hash common.Hash) (*TransactionResponse, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	var response TransactionResponse
	b, e := q.store.Get(append(prefixTxResponse, hash.Bytes()...))
	if e != nil {
		return nil, e
	}
	if b == nil {
		return nil, errNotFound
	}
	e = json.Unmarshal(b, &response)
	if e != nil {
		return nil, e
	}

	return &response, nil
}

func (q Querier) GetBlockByHash(hash common.Hash, fullTx bool) (*evmtypes.Block, error) {
	if !q.enabled() {
		return nil, errDisable
	}
	var block evmtypes.Block
	var err error
	var blockHashKey []byte
	if blockHashKey, err = getHashPrefixKey(prefixBlock, hash.Bytes()); err != nil {
		blockHashKey = append(prefixBlock, hash.Bytes()...)
	} else {
		defer putHashPrefixKey(blockHashKey)
	}

	_, err = q.store.GetUnsafe(blockHashKey, func(value []byte) (interface{}, error) {
		if value == nil {
			return nil, errNotFound
		}
		e := json.Unmarshal(value, &block)
		if e != nil {
			return nil, e
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	if fullTx && block.Transactions != nil {
		txsHash := block.Transactions.([]interface{})
		txList := make([]*Transaction, 0, len(txsHash))
		for _, tx := range txsHash {
			transaction, e := q.GetTransactionByHash(common.HexToHash(tx.(string)))
			if e == nil && transaction != nil {
				txList = append(txList, transaction)
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
		return common.Hash{}, errDisable
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

func (q Querier) GetBlockByNumber(number uint64, fullTx bool) (*evmtypes.Block, error) {
	if !q.enabled() {
		return nil, errDisable
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
		return nil, errDisable
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
		return nil, errDisable
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
		return 0, errDisable
	}
	height, e := q.store.Get(keyLatestBlockHeight)
	if e != nil {
		return 0, e
	}
	if height == nil {
		return 0, errNotFound
	}
	h, e := strconv.Atoi(string(height))
	return uint64(h), e
}

func (q Querier) GetTransactionByHash(hash common.Hash) (*Transaction, error) {
	if !q.enabled() {
		return nil, errDisable
	}
	var protoTx prototypes.Transaction
	var txHashKey []byte
	var err error
	if txHashKey, err = getHashPrefixKey(prefixTx, hash.Bytes()); err != nil {
		txHashKey = append(prefixTx, hash.Bytes()...)
	} else {
		defer putHashPrefixKey(txHashKey)
	}

	_, err = q.store.GetUnsafe(txHashKey, func(value []byte) (interface{}, error) {
		if value == nil {
			return nil, errNotFound
		}
		e := proto.Unmarshal(value, &protoTx)
		if e != nil {
			return nil, e
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	tx := protoToTransaction(&protoTx)
	return tx, nil
}

func (q Querier) GetTransactionByBlockNumberAndIndex(number uint64, idx uint) (*Transaction, error) {
	if !q.enabled() {
		return nil, errDisable
	}
	block, e := q.GetBlockByNumber(number, true)
	if e != nil {
		return nil, e
	}
	return q.getTransactionByBlockAndIndex(block, idx)
}

func (q Querier) GetTransactionByBlockHashAndIndex(hash common.Hash, idx uint) (*Transaction, error) {
	if !q.enabled() {
		return nil, errDisable
	}
	block, e := q.GetBlockByHash(hash, true)
	if e != nil {
		return nil, e
	}
	return q.getTransactionByBlockAndIndex(block, idx)
}

func (q Querier) getTransactionByBlockAndIndex(block *evmtypes.Block, idx uint) (*Transaction, error) {
	if block.Transactions == nil {
		return nil, errors.New("no such transaction in target block")
	}
	txs, ok := block.Transactions.([]*Transaction)
	if ok {
		for _, tx := range txs {
			rawTx := *tx
			if idx == uint(*rawTx.TransactionIndex) {
				return &rawTx, nil
			}
		}
	}
	return nil, errors.New("no such transaction in target block")
}

func (q Querier) GetTransactionsByBlockNumber(number, offset, limit uint64) ([]*Transaction, error) {
	if !q.enabled() {
		return nil, errDisable
	}
	block, err := q.GetBlockByNumber(number, true)
	if err != nil {
		return nil, err
	}
	if block.Transactions == nil {
		return nil, errors.New("no such transaction in target block")
	}

	rawTxs, ok := block.Transactions.([]*Transaction)
	if ok {
		var txs []*Transaction
		for idx := offset; idx < offset+limit && int(idx) < len(rawTxs); idx++ {
			rawTx := *rawTxs[idx]
			txs = append(txs, &rawTx)
		}
		return txs, nil
	}
	return nil, errors.New("no such transaction in target block")
}

func (q Querier) GetTxResultByBlock(clientCtx clientcontext.CLIContext,
	height, offset, limit uint64) ([]*TransactionResult, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}

	// get block hash
	rawBlockHash, err := q.store.Get(append(prefixBlockInfo, []byte(strconv.Itoa(int(height)))...))
	if err != nil {
		return nil, err
	}
	if rawBlockHash == nil {
		return nil, errNotFound
	}

	blockHash := common.HexToHash(string(rawBlockHash))

	// get block by hash
	var block evmtypes.Block
	var blockHashKey []byte
	if blockHashKey, err = getHashPrefixKey(prefixBlock, blockHash.Bytes()); err != nil {
		blockHashKey = append(prefixBlock, blockHash.Bytes()...)
	} else {
		defer putHashPrefixKey(blockHashKey)
	}

	_, err = q.store.GetUnsafe(blockHashKey, func(value []byte) (interface{}, error) {
		if value == nil {
			return nil, errNotFound
		}
		e := json.Unmarshal(value, &block)
		if e != nil {
			return nil, e
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	results := make([]*TransactionResult, 0, limit)
	var ethStart, ethEnd, ethTxLen uint64

	// get result from eth tx
	if block.Transactions != nil {
		txsHash := block.Transactions.([]interface{})
		ethTxLen = uint64(len(txsHash))

		if offset < ethTxLen {
			ethStart = offset
			if ethEnd = ethStart + limit; ethEnd > ethTxLen {
				ethEnd = ethTxLen
			}
		}

		for i := ethStart; i < ethEnd; i++ {
			txHash := common.HexToHash(txsHash[i].(string))
			//Get Eth Tx
			tx, err := q.GetTransactionByHash(txHash)
			if err != nil {
				return nil, err
			}
			//Get Eth Receipt
			receipt, err := q.GetTransactionReceipt(txHash)
			if err != nil {
				return nil, err
			}

			// Get tx Response
			var txLog string
			txResult, err := q.GetTransactionResponse(txHash)
			if err == nil {
				txLog = txResult.TxResult.Log
			}

			r := &TransactionResult{TxType: hexutil.Uint64(EthReceipt), EthTx: tx, Receipt: receipt,
				EthTxLog: txLog}
			results = append(results, r)
		}
	}

	// enough Tx by Eth
	ethTxNums := ethEnd - ethStart
	if ethTxNums == limit {
		return results, nil
	}
	// calc remain std txs
	remainTxs := limit - ethTxNums
	// get result from Std tx
	var stdTxsHash []common.Hash
	b, err := q.store.Get(append(prefixStdTxHash, blockHash.Bytes()...))
	if err != nil {
		return nil, err
	}

	if b == nil {
		return results, nil

	}
	err = json.Unmarshal(b, &stdTxsHash)
	if err != nil {
		return nil, err
	}

	if stdTxsHash != nil && len(stdTxsHash) != 0 {
		stdTxsLen := uint64(len(stdTxsHash))
		var stdStart, stdEnd uint64
		stdStart = offset + ethTxNums - ethTxLen
		if stdEnd = stdStart + remainTxs; stdEnd > stdTxsLen {
			stdEnd = stdTxsLen
		}

		for i := stdStart; i < stdEnd; i++ {
			stdResponse, err := q.GetTransactionResponse(stdTxsHash[i])
			if err != nil {
				return nil, err
			}

			res, err := RawTxResultToStdResponse(clientCtx, stdResponse.ResultTx, nil, stdResponse.Timestamp)
			if err != nil {
				return nil, err
			}
			results = append(results, res)
		}
	}

	return results, nil
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
		return nil, errDisable
	}
	b, e := q.store.Get([]byte(GetMsgAccountKey(addr.Bytes())))
	if e != nil {
		return nil, e
	}
	if b == nil {
		return nil, errNotFound
	}
	acc, err := DecodeAccount(b)
	if err != nil {
		return nil, e
	}
	return acc, nil
}

func (q Querier) GetAccountFromRdb(addr sdk.AccAddress) (*types.EthAccount, error) {
	if !q.enabled() {
		return nil, errDisable
	}
	key := append(prefixRpcDb, GetMsgAccountKey(addr.Bytes())...)

	b, e := q.store.Get(key)
	if e != nil {
		return nil, e
	}
	if b == nil {
		return nil, errNotFound
	}
	acc, err := DecodeAccount(b)
	if err != nil {
		return nil, e
	}
	return acc, nil
}

func (q Querier) DeleteAccountFromRdb(addr sdk.AccAddress) {
	if !q.enabled() {
		return
	}
	q.store.Delete(append(prefixRpcDb, GetMsgAccountKey(addr.Bytes())...))
}

func (q Querier) MustGetState(addr common.Address, key []byte) ([]byte, error) {
	orgKey := GetMsgStateKey(addr, key)
	data := state.GetStateFromLru(orgKey)
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
		state.SetStateToLru(orgKey, b)
	}
	return b, e
}

func (q Querier) GetState(key []byte) ([]byte, error) {
	if !q.enabled() {
		return nil, errDisable
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
		return nil, errDisable
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
		return nil, errDisable
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
		return nil, errDisable
	}
	return q.store.Get(append(prefixBlackList, key...))
}

func (q Querier) HasContractDeploymentWhitelist(key []byte) bool {
	if !q.enabled() {
		return false
	}
	return q.store.Has(append(prefixWhiteList, key...))
}

func (q Querier) GetStdTxHashByBlockHash(hash common.Hash) ([]common.Hash, error) {
	if !q.enabled() {
		return nil, errors.New(MsgFunctionDisable)
	}
	var stdTxHash []common.Hash
	b, e := q.store.Get(append(prefixStdTxHash, hash.Bytes()...))
	if e != nil {
		return nil, e
	}
	if b == nil {
		return nil, errNotFound
	}
	e = json.Unmarshal(b, &stdTxHash)
	if e != nil {
		return nil, e
	}

	return stdTxHash, nil
}

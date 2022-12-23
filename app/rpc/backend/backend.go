package backend

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"

	"github.com/okex/exchain/libs/tendermint/global"

	lru "github.com/hashicorp/golang-lru"
	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"

	"github.com/spf13/viper"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/evm/watcher"
	"golang.org/x/time/rate"

	rpctypes "github.com/okex/exchain/app/rpc/types"
	evmtypes "github.com/okex/exchain/x/evm/types"

	clientcontext "github.com/okex/exchain/libs/cosmos-sdk/client/context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/bitutil"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/bloombits"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"
)

const (
	FlagLogsLimit   = "rpc.logs-limit"
	FlagLogsTimeout = "rpc.logs-timeout"
	blockCacheSize  = 1024
)

var ErrTimeout = errors.New("query timeout exceeded")

// Backend implements the functionality needed to filter changes.
// Implemented by EthermintBackend.
type Backend interface {
	// Used by block filter; also used for polling
	BlockNumber() (hexutil.Uint64, error)
	LatestBlockNumber() (int64, error)
	HeaderByNumber(blockNum rpctypes.BlockNumber) (*ethtypes.Header, error)
	HeaderByHash(blockHash common.Hash) (*ethtypes.Header, error)
	GetBlockByNumber(blockNum rpctypes.BlockNumber, fullTx bool) (*watcher.Block, error)
	GetBlockByHash(hash common.Hash, fullTx bool) (*watcher.Block, error)

	GetTransactionByHash(hash common.Hash) (*watcher.Transaction, error)

	// returns the logs of a given block
	GetLogs(height int64) ([][]*ethtypes.Log, error)

	// Used by pending transaction filter
	PendingTransactions() ([]*watcher.Transaction, error)
	PendingTransactionCnt() (int, error)
	PendingTransactionsByHash(target common.Hash) (*watcher.Transaction, error)
	UserPendingTransactionsCnt(address string) (int, error)
	UserPendingTransactions(address string, limit int) ([]*watcher.Transaction, error)
	PendingAddressList() ([]string, error)
	GetPendingNonce(address string) (uint64, bool)

	// Used by log filter
	GetTransactionLogs(txHash common.Hash) ([]*ethtypes.Log, error)
	BloomStatus() (uint64, uint64)
	ServiceFilter(ctx context.Context, session *bloombits.MatcherSession)

	// Used by eip-1898
	ConvertToBlockNumber(rpctypes.BlockNumberOrHash) (rpctypes.BlockNumber, error)
	// Block returns the block at the given block number, block data is readonly
	Block(height *int64) (*coretypes.ResultBlock, error)
	PruneEverything() bool
}

var _ Backend = (*EthermintBackend)(nil)

// EthermintBackend implements the Backend interface
type EthermintBackend struct {
	ctx               context.Context
	clientCtx         clientcontext.CLIContext
	logger            log.Logger
	gasLimit          int64
	bloomRequests     chan chan *bloombits.Retrieval
	closeBloomHandler chan struct{}
	wrappedBackend    *watcher.Querier
	rateLimiters      map[string]*rate.Limiter
	disableAPI        map[string]bool
	backendCache      Cache
	logsLimit         int
	logsTimeout       int // timeout second
	blockCache        *lru.Cache
	pruneEverything   bool
}

// New creates a new EthermintBackend instance
func New(clientCtx clientcontext.CLIContext, log log.Logger, rateLimiters map[string]*rate.Limiter, disableAPI map[string]bool) *EthermintBackend {
	b := &EthermintBackend{
		ctx:               context.Background(),
		clientCtx:         clientCtx,
		logger:            log.With("module", "json-rpc"),
		gasLimit:          int64(^uint32(0)),
		bloomRequests:     make(chan chan *bloombits.Retrieval),
		closeBloomHandler: make(chan struct{}),
		wrappedBackend:    watcher.NewQuerier(),
		rateLimiters:      rateLimiters,
		disableAPI:        disableAPI,
		backendCache:      NewLruCache(),
		logsLimit:         viper.GetInt(FlagLogsLimit),
		logsTimeout:       viper.GetInt(FlagLogsTimeout),
		pruneEverything:   viper.GetString(server.FlagPruning) == types.PruningOptionEverything,
	}
	b.blockCache, _ = lru.New(blockCacheSize)
	return b
}

func (b *EthermintBackend) PruneEverything() bool {
	return b.pruneEverything
}

func (b *EthermintBackend) LogsLimit() int {
	return b.logsLimit
}

func (b *EthermintBackend) LogsTimeout() time.Duration {
	return time.Duration(b.logsTimeout) * time.Second
}

// BlockNumber returns the current block number.
func (b *EthermintBackend) BlockNumber() (hexutil.Uint64, error) {
	committedHeight := global.GetGlobalHeight()
	return hexutil.Uint64(committedHeight), nil
}

// GetBlockByNumber returns the block identified by number.
func (b *EthermintBackend) GetBlockByNumber(blockNum rpctypes.BlockNumber, fullTx bool) (*watcher.Block, error) {
	//query block in cache first
	block, err := b.backendCache.GetBlockByNumber(uint64(blockNum), fullTx)
	if err == nil {
		return block, nil
	}
	//query block from watch db
	block, err = b.wrappedBackend.GetBlockByNumber(uint64(blockNum), fullTx)
	if err == nil {
		b.backendCache.AddOrUpdateBlock(block.Hash, block, fullTx)
		return block, nil
	}
	//query block from db
	height := blockNum.Int64()
	if height <= 0 {
		// get latest block height
		num, err := b.BlockNumber()
		if err != nil {
			return nil, err
		}
		height = int64(num)
	}

	resBlock, err := b.Block(&height)
	if err != nil {
		return nil, nil
	}

	block, err = rpctypes.RpcBlockFromTendermint(b.clientCtx, resBlock.Block, fullTx)
	if err != nil {
		return nil, err
	}
	b.backendCache.AddOrUpdateBlock(block.Hash, block, fullTx)
	return block, nil
}

// GetBlockByHash returns the block identified by hash.
func (b *EthermintBackend) GetBlockByHash(hash common.Hash, fullTx bool) (*watcher.Block, error) {
	//query block in cache first
	block, err := b.backendCache.GetBlockByHash(hash, fullTx)
	if err == nil {
		return block, err
	}
	//query block from watch db
	block, err = b.wrappedBackend.GetBlockByHash(hash, fullTx)
	if err == nil {
		b.backendCache.AddOrUpdateBlock(hash, block, fullTx)
		return block, nil
	}
	//query block from tendermint
	res, _, err := b.clientCtx.Query(fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryHashToHeight, hash.Hex()))
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResBlockNumber
	if err := b.clientCtx.Codec.UnmarshalJSON(res, &out); err != nil {
		return nil, err
	}

	resBlock, err := b.Block(&out.Number)
	if err != nil {
		return nil, nil
	}

	block, err = rpctypes.RpcBlockFromTendermint(b.clientCtx, resBlock.Block, fullTx)
	if err != nil {
		return nil, err
	}
	b.backendCache.AddOrUpdateBlock(hash, block, fullTx)
	return block, nil
}

// HeaderByNumber returns the block header identified by height.
func (b *EthermintBackend) HeaderByNumber(blockNum rpctypes.BlockNumber) (*ethtypes.Header, error) {
	height := blockNum.Int64()
	if height <= 0 {
		// get latest block height
		num, err := b.BlockNumber()
		if err != nil {
			return nil, err
		}

		height = int64(num)
	}

	resBlock, err := b.Block(&height)
	if err != nil {
		return nil, err
	}

	res, _, err := b.clientCtx.Query(fmt.Sprintf("custom/%s/%s/%d", evmtypes.ModuleName, evmtypes.QueryBloom, resBlock.Block.Height))
	if err != nil {
		return nil, err
	}

	var bloomRes evmtypes.QueryBloomFilter
	b.clientCtx.Codec.MustUnmarshalJSON(res, &bloomRes)

	ethHeader := rpctypes.EthHeaderFromTendermint(resBlock.Block.Header)
	ethHeader.Bloom = bloomRes.Bloom
	return ethHeader, nil
}

// HeaderByHash returns the block header identified by hash.
func (b *EthermintBackend) HeaderByHash(blockHash common.Hash) (*ethtypes.Header, error) {
	res, _, err := b.clientCtx.Query(fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryHashToHeight, blockHash.Hex()))
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResBlockNumber
	if err := b.clientCtx.Codec.UnmarshalJSON(res, &out); err != nil {
		return nil, err
	}

	resBlock, err := b.Block(&out.Number)
	if err != nil {
		return nil, err
	}

	res, _, err = b.clientCtx.Query(fmt.Sprintf("custom/%s/%s/%d", evmtypes.ModuleName, evmtypes.QueryBloom, resBlock.Block.Height))
	if err != nil {
		return nil, err
	}

	var bloomRes evmtypes.QueryBloomFilter
	b.clientCtx.Codec.MustUnmarshalJSON(res, &bloomRes)

	ethHeader := rpctypes.EthHeaderFromTendermint(resBlock.Block.Header)
	ethHeader.Bloom = bloomRes.Bloom
	return ethHeader, nil
}

// GetTransactionLogs returns the logs given a transaction hash.
// It returns an error if there's an encoding error.
// If no logs are found for the tx hash, the error is nil.
func (b *EthermintBackend) GetTransactionLogs(txHash common.Hash) ([]*ethtypes.Log, error) {
	txRes, err := b.clientCtx.Client.Tx(txHash.Bytes(), !b.clientCtx.TrustNode)
	if err != nil {
		return nil, err
	}

	execRes, err := evmtypes.DecodeResultData(txRes.TxResult.Data)
	if err != nil {
		return nil, err
	}

	return execRes.Logs, nil
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (b *EthermintBackend) PendingTransactions() ([]*watcher.Transaction, error) {
	lastHeight, err := b.clientCtx.Client.LatestBlockNumber()
	if err != nil {
		return nil, err
	}
	pendingTxs, err := b.clientCtx.Client.UnconfirmedTxs(-1)
	if err != nil {
		return nil, err
	}

	transactions := make([]*watcher.Transaction, 0, len(pendingTxs.Txs))
	for _, tx := range pendingTxs.Txs {
		ethTx, err := rpctypes.RawTxToEthTx(b.clientCtx, tx, lastHeight)
		if err != nil {
			// ignore non Ethermint EVM transactions
			continue
		}

		// TODO: check signer and reference against accounts the node manages
		rpcTx, err := watcher.NewTransaction(ethTx, common.BytesToHash(ethTx.Hash), common.Hash{}, 0, 0)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, rpcTx)
	}

	return transactions, nil
}

func (b *EthermintBackend) PendingTransactionCnt() (int, error) {
	result, err := b.clientCtx.Client.NumUnconfirmedTxs()
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

func (b *EthermintBackend) UserPendingTransactionsCnt(address string) (int, error) {
	result, err := b.clientCtx.Client.UserNumUnconfirmedTxs(address)
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

func (b *EthermintBackend) GetPendingNonce(address string) (uint64, bool) {
	result, ok := b.clientCtx.Client.GetPendingNonce(address)
	if !ok {
		return 0, false
	}
	return result.Nonce, true
}

func (b *EthermintBackend) UserPendingTransactions(address string, limit int) ([]*watcher.Transaction, error) {
	lastHeight, err := b.clientCtx.Client.LatestBlockNumber()
	if err != nil {
		return nil, err
	}
	result, err := b.clientCtx.Client.UserUnconfirmedTxs(address, limit)
	if err != nil {
		return nil, err
	}
	transactions := make([]*watcher.Transaction, 0, len(result.Txs))
	for _, tx := range result.Txs {
		ethTx, err := rpctypes.RawTxToEthTx(b.clientCtx, tx, lastHeight)
		if err != nil {
			// ignore non Ethermint EVM transactions
			continue
		}

		// TODO: check signer and reference against accounts the node manages
		rpcTx, err := watcher.NewTransaction(ethTx, common.BytesToHash(ethTx.Hash), common.Hash{}, 0, 0)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, rpcTx)
	}

	return transactions, nil
}

func (b *EthermintBackend) PendingAddressList() ([]string, error) {
	res, err := b.clientCtx.Client.GetAddressList()
	if err != nil {
		return nil, err
	}
	return res.Addresses, nil
}

// PendingTransactions returns the transaction that is in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (b *EthermintBackend) PendingTransactionsByHash(target common.Hash) (*watcher.Transaction, error) {
	lastHeight, err := b.clientCtx.Client.LatestBlockNumber()
	if err != nil {
		return nil, err
	}
	pendingTx, err := b.clientCtx.Client.GetUnconfirmedTxByHash(target)
	if err != nil {
		return nil, err
	}
	ethTx, err := rpctypes.RawTxToEthTx(b.clientCtx, pendingTx, lastHeight)
	if err != nil {
		// ignore non Ethermint EVM transactions
		return nil, err
	}
	rpcTx, err := watcher.NewTransaction(ethTx, common.BytesToHash(ethTx.Hash), common.Hash{}, 0, 0)
	if err != nil {
		return nil, err
	}
	return rpcTx, nil
}

func (b *EthermintBackend) GetTransactionByHash(hash common.Hash) (tx *watcher.Transaction, err error) {
	// query tx in cache first
	tx, err = b.backendCache.GetTransaction(hash)
	if err == nil {
		return tx, err
	}
	// query tx in watch db
	tx, err = b.wrappedBackend.GetTransactionByHash(hash)
	if err == nil {
		b.backendCache.AddOrUpdateTransaction(hash, tx)
		return tx, nil
	}
	// query tx in tendermint
	txRes, err := b.clientCtx.Client.Tx(hash.Bytes(), false)
	if err != nil {
		return nil, err
	}

	// Can either cache or just leave this out if not necessary
	block, err := b.Block(&txRes.Height)
	if err != nil {
		return nil, err
	}

	blockHash := common.BytesToHash(block.Block.Hash())

	ethTx, err := rpctypes.RawTxToEthTx(b.clientCtx, txRes.Tx, txRes.Height)
	if err != nil {
		return nil, err
	}

	tx, err = watcher.NewTransaction(ethTx, common.BytesToHash(ethTx.Hash), blockHash, uint64(txRes.Height), uint64(txRes.Index))
	if err != nil {
		return nil, err
	}
	b.backendCache.AddOrUpdateTransaction(hash, tx)
	return tx, nil
}

// GetLogs returns all the logs from all the ethereum transactions in a block.
func (b *EthermintBackend) GetLogs(height int64) ([][]*ethtypes.Log, error) {
	block, err := b.Block(&height)
	if err != nil {
		return nil, err
	}
	// return empty directly when block was produced during stress testing.
	var blockLogs = [][]*ethtypes.Log{}
	if b.logsLimit > 0 && len(block.Block.Txs) > b.logsLimit {
		return blockLogs, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(b.logsTimeout)*time.Second)
	defer cancel()
	for _, tx := range block.Block.Txs {
		select {
		case <-ctx.Done():
			return nil, ErrTimeout
		default:
			// NOTE: we query the state in case the tx result logs are not persisted after an upgrade.
			txRes, err := b.clientCtx.Client.Tx(tx.Hash(block.Block.Height), !b.clientCtx.TrustNode)
			if err != nil {
				continue
			}
			execRes, err := evmtypes.DecodeResultData(txRes.TxResult.Data)
			if err != nil {
				continue
			}
			var validLogs []*ethtypes.Log
			for _, log := range execRes.Logs {
				if int64(log.BlockNumber) == block.Block.Height {
					validLogs = append(validLogs, log)
				}
			}
			blockLogs = append(blockLogs, validLogs)
		}
	}

	return blockLogs, nil
}

// BloomStatus returns the BloomBitsBlocks and the number of processed sections maintained
// by the chain indexer.
func (b *EthermintBackend) BloomStatus() (uint64, uint64) {
	sections := evmtypes.GetIndexer().StoredSection()
	return evmtypes.BloomBitsBlocks, sections
}

// LatestBlockNumber gets the latest block height in int64 format.
func (b *EthermintBackend) LatestBlockNumber() (int64, error) {
	return b.clientCtx.Client.LatestBlockNumber()
}

func (b *EthermintBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < evmtypes.BloomFilterThreads; i++ {
		go session.Multiplex(evmtypes.BloomRetrievalBatch, evmtypes.BloomRetrievalWait, b.bloomRequests)
	}
}

// startBloomHandlers starts a batch of goroutines to accept bloom bit database
// retrievals from possibly a range of filters and serving the data to satisfy.
func (b *EthermintBackend) StartBloomHandlers(sectionSize uint64, db dbm.DB) {
	for i := 0; i < evmtypes.BloomServiceThreads; i++ {
		go func() {
			for {
				select {
				case <-b.closeBloomHandler:
					return

				case request := <-b.bloomRequests:
					task := <-request
					task.Bitsets = make([][]byte, len(task.Sections))
					for i, section := range task.Sections {
						height := int64((section+1)*sectionSize-1) + tmtypes.GetStartBlockHeight()
						hash, err := b.GetBlockHashByHeight(rpctypes.BlockNumber(height))
						if err != nil {
							task.Error = err
						}
						if compVector, err := evmtypes.ReadBloomBits(db, task.Bit, section, hash); err == nil {
							if blob, err := bitutil.DecompressBytes(compVector, int(sectionSize/8)); err == nil {
								task.Bitsets[i] = blob
							} else {
								task.Error = err
							}
						} else {
							task.Error = err
						}
					}
					request <- task
				}
			}
		}()
	}
}

// GetBlockHashByHeight returns the block hash by height.
func (b *EthermintBackend) GetBlockHashByHeight(height rpctypes.BlockNumber) (common.Hash, error) {
	res, _, err := b.clientCtx.Query(fmt.Sprintf("custom/%s/%s/%d", evmtypes.ModuleName, evmtypes.QueryHeightToHash, height))
	if err != nil {
		return common.Hash{}, err
	}

	hash := common.BytesToHash(res)
	return hash, nil
}

// Close
func (b *EthermintBackend) Close() {
	close(b.closeBloomHandler)
}

func (b *EthermintBackend) GetRateLimiter(apiName string) *rate.Limiter {
	if b.rateLimiters == nil {
		return nil
	}
	return b.rateLimiters[apiName]
}

func (b *EthermintBackend) IsDisabled(apiName string) bool {
	if b.disableAPI == nil {
		return false
	}
	return b.disableAPI[apiName]
}

func (b *EthermintBackend) ConvertToBlockNumber(blockNumberOrHash rpctypes.BlockNumberOrHash) (rpctypes.BlockNumber, error) {
	if blockNumber, ok := blockNumberOrHash.Number(); ok {
		return blockNumber, nil
	}
	hash, ok := blockNumberOrHash.Hash()
	if !ok {
		return rpctypes.LatestBlockNumber, nil
	}
	ethBlock, err := b.wrappedBackend.GetBlockByHash(hash, false)
	if err == nil {
		return rpctypes.BlockNumber(ethBlock.Number), nil
	}

	res, _, err := b.clientCtx.Query(fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryHashToHeight, hash.Hex()))
	if err != nil {
		return rpctypes.LatestBlockNumber, rpctypes.ErrResourceNotFound
	}

	var out evmtypes.QueryResBlockNumber
	if err := b.clientCtx.Codec.UnmarshalJSON(res, &out); err != nil {
		return rpctypes.LatestBlockNumber, rpctypes.ErrResourceNotFound
	}
	return rpctypes.BlockNumber(out.Number), nil
}

func (b *EthermintBackend) cacheBlock(block *coretypes.ResultBlock) {
	if b.blockCache != nil {
		b.blockCache.Add(block.Block.Height, block)
	}
}

func (b *EthermintBackend) getBlockFromCache(height int64) *coretypes.ResultBlock {
	if b.blockCache != nil {
		if v, ok := b.blockCache.Get(height); ok {
			return v.(*coretypes.ResultBlock)
		}
	}
	return nil
}

func (b *EthermintBackend) Block(height *int64) (block *coretypes.ResultBlock, err error) {
	if height != nil {
		block = b.getBlockFromCache(*height)
	}
	if block == nil {
		block, err = b.clientCtx.Client.Block(height)
		if err != nil {
			return nil, err
		}
		b.cacheBlock(block)
	}
	return block, nil
}

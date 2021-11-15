package eth

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/spf13/viper"

	"github.com/okex/exchain/app"
	"github.com/okex/exchain/app/config"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/app/crypto/hd"
	"github.com/okex/exchain/app/rpc/backend"
	"github.com/okex/exchain/app/rpc/monitor"
	"github.com/okex/exchain/app/rpc/namespaces/eth/simulation"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	ethermint "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/app/utils"
	clientcontext "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	cmserver "github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authclient "github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/merkle"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
)

const (
	CacheOfEthCallLru = 40960

	FlagEnableMultiCall = "rpc.enable-multi-call"
)

// PublicEthereumAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicEthereumAPI struct {
	ctx            context.Context
	clientCtx      clientcontext.CLIContext
	chainIDEpoch   *big.Int
	logger         log.Logger
	backend        backend.Backend
	keys           []ethsecp256k1.PrivKey // unlocked keys
	nonceLock      *rpctypes.AddrLocker
	keyringLock    sync.Mutex
	gasPrice       *hexutil.Big
	wrappedBackend *watcher.Querier
	watcherBackend *watcher.Watcher
	evmFactory     simulation.EvmFactory
	txPool         *TxPool
	Metrics        map[string]*monitor.RpcMetrics
	callCache      *lru.Cache
}

// NewAPI creates an instance of the public ETH Web3 API.
func NewAPI(
	clientCtx clientcontext.CLIContext, log log.Logger, backend backend.Backend, nonceLock *rpctypes.AddrLocker,
	keys ...ethsecp256k1.PrivKey,
) *PublicEthereumAPI {

	epoch, err := ethermint.ParseChainID(clientCtx.ChainID)
	if err != nil {
		panic(err)
	}

	api := &PublicEthereumAPI{
		ctx:            context.Background(),
		clientCtx:      clientCtx,
		chainIDEpoch:   epoch,
		logger:         log.With("module", "json-rpc", "namespace", "eth"),
		backend:        backend,
		keys:           keys,
		nonceLock:      nonceLock,
		gasPrice:       ParseGasPrice(),
		wrappedBackend: watcher.NewQuerier(),
		watcherBackend: watcher.NewWatcher(),
	}
	api.evmFactory = simulation.NewEvmFactory(clientCtx.ChainID, api.wrappedBackend)

	if watcher.IsWatcherEnabled() {
		callCache, err := lru.New(CacheOfEthCallLru)
		if err != nil {
			panic(err)
		}
		api.callCache = callCache
	}

	if err := api.GetKeyringInfo(); err != nil {
		api.logger.Error("failed to get keybase info", "error", err)
	}

	if viper.GetBool(FlagEnableTxPool) {
		api.txPool = NewTxPool(clientCtx, api)
		go api.txPool.broadcastPeriod(api)
	}

	return api
}

// GetKeyringInfo checks if the keyring is present on the client context. If not, it creates a new
// instance and sets it to the client context for later usage.
func (api *PublicEthereumAPI) GetKeyringInfo() error {
	api.keyringLock.Lock()
	defer api.keyringLock.Unlock()

	if api.clientCtx.Keybase != nil {
		return nil
	}

	keybase, err := keys.NewKeyring(
		sdk.KeyringServiceName(),
		viper.GetString(flags.FlagKeyringBackend),
		viper.GetString(cmserver.FlagUlockKeyHome),
		api.clientCtx.Input,
		hd.EthSecp256k1Options()...,
	)
	if err != nil {
		return err
	}

	api.clientCtx.Keybase = keybase
	return nil
}

// ClientCtx returns the Cosmos SDK client context.
func (api *PublicEthereumAPI) ClientCtx() clientcontext.CLIContext {
	return api.clientCtx
}

// GetKeys returns the Cosmos SDK client context.
func (api *PublicEthereumAPI) GetKeys() []ethsecp256k1.PrivKey {
	return api.keys
}

// SetKeys sets the given key slice to the set of private keys
func (api *PublicEthereumAPI) SetKeys(keys []ethsecp256k1.PrivKey) {
	api.keys = keys
}

// ProtocolVersion returns the supported Ethereum protocol version.
func (api *PublicEthereumAPI) ProtocolVersion() hexutil.Uint {
	monitor := monitor.GetMonitor("eth_protocolVersion", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd()
	return hexutil.Uint(ethermint.ProtocolVersion)
}

// ChainId returns the chain's identifier in hex format
func (api *PublicEthereumAPI) ChainId() (hexutil.Uint, error) { // nolint
	api.logger.Debug("eth_chainId")
	return hexutil.Uint(uint(api.chainIDEpoch.Uint64())), nil
}

// Syncing returns whether or not the current node is syncing with other peers. Returns false if not, or a struct
// outlining the state of the sync if it is.
func (api *PublicEthereumAPI) Syncing() (interface{}, error) {
	monitor := monitor.GetMonitor("eth_syncing", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd()
	status, err := api.clientCtx.Client.Status()
	if err != nil {
		return false, err
	}

	if !status.SyncInfo.CatchingUp {
		return false, nil
	}

	return map[string]interface{}{
		"startingBlock": hexutil.Uint64(status.SyncInfo.EarliestBlockHeight),
		"currentBlock":  hexutil.Uint64(status.SyncInfo.LatestBlockHeight),
		"highestBlock":  hexutil.Uint64(0), // NA
		// "pulledStates":  nil, // NA
		// "knownStates":   nil, // NA
	}, nil
}

// Coinbase is the address that staking rewards will be send to (alias for Etherbase).
func (api *PublicEthereumAPI) Coinbase() (common.Address, error) {
	api.logger.Debug("eth_coinbase")

	node, err := api.clientCtx.GetNode()
	if err != nil {
		return common.Address{}, err
	}

	status, err := node.Status()
	if err != nil {
		return common.Address{}, err
	}

	return common.BytesToAddress(status.ValidatorInfo.Address.Bytes()), nil
}

// Mining returns whether or not this node is currently mining. Always false.
func (api *PublicEthereumAPI) Mining() bool {
	api.logger.Debug("eth_mining")
	return false
}

// Hashrate returns the current node's hashrate. Always 0.
func (api *PublicEthereumAPI) Hashrate() hexutil.Uint64 {
	api.logger.Debug("eth_hashrate")
	return 0
}

// GasPrice returns the current gas price based on Ethermint's gas price oracle.
func (api *PublicEthereumAPI) GasPrice() *hexutil.Big {
	monitor := monitor.GetMonitor("eth_gasPrice", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd()

	if app.GlobalGpIndex.RecommendGp != nil {
		return (*hexutil.Big)(app.GlobalGpIndex.RecommendGp)
	}

	return api.gasPrice
}

// Accounts returns the list of accounts available to this node.
func (api *PublicEthereumAPI) Accounts() ([]common.Address, error) {
	monitor := monitor.GetMonitor("eth_accounts", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd()
	return api.accounts()
}

func (api *PublicEthereumAPI) accounts() ([]common.Address, error) {
	api.keyringLock.Lock()
	defer api.keyringLock.Unlock()

	addresses := make([]common.Address, 0) // return [] instead of nil if empty

	infos, err := api.clientCtx.Keybase.List()
	if err != nil {
		return addresses, err
	}

	for _, info := range infos {
		addressBytes := info.GetPubKey().Address().Bytes()
		addresses = append(addresses, common.BytesToAddress(addressBytes))
	}

	return addresses, nil
}

// BlockNumber returns the current block number.
func (api *PublicEthereumAPI) BlockNumber() (hexutil.Uint64, error) {
	monitor := monitor.GetMonitor("eth_blockNumber", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd()
	return api.backend.BlockNumber()
}

// GetBalance returns the provided account's balance up to the provided block number.
func (api *PublicEthereumAPI) GetBalance(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (*hexutil.Big, error) {
	monitor := monitor.GetMonitor("eth_getBalance", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("address", address, "block number", blockNrOrHash)
	acc, err := api.wrappedBackend.MustGetAccount(address.Bytes())
	if err == nil {
		balance := acc.GetCoins().AmountOf(sdk.DefaultBondDenom).BigInt()
		if balance == nil {
			return (*hexutil.Big)(sdk.ZeroInt().BigInt()), nil
		}
		return (*hexutil.Big)(balance), nil
	}

	blockNum, err := api.backend.ConvertToBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}
	clientCtx := api.clientCtx
	if !(blockNum == rpctypes.PendingBlockNumber || blockNum == rpctypes.LatestBlockNumber) {
		clientCtx = api.clientCtx.WithHeight(blockNum.Int64())
	}

	bs, err := api.clientCtx.Codec.MarshalJSON(auth.NewQueryAccountParams(address.Bytes()))
	if err != nil {
		return nil, err
	}

	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", auth.QuerierRoute, auth.QueryAccount), bs)
	if err != nil {
		api.saveZeroAccount(address)
		return (*hexutil.Big)(sdk.ZeroInt().BigInt()), nil
	}

	var account ethermint.EthAccount
	if err := api.clientCtx.Codec.UnmarshalJSON(res, &account); err != nil {
		return nil, err
	}

	val := account.Balance(sdk.DefaultBondDenom).BigInt()
	api.watcherBackend.CommitAccountToRpcDb(account)
	if blockNum != rpctypes.PendingBlockNumber {
		return (*hexutil.Big)(val), nil
	}

	// update the address balance with the pending transactions value (if applicable)
	pendingTxs, err := api.backend.UserPendingTransactions(address.String(), -1)
	if err != nil {
		return nil, err
	}

	for _, tx := range pendingTxs {
		if tx == nil {
			continue
		}

		if tx.From == address {
			val = new(big.Int).Sub(val, tx.Value.ToInt())
		}
		if *tx.To == address {
			val = new(big.Int).Add(val, tx.Value.ToInt())
		}
	}

	return (*hexutil.Big)(val), nil
}

// GetAccount returns the provided account's balance up to the provided block number.
func (api *PublicEthereumAPI) GetAccount(address common.Address) (*ethermint.EthAccount, error) {
	acc, err := api.wrappedBackend.MustGetAccount(address.Bytes())
	if err == nil {
		return acc, nil
	}
	clientCtx := api.clientCtx

	bs, err := api.clientCtx.Codec.MarshalJSON(auth.NewQueryAccountParams(address.Bytes()))
	if err != nil {
		return nil, err
	}

	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", auth.QuerierRoute, auth.QueryAccount), bs)
	if err != nil {
		return nil, err
	}

	var account ethermint.EthAccount
	if err := api.clientCtx.Codec.UnmarshalJSON(res, &account); err != nil {
		return nil, err
	}

	api.watcherBackend.CommitAccountToRpcDb(account)

	return &account, nil
}

func (api *PublicEthereumAPI) getStorageAt(address common.Address, key []byte, blockNum rpctypes.BlockNumber, directlyKey bool) (hexutil.Bytes, error) {
	clientCtx := api.clientCtx.WithHeight(blockNum.Int64())
	res, err := api.wrappedBackend.MustGetState(address, key)
	if err == nil {
		return res, nil
	}
	var queryStr = ""
	if !directlyKey {
		queryStr = fmt.Sprintf("custom/%s/storage/%s/%X", evmtypes.ModuleName, address.Hex(), key)
	} else {
		queryStr = fmt.Sprintf("custom/%s/storageKey/%s/%X", evmtypes.ModuleName, address.Hex(), key)
	}

	res, _, err = clientCtx.QueryWithData(queryStr, nil)
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResStorage
	api.clientCtx.Codec.MustUnmarshalJSON(res, &out)

	api.watcherBackend.CommitStateToRpcDb(address, key, out.Value)
	return out.Value, nil
}

// GetStorageAt returns the contract storage at the given address, block number, and key.
func (api *PublicEthereumAPI) GetStorageAt(address common.Address, key string, blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Bytes, error) {
	monitor := monitor.GetMonitor("eth_getStorageAt", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("address", address, "key", key, "block number", blockNrOrHash)
	blockNum, err := api.backend.ConvertToBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}
	return api.getStorageAt(address, common.HexToHash(key).Bytes(), blockNum, false)
}

// GetStorageAtInternal returns the contract storage at the given address, block number, and key.
func (api *PublicEthereumAPI) GetStorageAtInternal(address common.Address, key []byte) (hexutil.Bytes, error) {
	return api.getStorageAt(address, key, 0, true)
}

// GetTransactionCount returns the number of transactions at the given address up to the given block number.
func (api *PublicEthereumAPI) GetTransactionCount(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (*hexutil.Uint64, error) {
	monitor := monitor.GetMonitor("eth_getTransactionCount", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("address", address, "block number", blockNrOrHash)

	blockNum, err := api.backend.ConvertToBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}
	clientCtx := api.clientCtx
	pending := blockNum == rpctypes.PendingBlockNumber
	// pass the given block height to the context if the height is not pending or latest
	if !pending && blockNum != rpctypes.LatestBlockNumber {
		clientCtx = api.clientCtx.WithHeight(blockNum.Int64())
	}

	nonce, err := api.accountNonce(clientCtx, address, pending)
	if err != nil {
		return nil, err
	}

	n := hexutil.Uint64(nonce)
	return &n, nil
}

// GetBlockTransactionCountByHash returns the number of transactions in the block identified by hash.
func (api *PublicEthereumAPI) GetBlockTransactionCountByHash(hash common.Hash) *hexutil.Uint {
	monitor := monitor.GetMonitor("eth_getBlockTransactionCountByHash", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("hash", hash)
	res, _, err := api.clientCtx.Query(fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryHashToHeight, hash.Hex()))
	if err != nil {
		return nil
	}

	var out evmtypes.QueryResBlockNumber
	if err := api.clientCtx.Codec.UnmarshalJSON(res, &out); err != nil {
		return nil
	}

	resBlock, err := api.clientCtx.Client.Block(&out.Number)
	if err != nil {
		return nil
	}

	n := hexutil.Uint(len(resBlock.Block.Txs))
	return &n
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block identified by its height.
func (api *PublicEthereumAPI) GetBlockTransactionCountByNumber(blockNum rpctypes.BlockNumber) *hexutil.Uint {
	monitor := monitor.GetMonitor("eth_getBlockTransactionCountByNumber", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("block number", blockNum)
	var (
		height  int64
		err     error
		txCount hexutil.Uint
		txs     int
	)

	switch blockNum {
	case rpctypes.PendingBlockNumber:
		height, err = api.backend.LatestBlockNumber()
		if err != nil {
			return nil
		}
		resBlock, err := api.clientCtx.Client.Block(&height)
		if err != nil {
			return nil
		}
		// get the pending transaction count
		pendingCnt, err := api.backend.PendingTransactionCnt()
		if err != nil {
			return nil
		}
		txs = len(resBlock.Block.Txs) + pendingCnt
	case rpctypes.LatestBlockNumber:
		height, err = api.backend.LatestBlockNumber()
		if err != nil {
			return nil
		}
		resBlock, err := api.clientCtx.Client.Block(&height)
		if err != nil {
			return nil
		}
		txs = len(resBlock.Block.Txs)
	default:
		height = blockNum.Int64()
		resBlock, err := api.clientCtx.Client.Block(&height)
		if err != nil {
			return nil
		}
		txs = len(resBlock.Block.Txs)
	}

	txCount = hexutil.Uint(txs)
	return &txCount
}

// GetUncleCountByBlockHash returns the number of uncles in the block idenfied by hash. Always zero.
func (api *PublicEthereumAPI) GetUncleCountByBlockHash(_ common.Hash) hexutil.Uint {
	return 0
}

// GetUncleCountByBlockNumber returns the number of uncles in the block idenfied by number. Always zero.
func (api *PublicEthereumAPI) GetUncleCountByBlockNumber(_ rpctypes.BlockNumber) hexutil.Uint {
	return 0
}

// GetCode returns the contract code at the given address and block number.
func (api *PublicEthereumAPI) GetCode(address common.Address, blockNrOrHash rpctypes.BlockNumberOrHash) (hexutil.Bytes, error) {
	monitor := monitor.GetMonitor("eth_getCode", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("address", address, "block number", blockNrOrHash)
	blockNumber, err := api.backend.ConvertToBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	code, err := api.wrappedBackend.GetCode(address, uint64(blockNumber))
	if err == nil {
		return code, nil
	}

	clientCtx := api.clientCtx
	if !(blockNumber == rpctypes.PendingBlockNumber || blockNumber == rpctypes.LatestBlockNumber) {
		clientCtx = api.clientCtx.WithHeight(blockNumber.Int64())
	}
	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryCode, address.Hex()), nil)
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResCode
	api.clientCtx.Codec.MustUnmarshalJSON(res, &out)
	return out.Code, nil
}

// GetCodeByHash returns the contract code at the given address and block number.
func (api *PublicEthereumAPI) GetCodeByHash(hash common.Hash) (hexutil.Bytes, error) {
	code, err := api.wrappedBackend.GetCodeByHash(hash.Bytes())
	if err == nil {
		return code, nil
	}
	clientCtx := api.clientCtx
	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryCodeByHash, hash.Hex()), nil)
	if err != nil {
		return nil, err
	}

	var out evmtypes.QueryResCode
	api.clientCtx.Codec.MustUnmarshalJSON(res, &out)

	api.watcherBackend.CommitCodeHashToDb(hash.Bytes(), out.Code)

	return out.Code, nil
}

// GetTransactionLogs returns the logs given a transaction hash.
func (api *PublicEthereumAPI) GetTransactionLogs(txHash common.Hash) ([]*ethtypes.Log, error) {
	api.logger.Debug("eth_getTransactionLogs", "hash", txHash)
	return api.backend.GetTransactionLogs(txHash)
}

// Sign signs the provided data using the private key of address via Geth's signature standard.
func (api *PublicEthereumAPI) Sign(address common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	monitor := monitor.GetMonitor("eth_sign", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("address", address, "data", data)
	// TODO: Change this functionality to find an unlocked account by address

	key, exist := rpctypes.GetKeyByAddress(api.keys, address)
	if !exist {
		return nil, keystore.ErrLocked
	}

	// Sign the requested hash with the wallet
	sig, err := crypto.Sign(accounts.TextHash(data), key.ToECDSA())
	if err != nil {
		return nil, err
	}

	sig[crypto.RecoveryIDOffset] += 27 // transform V from 0/1 to 27/28

	return sig, nil
}

// SendTransaction sends an Ethereum transaction.
func (api *PublicEthereumAPI) SendTransaction(args rpctypes.SendTxArgs) (common.Hash, error) {
	monitor := monitor.GetMonitor("eth_sendTransaction", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("args", args)
	// TODO: Change this functionality to find an unlocked account by address

	key, exist := rpctypes.GetKeyByAddress(api.keys, *args.From)
	if !exist {
		api.logger.Debug("failed to find key in keyring", "key", args.From)
		return common.Hash{}, keystore.ErrLocked
	}

	// Mutex lock the address' nonce to avoid assigning it to multiple requests
	if args.Nonce == nil {
		api.nonceLock.LockAddr(*args.From)
		defer api.nonceLock.UnlockAddr(*args.From)
	}

	// Assemble transaction from fields
	tx, err := api.generateFromArgs(args)
	if err != nil {
		api.logger.Debug("failed to generate tx", "error", err)
		return common.Hash{}, err
	}

	if err := tx.ValidateBasic(); err != nil {
		api.logger.Debug("tx failed basic validation", "error", err)
		return common.Hash{}, err
	}

	// Sign transaction
	if err := tx.Sign(api.chainIDEpoch, key.ToECDSA()); err != nil {
		api.logger.Debug("failed to sign tx", "error", err)
		return common.Hash{}, err
	}

	// Encode transaction by default Tx encoder
	txEncoder := authclient.GetTxEncoder(api.clientCtx.Codec)
	txBytes, err := txEncoder(tx)
	if err != nil {
		return common.Hash{}, err
	}

	// send chanData to txPool
	if api.txPool != nil {
		return broadcastTxByTxPool(api, tx, txBytes)
	}

	// Broadcast transaction in sync mode (default)
	// NOTE: If error is encountered on the node, the broadcast will not return an error
	res, err := api.clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return common.Hash{}, err
	}

	if res.Code != abci.CodeTypeOK {
		return CheckError(res)
	}

	// Return transaction hash
	return common.HexToHash(res.TxHash), nil
}

// SendRawTransaction send a raw Ethereum transaction.
func (api *PublicEthereumAPI) SendRawTransaction(data hexutil.Bytes) (common.Hash, error) {
	monitor := monitor.GetMonitor("eth_sendRawTransaction", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("data", data)
	tx := new(evmtypes.MsgEthereumTx)

	// RLP decode raw transaction bytes
	if err := rlp.DecodeBytes(data, tx); err != nil {
		// Return nil is for when gasLimit overflows uint64
		return common.Hash{}, err
	}

	// Encode transaction by default Tx encoder
	txEncoder := authclient.GetTxEncoder(api.clientCtx.Codec)
	txBytes, err := txEncoder(tx)
	if err != nil {
		return common.Hash{}, err
	}

	// send chanData to txPool
	if api.txPool != nil {
		return broadcastTxByTxPool(api, tx, txBytes)
	}

	// TODO: Possibly log the contract creation address (if recipient address is nil) or tx data
	// If error is encountered on the node, the broadcast will not return an error
	res, err := api.clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return common.Hash{}, err
	}

	if res.Code != abci.CodeTypeOK {
		return CheckError(res)
	}
	// Return transaction hash
	return common.HexToHash(res.TxHash), nil
}

func (api *PublicEthereumAPI) buildKey(args rpctypes.CallArgs) common.Hash {
	latest, e := api.wrappedBackend.GetLatestBlockNumber()
	if e != nil {
		return common.Hash{}
	}
	return sha256.Sum256([]byte(args.String() + strconv.Itoa(int(latest))))
}

func (api *PublicEthereumAPI) getFromCallCache(key common.Hash) ([]byte, bool) {
	if api.callCache == nil {
		return nil, false
	}
	nilKey := common.Hash{}
	if key == nilKey {
		return nil, false
	}
	cacheData, ok := api.callCache.Get(key)
	if ok {
		ret, ok := cacheData.([]byte)
		if ok {
			return ret, true
		}
	}
	return nil, false
}

func (api *PublicEthereumAPI) addCallCache(key common.Hash, data []byte) {
	if api.callCache == nil {
		return
	}
	api.callCache.Add(key, data)
}

// Call performs a raw contract call.
func (api *PublicEthereumAPI) Call(args rpctypes.CallArgs, blockNrOrHash rpctypes.BlockNumberOrHash, _ *map[common.Address]rpctypes.Account) (hexutil.Bytes, error) {
	monitor := monitor.GetMonitor("eth_call", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("args", args, "block number", blockNrOrHash)
	key := api.buildKey(args)
	cacheData, ok := api.getFromCallCache(key)
	if ok {
		return cacheData, nil
	}
	blockNr, err := api.backend.ConvertToBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}
	simRes, err := api.doCall(args, blockNr, big.NewInt(ethermint.DefaultRPCGasLimit), false)
	if err != nil {
		return []byte{}, TransformDataError(err, "eth_call")
	}

	data, err := evmtypes.DecodeResultData(simRes.Result.Data)
	if err != nil {
		return []byte{}, TransformDataError(err, "eth_call")
	}
	api.addCallCache(key, data.Ret)
	return data.Ret, nil
}

// MultiCall performs multiple raw contract call.
func (api *PublicEthereumAPI) MultiCall(args []rpctypes.CallArgs, blockNr rpctypes.BlockNumber, _ *map[common.Address]rpctypes.Account) ([]hexutil.Bytes, error) {
	if !viper.GetBool(FlagEnableMultiCall) {
		return nil, errors.New("the method is not allowed")
	}

	monitor := monitor.GetMonitor("eth_multiCall", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("args", args, "block number", blockNr)

	blockNrOrHash := rpctypes.BlockNumberOrHashWithNumber(blockNr)
	rets := make([]hexutil.Bytes, 0, len(args))
	for _, arg := range args {
		ret, err := api.Call(arg, blockNrOrHash, nil)
		if err != nil {
			return rets, err
		}
		rets = append(rets, ret)
	}
	return rets, nil
}

// DoCall performs a simulated call operation through the evmtypes. It returns the
// estimated gas used on the operation or an error if fails.
func (api *PublicEthereumAPI) doCall(
	args rpctypes.CallArgs, blockNum rpctypes.BlockNumber, globalGasCap *big.Int, isEstimate bool,
) (*sdk.SimulationResponse, error) {

	clientCtx := api.clientCtx
	// pass the given block height to the context if the height is not pending or latest
	if !(blockNum == rpctypes.PendingBlockNumber || blockNum == rpctypes.LatestBlockNumber) {
		clientCtx = api.clientCtx.WithHeight(blockNum.Int64())
	}

	// Set sender address or use a default if none specified
	var addr common.Address
	if args.From != nil {
		addr = *args.From
	}

	nonce := uint64(0)
	if isEstimate && args.To == nil && args.Data != nil {
		//only get real nonce when estimate gas and the action is contract deploy
		nonce, _ = api.accountNonce(api.clientCtx, addr, true)
	}

	// Set default gas & gas price if none were set
	// Change this to uint64(math.MaxUint64 / 2) if gas cap can be configured
	gas := uint64(ethermint.DefaultRPCGasLimit)
	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}
	if globalGasCap != nil && globalGasCap.Uint64() < gas {
		api.logger.Debug("Caller gas above allowance, capping", "requested", gas, "cap", globalGasCap)
		gas = globalGasCap.Uint64()
	}

	// Set gas price using default or parameter if passed in
	gasPrice := new(big.Int).SetUint64(ethermint.DefaultGasPrice)
	if args.GasPrice != nil {
		gasPrice = args.GasPrice.ToInt()
	}

	// Set value for transaction
	value := new(big.Int)
	if args.Value != nil {
		value = args.Value.ToInt()
	}

	// Set Data if provided
	var data []byte
	if args.Data != nil {
		data = []byte(*args.Data)
	}

	// Set destination address for call
	var toAddr *sdk.AccAddress
	if args.To != nil {
		pTemp := sdk.AccAddress(args.To.Bytes())
		toAddr = &pTemp
	}

	var msgs []sdk.Msg
	// Create new call message
	msg := evmtypes.NewMsgEthermint(nonce, toAddr, sdk.NewIntFromBigInt(value), gas,
		sdk.NewIntFromBigInt(gasPrice), data, sdk.AccAddress(addr.Bytes()))
	msgs = append(msgs, msg)

	sim := api.evmFactory.BuildSimulator(api)
	//only worked when fast-query has been enabled
	if sim != nil {
		return sim.DoCall(msg)
	}

	//convert the pending transactions into ethermint msgs
	if blockNum == rpctypes.PendingBlockNumber {
		pendingMsgs, err := api.pendingMsgs()
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, pendingMsgs...)
	}

	//Generate tx to be used to simulate (signature isn't needed)
	var stdSig authtypes.StdSignature
	stdSigs := []authtypes.StdSignature{stdSig}

	tx := authtypes.NewStdTx(msgs, authtypes.StdFee{}, stdSigs, "")
	if err := tx.ValidateBasic(); err != nil {
		return nil, err
	}

	txEncoder := authclient.GetTxEncoder(clientCtx.Codec)
	txBytes, err := txEncoder(tx)
	if err != nil {
		return nil, err
	}

	// Transaction simulation through query
	res, _, err := clientCtx.QueryWithData("app/simulate", txBytes)
	if err != nil {
		return nil, err
	}

	var simResponse sdk.SimulationResponse
	if err := clientCtx.Codec.UnmarshalBinaryBare(res, &simResponse); err != nil {
		return nil, err
	}

	return &simResponse, nil
}

// EstimateGas returns an estimate of gas usage for the given smart contract call.
// It adds 1,000 gas to the returned value instead of using the gas adjustment
// param from the SDK.
func (api *PublicEthereumAPI) EstimateGas(args rpctypes.CallArgs) (hexutil.Uint64, error) {
	monitor := monitor.GetMonitor("eth_estimateGas", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("args", args)

	simResponse, err := api.doCall(args, 0, big.NewInt(ethermint.DefaultRPCGasLimit), true)
	if err != nil {
		return 0, TransformDataError(err, "eth_estimateGas")
	}

	// TODO: change 1000 buffer for more accurate buffer (eg: SDK's gasAdjusted)
	estimatedGas := simResponse.GasInfo.GasUsed
	gasBuffer := estimatedGas / 100 * config.GetOecConfig().GetGasLimitBuffer()
	gas := estimatedGas + gasBuffer

	return hexutil.Uint64(gas), nil
}

// GetBlockByHash returns the block identified by hash.
func (api *PublicEthereumAPI) GetBlockByHash(hash common.Hash, fullTx bool) (interface{}, error) {
	monitor := monitor.GetMonitor("eth_getBlockByHash", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("hash", hash, "full", fullTx)
	block, err := api.backend.GetBlockByHash(hash, fullTx)
	if err != nil {
		return nil, TransformDataError(err, RPCEthGetBlockByHash)
	}
	return block, nil
}

// GetBlockByNumber returns the block identified by number.
func (api *PublicEthereumAPI) GetBlockByNumber(blockNum rpctypes.BlockNumber, fullTx bool) (interface{}, error) {
	monitor := monitor.GetMonitor("eth_getBlockByNumber", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("number", blockNum, "full", fullTx)
	var blockTxs interface{}
	if blockNum != rpctypes.PendingBlockNumber {
		return api.backend.GetBlockByNumber(blockNum, fullTx)
	}

	height, err := api.backend.LatestBlockNumber()
	if err != nil {
		return nil, err
	}

	// latest block info
	latestBlock, err := api.clientCtx.Client.Block(&height)
	if err != nil {
		return nil, err
	}

	// number of pending txs queried from the mempool
	unconfirmedTxs, err := api.clientCtx.Client.UnconfirmedTxs(1000)
	if err != nil {
		return nil, err
	}

	pendingTxs, gasUsed, ethTxs, err := rpctypes.EthTransactionsFromTendermint(api.clientCtx, unconfirmedTxs.Txs, common.BytesToHash(latestBlock.Block.Hash()), uint64(height))
	if err != nil {
		return nil, err
	}

	if fullTx {
		blockTxs = ethTxs
	} else {
		blockTxs = pendingTxs
	}

	return rpctypes.FormatBlock(
		tmtypes.Header{
			Version:         latestBlock.Block.Version,
			ChainID:         api.clientCtx.ChainID,
			Height:          height + 1,
			Time:            time.Unix(0, 0),
			LastBlockID:     latestBlock.Block.LastBlockID,
			ValidatorsHash:  latestBlock.Block.NextValidatorsHash,
			ProposerAddress: latestBlock.Block.ProposerAddress,
		},
		0,
		latestBlock.Block.Hash(),
		0,
		gasUsed,
		blockTxs,
		ethtypes.Bloom{},
	), nil

}

// GetTransactionByHash returns the transaction identified by hash.
func (api *PublicEthereumAPI) GetTransactionByHash(hash common.Hash) (*rpctypes.Transaction, error) {
	monitor := monitor.GetMonitor("eth_getTransactionByHash", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("hash", hash)
	rawTx, err := api.wrappedBackend.GetTransactionByHash(hash)
	if err == nil {
		return rawTx, nil
	}
	tx, err := api.clientCtx.Client.Tx(hash.Bytes(), false)
	if err != nil {
		// check if the tx is on the mempool
		pendingTx, pendingErr := api.PendingTransactionsByHash(hash)
		if pendingErr != nil {
			//to keep consistent with rpc of ethereum, should be return nil
			return nil, nil
		}
		return pendingTx, nil
	}

	// Can either cache or just leave this out if not necessary
	block, err := api.clientCtx.Client.Block(&tx.Height)
	if err != nil {
		return nil, err
	}

	blockHash := common.BytesToHash(block.Block.Hash())

	ethTx, err := rpctypes.RawTxToEthTx(api.clientCtx, tx.Tx)
	if err != nil {
		return nil, err
	}

	height := uint64(tx.Height)
	return rpctypes.NewTransaction(ethTx, common.BytesToHash(tx.Tx.Hash()), blockHash, height, uint64(tx.Index))
}

// GetTransactionByBlockHashAndIndex returns the transaction identified by hash and index.
func (api *PublicEthereumAPI) GetTransactionByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) (*rpctypes.Transaction, error) {
	monitor := monitor.GetMonitor("eth_getTransactionByBlockHashAndIndex", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("hash", hash, "index", idx)
	res, _, err := api.clientCtx.Query(fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryHashToHeight, hash.Hex()))
	if err != nil {
		return nil, nil
	}

	var out evmtypes.QueryResBlockNumber
	api.clientCtx.Codec.MustUnmarshalJSON(res, &out)

	resBlock, err := api.clientCtx.Client.Block(&out.Number)
	if err != nil {
		return nil, nil
	}

	return api.getTransactionByBlockAndIndex(resBlock.Block, idx)
}

// GetTransactionByBlockNumberAndIndex returns the transaction identified by number and index.
func (api *PublicEthereumAPI) GetTransactionByBlockNumberAndIndex(blockNum rpctypes.BlockNumber, idx hexutil.Uint) (*rpctypes.Transaction, error) {
	monitor := monitor.GetMonitor("eth_getTransactionByBlockNumberAndIndex", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("blockNum", blockNum, "index", idx)
	tx, e := api.wrappedBackend.GetTransactionByBlockNumberAndIndex(uint64(blockNum), uint(idx))
	if e == nil && tx != nil {
		return tx, nil
	}
	var (
		height int64
		err    error
	)

	switch blockNum {
	case rpctypes.PendingBlockNumber:
		// get all the EVM pending txs
		pendingTxs, err := api.backend.PendingTransactions()
		if err != nil {
			return nil, err
		}

		// return if index out of bounds
		if uint64(idx) >= uint64(len(pendingTxs)) {
			return nil, nil
		}

		// change back to pendingTxs[idx] once pending queue is fixed.
		return pendingTxs[int(idx)], nil

	case rpctypes.LatestBlockNumber:
		height, err = api.backend.LatestBlockNumber()
		if err != nil {
			return nil, err
		}

	default:
		height = blockNum.Int64()
	}

	resBlock, err := api.clientCtx.Client.Block(&height)
	if err != nil {
		return nil, err
	}

	return api.getTransactionByBlockAndIndex(resBlock.Block, idx)
}

func (api *PublicEthereumAPI) getTransactionByBlockAndIndex(block *tmtypes.Block, idx hexutil.Uint) (*rpctypes.Transaction, error) {
	// return if index out of bounds
	if uint64(idx) >= uint64(len(block.Txs)) {
		return nil, nil
	}

	ethTx, err := rpctypes.RawTxToEthTx(api.clientCtx, block.Txs[idx])
	if err != nil {
		// return nil error if the transaction is not a MsgEthereumTx
		return nil, nil
	}

	height := uint64(block.Height)
	txHash := common.BytesToHash(block.Txs[idx].Hash())
	blockHash := common.BytesToHash(block.Hash())
	return rpctypes.NewTransaction(ethTx, txHash, blockHash, height, uint64(idx))
}

// GetTransactionsByBlock returns some transactions identified by number or hash.
func (api *PublicEthereumAPI) GetTransactionsByBlock(blockNrOrHash rpctypes.BlockNumberOrHash, offset, limit hexutil.Uint) ([]*rpctypes.Transaction, error) {
	monitor := monitor.GetMonitor("eth_getTransactionsByBlock", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("block number", blockNrOrHash, "offset", offset, "limit", limit)

	blockNum, err := api.backend.ConvertToBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}

	txs, e := api.wrappedBackend.GetTransactionsByBlockNumber(uint64(blockNum), uint64(offset), uint64(limit))
	if e == nil && txs != nil {
		return txs, nil
	}

	height := blockNum.Int64()
	switch blockNum {
	case rpctypes.PendingBlockNumber:
		// get all the EVM pending txs
		pendingTxs, err := api.backend.PendingTransactions()
		if err != nil {
			return nil, err
		}
		switch {
		case len(pendingTxs) <= int(offset):
			return nil, nil
		case len(pendingTxs) < int(offset+limit):
			return pendingTxs[offset:], nil
		default:
			return pendingTxs[offset : offset+limit], nil
		}
	case rpctypes.LatestBlockNumber:
		height, err = api.backend.LatestBlockNumber()
		if err != nil {
			return nil, err
		}
	}

	resBlock, err := api.clientCtx.Client.Block(&height)
	if err != nil {
		return nil, err
	}
	for idx := offset; idx < offset+limit && int(idx) < len(resBlock.Block.Txs); idx++ {
		tx, _ := api.getTransactionByBlockAndIndex(resBlock.Block, idx)
		if tx != nil {
			txs = append(txs, tx)
		}
	}
	return txs, nil
}

// GetTransactionReceipt returns the transaction receipt identified by hash.
func (api *PublicEthereumAPI) GetTransactionReceipt(hash common.Hash) (*watcher.TransactionReceipt, error) {
	monitor := monitor.GetMonitor("eth_getTransactionReceipt", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("hash", hash)
	res, e := api.wrappedBackend.GetTransactionReceipt(hash)
	if e == nil {
		return res, nil
	}

	tx, err := api.clientCtx.Client.Tx(hash.Bytes(), false)
	if err != nil {
		// Return nil for transaction when not found
		return nil, nil
	}

	// Query block for consensus hash
	block, err := api.clientCtx.Client.Block(&tx.Height)
	if err != nil {
		return nil, err
	}

	blockHash := common.BytesToHash(block.Block.Hash())

	// Convert tx bytes to eth transaction
	ethTx, err := rpctypes.RawTxToEthTx(api.clientCtx, tx.Tx)
	if err != nil {
		return nil, err
	}

	fromSigCache, err := ethTx.VerifySig(ethTx.ChainID(), tx.Height, sdk.EmptyContext().SigCache())
	if err != nil {
		return nil, err
	}

	from := fromSigCache.GetFrom()
	cumulativeGasUsed := uint64(tx.TxResult.GasUsed)
	if tx.Index != 0 {
		cumulativeGasUsed += rpctypes.GetBlockCumulativeGas(api.clientCtx.Codec, block.Block, int(tx.Index))
	}

	// Set status codes based on tx result
	var status hexutil.Uint64
	if tx.TxResult.IsOK() {
		status = hexutil.Uint64(1)
	} else {
		status = hexutil.Uint64(0)
	}

	txData := tx.TxResult.GetData()

	data, err := evmtypes.DecodeResultData(txData)
	if err != nil {
		status = 0 // transaction failed
	}

	if len(data.Logs) == 0 {
		data.Logs = []*ethtypes.Log{}
	}
	contractAddr := &data.ContractAddress
	if data.ContractAddress == common.HexToAddress("0x00000000000000000000") {
		contractAddr = nil
	}

	receipt := &watcher.TransactionReceipt{
		Status:            status,
		CumulativeGasUsed: hexutil.Uint64(cumulativeGasUsed),
		LogsBloom:         data.Bloom,
		Logs:              data.Logs,
		TransactionHash:   hash.String(),
		ContractAddress:   contractAddr,
		GasUsed:           hexutil.Uint64(tx.TxResult.GasUsed),
		BlockHash:         blockHash.String(),
		BlockNumber:       hexutil.Uint64(tx.Height),
		TransactionIndex:  hexutil.Uint64(tx.Index),
		From:              from.String(),
		To:                ethTx.To(),
	}

	return receipt, nil
}

// GetTransactionReceiptsByBlock returns the transaction receipt identified by block hash or number.
func (api *PublicEthereumAPI) GetTransactionReceiptsByBlock(blockNrOrHash rpctypes.BlockNumberOrHash, offset, limit hexutil.Uint) ([]*watcher.TransactionReceipt, error) {
	monitor := monitor.GetMonitor("eth_getTransactionReceiptsByBlock", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("block number", blockNrOrHash, "offset", offset, "limit", limit)

	txs, err := api.GetTransactionsByBlock(blockNrOrHash, offset, limit)
	if err != nil || len(txs) == 0 {
		return nil, err
	}

	var receipts []*watcher.TransactionReceipt
	for _, tx := range txs {
		res, _ := api.wrappedBackend.GetTransactionReceipt(tx.Hash)
		if res != nil {
			receipts = append(receipts, res)
			continue
		}

		tx, err := api.clientCtx.Client.Tx(tx.Hash.Bytes(), false)
		if err != nil {
			// Return nil for transaction when not found
			return nil, nil
		}

		// Query block for consensus hash
		block, err := api.clientCtx.Client.Block(&tx.Height)
		if err != nil {
			return nil, err
		}

		blockHash := common.BytesToHash(block.Block.Hash())

		// Convert tx bytes to eth transaction
		ethTx, err := rpctypes.RawTxToEthTx(api.clientCtx, tx.Tx)
		if err != nil {
			return nil, err
		}

		fromSigCache, err := ethTx.VerifySig(ethTx.ChainID(), tx.Height, sdk.EmptyContext().SigCache())
		if err != nil {
			return nil, err
		}

		from := fromSigCache.GetFrom()
		cumulativeGasUsed := uint64(tx.TxResult.GasUsed)
		if tx.Index != 0 {
			cumulativeGasUsed += rpctypes.GetBlockCumulativeGas(api.clientCtx.Codec, block.Block, int(tx.Index))
		}

		// Set status codes based on tx result
		var status hexutil.Uint64
		if tx.TxResult.IsOK() {
			status = hexutil.Uint64(1)
		} else {
			status = hexutil.Uint64(0)
		}

		txData := tx.TxResult.GetData()
		data, err := evmtypes.DecodeResultData(txData)
		if err != nil {
			status = 0 // transaction failed
		}

		if len(data.Logs) == 0 {
			data.Logs = []*ethtypes.Log{}
		}
		contractAddr := &data.ContractAddress
		if data.ContractAddress == common.HexToAddress("0x00000000000000000000") {
			contractAddr = nil
		}

		receipt := &watcher.TransactionReceipt{
			Status:            status,
			CumulativeGasUsed: hexutil.Uint64(cumulativeGasUsed),
			LogsBloom:         data.Bloom,
			Logs:              data.Logs,
			TransactionHash:   common.BytesToHash(tx.Hash.Bytes()).String(),
			ContractAddress:   contractAddr,
			GasUsed:           hexutil.Uint64(tx.TxResult.GasUsed),
			BlockHash:         blockHash.String(),
			BlockNumber:       hexutil.Uint64(tx.Height),
			TransactionIndex:  hexutil.Uint64(tx.Index),
			From:              from.String(),
			To:                ethTx.To(),
		}
		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (api *PublicEthereumAPI) PendingTransactions() ([]*rpctypes.Transaction, error) {
	api.logger.Debug("eth_pendingTransactions")
	return api.backend.PendingTransactions()
}

func (api *PublicEthereumAPI) PendingTransactionsByHash(target common.Hash) (*rpctypes.Transaction, error) {
	api.logger.Debug("eth_pendingTransactionsByHash")
	return api.backend.PendingTransactionsByHash(target)
}

// GetUncleByBlockHashAndIndex returns the uncle identified by hash and index. Always returns nil.
func (api *PublicEthereumAPI) GetUncleByBlockHashAndIndex(hash common.Hash, idx hexutil.Uint) map[string]interface{} {
	return nil
}

// GetUncleByBlockNumberAndIndex returns the uncle identified by number and index. Always returns nil.
func (api *PublicEthereumAPI) GetUncleByBlockNumberAndIndex(number hexutil.Uint, idx hexutil.Uint) map[string]interface{} {
	return nil
}

// GetProof returns an account object with proof and any storage proofs
func (api *PublicEthereumAPI) GetProof(address common.Address, storageKeys []string, blockNrOrHash rpctypes.BlockNumberOrHash) (*rpctypes.AccountResult, error) {
	monitor := monitor.GetMonitor("eth_getProof", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("address", address, "keys", storageKeys, "number", blockNrOrHash)
	blockNum, err := api.backend.ConvertToBlockNumber(blockNrOrHash)
	if err != nil {
		return nil, err
	}
	clientCtx := api.clientCtx.WithHeight(int64(blockNum))
	path := fmt.Sprintf("custom/%s/%s/%s", evmtypes.ModuleName, evmtypes.QueryAccount, address.Hex())

	// query eth account at block height
	resBz, _, err := clientCtx.Query(path)
	if err != nil {
		return nil, err
	}

	var account evmtypes.QueryResAccount
	clientCtx.Codec.MustUnmarshalJSON(resBz, &account)

	storageProofs := make([]rpctypes.StorageResult, len(storageKeys))
	for i, k := range storageKeys {
		data := append(evmtypes.AddressStoragePrefix(address), getStorageByAddressKey(address, common.HexToHash(k).Bytes()).Bytes()...)
		// Get value for key
		req := abci.RequestQuery{
			Path:   fmt.Sprintf("store/%s/key", evmtypes.StoreKey),
			Data:   data,
			Height: int64(blockNum),
			Prove:  true,
		}

		vRes, err := clientCtx.QueryABCI(req)
		if err != nil {
			return nil, err
		}

		var value evmtypes.QueryResStorage
		value.Value = vRes.GetValue()

		// check for proof
		proof := vRes.GetProof()
		proofStr := new(merkle.Proof).String()
		if proof != nil {
			proofStr = proof.String()
		}

		storageProofs[i] = rpctypes.StorageResult{
			Key:   k,
			Value: (*hexutil.Big)(common.BytesToHash(value.Value).Big()),
			Proof: []string{proofStr},
		}
	}

	req := abci.RequestQuery{
		Path:   fmt.Sprintf("store/%s/key", auth.StoreKey),
		Data:   auth.AddressStoreKey(sdk.AccAddress(address.Bytes())),
		Height: int64(blockNum),
		Prove:  true,
	}

	res, err := clientCtx.QueryABCI(req)
	if err != nil {
		return nil, err
	}

	// check for proof
	accountProof := res.GetProof()
	accProofStr := new(merkle.Proof).String()
	if accountProof != nil {
		accProofStr = accountProof.String()
	}

	return &rpctypes.AccountResult{
		Address:      address,
		AccountProof: []string{accProofStr},
		Balance:      (*hexutil.Big)(utils.MustUnmarshalBigInt(account.Balance)),
		CodeHash:     common.BytesToHash(account.CodeHash),
		Nonce:        hexutil.Uint64(account.Nonce),
		StorageHash:  common.Hash{}, // Ethermint doesn't have a storage hash
		StorageProof: storageProofs,
	}, nil
}

// generateFromArgs populates tx message with args (used in RPC API)
func (api *PublicEthereumAPI) generateFromArgs(args rpctypes.SendTxArgs) (*evmtypes.MsgEthereumTx, error) {
	var (
		nonce, gasLimit uint64
		err             error
	)

	amount := (*big.Int)(args.Value)
	gasPrice := (*big.Int)(args.GasPrice)

	if args.GasPrice == nil {
		// Set default gas price
		// TODO: Change to min gas price from context once available through server/daemon
		gasPrice = ParseGasPrice().ToInt()
	}

	if args.Nonce != nil && (uint64)(*args.Nonce) > 0 {
		nonce = (uint64)(*args.Nonce)
	} else {
		// get the nonce from the account retriever and the pending transactions
		nonce, err = api.accountNonce(api.clientCtx, *args.From, true)
		if err != nil {
			return nil, err
		}
	}

	if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
		return nil, errors.New("both 'data' and 'input' are set and not equal. Please use 'input' to pass transaction call data")
	}

	// Sets input to either Input or Data, if both are set and not equal error above returns
	var input hexutil.Bytes
	if args.Input != nil {
		input = *args.Input
	} else if args.Data != nil {
		input = *args.Data
	}

	if args.To == nil && len(input) == 0 {
		// Contract creation
		return nil, fmt.Errorf("contract creation without any data provided")
	}

	if args.Gas == nil {
		callArgs := rpctypes.CallArgs{
			From:     args.From,
			To:       args.To,
			Gas:      args.Gas,
			GasPrice: args.GasPrice,
			Value:    args.Value,
			Data:     &input,
		}
		gl, err := api.EstimateGas(callArgs)
		if err != nil {
			return nil, err
		}
		gasLimit = uint64(gl)
	} else {
		gasLimit = (uint64)(*args.Gas)
	}
	msg := evmtypes.NewMsgEthereumTx(nonce, args.To, amount, gasLimit, gasPrice, input)

	return &msg, nil
}

// pendingMsgs constructs an array of sdk.Msg. This method will check pending transactions and convert
// those transactions into ethermint messages.
func (api *PublicEthereumAPI) pendingMsgs() ([]sdk.Msg, error) {
	// nolint: prealloc
	var msgs []sdk.Msg

	pendingTxs, err := api.PendingTransactions()
	if err != nil {
		return nil, err
	}

	for _, pendingTx := range pendingTxs {
		// NOTE: we have to construct the EVM transaction instead of just casting from the tendermint
		// transactions because PendingTransactions only checks for MsgEthereumTx messages.

		pendingTo := sdk.AccAddress(pendingTx.To.Bytes())
		pendingFrom := sdk.AccAddress(pendingTx.From.Bytes())
		pendingGas, err := hexutil.DecodeUint64(pendingTx.Gas.String())
		if err != nil {
			return nil, err
		}

		pendingValue := pendingTx.Value.ToInt()
		pendingGasPrice := new(big.Int).SetUint64(ethermint.DefaultGasPrice)
		if pendingTx.GasPrice != nil {
			pendingGasPrice = pendingTx.GasPrice.ToInt()
		}

		pendingData := pendingTx.Input
		nonce, _ := api.accountNonce(api.clientCtx, pendingTx.From, true)

		msg := evmtypes.NewMsgEthermint(nonce, &pendingTo, sdk.NewIntFromBigInt(pendingValue), pendingGas,
			sdk.NewIntFromBigInt(pendingGasPrice), pendingData, pendingFrom)

		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// accountNonce returns looks up the transaction nonce count for a given address. If the pending boolean
// is set to true, it will add to the counter all the uncommitted EVM transactions sent from the address.
// NOTE: The function returns no error if the account doesn't exist.
func (api *PublicEthereumAPI) accountNonce(
	clientCtx clientcontext.CLIContext, address common.Address, pending bool,
) (uint64, error) {
	// Get nonce (sequence) from sender account
	nonce := uint64(0)
	acc, err := api.wrappedBackend.MustGetAccount(address.Bytes())
	if err == nil { // account in watch db
		nonce = acc.GetSequence()
	} else {
		// use a the given client context in case its wrapped with a custom height
		accRet := authtypes.NewAccountRetriever(clientCtx)
		from := sdk.AccAddress(address.Bytes())
		account, err := accRet.GetAccount(from)
		if err != nil {
			// account doesn't exist yet, return 0
			return 0, nil
		}
		nonce = account.GetSequence()
		api.watcherBackend.CommitAccountToRpcDb(account)
	}

	if !pending {
		return nonce, nil
	}

	// the account retriever doesn't include the uncommitted transactions on the nonce so we need to
	// to manually add them.
	pendingTxs, err := api.backend.UserPendingTransactionsCnt(address.String())
	if err == nil {
		nonce += uint64(pendingTxs)
	}

	return nonce, nil
}

// GetTxTrace returns the trace of tx execution by txhash.
func (api *PublicEthereumAPI) GetTxTrace(txHash common.Hash) json.RawMessage {
	monitor := monitor.GetMonitor("eth_getTxTrace", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("hash", txHash)

	return json.RawMessage(evmtypes.GetTracesFromDB(txHash.Bytes()))
}

// DeleteTxTrace delete the trace of tx execution by txhash.
func (api *PublicEthereumAPI) DeleteTxTrace(txHash common.Hash) string {
	monitor := monitor.GetMonitor("eth_deleteTxTrace", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd("hash", txHash)

	if err := evmtypes.DeleteTracesFromDB(txHash.Bytes()); err != nil {
		return "delete trace failed"
	}
	return "delete trace succeed"
}

func (api *PublicEthereumAPI) saveZeroAccount(address common.Address) {
	zeroAccount := ethermint.EthAccount{BaseAccount: &auth.BaseAccount{}}
	zeroAccount.SetAddress(address.Bytes())
	zeroAccount.SetBalance(sdk.DefaultBondDenom, sdk.ZeroDec())
	api.watcherBackend.CommitAccountToRpcDb(zeroAccount)
}

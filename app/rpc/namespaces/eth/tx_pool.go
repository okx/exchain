package eth

import (
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	clientcontext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	ethermint "github.com/okex/exchain/app/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmdb "github.com/tendermint/tm-db"
)

const (
	FlagEnableTxPool      = "enable-tx-pool"
	TxPoolCap             = "tx-pool-cap"
	BroadcastPeriodSecond = "broadcast-period-second"
	txPoolDb              = "tx_pool"
)

var broadcastErrors = map[uint32]*sdkerrors.Error{
	sdkerrors.ErrTxInMempoolCache.ABCICode(): sdkerrors.ErrTxInMempoolCache,
	sdkerrors.ErrMempoolIsFull.ABCICode():    sdkerrors.ErrMempoolIsFull,
	sdkerrors.ErrTxTooLarge.ABCICode():       sdkerrors.ErrTxTooLarge,
}

type TxPool struct {
	addressTxsPool    map[common.Address][]*evmtypes.MsgEthereumTx // All currently processable transactions
	clientCtx         clientcontext.CLIContext
	db                tmdb.DB
	mu                sync.Mutex
	cap               uint64
	broadcastInterval time.Duration
}

func NewTxPool(clientCtx clientcontext.CLIContext, api *PublicEthereumAPI) *TxPool {
	db, err := openDB()
	if err != nil {
		panic(err)
	}
	interval := time.Second * time.Duration(viper.GetInt(BroadcastPeriodSecond))
	pool := &TxPool{
		addressTxsPool:    make(map[common.Address][]*evmtypes.MsgEthereumTx),
		clientCtx:         clientCtx,
		db:                db,
		cap:               viper.GetUint64(TxPoolCap),
		broadcastInterval: interval,
	}

	if err = pool.initDB(api); err != nil {
		panic(err)
	}

	return pool
}

func openDB() (tmdb.DB, error) {
	rootDir := viper.GetString("home")
	dataDir := filepath.Join(rootDir, "data")
	return sdk.NewLevelDB(txPoolDb, dataDir)
}

func (pool *TxPool) initDB(api *PublicEthereumAPI) error {
	itr, err := pool.db.Iterator(nil, nil)
	if err != nil {
		return err
	}
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		key := string(itr.Key())
		txBytes := itr.Value()

		tmp := strings.Split(key, "|")
		if len(tmp) != 2 {
			continue
		}
		address := common.HexToAddress(tmp[0])
		txNonce, err := strconv.Atoi(tmp[1])
		if err != nil {
			return err
		}

		tx := new(evmtypes.MsgEthereumTx)
		if err = rlp.DecodeBytes(txBytes, tx); err != nil {
			return err
		}
		if int(tx.Data.AccountNonce) != txNonce {
			return fmt.Errorf("nonce[%d] in key is not equal to nonce[%d] in value", tx.Data.AccountNonce, txNonce)
		}

		pCurrentNonce, err := api.GetTransactionCount(address, rpctypes.PendingBlockNumber)
		if err != nil {
			return err
		}
		currentNonce := int(*pCurrentNonce)
		if txNonce < currentNonce {
			continue
		}

		if err = pool.insertTx(address, tx); err != nil {
			return err
		}
	}

	return nil
}

func broadcastTxByTxPool(api *PublicEthereumAPI, tx *evmtypes.MsgEthereumTx, txBytes []byte) (common.Hash, error) {
	// Get sender address
	chainIDEpoch, err := ethermint.ParseChainID(api.clientCtx.ChainID)
	if err != nil {
		return common.Hash{}, err
	}
	from, err := tx.VerifySig(chainIDEpoch, api.clientCtx.Height)
	if err != nil {
		return common.Hash{}, err
	}

	api.txPool.mu.Lock()
	defer api.txPool.mu.Unlock()
	if err = api.txPool.CacheAndBroadcastTx(api, from, tx); err != nil {
		api.logger.Error("eth_sendRawTransaction txPool err:", err.Error())
		return common.Hash{}, err
	}

	return common.HexToHash(strings.ToUpper(hex.EncodeToString(tmhash.Sum(txBytes)))), nil
}

func (pool *TxPool) CacheAndBroadcastTx(api *PublicEthereumAPI, address common.Address, tx *evmtypes.MsgEthereumTx) error {
	// get currentNonce
	pCurrentNonce, err := api.GetTransactionCount(address, rpctypes.PendingBlockNumber)
	if err != nil {
		return err
	}
	currentNonce := uint64(*pCurrentNonce)

	if tx.Data.AccountNonce < currentNonce {
		return fmt.Errorf("AccountNonce of tx is less than currentNonce in memPool: AccountNonce[%d], currentNonce[%d]", tx.Data.AccountNonce, currentNonce)
	}

	if tx.Data.AccountNonce > currentNonce+pool.cap {
		return fmt.Errorf("AccountNonce of tx is bigger than txPool capacity, please try later: AccountNonce[%d]", tx.Data.AccountNonce)
	}

	if err = pool.insertTx(address, tx); err != nil {
		return err
	}

	// update DB
	if err = pool.writeTxInDB(address, tx); err != nil {
		return err
	}

	_ = pool.continueBroadcast(api, currentNonce, address)

	return nil
}

func (pool *TxPool) update(index int, address common.Address, tx *evmtypes.MsgEthereumTx) error {
	txsLen := len(pool.addressTxsPool[address])
	if index >= txsLen {
		pool.addressTxsPool[address] = append(pool.addressTxsPool[address], tx)
	} else {
		tmpTx := make([]*evmtypes.MsgEthereumTx, len(pool.addressTxsPool[address][index:]))
		copy(tmpTx, pool.addressTxsPool[address][index:])

		pool.addressTxsPool[address] =
			append(append(pool.addressTxsPool[address][:index], tx), tmpTx...)
	}
	return nil
}

// insert the tx into the txPool
func (pool *TxPool) insertTx(address common.Address, tx *evmtypes.MsgEthereumTx) error {
	// if this is the first time to insertTx, make the cap of txPool be TxPoolSliceMaxLen
	if _, ok := pool.addressTxsPool[address]; !ok {
		pool.addressTxsPool[address] = make([]*evmtypes.MsgEthereumTx, 0, pool.cap)
	}
	index := 0
	for index < len(pool.addressTxsPool[address]) {
		// the tx nonce has in txPool, drop duplicate tx
		if tx.Data.AccountNonce == pool.addressTxsPool[address][index].Data.AccountNonce {
			return fmt.Errorf("duplicate tx, this AccountNonce of tx has been in txPool. AccountNonce[%d]", tx.Data.AccountNonce)
		}
		// find the index to insert
		if tx.Data.AccountNonce < pool.addressTxsPool[address][index].Data.AccountNonce {
			break
		}
		index++
	}
	// update txPool
	return pool.update(index, address, tx)
}

// iterate through the txPool map, check if need to continue broadcast tx and do it
func (pool *TxPool) continueBroadcast(api *PublicEthereumAPI, currentNonce uint64, address common.Address) error {
	i := 0
	txsLen := len(pool.addressTxsPool[address])
	var err error
	for ; i < txsLen; i++ {
		if pool.addressTxsPool[address][i].Data.AccountNonce == currentNonce {
			// do broadcast
			if err = pool.broadcast(pool.addressTxsPool[address][i]); err != nil {
				// delete the tx when broadcast failed
				pool.delTxInDB(address, pool.addressTxsPool[address][i].Data.AccountNonce)
				break
			}
			// update currentNonce
			currentNonce++
		} else if pool.addressTxsPool[address][i].Data.AccountNonce < currentNonce {
			continue
		} else {
			break
		}
	}
	// i is the start index of txs that don't need to be dropped
	if err != nil {
		if !strings.Contains(err.Error(), sdkerrors.ErrMempoolIsFull.Error()) {
			// tx has err, and err is not mempoolfull, the tx should be dropped
			err = fmt.Errorf("%s, nonce %d of tx has been dropped, please send again",
				err.Error(), pool.addressTxsPool[address][i].Data.AccountNonce)
			pool.dropTxs(i+1, address)
		} else {
			err = fmt.Errorf("%s, nonce %d :", err.Error(), pool.addressTxsPool[address][i].Data.AccountNonce)
			pool.dropTxs(i, address)
		}
		api.logger.Error(err.Error())
	}

	return err
}

// drop [0:index) txs in txpool
func (pool *TxPool) dropTxs(index int, address common.Address) {
	tmp := make([]*evmtypes.MsgEthereumTx, len(pool.addressTxsPool[address][index:]), pool.cap)
	copy(tmp, pool.addressTxsPool[address][index:])
	pool.addressTxsPool[address] = tmp
}

func (pool *TxPool) broadcast(tx *evmtypes.MsgEthereumTx) error {
	txEncoder := authclient.GetTxEncoder(pool.clientCtx.Codec)
	txBytes, err := txEncoder(tx)
	if err != nil {
		return err
	}
	res, err := pool.clientCtx.BroadcastTx(txBytes)
	if res.Code != sdk.CodeOK {
		if broadcastErrors[res.Code] == nil {
			return fmt.Errorf("broadcast tx failed, code : %d", res.Code)
		} else {
			return fmt.Errorf("broadcast tx failed, err:%s", broadcastErrors[res.Code].Error())
		}
	}
	return nil
}

func (pool *TxPool) writeTxInDB(address common.Address, tx *evmtypes.MsgEthereumTx) error {
	key := []byte(address.Hex() + "|" + strconv.Itoa(int(tx.Data.AccountNonce)))

	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return err
	}

	ok, err := pool.db.Has(key)
	if err != nil {
		return err
	}
	if ok {
		return fmt.Errorf("this AccountNonce of tx has been in DB. AccountNonce[%d]", tx.Data.AccountNonce)
	}

	return pool.db.Set(key, txBytes)
}

func (pool *TxPool) delTxInDB(address common.Address, txNonce uint64) error {
	key := []byte(address.Hex() + "|" + strconv.Itoa(int(txNonce)))
	ok, err := pool.db.Has(key)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("this AccontNonce is not found in DB. AccountNonce[%d]", txNonce)
	}

	return pool.db.Delete(key)
}

func (pool *TxPool) broadcastPeriod(api *PublicEthereumAPI) {
	for {
		time.Sleep(pool.broadcastInterval)
		pool.broadcastPeriodCore(api)
	}
}
func (pool *TxPool) broadcastPeriodCore(api *PublicEthereumAPI) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	for address, _ := range pool.addressTxsPool {
		pCurrentNonce, err := api.GetTransactionCount(address, rpctypes.PendingBlockNumber)
		if err != nil {
			continue
		}
		currentNonce := uint64(*pCurrentNonce)

		pool.continueBroadcast(api, currentNonce, address)
	}
}

func (pool *TxPool) broadcastOnce(api *PublicEthereumAPI) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	for address, _ := range pool.addressTxsPool {
		pCurrentNonce, err := api.GetTransactionCount(address, rpctypes.PendingBlockNumber)
		if err != nil {
			continue
		}
		currentNonce := uint64(*pCurrentNonce)

		err = pool.continueBroadcast(api, currentNonce, address)
	}
}

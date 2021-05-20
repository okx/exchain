package eth

import (
	"fmt"
	clientcontext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/viper"
	db "github.com/tendermint/tm-db"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	FlagEnableTxPool  = "enable-tx-pool"
	TxPoolSliceMaxLen = "tx-pool-cap"
	txPoolDb          = "tx_pool"
)

var broadcastErrors = map[uint32]*sdkerrors.Error{
	sdkerrors.ErrTxInMempoolCache.ABCICode(): sdkerrors.ErrTxInMempoolCache,
	sdkerrors.ErrMempoolIsFull.ABCICode():    sdkerrors.ErrMempoolIsFull,
	sdkerrors.ErrTxTooLarge.ABCICode():       sdkerrors.ErrTxTooLarge,
}

type TxPool struct {
	addressTxsPool map[common.Address][]*evmtypes.MsgEthereumTx // All currently processable transactions
	clientCtx      clientcontext.CLIContext
	db             db.DB
	mu             sync.Mutex
}

func NewTxPool(clientCtx clientcontext.CLIContext) *TxPool {
	db, err := initDb()
	if err != nil {
		panic(err)
	}

	pool := &TxPool{
		addressTxsPool: make(map[common.Address][]*evmtypes.MsgEthereumTx),
		clientCtx:      clientCtx,
		db:             db,
	}

	itr, err := db.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		key := string(itr.Key())
		txBytes := itr.Value()

		tmp := strings.Split(key, "|")
		if len(tmp) != 2 {
			continue
		}
		address := tmp[0]
		txNonce, err := strconv.Atoi(tmp[1])
		if err != nil {
			panic(err)
		}

		tx := new(evmtypes.MsgEthereumTx)
		if err = rlp.DecodeBytes(txBytes, tx); err != nil {
			// Return nil is for when gasLimit overflows uint64
			panic(err)
		}
		if int(tx.Data.AccountNonce) != txNonce {
			panic(fmt.Errorf("nonce[%d] in key is not equal to nonce[%d] in value", tx.Data.AccountNonce, txNonce))
		}

		if err = pool.insertTx(common.HexToAddress(address), tx); err != nil {
			panic(err)
		}
	}

	return pool
}

func initDb() (db.DB, error) {
	rootDir := viper.GetString("home")
	dataDir := filepath.Join(rootDir, "data")
	return sdk.NewLevelDB(txPoolDb, dataDir)
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

	if tx.Data.AccountNonce > currentNonce+viper.GetUint64(TxPoolSliceMaxLen) {
		return fmt.Errorf("AccountNonce of tx is bigger than txPool capacity, please try later: AccountNonce[%d]", tx.Data.AccountNonce)
	}

	if err = pool.insertTx(address, tx); err != nil {
		return err
	}

	// update DB
	if err = pool.writeTxInDB(address, tx); err != nil {
		return err
	}

	pool.continueBroadcast(api, currentNonce, address)

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
		pool.addressTxsPool[address] = make([]*evmtypes.MsgEthereumTx, 0, viper.GetUint64(TxPoolSliceMaxLen))
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
func (pool *TxPool) continueBroadcast(api *PublicEthereumAPI, currentNonce uint64, address common.Address) {
	i := 0
	txsLen := len(pool.addressTxsPool[address])
	var err error
	for i < txsLen {
		if pool.addressTxsPool[address][i].Data.AccountNonce != currentNonce {
			break
		}
		// do broadcast
		if err = pool.broadcast(pool.addressTxsPool[address][i]); err != nil {
			break
		}
		// update DB
		if err = pool.delTxInDB(address, pool.addressTxsPool[address][i].Data.AccountNonce); err != nil {
			break
		}
		// update currentNonce
		currentNonce++
		i++
	}
	if err != nil {
		api.logger.Error(err.Error())
	}

	// update txPool
	if i != 0 {
		tmp := make([]*evmtypes.MsgEthereumTx, len(pool.addressTxsPool[address][i:]), viper.GetUint64(TxPoolSliceMaxLen))
		copy(tmp, pool.addressTxsPool[address][i:])
		pool.addressTxsPool[address] = tmp
	}
}

func (pool *TxPool) broadcast(tx *evmtypes.MsgEthereumTx) error {
	txEncoder := authclient.GetTxEncoder(pool.clientCtx.Codec)
	txBytes, err := txEncoder(tx)
	if err != nil {
		return err
	}
	res, err := pool.clientCtx.BroadcastTx(txBytes)
	if res.Code != sdk.CodeOK {
		return broadcastErrors[res.Code]
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
		pool.mu.Lock()
		for address, _ := range pool.addressTxsPool {
			pCurrentNonce, err := api.GetTransactionCount(address, rpctypes.PendingBlockNumber)
			if err != nil {
				continue
			}
			currentNonce := uint64(*pCurrentNonce)

			pool.continueBroadcast(api, currentNonce, address)
		}
		pool.mu.Unlock()
		time.Sleep(time.Second * 10)
	}
}

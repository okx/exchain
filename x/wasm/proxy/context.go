package proxy

import (
	"github.com/okex/exchain/x/wasm/types"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/store"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmlog "github.com/okex/exchain/libs/tendermint/libs/log"
	dbm "github.com/okex/exchain/libs/tm-db"
	evmwatcher "github.com/okex/exchain/x/evm/watcher"
)

func MakeContext(storeKey sdk.StoreKey, chainID string) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	paramsTKey := sdk.NewTransientStoreKey(params.TStoreKey)
	cms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsTKey, sdk.StoreTypeTransient, db)

	err := cms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}

	header := getHeader(chainID)

	ctx := sdk.NewContext(cms, header, true, tmlog.NewNopLogger())
	ctx.SetGasMeter(sdk.NewGasMeter(*types.DefaultWasmConfig().SimulationGasLimit))
	return ctx
}

var (
	qOnce      sync.Once
	evmQuerier *evmwatcher.Querier
)

func getHeader(chainID string) abci.Header {
	qOnce.Do(func() {
		evmQuerier = evmwatcher.NewQuerier()
	})
	timestamp := time.Now()
	latest, _ := evmQuerier.GetLatestBlockNumber()
	hash, e := evmQuerier.GetBlockHashByNumber(latest)
	if e != nil {
		hash = common.HexToHash("0x000000000000000000000000000000")
	}

	block, e := evmQuerier.GetBlockByHash(hash, false)
	if e == nil {
		timestamp = time.Unix(int64(block.Timestamp), 0)
	}

	header := abci.Header{
		ChainID: chainID,
		LastBlockId: abci.BlockID{
			Hash: hash.Bytes(),
		},
		Height: int64(latest),
		Time:   timestamp,
	}
	return header
}

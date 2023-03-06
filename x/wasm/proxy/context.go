package proxy

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	clientcontext "github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/store"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/params"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	tmlog "github.com/okx/okbchain/libs/tendermint/libs/log"
	dbm "github.com/okx/okbchain/libs/tm-db"
	evmwatcher "github.com/okx/okbchain/x/evm/watcher"
)

const (
	simulationGasLimit = 3000000
)

var clientCtx clientcontext.CLIContext

func SetCliContext(ctx clientcontext.CLIContext) {
	clientCtx = ctx
}

func MakeContext(storeKey sdk.StoreKey) sdk.Context {
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

	header := getHeader()

	ctx := sdk.NewContext(cms, header, true, tmlog.NewNopLogger())
	ctx.SetGasMeter(sdk.NewGasMeter(simulationGasLimit))
	return ctx
}

var (
	qOnce      sync.Once
	evmQuerier *evmwatcher.Querier
)

func getHeader() abci.Header {
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
		LastBlockId: abci.BlockID{
			Hash: hash.Bytes(),
		},
		Height: int64(latest),
		Time:   timestamp,
	}
	return header
}

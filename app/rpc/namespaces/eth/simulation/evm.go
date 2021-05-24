package simulation

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	abci "github.com/tendermint/tendermint/abci/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type EvmFactory struct {
	ChainId        string
	WrappedQuerier *watcher.Querier
}

func NewEvmFactory(chainId string, q *watcher.Querier) EvmFactory {
	return EvmFactory{ChainId: chainId, WrappedQuerier: q}
}

func (ef EvmFactory) BuildSimulator(qoc QueryOnChainProxy) *EvmSimulator {
	keeper := ef.makeEvmKeeper(qoc)

	if !watcher.IsWatcherEnabled() {
		return nil
	}
	timestamp := time.Now()

	latest, _ := ef.WrappedQuerier.GetLatestBlockNumber()
	hash, e := ef.WrappedQuerier.GetBlockHashByNumber(latest)
	if e != nil {
		hash = common.HexToHash("0x000000000000000000000000000000").Bytes()
	} else {
		hash = common.HexToHash(string(hash)).Bytes()
	}

	block, e := ef.WrappedQuerier.GetBlockByHash(common.BytesToHash(hash), false)

	if e == nil {
		timestamp = time.Unix(int64(block.Timestamp), 0)
	}
	req := abci.RequestBeginBlock{
		Header: abci.Header{
			ChainID: ef.ChainId,
			LastBlockId: abci.BlockID{
				Hash: hash,
			},
			Height: int64(latest),
			Time:   timestamp,
		},
		Hash: hash,
	}

	ctx := ef.makeContext(keeper, req.Header)

	keeper.BeginBlock(ctx, req)

	return &EvmSimulator{
		handler: evm.NewHandler(keeper),
		ctx:     ctx,
	}
}

type EvmSimulator struct {
	handler sdk.Handler
	ctx     sdk.Context
}

func (es *EvmSimulator) DoCall(msg evmtypes.MsgEthermint) (*sdk.SimulationResponse, error) {
	r, e := es.handler(es.ctx, msg)
	if e != nil {
		return nil, e
	}
	return &sdk.SimulationResponse{
		GasInfo: sdk.GasInfo{
			GasWanted: es.ctx.GasMeter().Limit(),
			GasUsed:   es.ctx.GasMeter().GasConsumed(),
		},
		Result: r,
	}, nil
}

func (ef EvmFactory) makeEvmKeeper(qoc QueryOnChainProxy) *evm.Keeper {
	module := evm.AppModuleBasic{}
	cdc := codec.New()
	module.RegisterCodec(cdc)
	return evm.NewSimulateKeeper(cdc, sdk.NewKVStoreKey(evm.StoreKey), NewSubspaceProxy(), NewAccountKeeperProxy(qoc), SupplyKeeperProxy{}, NewBankKeeperProxy(), NewInternalDba(qoc))
}

func (ef EvmFactory) makeContext(k *evm.Keeper, header abci.Header) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	paramsTKey := sdk.NewTransientStoreKey(params.TStoreKey)
	cms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(k.GetStoreKey(), sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsTKey, sdk.StoreTypeTransient, db)

	cms.LoadLatestVersion()

	ctx := sdk.NewContext(cms, header, true, tmlog.NewNopLogger()).WithGasMeter(sdk.NewGasMeter(evmtypes.DefaultMaxGasLimitPerTx))
	return ctx
}

package simulation

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	abci "github.com/tendermint/tendermint/abci/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type EvmFactory struct {
	ChainId string
}

func (ef EvmFactory) BuildSimulator() *EvmSimulator {
	keeper := ef.makeEvmKeeper()
	if !watcher.IsWatcherEnabled() {
		return nil
	}
	return &EvmSimulator{
		handler: evm.NewHandler(keeper),
		ctx:     ef.makeContext(keeper),
	}
}

type EvmSimulator struct {
	handler sdk.Handler
	ctx     sdk.Context
}

func (es *EvmSimulator) DoCall(msg evmtypes.MsgEthermint) (*sdk.SimulationResponse, error) {
	defer func() {
		if e := recover(); e != nil {
			panic(e)
		}
	}()
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

func (ef EvmFactory) makeEvmKeeper() *evm.Keeper {
	module := evm.AppModuleBasic{}
	cdc := codec.New()
	module.RegisterCodec(cdc)
	return evm.NewSimulateKeeper(cdc, sdk.NewKVStoreKey(evm.StoreKey), NewSubspaceProxy(), NewAccountKeeperProxy(), SupplyKeeperProxy{}, BankKeeperProxy{}, InternalDba{})
}

func (ef EvmFactory) makeContext(k *evm.Keeper) sdk.Context {
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
	ctx := sdk.NewContext(cms, abci.Header{ChainID: ef.ChainId}, true, tmlog.NewNopLogger()).WithGasMeter(sdk.NewGasMeter(evmtypes.DefaultMaxGasLimitPerTx))
	return ctx
}

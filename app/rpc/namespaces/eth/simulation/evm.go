package simulation

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type EvmFactory struct {
	chainId string
}

func (ef EvmFactory) BuildSimulator() *EvmSimulator {
	return &EvmSimulator{
		handler: evm.NewHandler(ef.makeEvmKeeper()),
		ctx:     ef.makeContext(),
	}
}

type EvmSimulator struct {
	handler sdk.Handler
	ctx     sdk.Context
}

func (es *EvmSimulator) doCall(msg evmtypes.MsgEthermint) (*sdk.SimulationResponse, error) {
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
	return evm.NewSimulateKeeper(nil, sdk.NewKVStoreKey(evm.StoreKey), SubspaceProxy{}, AccountKeeperProxy{}, SupplyKeeperProxy{}, BankKeeperProxy{}, InternalDba{})
}

func (ef EvmFactory) makeContext() sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	storeKey := sdk.NewKVStoreKey(evmtypes.StoreKey)
	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	paramsTKey := sdk.NewTransientStoreKey(params.TStoreKey)
	cms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsTKey, sdk.StoreTypeTransient, db)

	return sdk.NewContext(cms, abci.Header{ChainID: ef.chainId}, true, tmlog.NewNopLogger()).WithGasMeter(sdk.NewGasMeter(evmtypes.DefaultMaxGasLimitPerTx))
}

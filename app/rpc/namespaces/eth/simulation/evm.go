package simulation

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/okex/exchain/app/config"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/store"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmlog "github.com/okex/exchain/libs/tendermint/libs/log"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
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
		hash = common.HexToHash("0x000000000000000000000000000000")
	}

	block, e := ef.WrappedQuerier.GetBlockByHash(hash, false)

	if e == nil {
		timestamp = time.Unix(int64(block.Timestamp), 0)
	}
	req := abci.RequestBeginBlock{
		Header: abci.Header{
			ChainID: ef.ChainId,
			LastBlockId: abci.BlockID{
				Hash: hash.Bytes(),
			},
			Height: int64(latest),
			Time:   timestamp,
		},
		Hash: hash.Bytes(),
	}

	ctx := ef.makeContext(keeper, req.Header)

	keeper.BeginBlock(ctx, req)

	return &EvmSimulator{
		handler: evm.NewHandler(keeper),
		ctx:     ctx,
		keeper:  keeper,
	}
}

type EvmSimulator struct {
	handler sdk.Handler
	ctx     sdk.Context
	keeper  *evm.Keeper
}

// DoCall call simulate tx. we pass the sender by args to reduce address convert
func (es *EvmSimulator) DoCall(msg *evmtypes.MsgEthereumTx, sender string, overridesBytes []byte, estimateGas bool) (*sdk.SimulationResponse, error) {
	es.ctx.SetFrom(sender)
	if overridesBytes != nil {
		es.ctx.SetOverrideBytes(overridesBytes)
	}

	r, e := es.handler(es.ctx, msg)
	if e != nil {
		return nil, e
	}

	maxGasLimitPerTx := es.keeper.GetParams(es.ctx).MaxGasLimitPerTx
	checkedGas, err := CheckEstimatedGas(es.ctx.GasMeter().GasConsumed(), maxGasLimitPerTx)
	if err != nil {
		return nil, err
	}

	return &sdk.SimulationResponse{
		GasInfo: sdk.GasInfo{
			GasWanted: es.ctx.GasMeter().Limit(),
			GasUsed:   checkedGas,
		},
		Result: r,
	}, nil
}

func (ef EvmFactory) makeEvmKeeper(qoc QueryOnChainProxy) *evm.Keeper {
	module := evm.AppModuleBasic{}
	cdc := codec.New()
	module.RegisterCodec(cdc)
	return evm.NewSimulateKeeper(cdc, sdk.NewKVStoreKey(evm.StoreKey), sdk.NewKVStoreKey(evm.LegacyStoreKey), NewSubspaceProxy(), NewAccountKeeperProxy(qoc), SupplyKeeperProxy{}, NewBankKeeperProxy(), NewInternalDba(qoc), tmlog.NewNopLogger())
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
	cms.MountStoreWithDB(k.GetLegacyStoreKey(), sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(paramsTKey, sdk.StoreTypeTransient, db)

	cms.LoadLatestVersion()

	ctx := sdk.NewContext(cms, header, true, tmlog.NewNopLogger())
	ctx.SetGasMeter(sdk.NewGasMeter(k.GetParams(ctx).MaxGasLimitPerTx))
	return ctx
}

func CheckEstimatedGas(estimatedGas, maxGasLimitPerTx uint64) (uint64, error) {
	if estimatedGas > maxGasLimitPerTx {
		return 0, fmt.Errorf("out of gas: estimate gas is %v greater than system's max gas limit per tx %v", estimatedGas, maxGasLimitPerTx)
	}
	gasBuffer := estimatedGas / 100 * config.GetOecConfig().GetGasLimitBuffer()
	gas := estimatedGas + gasBuffer
	if gas > maxGasLimitPerTx {
		gas = maxGasLimitPerTx
	}

	return gas, nil
}

package wasm

import (
	wasmvm "github.com/CosmWasm/wasmvm"
	"github.com/okex/exchain/app/rpc/simulator"
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	"github.com/okex/exchain/x/wasm/keeper"
	"github.com/okex/exchain/x/wasm/proxy"
	"github.com/okex/exchain/x/wasm/types"
	"github.com/okex/exchain/x/wasm/watcher"
	"path/filepath"
	"sync"
)

type Simulator struct {
	handler sdk.Handler
	ctx     sdk.Context
	k       *keeper.Keeper
}

func NewWasmSimulator() simulator.Simulator {
	k := NewSimWasmKeeper()
	h := NewHandler(keeper.NewDefaultPermissionKeeper(k))
	ctx := proxy.MakeContext(k.GetStoreKey())
	return &Simulator{
		handler: h,
		k:       &k,
		ctx:     ctx,
	}
}

func (w *Simulator) Simulate(msgs []sdk.Msg, ms sdk.CacheMultiStore) (*sdk.Result, error) {
	defer func() {
		w.ctx.MoveWasmSimulateCacheToPool()
	}()
	w.ctx.SetWasmSimulateCache()
	//wasm Result has no Logs
	data := make([]byte, 0, len(msgs))
	events := sdk.EmptyEvents()

	for _, msg := range msgs {
		w.ctx.SetMultiStore(ms)
		res, err := w.handler(w.ctx, msg)
		if err != nil {
			return nil, err
		}
		data = append(data, res.Data...)
		events = events.AppendEvents(res.Events)
	}
	return &sdk.Result{
		Data:   data,
		Events: events,
	}, nil
}

func (w *Simulator) Context() *sdk.Context {
	return &w.ctx
}

func (w *Simulator) Release() {
	if !watcher.Enable() {
		return
	}
	proxy.PutBackStorePool(w.ctx.MultiStore().(sdk.CacheMultiStore))
}

func NewProxyKeeper() keeper.Keeper {
	cdc := codec.New()
	RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	interfaceReg := types2.NewInterfaceRegistry()
	RegisterInterfaces(interfaceReg)
	bank.RegisterInterface(interfaceReg)
	protoCdc := codec.NewProtoCodec(interfaceReg)

	ss := proxy.SubspaceProxy{}
	akp := proxy.NewAccountKeeperProxy()
	bkp := proxy.NewBankKeeperProxy(akp)
	pkp := proxy.PortKeeperProxy{}
	ckp := proxy.CapabilityKeeperProxy{}
	skp := proxy.SupplyKeeperProxy{}
	msgRouter := baseapp.NewMsgServiceRouter()
	msgRouter.SetInterfaceRegistry(interfaceReg)
	queryRouter := baseapp.NewGRPCQueryRouter()
	queryRouter.SetInterfaceRegistry(interfaceReg)

	k := keeper.NewSimulateKeeper(nil, codec.NewCodecProxy(protoCdc, cdc), ss, akp, bkp, nil, pkp, ckp, nil, msgRouter, queryRouter, WasmDir(), WasmConfig(), SupportedFeatures)
	types.RegisterMsgServer(msgRouter, keeper.NewMsgServerImpl(keeper.NewDefaultPermissionKeeper(k)))
	types.RegisterQueryServer(queryRouter, NewQuerier(&k))
	bank.RegisterBankMsgServer(msgRouter, bank.NewMsgServerImpl(bkp))
	bank.RegisterQueryServer(queryRouter, bank.NewBankQueryServer(bkp, skp))
	return k
}

var (
	wasmerVMCache *wasmvm.VM
	initwasmerVM  sync.Once
)

func NewWasmerVM(homeDir string, supportedFeatures string, wasmConfig types.WasmConfig) *wasmvm.VM {
	initwasmerVM.Do(func() {
		wasmer, err := wasmvm.NewVM(filepath.Join(homeDir, "wasm"), supportedFeatures, keeper.ContractMemoryLimit, wasmConfig.ContractDebugMode, wasmConfig.MemoryCacheSize)
		if err != nil {
			panic(err)
		}
		wasmerVMCache = wasmer
	})

	return wasmerVMCache
}

func NewSimWasmKeeper() keeper.Keeper {
	cdc := codec.New()
	interfaceReg := types2.NewInterfaceRegistry()
	protoCdc := codec.NewProtoCodec(interfaceReg)
	ss := proxy.SubspaceProxy{}
	akp := proxy.NewAccountKeeperProxy()
	bkp := proxy.NewBankKeeperProxy(akp)
	pkp := proxy.PortKeeperProxy{}
	ckp := proxy.CapabilityKeeperProxy{}
	msgRouter := baseapp.NewMsgServiceRouter()
	queryRouter := baseapp.NewGRPCQueryRouter()
	k := keeper.NewSimulateKeeper(NewWasmerVM(WasmDir(), SupportedFeatures, WasmConfig()), codec.NewCodecProxy(protoCdc, cdc), ss, akp, bkp, nil, pkp, ckp, nil, msgRouter, queryRouter, WasmDir(), WasmConfig(), SupportedFeatures)
	return k
}

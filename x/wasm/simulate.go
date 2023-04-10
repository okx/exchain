package wasm

import (
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
	"sync"
)

type Simulator struct {
	handler sdk.Handler
	ctx     sdk.Context
	k       *keeper.Keeper
}

func NewWasmSimulator() simulator.Simulator {
	k := NewProxyKeeper()
	h := NewHandler(keeper.NewDefaultPermissionKeeper(k))
	ctx := proxy.MakeContext(k.GetStoreKey())
	return &Simulator{
		handler: h,
		k:       &k,
		ctx:     ctx,
	}
}

func (w *Simulator) Simulate(msgs []sdk.Msg) (*sdk.Result, error) {
	//wasm Result has no Logs
	data := make([]byte, 0, len(msgs))
	events := sdk.EmptyEvents()

	for _, msg := range msgs {
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

var (
	cdc          *codec.Codec
	interfaceReg types2.InterfaceRegistry
	protoCdc     *codec.ProtoCodec
	proxyOnce    sync.Once
)

func NewProxyKeeper() keeper.Keeper {
	proxyOnce.Do(func() {
		cdc = codec.New()
		RegisterCodec(cdc)
		bank.RegisterCodec(cdc)
		interfaceReg = types2.NewInterfaceRegistry()
		RegisterInterfaces(interfaceReg)
		bank.RegisterInterface(interfaceReg)
		protoCdc = codec.NewProtoCodec(interfaceReg)
	})

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

	k := keeper.NewSimulateKeeper(codec.NewCodecProxy(protoCdc, cdc), getStoreKey(), ss, akp, bkp, nil, pkp, ckp, nil, msgRouter, queryRouter, WasmDir(), WasmConfig(), SupportedFeatures)
	types.RegisterMsgServer(msgRouter, keeper.NewMsgServerImpl(keeper.NewDefaultPermissionKeeper(k)))
	types.RegisterQueryServer(queryRouter, NewQuerier(&k))
	bank.RegisterBankMsgServer(msgRouter, bank.NewMsgServerImpl(bkp))
	bank.RegisterQueryServer(queryRouter, bank.NewBankQueryServer(bkp, skp))

	return k
}

var (
	storeKeyOnce sync.Once
	gStoreKey    sdk.StoreKey
)

func getStoreKey() sdk.StoreKey {
	storeKeyOnce.Do(
		func() {
			gStoreKey = sdk.NewKVStoreKey(StoreKey)
		},
	)

	return gStoreKey
}

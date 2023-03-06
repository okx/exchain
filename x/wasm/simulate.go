package wasm

import (
	"github.com/okx/okbchain/app/rpc/simulator"
	"github.com/okx/okbchain/libs/cosmos-sdk/baseapp"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	types2 "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/bank"
	"github.com/okx/okbchain/x/wasm/keeper"
	"github.com/okx/okbchain/x/wasm/proxy"
	"github.com/okx/okbchain/x/wasm/types"
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
	paramKP := proxy.ParamsKeeperProxy{}
	pkp := proxy.PortKeeperProxy{}
	ckp := proxy.CapabilityKeeperProxy{}
	skp := proxy.SupplyKeeperProxy{}
	msgRouter := baseapp.NewMsgServiceRouter()
	msgRouter.SetInterfaceRegistry(interfaceReg)
	queryRouter := baseapp.NewGRPCQueryRouter()
	queryRouter.SetInterfaceRegistry(interfaceReg)

	k := keeper.NewSimulateKeeper(codec.NewCodecProxy(protoCdc, cdc), sdk.NewKVStoreKey(StoreKey), ss, akp, bkp, paramKP, nil, pkp, ckp, nil, msgRouter, queryRouter, WasmDir(), WasmConfig(), SupportedFeatures)
	types.RegisterMsgServer(msgRouter, keeper.NewMsgServerImpl(keeper.NewDefaultPermissionKeeper(k)))
	types.RegisterQueryServer(queryRouter, NewQuerier(&k))
	bank.RegisterBankMsgServer(msgRouter, bank.NewMsgServerImpl(bkp))
	bank.RegisterQueryServer(queryRouter, bank.NewBankQueryServer(bkp, skp))
	return k
}

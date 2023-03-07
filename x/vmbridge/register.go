package vmbridge

import (
	"github.com/okx/exchain/libs/cosmos-sdk/codec"
	"github.com/okx/exchain/libs/cosmos-sdk/types/module"
	"github.com/okx/exchain/x/vmbridge/keeper"
	"github.com/okx/exchain/x/wasm"
)

func RegisterServices(cfg module.Configurator, keeper keeper.Keeper) {
	RegisterMsgServer(cfg.MsgServer(), NewMsgServerImpl(keeper))
}

func GetWasmOpts(cdc *codec.ProtoCodec) wasm.Option {
	return wasm.WithMessageEncoders(RegisterSendToEvmEncoder(cdc))
}

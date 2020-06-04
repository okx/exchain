package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	// "github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&MsgStoreCode{}, "wasm/store-code", nil)
	cdc.RegisterConcrete(&MsgInstantiateContract{}, "wasm/instantiate", nil)
	cdc.RegisterConcrete(&MsgExecuteContract{}, "wasm/execute", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}

package types

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
	"github.com/tendermint/go-amino"
)

const (
	// MudulleAccountName is the amino encoding name for ModuleAccount
	MuduleAccountName = "cosmos-sdk/ModuleAccount"
)

// RegisterCodec registers the account types and interface
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.ModuleAccountI)(nil), nil)
	cdc.RegisterInterface((*exported.SupplyI)(nil), nil)
	cdc.RegisterConcrete(&ModuleAccount{}, MuduleAccountName, nil)
	cdc.RegisterConcreteUnmarshaller(MuduleAccountName, func(cdc *amino.Codec, data []byte) (interface{}, int, error) {
		var acc ModuleAccount
		err := acc.UnmarshalFromAmino(cdc, data)
		if err != nil {
			return nil, 0, err
		}
		return &acc, len(data), nil
	})
	cdc.RegisterConcreteMarshaller(MuduleAccountName, func(cdc *amino.Codec, v interface{}) ([]byte, error) {
		if m, ok := v.(*ModuleAccount); ok {
			return m.MarshalToAmino(cdc)
		} else if m, ok := v.(ModuleAccount); ok {
			return m.MarshalToAmino(cdc)
		} else {
			return nil, fmt.Errorf("%T is not a ModuleAccount", v)
		}
	})

	cdc.RegisterConcrete(&Supply{}, "cosmos-sdk/Supply", nil)
}

// ModuleCdc generic sealed codec to be used throughout module
var ModuleCdc *codec.Codec

func init() {
	cdc := codec.New()
	RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ModuleCdc = cdc.Seal()
}

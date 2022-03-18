package types

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/tendermint/go-amino"
)

const (
	// EthAccountName is the amino encoding name for EthAccount
	EthAccountName = "okexchain/EthAccount"
)

// RegisterCodec registers the account interfaces and concrete types on the
// provided Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&EthAccount{}, EthAccountName, nil)

	cdc.RegisterConcreteUnmarshaller(EthAccountName, func(cdc *amino.Codec, data []byte) (interface{}, int, error) {
		var acc EthAccount
		err := acc.UnmarshalFromAmino(cdc, data)
		if err != nil {
			return nil, 0, err
		}
		return &acc, len(data), nil
	})
	cdc.RegisterConcreteMarshaller(EthAccountName, func(cdc *amino.Codec, v interface{}) ([]byte, error) {
		if acc, ok := v.(*EthAccount); ok {
			return acc.MarshalToAmino(cdc)
		} else if acc, ok := v.(EthAccount); ok {
			return acc.MarshalToAmino(cdc)
		} else {
			return nil, fmt.Errorf("%T is not an EthAccount", v)
		}
	})
	cdc.EnableBufferMarshaler(EthAccount{})
}

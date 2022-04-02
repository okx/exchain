package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/tendermint/go-amino"
)

// ModuleCdc auth module wide codec
var ModuleCdc = codec.New()

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.GenesisAccount)(nil), nil)
	cdc.RegisterInterface((*exported.Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "cosmos-sdk/Account", nil)
	cdc.RegisterConcrete(StdTx{}, "cosmos-sdk/StdTx", nil)

	cdc.RegisterConcreteUnmarshaller("cosmos-sdk/StdTx", func(c *amino.Codec, bytes []byte) (interface{}, int, error) {
		var tx StdTx
		err := tx.UnmarshalFromAmino(c, bytes)
		if err != nil {
			return nil, 0, err
		}
		return tx, len(bytes), nil
	})
	cdc.RegisterConcrete(&RawIBCViewTx{}, "cosmos-sdk/IbcViewTx", nil)
	cdc.RegisterConcrete(&IbcViewMsg{}, "cosmos-sdk/IbcViewMsg", nil)
}

// RegisterAccountTypeCodec registers an external account type defined in
// another module for the internal ModuleCdc.
func RegisterAccountTypeCodec(o interface{}, name string) {
	ModuleCdc.RegisterConcrete(o, name, nil)
}

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
}

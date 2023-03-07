package codec

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/crypto/keys"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/module"
	cosmoscryptocodec "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/ibc-tx"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/vesting"

	cryptocodec "github.com/okx/okbchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okx/okbchain/app/types"
)

// MakeCodec registers the necessary types and interfaces for an sdk.App. This
// codec is provided to all the modules the application depends on.
//
// NOTE: This codec will be deprecated in favor of AppCodec once all modules are
// migrated to protobuf.
func MakeCodec(bm module.BasicManager) *codec.Codec {
	cdc := codec.New()

	bm.RegisterCodec(cdc)
	vesting.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	cryptocodec.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	ethermint.RegisterCodec(cdc)
	keys.RegisterCodec(cdc) // temporary. Used to register keyring.Info

	return cdc
}

func MakeIBC(bm module.BasicManager) types.InterfaceRegistry {
	interfaceReg := types.NewInterfaceRegistry()
	bm.RegisterInterfaces(interfaceReg)
	cosmoscryptocodec.PubKeyRegisterInterfaces(interfaceReg)
	return interfaceReg
}

func MakeCodecSuit(bm module.BasicManager) (*codec.CodecProxy, types.InterfaceRegistry) {
	aminoCodec := MakeCodec(bm)
	interfaceReg := MakeIBC(bm)
	protoCdc := codec.NewProtoCodec(interfaceReg)
	return codec.NewCodecProxy(protoCdc, aminoCodec), interfaceReg
}

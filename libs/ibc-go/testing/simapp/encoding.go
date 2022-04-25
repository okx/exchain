package simapp

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	cosmoscryptocodec "github.com/okex/exchain/libs/cosmos-sdk/crypto/ibc-codec"
	simappparams "github.com/okex/exchain/libs/ibc-go/testing/simapp/params"
)

// MakeTestEncodingConfig creates an EncodingConfig for testing. This function
// should be used only in tests or when creating a new app instance (NewApp*()).
// App user shouldn't create new codecs - use the app.AppCodec instead.
// [DEPRECATED]
func MakeTestEncodingConfig() simappparams.EncodingConfig {
	encodingConfig := simappparams.MakeTestEncodingConfig()
	//std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	cosmoscryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	//ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	interfaceReg := types.NewInterfaceRegistry()
	ModuleBasics.RegisterInterfaces(interfaceReg)
	return encodingConfig
}

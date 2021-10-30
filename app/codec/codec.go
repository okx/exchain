package codec

import (
	"github.com/okex/exchain/dependence/cosmos-sdk/codec"
	"github.com/okex/exchain/dependence/cosmos-sdk/crypto/keys"
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/dependence/cosmos-sdk/types/module"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/auth/vesting"

	cryptocodec "github.com/okex/exchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/exchain/app/types"
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

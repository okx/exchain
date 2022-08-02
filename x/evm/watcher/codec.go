package watcher

import (
	cryptocodec "github.com/okex/exchain/app/crypto/ethsecp256k1"
	app "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

var watchCdc *codec.Codec

func init() {
	watchCdc = codec.New()
	cryptocodec.RegisterCodec(watchCdc)
	codec.RegisterCrypto(watchCdc)
	watchCdc.RegisterInterface((*exported.Account)(nil), nil)
	app.RegisterCodec(watchCdc)
}

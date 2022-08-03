package watcher

import (
	cryptocodec "github.com/okex/exchain/app/crypto/ethsecp256k1"
	app "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

var WatchCdc *codec.Codec

func init() {
	WatchCdc = codec.New()
	cryptocodec.RegisterCodec(WatchCdc)
	codec.RegisterCrypto(WatchCdc)
	WatchCdc.RegisterInterface((*exported.Account)(nil), nil)
	app.RegisterCodec(WatchCdc)
}

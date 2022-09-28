package coretypes

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	amino "github.com/tendermint/go-amino"

	"github.com/okex/exchain/libs/tendermint/types"
)

func RegisterAmino(cdc *amino.Codec) {
	types.RegisterEventDatas(cdc)
	types.RegisterBlockAmino(cdc)
}

func RegisterCM40Codec(cdc *codec.CodecProxy) {
	types.RegisterCM40EventDatas(cdc.GetCdc())
	types.RegisterBlockAmino(cdc.GetCdc())
}

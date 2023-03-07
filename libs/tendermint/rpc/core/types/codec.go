package coretypes

import (
	amino "github.com/tendermint/go-amino"

	"github.com/okx/exchain/libs/tendermint/types"
)

func RegisterAmino(cdc *amino.Codec) {
	types.RegisterEventDatas(cdc)
	types.RegisterBlockAmino(cdc)
}

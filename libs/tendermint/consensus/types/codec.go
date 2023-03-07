package types

import (
	amino "github.com/tendermint/go-amino"

	"github.com/okx/exchain/libs/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	types.RegisterBlockAmino(cdc)
}

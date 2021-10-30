package store

import (
	amino "github.com/tendermint/go-amino"

	"github.com/okex/exchain/dependence/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	types.RegisterBlockAmino(cdc)
}

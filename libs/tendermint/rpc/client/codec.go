package client

import (
	amino "github.com/tendermint/go-amino"

	"github.com/okex/exchain/libs/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	types.RegisterEvidences(cdc)
}

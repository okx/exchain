package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

const (
	// EthAccountName is the amino encoding name for EthAccount
	EthAccountName = "okexchain/EthAccount"
)

// RegisterCodec registers the account interfaces and concrete types on the
// provided Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&EthAccount{}, EthAccountName, nil)
	exported.RegisterConcreteAccountInfo(uint(exported.EthAcc), func() exported.MptAccount {
		return &EthAccount{}
	})
}

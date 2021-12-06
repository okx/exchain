package types

import (
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/wrap"
)

const (
	// EthAccountName is the amino encoding name for EthAccount
	EthAccountName = "okexchain/EthAccount"
)

// RegisterCodec registers the account interfaces and concrete types on the
// provided Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&EthAccount{}, EthAccountName, nil)
	wrap.RegisterConcreteAccountInfo(uint(exported.EthAcc), &EthAccount{})
}
